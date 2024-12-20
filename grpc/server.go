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
	"github.com/piotrkowalczuk/promgrpc/v4"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
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

	ssh := promgrpc.ServerStatsHandler()

	srvOpts := []grpc.ServerOption{
		grpc.StatsHandler(ssh),
	}
	if config.Config.GrpcHighScale {
		srvOpts = append(srvOpts,
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionAge: 5 * time.Minute,
				Time:             60 * time.Second,
				Timeout:          20 * time.Second,
			}),
			grpc.MaxConcurrentStreams(50_000),
		)
	}
	s := grpc.NewServer(srvOpts...)
	prometheus.DefaultRegisterer.MustRegister(ssh)

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

	slog.Info("GRPC server listening at", "address", lis.Addr())
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve:", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
