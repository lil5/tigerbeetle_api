package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/lil5/tigerbeetle_api/app"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func main() {
	// Parse flags
	fFile := flag.String("c", "config.yml", "Override config file")
	flag.Parse()

	fpath, name, ext := app.ReadFlag(lo.FromPtrOr(fFile, "config.yml"))
	slog.Info(fmt.Sprintf("config file: %s/%s.%s", fpath, name, ext))

	// Parse command line arguments
	viper.SetConfigName(name)
	viper.SetConfigType(ext)
	viper.AddConfigPath(fpath)
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

	slog.Info("Connecting to tigerbeetle cluster", "addresses:", strings.Join(tbAddresses, ", "))

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
