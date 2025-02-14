package transactions

import (
	"encoding/base64"
	"github.com/okx/go-wallet-sdk/coins/helium/keypair"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	alice    = "148d8KTRcKA5JKPekBcKFd4KfvprvFRpjGtivhtmRmnZ8MFYnP3"
	bob      = "13Lqwnbh427csevUveZF9n3ra1LnVYQug31RFeENaYgXuK2s8UC"
	TmpSig   = make([]byte, 64)
	from     = keypair.NewAddressable(bob)
	to       = keypair.NewAddressable(alice)
	toAmount = map[string]uint64{alice: 10}
	v2       = NewPaymentV2Tx(from, toAmount, 0, 1, "hnt", true, TmpSig)
	kp       = keypair.NewKeypairFromHex(1, "f5e029dd6cca805047ca64e131c0a6cf3bf45c7ad03a7a1e7681963c9b1f3043")
)

func TestPaymentV2Tx_Serialize(t *testing.T) {
	v2.Fee = 3
	v2.SetSignature([]byte("bob`s signature"))
	data, err := v2.Serialize()
	require.NoError(t, err)
	expected := "wgFhCiEBNHM7ubKaxcXmamFwaK3pL5W2MrH8hG4iu+k5KoRZIswSJwohAZxlnXI8wegQpy54996vRzaofxDvj8/IAQC1Myfn7kmkEAogARgDIAEqD2JvYmBzIHNpZ25hdHVyZQ=="
	require.Equal(t, expected, base64.StdEncoding.EncodeToString(data))
}

func TestPaymentV2Tx_SignTransaction(t *testing.T) {
	v2Tx, err := v2.BuildTransaction(true)
	require.NoError(t, err)
	sig, err := kp.Sign(v2Tx)
	require.NoError(t, err)
	expected := "5i4xKIeYNbhVAsAF5m8IXzgfssait5ja1MghBnd4mgoY9VsJ9FUSUwiNPcx5sFOdHUgY1MCdCt12sCdPKYV8AQ=="
	require.Equal(t, expected, base64.StdEncoding.EncodeToString(sig))
}
