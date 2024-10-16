package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opplieam/bund-blockchain/internal/blockchain/state"
	"github.com/opplieam/bund-blockchain/internal/handler"
	"github.com/opplieam/bund-blockchain/internal/nameservice"
	slogecho "github.com/samber/slog-echo"
)

func setupRoutes(e *echo.Echo, log *slog.Logger, state *state.State, ns *nameservice.NameService) {
	e.Use(slogecho.New(log))
	e.Use(middleware.Recover())

	h := handler.New(log, state, ns)

	e.GET("/genesis/list", h.Genesis)
	e.GET("/accounts/list", h.Accounts)
	e.GET("/accounts/list/:account", h.Accounts)
	e.GET("/tx/uncommitted/list", h.Mempool)
	e.GET("/tx/uncommitted/list/:account", h.Mempool)
	e.POST("/tx/submit", h.SubmitWalletTransaction)
	//e.POST("/tx/proof/:block")

}

func setupPrivateRoutes(e *echo.Echo, log *slog.Logger, state *state.State, ns *nameservice.NameService) {
	e.Use(slogecho.New(log))
	e.Use(middleware.Recover())

	h := handler.New(log, state, ns)
	e.POST("/node/peers", h.SubmitPeer)
	e.GET("/node/status", h.Status)
	e.GET("/node/tx/list", h.PrivateMempool)
	e.POST("/node/tx/submit", h.SubmitNodeTransaction)
	e.POST("/node/block/propose", h.ProposeBlock)
	e.GET("/node/block/list/:from/:to", h.BlocksByNumber)
}
