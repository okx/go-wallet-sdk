/*
 * Copyright 2016 Stellar Development Foundation and contributors.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file includes portions of third-party code from [https://github.com/stellar/go].
 * The original code is licensed under the Apache License 2.0.
 */

package keypair

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"github.com/okx/go-wallet-sdk/coins/stellar/strkey"
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
)

type Full struct {
	address    string
	seed       string
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func newFull(seed string) (*Full, error) {
	rawSeed, err := strkey.Decode(strkey.VersionByteSeed, seed)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(rawSeed)
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return nil, err
	}
	address, err := strkey.Encode(strkey.VersionByteAccountID, pub)
	if err != nil {
		return nil, err
	}
	return &Full{
		address:    address,
		seed:       seed,
		publicKey:  pub,
		privateKey: priv,
	}, nil
}

func newFullFromRawSeed(rawSeed [32]byte) (*Full, error) {
	seed, err := strkey.Encode(strkey.VersionByteSeed, rawSeed[:])
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(rawSeed[:])
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return nil, err
	}
	address, err := strkey.Encode(strkey.VersionByteAccountID, pub)
	if err != nil {
		return nil, err
	}
	return &Full{
		address:    address,
		seed:       seed,
		publicKey:  pub,
		privateKey: priv,
	}, nil
}

func (kp *Full) Address() string {
	return kp.address
}

// FromAddress gets the address-only representation, or public key, of this
// Full keypair.
func (kp *Full) FromAddress() *FromAddress {
	return newFromAddressWithPublicKey(kp.address, kp.publicKey)
}

func (kp *Full) Hint() (r [4]byte) {
	copy(r[:], kp.publicKey[28:])
	return
}

func (kp *Full) Seed() string {
	return kp.seed
}

func (kp *Full) Verify(input []byte, sig []byte) error {
	if len(sig) != 64 {
		return ErrInvalidSignature
	}
	if !ed25519.Verify(kp.publicKey, input, sig) {
		return ErrInvalidSignature
	}
	return nil
}

func (kp *Full) Sign(input []byte) ([]byte, error) {
	return ed25519.Sign(kp.privateKey, input), nil
}

// SignBase64 signs the input data and returns a base64 encoded string, the
// common format in which signatures are exchanged.
func (kp *Full) SignBase64(input []byte) (string, error) {
	sig, err := kp.Sign(input)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func (kp *Full) SignDecorated(input []byte) (xdr.DecoratedSignature, error) {
	sig, err := kp.Sign(input)
	if err != nil {
		return xdr.DecoratedSignature{}, err
	}

	return xdr.NewDecoratedSignature(sig, kp.Hint()), nil
}

// SignPayloadDecorated returns a decorated signature that signs for a signed payload signer where the input is the payload being signed.
func (kp *Full) SignPayloadDecorated(input []byte) (xdr.DecoratedSignature, error) {
	sig, err := kp.Sign(input)
	if err != nil {
		return xdr.DecoratedSignature{}, err
	}

	return xdr.NewDecoratedSignatureForPayload(sig, kp.Hint(), input), nil
}

func (kp *Full) Equal(f *Full) bool {
	if kp == nil && f == nil {
		return true
	}
	if kp == nil || f == nil {
		return false
	}
	return kp.seed == f.seed
}
