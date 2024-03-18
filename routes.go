package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/samber/lo"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

var (
	ErrZeroAccounts  = errors.New("no accounts were specified")
	ErrZeroTransfers = errors.New("no transfers were specified")
)

type server struct {
	TB tb.Client
}

func (s *server) GetID(c *gin.Context) {
	c.JSON(http.StatusOK, GetIDResponse{ID: types.ID().String()})
}

func (s *server) CreateAccounts(c *gin.Context) {
	var req *CreateAccountsRequest
	if ok := bindJSON(c, req); !ok {
		return
	}
	if len(req.Accounts) == 0 {
		abort(c, http.StatusBadRequest, ErrZeroAccounts)
		return
	}
	accounts := []types.Account{}
	for _, inAccount := range req.Accounts {
		// id
		id, err := hexStringToUint128(inAccount.ID)
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		// userData64
		ud128, ud64, ud32, err := inAccount.UserData.ToUint()
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}

		flags := types.AccountFlags{}
		if inAccount.Flags != nil {
			flags.Linked = inAccount.Flags.Linked
			flags.DebitsMustNotExceedCredits = inAccount.Flags.DebitsMustNotExceedCredits
			flags.CreditsMustNotExceedDebits = inAccount.Flags.CreditsMustNotExceedDebits
			flags.History = inAccount.Flags.History
		}

		ud128, ud64, ud32, err = inAccount.UserData.ToUint()
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		accounts = append(accounts, types.Account{
			ID:             *id,
			DebitsPending:  types.ToUint128(uint64(inAccount.DebitsPending)),
			DebitsPosted:   types.ToUint128(uint64(inAccount.DebitsPosted)),
			CreditsPending: types.ToUint128(uint64(inAccount.CreditsPending)),
			CreditsPosted:  types.ToUint128(uint64(inAccount.CreditsPosted)),
			UserData128:    ud128,
			UserData64:     ud64,
			UserData32:     ud32,
			Ledger:         uint32(inAccount.Ledger),
			Code:           uint16(inAccount.Code),
			Flags:          flags.ToUint16(),
		})
	}

	resp, err := s.TB.CreateAccounts(accounts)
	if err != nil {
		abort(c, http.StatusInternalServerError, err)
		return
	}

	resArr := []string{}
	for _, r := range resp {
		resArr = append(resArr, r.Result.String())
	}
	c.JSON(http.StatusOK, CreateAccountsResponse{
		Results: resArr,
	})
}

func (s *server) CreateTransfers(c *gin.Context) {
	var req *CreateTransfersRequest
	if ok := bindJSON(c, req); !ok {
		return
	}
	if len(req.Transfers) == 0 {
		abort(c, http.StatusInternalServerError, ErrZeroTransfers)
		return
	}
	transfers := []types.Transfer{}
	for _, inTransfer := range req.Transfers {
		id, err := hexStringToUint128(inTransfer.ID)
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		flags := types.TransferFlags{}
		if inTransfer.TransferFlags != nil {
			inTransfer.TransferFlags.Linked = inTransfer.TransferFlags.Linked
			inTransfer.TransferFlags.Pending = inTransfer.TransferFlags.Pending
			inTransfer.TransferFlags.PostPendingTransfer = inTransfer.TransferFlags.PostPendingTransfer
			inTransfer.TransferFlags.VoidPendingTransfer = inTransfer.TransferFlags.VoidPendingTransfer
			inTransfer.TransferFlags.BalancingDebit = inTransfer.TransferFlags.BalancingDebit
			inTransfer.TransferFlags.BalancingCredit = inTransfer.TransferFlags.BalancingCredit
		}

		debitAccountID, err := hexStringToUint128(inTransfer.DebitAccountID)
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		creditAccountID, err := hexStringToUint128(inTransfer.CreditAccountID)
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		pendingID, err := hexStringToUint128(lo.FromPtrOr(inTransfer.PendingID, ""))
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		ud128, ud64, ud32, err := inTransfer.UserData.ToUint()
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		transfers = append(transfers, types.Transfer{
			ID:              *id,
			DebitAccountID:  *debitAccountID,
			CreditAccountID: *creditAccountID,
			Amount:          types.ToUint128(uint64(inTransfer.Amount)),
			PendingID:       *pendingID,
			UserData128:     ud128,
			UserData64:      ud64,
			UserData32:      ud32,
			Timeout:         0,
			Ledger:          uint32(inTransfer.Ledger),
			Code:            uint16(inTransfer.Ledger),
			Flags:           flags.ToUint16(),
			Timestamp:       0,
		})
	}

	resp, err := s.TB.CreateTransfers(transfers)
	if err != nil {
		abort(c, http.StatusInternalServerError, err)
		return
	}

	resArr := []string{}
	for _, r := range resp {
		resArr = append(resArr, r.Result.String())
	}
	c.JSON(http.StatusOK, CreateTransfersResponse{
		Results: resArr,
	})
}

