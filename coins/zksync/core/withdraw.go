/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
package core

import (
	"math/big"
)

const (
	TransactionTypeWithdraw   = "Withdraw"
	TransactionTypeForcedExit = "ForcedExit"
)

type Withdraw struct {
	Type      string     `json:"type"`
	AccountId uint32     `json:"accountId"`
	From      string     `json:"from"`
	To        string     `json:"to"`
	TokenId   uint32     `json:"token"`
	Amount    *big.Int   `json:"amount"`
	Fee       string     `json:"fee"`
	Nonce     uint32     `json:"nonce"`
	Signature *Signature `json:"signature"`
	*TimeRange
}

func (t *Withdraw) getType() string {
	return TransactionTypeWithdraw
}

type ForcedExit struct {
	Type      string     `json:"type"`
	AccountId uint32     `json:"initiatorAccountId"`
	Target    string     `json:"target"`
	TokenId   uint32     `json:"token"`
	Amount    *big.Int   `json:"amount"`
	Fee       string     `json:"fee"`
	Nonce     uint32     `json:"nonce"`
	Signature *Signature `json:"signature"`
	*TimeRange
}

func (t *ForcedExit) getType() string {
	return TransactionTypeForcedExit
}
