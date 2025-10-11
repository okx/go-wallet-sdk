package token

import (
	"encoding/json"
	"math/big"

	"github.com/emresenyuva/go-wallet-sdk/util/abi"
)

func ParseTrc20JsonAbi(data string) *abi.ABI {
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

var Abi20 = ParseTrc20JsonAbi(TRC20ABI)

func Transfer(to string, value *big.Int) ([]byte, error) {
	return Transact("transfer", to, value)
}

func Transact(name string, params ...interface{}) ([]byte, error) {
	input, err := Abi20.Pack(name, params...)
	if err != nil {
		return nil, err
	}
	return input, nil
}
