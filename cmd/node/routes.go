package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opplieam/bund-blockchain/internal/blockchain/state"
	"github.com/opplieam/bund-blockchain/internal/handler"
	slogecho "github.com/samber/slog-echo"
)

func setupRoutes(e *echo.Echo, log *slog.Logger, state *state.State) {
	e.Use(slogecho.New(log))
	e.Use(middleware.Recover())

	h := handler.New(log, state)
	e.GET("/genesis/list", h.Genesis)
}
