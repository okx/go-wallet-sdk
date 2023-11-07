package nist_p256

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

type NISTP256Curve struct {
	Version []byte
}

func (nc *NISTP256Curve) GenerateKey() ([]byte, []byte) {

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	private := priv.D.Bytes()
	x := priv.PublicKey.X.Bytes()
	y := priv.PublicKey.Y.Bytes()

	pub := append(x, y...)
	return private, pub
}

func (nc *NISTP256Curve) GetVersion() []byte {
	return nc.Version
}

func (nc *NISTP256Curve) NewKeyFromSeed(seed []byte) ([]byte, []byte) {

	return nil, nil
}
func NewNISTP256PrivateBySeed(seed []byte) *ecdsa.PrivateKey {
	d := new(big.Int).SetBytes(seed)
	priv := new(ecdsa.PrivateKey)
	priv.D = d
	priv.Curve = elliptic.P256()
	return priv
}
