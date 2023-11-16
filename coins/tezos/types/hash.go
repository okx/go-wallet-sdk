/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrUnknownHashType describes an error where a hash can not
	// decoded as a specific hash type because the string encoding
	// starts with an unknown identifier.
	ErrUnknownHashType = errors.New("tezos: unknown hash type")
)

type HashType byte

const (
	HashTypeInvalid HashType = iota
	HashTypeChainId
	HashTypeId
	HashTypePkhEd25519
	HashTypePkhSecp256k1
	HashTypePkhP256
	HashTypePkhNocurve
	HashTypePkhBlinded
	HashTypeBlock
	HashTypeOperation
	HashTypeOperationList
	HashTypeOperationListList
	HashTypeProtocol
	HashTypeContext
	HashTypeNonce
	HashTypeSeedEd25519
	HashTypePkEd25519
	HashTypeSkEd25519
	HashTypePkSecp256k1
	HashTypeSkSecp256k1
	HashTypePkP256
	HashTypeSkP256
	HashTypeScalarSecp256k1
	HashTypeElementSecp256k1
	HashTypeScriptExpr
	HashTypeEncryptedSeedEd25519
	HashTypeEncryptedSkSecp256k1
	HashTypeEncryptedSkP256
	HashTypeSigEd25519
	HashTypeSigSecp256k1
	HashTypeSigP256
	HashTypeSigGeneric

	HashTypeBlockPayload
	HashTypeBlockMetadata
	HashTypeOperationMetadata
	HashTypeOperationMetadataList
	HashTypeOperationMetadataListList
	HashTypeEncryptedSecp256k1Scalar
	HashTypeSaplingSpendingKey
	HashTypeSaplingAddress

	HashTypePkhBls12_381
	HashTypeSigGenericAggregate
	HashTypeSigBls12_381
	HashTypePkBls12_381
	HashTypeSkBls12_381
	HashTypeEncryptedSkBls12_381
	HashTypeToruAddress
	HashTypeToruInbox
	HashTypeToruMessage
	HashTypeToruCommitment
	HashTypeToruMessageResult
	HashTypeToruMessageResultList
	HashTypeToruWithdrawList
)

