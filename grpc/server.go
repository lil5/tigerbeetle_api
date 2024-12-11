package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/lil5/tigerbeetle_api/proto"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewServer(tb tigerbeetle_go.Client) {
	networkType := "tcp"
	if os.Getenv("ONLY_IPV4") == "true" {
		networkType = "ipv4"
	}
	lis, err := net.Listen(networkType, fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	app := NewApp(tb)
	defer app.TimedBuf.Close()
	proto.RegisterTigerBeetleServer(s, app)

	if os.Getenv("GRPC_HEALTH_SERVER") == "true" {
		healthServer := health.NewServer()
		healthpb.RegisterHealthServer(s, healthServer)
		healthServer.SetServingStatus("tigerbeetle.TigerBeetle", healthpb.HealthCheckResponse_SERVING)
	}

	if os.Getenv("GRPC_REFLECTION") == "true" {
		reflection.Register(s)
	}

	slog.Info("GRPC server listening at", "address", lis.Addr())
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve:", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
