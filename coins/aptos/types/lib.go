package types

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strings"

	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
)

type AccountAddress []byte

func (obj *AccountAddress) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeFixedBytes(*obj); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *AccountAddress) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type AccountAuthenticator interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type AccountAuthenticator__Ed25519 struct {
	PublicKey Ed25519PublicKey
	Signature Ed25519Signature
}

func (obj *AccountAuthenticator__Ed25519) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(0)
	if err := obj.PublicKey.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.Signature.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *AccountAuthenticator__Ed25519) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type ChainId uint8

func (obj *ChainId) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeU8(((uint8)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *ChainId) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type Ed25519PublicKey []byte

func (obj *Ed25519PublicKey) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeBytes((([]byte)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *Ed25519PublicKey) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type Ed25519Signature []byte

func (obj *Ed25519Signature) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeBytes((([]byte)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *Ed25519Signature) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type Identifier string

func (obj *Identifier) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeStr(((string)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *Identifier) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type ModuleId struct {
	Address AccountAddress
	Name    Identifier
}

func (obj *ModuleId) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := obj.Address.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.Name.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *ModuleId) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type RawTransaction struct {
	Sender                  AccountAddress
	SequenceNumber          uint64
	Payload                 TransactionPayload
	MaxGasAmount            uint64
	GasUnitPrice            uint64
	ExpirationTimestampSecs uint64
	ChainId                 ChainId
}

func (obj *RawTransaction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := obj.Sender.Serialize(serializer); err != nil {
		return err
	}
	if err := serializer.SerializeU64(obj.SequenceNumber); err != nil {
		return err
	}
	if err := obj.Payload.Serialize(serializer); err != nil {
		return err
	}
	if err := serializer.SerializeU64(obj.MaxGasAmount); err != nil {
		return err
	}
	if err := serializer.SerializeU64(obj.GasUnitPrice); err != nil {
		return err
	}
	if err := serializer.SerializeU64(obj.ExpirationTimestampSecs); err != nil {
		return err
	}
	if err := obj.ChainId.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *RawTransaction) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func Sha256Hash(bytes []byte) []byte {
	sha256 := sha3.New256()
	sha256.Write(bytes)
	return sha256.Sum(nil)
}

func (obj *RawTransaction) GetSigningMessage() ([]byte, error) {
	prefix := Sha256Hash([]byte("APTOS::RawTransaction"))
	bcsBytes, err := obj.BcsSerialize()
	if err != nil {
		return nil, err
	}
	// prefix + bcsBytes
	message := make([]byte, 0)
	message = append(message, prefix...)
	message = append(message, bcsBytes...)
	return message, nil
}

type Script struct {
	Code   []byte
	TyArgs []TypeTag
	Args   []TransactionArgument
}

func (obj *Script) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeBytes(obj.Code); err != nil {
		return err
	}
	if err := serialize_vector_TypeTag(obj.TyArgs, serializer); err != nil {
		return err
	}
	if err := serialize_vector_TransactionArgument(obj.Args, serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *Script) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type ScriptFunction struct {
	Module   ModuleId
	Function Identifier
	TyArgs   []TypeTag
	Args     [][]byte
}

func (obj *ScriptFunction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := obj.Module.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.Function.Serialize(serializer); err != nil {
		return err
	}
	if err := serialize_vector_TypeTag(obj.TyArgs, serializer); err != nil {
		return err
	}
	if err := serialize_vector_bytes(obj.Args, serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *ScriptFunction) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type SignedTransaction struct {
	RawTxn        RawTransaction
	Authenticator TransactionAuthenticator
}

func (obj *SignedTransaction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := obj.RawTxn.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.Authenticator.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *SignedTransaction) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type StructTag struct {
	Address    AccountAddress
	Module     Identifier
	Name       Identifier
	TypeParams []TypeTag
}

func (obj *StructTag) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := obj.Address.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.Module.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.Name.Serialize(serializer); err != nil {
		return err
	}
	if err := serialize_vector_TypeTag(obj.TypeParams, serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *StructTag) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type Transaction interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type Transaction__UserTransaction struct {
	Value SignedTransaction
}

func (obj *Transaction__UserTransaction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(0)
	if err := obj.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *Transaction__UserTransaction) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionAuthenticator interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type TransactionAuthenticator__Ed25519 struct {
	PublicKey Ed25519PublicKey
	Signature Ed25519Signature
}

func (obj *TransactionAuthenticator__Ed25519) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(0)
	if err := obj.PublicKey.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.Signature.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionAuthenticator__Ed25519) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionPayload interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type TransactionPayloadScript struct {
	Value Script
}

func (obj *TransactionPayloadScript) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(0)
	if err := obj.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionPayloadScript) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionPayloadEntryFunction struct {
	Value ScriptFunction
}

func (obj *TransactionPayloadEntryFunction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(2)
	if err := obj.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionPayloadEntryFunction) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgument interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type TransactionArgument__U8 uint8

func (obj *TransactionArgument__U8) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(0)
	if err := serializer.SerializeU8(((uint8)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionArgument__U8) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgument__U64 uint64

func (obj *TransactionArgument__U64) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(1)
	if err := serializer.SerializeU64(((uint64)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionArgument__U64) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgument__U128 serde.Uint128

func (obj *TransactionArgument__U128) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(2)
	if err := serializer.SerializeU128(((serde.Uint128)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionArgument__U128) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgument__Address struct {
	Value AccountAddress
}

func (obj *TransactionArgument__Address) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(3)
	if err := obj.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionArgument__Address) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgument__U8Vector []byte

func (obj *TransactionArgument__U8Vector) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(4)
	if err := serializer.SerializeBytes((([]byte)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

type TransactionArgument__Bool bool

func (obj *TransactionArgument__Bool) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(5)
	if err := serializer.SerializeBool(((bool)(*obj))); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionArgument__Bool) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (obj *TransactionArgument__U8Vector) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TypeTag interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type TypeTag__Bool struct {
}

func (obj *TypeTag__Bool) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(0)
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TypeTag__Bool) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TypeTag__U8 struct {
}

func (obj *TypeTag__U8) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(1)
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TypeTag__U8) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TypeTag__U64 struct {
}

func (obj *TypeTag__U64) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(2)
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TypeTag__U64) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TypeTag__U128 struct {
}

func (obj *TypeTag__U128) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(3)
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TypeTag__U128) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TypeTag__Address struct {
}

func (obj *TypeTag__Address) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(4)
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TypeTag__Address) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TypeTag__Signer struct {
}

