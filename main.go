package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/lil5/tigerbeetle_api/grpc"
	"github.com/lil5/tigerbeetle_api/rest"
)

func main() {
	godotenv.Load()
	if ok := grpc.NewConfig(); !ok {
		os.Exit(1)
	}
	if grpc.Config.UseGrpc {
		grpc.NewServer()
	} else {
		rest.NewServer()
	}
}
