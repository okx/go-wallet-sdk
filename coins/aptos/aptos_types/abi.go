package aptos_types

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
)

/**
 * Hex encoded 32 byte Aptos account address
 */
type Address = string

type IdentifierWrapper = string

/* Move module id is a string representation of Move module.
*
* Format: `{address}::{module name}`
*
* `address` should be hex-encoded 32 byte account address that is prefixed with `0x`.
*
* Module name is case-sensitive.
*
 */
type MoveModuleId = string

type MoveAbility = string

/**
* String representation of an on-chain Move type tag that is exposed in transaction payload.
* Values:
* - bool
* - u8
* - u64
* - u128
* - address
* - signer
* - vector: `vector<{non-reference MoveTypeId}>`
* - struct: `{address}::{module_name}::{struct_name}::<{generic types}>`
*
* Vector type value examples:
* - `vector<u8>`
* - `vector<vector<u64>>`
* - `vector<0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>>`
*
* Struct type value examples:
* - `0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>
 * - `0x1::account::Account`
 *
 * Note:
 * 1. Empty chars should be ignored when comparing 2 struct tag ids.
 * 2. When used in an URL path, should be encoded by url-encoding (AKA percent-encoding).
 *
*/
type MoveType = string

// PRIVATE = 'private',
// PUBLIC = 'public',
// FRIEND = 'friend',
type MoveFunctionVisibility = string

/**
 * All bytes (Vec<u8>) data is represented as hex-encoded string prefixed with `0x` and fulfilled with
 * two hex digits per byte.
 *
 * Unlike the `Address` type, HexEncodedBytes will not trim any zeros.
 *
 */
type HexEncodedBytes = string

type MoveFunctionGenericTypeParam struct {
	Constraints []MoveAbility `json:"constraints"`
}

type MoveFunctionFullName struct {
	FullName string
	MoveFunction
}
type MoveFunction struct {
	Name              IdentifierWrapper              `json:"name"`
	Visibility        MoveFunctionVisibility         `json:"visibility"`
	IsEntry           bool                           `json:"is_entry"`
	GenericTypeParams []MoveFunctionGenericTypeParam `json:"generic_type_params"`
	Params            []MoveType                     `json:"params"`
	Return            []MoveType                     `json:"return"`
}

type MoveModuleBytecode struct {
	Bytecode HexEncodedBytes `json:"bytecode"`
	Abi      MoveModule      `json:"abi"`
}

type MoveModule struct {
	Address          Address           `json:"address"`
	Name             IdentifierWrapper `json:"name"`
	Friends          []MoveModuleId    `json:"friends"`
	ExposedFunctions []MoveFunction    `json:"exposed_functions"`
	Structs          []MoveStruct      `json:"structs"`
}

type MoveStructGenericTypeParam struct {
	Constraints []MoveAbility `json:"constraints"`
}

type MoveStructField struct {
	Name IdentifierWrapper `json:"name"`
	Type MoveType          `json:"type"`
}

type MoveStruct struct {
	Name              IdentifierWrapper            `json:"name"`
	IsNative          bool                         `json:"is_native"`
	Abilities         []MoveAbility                `json:"abilities"`
	GenericTypeParams []MoveStructGenericTypeParam `json:"generic_type_params"`
	Fields            []MoveStructField            `json:"fields"`
}

type EntryFunctionPayload struct {
	Function      string        `json:"function"`
	TypeArguments []string      `json:"type_arguments"`
	Arguments     []interface{} `json:"arguments"`
	Type          string        `json:"type"`
}

type ArgumentABI struct {
	Name    string
	TypeTag TypeTag
}

func (o *ArgumentABI) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeStr(o.Name); err != nil {
		return err
	}
	if err := o.TypeTag.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *ArgumentABI) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}
