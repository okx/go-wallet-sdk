package axelar

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/require"
	"testing"
)

// https://k8s-testnet-axelarco-c0dd71f944-b4c8da2f814e7b8f.elb.us-east-2.amazonaws.com:1317/cosmos/auth/v1beta1/accounts/axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd
// curl -X POST -d '{"tx_bytes":"Cr4BCrsBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKNAQoIdHJhbnNmZXISCWNoYW5uZWwtMRoOCgR1YXhsEgYxMDAwMDAiLWF4ZWxhcjFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjZrcXZkZCorb3NtbzFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cmtyOTVzNzIAOICg2b6uxNuBFxJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDnCTAPeltdUKTRGnHRm5QS12GCMOZ3du0gUsZYfj3OlMSBAoCCAEYARISCgwKBHVheGwSBDEwMDAQoI0GGkBEmWixhNVWDgKgUJ4SVB/vYiWu69sdqmAp52ZDVmsZHlSuXNK/hISNeS7sTA1b5DdCEYHcYl+ieq6Z6ubOl2Z0","mode":"BROADCAST_MODE_SYNC"}' http://k8s-testnet-axelarco-c0dd71f944-b4c8da2f814e7b8f.elb.us-east-2.amazonaws.com:1317/cosmos/tx/v1beta1/txs
func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "axelar145q0tcdur4tcx2ya5cphqx96e54yflfyyu49q4"
	require.Equal(t, expected, address)

	param := cosmos.TransferParam{}
	param.FromAddress = "axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd"
	param.ToAddress = "axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd"
	param.Demon = "uaxl"
	param.Amount = "100000"
	param.CommonParam.ChainId = "axelar-testnet-lisbon-3"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 24569
	param.CommonParam.FeeDemon = "uaxl"
	param.CommonParam.FeeAmount = "1000"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	signedTx, err := cosmos.Transfer(param, privateKeyHex)
	require.Nil(t, err)
	expected = "CpEBCo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KLWF4ZWxhcjFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjZrcXZkZBItYXhlbGFyMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyNmtxdmRkGg4KBHVheGwSBjEwMDAwMBJkCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAESEgoMCgR1YXhsEgQxMDAwEKCNBhpAl6KVyR7I1pmnvra7mjADszcKrka7K7mM5EzLBWyDCUZvE+dXA9ZszSYzI4zZkuTVLqAIOA1QZVkvkxzprLzxXw=="
	require.Equal(t, expected, signedTx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "axelar-testnet-lisbon-3"
	p.CommonParam.Sequence = 1
	p.CommonParam.AccountNumber = 24569
	p.CommonParam.FeeDemon = "uaxl"
	p.CommonParam.FeeAmount = "1000"
	p.CommonParam.GasLimit = 100000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd"
	p.ToAddress = "osmo1rvs5xph4l3px2efynqsthus8p6r4exyrkr95s7"
	p.Demon = "uaxl"
	p.Amount = "100000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-1"
	p.TimeOutInSeconds = 1738641357
	signedIBCTx, err := cosmos.IbcTransfer(p, privateKeyHex)
	require.Nil(t, err)
	expected := "Cr4BCrsBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKNAQoIdHJhbnNmZXISCWNoYW5uZWwtMRoOCgR1YXhsEgYxMDAwMDAiLWF4ZWxhcjFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjZrcXZkZCorb3NtbzFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cmtyOTVzNzIAOICEp96jo7mQGBJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAEYARISCgwKBHVheGwSBDEwMDAQoI0GGkBOiY0QJ0/qRWhN661WWjDnpQR/hGrtbSjoUYYnADrRN2yNr5PqWnLWCZeR3vRnt/4e+SVPSBkxGB/O7uKeWB+k"
	require.Equal(t, expected, signedIBCTx)
}
