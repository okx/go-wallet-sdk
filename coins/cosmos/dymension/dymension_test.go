package dymension

import (
	"testing"
	"time"

	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "dym1jl389rqgh59lhf33j20pp082aj8utjtpppjh5j"
	require.Equal(t, expected, address)
	ret := ValidateAddress(address)
	require.True(t, ret)

	param := cosmos.TransferParam{}
	param.FromAddress = "dym1yc4q6svsl9xy9g2gplgnlpxwhnzr3y73wfs0xh"
	param.ToAddress = "dym1yc4q6svsl9xy9g2gplgnlpxwhnzr3y73wfs0xh"
	param.Demon = "adym"
	param.Amount = "10000000000000000" // 18
	param.CommonParam.ChainId = "dymension_1100-1"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 2091572
	param.CommonParam.FeeDemon = "adym"
	param.CommonParam.FeeAmount = "3500000000000000"
	param.CommonParam.GasLimit = 140000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	signedTx, err := cosmos.TransferAction(param, privateKeyHex, true)
	require.Nil(t, err)
	t.Log("signedTx : ", signedTx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "dymension_1100-1"
	p.CommonParam.Sequence = 1
	p.CommonParam.AccountNumber = 2091572
	p.CommonParam.FeeDemon = "adym"
	p.CommonParam.FeeAmount = "4000000000000000"
	p.CommonParam.GasLimit = 200000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "dym1yc4q6svsl9xy9g2gplgnlpxwhnzr3y73wfs0xh"
	p.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	p.Demon = "adym"
	p.Amount = "10000000000000000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-3"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	signedIBCTx, err := cosmos.IbcTransferAction(p, privateKeyHex, true)
	require.Nil(t, err)
	t.Log("signedIBCTx : ", signedIBCTx)
}
