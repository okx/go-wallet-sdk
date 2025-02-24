package harmony

import (
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/okx/go-wallet-sdk/coins/ethereum"
	"github.com/okx/go-wallet-sdk/util"
	"math/big"
	"strings"
)

var (
	ErrInvalidSign   = errors.New("invali sign")
	ErrInvalidPubKey = errors.New("invalid public key")
)

const HRP = "one"

func NewAddress(seedHex string, followETH bool) (string, error) {
	p, err := hex.DecodeString(seedHex)
	if err != nil {
		return "", err
	}
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	addr := ethereum.GetNewAddress(prvKey.PubKey())
	if followETH {
		return addr, nil
	}
	bytes, err := hex.DecodeString(addr[2:])
	if err != nil {
		return "", err
	}
	bech32Address, err := bech32.EncodeFromBase256(HRP, bytes)
	if err != nil {
		return "", err
	}
	return bech32Address, nil
}
func GetAddress(pub *btcec.PublicKey) (string, error) {
	if pub == nil {
		return "", ErrInvalidPubKey
	}
	ethAddressBytes := ethereum.GetNewAddressBytes(pub)
	bech32Address, err := bech32.EncodeFromBase256(HRP, ethAddressBytes)
	if err != nil {
		return "", err
	}
	return bech32Address, nil
}

func Transfer(transaction *ethereum.EthTransaction, chainId *big.Int, prvKey *btcec.PrivateKey) (string, error) {
	return transaction.SignTransaction(chainId, prvKey)
}
func VerifySignMsg(signature, message, address string, addPrefix bool) error {
	pub, err := ethereum.EcRecoverPubKey(signature, message, addPrefix)
	if err != nil {
		return err
	}
	addr, err := GetAddress(pub)
	if err != nil {
		return err
	}
	if addr == address {
		return nil
	}
	ethAddress := ethereum.GetNewAddress(pub)
	if ethAddress == address {
		return nil
	}
	return ErrInvalidSign
}

func oneAddressToEthAddress(address string) string {
	if strings.HasPrefix(strings.ToLower(address), HRP) {
		_, hexByte, err := bech32.DecodeToBase256(address)
		if err != nil {
			return ""
		}
		return hex.EncodeToString(hexByte)
	}
	return hex.EncodeToString(util.RemoveZeroHex(address))
}
func ValidateAddress(address string) bool {
	if strings.HasPrefix(address, "0x") {
		return ethereum.IsEthHexAddress(address)
	}
	return ethereum.IsEthHexAddress(oneAddressToEthAddress(address))
}
