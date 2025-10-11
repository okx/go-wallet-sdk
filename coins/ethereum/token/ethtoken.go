package token

import (
	"encoding/json"
	"math/big"

	"github.com/emresenyuva/go-wallet-sdk/util/abi"
)

func ParseErc20JsonAbi(data string) *abi.ABI {
	var fields []struct {
		Type      string
		Name      string
		Constant  bool
		Anonymous bool
		Inputs    []abi.Argument
		Outputs   []abi.Argument
	}
	if err := json.Unmarshal([]byte(data), &fields); err != nil {
		panic(err)
	}
	var inst = &abi.ABI{}
	inst.Methods = make(map[string]*abi.Method)
	inst.Events = make(map[string]*abi.Event)
	for _, field := range fields {
		switch field.Type {
		case "constructor":
			inst.Constructor = &abi.Method{
				Inputs: field.Inputs,
			}
		case "function":
			inst.Methods[field.Name] = &abi.Method{
				Name:    field.Name,
				Const:   field.Constant,
				Inputs:  field.Inputs,
				Outputs: field.Outputs,
			}
		case "event":
			inst.Events[field.Name] = &abi.Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    field.Inputs,
			}
		}
	}

	return inst
}

var Abi20 = ParseErc20JsonAbi(ERC20ABI)
var Abi721 = ParseErc20JsonAbi(ERC721ABI)

func Transfer(to string, value *big.Int) ([]byte, error) {
	return Transact("transfer", to, value)
}

func Approve(spender string, value *big.Int) ([]byte, error) {
	return Transact("approve", spender, value)
}

func Transfer721(from, to string, tokenId *big.Int) ([]byte, error) {
	return Transact721("safeTransferFrom", from, to, tokenId)
}

func Transact(name string, params ...interface{}) ([]byte, error) {
	input, err := Abi20.Pack(name, params...)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func Transact721(name string, params ...interface{}) ([]byte, error) {
	input, err := Abi721.Pack(name, params...)
	if err != nil {
		return nil, err
	}
	return input, nil
}
