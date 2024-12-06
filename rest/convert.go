package rest

import (
	"github.com/lil5/tigerbeetle_api/shared"
	"github.com/samber/lo"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

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
		Timestamp:      shared.TimestampFromUintToString(tbAccount.Timestamp),
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
		Timestamp:       shared.TimestampFromUintToString(tbTransfer.Timestamp),
	}
}

func AccountFilterFromJsonToTigerbeetle(pAccountFilter *AccountFilter) (*types.AccountFilter, error) {
	accountID, err := shared.HexStringToUint128(pAccountFilter.AccountID)
	if err != nil {
		return nil, err
	}

	timestampMin, err := shared.TimestampFromPstringToUint(pAccountFilter.TimestampMin)
	if err != nil {
		return nil, err
	}
	timestampMax, err := shared.TimestampFromPstringToUint(pAccountFilter.TimestampMax)
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
		Timestamp:      shared.TimestampFromUintToString(tbBalance.Timestamp),
	}
}

// either returns filled pointers or an error
func (ud UserData) ToUint() (ud128 types.Uint128, ud64 uint64, ud32 uint32, err error) {
	if ud.UserData128 != nil {
		var ud128P *types.Uint128
		ud128P, err = shared.HexStringToUint128(*ud.UserData128)
		if err != nil {
			return
		}
		ud128 = *ud128P
	}

	if ud.UserData64 != nil {
		var ud64P *uint64
		ud64P, err = shared.TimestampFromStringToUint(*ud.UserData64)
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
		ud.UserData64 = lo.ToPtr(shared.TimestampFromUintToString(ud64))
	}
	if ud32 != 0 {
		ud.UserData32 = lo.ToPtr(int32(ud32))
	}
	return
}