func (t HashType) PrefixBytes() []byte {
	switch t {
	case HashTypeChainId:
		return CHAIN_ID
	case HashTypeId:
		return ID_HASH_ID
	case HashTypePkhEd25519:
		return ED25519_PUBLIC_KEY_HASH_ID
	case HashTypePkhSecp256k1:
		return SECP256K1_PUBLIC_KEY_HASH_ID
	case HashTypePkhP256:
		return P256_PUBLIC_KEY_HASH_ID
	case HashTypePkhNocurve:
		return NOCURVE_PUBLIC_KEY_HASH_ID
	case HashTypePkhBlinded:
		return BLINDED_PUBLIC_KEY_HASH_ID
	case HashTypeBlock:
		return BLOCK_HASH_ID
	case HashTypeOperation:
		return OPERATION_HASH_ID
	case HashTypeOperationList:
		return OPERATION_LIST_HASH_ID
	case HashTypeOperationListList:
		return OPERATION_LIST_LIST_HASH_ID
	case HashTypeProtocol:
		return PROTOCOL_HASH_ID
	case HashTypeContext:
		return CONTEXT_HASH_ID
	case HashTypeNonce:
		return NONCE_HASH_ID
	case HashTypeSeedEd25519:
		return ED25519_SEED_ID
	case HashTypePkEd25519:
		return ED25519_PUBLIC_KEY_ID
	case HashTypeSkEd25519:
		return ED25519_SECRET_KEY_ID
	case HashTypePkSecp256k1:
		return SECP256K1_PUBLIC_KEY_ID
	case HashTypeSkSecp256k1:
		return SECP256K1_SECRET_KEY_ID
	case HashTypePkP256:
		return P256_PUBLIC_KEY_ID
	case HashTypeSkP256:
		return P256_SECRET_KEY_ID
	case HashTypeScalarSecp256k1:
		return SECP256K1_SCALAR_ID
	case HashTypeElementSecp256k1:
		return SECP256K1_ELEMENT_ID
	case HashTypeScriptExpr:
		return SCRIPT_EXPR_HASH_ID
	case HashTypeEncryptedSeedEd25519:
		return ED25519_ENCRYPTED_SEED_ID
	case HashTypeEncryptedSkSecp256k1:
		return SECP256K1_ENCRYPTED_SECRET_KEY_ID
	case HashTypeEncryptedSkP256:
		return P256_ENCRYPTED_SECRET_KEY_ID
	case HashTypeSigEd25519:
		return ED25519_SIGNATURE_ID
	case HashTypeSigSecp256k1:
		return SECP256K1_SIGNATURE_ID
	case HashTypeSigP256:
		return P256_SIGNATURE_ID
	case HashTypeSigGeneric:
		return GENERIC_SIGNATURE_ID
	case HashTypeBlockPayload:
		return BLOCK_PAYLOAD_HASH_ID
	case HashTypeBlockMetadata:
		return BLOCK_METADATA_HASH_ID
	case HashTypeOperationMetadata:
		return OPERATION_METADATA_HASH_ID
	case HashTypeOperationMetadataList:
		return OPERATION_METADATA_LIST_HASH_ID
	case HashTypeOperationMetadataListList:
		return OPERATION_METADATA_LIST_LIST_HASH_ID
	case HashTypeEncryptedSecp256k1Scalar:
		return SECP256K1_ENCRYPTED_SCALAR_ID
	case HashTypeSaplingSpendingKey:
		return SAPLING_SPENDING_KEY_ID
	case HashTypeSaplingAddress:
		return SAPLING_ADDRESS_ID
	case HashTypePkhBls12_381:
		return BLS12_381_PUBLIC_KEY_HASH_ID
	case HashTypeSigGenericAggregate:
		return GENERIC_AGGREGATE_SIGNATURE_ID
	case HashTypeSigBls12_381:
		return BLS12_381_SIGNATURE_ID
	case HashTypePkBls12_381:
		return BLS12_381_PUBLIC_KEY_ID
	case HashTypeSkBls12_381:
		return BLS12_381_SECRET_KEY_ID
	case HashTypeEncryptedSkBls12_381:
		return BLS12_381_ENCRYPTED_SECRET_KEY_ID
	case HashTypeToruAddress:
		return TORU_ADDRESS_ID
	case HashTypeToruInbox:
		return TORU_INBOX_HASH_ID
	case HashTypeToruMessage:
		return TORU_MESSAGE_HASH_ID
	case HashTypeToruCommitment:
		return TORU_COMMITMENT_HASH_ID
	case HashTypeToruMessageResult:
		return TORU_MESSAGE_RESULT_HASH_ID
	case HashTypeToruMessageResultList:
		return TORU_MESSAGE_RESULT_LIST_HASH_ID
	case HashTypeToruWithdrawList:
		return TORU_WITHDRAW_LIST_HASH_ID
	default:
		return nil
	}
}

type ChainIdHash struct {
	Hash
}

func (h ChainIdHash) Clone() ChainIdHash {
	return ChainIdHash{h.Hash.Clone()}
}

type Hash struct {
	Type HashType
	Hash []byte
}

// Bytes returns the raw byte representation of the hash without type info.
func (h Hash) Bytes() []byte {
	return h.Hash
}

func (h Hash) IsValid() bool {
	return h.Type != HashTypeInvalid && len(h.Hash) == h.Type.Len()
}

func (h Hash) Equal(h2 Hash) bool {
	return h.Type == h2.Type && bytes.Equal(h.Hash, h2.Hash)
}

func (h *Hash) UnmarshalText(data []byte) error {
	x, err := decodeHash(string(data))
	if err != nil {
		return err
	}
	*h = x
	return nil
}

func (h Hash) Clone() Hash {
	buf := make([]byte, len(h.Hash))
	copy(buf, h.Hash)
	return Hash{
		Type: h.Type,
		Hash: buf,
	}
}

