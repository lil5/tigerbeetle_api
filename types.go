package main

import (
	"tigerbeetle_grpc/proto"
	"time"

	"github.com/samber/lo"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"golang.org/x/exp/slog"
)

func hexStringToUint128(hex string) (*types.Uint128, error) {
	if hex == "" {
		return &types.Uint128{}, nil
	}

	res, err := types.HexStringToUint128(hex)
	if err != nil {
		slog.Error("hex string to Uint128 failed", "hex", hex, "error", err)
		return nil, err
	}
	return &res, nil
}

// set to zero if timestamp is nil
func timestampFromPstringToUint(timestamp *string) (*uint64, error) {
	if timestamp == nil {
		return lo.ToPtr[uint64](0), nil
	}

	return timestampFromStringToUint(*timestamp)
}
func timestampFromUintToString(timestamp uint64) string {
	return time.Unix(0, int64(timestamp)).Format(time.RFC3339Nano)
}

func timestampFromStringToUint(timestamp string) (*uint64, error) {
	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return nil, err
	}

	nano := t.UnixNano()

	return lo.ToPtr(uint64(nano)), nil
}

func AccountToProtoAccount(tbAccount types.Account) *proto.Account {
	tbFlags := tbAccount.AccountFlags()
	pFlags := proto.AccountFlags{
		Linked:                     lo.ToPtr(tbFlags.Linked),
		DebitsMustNotExceedCredits: lo.ToPtr(tbFlags.DebitsMustNotExceedCredits),
		CreditsMustNotExceedDebits: lo.ToPtr(tbFlags.CreditsMustNotExceedDebits),
		History:                    lo.ToPtr(tbFlags.History),
	}
	return &proto.Account{
		Id:             tbAccount.ID.String(),
		DebitsPending:  lo.ToPtr(tbAccount.DebitsPending.BigInt()).Int64(),
		DebitsPosted:   lo.ToPtr(tbAccount.DebitsPosted.BigInt()).Int64(),
		CreditsPending: lo.ToPtr(tbAccount.CreditsPending.BigInt()).Int64(),
		CreditsPosted:  lo.ToPtr(tbAccount.CreditsPosted.BigInt()).Int64(),
		UserData128:    lo.ToPtr(tbAccount.UserData128.BigInt()).Int64(),
		UserData64:     int64(tbAccount.UserData64),
		UserData32:     int32(tbAccount.UserData32),
		Ledger:         int64(tbAccount.Ledger),
		Code:           int32(tbAccount.Code),
		Flags:          &pFlags,
		Timestamp:      timestampFromUintToString(tbAccount.Timestamp),
	}
}

func TransferToProtoTransfer(tbTransfer types.Transfer) *proto.Transfer {
	tbFlags := tbTransfer.TransferFlags()
	pFlags := &proto.TransferFlags{
		Linked:              lo.ToPtr(tbFlags.Linked),
		Pending:             lo.ToPtr(tbFlags.Pending),
		PostPendingTransfer: lo.ToPtr(tbFlags.PostPendingTransfer),
		VoidPendingTransfer: lo.ToPtr(tbFlags.VoidPendingTransfer),
		BalancingDebit:      lo.ToPtr(tbFlags.BalancingDebit),
		BalancingCredit:     lo.ToPtr(tbFlags.BalancingCredit),
	}
	var pendingId string
	emptyUint128 := types.Uint128{}
	if tbTransfer.PendingID != emptyUint128 {
		pendingId = tbTransfer.PendingID.String()
	}
	return &proto.Transfer{
		Id:              tbTransfer.ID.String(),
		DebitAccountId:  tbTransfer.DebitAccountID.String(),
		CreditAccountId: tbTransfer.CreditAccountID.String(),
		Amount:          lo.ToPtr(tbTransfer.Amount.BigInt()).Int64(),
		PendingId:       lo.If[*string](pendingId == "", nil).Else(&pendingId),
		UserData128:     lo.ToPtr(tbTransfer.UserData128.BigInt()).Int64(),
		UserData64:      int64(tbTransfer.UserData64),
		UserData32:      int32(tbTransfer.UserData32),
		Ledger:          int64(tbTransfer.Ledger),
		Code:            int32(tbTransfer.Code),
		TransferFlags:   pFlags,
		Timestamp:       lo.ToPtr(timestampFromUintToString(tbTransfer.Timestamp)),
	}
}

func AccountFilterFromProtoToTigerbeetle(pAccountFilter *proto.AccountFilter) (*types.AccountFilter, error) {
	accountID, err := hexStringToUint128(pAccountFilter.AccountId)
	if err != nil {
		return nil, err
	}

	timestampMin, err := timestampFromPstringToUint(pAccountFilter.TimestampMin)
	if err != nil {
		return nil, err
	}
	timestampMax, err := timestampFromPstringToUint(pAccountFilter.TimestampMax)
	if err != nil {
		return nil, err
	}

	var tbFlags types.AccountFilterFlags
	if pAccountFilter.Flags != nil {
		tbFlags = types.AccountFilterFlags{
			Debits:   lo.FromPtrOr(pAccountFilter.Flags.Debits, false),
			Credits:  lo.FromPtrOr(pAccountFilter.Flags.Credits, false),
			Reversed: lo.FromPtrOr(pAccountFilter.Flags.Reserved, false),
		}
	}

	return &types.AccountFilter{
		AccountID:    *accountID,
		TimestampMin: *timestampMin,
		TimestampMax: *timestampMax,
		Limit:        uint32(pAccountFilter.Limit),
		Flags:        tbFlags.ToUint32(),
	}, nil
}

func AccountBalanceFromTigerbeetleToProto(tbBalance types.AccountBalance) *proto.AccountBalance {
	return &proto.AccountBalance{
		DebitsPending:  lo.ToPtr(tbBalance.DebitsPending.BigInt()).Int64(),
		DebitsPosted:   lo.ToPtr(tbBalance.DebitsPosted.BigInt()).Int64(),
		CreditsPending: lo.ToPtr(tbBalance.CreditsPending.BigInt()).Int64(),
		CreditsPosted:  lo.ToPtr(tbBalance.CreditsPosted.BigInt()).Int64(),
		Timestamp:      timestampFromUintToString(tbBalance.Timestamp),
	}
}
