package kujira

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/require"
	"testing"
)

// https://rest.kujira.ccvalidators.com/cosmos/auth/v1beta1/accounts/kujira1rvs5xph4l3px2efynqsthus8p6r4exyr0s5utx
// curl -X POST -d '{"tx_bytes":"Cr4BCrsBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKNAQoIdHJhbnNmZXISCWNoYW5uZWwtMxoOCgV1a3VqaRIFMTAwMDAiLWt1amlyYTFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjBzNXV0eCorb3NtbzFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cmtyOTVzNzIAOID4/dvm1cKCFxJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDnCTAPeltdUKTRGnHRm5QS12GCMOZ3du0gUsZYfj3OlMSBAoCCAEYARISCgwKBXVrdWppEgMyMDAQoI0GGkAXkTv/pTKApBpu19XKhQqMo9hlCSQqj8q5mId4xU727DlZDtlJLRaKJJ77UHuLUYUqkfRBZChpbjT0tuxhUcvK","mode":"BROADCAST_MODE_SYNC"}' https://rest.kujira.ccvalidators.com/cosmos/tx/v1beta1/txs
func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "kujira145q0tcdur4tcx2ya5cphqx96e54yflfy36p4x7"
	require.Equal(t, expected, address)

	ret := ValidateAddress(address)
	require.True(t, ret)

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
	signedTx, err := cosmos.Transfer(param, privateKeyHex)
	require.Nil(t, err)
	expected = "CpEBCo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KLWt1amlyYTE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeTM2cDR4NxIta3VqaXJhMTQ1cTB0Y2R1cjR0Y3gyeWE1Y3BocXg5NmU1NHlmbGZ5MzZwNHg3Gg4KBXVrdWppEgUxMDAwMBJkCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAESEgoMCgV1a3VqaRIDMjAwEKCNBhpAvsOP9pFW3dZZzYkrBn7viJp9/xy8G/JB8LNI0IJ6qtMDoaKqcyfMDys+vzBieSNY0oX2v7oDXJIui1oPe0H+og=="
	require.Equal(t, expected, signedTx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)

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
	p.TimeOutInSeconds = 1738641357
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	require.Equal(t, "Cr4BCrsBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKNAQoIdHJhbnNmZXISCWNoYW5uZWwtMxoOCgV1a3VqaRIFMTAwMDAiLWt1amlyYTE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeTM2cDR4Nyorb3NtbzFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cmtyOTVzNzIAOICEp96jo7mQGBJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAEYARISCgwKBXVrdWppEgMyMDAQoI0GGkD+a5C8Zyc1tcbphNy7MJc5fsp1Mx3OpZbIQ6uPzpi8xWrFnc17uLuK2aqyxQAQzhOx8mk0bHFNh1qHqxg/favN", tt)
}