func (obj *TypeTag__Signer) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(5)
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TypeTag__Signer) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TypeTag__Vector struct {
	Value TypeTag
}

func (obj *TypeTag__Vector) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(6)
	if err := obj.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TypeTag__Vector) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TypeTag__Struct struct {
	Value StructTag
}

func (obj *TypeTag__Struct) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(7)
	if err := obj.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TypeTag__Struct) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Cannot serialize null object")
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeUint64(t uint64) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeU64(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeBool(t bool) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeBool(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeU8(t uint8) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeU8(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeU128(t serde.Uint128) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeU128(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeBytes(t []byte) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeBytes(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeFixedBytes(t []byte) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeFixedBytes(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeStr(t string) ([]byte, error) {
	return BcsSerializeBytes([]byte(t))
}

func BcsSerializeLen(t uint64) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeLen(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BytesFromHex(t string) []byte {
	if strings.HasPrefix(t, "0x") {
		t = strings.TrimPrefix(t, "0x")
	}
	bytes, err := hex.DecodeString(t)
	if err != nil {
		panic(err)
	}
	return bytes
}

func serialize_vector_bytes(value [][]byte, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	for _, item := range value {
		if err := serializer.SerializeBytes(item); err != nil {
			return err
		}
	}
	return nil
}

func serialize_vector_TypeTag(value []TypeTag, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	for _, item := range value {
		if err := item.Serialize(serializer); err != nil {
			return err
		}
	}
	return nil
}

func serialize_vector_TransactionArgument(value []TransactionArgument, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	for _, item := range value {
		if err := item.Serialize(serializer); err != nil {
			return err
		}
	}
	return nil
}
