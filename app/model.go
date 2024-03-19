package app

// Request & Response
// ----------------------------------------------------------------

type GetIDResponse struct {
	ID string `json:"id"`
}

type CreateAccountsRequest struct {
	Accounts []Account `json:"accounts"`
}
type CreateAccountsResponse struct {
	Results []string `json:"results"`
}
type CreateTransfersRequest struct {
	Transfers []Transfer `json:"transfers"`
}
type CreateTransfersResponse struct {
	Results []string `json:"results"`
}
type LookupAccountsRequest struct {
	AccountIds []string `json:"account_ids"`
}
type LookupAccountsResponse struct {
	Accounts []Account `json:"accounts"`
}
type LookupTransfersRequest struct {
	TransferIds []string `json:"transfer_ids"`
}
type LookupTransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}
type GetAccountTransfersRequest struct {
	Filter AccountFilter `json:"filter"`
}
type GetAccountTransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}
type GetAccountBalancesRequest struct {
	Filter AccountFilter `json:"filter"`
}
type GetAccountBalancesResponse struct {
	AccountBalances []AccountBalance `json:"account_balances"`
}

// Types
// ----------------------------------------------------------------
type Account struct {
	UserData
	ID             string        `json:"id"`
	DebitsPending  int64         `json:"debits_pending"`
	DebitsPosted   int64         `json:"debits_posted"`
	CreditsPending int64         `json:"credits_pending"`
	CreditsPosted  int64         `json:"credits_posted"`
	Ledger         int64         `json:"ledger"`
	Code           int32         `json:"code"`
	Flags          *AccountFlags `json:"flags"`
	Timestamp      string        `json:"timestamp"`
}

type AccountFlags struct {
	Linked                     bool `json:"linked"`
	DebitsMustNotExceedCredits bool `json:"debits_must_not_exceed_credits"`
	CreditsMustNotExceedDebits bool `json:"credits_must_not_exceed_debits"`
	History                    bool `json:"history"`
}

type Transfer struct {
	UserData
	ID              string         `json:"id"`
	DebitAccountID  string         `json:"debit_account_id"`
	CreditAccountID string         `json:"credit_account_id"`
	Amount          int64          `json:"amount"`
	PendingID       *string        `json:"pending_id"`
	Ledger          int64          `json:"ledger"`
	Code            int32          `json:"code"`
	TransferFlags   *TransferFlags `json:"transfer_flags"`
	Timestamp       string         `json:"timestamp"`
}

type TransferFlags struct {
	Linked              bool `json:"linked"`
	Pending             bool `json:"pending"`
	PostPendingTransfer bool `json:"post_pending_transfer"`
	VoidPendingTransfer bool `json:"void_pending_transfer"`
	BalancingDebit      bool `json:"balancing_debit"`
	BalancingCredit     bool `json:"balancing_credit"`
}

type AccountFilter struct {
	AccountID    string              `json:"account_id"`
	TimestampMin *string             `json:"timestamp_min"`
	TimestampMax *string             `json:"timestamp_max"`
	Limit        int64               `json:"limit"`
	Flags        *AccountFilterFlags `json:"flags"`
}

type AccountFilterFlags struct {
	Debits   bool `json:"debits"`
	Credits  bool `json:"credits"`
	Reserved bool `json:"reserved"`
}

type AccountBalance struct {
	DebitsPending  int64  `json:"debits_pending"`
	DebitsPosted   int64  `json:"debits_posted"`
	CreditsPending int64  `json:"credits_pending"`
	CreditsPosted  int64  `json:"credits_posted"`
	Timestamp      string `json:"timestamp"`
}

type UserData struct {
	UserData128 *string `json:"user_data_128"`
	UserData64  *string `json:"user_data_64"`
	UserData32  *int32  `json:"user_data_32"`
}
