package main

import (
	"context"
	"errors"

	"tigerbeetle_grpc/proto"

	"github.com/samber/lo"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

var (
	ErrZeroAccounts  = errors.New("no accounts were specified")
	ErrZeroTransfers = errors.New("no transfers were specified")
)

type server struct {
	proto.UnimplementedTigerBeetleServer
	TB tb.Client
}

func (s *server) GetID(ctx context.Context, in *proto.GetIDRequest) (*proto.GetIDReply, error) {
	return &proto.GetIDReply{Id: types.ID().String()}, nil
}

func (s *server) CreateAccounts(ctx context.Context, in *proto.CreateAccountsRequest) (*proto.CreateAccountsReply, error) {
	if len(in.Accounts) == 0 {
		return nil, ErrZeroAccounts
	}
	accounts := []types.Account{}
	for _, inAccount := range in.Accounts {
		id, err := hexStringToUint128(inAccount.Id)
		if err != nil {
			return nil, err
		}
		flags := types.AccountFlags{
			Linked:                     lo.FromPtrOr(inAccount.Flags.Linked, false),
			DebitsMustNotExceedCredits: lo.FromPtrOr(inAccount.Flags.DebitsMustNotExceedCredits, false),
			CreditsMustNotExceedDebits: lo.FromPtrOr(inAccount.Flags.CreditsMustNotExceedDebits, false),
			History:                    lo.FromPtrOr(inAccount.Flags.History, false),
		}
		accounts = append(accounts, types.Account{
			ID:             *id,
			DebitsPending:  types.ToUint128(uint64(inAccount.DebitsPending)),
			DebitsPosted:   types.ToUint128(uint64(inAccount.DebitsPosted)),
			CreditsPending: types.ToUint128(uint64(inAccount.CreditsPending)),
			CreditsPosted:  types.ToUint128(uint64(inAccount.CreditsPosted)),
			UserData128:    types.ToUint128(uint64(inAccount.UserData128)),
			UserData64:     uint64(inAccount.UserData64),
			UserData32:     uint32(inAccount.UserData32),
			Ledger:         uint32(inAccount.Ledger),
			Code:           uint16(inAccount.Code),
			Flags:          flags.ToUint16(),
		})
	}

	resp, err := s.TB.CreateAccounts(accounts)
	if err != nil {
		return nil, err
	}

	resArr := []string{}
	for _, r := range resp {
		resArr = append(resArr, r.Result.String())
	}
	return &proto.CreateAccountsReply{
		Results: resArr,
	}, nil
}

func (s *server) CreateTransfers(ctx context.Context, in *proto.CreateTransfersRequest) (*proto.CreateTransfersReply, error) {
	if len(in.Transfers) == 0 {
		return nil, ErrZeroTransfers
	}
	transfers := []types.Transfer{}
	for _, inTransfer := range in.Transfers {
		id, err := hexStringToUint128(inTransfer.Id)
		if err != nil {
			return nil, err
		}
		flags := types.TransferFlags{
			Linked:              lo.FromPtrOr(inTransfer.TransferFlags.Linked, false),
			Pending:             lo.FromPtrOr(inTransfer.TransferFlags.Pending, false),
			PostPendingTransfer: lo.FromPtrOr(inTransfer.TransferFlags.PostPendingTransfer, false),
			VoidPendingTransfer: lo.FromPtrOr(inTransfer.TransferFlags.VoidPendingTransfer, false),
			BalancingDebit:      lo.FromPtrOr(inTransfer.TransferFlags.BalancingDebit, false),
			BalancingCredit:     lo.FromPtrOr(inTransfer.TransferFlags.BalancingCredit, false),
		}
		debitAccountID, err := hexStringToUint128(inTransfer.DebitAccountId)
		if err != nil {
			return nil, err
		}
		creditAccountID, err := hexStringToUint128(inTransfer.CreditAccountId)
		if err != nil {
			return nil, err
		}
		pendingID, err := hexStringToUint128(lo.FromPtrOr(inTransfer.PendingId, ""))
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, types.Transfer{
			ID:              *id,
			DebitAccountID:  *debitAccountID,
			CreditAccountID: *creditAccountID,
			Amount:          types.ToUint128(uint64(inTransfer.Amount)),
			PendingID:       *pendingID,
			UserData128:     types.ToUint128(uint64(inTransfer.UserData128)),
			UserData64:      uint64(inTransfer.UserData64),
			UserData32:      uint32(inTransfer.UserData32),
			Timeout:         0,
			Ledger:          uint32(inTransfer.Ledger),
			Code:            uint16(inTransfer.Ledger),
			Flags:           flags.ToUint16(),
			Timestamp:       0,
		})
	}

	resp, err := s.TB.CreateTransfers(transfers)
	if err != nil {
		return nil, err
	}

	resArr := []string{}
	for _, r := range resp {
		resArr = append(resArr, r.Result.String())
	}
	return &proto.CreateTransfersReply{
		Results: resArr,
	}, nil
}

