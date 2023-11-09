package tezos

import (
	"encoding/json"
	"github.com/okx/go-wallet-sdk/coins/tezos/types"
	"strconv"
)

var DefaultOptions = CallOptions{
	MaxFee: 1_000_000,
}

type CallOptions struct {
	MaxFee       int64 // max acceptable fee, optional (default = 0)
	IgnoreLimits bool  // ignore simulated limits and use user-defined limits from op
	// custom options
	BlockHash  types.BlockHash
	Counter    int64 // number of times counters
	NeedReveal bool
}

func NewCallOptions(blockHash string, counter int64, needReveal bool) *CallOptions {
	var hash types.BlockHash
	if err := json.Unmarshal([]byte(strconv.Quote(blockHash)), &hash); err != nil {
		hash = types.BlockHash{}
	}
	return &CallOptions{BlockHash: hash, Counter: counter, NeedReveal: needReveal}
}
