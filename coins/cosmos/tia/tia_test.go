package tia

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAddress(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "celestia145q0tcdur4tcx2ya5cphqx96e54yflfy3cja3e"
	require.Equal(t, expected, address)

	ret := ValidateAddress(address)
	require.True(t, ret)
}

func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "celestia145q0tcdur4tcx2ya5cphqx96e54yflfy3cja3e"
	require.Equal(t, expected, address)
	ret := ValidateAddress(address)
	require.True(t, ret)

	param := cosmos.TransferParam{}
	param.FromAddress = address
	param.ToAddress = address
	param.Demon = "utia"
	param.Amount = "1000"
	param.CommonParam.ChainId = "mocha-4"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 149288
	param.CommonParam.FeeDemon = "utia"
	param.CommonParam.FeeAmount = "20000"
	param.CommonParam.GasLimit = 200000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	signedTx, err := cosmos.Transfer(param, privateKeyHex)
	require.Nil(t, err)
	expected = "CpMBCpABChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnAKL2NlbGVzdGlhMTQ1cTB0Y2R1cjR0Y3gyeWE1Y3BocXg5NmU1NHlmbGZ5M2NqYTNlEi9jZWxlc3RpYTE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeTNjamEzZRoMCgR1dGlhEgQxMDAwEmUKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQMQU+nvApXTNLa7IuIMxxfrGhalRvaSVyyIMLS8FME2dhIECgIIARITCg0KBHV0aWESBTIwMDAwEMCaDBpAH2vJG5BisaeW04Jrw/62m9UWo/Me5+abwdEoo9Z+gTwcvvPp7TPSzOewQq5s5neFY4mhV75Z4XEzZU0aIwZdVg=="
	require.Equal(t, expected, signedTx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "mocha-4"
	p.CommonParam.Sequence = 2
	p.CommonParam.AccountNumber = 149288
	p.CommonParam.FeeDemon = "utia"
	p.CommonParam.FeeAmount = "20000"
	p.CommonParam.GasLimit = 200000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "celestia145q0tcdur4tcx2ya5cphqx96e54yflfy3cja3e"
	p.ToAddress = "cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5"
	p.Demon = "utia"
	p.Amount = "100000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-0"
	p.TimeOutInSeconds = 1738641357
	signedIBCTx, err := cosmos.IbcTransfer(p, privateKeyHex)
	require.Nil(t, err)
	require.Equal(t, "CsIBCr8BCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKRAQoIdHJhbnNmZXISCWNoYW5uZWwtMBoOCgR1dGlhEgYxMDAwMDAiL2NlbGVzdGlhMTQ1cTB0Y2R1cjR0Y3gyeWE1Y3BocXg5NmU1NHlmbGZ5M2NqYTNlKi1jb3Ntb3MxNDVxMHRjZHVyNHRjeDJ5YTVjcGhxeDk2ZTU0eWZsZnlxanJkdDUyADiAhKfeo6O5kBgSZwpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAxBT6e8CldM0trsi4gzHF+saFqVG9pJXLIgwtLwUwTZ2EgQKAggBGAISEwoNCgR1dGlhEgUyMDAwMBDAmgwaQGHNPuIOz6QOeZBStlwyJ4pWlzNa4eyTkNwnt2LtjaZII9IhlpH6y/838tBMgHd/Iw9joAMjAG7eSpOHsbCmQpI=", signedIBCTx)
}