func decodeHash(hstr string) (Hash, error) {
	typ := ParseHashType(hstr)
	if typ == HashTypeInvalid {
		return Hash{}, ErrUnknownHashType
	}
	decoded, version, err := CheckDecode(hstr, len(typ.PrefixBytes()), nil)
	if err != nil {
		if err == ErrChecksum {
			return Hash{}, ErrChecksumMismatch
		}
		return Hash{}, fmt.Errorf("tezos: unknown hash format: %w", err)
	}
	if !bytes.Equal(version, typ.PrefixBytes()) {
		return Hash{}, fmt.Errorf("tezos: invalid prefix '%x' for decoded hash type '%s'", version, typ)
	}
	if have, want := len(decoded), typ.Len(); have != want {
		return Hash{}, fmt.Errorf("tezos: invalid length for decoded hash have=%d want=%d", have, want)
	}
	return Hash{
		Type: typ,
		Hash: decoded,
	}, nil
}

func (t HashType) Len() int {
	switch t {
	case HashTypeChainId:
		return 4
	case HashTypeId:
		return 16
	case HashTypePkhEd25519,
		HashTypePkhSecp256k1,
		HashTypePkhP256,
		HashTypePkhNocurve,
		HashTypePkhBlinded,
		HashTypePkhBls12_381,
		HashTypeToruAddress:
		return 20
	case HashTypeBlock,
		HashTypeOperation,
		HashTypeOperationList,
		HashTypeOperationListList,
		HashTypeProtocol,
		HashTypeContext,
		HashTypeNonce,
		HashTypeSeedEd25519,
		HashTypePkEd25519,
		HashTypeSkSecp256k1,
		HashTypeSkP256,
		HashTypeScriptExpr,
		HashTypeBlockPayload,
		HashTypeBlockMetadata,
		HashTypeOperationMetadata,
		HashTypeOperationMetadataList,
		HashTypeOperationMetadataListList,
		HashTypeToruInbox,
		HashTypeToruMessage,
		HashTypeToruCommitment,
		HashTypeToruMessageResultList,
		HashTypeToruWithdrawList,
		HashTypeToruMessageResult,
		HashTypeSkBls12_381:
		return 32
	case HashTypePkSecp256k1,
		HashTypePkP256,
		HashTypeScalarSecp256k1,
		HashTypeElementSecp256k1:
		return 33
	case HashTypeSaplingAddress:
		return 43
	case HashTypePkBls12_381:
		return 48
	case HashTypeEncryptedSeedEd25519,
		HashTypeEncryptedSkSecp256k1,
		HashTypeEncryptedSkP256:
		return 56
	case HashTypeEncryptedSkBls12_381:
		return 58
	case HashTypeEncryptedSecp256k1Scalar:
		return 60
	case HashTypeSkEd25519,
		HashTypeSigEd25519,
		HashTypeSigSecp256k1,
		HashTypeSigP256,
		HashTypeSigGeneric:
		return 64
	case HashTypeSigGenericAggregate,
		HashTypeSigBls12_381:
		return 96
	case HashTypeSaplingSpendingKey:
		return 169
	default:
		return 0
	}
}

