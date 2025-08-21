package v2

import (
	"errors"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
)

type TransactionInnerPayload struct {
	Payload TransactionInnerPayloadImpl
}

func (txn *TransactionInnerPayload) PayloadType() TransactionPayloadVariant {
	return TransactionPayloadVariantPayload
}

func (txn *TransactionInnerPayload) MarshalBCS(ser *bcs.Serializer) {
	if txn == nil || txn.Payload == nil {
		ser.SetError(errors.New("nil transaction payload"))
		return
	}
	ser.Uleb128(uint32(txn.Payload.InnerPayloadType()))
	ser.Struct(txn.Payload)
}

func (txn *TransactionInnerPayload) UnmarshalBCS(des *bcs.Deserializer) {
	innerType := TransactionInnerPayloadVariant(des.Uleb128())
	switch innerType {
	case TransactionInnerPayloadVariantV1:
		txn.Payload = &TransactionInnerPayloadV1{}
	default:
		des.SetError(errors.New("unknown transaction inner payload variant"))
		return
	}
	txn.Payload.UnmarshalBCS(des)
}

type TransactionInnerPayloadImpl interface {
	bcs.Struct

	InnerPayloadType() TransactionInnerPayloadVariant
}

type TransactionInnerPayloadV1 struct {
	Executable  TransactionExecutable
	ExtraConfig TransactionExtraConfig
}

func (txn *TransactionInnerPayloadV1) InnerPayloadType() TransactionInnerPayloadVariant {
	return TransactionInnerPayloadVariantV1
}

func (txn *TransactionInnerPayloadV1) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(&txn.Executable)
	ser.Struct(&txn.ExtraConfig)
}

func (txn *TransactionInnerPayloadV1) UnmarshalBCS(des *bcs.Deserializer) {
	txn.Executable.UnmarshalBCS(des)
	txn.ExtraConfig.UnmarshalBCS(des)
}

type TransactionExecutable struct {
	Inner TransactionExecutableImpl
}

func (txn *TransactionExecutable) MarshalBCS(ser *bcs.Serializer) {
	if txn == nil || txn.Inner == nil {
		ser.SetError(errors.New("nil transaction executable"))
		return
	}
	ser.Uleb128(uint32(txn.Inner.ExecutableType()))
	ser.Struct(txn.Inner)
}

func (txn *TransactionExecutable) UnmarshalBCS(des *bcs.Deserializer) {
	innerType := TransactionExecutableVariant(des.Uleb128())
	switch innerType {
	case TransactionExecutableVariantScript:
		txn.Inner = &Script{}
	case TransactionExecutableVariantEntryFunction:
		txn.Inner = &EntryFunction{}
	case TransactionExecutableVariantEmpty:
		txn.Inner = &TransactionExecutableEmpty{}
	default:
		des.SetError(errors.New("unknown transaction executable variant"))
		return
	}
	txn.Inner.UnmarshalBCS(des)
}

type TransactionExecutableImpl interface {
	bcs.Struct

	ExecutableType() TransactionExecutableVariant
}

type TransactionExecutableEmpty struct{}

func (txn *TransactionExecutableEmpty) ExecutableType() TransactionExecutableVariant {
	return TransactionExecutableVariantEmpty
}

func (txn *TransactionExecutableEmpty) MarshalBCS(*bcs.Serializer)     {}
func (txn *TransactionExecutableEmpty) UnmarshalBCS(*bcs.Deserializer) {}

type TransactionExtraConfig struct {
	Inner TransactionExtraConfigImpl
}

func (txn *TransactionExtraConfig) MarshalBCS(ser *bcs.Serializer) {
	if txn == nil || txn.Inner == nil {
		ser.SetError(errors.New("nil transaction extra config"))
		return
	}
	ser.Uleb128(uint32(txn.Inner.ConfigType()))
	ser.Struct(txn.Inner)
}

func (txn *TransactionExtraConfig) UnmarshalBCS(des *bcs.Deserializer) {
	innerType := TransactionExtraConfigVariant(des.Uleb128())
	switch innerType {
	case TransactionExtraConfigVariantV1:
		txn.Inner = &TransactionExtraConfigV1{}
	default:
		des.SetError(errors.New("unknown transaction extra config variant"))
		return
	}
	des.Struct(txn.Inner)
}

type TransactionExtraConfigImpl interface {
	bcs.Struct

	ConfigType() TransactionExtraConfigVariant
}

type TransactionExtraConfigV1 struct {
	MultisigAddress       *AccountAddress // Optional
	ReplayProtectionNonce *uint64         // Optional
}

func (txn *TransactionExtraConfigV1) ConfigType() TransactionExtraConfigVariant {
	return TransactionExtraConfigVariantV1
}

func (txn *TransactionExtraConfigV1) MarshalBCS(ser *bcs.Serializer) {
	bcs.SerializeOption(ser, txn.MultisigAddress, func(ser *bcs.Serializer, item AccountAddress) {
		ser.Struct(&item)
	})
	bcs.SerializeOption(ser, txn.ReplayProtectionNonce, func(ser *bcs.Serializer, item uint64) {
		ser.U64(item)
	})
}

func (txn *TransactionExtraConfigV1) UnmarshalBCS(des *bcs.Deserializer) {
	bcs.DeserializeOption(des, func(des *bcs.Deserializer, out *AccountAddress) {
		des.Struct(out)
	})
	bcs.DeserializeOption(des, func(des *bcs.Deserializer, out *uint64) {
		*out = des.U64()
	})
}
