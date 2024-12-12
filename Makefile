default:
	@grep '^[^#[:space:].].*:' Makefile

start:
	go run .

.PHONY: build
build:
	go build .

buildexec:
	go build .
	./tigerbeetle_api

release: release_darwin_arm64 release_darwin_arm64 release_linux_amd64 release_linux_arm64
release_linux_amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o build/tigerbeetle_api_linux_amd64 main.go
release_linux_amd64_cross_compile:
	CC="zig cc -target x86_64-linux-musl" CXX="zig c++ -target x86_64-linux-musl" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o build/tigerbeetle_api_linux_amd64 main.go
release_darwin_arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o build/tigerbeetle_api_darwin_arm64 main.go
release_linux_arm64_cross_compile:
	CC="zig cc -target aarch64-linux" CXX="zig c++ -target aarch64-linux" CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o build/tigerbeetle_api_linux_arm64 main.go

release_run_darwin_arm64:
	./build/tigerbeetle_api_darwin_arm64
release_run_linux_amd64:
	./build/tigerbeetle_api_linux_amd64
release_run_linux_arm64:
	./build/tigerbeetle_api_linux_arm64

docker_start:
	docker compose up -d
docker_stop:
	docker compose stop
docker_setup:
	docker compose pull && docker compose run tigerbeetle format --cluster=0 --replica=0 --replica-count=1 /data/0_0.tigerbeetle
docker_remove:
	docker compose down -v --remove-orphans

e2e_test:
	docker compose up -d
	go test --tags e2e ./...
e2e_benchmark:
	docker compose up -d
	go build -o tigerbeetle_api .
	go test --tags e2e --bench . --count 5 --benchmem --run=^# ./benchmark_e2e_test.go

proto_setup_mac:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	brew install protobuf
proto_gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional proto/*.proto

ghz_benchmark_id:
	ghz --insecure --call proto.TigerBeetle.GetID -d {} -n 500000 --concurrency 20000 --connections=32 localhost:50051
ghz_benchmark_transfers:
	ghz --insecure --call proto.TigerBeetle.CreateTransfers --data-file transfers.json -n 500000 --concurrency 20000 --connections=64 localhost:50051