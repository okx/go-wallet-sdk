package atom

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types/ibc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestGetSignedTx(t *testing.T) {
	pk, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	require.Nil(t, err)
	k, _ := btcec.PrivKeyFromBytes(pk)
	address, err := NewAddress(k)
	require.Nil(t, err)
	expected := "cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5"
	require.Equal(t, expected, address)

	ret := ValidateAddress(address)
	require.True(t, ret)

	chainId := "cosmoshub-4"
	from := "cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5"
	to := "cosmos1jun53r4ycc8g2v6tffp4cmxjjhv6y7ntat62wn"
	demon := "uatom"
	memo := "memo"
	amount := big.NewInt(10000)
	sequence := 0
	accountNumber := 623151
	feeAmount := big.NewInt(10)
	gasLimit := 100
	hexStr, err := SignStart(chainId, from, to, demon, memo, amount, 0, uint64(sequence), uint64(accountNumber), feeAmount, uint64(gasLimit), k)
	require.Nil(t, err)
	expected = "0a97010a8e010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e64126e0a2d636f736d6f733134357130746364757234746378327961356370687178393665353479666c6679716a72647435122d636f736d6f73316a756e3533723479636338673276367466667034636d786a6a68763679376e7461743632776e1a0e0a057561746f6d1205313030303012046d656d6f12610a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21031053e9ef0295d334b6bb22e20cc717eb1a16a546f692572c8830b4bc14c1367612040a020801120f0a0b0a057561746f6d1202313010641a0b636f736d6f736875622d3420af8426"
	require.Equal(t, expected, hexStr)

	signedStr, err := Sign(hexStr, k)
	require.Nil(t, err)

	expected = "57fa782b2982e9119d285a2cbbb2e4a8d8c08a1c7b419bad2cf4a6d219046f2c6d9d42424417b23d551e35cc60f461f5caaf5739334b59a9e9458719e8920296"
	require.Equal(t, expected, signedStr)

	signedTx, err := SignEnd(hexStr, signedStr)
	require.Nil(t, err)
	expected = "CpcBCo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KLWNvc21vczE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeXFqcmR0NRItY29zbW9zMWp1bjUzcjR5Y2M4ZzJ2NnRmZnA0Y214ampodjZ5N250YXQ2MnduGg4KBXVhdG9tEgUxMDAwMBIEbWVtbxJhCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAESDwoLCgV1YXRvbRICMTAQZBpAV/p4KymC6RGdKFosu7LkqNjAihx7QZutLPSm0hkEbyxtnUJCRBeyPVUeNcxg9GH1yq9XOTNLWanpRYcZ6JIClg=="
	require.Equal(t, expected, signedTx)
}

func TestGetSignedTransaction(t *testing.T) {
	pk, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	require.Nil(t, err)

	k, _ := btcec.PrivKeyFromBytes(pk)
	param := cosmos.TransferParam{}
	param.FromAddress = "cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5"
	param.ToAddress = "cosmos1jun53r4ycc8g2v6tffp4cmxjjhv6y7ntat62wn"
	param.Demon = "uatom"
	param.Amount = "10000"
	param.CommonParam.ChainId = "cosmoshub-4"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 623151
	param.CommonParam.FeeDemon = "uatom"
	param.CommonParam.FeeAmount = "10"
	param.CommonParam.GasLimit = 100
	param.CommonParam.Memo = "memo"
	param.CommonParam.TimeoutHeight = 0
	doc, err := cosmos.GetRawTransaction(param, hex.EncodeToString(k.PubKey().SerializeCompressed()))
	require.Nil(t, err)
	t.Log(doc)

	signature, err := cosmos.SignRawTransaction(doc, k)
	require.Nil(t, err)

	signedTransaction, err := cosmos.GetSignedTransaction(doc, signature)
	require.Nil(t, err)

	expected := "CpcBCo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KLWNvc21vczE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeXFqcmR0NRItY29zbW9zMWp1bjUzcjR5Y2M4ZzJ2NnRmZnA0Y214ampodjZ5N250YXQ2MnduGg4KBXVhdG9tEgUxMDAwMBIEbWVtbxJhCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAESDwoLCgV1YXRvbRICMTAQZBpAV/p4KymC6RGdKFosu7LkqNjAihx7QZutLPSm0hkEbyxtnUJCRBeyPVUeNcxg9GH1yq9XOTNLWanpRYcZ6JIClg=="
	require.Equal(t, expected, signedTransaction)
}

