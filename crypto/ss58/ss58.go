package ss58

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
)

var (
	SSPrefix = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
)

func Encode(publicKeyHash []byte, prefix []byte) (string, error) {
	if len(publicKeyHash) != 32 {
		return "", errors.New("public hash length is not equal 32")
	}
	payload := AppendBytes(prefix, publicKeyHash)
	input := AppendBytes(SSPrefix, payload)
	ck := blake2b.Sum512(input)
	checkum := ck[:2]
	address := base58.Encode(AppendBytes(payload, checkum))
	if address == "" {
		return address, errors.New("base58 encode error")
	}
	return address, nil
}

func EncodeByPubHex(publicHex string, prefix []byte) (string, error) {
	publicKeyHash, err := hex.DecodeString(publicHex)
	if err != nil {
		return "", err
	}
	return Encode(publicKeyHash, prefix)
}

func DecodeToPub(address string) ([]byte, error) {
	data := base58.Decode(address)
	if len(data) != 35 {
		return nil, errors.New("base58 decode error")
	}
	return data[1 : len(data)-2], nil
}

func Decode(address string) ([]byte, error) {
	data := base58.Decode(address)
	if len(data) != 35 {
		return nil, errors.New("base58 decode error")
	}
	return data, nil
}

func VerityAddress(address string, prefix []byte) error {
	decodeBytes := base58.Decode(address)
	if len(decodeBytes) != 35 {
		return errors.New("base58 decode error")
	}
	if decodeBytes[0] != prefix[0] {
		return errors.New("prefix valid error")
	}
	pub := decodeBytes[1 : len(decodeBytes)-2]

	data := append(prefix, pub...)
	input := append(SSPrefix, data...)
	ck := blake2b.Sum512(input)
	checkSum := ck[:2]
	for i := 0; i < 2; i++ {
		if checkSum[i] != decodeBytes[33+i] {
			return errors.New("checksum valid error")
		}
	}
	if len(pub) != 32 {
		return errors.New("decode public key length is not equal 32")
	}
	return nil
}

func AppendBytes(data1, data2 []byte) []byte {
	if data2 == nil {
		return data1
	}
	return append(data1, data2...)
}
