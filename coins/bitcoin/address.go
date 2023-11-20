package bitcoin

import (
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

const (
	LEGACY        = "legacy"
	SEGWIT_NATIVE = "segwit_native"
	SEGWIT_NESTED = "segwit_nested"
	TAPROOT       = "taproot"
)

func PubKeyToAddr(publicKey []byte, addrType string, network *chaincfg.Params) (string, error) {
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	if addrType == LEGACY {
		p2pkh, err := btcutil.NewAddressPubKey(publicKey, network)
		if err != nil {
			return "", err
		}

		return p2pkh.EncodeAddress(), nil
	} else if addrType == SEGWIT_NATIVE {
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), network)
		if err != nil {
			return "", err
		}

		return p2wpkh.EncodeAddress(), nil
	} else if addrType == SEGWIT_NESTED {
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), network)
		if err != nil {
			return "", err
		}
		redeemScript, err := txscript.PayToAddrScript(p2wpkh)
		if err != nil {
			return "", err
		}
		p2sh, err := btcutil.NewAddressScriptHash(redeemScript, network)
		if err != nil {
			return "", err
		}

		return p2sh.EncodeAddress(), nil
	} else if addrType == TAPROOT {
		internalKey, err := btcec.ParsePubKey(publicKey)
		if err != nil {
			return "", err
		}
		p2tr, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(internalKey).SerializeCompressed()[1:], network)
		if err != nil {
			return "", err
		}

		return p2tr.EncodeAddress(), nil
	} else {
		return "", errors.New("address type not supported")
	}
}