// Check account details
// https://api.cosmos.network/cosmos/auth/v1beta1/accounts/cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv
// curl -X POST -d '{"tx_bytes":"CpABCo0BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm0KLWNvc21vczFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjdja3l4dhItY29zbW9zMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyN2NreXh2Gg0KBXVhdG9tEgQxMDAwEmYKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQOcJMA96W11QpNEacdGblBLXYYIw5nd27SBSxlh+Pc6UxIECgIIfxgCEhIKDAoFdWF0b20SAzEzMBCgjQYaQA04G6nhx6Zo8uYBHKhyw46t7RkvxLwDO0XfkRG3hVRRDmCg6xl+61FhXe3x2A/temH/hGsIt1bjs37vcDQAgg4=","mode":"BROADCAST_MODE_SYNC"}' https://api.cosmos.network/cosmos/tx/v1beta1/txs
func TestGetSignedJsonTransaction(t *testing.T) {
	pk, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	require.Nil(t, err)
	k, _ := btcec.PrivKeyFromBytes(pk)
	param := cosmos.TransferParam{}
	param.FromAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	param.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	param.Demon = "uatom"
	param.Amount = "1000"
	param.CommonParam.ChainId = "cosmoshub-4"
	param.CommonParam.Sequence = 2
	param.CommonParam.AccountNumber = 1225716
	param.CommonParam.FeeDemon = "uatom"
	param.CommonParam.FeeAmount = "130"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	doc, err := cosmos.GetRawJsonTransaction(param)
	require.Nil(t, err)

	signature, err := cosmos.SignRawJsonTransaction(doc, k)
	require.Nil(t, err)

	publicKey := hex.EncodeToString(k.PubKey().SerializeCompressed())
	signedJsonTransaction, err := cosmos.GetSignedJsonTransaction(doc, publicKey, signature)
	require.Nil(t, err)
	expected := "CpABCo0BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm0KLWNvc21vczFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjdja3l4dhItY29zbW9zMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyN2NreXh2Gg0KBXVhdG9tEgQxMDAwEmYKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQMQU+nvApXTNLa7IuIMxxfrGhalRvaSVyyIMLS8FME2dhIECgIIfxgCEhIKDAoFdWF0b20SAzEzMBCgjQYaQNbZuALq96PyACNflFUnnS5kd/fCZniLitsmRhCR092UHIJYGEUugnG89Be+v6BbWe0E0jPRqPfw36thibn6ix0="
	require.Equal(t, expected, signedJsonTransaction)
}
func TestInjectiveAddress(t *testing.T) {
	address, err := cosmos.NewAddress("9fd01e033b56c22acd8ecd1318300147a73bd48c29d195c53f3e7ea415d78d86", "inj", true)
	assert.NoError(t, err)
	assert.Equal(t, "inj1glnhek2ml2397fx9pdw8s4vt8dklqsfh508eaj", address)

	address, err = cosmos.NewAddress("9892a1712215fe4740a774581dd0f1966817d21cdf1b23022d28aac2e4761530", "osmo", false)
	assert.NoError(t, err)
	assert.Equal(t, "osmo1rlvaqq27e4c5jcnghgvdnd3739w0vvt3w6djda", address)
}