func ParseHashType(s string) HashType {
	switch len(s) {
	case 15:
		if strings.HasPrefix(s, CHAIN_ID_PREFIX) {
			return HashTypeChainId
		}
	case 30:
		if strings.HasPrefix(s, ID_HASH_PREFIX) {
			return HashTypeId
		}
	case 36:
		switch true {
		case strings.HasPrefix(s, ED25519_PUBLIC_KEY_HASH_PREFIX):
			return HashTypePkhEd25519
		case strings.HasPrefix(s, SECP256K1_PUBLIC_KEY_HASH_PREFIX):
			return HashTypePkhSecp256k1
		case strings.HasPrefix(s, P256_PUBLIC_KEY_HASH_PREFIX):
			return HashTypePkhP256
		case strings.HasPrefix(s, NOCURVE_PUBLIC_KEY_HASH_PREFIX):
			return HashTypePkhNocurve
		case strings.HasPrefix(s, BLINDED_PUBLIC_KEY_HASH_PREFIX):
			return HashTypePkhBlinded
		case strings.HasPrefix(s, BLS12_381_PUBLIC_KEY_HASH_PREFIX):
			return HashTypePkhBls12_381
		}
	case 37:
		switch true {
		case strings.HasPrefix(s, TORU_ADDRESS_PREFIX):
			return HashTypeToruAddress
		}
	case 43:
		switch true {
		case strings.HasPrefix(s, SAPLING_ADDRESS_PREFIX):
			return HashTypeSaplingAddress
		}
	case 51:
		switch true {
		case strings.HasPrefix(s, BLOCK_HASH_PREFIX):
			return HashTypeBlock
		case strings.HasPrefix(s, OPERATION_HASH_PREFIX):
			return HashTypeOperation
		case strings.HasPrefix(s, PROTOCOL_HASH_PREFIX):
			return HashTypeProtocol
		case strings.HasPrefix(s, OPERATION_METADATA_HASH_PREFIX):
			return HashTypeOperationMetadata
		}
	case 52:
		switch true {
		case strings.HasPrefix(s, OPERATION_LIST_HASH_PREFIX):
			return HashTypeOperationList
		case strings.HasPrefix(s, CONTEXT_HASH_PREFIX):
			return HashTypeContext
		case strings.HasPrefix(s, BLOCK_PAYLOAD_HASH_PREFIX):
			return HashTypeBlockPayload
		case strings.HasPrefix(s, BLOCK_METADATA_HASH_PREFIX):
			return HashTypeBlockMetadata
		case strings.HasPrefix(s, OPERATION_METADATA_LIST_HASH_PREFIX):
			return HashTypeOperationMetadataList
		}
	case 53:
		switch true {
		case strings.HasPrefix(s, OPERATION_LIST_LIST_HASH_PREFIX):
			return HashTypeOperationListList
		case strings.HasPrefix(s, SECP256K1_SCALAR_PREFIX):
			return HashTypeScalarSecp256k1
		case strings.HasPrefix(s, NONCE_HASH_PREFIX):
			return HashTypeNonce
		case strings.HasPrefix(s, OPERATION_METADATA_LIST_LIST_HASH_PREFIX):
			return HashTypeOperationMetadataListList
		case strings.HasPrefix(s, TORU_INBOX_HASH_PREFIX):
			return HashTypeToruInbox
		case strings.HasPrefix(s, TORU_MESSAGE_HASH_PREFIX):
			return HashTypeToruMessage
		case strings.HasPrefix(s, TORU_COMMITMENT_HASH_PREFIX):
			return HashTypeToruCommitment
		case strings.HasPrefix(s, TORU_MESSAGE_RESULT_LIST_HASH_PREFIX):
			return HashTypeToruMessageResultList
		case strings.HasPrefix(s, TORU_WITHDRAW_LIST_HASH_PREFIX):
			return HashTypeToruWithdrawList
		}
	case 54:
		switch true {
		case strings.HasPrefix(s, ED25519_SEED_PREFIX):
			return HashTypeSeedEd25519
		case strings.HasPrefix(s, ED25519_PUBLIC_KEY_PREFIX):
			return HashTypePkEd25519
		case strings.HasPrefix(s, SECP256K1_SECRET_KEY_PREFIX):
			return HashTypeSkSecp256k1
		case strings.HasPrefix(s, P256_SECRET_KEY_PREFIX):
			return HashTypeSkP256
		case strings.HasPrefix(s, SECP256K1_ELEMENT_PREFIX):
			return HashTypeElementSecp256k1
		case strings.HasPrefix(s, SCRIPT_EXPR_HASH_PREFIX):
			return HashTypeScriptExpr
		case strings.HasPrefix(s, BLS12_381_SECRET_KEY_PREFIX):
			return HashTypeSkBls12_381
		case strings.HasPrefix(s, TORU_MESSAGE_RESULT_HASH_PREFIX):
			return HashTypeToruMessageResult
		}
	case 55:
		switch true {
		case strings.HasPrefix(s, SECP256K1_PUBLIC_KEY_PREFIX):
			return HashTypePkSecp256k1
		case strings.HasPrefix(s, P256_PUBLIC_KEY_PREFIX):
			return HashTypePkP256
		}
	case 76:
		switch true {
		case strings.HasPrefix(s, BLS12_381_PUBLIC_KEY_PREFIX):
			return HashTypePkBls12_381
		}
	case 88:
		switch true {
		case strings.HasPrefix(s, ED25519_ENCRYPTED_SEED_PREFIX):
			return HashTypeEncryptedSeedEd25519
		case strings.HasPrefix(s, SECP256K1_ENCRYPTED_SECRET_KEY_PREFIX):
			return HashTypeEncryptedSkSecp256k1
		case strings.HasPrefix(s, P256_ENCRYPTED_SECRET_KEY_PREFIX):
			return HashTypeEncryptedSkP256
		case strings.HasPrefix(s, BLS12_381_ENCRYPTED_SECRET_KEY_PREFIX):
			return HashTypeEncryptedSkBls12_381
		}
	case 93:
		switch true {
		case strings.HasPrefix(s, SECP256K1_ENCRYPTED_SCALAR_PREFIX):
			return HashTypeEncryptedSecp256k1Scalar
		}
	case 96:
		if strings.HasPrefix(s, GENERIC_SIGNATURE_PREFIX) {
			return HashTypeSigGeneric
		}
	case 98:
		switch true {
		case strings.HasPrefix(s, ED25519_SECRET_KEY_PREFIX):
			return HashTypeSkEd25519
		case strings.HasPrefix(s, P256_SIGNATURE_PREFIX):
			return HashTypeSigP256
		}
	case 99:
		switch true {
		case strings.HasPrefix(s, ED25519_SIGNATURE_PREFIX):
			return HashTypeSigEd25519
		case strings.HasPrefix(s, SECP256K1_SIGNATURE_PREFIX):
			return HashTypeSigSecp256k1
		}
	case 141:
		switch true {
		case strings.HasPrefix(s, GENERIC_AGGREGATE_SIGNATURE_PREFIX):
			return HashTypeSigGenericAggregate
		}
	case 142:
		switch true {
		case strings.HasPrefix(s, BLS12_381_SIGNATURE_PREFIX):
			return HashTypeSigBls12_381
		}
	case 169:
		switch true {
		case strings.HasPrefix(s, SAPLING_SPENDING_KEY_PREFIX):
			return HashTypeSaplingSpendingKey
		}
	}
	return HashTypeInvalid
}

