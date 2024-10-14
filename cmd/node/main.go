package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"

	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/opplieam/bund-blockchain/internal/blockchain/database"
	"github.com/opplieam/bund-blockchain/internal/blockchain/genesis"
	"github.com/opplieam/bund-blockchain/internal/blockchain/state"
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

	// Load config
	cfg := NewConfig()

	// Need to load the private key file for the configured beneficiary so the
	// account can get credited with fees and tips.
	path := fmt.Sprintf("%s/%s.ecdsa", cfg.NameService.Folder, cfg.State.Beneficiary)
	privateKey, err := crypto.LoadECDSA(path)
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	ev := func(v string, args ...any) {
		s := fmt.Sprintf(v, args...)
		log.Info(s, "traceid", "00000000-0000-0000-0000-000000000000")
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
		Genesis:        genesisInfo,
		SelectStrategy: cfg.State.SelectStrategy,
		EvHandler:      ev,
	})
	if err != nil {
		return err
	}
	defer stateM.Shutdown()

	// ===========================================================================================
	log.Info("http service start", "addr", cfg.Web.Addr)
	e := echo.New()
	setupRoutes(e, log, stateM)

	srv := &http.Server{
		Addr:         cfg.Web.Addr,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      e,
	}
	srv.ListenAndServe()

	return nil
}
