package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/lil5/tigerbeetle_api/config"
	"github.com/lil5/tigerbeetle_api/metrics"
	"github.com/lil5/tigerbeetle_api/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
)

func NewServer() {
	networkType := "tcp"
	if config.Config.OnlyIpv4 {
		networkType = "ipv4"
	}
	lis, err := net.Listen(networkType, fmt.Sprintf("%s:%s", config.Config.Host, config.Config.Port))
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	var s *grpc.Server
	if config.Config.GrpcHighScale {
		s = grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionAge: 5 * time.Minute,
				Time:             60 * time.Second,
				Timeout:          20 * time.Second,
			}),
			grpc.MaxConcurrentStreams(50_000),
		)
	} else {
		s = grpc.NewServer()
	}

	app := NewApp()
	defer app.Close()
	proto.RegisterTigerBeetleServer(s, app)

	if config.Config.GrpcHealthServer {
		healthServer := health.NewServer()
		healthpb.RegisterHealthServer(s, healthServer)
		healthServer.SetServingStatus("tigerbeetle.TigerBeetle", healthpb.HealthCheckResponse_SERVING)
	}

	if config.Config.GrpcReflection {
		reflection.Register(s)
	}

	prometheusDeferClose := metrics.Register(config.Config.PrometheusAddr)
	defer prometheusDeferClose()
	srvMetrics := grpcprom.NewServerMetrics(grpcprom.WithServerHandlingTimeHistogram(
		grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
	))
	prometheus.DefaultRegisterer.MustRegister(srvMetrics)
	srvMetrics.InitializeMetrics(s)

	slog.Info("GRPC server listening at", "address", lis.Addr())
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve:", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
