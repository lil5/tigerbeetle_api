package app

import (
	"log/slog"
	"time"

	"github.com/samber/lo"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
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
	if *timestamp == "" {
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

func AccountToJsonAccount(tbAccount types.Account) *Account {
	tbFlags := tbAccount.AccountFlags()
	pFlags := AccountFlags{
		Linked:                     tbFlags.Linked,
		DebitsMustNotExceedCredits: tbFlags.DebitsMustNotExceedCredits,
		CreditsMustNotExceedDebits: tbFlags.CreditsMustNotExceedDebits,
		History:                    tbFlags.History,
	}
	return &Account{
		UserData:       toUserData(tbAccount.UserData128, tbAccount.UserData64, tbAccount.UserData32),
		ID:             tbAccount.ID.String(),
		DebitsPending:  lo.ToPtr(tbAccount.DebitsPending.BigInt()).Int64(),
		DebitsPosted:   lo.ToPtr(tbAccount.DebitsPosted.BigInt()).Int64(),
		CreditsPending: lo.ToPtr(tbAccount.CreditsPending.BigInt()).Int64(),
		CreditsPosted:  lo.ToPtr(tbAccount.CreditsPosted.BigInt()).Int64(),
		Ledger:         int64(tbAccount.Ledger),
		Code:           int32(tbAccount.Code),
		Flags:          &pFlags,
		Timestamp:      timestampFromUintToString(tbAccount.Timestamp),
	}
}

func TransferToJsonTransfer(tbTransfer types.Transfer) *Transfer {
	tbFlags := tbTransfer.TransferFlags()
	pFlags := &TransferFlags{
		Linked:              tbFlags.Linked,
		Pending:             tbFlags.Pending,
		PostPendingTransfer: tbFlags.PostPendingTransfer,
		VoidPendingTransfer: tbFlags.VoidPendingTransfer,
		BalancingDebit:      tbFlags.BalancingDebit,
		BalancingCredit:     tbFlags.BalancingCredit,
	}
	var pendingId string
	emptyUint128 := types.Uint128{}
	if tbTransfer.PendingID != emptyUint128 {
		pendingId = tbTransfer.PendingID.String()
	}
	return &Transfer{
		UserData:        toUserData(tbTransfer.UserData128, tbTransfer.UserData64, tbTransfer.UserData32),
		ID:              tbTransfer.ID.String(),
		DebitAccountID:  tbTransfer.DebitAccountID.String(),
		CreditAccountID: tbTransfer.CreditAccountID.String(),
		Amount:          lo.ToPtr(tbTransfer.Amount.BigInt()).Int64(),
		PendingID:       lo.If[*string](pendingId == "", nil).Else(&pendingId),
		Ledger:          int64(tbTransfer.Ledger),
		Code:            int32(tbTransfer.Code),
		TransferFlags:   pFlags,
		Timestamp:       timestampFromUintToString(tbTransfer.Timestamp),
	}
}

func AccountFilterFromJsonToTigerbeetle(pAccountFilter *AccountFilter) (*types.AccountFilter, error) {
	accountID, err := hexStringToUint128(pAccountFilter.AccountID)
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
			Debits:   pAccountFilter.Flags.Debits,
			Credits:  pAccountFilter.Flags.Credits,
			Reversed: pAccountFilter.Flags.Reversed,
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

func AccountBalanceFromTigerbeetleToJson(tbBalance types.AccountBalance) *AccountBalance {
	return &AccountBalance{
		DebitsPending:  lo.ToPtr(tbBalance.DebitsPending.BigInt()).Int64(),
		DebitsPosted:   lo.ToPtr(tbBalance.DebitsPosted.BigInt()).Int64(),
		CreditsPending: lo.ToPtr(tbBalance.CreditsPending.BigInt()).Int64(),
		CreditsPosted:  lo.ToPtr(tbBalance.CreditsPosted.BigInt()).Int64(),
		Timestamp:      timestampFromUintToString(tbBalance.Timestamp),
	}
}

// either returns filled pointers or an error
func (ud UserData) ToUint() (ud128 types.Uint128, ud64 uint64, ud32 uint32, err error) {
	if ud.UserData128 != nil {
		var ud128P *types.Uint128
		ud128P, err = hexStringToUint128(*ud.UserData128)
		if err != nil {
			return
		}
		ud128 = *ud128P
	}

	if ud.UserData64 != nil {
		var ud64P *uint64
		ud64P, err = timestampFromStringToUint(*ud.UserData64)
		if err != nil {
			return
		}
		ud64 = *ud64P
	}

	if ud.UserData32 != nil {
		ud32 = uint32(*ud.UserData32)
	}

	return
}

func toUserData(ud128 types.Uint128, ud64 uint64, ud32 uint32) (ud UserData) {
	if lo.ToPtr(ud128.BigInt()).Int64() != 0 {
		ud.UserData128 = lo.ToPtr(ud128.String())
	}
	if ud64 != 0 {
		ud.UserData64 = lo.ToPtr(timestampFromUintToString(ud64))
	}
	if ud32 != 0 {
		ud.UserData32 = lo.ToPtr(int32(ud32))
	}
	return
}
