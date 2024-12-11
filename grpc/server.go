package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"

	"github.com/lil5/tigerbeetle_api/proto"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewTbClientSet(tbAddresses []string, tbClusterId uint64) AppTBs {
	clusterUint128 := types.ToUint128(tbClusterId)
	set := AppTBs{}

	bufSize, _ := strconv.Atoi(os.Getenv("CLIENT_COUNT"))
	if bufSize < 1 {
		bufSize = 1
	}

	tbs := make([]tb.Client, int(bufSize))
	for i := range bufSize {
		var err error
		tbs[i], err = tigerbeetle_go.NewClient(clusterUint128, tbAddresses)
		if err != nil {
			slog.Error("unable to connect to tigerbeetle", "err", err)
			os.Exit(1)
		}
	}
	set.TB = tbs[0]
	set.TBs = tbs
	set.SizeTBs = int64(bufSize)
	return set
}

func NewServer(tbs AppTBs) {
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
	app := NewApp(tbs)
	defer app.Close()
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
