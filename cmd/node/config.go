package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/opplieam/bund-blockchain/internal/utils/getenv"
)

type Config struct {
	Web         WebConfig
	NameService NameService
	State       State
}

type NameService struct {
	Folder string
}

type State struct {
	Beneficiary    string
	DBPath         string
	SelectStrategy string
	OriginPeers    []string
	Consensus      string
}

type WebConfig struct {
	Addr            string
	PrivateAddr     string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func NewConfig() Config {
	writeTimeout, _ := strconv.Atoi(getenv.GetEnv("WEB_WRITE_TIMEOUT", "10"))
	readTimeout, _ := strconv.Atoi(getenv.GetEnv("WEB_READ_TIMEOUT", "5"))
	idleTimeout, _ := strconv.Atoi(getenv.GetEnv("WEB_IDLE_TIMEOUT", "120"))
	shutDownTimeout, _ := strconv.Atoi(getenv.GetEnv("WEB_SHUTDOWN_TIMEOUT", "20"))

	originPeers := strings.Split(getenv.GetEnv("ORIGIN_PEERS", "0.0.0.0:3030"), ",")

	return Config{
		Web: WebConfig{
			Addr:            getenv.GetEnv("WEB_ADDR", "0.0.0.0:3000"),
			PrivateAddr:     getenv.GetEnv("WEB_PRIVATE_ADDR", "0.0.0.0:3030"),
			WriteTimeout:    time.Duration(writeTimeout) * time.Second,
			ReadTimeout:     time.Duration(readTimeout) * time.Second,
			IdleTimeout:     time.Duration(idleTimeout) * time.Second,
			ShutdownTimeout: time.Duration(shutDownTimeout) * time.Second,
		},
		NameService: NameService{
			Folder: getenv.GetEnv("PRIVATE_KEY_PATH", "conf/accounts"),
		},
		State: State{
			Beneficiary:    getenv.GetEnv("BENEFICIARY", "miner1"),
			DBPath:         getenv.GetEnv("DB_PATH", "data/miner1/"),
			SelectStrategy: getenv.GetEnv("SELECT_STRATEGY", "Tip"),
			OriginPeers:    originPeers,
			Consensus:      getenv.GetEnv("CONSENSUS", "POW"),
		},
	}
}