func (s *server) LookupAccounts(c *gin.Context) {
	var req *LookupAccountsRequest
	if ok := bindJSON(c, req); !ok {
		return
	}
	if len(req.AccountIds) == 0 {
		abort(c, http.StatusInternalServerError, ErrZeroAccounts)
		return
	}
	ids := []types.Uint128{}
	for _, inID := range req.AccountIds {
		id, err := hexStringToUint128(inID)
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		ids = append(ids, *id)
	}

	res, err := s.TB.LookupAccounts(ids)
	if err != nil {
		abort(c, http.StatusInternalServerError, err)
		return
	}

	pAccounts := lo.Map(res, func(a types.Account, _ int) Account {
		return *AccountToJsonAccount(a)
	})

	c.JSON(http.StatusOK, LookupAccountsResponse{Accounts: pAccounts})
}

func (s *server) LookupTransfers(c *gin.Context) {
	var req *LookupTransfersRequest
	if ok := bindJSON(c, req); !ok {
		return
	}
	if len(req.TransferIds) == 0 {
		abort(c, http.StatusInternalServerError, ErrZeroTransfers)
		return
	}
	ids := []types.Uint128{}
	for _, inID := range req.TransferIds {
		id, err := hexStringToUint128(inID)
		if err != nil {
			abort(c, http.StatusInternalServerError, err)
			return
		}
		ids = append(ids, *id)
	}

	res, err := s.TB.LookupTransfers(ids)
	if err != nil {
		abort(c, http.StatusInternalServerError, err)
		return
	}

	pTransfers := lo.Map(res, func(a types.Transfer, _ int) Transfer {
		return *TransferToJsonTransfer(a)
	})

	c.JSON(http.StatusOK, LookupTransfersResponse{Transfers: pTransfers})
}

func (s *server) GetAccountTransfers(c *gin.Context) {
	var req *GetAccountTransfersRequest
	if ok := bindJSON(c, req); !ok {
		return
	}
	if req.Filter.AccountID == "" {
		abort(c, http.StatusInternalServerError, ErrZeroAccounts)
		return
	}
	tbFilter, err := AccountFilterFromJsonToTigerbeetle(&req.Filter)
	if err != nil {
		abort(c, http.StatusInternalServerError, err)
		return
	}
	res, err := s.TB.GetAccountTransfers(*tbFilter)
	if err != nil {
		abort(c, http.StatusInternalServerError, err)
		return
	}

	pTransfers := lo.Map(res, func(v types.Transfer, _ int) Transfer {
		return *TransferToJsonTransfer(v)
	})

	c.JSON(http.StatusOK, GetAccountTransfersResponse{Transfers: pTransfers})
}

func (s *server) GetAccountHistory(c *gin.Context) {
	var req *GetAccountHistoryRequest
	if ok := bindJSON(c, req); !ok {
		return
	}
	if req.Filter.AccountID == "" {
		abort(c, http.StatusInternalServerError, ErrZeroAccounts)
		return
	}
	tbFilter, err := AccountFilterFromJsonToTigerbeetle(&req.Filter)
	if err != nil {
		abort(c, http.StatusInternalServerError, err)
		return
	}
	res, err := s.TB.GetAccountHistory(*tbFilter)
	if err != nil {
		abort(c, http.StatusInternalServerError, err)
		return
	}

	pBalances := lo.Map(res, func(v types.AccountBalance, _ int) AccountBalance {
		return *AccountBalanceFromTigerbeetleToJson(v)
	})

	c.JSON(http.StatusOK, GetAccountHistoryResponse{AccountBalances: pBalances})
}

func bindJSON[V any](c *gin.Context, v *V) (ok bool) {
	err := c.MustBindWith(v, binding.JSON)
	return err == nil
}

func abort(c *gin.Context, status int, err error) {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	c.Error(err)
	c.String(status, binding.MIMEPlain, err.Error())
}
