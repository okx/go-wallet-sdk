package crypto

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/okx/go-wallet-sdk/coins/cardano/ed25519"
	"github.com/okx/go-wallet-sdk/crypto/bech32"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/pbkdf2"
)

// XPrvKey is the extended private key (64 bytes) appended with the chain code (32 bytes).
type XPrvKey []byte

// NewXPrvKey creates a new extended private key from a bech32 encoded private key.
func NewXPrvKey(bech string) (XPrvKey, error) {
	_, xsk, err := bech32.DecodeToBase256(bech)
	return xsk, err
}

func NewXPrvKeyFromEntropy(entropy []byte, password string) XPrvKey {
	key := pbkdf2.Key([]byte(password), entropy, 4096, 96, sha512.New)

	key[0] &= 0xf8
	key[31] = (key[31] & 0x1f) | 0x40

	return key
}

// Bech32 returns the private key encoded as bech32.
func (prv XPrvKey) Bech32(prefix string) string {
	bech, err := bech32.EncodeFromBase256(prefix, prv)
	if err != nil {
		panic(err)
	}
	return bech
}

func (prv XPrvKey) String() string {
	return hex.EncodeToString(prv)
}

// PrvKey returns the ed25519 extended private key.
func (prv XPrvKey) PrvKey() PrvKey {
	return PrvKey(prv[:64])
}

// XPubKey returns the XPubKey derived from the extended private key.
func (prv XPrvKey) XPubKey() XPubKey {
	xvk := make([]byte, 64)
	vk := prv.PrvKey().PubKey()
	cc := prv[64:]

	copy(xvk[:32], vk)
	copy(xvk[32:], cc)

	return xvk
}

// XPubKey returns the XPubKey derived from the extended private key.
func (prv XPrvKey) PubKey() PubKey {
	return PubKey(ed25519.PublicKeyFrom(ed25519.ExtendedPrivateKey(prv[:64])))
}

func (prv *XPrvKey) Sign(message []byte) []byte {
	pk := prv.PrvKey()
	return pk.Sign(message)
}

// XPubKey is the public key (32 bytes) appended with the chain code (32 bytes).
type XPubKey []byte

// NewXPubKey creates a new extended public key from a bech32 encoded extended public key.
func NewXPubKey(bech string) (XPubKey, error) {
	_, xsk, err := bech32.DecodeToBase256(bech)
	return xsk, err
}

// XPubKey returns the PubKey from the extended public key.
func (pub XPubKey) PubKey() PubKey {
	return PubKey(pub[:32])
}

// NewPubKey creates a new public key from a bech32 encoded public key.
func NewPubKey(bech string) (PubKey, error) {
	_, xsk, err := bech32.DecodeToBase256(bech)
	return xsk, err
}

// Verify reports whether sig is a valid signature of message by the extended public key.
func (pub XPubKey) Verify(message, sig []byte) bool {
	return pub.PubKey().Verify(message, sig)
}

func (pub XPubKey) String() string {
	return hex.EncodeToString(pub)
}

// PubKey is a edd25519 public key.
type PubKey []byte

// Verify reports whether sig is a valid signature of message by the public key.
func (pub PubKey) Verify(message, signature []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(pub), message, signature)
}

// Bech32 returns the public key encoded as bech32.
func (pub PubKey) Bech32(prefix string) string {
	bech, err := bech32.EncodeFromBase256(prefix, pub)
	if err != nil {
		panic(err)
	}
	return bech
}

func (pub PubKey) String() string {
	return hex.EncodeToString(pub)
}

// Hash returns the public key hash using blake2b224.
func (pub PubKey) Hash() ([]byte, error) {
	return blake224Hash(pub)
}

// PrvKey is a ed25519 extended private key.
type PrvKey []byte

// NewPrvKey creates a new private key from a bech32 encoded private key.
func NewPrvKey(bech string) (PrvKey, error) {
	_, xsk, err := bech32.DecodeToBase256(bech)
	return xsk, err
}

// XPubKey returns the XPubKey derived from the private key.
func (prv PrvKey) PubKey() PubKey {
	vk := make([]byte, 32)
	pk := ed25519.PublicKeyFrom(ed25519.ExtendedPrivateKey(prv[:64]))

	copy(vk[:32], pk)

	return vk
}

// Bech32 returns the private key encoded as bech32.
func (prv PrvKey) Bech32(prefix string) string {
	bech, err := bech32.EncodeFromBase256(prefix, prv)
	if err != nil {
		panic(err)
	}
	return bech
}

func (prv PrvKey) String() string {
	return hex.EncodeToString(prv)
}

func (prv *PrvKey) Sign(message []byte) []byte {
	pk := ed25519.ExtendedPrivateKey((*prv)[:64])
	return ed25519.SignExtended(pk, message)
}

func blake224Hash(b []byte) ([]byte, error) {
	hash, err := blake2b.New(224/8, nil)
	if err != nil {
		return nil, err
	}
	_, err = hash.Write(b)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), err
}
