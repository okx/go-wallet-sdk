/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"
	"math/big"
)

const (
	encIterations = 32768
	encKeyLen     = 32
)

var (
	// Digest is an alias for blake2b checksum algorithm
	Digest = blake2b.Sum256
)

func ecPrivateKeyFromBytes(b []byte, curve elliptic.Curve) (key *ecdsa.PrivateKey, err error) {
	k := new(big.Int).SetBytes(b)
	curveOrder := curve.Params().N
	if k.Cmp(curveOrder) >= 0 {
		return nil, fmt.Errorf("tezos: invalid private key for curve %s", curve.Params().Name)
	}

	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
		D: k,
	}

	// https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/crypto/ecdsa/ecdsa.go;l=149
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())
	return priv, nil
}

// ecNormalizeSignature ensures strict compliance with the EC spec by returning
// S mod n for the appropriate keys curve.
func ecNormalizeSignature(r, s *big.Int, c elliptic.Curve) (*big.Int, *big.Int) {
	r = new(big.Int).Set(r)
	s = new(big.Int).Set(s)

	order := c.Params().N
	quo := new(big.Int).Quo(order, new(big.Int).SetInt64(2))
	if s.Cmp(quo) > 0 {
		s = s.Sub(order, s)
	}
	return r, s
}

func ecSign(sk *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, sk, hash)
	if err != nil {
		return nil, err
	}
	r, s = ecNormalizeSignature(r, s, sk.Curve)
	buf := make([]byte, 64)
	r.FillBytes(buf[:32])
	s.FillBytes(buf[32:])
	return buf, nil
}

func decryptPrivateKey(enc []byte, fn PassphraseFunc) ([]byte, error) {
	if fn == nil {
		return nil, ErrPassphrase
	}
	passphrase, err := fn()
	if err != nil {
		return nil, err
	}
	if len(passphrase) == 0 {
		return nil, ErrPassphrase
	}

	salt, box := enc[:8], enc[8:]
	secretboxKey := pbkdf2.Key(passphrase, salt, encIterations, encKeyLen, sha512.New)
	var (
		tmp   [32]byte
		nonce [24]byte // implicitly 0x00..
	)
	copy(tmp[:], secretboxKey)
	dec, ok := secretbox.Open(nil, box, &nonce, &tmp)
	if !ok {
		return nil, fmt.Errorf("tezos: private key decrypt failed")
	}
	return dec, nil
}
