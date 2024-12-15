package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var BufferFullRatio = promauto.NewCounter(prometheus.CounterOpts{
	Name: "tigerbeetle_buffer_full_ratio",
})
