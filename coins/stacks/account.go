package stacks

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func ValidAddress(address string) bool {
	if _, err := c32addressDecode(address); err != nil {
		return false
	}
	return true
}

func GetPublicKey(privKeyHex string) (string, error) {
	privKey, err := createStacksPrivateKey(privKeyHex)
	if err != nil {
		return "", err
	}
	privateKey := secp256k1.PrivKeyFromBytes(privKey.Data)
	if privKey.Compressed {
		pubBytes := privateKey.PubKey().SerializeCompressed()
		return hex.EncodeToString(pubBytes), nil
	}
	pubBytes := privateKey.PubKey().SerializeUncompressed()
	return hex.EncodeToString(pubBytes), nil
}

func GetAddressFromPublicKey(pubKey string) (string, error) {
	hexStr := hashP2PKH(pubKey)
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	version := 22
	versionHex := fmt.Sprintf("%02x", version)
	checkSumHex, err := c32checksum(versionHex + hexStr)
	if err != nil {
		return "", err
	}
	encode := c32encode(hexStr + checkSumHex)
	p := c32[version : version+1]
	return "S" + p + encode, nil
}

func NewAddress(privKeyHex string) (string, error) {
	public, err := GetPublicKey(privKeyHex)
	if err != nil {
		return "", err
	}
	return GetAddressFromPublicKey(public)
}

func GetPoxAddress(poxAddress string) (*TupleCV, error) {
	bean, err := decodeBtcAddress(poxAddress)
	if err != nil {
		return nil, err
	}

	version := make([]byte, 0)
	version = append(version, byte(bean.hashMode))
	tupleData := make(map[string]ClarityValue)
	tupleData["hashbytes"] = &BufferCV{2, bean.data}
	tupleData["version"] = &BufferCV{2, version}

	return &TupleCV{
		Data: tupleData,
		Type: Tuple,
	}, nil
}

func createStacksPrivateKey(secretKey string) (*StacksPrivateKey, error) {
	data, err := hex.DecodeString(secretKey)
	base64.StdEncoding.EncodeToString(data[:])
	if err != nil {
		return nil, err
	}
	var compressed bool
	if len(data) == 33 {
		if data[len(data)-1] != 1 {
			return nil, errors.New("improperly formatted private-key. 33 byte length usually indicates compressed key, but last byte must be == 0x01")
		}
		compressed = true
	} else if len(data) == 32 {
		compressed = false
	} else {
		return nil, fmt.Errorf("improperly formatted private-key hex string: length should be 32 or 33 bytes, provided with length %d", len(data))
	}
	spk := &StacksPrivateKey{
		Data:       data,
		Compressed: compressed,
	}
	return spk, nil
}

func addressFromPublicKeys(version uint64, hashMode uint64, numSigs int, pubKeys []StacksPublicKey) (*Signer, error) {
	if len(pubKeys) == 0 {
		return nil, fmt.Errorf("invalid number of public keys")
	}
	if hashMode == 0 || hashMode == 2 {
		if len(pubKeys) != 1 || numSigs != 1 {
			return nil, fmt.Errorf("invalid number of public keys or signatures")
		}
	}

	if hashMode == 2 || hashMode == 3 {
		for _, pubKey := range pubKeys {
			if !isCompressed(pubKey) {
				return nil, fmt.Errorf("public keys must be compressed for segwit")
			}
		}
	}

	switch hashMode {
	case 0:
		hash := hashP2PKH(pubKeys[0].Data)
		return addressFromVersionHash(version, hash), nil
	case 1:
		return nil, fmt.Errorf("not yet implemented: address construction using public keys for hash mode: %v", hashMode)
	default:
		return nil, fmt.Errorf("not yet implemented: address construction using public keys for hash mode: %v", hashMode)
	}
}

func pubKeyFromPrivKey(privateKey string) (*StacksPublicKey, error) {
	privKey, err := createStacksPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privateKeyHex := hex.EncodeToString(privKey.Data)
	pubKeyHex, err := GetPublicKey(privateKeyHex)
	if err != nil {
		return nil, err
	}
	stacksPublicKey := &StacksPublicKey{}
	stacksPublicKey.Type_ = 6
	stacksPublicKey.Data = pubKeyHex
	return stacksPublicKey, nil
}
