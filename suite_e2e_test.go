//go:build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lil5/tigerbeetle_api/app"
	"github.com/stretchr/testify/suite"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

const (
	LEDGER        = 99
	TB_ADDRESSES  = "127.0.0.1:3033"
	TB_CLUSTER_ID = 0
)

type MyTestSuite struct {
	suite.Suite
	server app.Server
}

// listen for 'go test' command --> run test methods
func TestMyTestSuite(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("Starting test", "TB_ADDRESSES", TB_ADDRESSES, "TB_CLUSTER_ID", TB_CLUSTER_ID)

	suite.Run(t, new(MyTestSuite))
}

// run once, before test suite methods
func (s *MyTestSuite) SetupSuite() {
	slog.Info("SetupSuite()")

	// setup tigerbeetle server
	err := exec.Command("/bin/bash", "-c", "docker compose up -d").Run()
	if err != nil {
		log.Fatal(err)
	}

	// connect to tigerbeetle server
	tb, err := tigerbeetle_go.NewClient(types.ToUint128(uint64(TB_CLUSTER_ID)), strings.Split(TB_ADDRESSES, ","))
	if err != nil {
		slog.Error("unable to connect to tigerbeetle:", "err", err)
		os.Exit(1)
	}

	s.server = app.Server{TB: tb}

	gin.SetMode(gin.TestMode)
}

// run once, after test suite methods
func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")

	// stop the tb client
	s.server.TB.Close()
	// stop the tb server
	err := exec.Command("/bin/bash", "-c", "docker compose down").Run()
	if err != nil {
		log.Fatal(err)
	}

}

func (s *MyTestSuite) TestGetID() {
	id, err := s.RunGetID()
	s.Nil(err, "body: %s, err: %v", id, err)
	s.Len(id, 31)
}

func (s *MyTestSuite) RunGetID() (string, error) {
	c, resultFunc := MockGinContext(http.MethodPost, "/", nil)
	s.server.GetID(c)
	result := resultFunc()
	if result.Response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected status code: %d, body: %s", result.Response.StatusCode, result.Body)
	}
	json := result.BodyJSON()

	return json["id"].(string), nil
}

