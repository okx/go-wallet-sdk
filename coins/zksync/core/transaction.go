/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
package core

type TransactionType string

type TransactionStatus int

const (
	TransactionStatusSent TransactionStatus = iota
	TransactionStatusCommitted
	TransactionStatusVerified
)

func (t TransactionType) getType() interface{} {
	switch t {
	case TransactionTypeChangePubKeyOnchain, TransactionTypeChangePubKeyECDSA, TransactionTypeChangePubKeyCREATE2:
		// custom object instead of string
		return TransactionTypeChangePubKey{ChangePubKey: string(t)}
	default:
		return string(t)
	}
}

type ZksTransaction interface {
	getType() string
}

type SignedTransaction struct {
	Transaction       ZksTransaction `json:"tx"`
	EthereumSignature *EthSignature  `json:"signature"`
}

type TransactionDetails struct {
	Executed   bool       `json:"executed"`
	Success    bool       `json:"success"`
	FailReason string     `json:"failReason"`
	Block      *BlockInfo `json:"block"`
}

type BlockInfo struct {
	BlockNumber uint64 `json:"blockNumber"`
	Committed   bool   `json:"committed"`
	Verified    bool   `json:"verified"`
}

type EthOpInfo struct {
	Executed bool       `json:"executed"`
	Block    *BlockInfo `json:"block"`
}
