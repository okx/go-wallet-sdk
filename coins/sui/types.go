package sui

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrInvalidStruct              = errors.New("invalid struct tag")
	ErrEmptyCoins                 = errors.New("empty coins")
	ErrEmptyTokens                = errors.New("empty tokens")
	ErrInvalidUint64              = errors.New("invalid uint64")
	ErrInvalidString              = errors.New("invalid string")
	ErrInvalidPure                = errors.New("invalid pure")
	ErrInvalidBytes               = errors.New("invalid bytes")
	ErrInvalidSharedOrOwnedObject = errors.New("invalid shared or owned object")
	ErrInvalidMoveCall            = errors.New("invalid move call target")
	ErrInvalidParams              = errors.New("invalid params")
)

type Type string

const (
	Empty         = Type("")
	Transfer      = Type("transfer")      //transfer sui objects
	Split         = Type("split")         // split sui objects
	Merge         = Type("merge")         // merge sui objects
	Stake         = Type("stake")         //stake sui
	WithdrawStake = Type("withdrawStake") //withdraw stake of sui
	Raw           = Type("raw")           //sign the txbytes
)

type SuiObjectRef struct {
	Digest   string `json:"digest"`
	ObjectId string `json:"objectId"`
	Version  uint64 `json:"version"`
}

type SignedTransaction struct {
	TxBytes   string `json:"tx_bytes,omitempty"`
	Signature string `json:"signature,omitempty"`
}

func (s *SuiObjectRef) Check() bool {
	if len(s.ObjectId) == 0 || len(s.Digest) == 0 {
		return false
	}
	if _, err := DecodeHexString(s.ObjectId); err != nil {
		return false
	}
	if _, err := base64.StdEncoding.DecodeString(s.Digest); err != nil {
		return false
	}
	return true
}

func (s *SuiObjectRef) Write(buf *bytes.Buffer) error {
	if err := WriteAddress(buf, s.ObjectId); err != nil {
		return err
	}
	if err := WriteUint64(buf, s.Version); err != nil {
		return err
	}
	if err := WriteHash(buf, s.Digest); err != nil {
		return err
	}
	return nil
}

type Input struct {
	Kind  string      `json:"kind"`
	Value interface{} `json:"value"`
	Index uint16      `json:"index"`
	Type  string      `json:"type"`
}

func toUint64(i interface{}) (uint64, error) {
	switch i.(type) {
	case uint64:
		return i.(uint64), nil
	case float64:
		return uint64(i.(float64)), nil
	default:
		return 0, ErrInvalidUint64
	}
}

func toString(i interface{}) (string, error) {
	switch i.(type) {
	case string:
		return i.(string), nil
	default:
		return "", ErrInvalidString
	}
}

func toBool(i interface{}) (bool, error) {
	switch i.(type) {
	case string:
		return i.(string) == "true", nil
	case bool:
		return i.(bool), nil
	case int64:
		return i.(int64) != 0, nil
	case float64:
		return i.(float64) != 0, nil
	default:
		return false, ErrInvalidString
	}
}

