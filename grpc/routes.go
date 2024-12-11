package grpc

import (
	"context"
	"errors"
	"math/rand/v2"

	"github.com/lil5/tigerbeetle_api/proto"
	"github.com/samber/lo"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

var (
	ErrZeroAccounts  = errors.New("no accounts were specified")
	ErrZeroTransfers = errors.New("no transfers were specified")
)

type TimedPayloadResponse struct {
	Results []types.TransferEventResult
	Error   error
}
type TimedPayload struct {
	c       chan TimedPayloadResponse
	payload []types.Transfer
}

type AppTBs struct {
	TB      tb.Client
	TBs     []tb.Client
	SizeTBs int64
}

func (b *AppTBs) Close() {
	for _, tb := range b.TBs {
		tb.Close()
	}
}

type App struct {
	proto.UnimplementedTigerBeetleServer
	AppTBs
}

func NewApp(tbs AppTBs) *App {
	app := &App{AppTBs: tbs}
	return app
}

func (s *App) GetID(ctx context.Context, in *proto.GetIDRequest) (*proto.GetIDReply, error) {
	return &proto.GetIDReply{Id: types.ID().String()}, nil
}

func (s *App) CreateAccounts(ctx context.Context, in *proto.CreateAccountsRequest) (*proto.CreateAccountsReply, error) {
	if len(in.Accounts) == 0 {
		return nil, ErrZeroAccounts
	}
	accounts := []types.Account{}
	for _, inAccount := range in.Accounts {
		id, err := HexStringToUint128(inAccount.Id)
		if err != nil {
			return nil, err
		}
		userData128, err := types.HexStringToUint128(inAccount.UserData128)
		if err != nil {
			return nil, err
		}
		flags := types.AccountFlags{}
		if inAccount.Flags != nil {
			flags.Linked = lo.FromPtrOr(inAccount.Flags.Linked, false)
			flags.DebitsMustNotExceedCredits = lo.FromPtrOr(inAccount.Flags.DebitsMustNotExceedCredits, false)
			flags.CreditsMustNotExceedDebits = lo.FromPtrOr(inAccount.Flags.CreditsMustNotExceedDebits, false)
			flags.History = lo.FromPtrOr(inAccount.Flags.History, false)
		}
		accounts = append(accounts, types.Account{
			ID:             *id,
			DebitsPending:  types.ToUint128(uint64(inAccount.DebitsPending)),
			DebitsPosted:   types.ToUint128(uint64(inAccount.DebitsPosted)),
			CreditsPending: types.ToUint128(uint64(inAccount.CreditsPending)),
			CreditsPosted:  types.ToUint128(uint64(inAccount.CreditsPosted)),
			UserData128:    userData128,
			UserData64:     uint64(inAccount.UserData64),
			UserData32:     uint32(inAccount.UserData32),
			Ledger:         uint32(inAccount.Ledger),
			Code:           uint16(inAccount.Code),
			Flags:          flags.ToUint16(),
		})
	}

	results, err := s.TB.CreateAccounts(accounts)
	if err != nil {
		return nil, err
	}

	resArr := []*proto.CreateAccountsReplyItem{}
	for _, r := range results {
		resArr = append(resArr, &proto.CreateAccountsReplyItem{
			Index:  int32(r.Index),
			Result: proto.CreateAccountResult(r.Result),
		})
	}
	return &proto.CreateAccountsReply{
		Results: resArr,
	}, nil
}

func (s *App) CreateTransfers(ctx context.Context, in *proto.CreateTransfersRequest) (*proto.CreateTransfersReply, error) {
	if len(in.Transfers) == 0 {
		return nil, ErrZeroTransfers
	}
	transfers := []types.Transfer{}
	for _, inTransfer := range in.Transfers {
		id, err := HexStringToUint128(inTransfer.Id)
		if err != nil {
			return nil, err
		}
		flags := types.TransferFlags{}
		if inTransfer.TransferFlags != nil {
			flags.Linked = lo.FromPtrOr(inTransfer.TransferFlags.Linked, false)
			flags.Pending = lo.FromPtrOr(inTransfer.TransferFlags.Pending, false)
			flags.PostPendingTransfer = lo.FromPtrOr(inTransfer.TransferFlags.PostPendingTransfer, false)
			flags.VoidPendingTransfer = lo.FromPtrOr(inTransfer.TransferFlags.VoidPendingTransfer, false)
			flags.BalancingDebit = lo.FromPtrOr(inTransfer.TransferFlags.BalancingDebit, false)
			flags.BalancingCredit = lo.FromPtrOr(inTransfer.TransferFlags.BalancingCredit, false)
		}
		debitAccountID, err := HexStringToUint128(inTransfer.DebitAccountId)
		if err != nil {
			return nil, err
		}
		creditAccountID, err := HexStringToUint128(inTransfer.CreditAccountId)
		if err != nil {
			return nil, err
		}
		pendingID, err := HexStringToUint128(lo.FromPtrOr(inTransfer.PendingId, ""))
		if err != nil {
			return nil, err
		}
		userData128, err := types.HexStringToUint128(inTransfer.UserData128)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, types.Transfer{
			ID:              *id,
			DebitAccountID:  *debitAccountID,
			CreditAccountID: *creditAccountID,
			Amount:          types.ToUint128(uint64(inTransfer.Amount)),
			PendingID:       *pendingID,
			UserData128:     userData128,
			UserData64:      uint64(inTransfer.UserData64),
			UserData32:      uint32(inTransfer.UserData32),
			Timeout:         0,
			Ledger:          uint32(inTransfer.Ledger),
			Code:            uint16(inTransfer.Ledger),
			Flags:           flags.ToUint16(),
			Timestamp:       0,
		})
	}

	i := rand.Int64N(s.AppTBs.SizeTBs - 1)
	tb := s.AppTBs.TBs[i]
	results, err := tb.CreateTransfers(transfers)

	if err != nil {
		return nil, err
	}
	resArr := []*proto.CreateTransfersReplyItem{}
	for _, r := range results {
		resArr = append(resArr, &proto.CreateTransfersReplyItem{
			Index:  int32(r.Index),
			Result: proto.CreateTransferResult(r.Result),
		})
	}
	return &proto.CreateTransfersReply{
		Results: resArr,
	}, nil
}

