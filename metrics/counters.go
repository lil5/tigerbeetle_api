package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Reg = prometheus.NewRegistry()

	TotalBufferContents = promauto.With(Reg).NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_contents_total",
		Help: "Tigerbeetle requests buffered filled size sum",
	})

	TotalBufferMax = promauto.With(Reg).NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_max_total",
		Help: "Tigerbeetle requests buffer max size sum",
	})

	TotalBufferCount = promauto.With(Reg).NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_buffer_count_total",
		Help: "Tigerbeetle requests total buffers",
	})

	TotalCreateTransferTx = promauto.With(Reg).NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_create_transfers_tx_total",
		Help: "Created transfer transactions",
	})
)
