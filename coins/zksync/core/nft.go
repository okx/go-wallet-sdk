/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
package core

const (
	TransactionTypeMintNFT     = "MintNFT"
	TransactionTypeWithdrawNFT = "WithdrawNFT"
)

type NFT struct {
	Id             uint32 `json:"id"`
	Symbol         string `json:"symbol"`
	CreatorId      uint32 `json:"creatorId"`
	ContentHash    Hash   `json:"contentHash"`
	CreatorAddress string `json:"creatorAddress"`
	SerialId       uint32 `json:"serialId"`
	Address        string `json:"address"`
}

func (t *NFT) ToToken() *Token {
	return &Token{
		Id:       t.Id,
		Address:  t.Address,
		Symbol:   t.Symbol,
		Decimals: 0,
		IsNft:    true,
	}
}

type MintNFT struct {
	Type           string     `json:"type"`
	CreatorId      uint32     `json:"creatorId"`
	CreatorAddress string     `json:"creatorAddress"`
	ContentHash    Hash       `json:"contentHash"`
	Recipient      string     `json:"recipient"`
	Fee            string     `json:"fee"`
	FeeToken       uint32     `json:"feeToken"`
	Nonce          uint32     `json:"nonce"`
	Signature      *Signature `json:"signature"`
}

func (t *MintNFT) getType() string {
	return TransactionTypeMintNFT
}

type WithdrawNFT struct {
	Type      string     `json:"type"`
	AccountId uint32     `json:"accountId"`
	From      string     `json:"from"`
	To        string     `json:"to"`
	Token     uint32     `json:"token"`
	FeeToken  uint32     `json:"feeToken"`
	Fee       string     `json:"fee"`
	Nonce     uint32     `json:"nonce"`
	Signature *Signature `json:"signature"`
	*TimeRange
}

func (t *WithdrawNFT) getType() string {
	return TransactionTypeWithdrawNFT
}
