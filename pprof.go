//go:build pprof

package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	slog.Info("pprof enabled")
	go func() {
		slog.Warn("pprof stopped", "err", http.ListenAndServe("localhost:6060", nil))
	}()
}
