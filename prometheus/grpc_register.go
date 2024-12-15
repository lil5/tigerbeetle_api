package prometheus

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func GrpcRegister(s *grpc.Server, addr string) func() {
	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: addr, Handler: nil}
	go server.ListenAndServe()
	return func() { server.Shutdown(context.TODO()) }
}
