/**
Author： https://github.com/xssnick/tonutils-go
*/

package wallet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSeedWithPassword(t *testing.T) {
	seed := NewSeedWithPassword("123")
	fmt.Println(seed)
	prv, err := FromSeedWithPassword(seed, "123", V3)
	require.NoError(t, err)
	fmt.Println(prv)
	fmt.Println(hex.EncodeToString(prv.PrivateKey()))
	prv2, err := FromSeedWithPassword(seed, "123", V4R2)
	require.NoError(t, err)
	fmt.Println(hex.EncodeToString(prv2.PrivateKey()))
	assert.True(t, bytes.Equal(prv.PrivateKey(), prv2.PrivateKey()))

	_, err = FromSeedWithPassword(seed, "1234", V3)
	require.NotNil(t, err)

	_, err = FromSeedWithPassword(seed, "", V3)
	require.NotNil(t, err)

	_, err = FromSeedWithPassword([]string{"birth", "core"}, "", V3)
	require.NotNil(t, err)

	seed = NewSeed()
	seed[7] = "wat"
	_, err = FromSeedWithPassword(seed, "", V3)
	require.NotNil(t, err)

	seedNoPass := NewSeed()

	_, err = FromSeed(seedNoPass, V3)
	require.NoError(t, err)

	_, err = FromSeedWithPassword(seedNoPass, "123", V3)
	require.NotNil(t, err)
}

//func TestCheckoutOfficialSDK(t *testing.T) {
//	for i := 0; i < 100000; i++ {
//		seedNoPass := NewSeed()
//		w, err := FromSeed(seedNoPass, V4R2)
//		require.NoError(t, err)
//		addr := w.Address()
//		addr.SetBounce(false)
//		//t.Log(addr.String())/
//		tw, err := tonwallet.FromSeed(nil, seedNoPass, tonwallet.V4R2)
//		require.NoError(t, err)
//		twAddr := tw.Address()
//		twAddr.SetBounce(false)
//		require.Equal(t, addr.String(), twAddr.String())
//		t.Log("i = :", i)
//	}
//}
