package main

import (
	"os/exec"
	"testing"
	"time"

	"github.com/bojand/ghz/runner"
	"github.com/stretchr/testify/assert"
)

// Run `make build` and comment out any values in your `.env` file before running this test
func TestBenchmarkGhz(t *testing.T) {
	f := func(t *testing.T, env []string, port string) *runner.Report {
		// t.Helper()
		cmd := exec.Command("../tigerbeetle_api")
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stdout

		cmd.Env = env
		cmd.Start()
		t.Cleanup(func() {
			cmd.Process.Kill()
		})
		time.Sleep(500 * time.Millisecond)

		report, err := runner.Run(
			// --call proto.TigerBeetle.CreateTransfers
			"proto.TigerBeetle.CreateTransfers",
			// localhost:50051
			"localhost:"+port,
			// --call proto.TigerBeetle.CreateTransfers
			runner.WithProtoFile("../proto/tigerbeetle.proto", []string{}),
			// --insecure
			runner.WithInsecure(true),
			// --data-file transfers.json
			runner.WithDataFromFile("../transfers.json"),
			// -n 500000
			runner.WithTotalRequests(500000),
			// --concurrency 20000
			runner.WithConcurrency(20000),
			// --connections=32
			runner.WithConnections(64),
		)
		if err != nil {
			t.Error(err)
		}

		// printer := printer.ReportPrinter{
		// 	Out:    os.Stdout,
		// 	Report: report,
		// }
		// printer.Print("pretty")

		assert.Equal(t, runner.ReasonNormalEnd, report.EndReason)

		return report
	}

	var noBufferTps int64

	t.Run("one client", func(t *testing.T) {
		report := f(t, []string{
			"PORT=50052",
			"TB_ADDRESSES=3033",
			"TB_CLUSTER_ID=0",
			"USE_GRPC=true",
			"GRPC_REFLECTION=true",
			"CLIENT_COUNT=1",
			"MODE=production",
		}, "50052")

		noBufferTps = int64(report.Rps)
		assert.GreaterOrEqual(t, noBufferTps, int64(16_000), "m2 traffic maximum")
		// assert.GreaterOrEqual(t, noBufferTps, int64(200_000), "max traffic requirements")
		t.Logf("report rps: %d", int64(report.Rps))
	})

	t.Run("multi client", func(t *testing.T) {
		report := f(t, []string{
			"PORT=50053",
			"TB_ADDRESSES=3033",
			"TB_CLUSTER_ID=0",
			"USE_GRPC=true",
			"GRPC_REFLECTION=true",
			"CLIENT_COUNT=20",
			"BUFFER_SIZE=20",
			"BUFFER_DELAY=50ms",
			"MODE=production",
		}, "50053")

		assert.GreaterOrEqual(t, int64(report.Rps), noBufferTps, "1 gt multiple clients")
		// assert.GreaterOrEqual(t, int64(report.Rps), int64(200_000), "max traffic requirements")
		t.Logf("report rps: %d", int64(report.Rps))
	})
}
