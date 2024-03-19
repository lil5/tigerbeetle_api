export interface TigerBeetle {
  GetID: () => Promise<GetIDResponse>
  CreateAccounts: (req: CreateAccountsRequest) => Promise<CreateAccountsResponse>
  CreateTransfers: (req: CreateTransfersRequest) => Promise<CreateTransfersResponse>
  LookupAccounts: (req: LookupAccountsRequest) => Promise<LookupAccountsResponse>
  LookupTransfers: (req: LookupTransfersRequest) => Promise<LookupTransfersResponse>
  GetAccountTransfers: (req: GetAccountTransfersRequest) => Promise<GetAccountTransfersResponse>
  GetAccountBalances: (req: GetAccountBalancesRequest) => Promise<GetAccountBalancesResponse>
}

type int64 = number
type int32 = number
type bool = boolean

export interface GetIDResponse {
  id: string;
}

export interface CreateAccountsRequest {
  accounts: Account[];
}
export interface CreateAccountsResponse {
  results: string[];
}
export interface CreateTransfersRequest {
  transfers: Transfer[];

}
export interface CreateTransfersResponse {
  results: string[];
}
export interface LookupAccountsRequest {
  account_ids: string[];
}
export interface LookupAccountsResponse {
  accounts: Account[];
}
export interface LookupTransfersRequest {
  transfer_ids: string[];
}
export interface LookupTransfersResponse {
  transfers: Transfer[];
}
export interface GetAccountTransfersRequest {
  filter: AccountFilter;
}
export interface GetAccountTransfersResponse {
  transfers: Transfer[];
}
export interface GetAccountBalancesRequest {
  filter: AccountFilter;
}
export interface GetAccountBalancesResponse {
  account_balances: AccountBalance[];
}


// Types
// ----------------------------------------------------------------
export interface Account extends UserData {
  id: string;
  debits_pending: int64;
  debits_posted: int64;
  credits_pending: int64;
  credits_posted: int64;
  ledger: int64;
  code: int32;
  flags?: AccountFlags;
  timestamp?: string;
}

export interface AccountFlags {
  linked?: bool;
  debits_must_not_exceed_credits?: bool;
  credits_must_not_exceed_debits?: bool;
  history?: bool;
}

export interface Transfer extends UserData {
  id: string;
  debit_account_id: string;
  credit_account_id: string;
  amount: int64;
  pending_id?: string;
  ledger: int64;
  code: int32;
  transfer_flags?: TransferFlags;
  timestamp?: string;
}

export interface TransferFlags {
  linked?: bool;
  pending?: bool;
  post_pending_transfer?: bool;
  void_pending_transfer?: bool;
  balancing_debit?: bool;
  balancing_credit?: bool;
}

export interface AccountFilter {
  account_id: string;
  timestamp_min?: string;
  timestamp_max?: string;
  limit: int32;
  flags?: AccountFilterFlags;
}

export interface AccountFilterFlags {
  debits?: bool;
  credits?: bool;
  reserved?: bool;
}

export interface AccountBalance {
  debits_pending: int64;
  debits_posted: int64;
  credits_pending: int64;
  credits_posted: int64;
  timestamp: string;
}

interface UserData {
  user_data128?: string;
  user_data64?: int64;
  user_data32?: int32;
}