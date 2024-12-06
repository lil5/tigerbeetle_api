package rest

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
)

func NewServer(tb tigerbeetle_go.Client) {
	s := Server{TB: tb}
	if os.Getenv("MODE") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.GET("/id", s.GetID)
	r.GET("/ping", ping)
	r.POST("/accounts/create", s.CreateAccounts)
	r.POST("/transfers/create", s.CreateTransfers)
	r.POST("/accounts/lookup", s.LookupAccounts)
	r.POST("/transfers/lookup", s.LookupTransfers)
	r.POST("/account/transfers", s.GetAccountTransfers)
	r.POST("/account/balances", s.GetAccountBalances)

	slog.Info("server listening at", "host", os.Getenv("HOST"), "port", os.Getenv("PORT"))
	defer slog.Info("server exiting")

	addr := fmt.Sprintf("%s:%d", os.Getenv("HOST"), os.Getenv("PORT"))
	if os.Getenv("ONLY_IPV4") == "true" {
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

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
