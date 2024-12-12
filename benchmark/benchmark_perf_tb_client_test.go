package benchmark

import (
	"testing"

	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func BenchmarkTbClient(b *testing.B) {
	tbAddresses := []string{"3033"}
	clusterUint128 := types.Uint128{0}

	// Start the tb client
	tb, err := tigerbeetle_go.NewClient(clusterUint128, tbAddresses)
	if err != nil {
		b.Error(err)
	}
	defer tb.Close()

	b.ResetTimer()
	for range b.N {
		f := types.TransferFlags{}
		res, err := tb.CreateTransfers([]types.Transfer{

			{
				ID:              types.ID(),
				DebitAccountID:  types.Uint128{2},
				CreditAccountID: types.Uint128{3},
				Amount:          types.Uint128{1},
				Ledger:          999,
				Code:            1,
				Flags:           f.ToUint16(),
				UserData128:     types.Uint128{0},
				UserData64:      0,
				UserData32:      0,
			},
			{
				ID:              types.ID(),
				DebitAccountID:  types.Uint128{2},
				CreditAccountID: types.Uint128{3},
				Amount:          types.Uint128{1},
				Ledger:          999,
				Code:            1,
				Flags:           f.ToUint16(),
				UserData128:     types.Uint128{0},
				UserData64:      0,
				UserData32:      0,
			},
		})

		if err != nil {
			b.Error(err)
		}
		if len(res) != 0 {
			b.Error("Results gt 0", res)
		}
	}
}
