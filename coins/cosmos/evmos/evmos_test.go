package evmos

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// https://rest.bd.evmos.org:1317/cosmos/auth/v1beta1/accounts/evmos1rvs5xph4l3px2efynqsthus8p6r4exyrue82uy
// curl -X POST -d '{"tx_bytes":"CswBCskBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKbAQoIdHJhbnNmZXISCWNoYW5uZWwtMxobCgZhZXZtb3MSETEwMDAwMDAwMDAwMDAwMDAwIixldm1vczF5YzRxNnN2c2w5eHk5ZzJncGxnbmxweHdobnpyM3k3M3dmczB4aCotY29zbW9zMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyN2NreXh2MgA4gOCBuJnk3oEXEn0KWQpPCigvZXRoZXJtaW50LmNyeXB0by52MS5ldGhzZWNwMjU2azEuUHViS2V5EiMKIQOcJMA96W11QpNEacdGblBLXYYIw5nd27SBSxlh+Pc6UxIECgIIARgBEiAKGgoGYWV2bW9zEhA0MDAwMDAwMDAwMDAwMDAwEMCaDBpBxW58piSUv3r+MRmwIe3xllBkJxgF5I0QIlHLjym8amohsZWmyYCUzaux/pO2RNbB4K9VmJ2m8Y3/56w6Gpeh5wE=","mode":"BROADCAST_MODE_SYNC"}' https://rest.bd.evmos.org:1317/cosmos/tx/v1beta1/txs
func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "evmos1jl389rqgh59lhf33j20pp082aj8utjtp3a0lt5"
	require.Equal(t, expected, address)

	param := cosmos.TransferParam{}
	param.FromAddress = "evmos1yc4q6svsl9xy9g2gplgnlpxwhnzr3y73wfs0xh"
	param.ToAddress = "evmos1yc4q6svsl9xy9g2gplgnlpxwhnzr3y73wfs0xh"
	param.Demon = "aevmos"
	param.Amount = "10000000000000000" // 18
	param.CommonParam.ChainId = "evmos_9001-2"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 2091572
	param.CommonParam.FeeDemon = "aevmos"
	param.CommonParam.FeeAmount = "3500000000000000"
	param.CommonParam.GasLimit = 140000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	signedTx, err := cosmos.TransferAction(param, privateKeyHex, true)
	require.Nil(t, err)
	assert.Equal(t, "CpwBCpkBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnkKLGV2bW9zMXljNHE2c3ZzbDl4eTlnMmdwbGdubHB4d2huenIzeTczd2ZzMHhoEixldm1vczF5YzRxNnN2c2w5eHk5ZzJncGxnbmxweHdobnpyM3k3M3dmczB4aBobCgZhZXZtb3MSETEwMDAwMDAwMDAwMDAwMDAwEnsKVwpPCigvZXRoZXJtaW50LmNyeXB0by52MS5ldGhzZWNwMjU2azEuUHViS2V5EiMKIQMQU+nvApXTNLa7IuIMxxfrGhalRvaSVyyIMLS8FME2dhIECgIIARIgChoKBmFldm1vcxIQMzUwMDAwMDAwMDAwMDAwMBDgxQgaQdDLRmTF5ZANbBk2D2ZNjXOonztGsOJxH2o0G+I9WSuhYDLjcQ/qJd89Y1wQ7YWbBZEyUM4IgeogadLf6j6xKT4B", signedTx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "evmos_9001-2"
	p.CommonParam.Sequence = 1
	p.CommonParam.AccountNumber = 2091572
	p.CommonParam.FeeDemon = "aevmos"
	p.CommonParam.FeeAmount = "4000000000000000"
	p.CommonParam.GasLimit = 200000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "evmos1yc4q6svsl9xy9g2gplgnlpxwhnzr3y73wfs0xh"
	p.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	p.Demon = "aevmos"
	p.Amount = "10000000000000000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-3"
	p.TimeOutInSeconds = 1738641357
	signedIbcTx, _ := cosmos.IbcTransferAction(p, privateKeyHex, true)
	assert.Equal(t, "CswBCskBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKbAQoIdHJhbnNmZXISCWNoYW5uZWwtMxobCgZhZXZtb3MSETEwMDAwMDAwMDAwMDAwMDAwIixldm1vczF5YzRxNnN2c2w5eHk5ZzJncGxnbmxweHdobnpyM3k3M3dmczB4aCotY29zbW9zMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyN2NreXh2MgA4gISn3qOjuZAYEn0KWQpPCigvZXRoZXJtaW50LmNyeXB0by52MS5ldGhzZWNwMjU2azEuUHViS2V5EiMKIQMQU+nvApXTNLa7IuIMxxfrGhalRvaSVyyIMLS8FME2dhIECgIIARgBEiAKGgoGYWV2bW9zEhA0MDAwMDAwMDAwMDAwMDAwEMCaDBpBurCw/OtDiKAVDEJ6XSOVFqhZx+yOkrcU8OfblD8NyhQvOH0FrhPTqAbX+XkSU24qO9rPTvW0gRrcB59peNAGAAE=", signedIbcTx)
}
func TestNewPubAddress(t *testing.T) {
	addr, err := NewPubAddress("0450365da35a0a7cd463fe421b7bc78552d3ca47361af11fa5155f4b016f459bfaf80e63125fb4ee5d50f8b10cab5b00a413ff301352ee8b49c19ec3c917231b1f")
	assert.NoError(t, err)
	assert.Equal(t, "evmos1efvrk87rh334z0uv9sm3kensc033plzc5alcq9", addr)
}
