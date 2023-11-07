package ed25519

import (
	"crypto/rand"
	"golang.org/x/crypto/ed25519"
)

type Ed25519Curve struct {
	Version []byte
}

func (ec *Ed25519Curve) GenerateKey() ([]byte, []byte) {

	pub, priv, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		panic(err)
	}

	return priv.Seed(), pub
}
func (ec *Ed25519Curve) GetVersion() []byte {
	return ec.Version
}
