package main

import (
	"github.com/joho/godotenv"
	"github.com/lil5/tigerbeetle_api/grpc"
	"github.com/lil5/tigerbeetle_api/rest"
)

func main() {
	godotenv.Load()
	grpc.NewConfig()
	if grpc.Config.UseGrpc {
		grpc.NewServer()
	} else {
		rest.NewServer()
	}
}