func TestPubHex2AnyHex(t *testing.T) {
	res, err := cosmos.PubHex2AnyHex("04627540e5288e988813fdd8d2b0267a3343f9509899c342d02a12c6cc89056ea12df17a4e9d17eadd2d198877c224f3ef0ef8b81306a4540ae41e7845702ac9e9", true)
	assert.NoError(t, err)
	assert.Equal(t, "0a2103627540e5288e988813fdd8d2b0267a3343f9509899c342d02a12c6cc89056ea1", res)
}
func TestConvert2AnyPubKey(t *testing.T) {
	res, err := cosmos.Convert2AnyPubKey("0a2103627540e5288e988813fdd8d2b0267a3343f9509899c342d02a12c6cc89056ea1", false, true)
	assert.NoError(t, err)
	assert.Equal(t, "04627540e5288e988813fdd8d2b0267a3343f9509899c342d02a12c6cc89056ea12df17a4e9d17eadd2d198877c224f3ef0ef8b81306a4540ae41e7845702ac9e9", res)
}

func TestMakeTransactionWithSignDoc(t *testing.T) {
	body := "0x0a96040a242f636f736d7761736d2e7761736d2e76312e4d736745786563757465436f6e747261637412ed030a2b6f736d6f313935797372386e73763833716d7673326474636c617735656634646c637861676d3430653233123f6f736d6f31616a3261717a3034796674737365687433376d686775787874717161637330743376743333327536677472397a3472326c787971356836397a671aea027b22657865637574655f737761705f6f7065726174696f6e73223a7b226d696e696d756d5f72656365697665223a223239333734303537222c22726f75746573223a5b7b226f666665725f616d6f756e74223a2231303030303030222c226f7065726174696f6e73223a5b7b22745f665f6d5f73776170223a7b22706f6f6c5f6964223a3537332c226f666665725f61737365745f696e666f223a7b226e61746976655f746f6b656e223a7b2264656e6f6d223a22756f736d6f227d7d2c2261736b5f61737365745f696e666f223a7b226e61746976655f746f6b656e223a7b2264656e6f6d223a226962632f34453534343443333536313043433736464339344537463738383642393331323131373543323832363244444644444536463834453832424632343235343532227d7d7d7d5d7d5d2c22746f223a226f736d6f313935797372386e73763833716d7673326474636c617735656634646c637861676d3430653233227d7d2a100a05756f736d6f120731303030303030"
	auth := "0x0a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21023201511a7aba3b0c8c31b8db13f85f1fc91dda848c3e262917d61551aa1cc61312040a020801180d12130a0d0a05756f736d6f12043632353010aa931b"
	chainId := "osmosis-1"
	accountNumber := uint64(799205)
	res, err := cosmos.MakeTransactionWithSignDoc(body, auth, chainId, accountNumber)
	assert.NoError(t, err)
	assert.Equal(t, "0a99040a96040a242f636f736d7761736d2e7761736d2e76312e4d736745786563757465436f6e747261637412ed030a2b6f736d6f313935797372386e73763833716d7673326474636c617735656634646c637861676d3430653233123f6f736d6f31616a3261717a3034796674737365687433376d686775787874717161637330743376743333327536677472397a3472326c787971356836397a671aea027b22657865637574655f737761705f6f7065726174696f6e73223a7b226d696e696d756d5f72656365697665223a223239333734303537222c22726f75746573223a5b7b226f666665725f616d6f756e74223a2231303030303030222c226f7065726174696f6e73223a5b7b22745f665f6d5f73776170223a7b22706f6f6c5f6964223a3537332c226f666665725f61737365745f696e666f223a7b226e61746976655f746f6b656e223a7b2264656e6f6d223a22756f736d6f227d7d2c2261736b5f61737365745f696e666f223a7b226e61746976655f746f6b656e223a7b2264656e6f6d223a226962632f34453534343443333536313043433736464339344537463738383642393331323131373543323832363244444644444536463834453832424632343235343532227d7d7d7d5d7d5d2c22746f223a226f736d6f313935797372386e73763833716d7673326474636c617735656634646c637861676d3430653233227d7d2a100a05756f736d6f12073130303030303012670a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21023201511a7aba3b0c8c31b8db13f85f1fc91dda848c3e262917d61551aa1cc61312040a020801180d12130a0d0a05756f736d6f12043632353010aa931b1a096f736d6f7369732d3120e5e330", res)
}
func TestGetRawTransaction(t *testing.T) {
	publicKey := "02dc02bb89e72e0b8e596c2276e341734b98e35c52d8d0462147ebdf9a4b0d9a3c"
	param := cosmos.IbcTransferParam{
		CommonParam: cosmos.CommonParam{
			ChainId:       "crypto-org-chain-mainnet-1",
			Sequence:      2,
			AccountNumber: 2091572,
			FeeDemon:      "basecro",
			FeeAmount:     "12500",
			GasLimit:      500000,
			Memo:          "",
			TimeoutHeight: 0,
		},
		FromAddress:      "cosmos1w97axsme4h65u63z9malcnrgppr6sn6ynqr8v8",
		ToAddress:        "osmo1rvs5xph4l3px2efynqsthus8p6r4exyrkr95s7",
		Demon:            "basecro",
		Amount:           "100000",
		SourcePort:       "transfer",
		SourceChannel:    "channel-10",
		TimeOutHeight:    ibc.Height{},
		TimeOutInSeconds: 300,
	}
	res, err := cosmos.GetRawTransaction(param, publicKey)
	assert.NoError(t, err)
	assert.Equal(t, "0abf010abc010a292f6962632e6170706c69636174696f6e732e7472616e736665722e76312e4d73675472616e73666572128e010a087472616e73666572120a6368616e6e656c2d31301a110a076261736563726f1206313030303030222d636f736d6f73317739376178736d65346836357536337a396d616c636e726770707236736e36796e71723876382a2b6f736d6f3172767335787068346c337078326566796e7173746875733870367234657879726b723935733732003880f092cbdd08126a0a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a2102dc02bb89e72e0b8e596c2276e341734b98e35c52d8d0462147ebdf9a4b0d9a3c12040a020801180212160a100a076261736563726f1205313235303010a0c21e1a1a63727970746f2d6f72672d636861696e2d6d61696e6e65742d3120b4d47f", res)
}
func TestGetRawJsonTransaction(t *testing.T) {
	param := cosmos.IbcTransferParam{
		CommonParam: cosmos.CommonParam{
			ChainId:       "crypto-org-chain-mainnet-1",
			Sequence:      2,
			AccountNumber: 2091572,
			FeeDemon:      "basecro",
			FeeAmount:     "12500",
			GasLimit:      500000,
			Memo:          "",
			TimeoutHeight: 0,
		},
		FromAddress:      "cosmos1x9q0esjv8j7nyul043arwkquxmjzxx8cgfylvc",
		ToAddress:        "osmo1rvs5xph4l3px2efynqsthus8p6r4exyrkr95s7",
		Demon:            "basecro",
		Amount:           "100000",
		SourcePort:       "transfer",
		SourceChannel:    "channel-10",
		TimeOutHeight:    ibc.Height{},
		TimeOutInSeconds: 1708498314,
	}
	res, err := cosmos.GetRawJsonTransaction(param)
	assert.NoError(t, err)
	assert.Equal(t, "{\"account_number\":\"2091572\",\"chain_id\":\"crypto-org-chain-mainnet-1\",\"fee\":{\"amount\":[{\"amount\":\"12500\",\"denom\":\"basecro\"}],\"gas\":\"500000\"},\"memo\":\"\",\"msgs\":[{\"type\":\"cosmos-sdk/MsgTransfer\",\"value\":{\"receiver\":\"osmo1rvs5xph4l3px2efynqsthus8p6r4exyrkr95s7\",\"sender\":\"cosmos1x9q0esjv8j7nyul043arwkquxmjzxx8cgfylvc\",\"source_channel\":\"channel-10\",\"source_port\":\"transfer\",\"timeout_height\":{},\"timeout_timestamp\":1708498314000000000,\"token\":{\"amount\":\"100000\",\"denom\":\"basecro\"}}}],\"sequence\":\"2\"}", res)
}

