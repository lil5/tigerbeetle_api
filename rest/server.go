package rest

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lil5/tigerbeetle_api/config"
	"github.com/lil5/tigerbeetle_api/grpc"
	"github.com/lil5/tigerbeetle_api/metrics"
	"github.com/prometheus/client_golang/prometheus"

	metrics_prometheus "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

func NewServer() {
	if config.Config.Mode != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	r, app := Router()
	defer app.Close()
	slog.Info("Rest server listening at", "host", config.Config.Host, "port", config.Config.Port)
	defer slog.Info("Server exiting")

	mdlw := middleware.New(middleware.Config{
		Recorder: metrics_prometheus.NewRecorder(metrics_prometheus.Config{
			Registry: prometheus.DefaultRegisterer,
		}),
	})
	r.Use(ginmiddleware.Handler("", mdlw))

	prometheusDeferClose := metrics.Register(config.Config.PrometheusAddr)
	defer prometheusDeferClose()

	addr := fmt.Sprintf("%s:%s", config.Config.Host, config.Config.Port)
	if config.Config.OnlyIpv4 {
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

func Router() (*gin.Engine, *grpc.App) {
	s := grpc.NewApp()
	var r *gin.Engine
	if config.Config.Mode == "development" {
		r = gin.Default()
	} else {
		r = gin.New()
	}
	r.GET("/id", grpcHandle(s.GetID))
	r.GET("/ping", ping)
	r.POST("/accounts/create", grpcHandle(s.CreateAccounts))
	r.POST("/transfers/create", grpcHandle(s.CreateTransfers))
	r.POST("/accounts/lookup", grpcHandle(s.LookupAccounts))
	r.POST("/transfers/lookup", grpcHandle(s.LookupTransfers))
	r.POST("/account/transfers", grpcHandle(s.GetAccountTransfers))
	r.POST("/account/balances", grpcHandle(s.GetAccountBalances))
	r.POST("/transfers/query", grpcHandle(s.QueryTransfers))
	r.POST("/accounts/query", grpcHandle(s.QueryAccounts))
	return r, s
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func grpcHandle[In any, Out any](f func(ctx context.Context, in *In) (out *Out, err error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in In
		if err := c.ShouldBindBodyWithJSON(&in); err != nil {
			errStr := err.Error()
			slog.Warn(errStr)
			c.String(http.StatusBadRequest, errStr)
			return
		}
		out, err := f(c.Request.Context(), &in)
		if err != nil {
			errStr := err.Error()
			slog.Error(errStr)
			c.String(http.StatusInternalServerError, errStr)
			return
		}

		c.JSON(http.StatusOK, out)
	}
}
