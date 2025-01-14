package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

var Config config

type config struct {
	Host string
	Port string

	OnlyIpv4 bool
	Mode     string

	TbClusterID uint64
	TbAddresses []string

	UseGrpc          bool
	GrpcHealthServer bool
	GrpcReflection   bool
	GrpcHighScale    bool

	IsBuffered    bool
	BufferSize    int
	BufferDelay   time.Duration
	BufferCluster int

	IsDryRun bool

	PrometheusAddr string
}

func NewConfig() (ok bool) {
	useGrpc := os.Getenv("USE_GRPC") == "true"

	if host := os.Getenv("HOST"); host == "" {
		os.Setenv("HOST", "0.0.0.0")
	}

	if port, _ := strconv.Atoi(os.Getenv("PORT")); port == 0 {
		if useGrpc {
			os.Setenv("PORT", "50051")
		} else {
			os.Setenv("PORT", "8000")
		}
	}

	tbAddressesArr := os.Getenv("TB_ADDRESSES")
	if tbAddressesArr == "" {
		slog.Error("tb_addresses is empty")
		return false
	}
	tbAddresses := strings.Split(tbAddressesArr, ",")

	tbClusterId, _ := strconv.ParseUint(os.Getenv("TB_CLUSTER_ID"), 10, 64)

	isBuffered := os.Getenv("IS_BUFFERED") == "true"
	bufferSize := 0
	bufferCluster := 0
	var bufferDelay time.Duration
	if isBuffered {
		bufferSize, _ = strconv.Atoi(os.Getenv("BUFFER_SIZE"))
		if bufferSize == 0 {
			bufferSize = 1
		}
		var err error
		bufferDelay, err = time.ParseDuration(os.Getenv("BUFFER_DELAY"))
		if err != nil {
			slog.Error("BUFFER_DELAY is invalid duration", "error", err)
			return false
		}
		bufferCluster, _ = strconv.Atoi(os.Getenv("BUFFER_CLUSTER"))
		if bufferCluster == 0 {
			bufferCluster = 1
		}
	}

	prometheusAddr := os.Getenv("PROMETHEUS_ADDR")
	if prometheusAddr == "" {
		prometheusAddr = ":9323"
	}

	Config = config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),

		OnlyIpv4: os.Getenv("ONLY_IPV4") == "true",
		Mode:     os.Getenv("MODE"),

		TbClusterID: tbClusterId,
		TbAddresses: tbAddresses,

		UseGrpc:          useGrpc,
		GrpcHealthServer: os.Getenv("GRPC_HEALTH_SERVER") == "true",
		GrpcReflection:   os.Getenv("GRPC_REFLECTION") == "true",
		GrpcHighScale:    os.Getenv("GRPC_HIGH_SCALE") == "true",

		IsBuffered:    isBuffered,
		BufferSize:    bufferSize,
		BufferDelay:   bufferDelay,
		BufferCluster: bufferCluster,

		IsDryRun: os.Getenv("IS_DRY_RUN") == "true",

		PrometheusAddr: prometheusAddr,
	}

	slog.Info(fmt.Sprintf("%+v", Config))
	slog.Info("Config loaded", "version", version)
	return true
}
