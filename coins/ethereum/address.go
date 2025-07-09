package ethereum

import (
	"github.com/okx/go-wallet-sdk/util"

	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/sha3"
)

const (
	AddressLength = 20
)

func GetNewAddress(pubKey *btcec.PublicKey) string {
	return util.EncodeHexWithPrefix(GetNewAddressBytes(pubKey))
}

func GetNewAddressBytes(pubKey *btcec.PublicKey) []byte {
	pubBytes := pubKey.SerializeUncompressed()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:])
	addressByte := hash.Sum(nil)
	return addressByte[12:]
}

func PubKeyToAddr(publicKey []byte) (string, error) {
	pubKey, err := btcec.ParsePubKey(publicKey)
	if err != nil {
		return "", err
	}
	return GetNewAddress(pubKey), nil
}

func ValidateAddress(s string) bool {
	return IsEthHexAddress(s)
}

func IsEthHexAddress(s string) bool {
	if HasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && IsHex(s)
}

func HasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func IsHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}