func toBytes(i interface{}) ([]byte, error) {
	ii, ok := i.([]interface{})
	if !ok {
		return nil, ErrInvalidBytes
	}
	b := make([]byte, len(ii))
	for k, v := range ii {
		switch v.(type) {
		case uint64:
			b[k] = uint8(v.(uint64))
		case float64:
			b[k] = uint8(v.(float64))
		default:
			return nil, ErrInvalidBytes
		}
	}
	return b, nil
}
func (s *Input) Write(buf *bytes.Buffer) error {
	switch s.Kind {
	case "Input":
		switch s.Type {
		case "pure":
			//Array(3) [Pure,				Object,				ObjVec]
			if err := WriteUint8(buf, 0); err != nil {
				return err
			}
			switch s.Value.(type) {
			case uint64:
				ss := s.Value.(uint64)
				if err := WriteUint8(buf, 8); err != nil {
					return err
				}
				if err := WriteUint64(buf, ss); err != nil {
					return err
				}
			case float64:
				if err := WriteUint8(buf, 8); err != nil {
					return err
				}
				ss := uint64(s.Value.(float64))
				if err := WriteUint64(buf, ss); err != nil {
					return err
				}
			case []byte:
				d := s.Value.([]byte)
				if err := WriteUint8(buf, uint8(len(d))); err != nil {
					return err
				}
				if _, err := buf.Write(d); err != nil {
					return err
				}
			case map[string]interface{}:
				m := s.Value.(map[string]interface{})
				if !containKeys(m, "Pure") {
					return ErrInvalidPure
				}
				b, err := toBytes(m["Pure"])
				if err != nil {
					return err
				}
				if err := WriteLen(buf, len(b)); err != nil {
					return err
				}
				buf.Write(b)
			case string:
				ss := s.Value.(string)
				if err := WriteUint8(buf, 32); err != nil {
					return err
				}
				if err := WriteAddress(buf, ss); err != nil {
					return err
				}
			default:
				return ErrInvalidPure
			}
		case "object":
			switch s.Value.(type) {
			case map[string]interface{}:
				o, ok := s.Value.(map[string]interface{})
				if !ok {
					return ErrInvalidStruct
				}
				imm, err := fromImmOrOwnedJson(o)
				if err != nil {
					return err
				}
				//Array(3) [Pure,				Object,				ObjVec]
				if err := WriteUint8(buf, 1); err != nil {
					return err
				}
				//	[  "ImmOrOwned", "Shared"]
				if imm.Object.ImmOrOwned != nil {
					if err := WriteUint8(buf, 0); err != nil {
						return err
					}
					if err := imm.Object.ImmOrOwned.Write(buf); err != nil {
						return err
					}
				} else {
					if err := WriteUint8(buf, 1); err != nil {
						return err
					}
					if err := imm.Object.Shared.Write(buf); err != nil {
						return err
					}
				}
			case *SuiObject:
				o, ok := s.Value.(*SuiObject)
				if !ok || o == nil {
					return ErrInvalidStruct
				}
				if !o.Valid() {
					return ErrInvalidStruct
				}
				//Array(3) [Pure,				Object,				ObjVec]
				if err := WriteUint8(buf, 1); err != nil {
					return err
				}
				//	[  "ImmOrOwned", "Shared"]
				if o.Object.ImmOrOwned != nil {
					if err := WriteUint8(buf, 0); err != nil {
						return err
					}
					if err := o.Object.ImmOrOwned.Write(buf); err != nil {
						return err
					}
				} else {
					if err := WriteUint8(buf, 1); err != nil {
						return err
					}
					if err := o.Object.Shared.Write(buf); err != nil {
						return err
					}
				}
			}
		default:
			return ErrInvalidStruct
		}
	default:
		return ErrInvalidStruct
	}
	return nil
}

type Shared struct {
	Mutable              bool   `json:"mutable"`
	InitialSharedVersion uint64 `json:"initialSharedVersion"`
	ObjectId             string `json:"objectId"`
}

func (s *Shared) Write(buf *bytes.Buffer) error {
	if err := WriteAddress(buf, s.ObjectId); err != nil {
		return err
	}
	if err := WriteUint64(buf, s.InitialSharedVersion); err != nil {
		return err
	}
	if err := WriteBool(buf, s.Mutable); err != nil {
		return err
	}
	return nil
}

type SuiObject struct {
	Object *SharedOrOwned `json:"Object"`
}

func (s *SuiObject) Valid() bool {
	if s == nil {
		return false
	}
	if s.Object.ImmOrOwned == nil && s.Object.Shared == nil {
		return false
	}
	return true
}

func NewImmOrOwnedObject(ref *SuiObjectRef) *SuiObject {
	return &SuiObject{Object: &SharedOrOwned{
		ImmOrOwned: ref,
	}}
}

