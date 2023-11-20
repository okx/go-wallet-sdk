package polkadot

import (
	"encoding/hex"
)

const (
	Transfer    = 1 // balances-transfer  balances-transfer_keep_alive
	TransferAll = 2 //  balances-transfer_all
)

// BalanceTransfer balances-transfer_keep_alive
func BalanceTransfer(method, to string, amount uint64) ([]byte, error) {
	pubKey, err := AddressToPublicKey(to)
	if err != nil {
		return nil, err
	}
	pubBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	var amountBytes []byte
	if amount == 0 {
		amountBytes = []byte{0}
	} else {
		amountStr := Encode(uint64(amount))
		amountBytes, err = hex.DecodeString(amountStr)
		if err != nil {
			return nil, err
		}
	}
	ret, err := hex.DecodeString(method)
	if err != nil {
		return nil, err
	}
	ret = append(ret, 0x00)
	ret = append(ret, pubBytes...)
	ret = append(ret, amountBytes...)
	return ret, nil
}

// BalanceTransferAll
// false: keepAlive 00
func BalanceTransferAll(method, to, keepAlive string) ([]byte, error) {
	pubKey, err := AddressToPublicKey(to)
	if err != nil {
		return nil, err
	}
	pubBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	ret, err := hex.DecodeString(method)
	if err != nil {
		return nil, err
	}
	ret = append(ret, 0x00)
	ret = append(ret, pubBytes...)
	alive, err := hex.DecodeString(keepAlive)
	if err != nil {
		return nil, err
	}
	ret = append(ret, alive...)
	return ret, nil
}
