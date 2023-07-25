package sei

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"testing"
	"time"
)

func TestAddress(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	if address != "sei145q0tcdur4tcx2ya5cphqx96e54yflfyd7jmd4" {
		t.Errorf("NewAddress failed, want sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57, get %s", address)
	}

	ret := ValidateAddress(address)
	if !ret {
		t.Fatal("ValidateAddress failed")
	}
}

// txHash: https://sei.explorers.guru/transaction/08D076CFE1903AB697974E7CB756F5E9A3344FF07220892CE2A75EBB29494435
// https://rest.atlantic-2.seinetwork.io/cosmos/auth/v1beta1/accounts/sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57
// curl -X POST -d '{"tx_bytes":"CosBCogBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmgKKnNlaTFzOTV6dnB4d3hjcjB5a2RrajN5bXNjcmV2ZGFtN3d2czI0ZGs1NxIqc2VpMXVyZGVkZWVqMGZkNGtzbHpuM3VxNnM4bW5kaDh3dDd1c2s2YTR6Gg4KBHVzZWkSBjEwMDAwMBJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECWzs64TTLQ3sGP88eAUtzXoHtGHUYauDmYZWgBLyYUesSBAoCCAEYARISCgwKBHVzZWkSBDEwMDAQoI0GGkAPt3BsqAL807wgpPtKQdF8mYsPwM52HjRaScsc27rIdh30d6JxWnu9Zy1Tm9funsAYIOtStq7GKTfekctaIRK/","mode":"BROADCAST_MODE_SYNC"}' https://rest.atlantic-2.seinetwork.io/cosmos/tx/v1beta1/txs
func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	// sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57
	fmt.Println(address)

	param := cosmos.TransferParam{}
	param.FromAddress = "sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57"
	param.ToAddress = "sei1urdedeej0fd4kslzn3uq6s8mndh8wt7usk6a4z"
	param.Demon = "usei"
	param.Amount = "100000"
	param.CommonParam.ChainId = "atlantic-2"
	param.CommonParam.Sequence = 1
	param.CommonParam.AccountNumber = 4050874
	param.CommonParam.FeeDemon = "usei"
	param.CommonParam.FeeAmount = "1000"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	tt, _ := cosmos.Transfer(param, privateKeyHex)
	t.Log(tt)
}

// https://sei.explorers.guru/transaction/E07194819858ED6C8BF355AF55A1F57E6346C45E3A9C5CDEE2DFFDEA3992BE5B
// curl -X POST -d '{"tx_bytes":"CrsBCrgBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKKAQoIdHJhbnNmZXISCWNoYW5uZWwtNxoMCgR1c2VpEgQxMDAwIipzZWkxczk1enZweHd4Y3IweWtka2ozeW1zY3JldmRhbTd3dnMyNGRrNTcqLWF4ZWxhcjFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjZrcXZkZDIAOIDs+ImekvirFxJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECWzs64TTLQ3sGP88eAUtzXoHtGHUYauDmYZWgBLyYUesSBAoCCAEYBBISCgwKBHVzZWkSBDEwMDAQoI0GGkB2KRZkV1uDATVA9Tk7rpqCeSu6b6j1POSrOKQLX2Vijxoj5hpIRm32sKykpf+O+ED7QiZwzF24+WTowx9eAMwb","mode":"BROADCAST_MODE_SYNC"}' https://rest.atlantic-2.seinetwork.io/cosmos/tx/v1beta1/txs
func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "atlantic-2"
	p.CommonParam.Sequence = 4
	p.CommonParam.AccountNumber = 4050874
	p.CommonParam.FeeDemon = "usei"
	p.CommonParam.FeeAmount = "1000"
	p.CommonParam.GasLimit = 100000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57"
	p.ToAddress = "axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd"
	p.Demon = "usei"
	p.Amount = "1000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-7"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	t.Log(tt)
}
