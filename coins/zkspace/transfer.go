package zkspace

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/zksync/core"
	"math/big"
)

const (
	TransactionTypeTransfer core.TransactionType = "Transfer"
)

type Transfer struct {
	Type       string     `json:"type"`
	AccountId  uint32     `json:"accountId"`
	From       string     `json:"from"`
	To         string     `json:"to"`
	TokenId    uint16     `json:"token"`
	Amount     *big.Int   `json:"amount"`
	FeeTokenId uint8      `json:"feeToken"`
	Fee        *big.Int   `json:"fee"`
	ChainId    uint8      `json:"chainId"`
	Nonce      uint32     `json:"nonce"`
	Signature  *Signature `json:"signature"`
}

type SignTransfer struct {
	Type       string     `json:"type"`
	AccountId  uint32     `json:"accountId"`
	From       string     `json:"from"`
	To         string     `json:"to"`
	TokenId    uint16     `json:"token"`
	Amount     string     `json:"amount"`
	FeeTokenId uint8      `json:"feeToken"`
	Fee        string     `json:"fee"`
	ChainId    uint8      `json:"chainId"`
	Nonce      uint32     `json:"nonce"`
	Signature  *Signature `json:"signature"`
}

func (t *Transfer) getType() string {
	return "Transfer"
}

func (t *Transfer) GetTxHash() (string, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0x05)
	buf.Write(core.Uint32ToBytes(t.AccountId))
	fromBytes, err := hex.DecodeString(t.From[2:])
	if err != nil {
		return "", err
	}
	buf.Write(fromBytes)
	toBytes, err := hex.DecodeString(t.To[2:])
	if err != nil {
		return "", err
	}
	buf.Write(toBytes)
	buf.Write(core.Uint16ToBytes(t.TokenId))
	packedAmount, err := core.PackAmount(t.Amount)
	if err != nil {
		return "", err
	}
	buf.Write(packedAmount)
	buf.WriteByte(t.FeeTokenId)
	packedFee, err := core.PackFee(t.Fee)
	if err != nil {
		return "", err
	}
	buf.Write(packedFee)
	buf.WriteByte(t.ChainId)
	buf.Write(core.Uint32ToBytes(t.Nonce))
	hash := sha256.New()
	hash.Write(buf.Bytes())
	txHash := "0x" + hex.EncodeToString(hash.Sum(nil))
	return txHash, nil
}
