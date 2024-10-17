package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/opplieam/bund-blockchain/internal/blockchain/database"
	"github.com/opplieam/bund-blockchain/internal/blockchain/peer"
	"github.com/opplieam/bund-blockchain/internal/blockchain/state"
	"github.com/opplieam/bund-blockchain/internal/nameservice"
)

type Handler struct {
	Log   *slog.Logger
	State *state.State
	NS    *nameservice.NameService
}

func New(logger *slog.Logger, state *state.State, ns *nameservice.NameService) *Handler {
	return &Handler{
		Log:   logger,
		State: state,
		NS:    ns,
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

	resp := make([]act, 0, len(accounts))
	for account, info := range accounts {
		act := act{
			Account: account,
			Name:    h.NS.Lookup(account),
			Balance: info.Balance,
			Nonce:   info.Nonce,
		}
		resp = append(resp, act)
	}

	ai := actInfo{
		LastestBlock: h.State.LatestBlock().Hash(),
		Uncommitted:  len(h.State.Mempool()),
		Accounts:     resp,
	}

	return c.JSON(200, ai)
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
			FromName:    h.NS.Lookup(tran.FromID),
			To:          tran.ToID,
			ToName:      h.NS.Lookup(tran.ToID),
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

func (h *Handler) Cancel(c echo.Context) error {
	h.State.Worker.SignalCancelMining()
	return c.String(http.StatusOK, "cancelled")
}

func (h *Handler) Status(c echo.Context) error {
	latestBlock := h.State.LatestBlock()

	status := peer.PeerStatus{
		LatestBlockHash:   latestBlock.Hash(),
		LatestBlockNumber: latestBlock.Header.Number,
		KnownPeers:        h.State.KnownExternalPeers(),
	}
	return c.JSON(http.StatusOK, status)
}

func (h *Handler) PrivateMempool(c echo.Context) error {
	txs := h.State.Mempool()
	return c.JSON(http.StatusOK, txs)
}

func (h *Handler) BlocksByNumber(c echo.Context) error {
	fromStr := c.Param("from")
	if fromStr == "latest" || fromStr == "" {
		fromStr = fmt.Sprintf("%d", state.QueryLatest)
	}

	toStr := c.Param("to")
	if toStr == "latest" || toStr == "" {
		toStr = fmt.Sprintf("%d", state.QueryLatest)
	}

	from, err := strconv.ParseUint(fromStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	to, err := strconv.ParseUint(toStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if from > to {
		return c.String(http.StatusBadRequest, "from must be less than to")
	}

	blocks := h.State.QueryBlocksByNumber(from, to)
	if len(blocks) == 0 {
		return c.JSON(http.StatusNoContent, nil)
	}

	blockData := make([]database.BlockData, len(blocks))
	for i, block := range blocks {
		blockData[i] = database.NewBlockData(block)
	}

	return c.JSON(http.StatusOK, blockData)
}

func (h *Handler) SubmitPeer(c echo.Context) error {
	var peer peer.Peer
	if err := c.Bind(&peer); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if !h.State.AddKnownPeer(peer) {
		h.Log.Info("adding peer", "host", peer.Host)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) SubmitNodeTransaction(c echo.Context) error {
	// Decode the JSON in the post call into a block transaction.
	var tx database.BlockTx
	if err := c.Bind(&tx); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Ask the state package to add this transaction to the mempool and perform
	// any other business logic.
	h.Log.Info("add tran", "sig:nonce", tx, "from", tx.FromID, "to", tx.ToID, "value", tx.Value, "tip", tx.Tip)
	if err := h.State.UpsertNodeTransaction(tx); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	resp := struct {
		Status string `json:"status"`
	}{
		Status: "transactions added to mempool",
	}
	return c.JSON(http.StatusOK, resp)
}

// ProposeBlock takes a block received from a peer, validates it and
// if that passes, adds the block to the local blockchain.
func (h *Handler) ProposeBlock(c echo.Context) error {
	// Decode the JSON in the post call into a file system block.
	var blockData database.BlockData
	if err := c.Bind(&blockData); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Convert the block data into a block. This action will create a merkle
	// tree for the set of transactions required for blockchain operations.
	block, err := database.ToBlock(blockData)
	if err != nil {
		return fmt.Errorf("unable to decode block: %w", err)
	}

	// Ask the state package to validate the proposed block. If the block
	// passes validation, it will be added to the blockchain database.
	if err := h.State.ProcessProposedBlock(block); err != nil {
		//if errors.Is(err, database.ErrChainForked) {
		//	h.State.Reorganize()
		//}

		return c.String(http.StatusBadRequest, "block not accepted")
	}

	resp := struct {
		Status string `json:"status"`
	}{
		Status: "accepted",
	}

	return c.JSON(http.StatusOK, resp)
}
