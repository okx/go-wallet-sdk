package ethsecp256k1

import (
	"bytes"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/amino"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/tendermint"
	"golang.org/x/crypto/sha3"
)

// Amino encoding names
const (
	// PrivKeyName defines the amino encoding name for the EthSecp256k1 private key
	PrivKeyName = "ethermint/PrivKeyEthSecp256k1"
	// PubKeyName defines the amino encoding name for the EthSecp256k1 public key
	PubKeyName = "ethermint/PubKeyEthSecp256k1"
)

// ----------------------------------------------------------------------------
// secp256k1 Private Key

var _ tendermint.PrivKey = PrivKey{}

// PrivKey defines a type alias for an ecdsa.PrivateKey that implements
// Tendermint's PrivateKey interface.
type PrivKey []byte

// PubKey returns the ECDSA private key's public key.
func (privkey PrivKey) PubKey() tendermint.PubKey {
	_, ecPub := btcec.PrivKeyFromBytes(privkey)
	return PubKey(ecPub.SerializeCompressed())
}

// Bytes returns the raw ECDSA private key bytes.
func (privkey PrivKey) Bytes() []byte {
	return amino.GCodec.MustMarshalBinaryBare(privkey)
}

// Equals returns true if two ECDSA private keys are equal and false otherwise.
func (privkey PrivKey) Equals(other tendermint.PrivKey) bool {
	if other, ok := other.(PrivKey); ok {
		return bytes.Equal(privkey.Bytes(), other.Bytes())
	}

	return false
}

// ----------------------------------------------------------------------------
// secp256k1 Public Key

var _ tendermint.PubKey = (*PubKey)(nil)

// PubKey defines a type alias for an ecdsa.PublicKey that implements Tendermint's PubKey
// interface. It represents the 33-byte compressed public key format.
type PubKey []byte

// Address returns the address of the ECDSA public key.
// The function will panic if the public key is invalid.
func (key PubKey) Address() tendermint.Address {
	ecPubKey, _ := btcec.ParsePubKey(key)
	pubBytes := ecPubKey.SerializeUncompressed()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:])
	addressByte := hash.Sum(nil)
	return tendermint.Address(addressByte[12:])
}

// Bytes returns the raw bytes of the ECDSA public key.
// The function panics if the key cannot be marshaled to bytes.
func (key PubKey) Bytes() []byte {
	bz, err := amino.GCodec.MarshalBinaryBare(key)
	if err != nil {
		panic(err)
	}
	return bz
}

// Equals returns true if two ECDSA public keys are equal and false otherwise.
func (key PubKey) Equals(other tendermint.PubKey) bool {
	if other, ok := other.(PubKey); ok {
		return bytes.Equal(key.Bytes(), other.Bytes())
	}

	return false
}
