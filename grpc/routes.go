package grpc

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"os"

	"github.com/charithe/timedbuf/v2"
	"github.com/lil5/tigerbeetle_api/metrics"
	"github.com/lil5/tigerbeetle_api/proto"
	"github.com/samber/lo"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
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
	c         chan TimedPayloadResponse
	Transfers []types.Transfer
}

type App struct {
	proto.UnimplementedTigerBeetleServer

	TB tigerbeetle_go.Client

	TBuf  *timedbuf.TimedBuf[TimedPayload]
	TBufs []*timedbuf.TimedBuf[TimedPayload]
}

func (a *App) getRandomTBuf() *timedbuf.TimedBuf[TimedPayload] {
	if Config.BufferCluster > 1 {
		i := rand.IntN(Config.BufferCluster - 1)
		return a.TBufs[i]
	} else {
		return a.TBuf
	}
}

func (a *App) Close() {
	for _, b := range a.TBufs {
		b.Close()
	}
	a.TB.Close()
}

func NewApp() *App {
	tigerbeetle_go, err := tigerbeetle_go.NewClient(types.Uint128{uint8(Config.TbClusterID)}, Config.TbAddresses)
	if err != nil {
		slog.Error("unable to connect to tigerbeetle", "err", err)
		os.Exit(1)
	}

	var tbuf *timedbuf.TimedBuf[TimedPayload]
	var tbufs []*timedbuf.TimedBuf[TimedPayload]
	if Config.IsBuffered {
		tbufs = make([]*timedbuf.TimedBuf[TimedPayload], Config.BufferCluster)

		lenMaxBuf := float64(Config.BufferSize)
		lenMaxBufSlog := lenMaxBuf * 0.8
		flushFunc := func(payloads []TimedPayload) {
			transfers := []types.Transfer{}
			lenPayloads := float64(len(payloads))
			if Config.PrometheusEnabled {
				metrics.BufferFullRatio.Add(lenPayloads / lenMaxBuf)
			} else if lenPayloads < lenMaxBufSlog {
				slog.Warn("Flushing Buffer", "max buffer", Config.BufferSize, "buffer size collected", lenPayloads)
			}
			for _, payload := range payloads {
				transfers = append(transfers, payload.Transfers...)
			}
			results, err := tigerbeetle_go.CreateTransfers(transfers)
			res := TimedPayloadResponse{
				Results: results,
				Error:   err,
			}
			for _, payload := range payloads {
				payload.c <- res
			}
		}
		for i := range Config.BufferCluster {
			tbufs[i] = timedbuf.New(Config.BufferSize, Config.BufferDelay, flushFunc)
		}
		tbuf = tbufs[0]
	}

	app := &App{
		TB:    tigerbeetle_go,
		TBuf:  tbuf,
		TBufs: tbufs,
	}
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

	var results []types.TransferEventResult
	var err error
	if Config.IsBuffered {
		buf := s.getRandomTBuf()
		c := make(chan TimedPayloadResponse)
		buf.Put(TimedPayload{
			c:         c,
			Transfers: transfers,
		})
		res := <-c
		results = res.Results
		err = res.Error
	} else {
		results, err = s.TB.CreateTransfers(transfers)
	}

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
