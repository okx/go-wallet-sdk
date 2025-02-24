package serialize

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/okx/go-wallet-sdk/util"
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
)

const (
	CreateAccountAction = iota
	DeployContractAction
	FunctionCallAction
	TransferAction
	StakeAction
	AddKeyAction
	DeleteKeyAction
	DeleteAccountAccount
)

var (
	NearPrefix    = "ed25519:"
	Ed25519Prefix = "ed25519"
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

func TryParse(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("key decode error, key =%s", s)
	}
	if !strings.HasPrefix(s, NearPrefix) {
		publicKeyByte, err := util.DecodeHexString(s)
		if err != nil {
			return nil, fmt.Errorf("key decode error, key =%s", s)
		}
		return publicKeyByte, nil
	}
	args := strings.Split(s, ":")
	if len(args) != 2 || args[0] != Ed25519Prefix {
		return nil, fmt.Errorf("key decode error, key =%s", s)
	}
	return base58.Decode(args[1]), nil
}

func TryParsePubKey(s string) (*PublicKey, error) {
	publicKeyByte, err := TryParse(s)
	if err != nil {
		return nil, fmt.Errorf("public key decode error,public key =%s", s)
	}
	if len(publicKeyByte) != 32 {
		return nil, fmt.Errorf("public key len error,public key=%s", s)
	}
	return &PublicKey{KeyType: 0, Value: publicKeyByte}, nil
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

func CreateCreateAccount() (*CreateAccount, error) {
	return &CreateAccount{
		Action: CreateAccountAction,
	}, nil
}
func (s *CreateAccount) GetActionIndex() uint8 {
	return s.Action
}

func (s *CreateAccount) Serialize() ([]byte, error) {
	return []byte{s.Action}, nil
}

type DeployContract struct {
	Action uint8
	Code   []U8 //Uint8Array
}

func CreateDeployContract(code []byte) (*DeployContract, error) {
	c := make([]U8, len(code))
	for k, v := range code {
		c[k] = U8{v}
	}
	return &DeployContract{
		Action: DeployContractAction,
		Code:   c,
	}, nil
}

func (s *DeployContract) GetActionIndex() uint8 {
	return s.Action
}

func (s *DeployContract) Serialize() ([]byte, error) {
	data := []byte{s.Action}
	argLen := len(s.Code)
	argLenU32 := U32{
		Value: uint32(argLen),
	}

	argLenBytes, err := argLenU32.Serialize()
	if err != nil {
		return nil, err
	}

	data = append(data, argLenBytes...)
	for _, v := range s.Code {
		d, err := v.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, d...)
	}
	return data, nil
}

type Stake struct {
	Action    uint8
	Stake     U128
	PublicKey PublicKey
}

func CreateStake(publicKeyHex string, allowance string) (*Stake, error) {
	pub, err := TryParsePubKey(publicKeyHex)
	if err != nil {
		return nil, err
	}
	return &Stake{
		Action:    StakeAction,
		PublicKey: *pub,
		Stake:     U128{Value: convertToBigInt(allowance)},
	}, nil
}

func (s *Stake) GetActionIndex() uint8 {
	return s.Action
}

func (s *Stake) Serialize() ([]byte, error) {
	data := []byte{s.Action}
	v, err := s.Stake.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, v...)
	d, err := s.PublicKey.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	return data, nil
}

type AddKeyAct struct {
	Action    uint8
	PublicKey PublicKey
	AccessKey AccessKey
}

func (s *AddKeyAct) GetActionIndex() uint8 {
	return s.Action
}

func (s *AddKeyAct) Serialize() ([]byte, error) {
	data := []byte{s.Action}
	d, err := s.PublicKey.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	v, err := s.AccessKey.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, v...)
	return data, nil
}

type AccessKey struct {
	Nonce      U64
	Permission AccessKeyPermission
}

func CreateAddFullAccessKey(publicKeyHex string) (*AddKeyAct, error) {
	pub, err := TryParsePubKey(publicKeyHex)
	if err != nil {
		return nil, err
	}
	return &AddKeyAct{
		Action:    AddKeyAction,
		PublicKey: *pub,
		AccessKey: AccessKey{Nonce: U64{0}, Permission: AccessKeyPermission{FullAccess: &FullAccessPermission{}}},
	}, nil
}

func convertToBigInt(v string) *big.Int {
	b := new(big.Int)
	b.SetString(v, 10)
	return b
}

