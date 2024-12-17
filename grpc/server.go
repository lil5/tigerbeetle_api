package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/lil5/tigerbeetle_api/metrics"
	"github.com/lil5/tigerbeetle_api/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
)

func NewServer() {
	networkType := "tcp"
	if Config.OnlyIpv4 {
		networkType = "ipv4"
	}
	lis, err := net.Listen(networkType, fmt.Sprintf("%s:%s", Config.Host, Config.Port))
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	app := NewApp()
	defer app.Close()
	proto.RegisterTigerBeetleServer(s, app)

	if Config.GrpcHealthServer {
		healthServer := health.NewServer()
		healthpb.RegisterHealthServer(s, healthServer)
		healthServer.SetServingStatus("tigerbeetle.TigerBeetle", healthpb.HealthCheckResponse_SERVING)
	}

	if Config.GrpcReflection {
		reflection.Register(s)
	}

	prometheusDeferClose := metrics.Register(Config.PrometheusAddr, promhttp.HandlerFor(metrics.Reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))
	defer prometheusDeferClose()
	srvMetrics := grpcprom.NewServerMetrics(grpcprom.WithServerHandlingTimeHistogram(
		grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
	))
	metrics.Reg.MustRegister(srvMetrics)
	srvMetrics.InitializeMetrics(s)

	slog.Info("GRPC server listening at", "address", lis.Addr())
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve:", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
