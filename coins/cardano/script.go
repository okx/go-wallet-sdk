package cardano

import (
	"fmt"

	"github.com/okx/go-wallet-sdk/coins/cardano/crypto"
)

type ScriptHashNamespace uint8

const (
	NativeScriptNamespace ScriptHashNamespace = iota
	PlutusScriptNamespace
)

type NativeScriptType uint64

const (
	ScriptPubKey NativeScriptType = iota
	ScriptAll
	ScriptAny
	ScriptNofK
	ScriptInvalidBefore
	ScriptInvalidAfter
)

type scriptPubKey struct {
	_       struct{} `cbor:"_,toarray"`
	Type    NativeScriptType
	KeyHash AddrKeyHash
}

type scriptAll struct {
	_       struct{} `cbor:"_,toarray"`
	Type    NativeScriptType
	Scripts []NativeScript
}

type scriptAny struct {
	_       struct{} `cbor:"_,toarray"`
	Type    NativeScriptType
	Scripts []NativeScript
}

type scriptNofK struct {
	_       struct{} `cbor:"_,toarray"`
	Type    NativeScriptType
	N       uint64
	Scripts []NativeScript
}

type scriptInvalidBefore struct {
	_             struct{} `cbor:"_,toarray"`
	Type          NativeScriptType
	IntervalValue uint64
}

type scriptInvalidAfter struct {
	_             struct{} `cbor:"_,toarray"`
	Type          NativeScriptType
	IntervalValue uint64
}

// NativeScript is a Cardano Native Script.
type NativeScript struct {
	Type          NativeScriptType
	KeyHash       AddrKeyHash
	N             uint64
	Scripts       []NativeScript
	IntervalValue uint64
}

// NewScriptPubKey returns a new Script PubKey.
func NewScriptPubKey(publicKey crypto.PubKey) (NativeScript, error) {
	keyHash, err := publicKey.Hash()
	if err != nil {
		return NativeScript{}, err
	}
	return NativeScript{Type: ScriptPubKey, KeyHash: keyHash}, nil
}

// Hash returns the script hash using blake2b224.
func (ns *NativeScript) Hash() (Hash28, error) {
	bytes, err := ns.Bytes()
	if err != nil {
		return nil, err
	}
	bytes = append([]byte{byte(NativeScriptNamespace)}, bytes...)
	return Blake224Hash(append(bytes))
}

// Bytes returns the CBOR encoding of the script as bytes.
func (ns *NativeScript) Bytes() ([]byte, error) {
	return cborEnc.Marshal(ns)
}

// MarshalCBOR implements cbor.Marshaler.
func (ns *NativeScript) MarshalCBOR() ([]byte, error) {
	var script []interface{}
	switch ns.Type {
	case ScriptPubKey:
		script = append(script, ns.Type, ns.KeyHash)
	case ScriptAll, ScriptAny:
		script = append(script, ns.Type, ns.Scripts)
	case ScriptNofK:
		script = append(script, ns.Type, ns.N, ns.Scripts)
	case ScriptInvalidBefore, ScriptInvalidAfter:
		script = append(script, ns.Type, ns.IntervalValue)
	}
	return cborEnc.Marshal(script)
}

// UnmarshalCBOR implements cbor.Unmarshaler.
func (ns *NativeScript) UnmarshalCBOR(data []byte) error {
	nsType, err := getTypeFromCBORArray(data)
	if err != nil {
		return fmt.Errorf("cbor: cannot unmarshal CBOR array into StakeCredential (%v)", err)
	}

	switch NativeScriptType(nsType) {
	case ScriptPubKey:
		script := scriptPubKey{}
		if err := cborDec.Unmarshal(data, &script); err != nil {
			return err
		}
		ns.Type = script.Type
		ns.KeyHash = script.KeyHash
	case ScriptAll:
		script := scriptAll{}
		if err := cborDec.Unmarshal(data, &script); err != nil {
			return err
		}
		ns.Type = script.Type
		ns.Scripts = script.Scripts
	case ScriptAny:
		script := scriptAny{}
		if err := cborDec.Unmarshal(data, &script); err != nil {
			return err
		}
		ns.Type = script.Type
		ns.Scripts = script.Scripts
	case ScriptNofK:
		script := scriptNofK{}
		if err := cborDec.Unmarshal(data, &script); err != nil {
			return err
		}
		ns.Type = script.Type
		ns.N = script.N
		ns.Scripts = script.Scripts
	case ScriptInvalidBefore:
		script := scriptInvalidBefore{}
		if err := cborDec.Unmarshal(data, &script); err != nil {
			return err
		}
		ns.Type = script.Type
		ns.IntervalValue = script.IntervalValue
	case ScriptInvalidAfter:
		script := scriptInvalidAfter{}
		if err := cborDec.Unmarshal(data, &script); err != nil {
			return err
		}
		ns.Type = script.Type
		ns.IntervalValue = script.IntervalValue
	}

	return nil
}
