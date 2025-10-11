package tendermint

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/emresenyuva/go-wallet-sdk/coins/cosmos/okc/tx/amino"
	"golang.org/x/crypto/ripemd160"
)

// -------------------------------------
const (
	PrivKeyAminoName = "tendermint/PrivKeySecp256k1"
	PubKeyAminoName  = "tendermint/PubKeySecp256k1"
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterInterface((*PubKey)(nil), nil)
	cdc.RegisterConcrete(PubKeySecp256k1{},
		PubKeyAminoName, nil)

	cdc.RegisterInterface((*PrivKey)(nil), nil)
	cdc.RegisterConcrete(PrivKeySecp256k1{},
		PrivKeyAminoName, nil)
}

// -------------------------------------
var _ PrivKey = PrivKeySecp256k1{}

// PrivKeySecp256k1 implements PrivKey.
type PrivKeySecp256k1 [32]byte

// Bytes marshalls the private key using amino encoding.
func (privKey PrivKeySecp256k1) Bytes() []byte {
	return cdc.MustMarshalBinaryBare(privKey)
}

// PubKey performs the point-scalar multiplication from the privKey on the
// generator point to get the pubkey.
func (privKey PrivKeySecp256k1) PubKey() PubKey {
	_, pubkeyObject := btcec.PrivKeyFromBytes(privKey[:])
	var pubkeyBytes PubKeySecp256k1
	copy(pubkeyBytes[:], pubkeyObject.SerializeCompressed())
	return pubkeyBytes
}

// Equals - you probably don't need to use this.
// Runs in constant time based on length of the keys.
func (privKey PrivKeySecp256k1) Equals(other PrivKey) bool {
	if otherSecp, ok := other.(PrivKeySecp256k1); ok {
		return subtle.ConstantTimeCompare(privKey[:], otherSecp[:]) == 1
	}
	return false
}

//-------------------------------------

var _ PubKey = PubKeySecp256k1{}

// PubKeySecp256k1Size is comprised of 32 bytes for one field element
// (the x-coordinate), plus one byte for the parity of the y-coordinate.
const PubKeySecp256k1Size = 33

// PubKeySecp256k1 implements crypto.PubKey.
// It is the compressed form of the pubkey. The first byte depends is a 0x02 byte
// if the y-coordinate is the lexicographically largest of the two associated with
// the x-coordinate. Otherwise the first byte is a 0x03.
// This prefix is followed with the x-coordinate.
type PubKeySecp256k1 [PubKeySecp256k1Size]byte

func (pubKey PubKeySecp256k1) AminoSize(_ *amino.Codec) int {
	return 1 + PubKeySecp256k1Size
}

// Address returns a Bitcoin style addresses: RIPEMD160(SHA256(pubkey))
func (pubKey PubKeySecp256k1) Address() Address {
	hasherSHA256 := sha256.New()
	hasherSHA256.Write(pubKey[:]) // does not error
	sha := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha) // does not error
	return Address(hasherRIPEMD160.Sum(nil))
}

// Bytes returns the pubkey marshalled with amino encoding.
func (pubKey PubKeySecp256k1) Bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(pubKey)
	if err != nil {
		panic(err)
	}
	return bz
}

func (pubKey PubKeySecp256k1) String() string {
	return fmt.Sprintf("PubKeySecp256k1{%X}", pubKey[:])
}

func (pubKey PubKeySecp256k1) Equals(other PubKey) bool {
	if otherSecp, ok := other.(PubKeySecp256k1); ok {
		return bytes.Equal(pubKey[:], otherSecp[:])
	}
	return false
}