type SharedOrOwned struct {
	ImmOrOwned *SuiObjectRef `json:"ImmOrOwned"`
	Shared     *Shared       `json:"Shared"`
}

func containKeys(m map[string]interface{}, keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}

func fromImmOrOwnedJson(o interface{}) (*SuiObject, error) {
	if o == nil {
		return nil, ErrInvalidSharedOrOwnedObject
	}
	var err error
	switch o.(type) {
	case map[string]interface{}:
		m := o.(map[string]interface{})
		if !containKeys(m, "Object") {
			return nil, ErrInvalidSharedOrOwnedObject
		}
		m2, ok := m["Object"].(map[string]interface{})
		if !ok {
			return nil, ErrInvalidSharedOrOwnedObject
		}
		if !containKeys(m2, "ImmOrOwned") {
			if !containKeys(m2, "Shared") {
				return nil, ErrInvalidSharedOrOwnedObject
			}
			m3, ok := m2["Shared"].(map[string]interface{})
			if !ok {
				return nil, ErrInvalidSharedOrOwnedObject
			}
			if !containKeys(m3, "mutable", "initialSharedVersion", "objectId") {
				return nil, ErrInvalidSharedOrOwnedObject
			}
			ref := &Shared{}
			if ref.Mutable, err = toBool(m3["mutable"]); err != nil {
				return nil, err
			}
			if ref.InitialSharedVersion, err = toUint64(m3["initialSharedVersion"]); err != nil {
				return nil, err
			}
			if ref.ObjectId, err = toString(m3["objectId"]); err != nil {
				return nil, err
			}
			return &SuiObject{Object: &SharedOrOwned{Shared: ref}}, nil
		}
		m3, ok := m2["ImmOrOwned"].(map[string]interface{})
		if !ok {
			return nil, ErrInvalidSharedOrOwnedObject
		}
		if !containKeys(m3, "digest", "version", "objectId") {
			return nil, ErrInvalidSharedOrOwnedObject
		}
		ref := &SuiObjectRef{}
		if ref.Digest, err = toString(m3["digest"]); err != nil {
			return nil, err
		}
		if ref.Version, err = toUint64(m3["version"]); err != nil {
			return nil, err
		}
		if ref.ObjectId, err = toString(m3["objectId"]); err != nil {
			return nil, err
		}
		return &SuiObject{Object: &SharedOrOwned{ImmOrOwned: ref}}, nil

	default:
		return nil, ErrInvalidSharedOrOwnedObject
	}
}

type Command struct {
	Kind          string        `json:"kind"`
	Target        string        `json:"target,omitempty"`
	Coin          *Coin         `json:"coin,omitempty"`
	Destination   *Coin         `json:"destination,omitempty"`
	Sources       []*Coin       `json:"sources,omitempty"`
	Amounts       []*Amount     `json:"amounts,omitempty"`
	Objects       []*Object     `json:"objects,omitempty"`
	Address       *Address      `json:"address,omitempty"`
	Arguments     []*Coin       `json:"arguments,omitempty"`
	TypeArguments []interface{} `json:"typeArguments,omitempty"`
}