func (s *App) LookupAccounts(ctx context.Context, in *proto.LookupAccountsRequest) (*proto.LookupAccountsReply, error) {
	if len(in.AccountIds) == 0 {
		return nil, ErrZeroAccounts
	}
	ids := []types.Uint128{}
	for _, inID := range in.AccountIds {
		id, err := HexStringToUint128(inID)
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

func (s *App) LookupTransfers(ctx context.Context, in *proto.LookupTransfersRequest) (*proto.LookupTransfersReply, error) {
	if len(in.TransferIds) == 0 {
		return nil, ErrZeroTransfers
	}
	ids := []types.Uint128{}
	for _, inID := range in.TransferIds {
		id, err := HexStringToUint128(inID)
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

func (s *App) GetAccountTransfers(ctx context.Context, in *proto.GetAccountTransfersRequest) (*proto.GetAccountTransfersReply, error) {
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

func (s *App) GetAccountBalances(ctx context.Context, in *proto.GetAccountBalancesRequest) (*proto.GetAccountBalancesReply, error) {
	if in.Filter.AccountId == "" {
		return nil, ErrZeroAccounts
	}
	tbFilter, err := AccountFilterFromProtoToTigerbeetle(in.Filter)
	if err != nil {
		return nil, err
	}
	res, err := s.TB.GetAccountBalances(*tbFilter)
	if err != nil {
		return nil, err
	}

	pBalances := lo.Map(res, func(v types.AccountBalance, _ int) *proto.AccountBalance {
		return AccountBalanceFromTigerbeetleToProto(v)
	})
	return &proto.GetAccountBalancesReply{AccountBalances: pBalances}, nil
}
