package bitcoin

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/okx/go-wallet-sdk/util"
)

func GetRedeemScript(pubKeys []string, minSignNum int) ([]byte, error) {
	var allPubKeys []*btcutil.AddressPubKey
	for _, v := range pubKeys {
		pubKey, _ := hex.DecodeString(v)
		addressPubKey, _ := btcutil.NewAddressPubKey(pubKey, &chaincfg.MainNetParams)
		allPubKeys = append(allPubKeys, addressPubKey)
	}
	return txscript.MultiSigScript(allPubKeys, minSignNum)
}

func GenerateMultiAddress(redeemScript []byte, net *chaincfg.Params) (string, error) {
	if net == nil {
		net = &chaincfg.MainNetParams
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(redeemScript, net)
	if err != nil {
		return "", err
	}
	P2SHAddress := base58.CheckEncode(addressScriptHash.ScriptAddress(), net.ScriptHashAddrID)
	return P2SHAddress, nil
}

func GenerateAddress(pubKey string, net *chaincfg.Params) (string, error) {
	if net == nil {
		net = &chaincfg.MainNetParams
	}
	addressPubKey, err := btcutil.NewAddressPubKey(util.RemoveZeroHex(pubKey), &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return addressPubKey.EncodeAddress(), nil
}
