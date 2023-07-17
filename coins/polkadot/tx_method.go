package polkadot

import (
	"encoding/hex"
)

const (
	Transfer    = 1 // balances-transfer  balances-transfer_keep_alive
	TransferAll = 2 //  balances-transfer_all
)

// balances-transfer  balances-transfer_keep_alive
func BalanceTransfer(method, to string, amount uint64) ([]byte, error) {
	pubBytes, err := hex.DecodeString(AddressToPublicKey(to))
	var amountBytes []byte
	if amount == 0 {
		amountBytes = []byte{0}
	} else {
		amountStr := Encode(uint64(amount))
		amountBytes, err = hex.DecodeString(amountStr)
	}
	ret, err := hex.DecodeString(method)
	ret = append(ret, 0x00)
	ret = append(ret, pubBytes...)
	ret = append(ret, amountBytes...)
	return ret, err
}

// balances-transfer_all
// false: keepAlive 00
func BalanceTransferAll(method, to, keepAlive string) ([]byte, error) {
	pubBytes, err := hex.DecodeString(AddressToPublicKey(to))
	ret, err := hex.DecodeString(method)
	ret = append(ret, 0x00)
	ret = append(ret, pubBytes...)
	alive, err := hex.DecodeString(keepAlive)
	ret = append(ret, alive...)
	return ret, err
}
