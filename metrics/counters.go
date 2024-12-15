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

	TotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tigerbeetleapi_grpc_requests_total",
		Help: "The total number of grpc requests",
	})
)

func AddTotalRequests(enabled bool) {
	if enabled {
		TotalRequests.Inc()
	}
}
