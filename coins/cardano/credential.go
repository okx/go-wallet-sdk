package cardano

import (
	"bytes"
	"fmt"

	"github.com/okx/go-wallet-sdk/coins/cardano/crypto"
)

type StakeCredentialType uint64

const (
	KeyCredential StakeCredentialType = iota
	ScriptCredential
)

type keyStakeCredential struct {
	_       struct{} `cbor:",toarray"`
	Type    StakeCredentialType
	KeyHash AddrKeyHash
}

type scriptStakeCredential struct {
	_          struct{} `cbor:",toarray"`
	Type       StakeCredentialType
	ScriptHash Hash28
}

// StakeCredential is a Cardano credential.
type StakeCredential struct {
	Type       StakeCredentialType
	KeyHash    AddrKeyHash
	ScriptHash Hash28
}

func (s *StakeCredential) Hash() Hash28 {
	if s.Type == KeyCredential {
		return s.KeyHash
	} else {
		return s.ScriptHash
	}
}

// NewKeyCredential creates a Key Credential.
func NewKeyCredential(publicKey crypto.PubKey) (StakeCredential, error) {
	keyHash, err := Blake224Hash(publicKey)
	if err != nil {
		return StakeCredential{}, err
	}
	return StakeCredential{Type: KeyCredential, KeyHash: keyHash}, nil
}

// NewKeyCredential creates a Script Credential.
func NewScriptCredential(script []byte) (StakeCredential, error) {
	scriptHash, err := Blake224Hash(script)
	if err != nil {
		return StakeCredential{}, err
	}
	return StakeCredential{Type: ScriptCredential, ScriptHash: scriptHash}, nil
}

// Equal returns true if the credentials are equal.
func (s *StakeCredential) Equal(rhs StakeCredential) bool {
	if s.Type != rhs.Type {
		return false
	}

	if s.Type == KeyCredential {
		return bytes.Equal(s.KeyHash, rhs.KeyHash)
	} else {
		return bytes.Equal(s.ScriptHash, rhs.ScriptHash)
	}
}

// MarshalCBOR implements cbor.Marshaler.
func (s *StakeCredential) MarshalCBOR() ([]byte, error) {
	var cred []interface{}
	switch s.Type {
	case KeyCredential:
		cred = append(cred, s.Type, s.KeyHash)
	case ScriptCredential:
		cred = append(cred, s.Type, s.ScriptHash)
	}

	return cborEnc.Marshal(cred)

}

// UnmarshalCBOR implements cbor.Unmarshaler.
func (s *StakeCredential) UnmarshalCBOR(data []byte) error {
	credType, err := getTypeFromCBORArray(data)
	if err != nil {
		return fmt.Errorf("cbor: cannot unmarshal CBOR array into StakeCredential (%v)", err)
	}

	switch StakeCredentialType(credType) {
	case KeyCredential:
		cred := &keyStakeCredential{}
		if err := cborDec.Unmarshal(data, cred); err != nil {
			return err
		}
		s.Type = KeyCredential
		s.KeyHash = cred.KeyHash
	case ScriptCredential:
		cred := &scriptStakeCredential{}
		if err := cborDec.Unmarshal(data, cred); err != nil {
			return err
		}
		s.Type = ScriptCredential
		s.ScriptHash = cred.ScriptHash
	}

	return nil
}