func TestSignDoc(t *testing.T) {
	privateKey := "d34d39ecbc2494e6c5c94421d00d80c668839b134583ba5b82fa79383c96144f"
	body := "0x0a96040a242f636f736d7761736d2e7761736d2e76312e4d736745786563757465436f6e747261637412ed030a2b6f736d6f313935797372386e73763833716d7673326474636c617735656634646c637861676d3430653233123f6f736d6f31616a3261717a3034796674737365687433376d686775787874717161637330743376743333327536677472397a3472326c787971356836397a671aea027b22657865637574655f737761705f6f7065726174696f6e73223a7b226d696e696d756d5f72656365697665223a223239333734303537222c22726f75746573223a5b7b226f666665725f616d6f756e74223a2231303030303030222c226f7065726174696f6e73223a5b7b22745f665f6d5f73776170223a7b22706f6f6c5f6964223a3537332c226f666665725f61737365745f696e666f223a7b226e61746976655f746f6b656e223a7b2264656e6f6d223a22756f736d6f227d7d2c2261736b5f61737365745f696e666f223a7b226e61746976655f746f6b656e223a7b2264656e6f6d223a226962632f34453534343443333536313043433736464339344537463738383642393331323131373543323832363244444644444536463834453832424632343235343532227d7d7d7d5d7d5d2c22746f223a226f736d6f313935797372386e73763833716d7673326474636c617735656634646c637861676d3430653233227d7d2a100a05756f736d6f120731303030303030"
	auth := "0x0a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21023201511a7aba3b0c8c31b8db13f85f1fc91dda848c3e262917d61551aa1cc61312040a020801180d12130a0d0a05756f736d6f12043632353010aa931b"
	chainId := "osmosis-1"
	accountNumber := uint64(799205)
	tx, sig, err := cosmos.SignDoc(body, auth, privateKey, chainId, accountNumber)
	assert.NoError(t, err)
	assert.Equal(t, "CpkECpYECiQvY29zbXdhc20ud2FzbS52MS5Nc2dFeGVjdXRlQ29udHJhY3QS7QMKK29zbW8xOTV5c3I4bnN2ODNxbXZzMmR0Y2xhdzVlZjRkbGN4YWdtNDBlMjMSP29zbW8xYWoyYXF6MDR5ZnRzc2VodDM3bWhndXh4dHFxYWNzMHQzdnQzMzJ1Nmd0cjl6NHIybHh5cTVoNjl6ZxrqAnsiZXhlY3V0ZV9zd2FwX29wZXJhdGlvbnMiOnsibWluaW11bV9yZWNlaXZlIjoiMjkzNzQwNTciLCJyb3V0ZXMiOlt7Im9mZmVyX2Ftb3VudCI6IjEwMDAwMDAiLCJvcGVyYXRpb25zIjpbeyJ0X2ZfbV9zd2FwIjp7InBvb2xfaWQiOjU3Mywib2ZmZXJfYXNzZXRfaW5mbyI6eyJuYXRpdmVfdG9rZW4iOnsiZGVub20iOiJ1b3NtbyJ9fSwiYXNrX2Fzc2V0X2luZm8iOnsibmF0aXZlX3Rva2VuIjp7ImRlbm9tIjoiaWJjLzRFNTQ0NEMzNTYxMENDNzZGQzk0RTdGNzg4NkI5MzEyMTE3NUMyODI2MkRERkRERTZGODRFODJCRjI0MjU0NTIifX19fV19XSwidG8iOiJvc21vMTk1eXNyOG5zdjgzcW12czJkdGNsYXc1ZWY0ZGxjeGFnbTQwZTIzIn19KhAKBXVvc21vEgcxMDAwMDAwEmcKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIyAVEaero7DIwxuNsT+F8fyR3ahIw+JikX1hVRqhzGExIECgIIARgNEhMKDQoFdW9zbW8SBDYyNTAQqpMbGkCoLOiXcmMInlIqvi7RZAFLe0Ls+HnRH8h25mQYVJtwNnBgIeF59wPceY5tJa1SKcgIYwX74qPr6wOIiDfSC/id", tx)
	assert.Equal(t, "qCzol3JjCJ5SKr4u0WQBS3tC7Ph50R/IduZkGFSbcDZwYCHhefcD3HmObSWtUinICGMF++Kj6+sDiIg30gv4nQ==", sig)
}
