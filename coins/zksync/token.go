package zksync

import "github.com/okx/go-wallet-sdk/coins/zksync/core"

var ETH = core.CreateETH()

var USDT = &core.Token{
	Id:       4,
	Address:  "0xdac17f958d2ee523a2206206994597c13d831ec7",
	Symbol:   "USDT",
	Decimals: 6,
	IsNft:    false,
}
