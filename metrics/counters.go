package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalBufferContents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_contents_total",
		Help: "Tigerbeetle requests buffered filled size sum",
	})

	TotalBufferMax = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_max_total",
		Help: "Tigerbeetle requests buffer max size sum",
	})

	TotalBufferCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_count_total",
		Help: "Tigerbeetle requests total buffers",
	})

	TotalCreateTransferTx = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_create_transfers_tx_total",
		Help: "Created transfer transactions",
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
