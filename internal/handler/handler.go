package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opplieam/bund-blockchain/internal/blockchain/state"
)

type Handler struct {
	Log   *slog.Logger
	State *state.State
}

func New(logger *slog.Logger, state *state.State) *Handler {
	return &Handler{
		Log:   logger,
		State: state,
	}
}

func (h *Handler) Genesis(c echo.Context) error {

	return c.JSON(http.StatusOK, h.State.Genesis())
}
