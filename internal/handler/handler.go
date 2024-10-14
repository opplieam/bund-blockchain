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

// Genesis returns the genesis information.
func (h *Handler) Genesis(c echo.Context) error {
	return c.JSON(http.StatusOK, h.State.Genesis())
}

// Accounts returns the current balances for all users.
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

// Mempool returns the set of uncommitted transactions.
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

// SubmitWalletTransaction adds new transactions to the mempool.
func (h *Handler) SubmitWalletTransaction(c echo.Context) error {
	var signedTx database.SignedTx
	if err := c.Bind(&signedTx); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	h.Log.Info("add trans", "sig|nonce", signedTx, "from", signedTx.FromID, "to", signedTx.ToID, "value", signedTx.Value, "tip", signedTx.Tip)

	// Ask the state package to add this transaction to the mempool. Only the
	// checks are the transaction signature and the recipient account format.
	// It's up to the wallet to make sure the account has a proper balance and
	// nonce. Fees will be taken if this transaction is mined into a block.
	if err := h.State.UpsertWalletTransaction(signedTx); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: "transactions added to mempool",
	}
	return c.JSON(http.StatusOK, response)
}