func (s *Command) Write(buf *bytes.Buffer) error {
	//Array(6) [MoveCall,	TransferObjects,	SplitCoin,	MergeCoins,	Publish,	MakeMoveVec]
	switch s.Kind {
	case "SplitCoins":
		if err := WriteUint8(buf, 2); err != nil {
			return err
		}
		if err := s.Coin.Write(buf); err != nil {
			return err
		}
		if err := WriteLen(buf, len(s.Amounts)); err != nil {
			return err
		}
		for _, v := range s.Amounts {
			if err := v.Write(buf); err != nil {
				return err
			}
		}
	case "TransferObjects":
		if err := WriteUint8(buf, 1); err != nil {
			return err
		}
		if err := WriteLen(buf, len(s.Objects)); err != nil {
			return err
		}
		for _, v := range s.Objects {
			if err := v.Write(buf); err != nil {
				return err
			}
		}
		//Array(4) [GasCoin,			Input,			Result,			NestedResult]
		if err := WriteUint8(buf, 1); err != nil {
			return err
		}
		if err := WriteUint16(buf, s.Address.Index); err != nil {
			return err
		}
	case "MergeCoins":
		if err := WriteUint8(buf, 3); err != nil {
			return err
		}
		if err := s.Destination.Write(buf); err != nil {
			return err
		}
		if err := WriteLen(buf, len(s.Sources)); err != nil {
			return err
		}
		for _, v := range s.Sources {
			if err := v.Write(buf); err != nil {
				return err
			}
		}
	case "MoveCall":
		if err := WriteUint8(buf, 0); err != nil {
			return err
		}
		args := strings.Split(s.Target, "::")
		if len(args) != 3 {
			return ErrInvalidMoveCall
		}
		//package
		if err := WriteAddress(buf, args[0]); err != nil {
			return err
		}
		//	module
		if err := WriteString(buf, args[1]); err != nil {
			return err
		}
		//function
		if err := WriteString(buf, args[2]); err != nil {
			return err
		}
		//type_arguments
		if err := WriteLen(buf, 0); err != nil {
			return err
		}
		//arguments
		if err := WriteLen(buf, len(s.Arguments)); err != nil {
			return err
		}
		for _, v := range s.Arguments {
			if err := v.Write(buf); err != nil {
				return err
			}
		}

	default:
		return ErrInvalidStruct
	}
	return nil
}

type Coin struct {
	Kind  string `json:"kind"`
	Value string `json:"value,omitempty"`
	Index uint16 `json:"index"`
	Type  string `json:"type,omitempty"`
}

func (s *Coin) Write(buf *bytes.Buffer) error {
	switch s.Kind {
	//Array(4) [GasCoin,	//Input,	//Result,	//NestedResult]
	case "GasCoin":
		if err := WriteUint8(buf, 0); err != nil {
			return err
		}
	case "Input":
		if err := WriteUint8(buf, 1); err != nil {
			return err
		}
		if err := WriteUint16(buf, s.Index); err != nil {
			return err
		}
	case "Result":
		if err := WriteUint8(buf, 2); err != nil {
			return err
		}
		if err := WriteUint16(buf, s.Index); err != nil {
			return err
		}
	default:
		return ErrInvalidStruct
	}
	return nil
}

type Amount struct {
	Kind  string `json:"kind"`
	Value uint64 `json:"value"`
	Index uint16 `json:"index"`
	Type  string `json:"type"`
}

func (s *Amount) Write(buf *bytes.Buffer) error {
	switch s.Kind {
	//[	//  "GasCoin",	//  "Input",	//  "Result",	//  "NestedResult"	//]
	case "Input":
		if err := WriteUint8(buf, 1); err != nil {
			return err
		}
		if err := WriteUint16(buf, s.Index); err != nil {
			return err
		}
	default:
		return ErrInvalidStruct
	}
	return nil
}

type Object struct {
	Kind        string `json:"kind"`
	Index       uint16 `json:"index"`
	ResultIndex uint16 `json:"resultIndex"`
}

func (s *Object) Write(buf *bytes.Buffer) error {
	switch s.Kind {
	//Array(4) [GasCoin,		Input,		Result,		NestedResult]
	case "Result":
		if err := WriteUint8(buf, 2); err != nil {
			return err
		}
		if err := WriteUint16(buf, s.Index); err != nil {
			return err
		}
	case "NestedResult":
		if err := WriteUint8(buf, 3); err != nil {
			return err
		}
		if err := WriteUint16(buf, s.Index); err != nil {
			return err
		}
		if err := WriteUint16(buf, s.ResultIndex); err != nil {
			return err
		}
	default:
		return ErrInvalidStruct
	}
	return nil
}

