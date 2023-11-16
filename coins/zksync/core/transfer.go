/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
)

const (
	TransactionTypeTransfer = "Transfer"
)

type Transfer struct {
	Type      string     `json:"type"`
	AccountId uint32     `json:"accountId"`
	From      string     `json:"from"`
	To        string     `json:"to"`
	Token     *Token     `json:"-"`
	TokenId   uint32     `json:"token"`
	Amount    *big.Int   `json:"amount"`
	Fee       string     `json:"fee"`
	Nonce     uint32     `json:"nonce"`
	Signature *Signature `json:"signature"`
	*TimeRange
}

func (t *Transfer) getType() string {
	return TransactionTypeTransfer
}

func (t *Transfer) GetTxHash() (string, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xff - 0x05)
	buf.WriteByte(TransactionVersion)
	buf.Write(Uint32ToBytes(t.AccountId))
	buf.Write(ParseAddress(t.From))
	buf.Write(ParseAddress(t.To))
	buf.Write(Uint32ToBytes(t.Token.Id))
	packedAmount, err := packAmount(t.Amount)
	if err != nil {
		return "", errors.New("failed to pack amount")
	}
	buf.Write(packedAmount)
	fee, ok := big.NewInt(0).SetString(t.Fee, 10)
	if !ok {
		return "", ErrConvertBigInt
	}
	packedFee, err := packFee(fee)
	if err != nil {
		return "", err
	}
	buf.Write(packedFee)
	buf.Write(Uint32ToBytes(t.Nonce))
	buf.Write(Uint64ToBytes(t.TimeRange.ValidFrom))
	buf.Write(Uint64ToBytes(t.TimeRange.ValidUntil))
	hash := sha256.New()
	hash.Write(buf.Bytes())
	txHash := HEX_PREFIX + hex.EncodeToString(hash.Sum(nil))
	return txHash, nil
}
