package aptos_types

import (
	"encoding/hex"
	"errors"
	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
	"regexp"
	"strings"
)

var ErrNullObject = errors.New("cannot serialize null object")
var ErrInvalidTypeTag = errors.New("invalid type tag")

func BcsSerializeUint64(t uint64) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeU64(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeBool(t bool) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeBool(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeU8(t uint8) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeU8(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeU128(t serde.Uint128) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeU128(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeBytes(t []byte) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeBytes(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeFixedBytes(t []byte) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeFixedBytes(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BcsSerializeStr(t string) ([]byte, error) {
	return BcsSerializeBytes([]byte(t))
}

func BcsSerializeLen(t uint64) ([]byte, error) {
	serializer := bcs.NewSerializer()
	if err := serializer.SerializeLen(t); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

func BytesFromHex(t string) []byte {
	if strings.HasPrefix(t, "0x") {
		t = strings.TrimPrefix(t, "0x")
	}
	bytes, _ := hex.DecodeString(t)
	return bytes
}

func IsHexString(s string) bool {
	res, err := regexp.MatchString("0x[0-9a-fA-F]+", s)
	if err != nil {
		return false
	}
	return res
}
