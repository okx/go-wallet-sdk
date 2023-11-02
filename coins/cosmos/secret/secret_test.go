package secret

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Check account details
// https://api.scrt.network/cosmos/auth/v1beta1/accounts/secret1rvs5xph4l3px2efynqsthus8p6r4exyruazdms
// curl -X POST -d '{"tx_bytes":"CsABCr0BCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKPAQoIdHJhbnNmZXISCWNoYW5uZWwtMBoOCgV1c2NydBIFMTAwMDAiLXNlY3JldDFybHZhcXEyN2U0YzVqY25naGd2ZG5kMzczOXcwdnZ0M3l5MnR4biotY29zbW9zMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyN2NreXh2MgA4gIi76+T/34EXEmcKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQLdG2Z5xFONhbpELrZLirC3ub2brhef85dVRumxruYrehIECgIIARgCEhMKDQoFdXNjcnQSBDIwMDAQoI0GGkBo2WG6Qt6Na+sM7grAPuQNQQwLuWt6640i+kynQXQdMVzC4TBB06As/G4RQjarK5NX9HVZo7H3W5GTQmvSdr58","mode":"BROADCAST_MODE_SYNC"}' https://api.roninventures.io/cosmos/tx/v1beta1/txs

func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "secret145q0tcdur4tcx2ya5cphqx96e54yflfyzhhykg"
	require.Equal(t, expected, address)

	param := cosmos.TransferParam{}
	param.FromAddress = address
	param.ToAddress = address
	param.Demon = "uscrt"
	param.Amount = "1000"
	param.CommonParam.ChainId = "secret-4"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 247039
	param.CommonParam.FeeDemon = "uscrt"
	param.CommonParam.FeeAmount = "1250"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	signedTx, err := cosmos.Transfer(param, privateKeyHex)
	require.Nil(t, err)
	expected = "CpABCo0BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm0KLXNlY3JldDE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeXpoaHlrZxItc2VjcmV0MTQ1cTB0Y2R1cjR0Y3gyeWE1Y3BocXg5NmU1NHlmbGZ5emhoeWtnGg0KBXVzY3J0EgQxMDAwEmUKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQMQU+nvApXTNLa7IuIMxxfrGhalRvaSVyyIMLS8FME2dhIECgIIARITCg0KBXVzY3J0EgQxMjUwEKCNBhpAcv1on1fdBc/nVnLBHGPVh7ydrHwnt/3JM0bv3qOnzAMAvwMcm/Fh4HEzfp+aQAJgvWZ97RQCnVRHeO4d8sCSHQ=="
	require.Equal(t, expected, signedTx)
}

func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "secret-4"
	p.CommonParam.Sequence = 2
	p.CommonParam.AccountNumber = 236422
	p.CommonParam.FeeDemon = "uscrt"
	p.CommonParam.FeeAmount = "2000"
	p.CommonParam.GasLimit = 100000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "secret1rlvaqq27e4c5jcnghgvdnd3739w0vvt3yy2txn"
	p.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	p.Demon = "uscrt"
	p.Amount = "10000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-0"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	signedIBCTx, err := cosmos.IbcTransfer(p, privateKeyHex)
	require.Nil(t, err)
	t.Log("signedIBCTx : ", signedIBCTx)
}
