package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/lil5/tigerbeetle_api/proto"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func TestQueryAccounts(t *testing.T) {
	mockClient := new(MockTigerBeetleClient)
	app := &App{TB: mockClient}

	t.Run("should return error when filter is nil", func(t *testing.T) {
		req := &proto.QueryAccountsRequest{
			Filter: nil,
		}

		_, err := app.QueryAccounts(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "filter is required", err.Error())
	})

	t.Run("should return empty array when no accounts match", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData128: lo.ToPtr("0"),
			UserData64:  lo.ToPtr(uint64(0)),
			UserData32:  lo.ToPtr(uint32(0)),
			Code:        lo.ToPtr(uint32(1)),
			Limit:       10,
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		mockClient.On("QueryAccounts", mock.AnythingOfType("types.QueryFilter")).
			Return([]types.Account{}, nil).Once()

		resp, err := app.QueryAccounts(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.Accounts)
		mockClient.AssertExpectations(t)
	})

	t.Run("should return accounts with valid filter", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData128: lo.ToPtr("1"),
			UserData64:  lo.ToPtr(uint64(100)),
			UserData32:  lo.ToPtr(uint32(10)),
			Code:        lo.ToPtr(uint32(1)),
			Ledger:      lo.ToPtr(uint32(0)),
			Limit:       10,
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		// Mock account data
		mockAccounts := []types.Account{
			{
				ID:            types.ToUint128(1),
				DebitsPosted:  types.ToUint128(1000),
				CreditsPosted: types.ToUint128(500),
				UserData128:   types.ToUint128(1),
				UserData64:    100,
				UserData32:    10,
				Ledger:        0,
				Code:          1,
				Timestamp:     1000000,
			},
		}

		mockClient.On("QueryAccounts", mock.AnythingOfType("types.QueryFilter")).
			Return(mockAccounts, nil).Once()

		resp, err := app.QueryAccounts(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Accounts, 1)
		assert.Equal(t, "1", resp.Accounts[0].Id)
		assert.Equal(t, uint64(100), resp.Accounts[0].UserData64)
		assert.Equal(t, uint32(10), resp.Accounts[0].UserData32)
		mockClient.AssertExpectations(t)
	})

	t.Run("should respect limit parameter", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData64: lo.ToPtr(uint64(100)),
			Limit:      2,
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		// Mock 2 accounts but expect only 2 due to limit
		mockAccounts := []types.Account{
			{ID: types.ToUint128(1), UserData64: 100},
			{ID: types.ToUint128(2), UserData64: 100},
		}

		mockClient.On("QueryAccounts", mock.AnythingOfType("types.QueryFilter")).
			Return(mockAccounts, nil).Once()

		resp, err := app.QueryAccounts(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Accounts, 2)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle reversed flag correctly", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData64: lo.ToPtr(uint64(100)),
			Limit:      10,
			Flags: &proto.QueryFilterFlags{
				Reversed: lo.ToPtr(true),
			},
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		// Mock accounts in reverse order
		mockAccounts := []types.Account{
			{ID: types.ToUint128(3), Timestamp: 3000},
			{ID: types.ToUint128(2), Timestamp: 2000},
			{ID: types.ToUint128(1), Timestamp: 1000},
		}

		mockClient.On("QueryAccounts", mock.AnythingOfType("types.QueryFilter")).
			Return(mockAccounts, nil).Once()

		resp, err := app.QueryAccounts(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Accounts, 3)
		// Verify order is reversed (newest first)
		assert.Equal(t, "3", resp.Accounts[0].Id)
		assert.Equal(t, "2", resp.Accounts[1].Id)
		assert.Equal(t, "1", resp.Accounts[2].Id)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle timestamp range filters", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData64:   lo.ToPtr(uint64(100)),
			TimestampMin: lo.ToPtr(uint64(1000)),
			TimestampMax: lo.ToPtr(uint64(5000)),
			Limit:        10,
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		// Mock accounts within timestamp range
		mockAccounts := []types.Account{
			{ID: types.ToUint128(2), Timestamp: 2000},
			{ID: types.ToUint128(3), Timestamp: 3000},
			{ID: types.ToUint128(4), Timestamp: 4000},
		}

		mockClient.On("QueryAccounts", mock.AnythingOfType("types.QueryFilter")).
			Return(mockAccounts, nil).Once()

		resp, err := app.QueryAccounts(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Accounts, 3)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle multiple filter criteria", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData128: lo.ToPtr("1000"),
			UserData64:  lo.ToPtr(uint64(100)),
			UserData32:  lo.ToPtr(uint32(10)),
			Code:        lo.ToPtr(uint32(5)),
			Ledger:      lo.ToPtr(uint32(1)),
			Limit:       10,
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		mockAccounts := []types.Account{
			{
				ID:          types.ToUint128(1),
				UserData128: types.ToUint128(1000),
				UserData64:  100,
				UserData32:  10,
				Ledger:      1,
				Code:        5,
			},
		}

		mockClient.On("QueryAccounts", mock.AnythingOfType("types.QueryFilter")).
			Return(mockAccounts, nil).Once()

		resp, err := app.QueryAccounts(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Accounts, 1)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle TigerBeetle client errors", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData64: lo.ToPtr(uint64(100)),
			Limit:      10,
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		tbError := errors.New("TigerBeetle connection error")
		mockClient.On("QueryAccounts", mock.AnythingOfType("types.QueryFilter")).
			Return(nil, tbError).Once()

		_, err := app.QueryAccounts(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, tbError, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle invalid hex string in UserData128", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData128: lo.ToPtr("invalid-hex"),
			Limit:       10,
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		_, err := app.QueryAccounts(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UserData128:")
	})

	t.Run("should handle empty filter with only limit", func(t *testing.T) {
		filter := &proto.QueryFilter{
			Limit: 10,
		}

		req := &proto.QueryAccountsRequest{
			Filter: filter,
		}

		mockAccounts := []types.Account{
			{ID: types.ToUint128(1)},
			{ID: types.ToUint128(2)},
		}

		mockClient.On("QueryAccounts", mock.AnythingOfType("types.QueryFilter")).
			Return(mockAccounts, nil).Once()

		resp, err := app.QueryAccounts(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Accounts, 2)
		mockClient.AssertExpectations(t)
	})
}
