package bitcoin

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/emresenyuva/go-wallet-sdk/crypto/base58"
	"github.com/emresenyuva/go-wallet-sdk/util"
)

func GetRedeemScript(pubKeys []string, minSignNum int) ([]byte, error) {
	var allPubKeys []*btcutil.AddressPubKey
	for _, v := range pubKeys {
		pubKey, err := hex.DecodeString(v)
		if err != nil {
			return nil, err
		}
		addressPubKey, err := btcutil.NewAddressPubKey(pubKey, &chaincfg.MainNetParams)

		if err != nil {
			return nil, err
		}
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
