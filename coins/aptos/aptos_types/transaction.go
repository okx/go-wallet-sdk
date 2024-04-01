package aptos_types

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
	"golang.org/x/crypto/sha3"
)

type ModuleId struct {
	Address AccountAddress
	Name    Identifier
}

func (o *ModuleId) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := o.Address.Serialize(serializer); err != nil {
		return err
	}
	if err := o.Name.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *ModuleId) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
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

func (o *RawTransaction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := o.Sender.Serialize(serializer); err != nil {
		return err
	}
	if err := serializer.SerializeU64(o.SequenceNumber); err != nil {
		return err
	}
	if err := o.Payload.Serialize(serializer); err != nil {
		return err
	}
	if err := serializer.SerializeU64(o.MaxGasAmount); err != nil {
		return err
	}
	if err := serializer.SerializeU64(o.GasUnitPrice); err != nil {
		return err
	}
	if err := serializer.SerializeU64(o.ExpirationTimestampSecs); err != nil {
		return err
	}
	if err := o.ChainId.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *RawTransaction) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func Sha256Hash(bytes []byte) []byte {
	sha256 := sha3.New256()
	sha256.Write(bytes)
	return sha256.Sum(nil)
}

func (o *RawTransaction) GetSigningMessage() ([]byte, error) {
	prefix := Sha256Hash([]byte("APTOS::RawTransaction"))
	bcsBytes, err := o.BcsSerialize()
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

func (o *Script) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeBytes(o.Code); err != nil {
		return err
	}
	if err := serializeVectorTypeTag(o.TyArgs, serializer); err != nil {
		return err
	}
	if err := serializeVectorTransactionArgument(o.Args, serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *Script) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
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

func (o *ScriptFunction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := o.Module.Serialize(serializer); err != nil {
		return err
	}
	if err := o.Function.Serialize(serializer); err != nil {
		return err
	}
	if err := serializeVectorTypeTag(o.TyArgs, serializer); err != nil {
		return err
	}
	if err := serializeVectorBytes(o.Args, serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *ScriptFunction) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}
func serializeVectorBytes(value [][]byte, serializer serde.Serializer) error {
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

func serializeVectorTypeTag(value []TypeTag, serializer serde.Serializer) error {
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

func serializeVectorTransactionArgument(value []TransactionArgument, serializer serde.Serializer) error {
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

type SignedTransaction struct {
	RawTxn        RawTransaction
	Authenticator TransactionAuthenticator
}

func (o *SignedTransaction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := o.RawTxn.Serialize(serializer); err != nil {
		return err
	}
	if err := o.Authenticator.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *SignedTransaction) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type Transaction interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type TransactionUserTransaction struct {
	Value SignedTransaction
}

func (obj *TransactionUserTransaction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(0)
	if err != nil {
		return err
	}
	if err := obj.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionUserTransaction) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, ErrNullObject
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
	err := serializer.SerializeVariantIndex(0)
	if err != nil {
		return err
	}
	if err := obj.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionPayloadScript) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, ErrNullObject
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

func (o *TransactionPayloadEntryFunction) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(2)
	if err != nil {
		return err
	}
	if err := o.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TransactionPayloadEntryFunction) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgument interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type TransactionArgumentU8 uint8

func (obj *TransactionArgumentU8) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(0)
	if err != nil {
		return err
	}
	if err := serializer.SerializeU8((uint8)(*obj)); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionArgumentU8) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgumentU64 uint64

func (obj *TransactionArgumentU64) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(1)
	if err != nil {
		return err
	}
	if err := serializer.SerializeU64((uint64)(*obj)); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionArgumentU64) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgumentU128 serde.Uint128

func (obj *TransactionArgumentU128) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(2)
	if err != nil {
		return err
	}
	if err := serializer.SerializeU128((serde.Uint128)(*obj)); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionArgumentU128) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgumentAddress struct {
	Value AccountAddress
}

func (o *TransactionArgumentAddress) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(3)
	if err != nil {
		return err
	}
	if err := o.Value.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TransactionArgumentAddress) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionArgumentU8vector []byte

func (o *TransactionArgumentU8vector) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(4)
	if err != nil {
		return err
	}
	if err := serializer.SerializeBytes(*o); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

type TransactionArgumentBool bool

func (o *TransactionArgumentBool) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(5)
	if err != nil {
		return err
	}
	if err := serializer.SerializeBool((bool)(*o)); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *TransactionArgumentBool) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func (o *TransactionArgumentU8vector) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}
