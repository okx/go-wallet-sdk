package tezos

import (
	"bytes"
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/coin/tezos/types"
)

// TODO: fetch dynamic from /chains/main/mempool/filter
const (
	minFeeFixedNanoTez int64 = 100_000
	minFeeByteNanoTez  int64 = 1_000
	minFeeGasNanoTez   int64 = 100
)

// CalculateMinFee returns the minimum fee at/above which bakers will accept this operation under default config settings.
func CalculateMinFee(o types.Operation, gas int64, withHeader bool, p *types.Params) int64 {
	buf := bytes.NewBuffer(nil)
	_ = o.EncodeBuffer(buf, p)
	sz := int64(buf.Len())
	if withHeader {
		sz += 32 + 64 // branch + signature
	}
	fee := minFeeFixedNanoTez + sz*minFeeByteNanoTez + gas*minFeeGasNanoTez
	return fee / 1000 // nano -> micro
}
