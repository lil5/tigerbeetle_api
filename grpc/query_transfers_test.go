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

// MockTigerBeetleClient is a mock client for testing
type MockTigerBeetleClient struct {
	mock.Mock
}

func (m *MockTigerBeetleClient) QueryTransfers(filter types.QueryFilter) ([]types.Transfer, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.Transfer), args.Error(1)
}

func (m *MockTigerBeetleClient) QueryAccounts(filter types.QueryFilter) ([]types.Account, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.Account), args.Error(1)
}

func (m *MockTigerBeetleClient) Nop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTigerBeetleClient) Close() {
	m.Called()
}

func (m *MockTigerBeetleClient) CreateAccounts(accounts []types.Account) ([]types.AccountEventResult, error) {
	args := m.Called(accounts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.AccountEventResult), args.Error(1)
}

func (m *MockTigerBeetleClient) CreateTransfers(transfers []types.Transfer) ([]types.TransferEventResult, error) {
	args := m.Called(transfers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.TransferEventResult), args.Error(1)
}

func (m *MockTigerBeetleClient) LookupAccounts(ids []types.Uint128) ([]types.Account, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.Account), args.Error(1)
}

func (m *MockTigerBeetleClient) LookupTransfers(ids []types.Uint128) ([]types.Transfer, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.Transfer), args.Error(1)
}

func (m *MockTigerBeetleClient) GetAccountTransfers(filter types.AccountFilter) ([]types.Transfer, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.Transfer), args.Error(1)
}

func (m *MockTigerBeetleClient) GetAccountBalances(filter types.AccountFilter) ([]types.AccountBalance, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.AccountBalance), args.Error(1)
}

