package api

import (
	"encoding/json"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/types"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
)

// MoveBytecode describes a module, or script, and it's associated ABI as a [MoveModule]
//
// Example 0x1::coin:
//
//	{
//		"bytecode": "0xa11ceb0b123456...",
//		"abi": {
//			"address": "0x1",
//			"name": "coin",
//			"friends": [
//				"0x1::aptos_coin",
//				"0x1::genesis",
//				"0x1::transaction_fee"
//			],
//			"exposed_functions": [
//				{
//					"name": "balance",
//					"visibility": "public",
//					"is_entry": false,
//					"is_view": true,
//					"generic_type_params": [
//						{
//							"constraints": []
//						}
//					],
//					"params": [
//						"address"
//					],
//					"return": [
//						"u64"
//					]
//				}
//			],
//			"structs": [
//				{
//					"name": "Coin",
//					"is_native": false,
//					"abilities": [
//						"store"
//					],
//					"generic_type_params": [
//						{
//							"constraints": []
//						}
//					],
//					"fields": [
//						{
//							"name": "value",
//							"type": "u64"
//						}
//					]
//				},
//			],
//		}
//	}
type MoveBytecode struct {
	Bytecode HexBytes    `json:"bytecode"`      // Bytecode is the hex encoded version of the compiled module
	Abi      *MoveModule `json:"abi,omitempty"` // Abi is the ABI for the module, and is optional
}
type HexBytes []byte

func (u *HexBytes) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	bytes, err := util.ParseHex(str)
	if err != nil {
		return err
	}
	*u = bytes
	return nil
}

// MarshalJSON serializes a JSON data into blob [HexBytes]
//
// Example:
//
//	[]byte{0x12, 0x34, 0x56} -> "0x123456"
func (u *HexBytes) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	}
	return json.Marshal(util.BytesToHex(*u))
}

// MoveComponentId is an id for a struct, function, or other type e.g. 0x1::aptos_coin::AptosCoin
type MoveComponentId = string

// MoveModule describes the abilities and types associated with a specific module.
type MoveModule struct {
	Address          *types.AccountAddress `json:"address"`           // Address is the address of the module e.g. 0x1
	Name             string                `json:"name"`              // Name is the name of the module e.g. coin
	Friends          []MoveComponentId     `json:"friends"`           // Friends are other modules that can access this module in the same package.
	ExposedFunctions []*MoveFunction       `json:"exposed_functions"` // ExposedFunctions are the functions that can be called from outside the module.
	Structs          []*MoveStruct         `json:"structs"`           // Structs are the structs defined in the module.
}

// MoveScript is the representation of a compiled script.  The API may not fill in the ABI field.
//
// Example:
//
//	{
//		"bytecode": "0xa11ceb0b123456...",
//		"abi": {
//			"address": "0x1",
//			"name": "coin",
//			"friends": [
//				"0x1::aptos_coin",
//				"0x1::genesis",
//				"0x1::transaction_fee"
//			],
//			"exposed_functions": [
//				{
//					"name": "balance",
//					"visibility": "public",
//					"is_entry": false,
//					"is_view": true,
//					"generic_type_params": [
//						{
//							"constraints": []
//						}
//					],
//					"params": [
//						"address"
//					],
//					"return": [
//						"u64"
//					]
//				}
//			],
//			"structs": [
//				{
//					"name": "Coin",
//					"is_native": false,
//					"abilities": [
//						"store"
//					],
//					"generic_type_params": [
//						{
//							"constraints": []
//						}
//					],
//					"fields": [
//						{
//							"name": "value",
//							"type": "u64"
//						}
//					]
//				}
//			]
//		}
//	}
type MoveScript struct {
	Bytecode HexBytes      `json:"bytecode"`      // Bytecode is the hex encoded version of the compiled script.
	Abi      *MoveFunction `json:"abi,omitempty"` // Abi is the ABI for the module, and is optional.
}