type ProtocolHash struct {
	Hash
}

func (h ProtocolHash) Equal(h2 ProtocolHash) bool {
	return h.Hash.Equal(h2.Hash)
}

func (h *ProtocolHash) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if !strings.HasPrefix(string(data), PROTOCOL_HASH_PREFIX) {
		return fmt.Errorf("tezos: invalid prefix for protocol hash '%s'", string(data))
	}
	if err := h.Hash.UnmarshalText(data); err != nil {
		return err
	}
	if h.Type != HashTypeProtocol {
		return fmt.Errorf("tezos: invalid type.")
	}
	if len(h.Hash.Hash) != h.Type.Len() {
		return fmt.Errorf("tezos: invalid len %d for protocol hash", len(h.Hash.Hash))
	}
	return nil
}

func ParseProtocolHashSafe(s string) ProtocolHash {
	var h ProtocolHash
	h.UnmarshalText([]byte(s))
	return h
}

func (h *ChainIdHash) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if !strings.HasPrefix(string(data), CHAIN_ID_PREFIX) {
		return fmt.Errorf("tezos: invalid prefix for chain id hash '%s'", string(data))
	}
	if err := h.Hash.UnmarshalText(data); err != nil {
		return err
	}
	if h.Type != HashTypeChainId {
		return fmt.Errorf("tezos: invalid type for chain id hash")
	}
	if len(h.Hash.Hash) != h.Type.Len() {
		return fmt.Errorf("tezos: invalid len %d for chain id hash", len(h.Hash.Hash))
	}
	return nil
}

func (h ChainIdHash) Equal(h2 ChainIdHash) bool {
	return h.Hash.Equal(h2.Hash)
}

func MustParseChainIdHash(s string) ChainIdHash {
	h, err := ParseChainIdHash(s)
	if err != nil {
		// TODO:
		//panic(err)
		return h
	}
	return h
}

func ParseChainIdHash(s string) (ChainIdHash, error) {
	var h ChainIdHash
	if err := h.UnmarshalText([]byte(s)); err != nil {
		return h, err
	}
	return h, nil
}

type BlockHash struct {
	Hash
}
