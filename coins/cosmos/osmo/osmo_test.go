package osmo

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
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
	msgLockTokens := &tx.MsgLockTokens{}
	msgLockTokens.Owner = "osmo1dr8rh6pj78f6wzddjyruyj3ga3r0tsjk5cv3hl"
	msgLockTokens.Duration = 10 * time.Second
	msgLockTokens.Coins = types.NewCoins(types.NewCoin("osmo", types.NewIntFromUint64(1)))
	b, err := msgLockTokens.Marshal()
	require.Nil(t, err)
	expected := "0a2b6f736d6f31647238726836706a37386636777a64646a797275796a33676133723074736a6b35637633686c1202080a1a090a046f736d6f120131"
	require.Equal(t, expected, hex.EncodeToString(b))

	msgLockTokens2 := &tx.MsgLockTokens{}
	err = msgLockTokens2.Unmarshal(b)
	require.Nil(t, err)
	require.Equal(t, msgLockTokens, msgLockTokens2)
}

func TestAddress(t *testing.T) {
	privateKey := "2b9960183a7e94ed3686e758e5f853d4d4ddd1d5053525d0ce5e747ba69e9da3"
	address, err := NewAddress(privateKey)
	require.NoError(t, err)
	require.Equal(t, "osmo1rm6ql2ss7q4ech7wnaghzpczheafj0jkyjgnxf", address)
}

func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
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
	signedTx, err := cosmos.Transfer(param, privateKeyHex)
	require.Nil(t, err)
	expected := "Co4BCosBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmsKK29zbW8xcmx2YXFxMjdlNGM1amNuZ2hndmRuZDM3Mzl3MHZ2dDN3NmRqZGESK29zbW8xbHlqeGs0dDgzNXlqNnU4bDJtZzZhNnQydjl4M25qN3VsYWxqejIaDwoFdW9zbW8SBjEwMDAwMBJWCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAESBBCgjQYaQNuwNAIRJ7gPqXo72lgQ0hzHtJiMU8QnwjaoLnTgc/fhDn2U4/BKiJqqmHHLdkP778xJuo5VVPPmJXIuL/bd5vk="
	require.Equal(t, expected, signedTx)
}

// swap evmos to osmo
func TestSwap(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
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
	signedSwapTx, err := SwapExactAmountIn(param, privateKeyHex)
	require.Nil(t, err)
	expected := "CswBCskBCiovb3Ntb3Npcy5nYW1tLnYxYmV0YTEuTXNnU3dhcEV4YWN0QW1vdW50SW4SmgEKK29zbW8xbHlqeGs0dDgzNXlqNnU4bDJtZzZhNnQydjl4M25qN3VsYWxqejISSQjSBRJEaWJjLzZBRTk4ODgzRDRENUQ1RkY5RTUwRDcxMzBGMTMwNURBMkZGQTBDNjUyRDFERDlDMTIzNjU3QzZCNEVCMkRGOEEaDgoFdW9zbW8SBTEwMDAwIhAzODU0MTU0MTgwODEzMDE4ElYKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQMQU+nvApXTNLa7IuIMxxfrGhalRvaSVyyIMLS8FME2dhIECgIIARIEEJChDxpApLK4a9ykW+AMBWTOfVx/T0BETXQnkeoDncNl1P8VSutbJsUqJ4PzH5bGM2UNqGxJemGMq9Ot9QdMrC6BO6OhEg=="
	require.Equal(t, expected, signedSwapTx)
}

func TestSignMessage(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	data := "{\n  \"chain_id\": \"osmosis-1\",\n  \"account_number\": \"584406\",\n  \"sequence\": \"1\",\n  \"fee\": {\n    \"gas\": \"250000\",\n    \"amount\": [\n      {\n        \"denom\": \"uosmo\",\n        \"amount\": \"0\"\n      }\n    ]\n  },\n  \"msgs\": [\n    {\n      \"type\": \"osmosis/gamm/swap-exact-amount-in\",\n      \"value\": {\n        \"sender\": \"osmo1lyjxk4t835yj6u8l2mg6a6t2v9x3nj7ulaljz2\",\n        \"routes\": [\n          {\n            \"poolId\": \"722\",\n            \"tokenOutDenom\": \"ibc/6AE98883D4D5D5FF9E50D7130F1305DA2FFA0C652D1DD9C123657C6B4EB2DF8A\"\n          }\n        ],\n        \"tokenIn\": {\n          \"denom\": \"uosmo\",\n          \"amount\": \"10000\"\n        },\n        \"tokenOutMinAmount\": \"3854154180813018\"\n      }\n    }\n  ],\n  \"memo\": \"\"\n}"
	signedMessage, _, err := cosmos.SignMessage(data, privateKeyHex)
	require.Nil(t, err)
	expected := "CswBCskBCiovb3Ntb3Npcy5nYW1tLnYxYmV0YTEuTXNnU3dhcEV4YWN0QW1vdW50SW4SmgEKK29zbW8xbHlqeGs0dDgzNXlqNnU4bDJtZzZhNnQydjl4M25qN3VsYWxqejISSQjSBRJEaWJjLzZBRTk4ODgzRDRENUQ1RkY5RTUwRDcxMzBGMTMwNURBMkZGQTBDNjUyRDFERDlDMTIzNjU3QzZCNEVCMkRGOEEaDgoFdW9zbW8SBTEwMDAwIhAzODU0MTU0MTgwODEzMDE4ElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQMQU+nvApXTNLa7IuIMxxfrGhalRvaSVyyIMLS8FME2dhIECgIIARgBEgQQkKEPGkDgnemI608eStirAt1Nb5EZtbMSrytjFuEoExu+ShCD32XFAVZMWN6EvM/pIMXEy74pPUMcVY0I5Dx0kRFYLE7Q"
	require.Equal(t, expected, signedMessage)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
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
	p.TimeOutInSeconds = 1738641357
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	require.Equal(t, "Cr8BCrwBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKOAQoIdHJhbnNmZXISCWNoYW5uZWwtMBoPCgV1b3NtbxIGMTAwMDAwIitvc21vMWx5anhrNHQ4MzV5ajZ1OGwybWc2YTZ0MnY5eDNuajd1bGFsanoyKi1jb3Ntb3MxcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXI3Y2t5eHYyADiAhKfeo6O5kBgSWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAxBT6e8CldM0trsi4gzHF+saFqVG9pJXLIgwtLwUwTZ2EgQKAggBGAYSBBCgjQYaQGne0PQuSL5xfoa7gBzp0gC+C9worDbDYc4fidrsgaD4QSm0oVEXAZj+b7kPAboAxyC1nZAm9NQgRWdkoqIzGNc=", tt)
}
