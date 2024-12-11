package rest

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lil5/tigerbeetle_api/grpc"
)

func NewServer(tbs grpc.AppTBs) {
	if os.Getenv("MODE") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	r, app := Router(tbs)
	defer app.AppTBs.Close()
	slog.Info("Rest server listening at", "host", os.Getenv("HOST"), "port", os.Getenv("PORT"))
	defer slog.Info("Server exiting")

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
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

func Router(tbs grpc.AppTBs) (*gin.Engine, *grpc.App) {
	s := grpc.NewApp(tbs)
	r := gin.Default()
	r.GET("/id", grpcHandle(s.GetID))
	r.GET("/ping", ping)
	r.POST("/accounts/create", grpcHandle(s.CreateAccounts))
	r.POST("/transfers/create", grpcHandle(s.CreateTransfers))
	r.POST("/accounts/lookup", grpcHandle(s.LookupAccounts))
	r.POST("/transfers/lookup", grpcHandle(s.LookupTransfers))
	r.POST("/account/transfers", grpcHandle(s.GetAccountTransfers))
	r.POST("/account/balances", grpcHandle(s.GetAccountBalances))
	return r, s
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func grpcHandle[In any, Out any](f func(ctx context.Context, in *In) (out *Out, err error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in In
		if err := c.ShouldBindBodyWithJSON(&in); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		out, err := f(c.Request.Context(), &in)
		if err != nil {
			fmt.Fprint(c.Writer, err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, out)
	}
}
