package kujira

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"testing"
	"time"
)

/*
https://rest.kujira.ccvalidators.com/cosmos/auth/v1beta1/accounts/kujira1rvs5xph4l3px2efynqsthus8p6r4exyr0s5utx
curl -X POST -d '{"tx_bytes":"Cr4BCrsBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKNAQoIdHJhbnNmZXISCWNoYW5uZWwtMxoOCgV1a3VqaRIFMTAwMDAiLWt1amlyYTFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjBzNXV0eCorb3NtbzFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cmtyOTVzNzIAOID4/dvm1cKCFxJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDnCTAPeltdUKTRGnHRm5QS12GCMOZ3du0gUsZYfj3OlMSBAoCCAEYARISCgwKBXVrdWppEgMyMDAQoI0GGkAXkTv/pTKApBpu19XKhQqMo9hlCSQqj8q5mId4xU727DlZDtlJLRaKJJ77UHuLUYUqkfRBZChpbjT0tuxhUcvK","mode":"BROADCAST_MODE_SYNC"}' https://rest.kujira.ccvalidators.com/cosmos/tx/v1beta1/txs
*/
func TestKujira(t *testing.T) {
	privateKeyHex := "//todo please replace your hex cosmos key"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	// kujira1rvs5xph4l3px2efynqsthus8p6r4exyr0s5utx
	fmt.Println(address)

	ret := ValidateAddress(address)
	if !ret {
		t.Fatal("ValidateAddress failed")
	}

	param := cosmos.TransferParam{}
	param.FromAddress = address
	param.ToAddress = address
	param.Demon = "ukuji"
	// 0.01 kuji
	param.Amount = "10000"
	param.CommonParam.ChainId = "kaiyo-1"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 23140
	param.CommonParam.FeeDemon = "ukuji"
	param.CommonParam.FeeAmount = "200"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	tx, _ := cosmos.Transfer(param, privateKeyHex)
	fmt.Println(tx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "//todo please replace your hex cosmos key"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "kaiyo-1"
	p.CommonParam.Sequence = 1
	p.CommonParam.AccountNumber = 23140
	p.CommonParam.FeeDemon = "ukuji"
	p.CommonParam.FeeAmount = "200"
	p.CommonParam.GasLimit = 100000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = address
	p.ToAddress = "osmo1rvs5xph4l3px2efynqsthus8p6r4exyrkr95s7"
	p.Demon = "ukuji"
	p.Amount = "10000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-3"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	t.Log(tt)
}
