default:
	@grep '^[^#[:space:].].*:' Makefile

copy-proto:
	cp ../proto/tiger* proto/

run-ghz:
	go test --run TestBenchmarkGhz .
run-tbclient:
	go test --bench BenchmarkTbClient --run=^a
run-grpcclient:
	go test --bench BenchmarkGrpcTest --run=^a

run-compare-tb-vs-grpc:
	go test -run='^$$' -bench BenchmarkTbClient -count 10 > old.log
	go test -run='^$$' -bench BenchmarkGrpcTest -count 10 > new.log
	benchstat old.log new.log
