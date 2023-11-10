package zkspace

import (
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/zksync/core"
)

func GetPubKeyHash(ethPrivKeyHex string, chainId int) (string, error) {
	l1PrivateKeyBytes, err := hex.DecodeString(ethPrivKeyHex)
	if err != nil {
		return "", err
	}
	ethSigner, err := core.NewOkEthSignerFromPrivBytes(l1PrivateKeyBytes)
	if err != nil {
		return "", err
	}
	zkSigner, err := NewZkSignerFromEthSigner(ethSigner, core.ChainId(chainId))
	return zkSigner.GetPublicKeyHash(), nil
}

func GetAddress(ethPrivKeyHex string) (string, error) {
	privKeyBytes, err := hex.DecodeString(ethPrivKeyHex)
	if err != nil {
		return "", err
	}
	ethSigner, err := core.NewOkEthSignerFromPrivBytes(privKeyBytes)
	if err != nil {
		return "", err
	}
	return ethSigner.GetAddress(), nil
}
