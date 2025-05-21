package stargaze

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
https://api.stargaze.ezstaking.io/cosmos/auth/v1beta1/accounts/stars1rlvaqq27e4c5jcnghgvdnd3739w0vvt3jafls7
curl -X POST -d '{"tx_bytes":"Cr8BCrwBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKOAQoIdHJhbnNmZXISCWNoYW5uZWwtMBoQCgZ1c3RhcnMSBjEwMDAwMCIsc3RhcnMxcmx2YXFxMjdlNGM1amNuZ2hndmRuZDM3Mzl3MHZ2dDNqYWZsczcqK29zbW8xcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXJrcjk1czcyADiAmKaauZrggRcSWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAt0bZnnEU42FukQutkuKsLe5vZuuF5/zl1VG6bGu5it6EgQKAggBGAISBBCgjQYaQDkRCRKgJfb6it3sEskv6em75zlgZxBuAKre8CxkCABkQt3lrM1TcrNVCIn7ri4tcmWHv2MBFtf+fExlBe68Wsg=","mode":"BROADCAST_MODE_SYNC"}' https://api.stargaze.ezstaking.io/cosmos/tx/v1beta1/txs
*/
func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "stars145q0tcdur4tcx2ya5cphqx96e54yflfy5w5sq9"
	require.Equal(t, expected, address)

	ret := ValidateAddress(address)
	require.True(t, ret)

	param := cosmos.TransferParam{}
	param.FromAddress = address
	param.ToAddress = address
	param.Demon = "ustars"
	param.Amount = "1000"
	param.CommonParam.ChainId = "stargaze-1"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 149288
	param.CommonParam.FeeDemon = "ustars"
	param.CommonParam.FeeAmount = "0"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	signedTx, err := cosmos.Transfer(param, privateKeyHex)
	require.Nil(t, err)
	expected = "Co8BCowBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmwKLHN0YXJzMTQ1cTB0Y2R1cjR0Y3gyeWE1Y3BocXg5NmU1NHlmbGZ5NXc1c3E5EixzdGFyczE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeTV3NXNxORoOCgZ1c3RhcnMSBDEwMDASVgpOCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAxBT6e8CldM0trsi4gzHF+saFqVG9pJXLIgwtLwUwTZ2EgQKAggBEgQQoI0GGkDfp9lAfh2uz8vIYiwNeAniMphm2uHQLEyQua4fvxUh93ZN0ZuwnLINvbKjD7uX68UaQWKdfnJUEuL6O5mVgupn"
	require.Equal(t, expected, signedTx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "stargaze-1"
	p.CommonParam.Sequence = 2
	p.CommonParam.AccountNumber = 149288
	p.CommonParam.FeeDemon = "ustars"
	p.CommonParam.FeeAmount = "0"
	p.CommonParam.GasLimit = 100000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "stars1rlvaqq27e4c5jcnghgvdnd3739w0vvt3jafls7"
	p.ToAddress = "osmo1rvs5xph4l3px2efynqsthus8p6r4exyrkr95s7"
	p.Demon = "ustars"
	p.Amount = "100000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-0"
	p.TimeOutInSeconds = 1738641357
	tt, _ := cosmos.IbcTransfer(p, privateKeyHex)
	require.Equal(t, "Cr8BCrwBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKOAQoIdHJhbnNmZXISCWNoYW5uZWwtMBoQCgZ1c3RhcnMSBjEwMDAwMCIsc3RhcnMxcmx2YXFxMjdlNGM1amNuZ2hndmRuZDM3Mzl3MHZ2dDNqYWZsczcqK29zbW8xcnZzNXhwaDRsM3B4MmVmeW5xc3RodXM4cDZyNGV4eXJrcjk1czcyADiAhKfeo6O5kBgSWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAxBT6e8CldM0trsi4gzHF+saFqVG9pJXLIgwtLwUwTZ2EgQKAggBGAISBBCgjQYaQCk/kUYe+3uVrpTGIz4s4J37oafA3wsh+ab9pKc7gihLCBtxS8waCUeSwLiSCJaUyuzhnRN5/O5/z95XS2OYl/U=", tt)
}
