package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalBufferContentsFull = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_contents_full_total",
		Help: "Tigerbeetle buffer filled size is full",
	})
	TotalBufferContentsLt80 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_contents_lt80_total",
		Help: "Tigerbeetle buffer filled size is less than 80%",
	})
	TotalBufferContentsGt80 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_contents_gte80_total",
		Help: "Tigerbeetle buffer filled size is greater than or equal to 80%",
	})

	TotalBufferCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_count_total",
		Help: "Tigerbeetle requests total buffer instances created",
	})

	TotalCreateTransferTx = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_create_transfers_tx_total",
		Help: "Created transfer transactions",
	})

	TotalCreateTransferTxErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_create_transfers_tx_error_total",
		Help: "Created transfer error transactions",
	})

	TotalCreateAccountsTxErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_create_accounts_tx_error_total",
		Help: "Created account error transactions",
	})

	TotalTbCreateAccountsCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_create_accounts_total",
	})

	TotalTbCreateTransfersCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_create_transfers_total",
	})

	TotalTbLookupAccountsCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_lookup_accounts_total",
	})

	TotalTbLookupTransfersCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_lookup_transfers_total",
	})

	TotalTbGetAccountTransfersCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_get_account_transfers_total",
	})

	TotalTbGetAccountBalancesCall = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_tb_get_account_balances_total",
	})
)
