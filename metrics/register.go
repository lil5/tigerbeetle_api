package metrics

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Register(addr string) func() {
	h := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
	mux := http.NewServeMux()
	mux.Handle("/metrics", h)
	server := &http.Server{Addr: addr, Handler: mux}
	slog.Info("Prometheus server listening at", "address", addr, "path", "/metrics")
	go server.ListenAndServe()
	return func() { server.Shutdown(context.TODO()) }
}
