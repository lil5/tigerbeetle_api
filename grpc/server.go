package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/lil5/tigerbeetle_api/proto"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"google.golang.org/grpc"
)

func NewServer(tb tigerbeetle_go.Client) {
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}
	proto.RegisterTigerBeetleServer(s, &Server{TB: tb})
	slog.Info("server listening at", "address", lis.Addr())
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve:", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
