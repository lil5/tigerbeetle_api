package metrics

import (
	"context"
	"log/slog"
	"net/http"
)

func Register(addr string, h http.Handler) func() {
	http.Handle("/metrics", h)
	server := &http.Server{Addr: addr, Handler: nil}
	slog.Info("Prometheus server listening at", "address", addr, "path", "/metrics")
	go server.ListenAndServe()
	return func() { server.Shutdown(context.TODO()) }
}
