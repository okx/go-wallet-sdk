package zcash

import (
	"encoding/hex"
	"fmt"
	"github.com/okx/go-wallet-sdk/crypto/zec"
)

// This only works for transparent transactions
func CalTxHash(rawTx string) (string, error) {
	if len(rawTx)%2 != 0 {
		rawTx = "0" + rawTx
	}
	serializedTx, err := hex.DecodeString(rawTx)
	if err != nil {
		return "", err
	}

	tx, err := zec.DeserializeTx(serializedTx)

	if err != nil {
		return "", err
	}

	if tx.NSpendsSapling != 0 ||
		tx.NOutputsSapling != 0 ||
		tx.ValueBalanceSapling != 0 ||
		tx.NActionsOrchard != 0 ||
		tx.SizeProofsOrchard != 0 ||
		tx.NJoinSplit != 0 ||
		len(tx.VJoinSplit) != 0 ||
		tx.ValueBalanceOrchard != 0 {
		return "", fmt.Errorf("not a transparent transaction")
	}

	hash := tx.TxHash().String()
	return hash, nil
}
