package osmo

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	time "time"

	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/okx/go-wallet-sdk/coins/cosmos/osmo/tx"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types"
)

// Check account details
// https://lcd.osmosis.zone/cosmos/auth/v1beta1/accounts/osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2
// curl -X POST -d '{"tx_bytes":"Cr8BCrwBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKOAQoIdHJhbnNmZXISCWNoYW5uZWwtMBoPCgV1b3NtbxIGMTAwMDAwIitvc21vMWx5anhrNHQ4MzV5ajZ1OGwybWc2YTZ0MnY5eDNuajd1bGFsanoyKi1jb3Ntb3MxcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXI3Y2t5eHYyADiA5MH+6e3f/xYSWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohArgrNHdKxbe8aTAVVpSqtrcrHh47o6U3Bk4Pmwfxs2T+EgQKAggBGAYSBBCgjQYaQL+pAulAo/+OMqspiT/qHpXb71QLBfartDHj6gAGJFKUStImSdT4ltezhrcvKjVBKEVlf9+D+g197ewQ8wtGiJY=","mode":"BROADCAST_MODE_SYNC"}' https://lcd.osmosis.zone/cosmos/tx/v1beta1/txs

// curl -X POST -d '{"mode":"BROADCAST_MODE_BLOCK","tx_bytes":{"msg":[{"type":"/osmosis.gamm.v1beta1.MsgSwapExactAmountIn","value":{"routes":[{"poolId":"722","tokenOutDenom":"ibc/6AE98883D4D5D5FF9E50D7130F1305DA2FFA0C652D1DD9C123657C6B4EB2DF8A"}],"sender":"osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2","tokenIn":{"amount":"10000","denom":"uosmo"},"tokenOutMinAmount":"3854154180813018"}}],"fee":{"amount":[{"amount":"0","denom":"uosmo"}],"gas":"250000"},"signatures":[{"account_number":"584406","sequence":"1","signature":"kt90HM/P/ZioWum/r4g9/qx1i+/6xMDkXgjBA0zyHQZMILzQCSB1RSUMPQ1Ktvh4FabrjyXq5asaiWNYAzh4oA==","pub_key":{"type":"tendermint/PubKeySecp256k1","value":"ArgrNHdKxbe8aTAVVpSqtrcrHh47o6U3Bk4Pmwfxs2T+"}}]},"sequences":["1"]}' -H 'content-type:application/json;' https://lcd.osmosis.zone/cosmos/tx/v1beta1/txs
func TestProto(t *testing.T) {
	msgLockTokens := tx.MsgLockTokens{}
	msgLockTokens.Owner = "osmo1dr8rh6pj78f6wzddjyruyj3ga3r0tsjk5cv3hl"
	msgLockTokens.Duration = 10 * time.Second
	msgLockTokens.Coins = types.NewCoins(types.NewCoin("osmo", types.NewIntFromUint64(1)))
	b, _ := msgLockTokens.Marshal()
	t.Log(hex.EncodeToString(b))

	msgLockTokens2 := tx.MsgLockTokens{}
	_ = msgLockTokens2.Unmarshal(b)
	t.Log(msgLockTokens2)
}

func TestAddress(t *testing.T) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	t.Log(hex.EncodeToString(b))
	address, err := NewAddress(hex.EncodeToString(b))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(address)
}

func TestOsmo(t *testing.T) {
	privateKeyHex := "//todo please replace your hex key"

	param := cosmos.TransferParam{}
	param.FromAddress = "osmo1rlvaqq27e4c5jcnghgvdnd3739w0vvt3w6djda"
	param.ToAddress = "osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2"
	param.Demon = "uosmo"
	param.Amount = "100000"
	param.CommonParam.ChainId = "osmosis-1"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 584405
	param.CommonParam.FeeDemon = "uosmo"
	param.CommonParam.FeeAmount = "0"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	tt, _ := cosmos.Transfer(param, privateKeyHex)
	// Co4BCosBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmsKK29zbW8xcmx2YXFxMjdlNGM1amNuZ2hndmRuZDM3Mzl3MHZ2dDN3NmRqZGESK29zbW8xbHlqeGs0dDgzNXlqNnU4bDJtZzZhNnQydjl4M25qN3VsYWxqejIaDwoFdW9zbW8SBjEwMDAwMBJWCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEC3RtmecRTjYW6RC62S4qwt7m9m64Xn/OXVUbpsa7mK3oSBAoCCAESBBCgjQYaQJFWJYsuX7xmVKjGO9ugbQJ9bs9AY8JIH2rn9mN+/SXoKV+20jqKdDYeQvPTPaQ9Y7TtgDUtX22tJfB1GcZUpgM=
	// B2F55A971239FC6C747885A3BF46C258CBA8279C0F3463897EA5E0ED5DBB54EB
	t.Log(tt)
}

