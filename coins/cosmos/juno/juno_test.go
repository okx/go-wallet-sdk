package juno

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// https://lcd-juno.itastakers.com/cosmos/auth/v1beta1/accounts/juno1rvs5xph4l3px2efynqsthus8p6r4exyrg24lps
// curl -X POST -d '{"tx_bytes":"Cr4BCrsBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKNAQoIdHJhbnNmZXISCWNoYW5uZWwtMRoOCgV1anVubxIFMTAwMDAiK2p1bm8xcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXJnMjRscHMqLWNvc21vczFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjdja3l4djIAOICgo4iIo9+BFxJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDnCTAPeltdUKTRGnHRm5QS12GCMOZ3du0gUsZYfj3OlMSBAoCCAEYAhISCgwKBXVqdW5vEgMyNTAQkKEPGkBl30ljLCTSBg3+yQF0s/4c8G9/8uXaplKoZsDCQdwh1Uv2NeRUjVqDB1knbP4FYpmMy7epqLlja6dsKYnjvR74","mode":"BROADCAST_MODE_SYNC"}' https://lcd-juno.itastakers.com/cosmos/tx/v1beta1/txs
func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "juno145q0tcdur4tcx2ya5cphqx96e54yflfykqqkvg"
	require.Equal(t, expected, address)

	param := cosmos.TransferParam{}
	param.FromAddress = "juno1rvs5xph4l3px2efynqsthus8p6r4exyrg24lps"
	param.ToAddress = "juno1rvs5xph4l3px2efynqsthus8p6r4exyrg24lps"
	param.Demon = "ujuno"
	param.Amount = "10000"
	param.CommonParam.ChainId = "juno-1"
	param.CommonParam.Sequence = 1
	param.CommonParam.AccountNumber = 313126
	param.CommonParam.FeeDemon = "ujuno"
	param.CommonParam.FeeAmount = "250"
	param.CommonParam.GasLimit = 250000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	signedTx, _ := cosmos.Transfer(param, privateKeyHex)
	require.Nil(t, err)
	t.Log("signedTx : ", signedTx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "juno-1"
	p.CommonParam.Sequence = 2
	p.CommonParam.AccountNumber = 313126
	p.CommonParam.FeeDemon = "ujuno"
	p.CommonParam.FeeAmount = "250"
	p.CommonParam.GasLimit = 250000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "juno1rvs5xph4l3px2efynqsthus8p6r4exyrg24lps"
	p.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	p.Demon = "ujuno"
	p.Amount = "10000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-1"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	signedIBCTx, err := cosmos.IbcTransfer(p, privateKeyHex)
	require.Nil(t, err)
	t.Log("signedIBCTx : ", signedIBCTx)
}