type Address struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
	Index uint16 `json:"index"`
	Type  string `json:"type"`
}

func (s *Address) Write(buf *bytes.Buffer) error {
	switch s.Kind {
	case "Input":
		//Array(4) [GasCoin,		Input,		Result,		NestedResult]
		if err := WriteUint8(buf, 1); err != nil {
			return err
		}
		if err := WriteUint8(buf, 32); err != nil {
			return err
		}
		if err := WriteAddress(buf, s.Value); err != nil {
			return err
		}
	default:
		return ErrInvalidStruct
	}
	return nil
}

type Expiration struct {
	Epoch uint64 `json:"Epoch"`
}
type GasConfig struct {
	Price   string          `json:"price"`
	Budget  string          `json:"budget"`
	Payment []*SuiObjectRef `json:"payment"`
}

type Pay struct {
	Sender     string      `json:"sender"`
	Expiration *Expiration `json:"expiration,omitempty"`
	GasConfig  GasConfig   `json:"gasConfig"`
	Inputs     []*Input    `json:"inputs"`
	Commands   []*Command  `json:"transactions"`
}

func BuildTx(from, to string, coins []*SuiObjectRef, amount uint64, epoch, gasBudget, gasPrice uint64) ([]byte, error) {
	pay := Pay{Sender: from, Expiration: &Expiration{Epoch: epoch},
		GasConfig: GasConfig{Price: fmt.Sprintf("%d", gasPrice), Budget: fmt.Sprintf("%d", gasBudget), Payment: coins},
		Inputs:    []*Input{&Input{Kind: "Input", Value: amount, Index: 0, Type: "pure"}, &Input{Kind: "Input", Value: to, Index: 1, Type: "pure"}},
		Commands: []*Command{&Command{Kind: "SplitCoins", Coin: &Coin{Kind: "GasCoin"}, Amounts: []*Amount{{Kind: "Input", Value: amount, Index: 0, Type: "pure"}}},
			{Kind: "TransferObjects", Objects: []*Object{{Kind: "Result", Index: 0}}, Address: &Address{Kind: "Input", Value: to, Index: 1, Type: "pure"}}},
	}
	return pay.Build()
}

func BuildStakeTx(from, validator string, coins []*SuiObjectRef, amount uint64, epoch, gasBudget, gasPrice uint64) ([]byte, error) {
	if len(coins) == 0 || amount < 1000000000 {
		return nil, ErrInvalidParams
	}
	inputs := []*Input{{Kind: "Input", Value: amount, Index: 0, Type: "pure"},
		{Kind: "Input", Value: &SuiObject{Object: &SharedOrOwned{Shared: &Shared{Mutable: true, InitialSharedVersion: 1, ObjectId: "0x0000000000000000000000000000000000000000000000000000000000000005"}}}, Index: 1, Type: "object"},
		{Kind: "Input", Value: validator, Index: 2, Type: "pure"},
	}
	args := []*Coin{{Kind: "Input", Value: "0x0000000000000000000000000000000000000000000000000000000000000005", Index: 1, Type: "object"},
		{Kind: "Result", Index: 0}, {Kind: "Input", Value: validator, Index: 2, Type: "pure"},
	}
	pay := Pay{Sender: from, Expiration: &Expiration{Epoch: epoch},
		GasConfig: GasConfig{Price: fmt.Sprintf("%d", gasPrice), Budget: fmt.Sprintf("%d", gasBudget), Payment: coins},
		Inputs:    inputs,
		Commands: []*Command{{Kind: "SplitCoins", Coin: &Coin{Kind: "GasCoin"}, Amounts: []*Amount{{Kind: "Input", Value: amount, Index: 0, Type: "pure"}}},
			{Kind: "MoveCall", Target: "0x3::sui_system::request_add_stake", Arguments: args}},
	}
	return pay.Build()
}