func CreateAddFunctionCallAccessKey(publicKeyHex string, allowance, receiverId string, methodNames []string) (*AddKeyAct, error) {
	pub, err := TryParsePubKey(publicKeyHex)
	if err != nil {
		return nil, err
	}
	var allow *U128
	if len(allowance) > 0 {
		allow = &U128{Value: convertToBigInt(allowance)}
	}
	methods := make([]String, len(methodNames))
	for k, v := range methodNames {
		methods[k] = String{v}
	}
	return &AddKeyAct{
		Action:    AddKeyAction,
		PublicKey: *pub,
		AccessKey: AccessKey{Nonce: U64{0}, Permission: AccessKeyPermission{FunctionCall: &FunctionCallPermission{Allowance: allow, ReceiverId: String{receiverId}, MethodNames: methods}}},
	}, nil
}

func (s *AccessKey) Serialize() ([]byte, error) {
	data := []byte{}
	d, err := s.Nonce.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	v, err := s.Permission.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, v...)
	return data, nil
}

type FunctionCallPermission struct {
	Allowance   *U128
	ReceiverId  String
	MethodNames []String
}

func (s *FunctionCallPermission) Serialize() ([]byte, error) {
	data := []byte{}
	if s.Allowance == nil {
		data = append(data, byte(0))
	} else {
		data = append(data, byte(1))
		d, err := s.Allowance.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, d...)
	}
	v, err := s.ReceiverId.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, v...)

	argLen := len(s.MethodNames)
	argLenU32 := U32{
		Value: uint32(argLen),
	}

	argLenBytes, err := argLenU32.Serialize()
	if err != nil {
		return nil, err
	}

	data = append(data, argLenBytes...)

	for _, arg := range s.MethodNames {
		v, err = arg.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, v...)
	}
	return data, nil
}

type FullAccessPermission struct {
}

func (s *FullAccessPermission) Serialize() ([]byte, error) {
	return []byte{}, nil
}

type AccessKeyPermission struct {
	FunctionCall *FunctionCallPermission
	FullAccess   *FullAccessPermission
}

func (s *AccessKeyPermission) Serialize() ([]byte, error) {
	data := []byte{}
	if s.FullAccess != nil {
		data = append(data, byte(1))
	} else if s.FunctionCall != nil {
		data = append(data, byte(0))
		d, err := s.FunctionCall.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, d...)
	}
	return data, nil
}

type DeleteKeyAct struct {
	Action    uint8
	PublicKey PublicKey
}

func CreateDeleteKey(publicKeyHex string) (*DeleteKeyAct, error) {
	pub, err := TryParsePubKey(publicKeyHex)
	if err != nil {
		return nil, err
	}
	return &DeleteKeyAct{
		Action:    DeleteKeyAction,
		PublicKey: *pub,
	}, nil
}

func (s *DeleteKeyAct) GetActionIndex() uint8 {
	return s.Action
}

func (s *DeleteKeyAct) Serialize() ([]byte, error) {
	data := []byte{s.Action}
	d, err := s.PublicKey.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	return data, nil
}

type DeleteAccountAct struct {
	Action        uint8
	BeneficiaryId String
}

func (s *DeleteAccountAct) GetActionIndex() uint8 {
	return s.Action
}

func (s *DeleteAccountAct) Serialize() ([]byte, error) {
	data := []byte{s.Action}
	d, err := s.BeneficiaryId.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	return data, nil
}

func CreateDeleteAccount(beneficiaryId string) (*DeleteAccountAct, error) {
	return &DeleteAccountAct{
		Action:        DeleteAccountAccount,
		BeneficiaryId: String{Value: beneficiaryId},
	}, nil
}

type SignMessagePayload struct {
	tag         U32
	Message     String
	Nonce       []byte
	Recipient   String
	CallbackUrl *String
}

const def_tag = uint32(2147484061)

func NewSignMessagePayload(message string, nonce []byte, recipient string, callbackUrl string) *SignMessagePayload {
	var cbu *String
	if callbackUrl != "" {
		cbu = &String{Value: callbackUrl}
	}
	return &SignMessagePayload{
		tag:         U32{Value: def_tag},
		Message:     String{Value: message},
		Nonce:       nonce,
		Recipient:   String{Value: recipient},
		CallbackUrl: cbu,
	}
}

func (self *SignMessagePayload) Serialize() ([]byte, error) {
	data := make([]byte, 0)
	var err error
	data, err = serializeAndAppend(&self.tag, data)
	if err != nil {
		return nil, err
	}
	data, err = serializeAndAppend(&self.Message, data)
	if err != nil {
		return nil, err
	}
	if len(self.Nonce) != 32 {
		return nil, errors.New("Expected nonce to be a 32 bytes buffer")
	}
	data = append(data, self.Nonce...)
	data, err = serializeAndAppend(&self.Recipient, data)
	if err != nil {
		return nil, err
	}
	if self.CallbackUrl == nil {
		data = append(data, byte(0))
	} else {
		data = append(data, byte(1))
		data, err = serializeAndAppend(self.CallbackUrl, data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func serializeAndAppend(f ISerialize, data []byte) ([]byte, error) {
	d, err := f.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	return data, nil
}