func (s *MyTestSuite) TestCalls() {
	accountID1, _ := s.RunGetID()
	var accountID2 string
	s.Run("CreateAccounts", func() {
		creditsMustNotExceedDebits := false
		debitsMustNotExceedCredits := false
		c, resultFunc := MockGinContext(http.MethodPost, "/", &gin.H{
			"accounts": []gin.H{
				{
					"user_data_128":   nil,
					"user_data_64":    nil,
					"user_data_32":    nil,
					"id":              accountID1,
					"debits_pending":  0,
					"debits_posted":   0,
					"credits_pending": 0,
					"credits_posted":  0,
					"ledger":          LEDGER,
					"code":            1,
					"flags": gin.H{
						"linked":                         false,
						"credits_must_not_exceed_debits": creditsMustNotExceedDebits,
						"debits_must_not_exceed_credits": debitsMustNotExceedCredits,
						"history":                        true,
					},
					"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
				}, {
					"user_data_128": nil,
					"user_data_64":  nil,
					"user_data_32":  nil,
					// "id":              "",
					"debits_pending":  0,
					"debits_posted":   0,
					"credits_pending": 0,
					"credits_posted":  0,
					"ledger":          LEDGER,
					"code":            1,
					"flags": gin.H{
						"linked":                         false,
						"credits_must_not_exceed_debits": creditsMustNotExceedDebits,
						"debits_must_not_exceed_credits": debitsMustNotExceedCredits,
						"history":                        true,
					},
					"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
				},
			},
		})
		s.server.CreateAccounts(c)
		result := resultFunc()
		s.Equal(http.StatusOK, result.Response.StatusCode, result.Body)
		json := result.BodyJSON()
		resAccountIDs := json["account_ids"].([]any)
		s.Equal(accountID1, resAccountIDs[0].(string))
		accountID2 = resAccountIDs[1].(string)
	})

	s.Run("LookupAccounts empty", func() {
		c, resultFunc := MockGinContext(http.MethodPost, "/", &gin.H{
			"account_ids": []string{accountID1, accountID2},
		})
		s.server.LookupAccounts(c)
		result := resultFunc()

		json := result.BodyJSON()
		s.Equal(http.StatusOK, result.Response.StatusCode)
		jsonAccounts := json["accounts"].([]any)
		s.Len(jsonAccounts, 2)
		s.Equal(0.0, (jsonAccounts[0].(map[string]any))["debits_posted"])
		s.Equal(0.0, (jsonAccounts[0].(map[string]any))["credits_posted"])
		s.Equal(0.0, (jsonAccounts[1].(map[string]any))["debits_posted"])
		s.Equal(0.0, (jsonAccounts[1].(map[string]any))["credits_posted"])
	})

	s.Run("CreateTransfer", func() {
		slog.Info("Creating transfer, take out 10 from account 2 and put 10 in account 1")
		id1, _ := s.RunGetID()
		c, resultFunc := MockGinContext(http.MethodPost, "/", &gin.H{
			"transfers": []gin.H{
				{
					"user_data_128":     nil,
					"user_data_64":      nil,
					"user_data_32":      nil,
					"id":                id1,
					"debit_account_id":  accountID1,
					"credit_account_id": accountID2,
					"amount":            5,
					"pending_id":        nil,
					"ledger":            LEDGER,
					"code":              1,
					"transfer_flags": gin.H{
						"linked":                false,
						"pending":               false,
						"post_pending_transfer": false,
						"void_pending_transfer": false,
						"balancing_debit":       false,
						"balancing_credit":      false,
					},
					"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
				},
				{
					"user_data_128": nil,
					"user_data_64":  nil,
					"user_data_32":  nil,
					// "id":                "",
					"debit_account_id":  accountID1,
					"credit_account_id": accountID2,
					"amount":            5,
					"pending_id":        nil,
					"ledger":            LEDGER,
					"code":              1,
					"transfer_flags": gin.H{
						"linked":                false,
						"pending":               false,
						"post_pending_transfer": false,
						"void_pending_transfer": false,
						"balancing_debit":       false,
						"balancing_credit":      false,
					},
					"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
				},
			},
		})
		s.server.CreateTransfers(c)
		result := resultFunc()
		s.Equal(http.StatusOK, result.Response.StatusCode)
		json := result.BodyJSON()
		resTransferIDs := json["transfer_ids"].([]any)
		s.Equal(id1, resTransferIDs[0].(string))
	})

	s.Run("LookupAccounts after 1 transfer", func() {
		c, resultFunc := MockGinContext(http.MethodPost, "/", &gin.H{
			"account_ids": []string{accountID1, accountID2},
		})
		s.server.LookupAccounts(c)
		result := resultFunc()

		json := result.BodyJSON()
		s.Equal(http.StatusOK, result.Response.StatusCode)
		jsonAccounts := json["accounts"].([]any)
		s.Len(jsonAccounts, 2)

		s.Equal(10.0, (jsonAccounts[0].(map[string]any))["debits_posted"])
		s.Equal(0.0, (jsonAccounts[0].(map[string]any))["credits_posted"])
		s.Equal(0.0, (jsonAccounts[1].(map[string]any))["debits_posted"])
		s.Equal(10.0, (jsonAccounts[1].(map[string]any))["credits_posted"])
	})
}

// utility functions
// ----------------------------------------------------------------
type mockGinContextResponse struct {
	Response *http.Response
	Body     string
}

func (r mockGinContextResponse) BodyJSON() gin.H {
	body := gin.H{}
	json.Unmarshal([]byte(r.Body), &body)

	return body
}

func MockGinContext(method string, url string, bodyJSON *gin.H) (*gin.Context, func() mockGinContextResponse) {
	body := bytes.NewBuffer([]byte{})
	if bodyJSON != nil {
		json_data, _ := json.Marshal(bodyJSON)
		body = bytes.NewBuffer(json_data)
	}

	r := httptest.NewRequest(method, url, body)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = r

	// ro.ServeHTTP(rr, c.Request)
	resultFunc := func() mockGinContextResponse {
		body := rr.Body.String()
		res := rr.Result()
		return mockGinContextResponse{
			Response: res,
			Body:     body,
		}
	}

	return c, resultFunc
}
