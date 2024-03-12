package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"tigerbeetle_grpc/proto"

	"github.com/spf13/viper"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.

func main() {
	// Parse command line arguments
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("fatal error config file:", err)
		os.Exit(1)
	}
	viper.SetDefault("port", 50051)
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("tb_cluster_id", 0)
	viper.SetDefault("tb_addresses", []string{})
	viper.SetDefault("tb_concurrency_max", 2)

	port := viper.GetInt("port")
	host := viper.GetString("host")
	tbClusterId := viper.GetUint64("tb_cluster")
	tbAddresses := viper.GetStringSlice("tb_addresses")
	tbConcurrencyMax := viper.GetUint("tb_concurrency_max")

	if len(tbAddresses) == 0 {
		slog.Error("tb_addresses is empty")
		os.Exit(1)
	}

	// Connect to tigerbeetle
	tb, err := tigerbeetle_go.NewClient(types.ToUint128(tbClusterId), tbAddresses, tbConcurrencyMax)
	if err != nil {
		slog.Error("unable to connect to tigerbeetle:", err)
		os.Exit(1)
	}
	defer tb.Close()

	// Create grpc server
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		slog.Error("failed to listen:", err)
		os.Exit(1)
	}
	proto.RegisterTigerBeetleServer(s, &server{TB: tb})
	slog.Info("server listening at", "address", lis.Addr())
	if err := s.Serve(lis); err != nil {
		slog.Error("failed to serve:", err)
		os.Exit(1)
	}

	slog.Info("server exiting")
}