func BuildWithdrawStakeTx(from string, coins []*SuiObjectRef, stake *SuiObjectRef, epoch, gasBudget, gasPrice uint64) ([]byte, error) {
	if len(coins) == 0 || stake == nil {
		return nil, ErrInvalidParams
	}
	inputs := []*Input{
		{Kind: "Input", Value: &SuiObject{Object: &SharedOrOwned{Shared: &Shared{Mutable: true, InitialSharedVersion: 1, ObjectId: "0x0000000000000000000000000000000000000000000000000000000000000005"}}}, Index: 0, Type: "object"},
		{Kind: "Input", Value: &SuiObject{Object: &SharedOrOwned{ImmOrOwned: stake}}, Index: 1, Type: "object"},
	}
	args := []*Coin{{Kind: "Input", Value: "0x0000000000000000000000000000000000000000000000000000000000000005", Index: 0, Type: "object"},
		{Kind: "Input", Value: stake.ObjectId, Index: 1, Type: "object"},
	}
	pay := Pay{Sender: from, Expiration: &Expiration{Epoch: epoch},
		GasConfig: GasConfig{Price: fmt.Sprintf("%d", gasPrice), Budget: fmt.Sprintf("%d", gasBudget), Payment: coins},
		Inputs:    inputs,
		Commands:  []*Command{{Kind: "MoveCall", Target: "0x3::sui_system::request_withdraw_stake", Arguments: args}},
	}
	return pay.Build()
}

func BuildSplitTx(from, to string, coins []*SuiObjectRef, amounts []uint64, epoch, gasBudget, gasPrice uint64) ([]byte, error) {
	if len(amounts) < 2 {
		return nil, errors.New("too few amounts")
	}
	if len(amounts) > 1023 {
		return nil, errors.New("too more amounts")
	}
	inputs := make([]*Input, 0, len(amounts)+1)
	as := make([]*Amount, 0, len(amounts))
	objects := make([]*Object, 0, len(amounts))
	for k, v := range amounts {
		inputs = append(inputs, &Input{Kind: "Input", Value: v, Index: uint16(k), Type: "pure"})
		as = append(as, &Amount{Kind: "Input", Value: v, Index: uint16(k), Type: "pure"})
		objects = append(objects, &Object{Kind: "NestedResult", ResultIndex: uint16(k)})
	}
	inputs = append(inputs, &Input{Kind: "Input", Value: to, Index: uint16(len(inputs)), Type: "pure"})
	pay := Pay{Sender: from, Expiration: &Expiration{Epoch: epoch},
		GasConfig: GasConfig{Price: fmt.Sprintf("%d", gasPrice), Budget: fmt.Sprintf("%d", gasBudget), Payment: coins},
		Inputs:    inputs,
		Commands: []*Command{&Command{Kind: "SplitCoins", Coin: &Coin{Kind: "GasCoin"}, Amounts: as},
			{Kind: "TransferObjects", Objects: objects, Address: &Address{Kind: "Input", Value: to, Index: uint16(len(objects)), Type: "pure"}}},
	}
	return pay.Build()
}

