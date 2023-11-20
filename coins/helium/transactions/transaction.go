package transactions

import "math"

type ITransaction interface {
	Serialize() ([]byte, error)
	BuildTransaction() ([]byte, error)
	SetSignature(sig []byte)
}

var (
	DC_Payload_Size    = 24
	Txn_Fee_Multiplier = 0
)

func CalculateFee(length, dc_payload_size, txn_fee_multiplier int64) uint64 {
	return uint64(math.Ceil(float64(length)/float64(dc_payload_size))) * uint64(txn_fee_multiplier)
}
