package compatibility

import "github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"

// TypeScriptCompatible means the type is able to marshal and unmarshal in BCS to a specific type in TypeScript
type TypeScriptCompatible interface {
	TypeScriptBCSMarshaler
	TypeScriptBCSUnmarshaler
}

// TypeScriptBCSMarshaler means the type is able to marshal in BCS from a matching type in TypeScript, the comment on the function
// will need to explicit about which one.
type TypeScriptBCSMarshaler interface {
	MarshalTSBCS(ser *bcs.Marshaler)
}

// TypeScriptBCSUnmarshaler means the type is able to unmarshal in BCS from a matching type in TypeScript, the comment on the
// function will need to explicit about which one.
type TypeScriptBCSUnmarshaler interface {
	UnmarshalTSBCS(des *bcs.Deserializer)
}