func (s *server) LookupAccounts(ctx context.Context, in *proto.LookupAccountsRequest) (*proto.LookupAccountsReply, error) {
	if len(in.AccountIds) == 0 {
		return nil, ErrZeroAccounts
	}
	ids := []types.Uint128{}
	for _, inID := range in.AccountIds {
		id, err := hexStringToUint128(inID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, *id)
	}

	res, err := s.TB.LookupAccounts(ids)
	if err != nil {
		return nil, err
	}

	pAccounts := lo.Map(res, func(a types.Account, _ int) *proto.Account {
		return AccountToProtoAccount(a)
	})
	return &proto.LookupAccountsReply{Accounts: pAccounts}, nil
}

func (s *server) LookupTransfers(ctx context.Context, in *proto.LookupTransfersRequest) (*proto.LookupTransfersReply, error) {
	if len(in.TransferIds) == 0 {
		return nil, ErrZeroTransfers
	}
	ids := []types.Uint128{}
	for _, inID := range in.TransferIds {
		id, err := hexStringToUint128(inID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, *id)
	}

	res, err := s.TB.LookupTransfers(ids)
	if err != nil {
		return nil, err
	}

	pTransfers := lo.Map(res, func(a types.Transfer, _ int) *proto.Transfer {
		return TransferToProtoTransfer(a)
	})
	return &proto.LookupTransfersReply{Transfers: pTransfers}, nil
}

func (s *server) GetAccountTransfers(ctx context.Context, in *proto.GetAccountTransfersRequest) (*proto.GetAccountTransfersReply, error) {
	if in.Filter.AccountId == "" {
		return nil, ErrZeroAccounts
	}
	tbFilter, err := AccountFilterFromProtoToTigerbeetle(in.Filter)
	if err != nil {
		return nil, err
	}
	res, err := s.TB.GetAccountTransfers(*tbFilter)
	if err != nil {
		return nil, err
	}

	pTransfers := lo.Map(res, func(v types.Transfer, _ int) *proto.Transfer {
		return TransferToProtoTransfer(v)
	})
	return &proto.GetAccountTransfersReply{Transfers: pTransfers}, nil
}

func (s *server) GetAccountHistory(ctx context.Context, in *proto.GetAccountHistoryRequest) (*proto.GetAccountHistoryReply, error) {
	if in.Filter.AccountId == "" {
		return nil, ErrZeroAccounts
	}
	tbFilter, err := AccountFilterFromProtoToTigerbeetle(in.Filter)
	if err != nil {
		return nil, err
	}
	res, err := s.TB.GetAccountHistory(*tbFilter)
	if err != nil {
		return nil, err
	}

	pBalances := lo.Map(res, func(v types.AccountBalance, _ int) *proto.AccountBalance {
		return AccountBalanceFromTigerbeetleToProto(v)
	})
	return &proto.GetAccountHistoryReply{AccountBalances: pBalances}, nil
}
