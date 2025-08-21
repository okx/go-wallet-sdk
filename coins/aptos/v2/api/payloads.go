package api

import (
	"encoding/json"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/types"
)

// TransactionPayloadVariant is the type of payload represented in JSON
type TransactionPayloadVariant string

const (
	TransactionPayloadVariantEntryFunction TransactionPayloadVariant = "entry_function_payload" // TransactionPayloadVariantEntryFunction maps to TransactionPayloadEntryFunction
	TransactionPayloadVariantScript        TransactionPayloadVariant = "script_payload"         // TransactionPayloadVariantScript maps to TransactionPayloadScript
	TransactionPayloadVariantMultisig      TransactionPayloadVariant = "multisig_payload"       // TransactionPayloadVariantMultisig maps to TransactionPayloadMultisig
	TransactionPayloadVariantWriteSet      TransactionPayloadVariant = "write_set_payload"      // TransactionPayloadVariantWriteSet maps to TransactionPayloadWriteSet
	TransactionPayloadVariantModuleBundle  TransactionPayloadVariant = "module_bundle_payload"  // TransactionPayloadVariantModuleBundle maps to TransactionPayloadModuleBundle and is deprecated
	TransactionPayloadVariantUnknown       TransactionPayloadVariant = "unknown"                // TransactionPayloadVariantUnknown maps to TransactionPayloadUnknown for unknown types
)

// TransactionPayload is an enum of all possible transaction payloads
//
// Unknown types will have the Type set to [TransactionPayloadVariantUnknown] and the Inner set to [TransactionPayloadUnknown]
type TransactionPayload struct {
	Type  TransactionPayloadVariant // Type of the payload, if the payload isn't recognized, it will be [TransactionPayloadVariantUnknown]
	Inner TransactionPayloadImpl    // Inner is the actual payload
}

// UnmarshalJSON unmarshals the [TransactionPayload] from JSON handling conversion between types
func (o *TransactionPayload) UnmarshalJSON(b []byte) error {
	type inner struct {
		Type string `json:"type"`
	}
	data := &inner{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	o.Type = TransactionPayloadVariant(data.Type)
	switch o.Type {
	case TransactionPayloadVariantEntryFunction:
		o.Inner = &TransactionPayloadEntryFunction{}
	case TransactionPayloadVariantScript:
		o.Inner = &TransactionPayloadScript{}
	case TransactionPayloadVariantMultisig:
		o.Inner = &TransactionPayloadMultisig{}
	//case TransactionPayloadVariantWriteSet:
	//	o.Inner = &TransactionPayloadWriteSet{}
	//case TransactionPayloadVariantModuleBundle:
	//	o.Inner = &TransactionPayloadModuleBundle{}
	default:
		// Make sure it doesn't crash with new types
		o.Inner = &TransactionPayloadUnknown{Type: string(o.Type)}
		o.Type = TransactionPayloadVariantUnknown
		return json.Unmarshal(b, &o.Inner.(*TransactionPayloadUnknown).Payload)
	}
	return json.Unmarshal(b, o.Inner)
}

// TransactionPayloadImpl is all the interfaces required for all transaction payloads
//
// Current implementations are:
//
//   - [TransactionPayloadEntryFunction]
//   - [TransactionPayloadScript]
//   - [TransactionPayloadMultisig]
//   - [TransactionPayloadWriteSet]
//   - [TransactionPayloadModuleBundle]
//   - [TransactionPayloadUnknown]
type TransactionPayloadImpl interface{}

// TransactionPayloadUnknown is to handle new types gracefully.
//
// This is a fallback type for unknown transaction payloads.
type TransactionPayloadUnknown struct {
	Type    string         `json:"type"`    // Type is the actual type field from the JSON.
	Payload map[string]any `json:"payload"` // Payload is the raw JSON payload.
}

// TransactionPayloadEntryFunction describes an entry function call by a transaction.
type TransactionPayloadEntryFunction struct {
	Function      string   `json:"function"`       // Function is the name of the function called e.g. 0x1::coin::transfer
	TypeArguments []string `json:"type_arguments"` // TypeArguments are the type arguments for the function as a string representation of the TypeTag.
	Arguments     []any    `json:"arguments"`      // Arguments are the arguments for the function.  The order should match the order in the Move source.
}

// TransactionPayloadScript describes a script payload along with associated.
//
// See more information about scripts at the [MoveScript Documentation].
//
// [MoveScript Documentation]: https://aptos.dev/en/build/smart-contracts/scripts
type TransactionPayloadScript struct {
	Code          *MoveScript `json:"code"`           // Code is the Move bytecode for the script.
	TypeArguments []string    `json:"type_arguments"` // TypeArguments are the type arguments for the script as a string representation of the TypeTag.
	Arguments     []any       `json:"arguments"`      // Arguments are the arguments for the script.  The order should match the order in the Move source.
}

// TransactionPayloadMultisig describes a multi-sig running an entry function
//
// TODO: This isn't ever a top level transaction payload, it is always nested in a TransactionPayload so it may not apply here
type TransactionPayloadMultisig struct {
	MultisigAddress    *types.AccountAddress `json:"multisig_address"`              // MultisigAddress is the address of the multi-sig account
	TransactionPayload *TransactionPayload   `json:"transaction_payload,omitempty"` // TransactionPayload is the payload of the transaction, optional
}
