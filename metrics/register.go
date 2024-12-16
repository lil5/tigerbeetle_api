package metrics

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Register(addr string) func() {
	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: addr, Handler: nil}
	slog.Info("Prometheus server listening at", "address", addr, "path", "/metrics")
	go server.ListenAndServe()
	return func() { server.Shutdown(context.TODO()) }
}
