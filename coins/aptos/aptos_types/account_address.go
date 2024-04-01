package aptos_types

import (
	"encoding/hex"
	"errors"
	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
	"strings"
)

const LENGTH = 32

type AccountAddress [LENGTH]byte

var CORE_CODE_ADDRESS, _ = FromHex("0x1")

func FromHex(a string) (*AccountAddress, error) {
	if strings.HasPrefix(a, "0x") || strings.HasPrefix(a, "0X") {
		a = a[2:]
	}
	if len(a)%2 != 0 {
		a = "0" + a
	}

	bytes, err := hex.DecodeString(a)
	if err != nil {
		return nil, err
	}
	if len(bytes) > LENGTH {
		return nil, errors.New("hex string too long")
	}

	res := AccountAddress{}
	copy(res[LENGTH-len(bytes):], bytes[:])
	return &res, nil
}
func (a *AccountAddress) ToString() string {
	//return "0x" + hex.EncodeToString(a[:])
	return "0x" + hex.EncodeToString(a[:])
}

func (a *AccountAddress) ToShortString() string {
	//return "0x" + strings.TrimLeft(hex.EncodeToString(a[:]), "0")
	return "0x" + strings.TrimLeft(hex.EncodeToString(a[:]), "0")
}
func (a *AccountAddress) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeFixedBytes(a[:]); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (a *AccountAddress) BcsSerialize() ([]byte, error) {
	if a == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := a.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}
