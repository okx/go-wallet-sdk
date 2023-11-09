package zksync

import "github.com/okx/go-wallet-sdk/coins/zksync/core"

var RinkebyUSDT = &core.Token{
	Id:       1,
	Address:  "0x3b00ef435fa4fcff5c209a37d1f3dcff37c705ad",
	Symbol:   "USDT",
	Decimals: 6,
	IsNft:    false,
}

var RinkebyUSDC = &core.Token{
	Id:       2,
	Address:  "0xeb8f08a975ab53e34d8a0330e0d34de942c95926",
	Symbol:   "USDC",
	Decimals: 6,
	IsNft:    false,
}
