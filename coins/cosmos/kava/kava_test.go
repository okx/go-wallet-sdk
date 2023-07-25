package kava

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"testing"
	"time"
)

// Check account details
// https://api.data.kava.io/cosmos/auth/v1beta1/accounts/kava1m7mutcn7h3uccjhd5q7e8adxkl7wny59739vuq
// curl -X POST -d '{"tx_bytes":"Cr0BCroBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKMAQoIdHJhbnNmZXISCWNoYW5uZWwtMBoNCgV1a2F2YRIEMTAwMCIra2F2YTFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cnpkemVzdCotY29zbW9zMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyN2NreXh2MgA4gJDco9HW34EXEmQKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQOcJMA96W11QpNEacdGblBLXYYIw5nd27SBSxlh+Pc6UxIECgIIARISCgwKBXVrYXZhEgM1MDAQ4MUIGkB9RYvHfuZ8t+2hXHEeMgWPoKPPny62KFuuedHOzCgYSE8yWEf3r5KwnkiIJZRD0cDcso4PZbEyiwgvkmdvrnLR","mode":"BROADCAST_MODE_SYNC"}' https://api.data.kava.io/cosmos/tx/v1beta1/txs

func TestKava(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	privateKeyHex2 := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	if address != "kava145q0tcdur4tcx2ya5cphqx96e54yflfyu8hsan" {
		t.Fatal("NewAddress failed", address)
	}

	address2, err := NewAddress(privateKeyHex2)
	if err != nil {
		t.Fatal(err)
	}
	ret := ValidateAddress(address)
	if !ret {
		t.Fatal("ValidateAddress failed")
	}

	param := cosmos.TransferParam{}
	param.FromAddress = address
	param.ToAddress = address2
	param.Demon = "ukava"
	param.Amount = "100000"
	param.CommonParam.ChainId = "kava_2222-10"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 2206349
	param.CommonParam.FeeDemon = "ukava"
	param.CommonParam.FeeAmount = "7000"
	param.CommonParam.GasLimit = 140000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	tx, _ := cosmos.Transfer(param, privateKeyHex)
	fmt.Println(tx)
	// 94763169DDD220F2109AB4A5C619C60D583C33FD767BA9E413ED811ACA68AA48
	if tx != "Co4BCosBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmsKK2thdmExNDVxMHRjZHVyNHRjeDJ5YTVjcGhxeDk2ZTU0eWZsZnl1OGhzYW4SK2thdmExNDVxMHRjZHVyNHRjeDJ5YTVjcGhxeDk2ZTU0eWZsZnl1OGhzYW4aDwoFdWthdmESBjEwMDAwMBJlCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAESEwoNCgV1a2F2YRIENzAwMBDgxQgaQPJCdq9SjuciSGv/yDfZEVXdHWuQeK6Eh0c7/7zcklN1Hh6z06bbtVtl/uByCdoTtUu+gYmZM+MuLqv9JPu3YAw=" {
		t.Error("build transfer tx failed, tx: ", tx)
	}
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}

	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "kava_2222-10"
	p.CommonParam.Sequence = 0
	p.CommonParam.AccountNumber = 2211629
	p.CommonParam.FeeDemon = "ukava"
	p.CommonParam.FeeAmount = "500"
	p.CommonParam.GasLimit = 140000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = address
	p.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	p.Demon = "ukava"
	p.Amount = "1000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-0"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	t.Log(tt)
}
