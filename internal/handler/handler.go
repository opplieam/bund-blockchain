package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opplieam/bund-blockchain/internal/blockchain/database"
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

func (h *Handler) Accounts(c echo.Context) error {
	accountStr := c.Param("account")

	var accounts map[database.AccountID]database.Account
	switch accountStr {
	case "":
		accounts = h.State.Accounts()
	default:
		accountID, err := database.ToAccountID(accountStr)
		if err != nil {
			return err
		}
		account, err := h.State.QueryAccount(accountID)
		if err != nil {
			return err
		}
		accounts = map[database.AccountID]database.Account{accountID: account}
	}

	return c.JSON(200, accounts)
}

func (h *Handler) Mempool(c echo.Context) error {
	accountStr := c.Param("account")
	mempool := h.State.Mempool()

	txResult := []tx{} // If there is no mempool return empty slice
	for _, tran := range mempool {
		if accountStr != "" && ((accountStr != string(tran.FromID)) && (accountStr != string(tran.ToID))) {
			continue
		}

		txResult = append(txResult, tx{
			FromAccount: tran.FromID,
			To:          tran.ToID,
			ChainID:     tran.ChainID,
			Nonce:       tran.Nonce,
			Value:       tran.Value,
			Tip:         tran.Tip,
			Data:        tran.Data,
			TimeStamp:   tran.TimeStamp,
			GasPrice:    tran.GasPrice,
			GasUnits:    tran.GasUnits,
			Sig:         tran.SignatureString(),
		})
	}

	return c.JSON(http.StatusOK, txResult)
}
