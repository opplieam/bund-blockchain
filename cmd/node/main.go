package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/opplieam/bund-blockchain/internal/blockchain/database"
	"github.com/opplieam/bund-blockchain/internal/blockchain/genesis"
	"github.com/opplieam/bund-blockchain/internal/blockchain/peer"
	"github.com/opplieam/bund-blockchain/internal/blockchain/state"
	"github.com/opplieam/bund-blockchain/internal/blockchain/storage/disk"
	"github.com/opplieam/bund-blockchain/internal/blockchain/worker"
	"github.com/opplieam/bund-blockchain/internal/nameservice"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("service", "NODE")
	if err := run(logger); err != nil {
		logger.Error("Run Node", "error", err)
	}
}

func run(log *slog.Logger) error {
	log.Info("start up", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// Load Env
	args := os.Args[1]
	if args == "" {
		return fmt.Errorf("missing required argument (eg. miner1)")
	}
	pathToLoad := fmt.Sprintf("conf/%s.env", args)
	if err := godotenv.Load(pathToLoad); err != nil {
		return fmt.Errorf("could not load .env file: %w", err)
	}

	// Load config
	cfg := NewConfig()

	// =========================================================================
	// Name Service Support. ONLY FOR DEV

	// The nameservice package provides name resolution for account addresses.
	// The names come from the file names in the conf/accounts folder.
	ns, err := nameservice.New(cfg.NameService.Folder)
	if err != nil {
		return fmt.Errorf("unable to load account name service: %w", err)
	}

	// Logging the accounts for documentation in the logs.
	for account, name := range ns.Copy() {
		log.Info("startup", "status", "nameservice", "name", name, "account", account)
	}

	// =========================================================================
	// Need to load the private key file for the configured beneficiary so the
	// account can get credited with fees and tips.
	path := fmt.Sprintf("%s/%s.ecdsa", cfg.NameService.Folder, cfg.State.Beneficiary)
	privateKey, err := crypto.LoadECDSA(path)
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	// A peer set is a collection of known nodes in the network so transactions
	// and blocks can be shared.
	peerSet := peer.NewPeerSet()
	for _, host := range cfg.State.OriginPeers {
		peerSet.Add(peer.New(host))
	}
	peerSet.Add(peer.New(cfg.Web.PrivateAddr))

	ev := func(v string, args ...any) {
		s := fmt.Sprintf(v, args...)
		log.Info(s, "trace_id", "00000000-0000-0000-0000-000000000000")
	}

	// Construct the use of disk storage
	storage, err := disk.New(cfg.State.DBPath)
	if err != nil {
		return err
	}

	// Load the genesis file for blockchain settings and origin balances.
	genesisInfo, err := genesis.Load()
	if err != nil {
		return err
	}

	// The state value represents the blockchain node and manages the blockchain
	// database and provides an API for application support.
	stateM, err := state.New(state.Config{
		BeneficiaryID:  database.PublicKeyToAccountID(privateKey.PublicKey),
		Host:           cfg.Web.PrivateAddr,
		Storage:        storage,
		Genesis:        genesisInfo,
		SelectStrategy: cfg.State.SelectStrategy,
		KnownPeers:     peerSet,
		EvHandler:      ev,
	})
	if err != nil {
		return err
	}
	defer stateM.Shutdown()

	// The worker package implements the different workflows such as mining,
	// transaction peer sharing, and peer updates. The worker will register
	// itself with the state.
	worker.Run(stateM, ev)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// ===========================================================================================
	e := echo.New()
	setupRoutes(e, log, stateM, ns)

	publicSrv := &http.Server{
		Addr:         cfg.Web.Addr,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      e,
	}
	go func() {
		log.Info("http public service start", "addr", cfg.Web.Addr)
		serverErrors <- publicSrv.ListenAndServe()
	}()

	// ===========================================================================================
	pe := echo.New()
	setupPrivateRoutes(pe, log, stateM, ns)

	privateSrv := &http.Server{
		Addr:         cfg.Web.PrivateAddr,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      pe,
	}
	go func() {
		log.Info("http private service start", "addr", cfg.Web.PrivateAddr)
		serverErrors <- privateSrv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancelPub := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancelPub()
		if err := publicSrv.Shutdown(ctx); err != nil {
			publicSrv.Close()
		}

		ctx, cancelPri := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancelPri()
		if err := privateSrv.Shutdown(ctx); err != nil {
			privateSrv.Close()
		}

	}

	return nil
}
