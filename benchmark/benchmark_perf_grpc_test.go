package benchmark

import (
	"benchmark/proto"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func BenchmarkGrpcTest(b *testing.B) {
	close := tbApiStart([]string{
		"PORT=50054",
		"TB_ADDRESSES=3033",
		"TB_CLUSTER_ID=0",
		"USE_GRPC=true",
		"GRPC_REFLECTION=true",
		"CLIENT_COUNT=1",
		"MODE=production",
	})
	defer close()

	conn, err := grpc.NewClient("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Error(err)
		return
	}
	defer conn.Close()
	c := proto.NewTigerBeetleClient(conn)

	b.ResetTimer()
	for range b.N {
		res, err := c.CreateTransfers(context.Background(), &proto.CreateTransfersRequest{
			Transfers: []*proto.Transfer{
				{
					Id:              types.ID().String(),
					DebitAccountId:  types.Uint128{2}.String(),
					CreditAccountId: types.Uint128{3}.String(),
					Amount:          1,
					Ledger:          999,
					Code:            1,
					TransferFlags:   &proto.TransferFlags{},
					UserData128:     types.Uint128{0}.String(),
					UserData64:      0,
					UserData32:      0,
				},
				{
					Id:              types.ID().String(),
					DebitAccountId:  types.Uint128{2}.String(),
					CreditAccountId: types.Uint128{3}.String(),
					Amount:          1,
					Ledger:          999,
					Code:            1,
					TransferFlags:   &proto.TransferFlags{},
					UserData128:     types.Uint128{0}.String(),
					UserData64:      0,
					UserData32:      0,
				},
			},
		})
		if err == nil {
			if res == nil {
				err = errors.New("Grpc result is nil")
			} else if len(res.Results) != 0 {
				err = fmt.Errorf("Results gt 0 %v", res)
			}
		}
		if err != nil {
			b.Error(err)
			return
		}
	}
}
