package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lil5/tigerbeetle_api/app"

	"github.com/gin-gonic/gin"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

const (
	defaultPort = 8000
	defaultHost = "0.0.0.0"
)

func main() {
	godotenv.Load()

	isDev := os.Getenv("MODE") == "development"
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = defaultPort
	}
	tbClusterId, _ := strconv.Atoi(os.Getenv("TB_CLUSTER_ID"))
	host := os.Getenv("HOST")
	if host == "" {
		host = defaultHost
	}
	tbAddressesArr := os.Getenv("TB_ADDRESSES")
	onlyIpv4 := os.Getenv("ONLY_IPV4") == "true"

	if tbAddressesArr == "" {
		slog.Error("tb_addresses is empty")
		os.Exit(1)
	}
	tbAddresses := strings.Split(tbAddressesArr, ",")

	slog.Info("Connecting to tigerbeetle cluster", "addresses", strings.Join(tbAddresses, ", "))

	// Connect to tigerbeetle
	tb, err := tigerbeetle_go.NewClient(types.ToUint128(uint64(tbClusterId)), tbAddresses)
	if err != nil {
		slog.Error("unable to connect to tigerbeetle", "err", err)
		os.Exit(1)
	}
	defer tb.Close()

	// Create rest server
	s := app.Server{TB: tb}
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.GET("/id", s.GetID)
	r.GET("/ping", Ping)
	r.POST("/accounts/create", s.CreateAccounts)
	r.POST("/transfers/create", s.CreateTransfers)
	r.POST("/accounts/lookup", s.LookupAccounts)
	r.POST("/transfers/lookup", s.LookupTransfers)
	r.POST("/account/transfers", s.GetAccountTransfers)
	r.POST("/account/balances", s.GetAccountBalances)

	slog.Info("server listening at", "host", host, "port", port)
	defer slog.Info("server exiting")

	addr := fmt.Sprintf("%s:%d", host, port)
	if onlyIpv4 {
		server := &http.Server{Handler: r}
		l, err := net.Listen("tcp4", addr)
		if err != nil {
			log.Fatal(err)
		}
		server.Serve(l)
	} else {
		r.Run(addr)
	}
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
