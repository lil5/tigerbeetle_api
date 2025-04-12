//go:build e2e

package main

import (
	"bytes"
	"encoding/json"
	"io"
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
	BENCH_TB_API_PORT   = "8001"
	BENCH_TB_API_ADDR   = "http://127.0.0.1:8001"
)

func BenchmarkMockHttpTest(b *testing.B) {
	cmd := exec.Command("./tigerbeetle_api")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	cmd.Env = []string{"PORT=" + BENCH_TB_API_PORT, "TB_ADDRESSES=" + BENCH_TB_ADDRESSES}
	cmd.Start()
	b.Cleanup(func() {
		cmd.Process.Kill()
	})
	time.Sleep(2 * time.Second)

	account1ID := types.ID().String()
	account2ID := types.ID().String()

	httpRequest(b, http.MethodPost, BENCH_TB_API_ADDR+"/accounts/create", gin.H{
		"accounts": []gin.H{{
			"user_data_128":   nil,
			"user_data_64":    nil,
			"user_data_32":    nil,
			"id":              account1ID,
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
			"timestamp": time.Now().UTC().UnixNano(),
		}, {
			"user_data_128":   nil,
			"user_data_64":    nil,
			"user_data_32":    nil,
			"id":              account2ID,
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
			"timestamp": time.Now().UTC().UnixNano(),
		}}})

	b.ResetTimer()
	b.Run("CreateTransfer", func(b *testing.B) {
		for range b.N {
			httpRequest(b, http.MethodPost, BENCH_TB_API_ADDR+"/transfers/create", gin.H{
				"transfers": []gin.H{
					{
						"user_data_128":     nil,
						"user_data_64":      nil,
						"user_data_32":      nil,
						"id":                "",
						"debit_account_id":  account1ID,
						"credit_account_id": account2ID,
						"amount":            5,
						"pending_id":        nil,
						"ledger":            BENCH_LEDGER,
						"code":              1,
						"transfer_flags": gin.H{
							"linked":                false,
							"pending":               false,
							"post_pending_transfer": false,
							"void_pending_transfer": false,
							"balancing_debit":       false,
							"balancing_credit":      false,
						},
						"timestamp": time.Now().UTC().UnixNano(),
					},
				},
			})
		}
	})
}

func BenchmarkDirectTbClientComparison(b *testing.B) {
	tb, _ := tigerbeetle_go.NewClient(types.ToUint128(uint64(BENCH_TB_CLUSTER_ID)), []string{BENCH_TB_ADDRESSES})

	account1ID := types.ID()
	account2ID := types.ID()

	f := types.AccountFlags{
		History: true,
	}
	_, err := tb.CreateAccounts([]types.Account{{
		ID:             account1ID,
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
	},
		{
			ID:             account2ID,
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
		},
	})
	if err != nil {
		b.Error(err)
	}
	b.Cleanup(func() {
		tb.Close()
	})

	b.ResetTimer()
	b.Run("CreateTransfer", func(b *testing.B) {
		for range b.N {
			_, err := tb.CreateTransfers([]types.Transfer{{
				ID:              types.ID(),
				DebitAccountID:  account1ID,
				CreditAccountID: account2ID,
				Amount:          [16]uint8{15},
				PendingID:       [16]uint8{},
				UserData128:     [16]uint8{},
				UserData64:      0,
				UserData32:      0,
				Timeout:         0,
				Ledger:          BENCH_LEDGER,
				Code:            1,
				Flags:           0,
				Timestamp:       uint64(time.Now().UnixNano()),
			}})
			if err != nil {
				b.Error(err)
			}
		}
	})
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
	res, _ := http.DefaultClient.Do(req)
	// if err != nil {
	// 	b.Error("Unable to send request", err)
	// 	return
	// }
	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)
		b.Error(res.StatusCode, string(resBody))
	}
}
