package osmo

import (
	"errors"
	"github.com/emresenyuva/go-wallet-sdk/coins/cosmos"
	"github.com/emresenyuva/go-wallet-sdk/coins/cosmos/osmo/tx/gamm"
	"github.com/emresenyuva/go-wallet-sdk/coins/cosmos/types"
)

const (
	HRP = "osmo"
)

type SwapExactAmountInParam struct {
	cosmos.CommonParam
	Sender      string
	PoolId      uint64
	ToDemon     string
	FromDemon   string
	FromAmount  string
	MinToAmount string
}

func NewAddress(privateKeyHex string) (string, error) {
	return cosmos.NewAddress(privateKeyHex, HRP, false)
}

func ValidateAddress(address string) bool {
	return cosmos.ValidateAddress(address, HRP)
}

func SwapExactAmountIn(param SwapExactAmountInParam, privateKeyHex string) (string, error) {
	fa, ok := types.NewIntFromString(param.FromAmount)
	if !ok {
		return "", errors.New("invalid from amount")
	}
	inCoin := types.NewCoin(param.FromDemon, fa)

	tmo, ok := types.NewIntFromString(param.MinToAmount)
	if !ok {
		return "", errors.New("invalid min to  amount")
	}

	routes := make([]gamm.SwapAmountInRoute, 0)
	routes = append(routes, gamm.SwapAmountInRoute{PoolId: param.PoolId, TokenOutDenom: param.ToDemon})
	sendMsg := gamm.MsgSwapExactAmountIn{Sender: param.Sender, Routes: routes, TokenIn: inCoin, TokenOutMinAmount: tmo}

	messages := make([]*types.Any, 0)
	anySend, err := types.NewAnyWithValue(&sendMsg)
	if err != nil {
		return "", err
	}
	messages = append(messages, anySend)
	return cosmos.BuildTx(param.CommonParam, messages, privateKeyHex)
}