// MoveFunction describes a move function and its associated properties
//
// Example 0x1::coin::balance:
//
//	{
//		"name": "balance",
//		"visibility": "public",
//		"is_entry": false,
//		"is_view": true,
//		"generic_type_params": [
//			{
//				"constraints": []
//			}
//		],
//		"params": [
//			"address"
//		],
//		"return": [
//			"u64"
//		]
//	}
type MoveFunction struct {
	Name              MoveComponentId     `json:"name"`                // Name is the name of the function e.g. balance
	Visibility        MoveVisibility      `json:"visibility"`          // Visibility is the visibility of the function e.g. public
	IsEntry           bool                `json:"is_entry"`            // IsEntry is true if the function is an entry function
	IsView            bool                `json:"is_view"`             // IsView is true if the function is a view function
	GenericTypeParams []*GenericTypeParam `json:"generic_type_params"` // GenericTypeParams are the generic type parameters for the function
	Params            []string            `json:"params"`              // Params are the parameters for the function in string format for the TypeTag
	Return            []string            `json:"return"`              // Return is the return type for the function in string format for the TypeTag
}

// GenericTypeParam is a set of requirements for a generic.  These can be applied via different
// [MoveAbility] constraints required on the type.
//
// Example:
//
//	{
//		"constraints": [
//			"copy"
//		]
//	}
type GenericTypeParam struct {
	Constraints []MoveAbility `json:"constraints"` // Constraints are the constraints required for the generic type e.g. copy.
}

// MoveAbility are the types of abilities applied to structs, the possible types are listed
// as [MoveAbilityStore] and others.
//
// See more at the [Move Ability Documentation].
//
// [Move Ability Documentation]: https://aptos.dev/en/build/smart-contracts/book/abilities
type MoveAbility string

const (
	MoveAbilityStore MoveAbility = "store" // MoveAbilityStore is the ability to store the type
	MoveAbilityDrop  MoveAbility = "drop"  // MoveAbilityDrop is the ability to drop the type
	MoveAbilityKey   MoveAbility = "key"   // MoveAbilityKey is the ability to use the type as a key in global storage
	MoveAbilityCopy  MoveAbility = "copy"  // MoveAbilityCopy is the ability to copy the type
)

// MoveVisibility is the visibility of a function or struct, the possible types are listed
// as [MoveVisibilityPublic] and others
//
// See more at the [Move Visibility Documentation].
//
// [Move Visibility Documentation]: https://aptos.dev/en/build/smart-contracts/book/functions#visibility
type MoveVisibility string

const (
	MoveVisibilityPublic  MoveVisibility = "public"  // MoveVisibilityPublic is a function that is accessible anywhere
	MoveVisibilityPrivate MoveVisibility = "private" // MoveVisibilityPrivate is a function that is only accessible within the module
	MoveVisibilityFriend  MoveVisibility = "friend"  // MoveVisibilityFriend is a function that is only accessible to friends of the module
)

// MoveStruct describes the layout for a struct, and its constraints
//
// Example 0x1::coin::Coin:
//
//	{
//		"name": "Coin",
//		"is_native": false,
//		"abilities": [
//			"store"
//		],
//		"generic_type_params": [
//			{
//				"constraints": []
//			}
//		],
//		"fields": [
//			{
//				"name": "value",
//				"type": "u64"
//			}
//		]
//	}
type MoveStruct struct {
	Name              string              `json:"name"`                // Name is the name of the struct e.g. Coin
	IsNative          bool                `json:"is_native"`           // IsNative is true if the struct is native e.g. u64
	Abilities         []MoveAbility       `json:"abilities"`           // Abilities are the abilities applied to the struct e.g. copy or store
	GenericTypeParams []*GenericTypeParam `json:"generic_type_params"` // GenericTypeParams are the generic type parameters for the struct
	Fields            []*MoveStructField  `json:"fields"`              // Fields are the fields in the struct
}

// MoveStructField represents a single field in a struct, and it's associated type.
//
// Example:
//
//	{
//		"name": "value",
//		"type": "u64"
//	}
type MoveStructField struct {
	Name string `json:"name"` // Name of the field e.g. value
	Type string `json:"type"` // Type of the field in string format for the TypeTag e.g. u64
}
