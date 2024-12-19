package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/lil5/tigerbeetle_api/config"
	"github.com/lil5/tigerbeetle_api/grpc"
	"github.com/lil5/tigerbeetle_api/rest"
)

func main() {
	godotenv.Load()
	if ok := config.NewConfig(); !ok {
		os.Exit(1)
	}
	// log add version
	if config.Config.UseGrpc {
		grpc.NewServer()
	} else {
		rest.NewServer()
	}
}
