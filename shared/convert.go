package shared

import (
	"log/slog"
	"time"

	"github.com/samber/lo"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func HexStringToUint128(hex string) (*types.Uint128, error) {
	if hex == "" {
		return &types.Uint128{}, nil
	}

	res, err := types.HexStringToUint128(hex)
	if err != nil {
		slog.Error("hex string to Uint128 failed", "hex", hex, "error", err)
		return nil, err
	}
	return &res, nil

}

func GetOrCreateID(id string) (idStr string, idUint128 types.Uint128, err error) {
	if id == "" {
		idUint128 = types.ID()
		idStr = idUint128.String()
	} else {
		idStr = id
		idUint128, err = types.HexStringToUint128(idStr)
	}
	return
}

// set to zero if timestamp is nil
func TimestampFromPstringToUint(timestamp *string) (*uint64, error) {
	if timestamp == nil {
		return lo.ToPtr[uint64](0), nil
	}
	if *timestamp == "" {
		return lo.ToPtr[uint64](0), nil
	}

	return TimestampFromStringToUint(*timestamp)
}

func TimestampFromUintToString(timestamp uint64) string {
	return time.Unix(0, int64(timestamp)).Format(time.RFC3339Nano)
}

func TimestampFromStringToUint(timestamp string) (*uint64, error) {
	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return nil, err
	}

	nano := t.UnixNano()

	return lo.ToPtr(uint64(nano)), nil
}
