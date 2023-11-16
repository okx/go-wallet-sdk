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

type TransactionTypeChangePubKey struct {
	ChangePubKey string `json:"ChangePubKey"`
}

const (
	TransactionTypeChangePubKey_                       = "ChangePubKey"
	TransactionTypeChangePubKeyOnchain TransactionType = "Onchain"
	TransactionTypeChangePubKeyECDSA   TransactionType = "ECDSA"
	TransactionTypeChangePubKeyCREATE2 TransactionType = "CREATE2"
)

type ChangePubKey struct {
	Type        string              `json:"type"`
	AccountId   uint32              `json:"accountId"`
	Account     string              `json:"account"`
	NewPkHash   string              `json:"newPkHash"`
	FeeToken    uint32              `json:"feeToken"`
	Fee         string              `json:"fee"`
	Nonce       uint32              `json:"nonce"`
	Signature   *Signature          `json:"signature"`
	EthAuthData ChangePubKeyVariant `json:"ethAuthData"`
	*TimeRange
}

func (t *ChangePubKey) getType() string {
	return TransactionTypeChangePubKey_
}

type ChangePubKeyAuthType string

const (
	ChangePubKeyAuthTypeOnchain ChangePubKeyAuthType = `Onchain`
	ChangePubKeyAuthTypeECDSA   ChangePubKeyAuthType = `ECDSA`
	ChangePubKeyAuthTypeCREATE2 ChangePubKeyAuthType = `CREATE2`
)

type ChangePubKeyVariant interface {
	getType() ChangePubKeyAuthType
	getBytes() []byte
}

type ChangePubKeyOnchain struct {
	Type ChangePubKeyAuthType `json:"type"`
}

func (t *ChangePubKeyOnchain) getType() ChangePubKeyAuthType {
	return ChangePubKeyAuthTypeOnchain
}

func (t *ChangePubKeyOnchain) getBytes() []byte {
	return make([]byte, 32)
}

type ChangePubKeyECDSA struct {
	Type         ChangePubKeyAuthType `json:"type"`
	EthSignature string               `json:"ethSignature"`
	BatchHash    string               `json:"batchHash"`
}

func (t *ChangePubKeyECDSA) getType() ChangePubKeyAuthType {
	return ChangePubKeyAuthTypeECDSA
}

func (t *ChangePubKeyECDSA) getBytes() []byte {
	res, _ := hex.DecodeString(t.BatchHash[2:])
	return res
}

type ChangePubKeyCREATE2 struct {
	Type           ChangePubKeyAuthType `json:"type"`
	CreatorAddress string               `json:"creatorAddress"`
	SaltArg        string               `json:"saltArg"`
	CodeHash       string               `json:"codeHash"`
}

func (t *ChangePubKeyCREATE2) getType() ChangePubKeyAuthType {
	return ChangePubKeyAuthTypeCREATE2
}

func (t *ChangePubKeyCREATE2) getBytes() []byte {
	return make([]byte, 32)
}

type Signature struct {
	PubKey    string `json:"pubKey"`
	Signature string `json:"signature"`
}

func (t *ChangePubKey) GetTxHash() (string, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xff - 0x07)
	buf.WriteByte(TransactionVersion)
	buf.Write(Uint32ToBytes(t.AccountId))
	buf.Write(ParseAddress(t.Account))
	pkhBytes, err := pkhToBytes(t.NewPkHash)
	if err != nil {
		return "", err
	}
	buf.Write(pkhBytes)
	buf.Write(Uint32ToBytes(t.FeeToken))
	fee, ok := big.NewInt(0).SetString(t.Fee, 10)
	if !ok {
		return "", errors.New("failed to convert string fee to big.Int")
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
