package serialize

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/util"
	"github.com/shopspring/decimal"
	"math/big"
)

const (
	CreateAccountAction = iota
	DeployContractAction
	FunctionCallAction
	TransferAction
	StakeAction
	AddKey
	DeleteKey
	DeleteAccount
)

type ISerialize interface {
	Serialize() ([]byte, error)
}

type IAction interface {
	ISerialize
	GetActionIndex() uint8
}

type U8 struct {
	Value uint8
}

func (u *U8) Serialize() ([]byte, error) {
	return []byte{u.Value}, nil
}

type U32 struct {
	Value uint32
}

func (u *U32) Serialize() ([]byte, error) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, u.Value)
	return data, nil
}

type U64 struct {
	Value uint64
}

func (u *U64) Serialize() ([]byte, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, u.Value)
	return data, nil
}

type U128 struct {
	Value *big.Int
}

func (u *U128) Serialize() ([]byte, error) {
	data, err := util.BigIntToUintBytes(u.Value, 16)
	if err != nil {
		return nil, err
	}
	util.Reverse(data)
	return data, nil
}

type String struct {
	Value string
}

func (s *String) Serialize() ([]byte, error) {
	if s.Value == "" {
		return nil, errors.New("string is null")
	}
	length := len(s.Value)
	uL := U32{
		Value: uint32(length),
	}

	data, err := uL.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, []byte(s.Value)...)
	return data, nil
}

type PublicKey struct {
	KeyType uint8
	Value   []byte
}

func (s *PublicKey) Serialize() ([]byte, error) {
	data := []byte{s.KeyType}
	if len(s.Value) != 32 {
		return nil, fmt.Errorf("publickey length is not equal 32,length=%d", len(s.Value))
	}
	data = append(data, s.Value...)
	return data, nil
}

type Signature struct {
	KeyType uint8
	Value   []byte
}

func (s *Signature) Serialize() ([]byte, error) {
	data := []byte{s.KeyType}
	if len(s.Value) != 64 {
		return nil, fmt.Errorf("signature length is not equal 64,length=%d", len(s.Value))
	}
	data = append(data, s.Value...)
	return data, nil
}

type BlockHash struct {
	Value []byte
}

func (s *BlockHash) Serialize() ([]byte, error) {
	if len(s.Value) != 32 {
		return nil, fmt.Errorf("blockhash length is not equal 32,length=%d", len(s.Value))
	}
	return s.Value, nil
}

// Action
type Transfer struct {
	Action uint8
	Value  U128
}

func CreateTransfer(amount string) (*Transfer, error) {
	dec, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, err
	}
	return &Transfer{
		Action: TransferAction,
		Value:  U128{Value: dec.BigInt()},
	}, nil
}
func (s *Transfer) GetActionIndex() uint8 {
	return s.Action
}
func (s *Transfer) Serialize() ([]byte, error) {
	data := []byte{s.Action}
	v, err := s.Value.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, v...)
	return data, nil
}

// FunctionCall
type FunctionCall struct {
	Action     uint8
	MethodName String
	Args       []U8
	Gas        U64
	Deposit    U128
}

func CreateFunctionCall(methodName string, args []uint8, gas, deposit *big.Int) (*FunctionCall, error) {
	argsU8 := []U8{}
	for _, arg := range args {
		argsU8 = append(argsU8, U8{Value: arg})
	}

	return &FunctionCall{
		Action:     FunctionCallAction,
		MethodName: String{Value: methodName},
		Args:       argsU8,
		Gas:        U64{Value: gas.Uint64()},
		Deposit:    U128{Value: deposit},
	}, nil
}

func (s *FunctionCall) GetActionIndex() uint8 {
	return s.Action
}

func (s *FunctionCall) Serialize() ([]byte, error) {
	data := []byte{s.Action}

	v, err := s.MethodName.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, v...)

	argLen := len(s.Args)
	argLenU32 := U32{
		Value: uint32(argLen),
	}

	argLenBytes, err := argLenU32.Serialize()
	if err != nil {
		return nil, err
	}

	data = append(data, argLenBytes...)

	for _, arg := range s.Args {
		v, err = arg.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, v...)
	}

	v, err = s.Gas.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, v...)

	v, err = s.Deposit.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, v...)

	return data, nil
}

type CreateAccount struct {
	Action uint8
}

func (s *CreateAccount) GetActionIndex() uint8 {
	return s.Action
}

func (s *CreateAccount) Serialize() ([]byte, error) {
	return []byte{s.Action}, nil
}
