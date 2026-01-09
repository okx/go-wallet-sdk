package cardano

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	_ "crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	"github.com/veraison/go-cose"
	ed255192 "github.com/okx/go-wallet-sdk/coins/cardano/ed25519"
)

type OkExtendedPrivateKey []byte

func (okxpriv OkExtendedPrivateKey) Public() crypto.PublicKey {
	xpriv := ed255192.ExtendedPrivateKey(okxpriv)
	pubBytes := []byte(xpriv.Public().(ed255192.PublicKey))
	return ed25519.PublicKey(pubBytes)
}
func (okxpriv OkExtendedPrivateKey) Sign(rand io.Reader, message []byte, _ crypto.SignerOpts) (signature []byte, err error) {
	xpriv := ed255192.ExtendedPrivateKey(okxpriv)
	opts := crypto.SignerOpts(crypto.Hash(0))
	return xpriv.Sign(rand, message, opts)
}

type SignMessageResult struct {
	Signature string `json:"signature"`
	Key       string `json:"key"`
}

func SignMessage(prvKey string, address string, message string) (*SignMessageResult, error) {
	pBytes, err := hex.DecodeString(prvKey)
	if err != nil {
		return nil, err
	}

	// stake private coseKey
	pBytes = pBytes[:64] // payment coseKey

	privateKey := OkExtendedPrivateKey(pBytes)

	signer, err := cose.NewSigner(cose.AlgorithmEdDSA, privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return SignMessageWithSigner(address, message, signer, publicKey)
}

type SignFunc func(data []byte) ([]byte, error)
type SecureSigner struct {
	sign SignFunc
}

func (s SecureSigner) Algorithm() cose.Algorithm {
	return cose.AlgorithmEdDSA
}

func (s SecureSigner) Sign(rand io.Reader, content []byte) ([]byte, error) {
	return s.sign(content)
}

func NewSecureSigner(sign SignFunc) (cose.Signer, error) {
	return &SecureSigner{sign: sign}, nil

}

func SignMessageWithSigner(address string, message string, signer cose.Signer, publicKey []byte) (*SignMessageResult, error) {
	addr, err := NewAddress(address)
	if err != nil {
		return nil, err
	}

	data, err := hex.DecodeString(message)
	if err != nil {
		return nil, err
	}

	headers := cose.Headers{
		Protected: cose.ProtectedHeader{
			cose.HeaderLabelAlgorithm: cose.AlgorithmEdDSA,
			"address":                 addr.Bytes(),
		},
		Unprotected: cose.UnprotectedHeader{
			"hashed": false,
		},
	}

	cbor, err := cose.Sign1(rand.Reader, signer, headers, data, nil)
	if err != nil {
		return nil, err
	}

	// signed message is tagged, convert to untagged
	var m cose.Sign1Message
	err = m.UnmarshalCBOR(cbor)
	if err != nil {
		return nil, err
	}

	u := cose.UntaggedSign1Message(m)

	signature, err := u.MarshalCBOR()
	if err != nil {
		return nil, err
	}

	coseKey, err := cose.NewKeyOKP(cose.AlgorithmEdDSA, publicKey, nil)
	if err != nil {
		return nil, err
	}
	key, err := coseKey.MarshalCBOR()
	if err != nil {
		return nil, err
	}
	return &SignMessageResult{Signature: hex.EncodeToString(signature), Key: hex.EncodeToString(key)}, nil

}

func VerifyMessage(signature string, key string, publicKey string, address string, message string) (bool, error) {
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	var msg cose.UntaggedSign1Message
	err = msg.UnmarshalCBOR(sig)
	if err != nil {
		return false, err
	}

	ckey, err := hex.DecodeString(key)
	if err != nil {
		return false, err
	}
	var coseKey cose.Key
	err = coseKey.UnmarshalCBOR(ckey)
	if err != nil {
		return false, err
	}

	pbytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	// stake pub key
	pbytes = pbytes[:32]
	pubKey := ed25519.PublicKey(pbytes)

	addr, err := NewAddress(address)
	if err != nil {
		return false, err
	}

	data, err := hex.DecodeString(message)
	if err != nil {
		return false, err
	}

	// Verify key data
	if coseKey.Algorithm != cose.AlgorithmEdDSA {
		return false, errors.New("key algorithm not EdDSA")
	}
	if coseKey.Params[cose.KeyLabelEC2Curve] != cose.CurveEd25519 {
		return false, errors.New("key curve not Ed25519")
	}
	if !bytes.Equal(coseKey.Params[cose.KeyLabelOKPX].([]byte), pubKey) {
		return false, errors.New("key public key mismatch")
	}

	// Verify signature message data
	if !bytes.Equal(msg.Headers.Protected["address"].([]byte), addr.Bytes()) {
		return false, errors.New("signature address mismatch")
	}
	if msg.Headers.Unprotected["hashed"] != false {
		return false, errors.New("signature hashed not false")
	}
	if !bytes.Equal(msg.Payload, data) {
		return false, err
	}

	// Verify signature
	verifier, err := cose.NewVerifier(cose.AlgorithmEdDSA, pubKey)
	if err != nil {
		return false, err
	}
	err = msg.Verify([]byte{}, verifier)
	return err == nil, err

}

func VerifyMessageNoAddr(signature string, key string, publicKey string, message string) (bool, error) {
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	var msg cose.UntaggedSign1Message
	err = msg.UnmarshalCBOR(sig)
	if err != nil {
		return false, err
	}

	ckey, err := hex.DecodeString(key)
	if err != nil {
		return false, err
	}
	var coseKey cose.Key
	err = coseKey.UnmarshalCBOR(ckey)
	if err != nil {
		return false, err
	}

	pbytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	// stake pub key
	pbytes = pbytes[:32]
	pubKey := ed25519.PublicKey(pbytes)

	data, err := hex.DecodeString(message)
	if err != nil {
		return false, err
	}

	// Verify key data
	if coseKey.Algorithm != cose.AlgorithmEdDSA {
		return false, errors.New("key algorithm not EdDSA")
	}
	if coseKey.Params[cose.KeyLabelEC2Curve] != cose.CurveEd25519 {
		return false, errors.New("key curve not Ed25519")
	}
	if !bytes.Equal(coseKey.Params[cose.KeyLabelOKPX].([]byte), pubKey) {
		return false, errors.New("key public key mismatch")
	}

	// Verify signature message data
	if msg.Headers.Unprotected["hashed"] != false {
		return false, errors.New("signature hashed not false")
	}
	if !bytes.Equal(msg.Payload, data) {
		return false, err
	}

	// Verify signature
	verifier, err := cose.NewVerifier(cose.AlgorithmEdDSA, pubKey)
	if err != nil {
		return false, err
	}
	err = msg.Verify([]byte{}, verifier)
	return err == nil, err

}
