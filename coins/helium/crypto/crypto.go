/**
Author： https://github.com/hecodev007/block_sign
*/

package crypto

import (
	ed25519_helium "github.com/okx/go-wallet-sdk/coins/helium/crypto/ed25519"
	nist_p256 "github.com/okx/go-wallet-sdk/coins/helium/crypto/nist-p256"
)

type Curves interface {
	GenerateKey() ([]byte, []byte)
	GetVersion() []byte
}

func NewCurve(version int) Curves {
	var c Curves
	if version == 0 {
		nc := &nist_p256.NISTP256Curve{Version: []byte{byte(version)}}
		c = nc
	} else if version == 1 {
		ec := &ed25519_helium.Ed25519Curve{Version: []byte{byte(version)}}
		c = ec
	}
	return c
}
