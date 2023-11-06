package eos

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/okx/go-wallet-sdk/coins/eos/types"
)

type Signer struct {
	Keys []*ecc.PrivateKey `json:"keys"`
}

func NewSigner(keys []*ecc.PrivateKey) *Signer {
	return &Signer{keys}
}

func NewSignerFromWIFs(wifs []string) (*Signer, error) {
	keys := make([]*ecc.PrivateKey, len(wifs))
	for i, wif := range wifs {
		key, err := ecc.NewPrivateKey(wif)
		if err != nil {
			return nil, err
		}
		keys[i] = key
	}
	return &Signer{keys}, nil
}

func (b *Signer) Add(wifKey string) error {
	privateKey, err := ecc.NewPrivateKey(wifKey)
	if err != nil {
		return err
	}
	if privateKey == nil {
		// 目前理论走不到这里
		return errors.New("appending a nil private key is forbidden")
	}
	b.Keys = append(b.Keys, privateKey)
	return nil
}

func (b *Signer) Sign(tx *types.SignedTransaction, chainID []byte, requiredKeys ...ecc.PublicKey) (*types.SignedTransaction, error) {
	// TODO: probably want to use `tx.packed` and hash the ContextFreeData also.
	if tx == nil {
		return nil, errors.New("cannot sign a nil transaction")
	}
	txdata, cfd, err := tx.PackedTransactionAndCFD()
	if err != nil {
		return nil, err
	}

	sigDigest := SigDigest(chainID, txdata, cfd)

	keyMap := b.keyMap()
	for _, key := range requiredKeys {
		privKey := keyMap[key.String()]
		if privKey == nil {
			return nil, fmt.Errorf("private key for %q not in signer", key)
		}
		sig, err := privKey.Sign(sigDigest)
		if err != nil {
			return nil, err
		}

		tx.Signatures = append(tx.Signatures, sig)
	}

	return tx, nil
}

func (b *Signer) keyMap() map[string]*ecc.PrivateKey {
	out := map[string]*ecc.PrivateKey{}
	for _, key := range b.Keys {
		out[key.PublicKey().String()] = key
	}
	return out
}

// SigDigest computes the hash of the packed transaction
func SigDigest(chainID, payload, contextFreeData []byte) []byte {
	h := sha256.New()
	if len(chainID) == 0 {
		_, _ = h.Write(make([]byte, 32, 32))
	} else {
		_, _ = h.Write(chainID)
	}
	_, _ = h.Write(payload)

	if len(contextFreeData) > 0 {
		h2 := sha256.New()
		_, _ = h2.Write(contextFreeData)
		_, _ = h.Write(h2.Sum(nil)) // add the hash of CFD to the payload
	} else {
		_, _ = h.Write(make([]byte, 32, 32))
	}
	return h.Sum(nil)
}
