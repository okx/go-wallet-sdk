/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
package core

import "math/big"

type AccountState struct {
	Address    string           `json:"address"`
	Id         uint32           `json:"id"`
	Depositing *DepositingState `json:"depositing"`
	Committed  *State           `json:"committed"`
	Verified   *State           `json:"verified"`
}

type DepositingState struct {
	Balances map[string]*DepositingBalance `json:"balances"`
}

type DepositingBalance struct {
	Amount              string `json:"amount"`
	ExpectedBlockNumber string `json:"expectedBlockNumber"` // to *big.Int
}

type State struct {
	Balances   map[string]string `json:"balances"`
	Nonce      uint32            `json:"nonce"`
	PubKeyHash string            `json:"pubKeyHash"`
	Nfts       map[string]*NFT   `json:"nfts"`
	MintedNfts map[string]*NFT   `json:"mintedNfts"`
}

func (s *State) GetBalanceOf(token string) (*big.Int, bool) {
	n := new(big.Int)
	if v, ok := s.Balances[token]; ok {
		if n, ok := n.SetString(v, 10); ok {
			return n, true
		}
		return new(big.Int), false
	}
	return n, false
}
