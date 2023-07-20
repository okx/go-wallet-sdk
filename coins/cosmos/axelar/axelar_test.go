package axelar

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"testing"
	"time"
)

// https://k8s-testnet-axelarco-c0dd71f944-b4c8da2f814e7b8f.elb.us-east-2.amazonaws.com:1317/cosmos/auth/v1beta1/accounts/axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd
// curl -X POST -d '{"tx_bytes":"Cr4BCrsBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKNAQoIdHJhbnNmZXISCWNoYW5uZWwtMRoOCgR1YXhsEgYxMDAwMDAiLWF4ZWxhcjFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjZrcXZkZCorb3NtbzFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cmtyOTVzNzIAOICg2b6uxNuBFxJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDnCTAPeltdUKTRGnHRm5QS12GCMOZ3du0gUsZYfj3OlMSBAoCCAEYARISCgwKBHVheGwSBDEwMDAQoI0GGkBEmWixhNVWDgKgUJ4SVB/vYiWu69sdqmAp52ZDVmsZHlSuXNK/hISNeS7sTA1b5DdCEYHcYl+ieq6Z6ubOl2Z0","mode":"BROADCAST_MODE_SYNC"}' http://k8s-testnet-axelarco-c0dd71f944-b4c8da2f814e7b8f.elb.us-east-2.amazonaws.com:1317/cosmos/tx/v1beta1/txs
func TestTransfer(t *testing.T) {
	privateKeyHex := "//todo please replace your hex cosmos key"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	// axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd
	fmt.Println(address)

	param := cosmos.TransferParam{}
	param.FromAddress = "axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd"
	param.ToAddress = "axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd"
	param.Demon = "uaxl"
	param.Amount = "100000"
	param.CommonParam.ChainId = "axelar-testnet-lisbon-3"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 24569
	param.CommonParam.FeeDemon = "uaxl"
	param.CommonParam.FeeAmount = "1000"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	tt, _ := cosmos.Transfer(param, privateKeyHex)
	t.Log(tt)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "//todo please replace your hex cosmos key"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "axelar-testnet-lisbon-3"
	p.CommonParam.Sequence = 1
	p.CommonParam.AccountNumber = 24569
	p.CommonParam.FeeDemon = "uaxl"
	p.CommonParam.FeeAmount = "1000"
	p.CommonParam.GasLimit = 100000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd"
	p.ToAddress = "osmo1rvs5xph4l3px2efynqsthus8p6r4exyrkr95s7"
	p.Demon = "uaxl"
	p.Amount = "100000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-1"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	t.Log(tt)
}
