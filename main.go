package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jinzhu/configor"
	"github.com/lil5/tigerbeetle_api/app"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type Config struct {
	Port             int    `default:"50051" yaml:"port" env:"PORT"`
	Host             string `default:"0.0.0.0" yaml:"host" env:"HOST"`
	TbClusterId      int    `default:"0" yaml:"tb_cluster_id" env:"TB_CLUSTER_ID"`
	TbAddresses      string `required:"true" yaml:"tb_addresses" env:"TB_ADDRESSES"`
	TbConcurrencyMax int    `default:"2" yaml:"tb_concurrency_max" env:"TB_CONCURRENCY_MAX"`
}

func main() {
	// Parse flags
	fFile := flag.String("c", "config.yml", "Override config file")
	flag.Parse()

	fpath, name, ext := app.ReadFlag(lo.FromPtrOr(fFile, "config.yml"))
	slog.Info(fmt.Sprintf("config file: %s/%s.%s", fpath, name, ext))

	config := &Config{}
	err := configor.Load(config, fpath)
	if err != nil {
		slog.Error("fatal error config file:", err)
		os.Exit(1)
	}
	if config.TbAddresses == "" {
		slog.Error("tb_addresses is empty")
		os.Exit(1)
	}
	tbAddresses := strings.Split(config.TbAddresses, ",")

	slog.Info("Connecting to tigerbeetle cluster", "addresses:", strings.Join(tbAddresses, ", "))

	// Connect to tigerbeetle
	tb, err := tigerbeetle_go.NewClient(types.ToUint128(uint64(config.TbClusterId)), tbAddresses, uint(config.TbConcurrencyMax))
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

	slog.Info("server listening at", "host", config.Host, "port", config.Port)
	defer slog.Info("server exiting")
	r.Run(fmt.Sprintf("%s:%d", config.Host, config.Port))
}
