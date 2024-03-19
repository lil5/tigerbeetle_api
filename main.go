package main

import (
	"fmt"
	"log/slog"
	"os"
	"tigerbeetle_grpc/app"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

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

	// Create rest server
	s := app.Server{TB: tb}
	r := gin.New()
	r.GET("/id", s.GetID)
	r.POST("/accounts/create", s.CreateAccounts)
	r.POST("/transfers/create", s.CreateTransfers)
	r.POST("/accounts/lookup", s.LookupAccounts)
	r.POST("/transfers/lookup", s.LookupTransfers)
	r.POST("/account/transfers", s.GetAccountTransfers)
	r.POST("/account/balances", s.GetAccountBalances)

	slog.Info("server listening at", "host", host, "port", port)
	defer slog.Info("server exiting")
	r.Run(fmt.Sprintf("%s:%d", host, port))
}
