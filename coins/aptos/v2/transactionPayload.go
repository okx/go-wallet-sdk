package v2

import (
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
)

// TransactionPayload the actual instructions of which functions to call on chain
type TransactionPayload struct {
	Payload bcs.Struct
}

const (
	TransactionPayloadScript        = 0
	TransactionPayloadModuleBundle  = 1 // Deprecated
	TransactionPayloadEntryFunction = 2
	TransactionPayloadMultisig      = 3 // TODO? defined in aptos-core/types/src/transaction/mod.rs
)

func (txn *TransactionPayload) MarshalBCS(bcs *bcs.Serializer) {
	switch p := txn.Payload.(type) {
	case *Script:
		bcs.Uleb128(TransactionPayloadScript)
		p.MarshalBCS(bcs)
	case *ModuleBundle:
		// Deprecated, this never was used in production, and we will just drop it
		bcs.SetError(fmt.Errorf("module bundle is not supported as a transaction payload"))
	case *EntryFunction:
		bcs.Uleb128(TransactionPayloadEntryFunction)
		p.MarshalBCS(bcs)
	default:
		bcs.SetError(fmt.Errorf("bad txn payload, %T", txn.Payload))
	}
}
func (txn *TransactionPayload) UnmarshalBCS(bcs *bcs.Deserializer) {
	kind := bcs.Uleb128()
	switch kind {
	case TransactionPayloadScript:
		xs := &Script{}
		xs.UnmarshalBCS(bcs)
		txn.Payload = xs
	case TransactionPayloadModuleBundle:
		// Deprecated, should never be in production
		bcs.SetError(fmt.Errorf("module bundle is not supported as a transaction payload"))
	case TransactionPayloadEntryFunction:
		xs := &EntryFunction{}
		xs.UnmarshalBCS(bcs)
		txn.Payload = xs
	default:
		bcs.SetError(fmt.Errorf("bad txn payload kind, %d", kind))
	}
}

// ModuleBundle is long deprecated and no longer used, but exist as an enum position in TransactionPayload
type ModuleBundle struct {
}

func (txn *ModuleBundle) MarshalBCS(bcs *bcs.Serializer) {
	bcs.SetError(errors.New("ModuleBundle unimplemented"))
}
func (txn *ModuleBundle) UnmarshalBCS(bcs *bcs.Deserializer) {
	bcs.SetError(errors.New("ModuleBundle unimplemented"))
}

// EntryFunction call a single published entry function
type EntryFunction struct {
	Module   ModuleId
	Function string
	ArgTypes []TypeTag
	Args     [][]byte
}

func (sf *EntryFunction) MarshalBCS(serializer *bcs.Serializer) {
	sf.Module.MarshalBCS(serializer)
	serializer.WriteString(sf.Function)
	bcs.SerializeSequence(sf.ArgTypes, serializer)
	serializer.Uleb128(uint32(len(sf.Args)))
	for _, a := range sf.Args {
		serializer.WriteBytes(a)
	}
}
func (sf *EntryFunction) UnmarshalBCS(deserializer *bcs.Deserializer) {
	sf.Module.UnmarshalBCS(deserializer)
	sf.Function = deserializer.ReadString()
	sf.ArgTypes = bcs.DeserializeSequence[TypeTag](deserializer)
	alen := deserializer.Uleb128()
	sf.Args = make([][]byte, alen)
	// go 1.22
	//for i := range alen {
	//	sf.Args[i] = deserializer.ReadBytes()
	//}
	for i := uint32(0); i < alen; i++ {
		sf.Args[i] = deserializer.ReadBytes()
	}
}