// swap evmos to osmo
func TestOsmoSwap(t *testing.T) {
	privateKeyHex := "//todo please replace your hex key"
	// //todo please replace your hex key osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2

	param := SwapExactAmountInParam{}
	param.Sender = "osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2"
	// evmos - osmo
	param.PoolId = 722
	// evmos
	param.ToDemon = "ibc/6AE98883D4D5D5FF9E50D7130F1305DA2FFA0C652D1DD9C123657C6B4EB2DF8A"
	param.FromDemon = "uosmo"
	param.FromAmount = "10000"
	param.MinToAmount = "3854154180813018"

	param.CommonParam.ChainId = "osmosis-1"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 584406
	param.CommonParam.FeeDemon = "uosmo"
	param.CommonParam.FeeAmount = "0"
	param.CommonParam.GasLimit = 250000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0

	// CswBCskBCiovb3Ntb3Npcy5nYW1tLnYxYmV0YTEuTXNnU3dhcEV4YWN0QW1vdW50SW4SmgEKK29zbW8xbHlqeGs0dDgzNXlqNnU4bDJtZzZhNnQydjl4M25qN3VsYWxqejISSQjSBRJEaWJjLzZBRTk4ODgzRDRENUQ1RkY5RTUwRDcxMzBGMTMwNURBMkZGQTBDNjUyRDFERDlDMTIzNjU3QzZCNEVCMkRGOEEaDgoFdW9zbW8SBTEwMDAwIhAzODU0MTU0MTgwODEzMDE4ElYKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQK4KzR3SsW3vGkwFVaUqra3Kx4eO6OlNwZOD5sH8bNk/hIECgIIARIEEJChDxpAzM8IlsvQkgCTTbTB5BQGBEt8N9ZzhHAJmASzfVI4KBIYlRMeMEKTbFX6lkVzwyvKxi1zvWa04w84Lk+sGj1D0Q==
	// 0A2CA0827F37283727B35B4FDDC4B2F3C1A384B5B6093EC58D40DAAA4F344E15
	tt, _ := SwapExactAmountIn(param, privateKeyHex)
	t.Log(tt)
}

func TestSignMessage(t *testing.T) {
	privateKeyHex := "//todo please replace your hex key"
	// //todo please replace your hex key osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2
	data := "{\n  \"chain_id\": \"osmosis-1\",\n  \"account_number\": \"584406\",\n  \"sequence\": \"1\",\n  \"fee\": {\n    \"gas\": \"250000\",\n    \"amount\": [\n      {\n        \"denom\": \"uosmo\",\n        \"amount\": \"0\"\n      }\n    ]\n  },\n  \"msgs\": [\n    {\n      \"type\": \"osmosis/gamm/swap-exact-amount-in\",\n      \"value\": {\n        \"sender\": \"osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2\",\n        \"routes\": [\n          {\n            \"poolId\": \"722\",\n            \"tokenOutDenom\": \"ibc/6AE98883D4D5D5FF9E50D7130F1305DA2FFA0C652D1DD9C123657C6B4EB2DF8A\"\n          }\n        ],\n        \"tokenIn\": {\n          \"denom\": \"uosmo\",\n          \"amount\": \"10000\"\n        },\n        \"tokenOutMinAmount\": \"3854154180813018\"\n      }\n    }\n  ],\n  \"memo\": \"\"\n}"
	tt, _ := cosmos.SignMessage(data, privateKeyHex)
	t.Log(tt)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "//todo please replace your hex key"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "osmosis-1"
	p.CommonParam.Sequence = 6
	p.CommonParam.AccountNumber = 584406
	p.CommonParam.FeeDemon = "uosmo"
	p.CommonParam.FeeAmount = "0"
	p.CommonParam.GasLimit = 100000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2"
	p.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	p.Demon = "uosmo"
	p.Amount = "100000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-0"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	t.Log(tt)
}
