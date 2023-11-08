package kaspa

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/kaspanet/kaspad/domain/dagconfig"
	"github.com/kaspanet/kaspad/util"
)

func NewAddress(prvKeyHex string) (string, error) {
	return NewAddressWithNetParams(prvKeyHex, dagconfig.MainnetParams)
}

func NewAddressWithNetParams(prvKeyHex string, params dagconfig.Params) (string, error) {
	prvKeyBytes, err := hex.DecodeString(prvKeyHex)
	if err != nil {
		return "", err
	}
	_, pubKey := btcec.PrivKeyFromBytes(prvKeyBytes)
	pubKeyAddress, err := util.NewAddressPublicKey(schnorr.SerializePubKey(pubKey), params.Prefix)
	if err != nil {
		return "", err
	}

	return pubKeyAddress.EncodeAddress(), nil
}

func ValidateAddress(address string) bool {
	return ValidateAddressWithNetParams(address, dagconfig.MainnetParams)
}

func ValidateAddressWithNetParams(address string, params dagconfig.Params) bool {
	_, err := util.DecodeAddress(address, params.Prefix)
	return err == nil
}

func PrvKeyToPubKey(prvKeyHex string) (string, error) {
	prvKeyBytes, err := hex.DecodeString(prvKeyHex)
	if err != nil {
		return "", err
	}

	_, pubKey := btcec.PrivKeyFromBytes(prvKeyBytes)

	return hex.EncodeToString(schnorr.SerializePubKey(pubKey)), nil
}
