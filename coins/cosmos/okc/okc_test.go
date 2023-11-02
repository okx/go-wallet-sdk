package okc

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/common/types"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/token"
	"testing"
	"time"
)

/*
注意一个公钥可能对应两个账户。

https://exchaintestrpc.okex.org/cosmos/auth/v1beta1/accounts/ex1za3n0p07wvzzunm6292d9qmdtrrxf2d9zzzzd3

// Send ordinary transaction
curl -X POST -H 'Content-Type: application/json' -d '{"jsonrpc":"2.0","id":0,"method":"broadcast_tx_sync","params":{"tx":"4QEoKBapCktWBam6ChQXYzeF/nMELk96UVTSg21YxmSppRIUjeIo7mujJj7jjn3HOHIwOy7SYFYaGQoDb2t0EhIxMDAwMDAwMDAwMDAwMDAwMDASGwoVCgNva3QSDjEwMDAwMDAwMDAwMDAwEKCNBhprCibzs80DIQLGFvfgc0gnouLHkLE78WbcJM/dttkOtSuczT9vBMqSehJBJzsJxdLJ8tqWiwDvg2gFqJGBWIxPH7WUizC283Uo51UzQGV0MtQn7YkI94VRfgubYvKF5sEBMzUTrJFJ9JXXZwEiBHRlc3Q="}}' http://18.177.169.67:26657

// Send IBC transactions
curl --location --request POST 'https://exchaintestrpc.okex.org/cosmos/tx/v1beta1/txs' \
--header 'Content-Type: application/json' \
--data-raw '{
    "tx_bytes": "Cr8BCrYBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKIAQoIdHJhbnNmZXISCWNoYW5uZWwtMBoLCgN3ZWkSBDEwMDAiKWV4MWhyMjZjeWMzMzVnN3A1ZTk0OGE3dmttd254M2ZteGZ6d2R5cnlmKi1jb3Ntb3MxcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXI3Y2t5eHYyADiA+MfvkIjRgRcSBHRlc3QSbgpOCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAsYW9+BzSCei4seQsTvxZtwkz9222Q61K5zNP28EypJ6EgQKAggBEhwKFgoDd2VpEg8xMDAwMDAwMDAwMDAwMDAQwJoMGkDrtp0i1R3DtrGIYUitz+qDd9FPYIUVS0SA3uN+AVnbPyS2UJ4Vbm4nhPlAZIwcrkavoFD1db+qP3Fz5trXWx2G",
    "mode": "BROADCAST_MODE_SYNC"
}'
*/

// Ordinary transfer
// Use ethsecp256k1 address (keccak256 for public key processing)
// The ethsecp256k1 signature format is R+S+V
func TestOkc(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	addressFrom, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addressFrom)

	addressTo, err := NewAddress("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	if err != nil {
		t.Fatal(err)
	}

	coins, _ := types.ParseDecCoins("0.1okt")
	feeCoins, _ := types.ParseDecCoins("0.00001okt")
	ad1, _ := types.AccAddressFromBech32(addressFrom)
	ad2, _ := types.AccAddressFromBech32(addressTo)
	msg := token.NewMsgTokenSend(ad1, ad2, coins)
	chainId := "exchain-65"
	memo := "test"
	msgs := make([]types.Msg, 0)
	msgs = append(msgs, msg)
	gas := uint64(100000)
	accNumber := uint64(32191948)
	seqNumber := uint64(0)

	stdTx, _ := tx.BuildStdTx(privateKeyHex, chainId, memo, msgs, feeCoins, gas, accNumber, seqNumber)
	tx, _ := tx.MarshalStdTx(stdTx)
	// 4QEoKBapCktWBam6ChQXYzeF/nMELk96UVTSg21YxmSppRIUjeIo7mujJj7jjn3HOHIwOy7SYFYaGQoDb2t0EhIxMDAwMDAwMDAwMDAwMDAwMDASGwoVCgNva3QSDjEwMDAwMDAwMDAwMDAwEKCNBhprCibzs80DIQLGFvfgc0gnouLHkLE78WbcJM/dttkOtSuczT9vBMqSehJBJzsJxdLJ8tqWiwDvg2gFqJGBWIxPH7WUizC283Uo51UzQGV0MtQn7YkI94VRfgubYvKF5sEBMzUTrJFJ9JXXZwEiBHRlc3Q=
	t.Log(tx)
}

// IBC transfer
// Use secp256k1 address (hash160 for public keys)
// Use secp256k1 signature format R+S
// The unit must be wei
func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	addressFrom, err := cosmos.NewAddress(privateKeyHex, "ex", false)
	if err != nil {
		t.Fatal(err)
	}
	// ex1hr26cyc335g7p5e948a7vkmwnx3fmxfzwdyryf
	t.Log(addressFrom)

	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "exchain-65"
	p.CommonParam.Sequence = 0
	p.CommonParam.AccountNumber = 32190628
	p.CommonParam.FeeDemon = "wei"
	p.CommonParam.FeeAmount = "100000000000000"
	p.CommonParam.GasLimit = 200000
	p.CommonParam.Memo = "test"
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = addressFrom
	p.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	p.Demon = "wei"
	p.Amount = "1000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-0"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	t.Log(tt)
}