func BuildMulTx(from string, coin []*SuiObjectRef, amounts map[string]uint64, epoch, gasBudget, gasPrice uint64) ([]byte, error) {
	if len(amounts) < 1 {
		return nil, errors.New("too few amounts")
	}
	if len(amounts) > 512 {
		return nil, errors.New("too more amounts")
	}
	inputs := make([]*Input, 0, len(amounts)+1)
	as := make([]*Amount, 0, len(amounts))
	objects := make([][]*Object, 0, len(amounts))
	cmds := make([]*Command, 0, len(amounts)+1)
	addrs := make([]*Address, 0, len(amounts))
	keys := make(sort.StringSlice, 0, len(amounts))
	for k, _ := range amounts {
		keys = append(keys, k)
	}
	for k, v := range keys {
		inputs = append(inputs, &Input{Kind: "Input", Value: amounts[v], Index: uint16(k), Type: "pure"})
		as = append(as, &Amount{Kind: "Input", Value: amounts[v], Index: uint16(k), Type: "pure"})
		objects = append(objects, []*Object{{Kind: "NestedResult", Index: 0, ResultIndex: uint16(k)}})
	}
	for k, v := range keys {
		inputs = append(inputs, &Input{Kind: "Input", Value: v, Index: uint16(k + len(keys)), Type: "pure"})
		addrs = append(addrs, &Address{Kind: "Input", Value: v, Index: uint16(k + len(keys)), Type: "pure"})
	}
	cmds = append(cmds, &Command{Kind: "SplitCoins", Coin: &Coin{Kind: "GasCoin"}, Amounts: as})
	for k, _ := range keys {
		cmds = append(cmds, &Command{Kind: "TransferObjects", Objects: objects[k], Address: addrs[k]})
	}
	pay := Pay{Sender: from, Expiration: &Expiration{Epoch: epoch},
		GasConfig: GasConfig{Price: fmt.Sprintf("%d", gasPrice), Budget: fmt.Sprintf("%d", gasBudget), Payment: coin},
		Inputs:    inputs,
		Commands:  cmds,
	}
	return pay.Build()
}

func BuildMergeTx(from string, coin []*SuiObjectRef, objects []*SuiObjectRef, epoch, gasBudget, gasPrice uint64) ([]byte, error) {
	if len(objects) < 2 {
		return nil, errors.New("too few objects")
	}
	if len(objects) > 1023 {
		return nil, errors.New("too more amounts")
	}
	inputs := make([]*Input, 0, len(objects))
	cs := make([]*Coin, 0, len(objects))
	for k, v := range objects {
		inputs = append(inputs, &Input{Kind: "Input", Value: NewImmOrOwnedObject(v), Index: uint16(k), Type: "object"})
		cs = append(cs, &Coin{Kind: "Input", Value: v.ObjectId, Index: uint16(k), Type: "object"})
	}
	pay := Pay{Sender: from, Expiration: &Expiration{Epoch: epoch},
		GasConfig: GasConfig{Price: fmt.Sprintf("%d", gasPrice), Budget: fmt.Sprintf("%d", gasBudget), Payment: coin},
		Inputs:    inputs,
		Commands:  []*Command{{Kind: "MergeCoins", Destination: cs[len(cs)-1], Sources: cs[0 : len(cs)-1]}}}
	return pay.Build()
}

