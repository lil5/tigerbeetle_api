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
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lil5/tigerbeetle_api/grpc"
	"github.com/lil5/tigerbeetle_api/rest"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
)

const (
	LEDGER        = 99
	TB_ADDRESSES  = "127.0.0.1:3033"
	TB_CLUSTER_ID = "0"
)

type MyTestSuite struct {
	suite.Suite
	tb     tigerbeetle_go.Client
	router *gin.Engine
	app    *grpc.App
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

	os.Setenv("TB_ADDRESSES", TB_ADDRESSES)
	os.Setenv("TB_CLUSTER_ID", TB_CLUSTER_ID)
	if ok := grpc.NewConfig(); !ok {
		s.FailNow("SetupSuite failed to initialize creating configuration")
		return
	}

	gin.SetMode(gin.TestMode)
	s.router, s.app = rest.Router()
}

// run once, after test suite methods
func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")

	// stop the tb client
	s.app.Close()
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
	_, resultFunc := MockGinContext(s.router, http.MethodGet, "/id", &gin.H{})
	result := resultFunc()
	if result.Response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected status code: %d, body: %s", result.Response.StatusCode, result.Body)
	}
	json := result.BodyJSON()

	return json["id"].(string), nil
}

func (s *MyTestSuite) TestCalls() {
	var accountID1 string
	var accountID2 string
	s.Run("GetID", func() {
		accountID1, _ = s.RunGetID()
		accountID2, _ = s.RunGetID()
	})
	s.Run("CreateAccounts", func() {
		creditsMustNotExceedDebits := false
		debitsMustNotExceedCredits := false
		_, resultFunc := MockGinContext(s.router, http.MethodPost, "/accounts/create", &gin.H{
			"accounts": []gin.H{
				{
					"id":     accountID1,
					"ledger": LEDGER,
					"code":   1,
					"flags": gin.H{
						"credits_must_not_exceed_debits": creditsMustNotExceedDebits,
						"debits_must_not_exceed_credits": debitsMustNotExceedCredits,
						"history":                        true,
					},
					"timestamp": time.Now().Unix(),
				}, {
					"user_data_128":   nil,
					"user_data_64":    nil,
					"user_data_32":    nil,
					"id":              accountID2,
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
					"timestamp": time.Now().Unix(),
				},
			},
		})
		result := resultFunc()
		s.Equal(http.StatusOK, result.Response.StatusCode, result.Body)
		json := gjson.Parse(result.Body)

		s.Len(json.Get("results").Array(), 0, result.Body)
	})

	s.Run("LookupAccounts empty", func() {
		_, resultFunc := MockGinContext(s.router, http.MethodPost, "/accounts/lookup", &gin.H{
			"account_ids": []string{accountID1, accountID2},
		})
		result := resultFunc()

		s.Equal(http.StatusOK, result.Response.StatusCode)
		json := gjson.Parse(result.Body)
		s.Len(json.Get("accounts").Array(), 2)
		s.Equal(0.0, json.Get("accounts.0.debits_posted").Num)
		s.Equal(0.0, json.Get("accounts.0.credits_posted").Num)
		s.Equal(0.0, json.Get("accounts.1.debits_posted").Num)
		s.Equal(0.0, json.Get("accounts.1.credits_posted").Num)
	})

	s.Run("CreateTransfer", func() {
		slog.Info("Creating transfer, take out 10 from account 2 and put 10 in account 1")
		id1, _ := s.RunGetID()
		id2, _ := s.RunGetID()
		_, resultFunc := MockGinContext(s.router, http.MethodPost, "/transfers/create", &gin.H{
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
					"timestamp": time.Now().Unix(),
				},
				{
					"user_data_128":     nil,
					"user_data_64":      nil,
					"user_data_32":      nil,
					"id":                id2,
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
					"timestamp": time.Now().Unix(),
				},
			},
		})
		result := resultFunc()
		s.Equal(http.StatusOK, result.Response.StatusCode)
		s.Len(gjson.Get(result.Body, "results").Array(), 0)
	})

	s.Run("LookupAccounts after 1 transfer", func() {
		_, resultFunc := MockGinContext(s.router, http.MethodPost, "/accounts/lookup", &gin.H{
			"account_ids": []string{accountID1, accountID2},
		})
		result := resultFunc()

		s.Equal(http.StatusOK, result.Response.StatusCode)

		json := gjson.Parse(result.Body)
		s.Len(json.Get("accounts").Array(), 2)
		s.Equal(10.0, json.Get("accounts.0.debits_posted").Num)
		s.Equal(0.0, json.Get("accounts.0.credits_posted").Num)
		s.Equal(0.0, json.Get("accounts.1.debits_posted").Num)
		s.Equal(10.0, json.Get("accounts.1.credits_posted").Num)
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

func MockGinContext(router *gin.Engine, method string, url string, bodyJSON *gin.H) (*httptest.ResponseRecorder, func() mockGinContextResponse) {
	body := bytes.NewBuffer([]byte{})
	if bodyJSON != nil {
		json_data, _ := json.Marshal(bodyJSON)
		body = bytes.NewBuffer(json_data)
	}

	r := httptest.NewRequest(method, url, body)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// ro.ServeHTTP(rr, c.Request)
	resultFunc := func() mockGinContextResponse {
		body := w.Body.String()
		res := w.Result()
		return mockGinContextResponse{
			Response: res,
			Body:     body,
		}
	}

	return w, resultFunc
}
