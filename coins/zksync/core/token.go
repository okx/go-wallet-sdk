/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
package core

import (
	"errors"
	"math/big"
)

type Token struct {
	Id       uint32 `json:"id"`
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals uint   `json:"decimals"`
	IsNft    bool   `json:"is_nft"`
}

func CreateETH() *Token {
	return &Token{
		Id:       0,
		Address:  "0x0000000000000000000000000000000000000000",
		Symbol:   `ETH`,
		Decimals: 18,
	}
}

func (t Token) IsETH() bool {
	return t.Address == "0x0000000000000000000000000000000000000000" && t.Symbol == `ETH`
}

func (t Token) GetAddress() string {
	return t.Address
}

func (t Token) ToDecimalString(amount *big.Int) string {
	amountFloat := big.NewFloat(0).SetInt(amount)
	if t.IsNft {
		// return origin int value in "XXX.0" format
		return amountFloat.Text('f', 1)
	}
	// convert to pointed value considering decimals scale, like wei => ETH (10^18 wei == 1 ETH)
	divider := big.NewFloat(0).SetInt(big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(t.Decimals)), nil)) // = 10^decimals
	res := big.NewFloat(0).SetPrec(t.Decimals*8).Quo(amountFloat, divider)                                   // = amount / 10^decimals
	if res.IsInt() {
		return res.Text('f', 1) // return int value in "XXX.0" format
	}
	return res.Text('f', -int(t.Decimals)) // format as numeric with specified decimals limit
}

type Tokens struct {
	Tokens map[string]*Token
}

func (ts *Tokens) GetToken(id string) (*Token, error) {
	if t, ok := ts.Tokens[id]; ok {
		return t, nil
	}
	// suppose id is address
	for _, t := range ts.Tokens {
		if t.Address == id {
			return t, nil
		}
	}
	return nil, errors.New("token not found")
}
