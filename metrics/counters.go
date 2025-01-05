package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalBufferContentsFull = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_contents_full_total",
		Help: "Amount of payloads / max size of the buffer = 100 percent",
	})
	TotalBufferContentsLt80 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_contents_lt80_total",
		Help: "Amount of payloads / max size of the buffer = 80 percent or less",
	})
	TotalBufferContentsGte80 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_contents_gte80_total",
		Help: "Amount of payloads / max size of the buffer = 80 percent or more, but not 100 percent",
	})

	TotalBufferCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_count_total",
		Help: "Counter for each time the buffer is flushed",
	})

	TotalCreateTransferTx = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_create_transfers_tx_total",
		Help: "Counter for each tranfer created",
	})

	TotalCreateTransferTxErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_create_transfers_tx_error_total",
		Help: "Counter for each error sent back for each transfer",
	})

	TotalCreateAccountsTxErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_create_accounts_tx_error_total",
		Help: "Counter for each account create error",
	})

	TotalTbCreateAccountsCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_create_accounts_total",
		Help: "Called when tigerbeetle client create_accounts is run",
	})

	TotalTbCreateTransfersCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_create_transfers_total",
		Help: "Called when tigerbeetle client create_transfers is run",
	})

	TotalTbLookupAccountsCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_lookup_accounts_total",
		Help: "Called when tigerbeetle client lookup_accounts is run",
	})

	TotalTbLookupTransfersCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_lookup_transfers_total",
		Help: "Called when tigerbeetle client lookup_transfers is run",
	})

	TotalTbGetAccountTransfersCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_get_account_transfers_total",
		Help: "Called when tigerbeetle client get_account_transfers is run",
	})

	TotalTbGetAccountBalancesCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_get_account_balances_total",
		Help: "Called when tigerbeetle client get_account_balances is run",
	})
)
