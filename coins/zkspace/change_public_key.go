package zkspace

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/emresenyuva/go-wallet-sdk/coins/zksync/core"
)

type TransactionTypeChangePubKey struct {
	ChangePubKey string `json:"ChangePubKey"`
}

const (
	TransactionTypeChangePubKeyOnchain core.TransactionType = "Onchain"
	TransactionTypeChangePubKeyECDSA   core.TransactionType = "ECDSA"
	TransactionTypeChangePubKeyCREATE2 core.TransactionType = "CREATE2"
)

type ChangePubKey struct {
	Type         string `json:"type"`
	AccountId    uint32 `json:"accountId"`
	Account      string `json:"account"`
	NewPkHash    string `json:"newPkHash"`
	Nonce        uint32 `json:"nonce"`
	EthSignature string `json:"ethSignature"`
}

func (t *ChangePubKey) getType() string {
	return "ChangePubKey"
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
	buf.WriteByte(0x07)
	buf.Write(core.Uint32ToBytes(t.AccountId))
	accountBytes, err := hex.DecodeString(t.Account[2:])
	if err != nil {
		return "", err
	}
	buf.Write(accountBytes)
	newPubKeyHashBytes, err := hex.DecodeString(t.NewPkHash[5:])
	if err != nil {
		return "", err
	}
	buf.Write(newPubKeyHashBytes)
	buf.Write(core.Uint32ToBytes(t.Nonce))
	signatureBytes, err := hex.DecodeString(t.EthSignature[2:])
	if err != nil {
		return "", err
	}
	buf.Write(signatureBytes)
	hash := sha256.New()
	hash.Write(buf.Bytes())
	txHash := "0x" + hex.EncodeToString(hash.Sum(nil))
	return txHash, nil
}
