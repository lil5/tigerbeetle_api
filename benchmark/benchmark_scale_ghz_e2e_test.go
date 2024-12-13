package benchmark

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/bojand/ghz/runner"
	"github.com/stretchr/testify/assert"
)

// Run `make build` and comment out any values in your `.env` file before running this test
func TestBenchmarkGhz(t *testing.T) {
	clientCount := 4
	f := func(env []string, port string) (*runner.Report, func()) {
		// t.Helper()
		close := tbApiStart(env)
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
			runner.WithTotalRequests(500_000),
			// --concurrency 20000
			runner.WithConcurrency(20000),
			// --connections=32
			runner.WithConnections(100),
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

		return report, close
	}

	var rpsBufMulti int64

	t.Run("one client", func(t *testing.T) {
		report, close := f([]string{
			"PORT=50052",
			"TB_ADDRESSES=3033",
			"TB_CLUSTER_ID=0",
			"USE_GRPC=true",
			"GRPC_REFLECTION=true",
			"IS_BUFFERED=true",
			fmt.Sprintf("BUFFER_CLUSTER=%d", clientCount),
			"BUFFER_SIZE=100",
			"BUFFER_DELAY=100ms",
			"MODE=production",
		}, "50052")
		defer close()

		rpsBufMulti = int64(report.Rps)
		assert.GreaterOrEqual(t, rpsBufMulti, int64(16_000), "m2 traffic maximum")
		assert.GreaterOrEqual(t, rpsBufMulti, int64(200_000), "max traffic requirements")
		t.Logf("report rps: %d", int64(report.Rps))
	})

	time.Sleep(1 * time.Second)

	t.Run("no buffer", func(t *testing.T) {
		report, close := f([]string{
			"PORT=50053",
			"TB_ADDRESSES=3033",
			"TB_CLUSTER_ID=0",
			"USE_GRPC=true",
			"GRPC_REFLECTION=true",
			"CLIENT_COUNT=1",
			"MODE=production",
			"IS_BUFFERED=false",
		}, "50053")
		defer close()

		rpsNoBuf := int64(report.Rps)
		t.Log("one client", rpsBufMulti)
		t.Log("multi client", rpsNoBuf)
		assert.GreaterOrEqual(t, rpsNoBuf, rpsBufMulti, "1 gt multiple clients")
		// assert.GreaterOrEqual(t, int64(report.Rps), int64(200_000), "max traffic requirements")
		t.Logf("report rps: %d", int64(report.Rps))
	})
}

func tbApiStart(env []string) func() {
	cmd := exec.Command("../tigerbeetle_api")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stdout

	cmd.Env = env
	cmd.Start()
	time.Sleep(200 * time.Millisecond) // sleep waiting for grpc server to start
	return func() {
		cmd.Process.Kill()
	}
}