func TestQueryTransfers(t *testing.T) {
	mockClient := new(MockTigerBeetleClient)
	app := &App{TB: mockClient}

	t.Run("should return error when filter is nil", func(t *testing.T) {
		req := &proto.QueryTransfersRequest{
			Filter: nil,
		}
		
		_, err := app.QueryTransfers(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "filter is required", err.Error())
	})

	t.Run("should return empty array when no transfers match", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData128: lo.ToPtr("0"),
			UserData64:  lo.ToPtr(uint64(0)),
			UserData32:  lo.ToPtr(uint32(0)),
			Code:        lo.ToPtr(uint32(1)),
			Limit:       10,
		}

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		mockClient.On("QueryTransfers", mock.AnythingOfType("types.QueryFilter")).
			Return([]types.Transfer{}, nil).Once()

		resp, err := app.QueryTransfers(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.Transfers)
		mockClient.AssertExpectations(t)
	})

	t.Run("should return transfers with valid filter", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData128: lo.ToPtr("1"),
			UserData64:  lo.ToPtr(uint64(100)),
			UserData32:  lo.ToPtr(uint32(10)),
			Code:        lo.ToPtr(uint32(1)),
			Ledger:      lo.ToPtr(uint32(0)),
			Limit:       10,
		}

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		// Mock transfer data
		mockTransfers := []types.Transfer{
			{
				ID:              types.ToUint128(1),
				DebitAccountID:  types.ToUint128(100),
				CreditAccountID: types.ToUint128(200),
				Amount:          types.ToUint128(1000),
				UserData128:     types.ToUint128(1),
				UserData64:      100,
				UserData32:      10,
				Ledger:          0,
				Code:            1,
				Timestamp:       1000000,
			},
		}

		mockClient.On("QueryTransfers", mock.AnythingOfType("types.QueryFilter")).
			Return(mockTransfers, nil).Once()

		resp, err := app.QueryTransfers(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Transfers, 1)
		assert.Equal(t, "1", resp.Transfers[0].Id)
		assert.Equal(t, uint64(100), resp.Transfers[0].UserData64)
		assert.Equal(t, uint32(10), resp.Transfers[0].UserData32)
		mockClient.AssertExpectations(t)
	})

	t.Run("should respect limit parameter", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData64: lo.ToPtr(uint64(100)),
			Limit:      2,
		}

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		// Mock 5 transfers but expect only 2 due to limit
		mockTransfers := []types.Transfer{
			{ID: types.ToUint128(1), UserData64: 100},
			{ID: types.ToUint128(2), UserData64: 100},
		}

		mockClient.On("QueryTransfers", mock.AnythingOfType("types.QueryFilter")).
			Return(mockTransfers, nil).Once()

		resp, err := app.QueryTransfers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Transfers, 2)
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

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		// Mock transfers in reverse order
		mockTransfers := []types.Transfer{
			{ID: types.ToUint128(3), Timestamp: 3000},
			{ID: types.ToUint128(2), Timestamp: 2000},
			{ID: types.ToUint128(1), Timestamp: 1000},
		}

		mockClient.On("QueryTransfers", mock.AnythingOfType("types.QueryFilter")).
			Return(mockTransfers, nil).Once()

		resp, err := app.QueryTransfers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Transfers, 3)
		// Verify order is reversed (newest first)
		assert.Equal(t, "3", resp.Transfers[0].Id)
		assert.Equal(t, "2", resp.Transfers[1].Id)
		assert.Equal(t, "1", resp.Transfers[2].Id)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle timestamp range filters", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData64:   lo.ToPtr(uint64(100)),
			TimestampMin: lo.ToPtr(uint64(1000)),
			TimestampMax: lo.ToPtr(uint64(5000)),
			Limit:        10,
		}

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		// Mock transfers within timestamp range
		mockTransfers := []types.Transfer{
			{ID: types.ToUint128(2), Timestamp: 2000},
			{ID: types.ToUint128(3), Timestamp: 3000},
			{ID: types.ToUint128(4), Timestamp: 4000},
		}

		mockClient.On("QueryTransfers", mock.AnythingOfType("types.QueryFilter")).
			Return(mockTransfers, nil).Once()

		resp, err := app.QueryTransfers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Transfers, 3)
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

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		mockTransfers := []types.Transfer{
			{
				ID:              types.ToUint128(1),
				DebitAccountID:  types.ToUint128(100),
				CreditAccountID: types.ToUint128(200),
				Amount:          types.ToUint128(1000),
				UserData128:     types.ToUint128(1000),
				UserData64:      100,
				UserData32:      10,
				Ledger:          1,
				Code:            5,
			},
		}

		mockClient.On("QueryTransfers", mock.AnythingOfType("types.QueryFilter")).
			Return(mockTransfers, nil).Once()

		resp, err := app.QueryTransfers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Transfers, 1)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle TigerBeetle client errors", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData64: lo.ToPtr(uint64(100)),
			Limit:      10,
		}

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		tbError := errors.New("TigerBeetle connection error")
		mockClient.On("QueryTransfers", mock.AnythingOfType("types.QueryFilter")).
			Return(nil, tbError).Once()

		_, err := app.QueryTransfers(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, tbError, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle invalid hex string in UserData128", func(t *testing.T) {
		filter := &proto.QueryFilter{
			UserData128: lo.ToPtr("invalid-hex"),
			Limit:       10,
		}

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		_, err := app.QueryTransfers(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UserData128:")
	})

	t.Run("should handle empty filter with only limit", func(t *testing.T) {
		filter := &proto.QueryFilter{
			Limit: 10,
		}

		req := &proto.QueryTransfersRequest{
			Filter: filter,
		}

		mockTransfers := []types.Transfer{
			{ID: types.ToUint128(1)},
			{ID: types.ToUint128(2)},
		}

		mockClient.On("QueryTransfers", mock.AnythingOfType("types.QueryFilter")).
			Return(mockTransfers, nil).Once()

		resp, err := app.QueryTransfers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, resp.Transfers, 2)
		mockClient.AssertExpectations(t)
	})
}