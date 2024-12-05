//go:build e2e

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

const (
	BENCH_LEDGER        = 98
	BENCH_TB_ADDRESSES  = "127.0.0.1:3033"
	BENCH_TB_CLUSTER_ID = 0
	BENCH_TB_API_PORT   = "8080"
	BENCH_TB_API_ADDR   = "127.0.0.1:8080"
)

func BenchmarkTest(b *testing.B) {
	cmd := exec.Command("./tigerbeetle_api")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	cmd.Env = []string{"PORT=" + BENCH_TB_API_PORT, "TB_ADDRESSES=" + BENCH_TB_ADDRESSES}
	cmd.Start()
	time.Sleep(2 * time.Second)
	b.ResetTimer()
	b.Run("CreateAccounts", func(b *testing.B) {
		for range b.N {
			httpRequest(b, http.MethodPost, "http://127.0.0.1:8080/accounts/create", gin.H{
				"accounts": []gin.H{{
					"user_data_128":   nil,
					"user_data_64":    nil,
					"user_data_32":    nil,
					"id":              "",
					"debits_pending":  0,
					"debits_posted":   0,
					"credits_pending": 0,
					"credits_posted":  0,
					"ledger":          BENCH_LEDGER,
					"code":            1,
					"flags": gin.H{
						"linked":                         false,
						"credits_must_not_exceed_debits": false,
						"debits_must_not_exceed_credits": false,
						"history":                        true,
					},
					"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
				}}})
		}
	})

	cmd.Process.Kill()
}

func BenchmarkComparison(b *testing.B) {
	tb, _ := tigerbeetle_go.NewClient(types.ToUint128(uint64(BENCH_TB_CLUSTER_ID)), []string{BENCH_TB_ADDRESSES})

	b.ResetTimer()
	b.Run("CreateAccounts", func(b *testing.B) {
		for range b.N {
			f := types.AccountFlags{
				History: true,
			}
			tb.CreateAccounts([]types.Account{{
				ID:             types.ID(),
				DebitsPending:  [16]uint8{0},
				DebitsPosted:   [16]uint8{0},
				CreditsPending: [16]uint8{0},
				CreditsPosted:  [16]uint8{0},
				UserData128:    [16]uint8{},
				UserData64:     0,
				UserData32:     0,
				Reserved:       0,
				Ledger:         BENCH_LEDGER,
				Code:           1,
				Flags:          f.ToUint16(),
				Timestamp:      uint64(time.Now().UnixNano()),
			}})
		}
	})
	tb.Close()
}

// utility functions
// ----------------------------------------------------------------
func httpRequest(b *testing.B, method, url string, reqJson gin.H) {
	reqBody, _ := json.Marshal(reqJson)
	// if err != nil {
	// 	b.Error("Failed to marshal", err)
	// 	return
	// }
	r := bytes.NewReader(reqBody)
	req, _ := http.NewRequest(method, url, r)
	// if err != nil {
	// 	b.Error("Failed to create request", err)
	// }
	http.DefaultClient.Do(req)
	// if err != nil {
	// 	b.Error("Unable to send request", err)
	// 	return
	// }
	// if res.StatusCode != http.StatusOK {
	// 	resBody, _ := io.ReadAll(res.Body)
	// 	b.Error(string(resBody))
	// }
}