func BuildTokenTx(from, to string, coins []*SuiObjectRef, tokens []*SuiObjectRef, amount uint64, epoch, gasBudget, gasPrice uint64) ([]byte, error) {
	if len(coins) == 0 {
		return nil, ErrEmptyCoins
	}
	if len(tokens) == 0 {
		return nil, ErrEmptyTokens
	}
	for _, v := range tokens {
		if v == nil {
			return nil, ErrEmptyTokens
		}
	}
	for _, v := range coins {
		if v == nil {
			return nil, ErrEmptyCoins
		}
	}
	if len(tokens) == 1 {
		pay := Pay{Sender: from, Expiration: &Expiration{Epoch: epoch},
			GasConfig: GasConfig{Price: fmt.Sprintf("%d", gasPrice), Budget: fmt.Sprintf("%d", gasBudget), Payment: coins},
			Inputs: []*Input{{Kind: "Input", Value: NewImmOrOwnedObject(tokens[0]), Index: 0, Type: "object"}, //Object to be split.
				{Kind: "Input", Value: amount, Index: 1, Type: "pure"}, // amount
				{Kind: "Input", Value: to, Index: 2, Type: "pure"}},    //to
			Commands: []*Command{{Kind: "SplitCoins", Coin: &Coin{Kind: "Input", Type: "object", Index: 0, Value: tokens[0].ObjectId}, Amounts: []*Amount{{Kind: "Input", Value: amount, Index: 1, Type: "pure"}}},
				{Kind: "TransferObjects", Objects: []*Object{{Kind: "Result", Index: 0}}, Address: &Address{Kind: "Input", Value: to, Index: 2, Type: "pure"}}},
		}
		return pay.Build()
	}
	inputs := make([]*Input, len(tokens))
	cs := make([]*Coin, len(tokens))
	for k, v := range tokens {
		inputs[k] = &Input{Kind: "Input", Value: NewImmOrOwnedObject(v), Index: uint16(k), Type: "object"}
		cs[k] = &Coin{Kind: "Input", Value: v.ObjectId, Index: uint16(k), Type: "object"}
	}
	l := len(inputs)
	inputs = append(inputs, &Input{Kind: "Input", Value: amount, Index: uint16(l), Type: "pure"}, &Input{Kind: "Input", Value: to, Index: uint16(l + 1), Type: "pure"})
	pay := Pay{Sender: from, Expiration: &Expiration{Epoch: epoch},
		GasConfig: GasConfig{Price: fmt.Sprintf("%d", gasPrice), Budget: fmt.Sprintf("%d", gasBudget), Payment: coins},
		Inputs:    inputs,
		Commands: []*Command{{Kind: "MergeCoins", Destination: cs[len(cs)-1], Sources: cs[0 : len(cs)-1]},
			{Kind: "SplitCoins", Coin: &Coin{Kind: "Input", Value: cs[len(cs)-1].Value, Index: uint16(l - 1), Type: "object"}, Amounts: []*Amount{{Kind: "Input", Value: amount, Index: uint16(l), Type: "pure"}}},
			{Kind: "TransferObjects", Objects: []*Object{{Kind: "Result", Index: 1}}, Address: &Address{Kind: "Input", Value: to, Index: uint16(l + 1), Type: "pure"}}},
	}
	return pay.Build()
}

func (p *Pay) Build() ([]byte, error) {
	var b bytes.Buffer
	//Array(1) [V1]
	if err := WriteUint8(&b, 0); err != nil {
		return nil, err
	}
	//Array(4) [kind, sender,gasData,expiration]
	//Array(4) [ProgrammableTransaction, ChangeEpoch, Genesis, ConsensusCommitPrologue]
	if err := WriteUint8(&b, 0); err != nil {
		return nil, err
	}
	if err := WriteLen(&b, len(p.Inputs)); err != nil {
		return nil, err
	}
	for _, v := range p.Inputs {
		if err := v.Write(&b); err != nil {
			return nil, err
		}
	}
	if err := WriteLen(&b, len(p.Commands)); err != nil {
		return nil, err
	}
	for _, v := range p.Commands {
		if err := v.Write(&b); err != nil {
			return nil, err
		}
	}
	if err := WriteAddress(&b, p.Sender); err != nil {
		return nil, err
	}
	if err := WriteLen(&b, len(p.GasConfig.Payment)); err != nil {
		return nil, err
	}
	//export type TransactionExpiration = {	None: null;	} | {	Epoch: number;	};
	for _, v := range p.GasConfig.Payment {
		if err := v.Write(&b); err != nil {
			return nil, err
		}
	}

	if err := WriteAddress(&b, p.Sender); err != nil {
		return nil, err
	}
	price, err := strconv.ParseUint(p.GasConfig.Price, 10, 64)
	if err != nil {
		return nil, err
	}
	if err := WriteUint64(&b, price); err != nil {
		return nil, err
	} // todo bigint le
	var budget uint64
	if len(p.GasConfig.Budget) > 0 {
		budget, err = strconv.ParseUint(p.GasConfig.Budget, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if err := WriteUint64(&b, budget); err != nil {
		return nil, err
	}
	if p.Expiration == nil || p.Expiration.Epoch == 0 {
		if err := WriteUint8(&b, 0); err != nil {
			return nil, err
		}
	} else {
		if err := WriteUint8(&b, 1); err != nil {
			return nil, err
		}
		if err := WriteUint64(&b, p.Expiration.Epoch); err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}
