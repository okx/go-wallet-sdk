package cronos

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"testing"
	"time"
)

// https://testnet-croeseid-4.crypto.org:1317/cosmos/auth/v1beta1/accounts/tcro1rvs5xph4l3px2efynqsthus8p6r4exyrgkhe6v
// curl -X POST -d '{"tx_bytes":"CsIBCr8BCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKRAQoIdHJhbnNmZXISCWNoYW5uZWwtMxoUCghiYXNldGNybxIIMTAwMDAwMDAiK3Rjcm8xcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXJna2hlNnYqK3Rjcm8xcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXJna2hlNnYyADiA6PzV5KHcgRcSagpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA5wkwD3pbXVCk0Rpx0ZuUEtdhgjDmd3btIFLGWH49zpTEgQKAggBGAISFgoQCghiYXNldGNybxIEMjUwMBCgjQYaQJTJDNWvwrVsslTM8IDh0tH8Eww+FvK7K+b3bbKdkn1VLEDgOW7gHAMN7E+tF8GsZnGV8Zxo3qhzKQG/BfwYlIE=","mode":"BROADCAST_MODE_SYNC"}' https://testnet-croeseid-4.crypto.org:1317/cosmos/tx/v1beta1/txs
// https://crypto.org/explorer/croeseid4/tx/A8F4F4953C3EF658079D538F006B5C583E55F08CE2F06662AD3199140ADD3D2D

// https://mainnet.crypto.org:1317/cosmos/auth/v1beta1/accounts/cro1rvs5xph4l3px2efynqsthus8p6r4exyrxr7a6a
// curl -X POST -d '{"tx_bytes":"Cr8BCrwBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKOAQoIdHJhbnNmZXISCmNoYW5uZWwtMTAaEQoHYmFzZWNybxIGMTAwMDAwIipjcm8xcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXJ4cjdhNmEqK29zbW8xcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXJrcjk1czcyADiAxK/DjLPggRcSagpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA5wkwD3pbXVCk0Rpx0ZuUEtdhgjDmd3btIFLGWH49zpTEgQKAggBGAESFgoQCgdiYXNlY3JvEgUyMDAwMBCgwh4aQLjCjqEX0jLLasNyVjWd4euYAlMFeC2mgKtFgQV4AXDKUOefdrBL9u2o2hJpdG28+aQZEjsenU7YPAIkG4r+of8=","mode":"BROADCAST_MODE_SYNC"}' https://mainnet.crypto.org:1317/cosmos/tx/v1beta1/txs
func TestTransfer(t *testing.T) {
	privateKeyHex := "//todo please replace your hex cosmos key"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	// cro1rvs5xph4l3px2efynqsthus8p6r4exyrxr7a6a
	fmt.Println(address)

	param := cosmos.TransferParam{}
	param.FromAddress = address
	param.ToAddress = address
	param.Demon = "basecro"
	param.Amount = "1000000"
	param.CommonParam.ChainId = "crypto-org-chain-mainnet-1"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 554047
	param.CommonParam.FeeDemon = "basecro"
	param.CommonParam.FeeAmount = "20000"
	param.CommonParam.GasLimit = 500000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	tt, _ := cosmos.Transfer(param, privateKeyHex)
	t.Log(tt)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "//todo please replace your hex cosmos key"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "crypto-org-chain-mainnet-1"
	p.CommonParam.Sequence = 1
	p.CommonParam.AccountNumber = 554047
	p.CommonParam.FeeDemon = "basecro"
	p.CommonParam.FeeAmount = "20000"
	p.CommonParam.GasLimit = 500000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "cro1rvs5xph4l3px2efynqsthus8p6r4exyrxr7a6a"
	p.ToAddress = "osmo1rvs5xph4l3px2efynqsthus8p6r4exyrkr95s7"
	p.Demon = "basecro"
	p.Amount = "100000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-10"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	t.Log(tt)
}
