package example

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okx/go-wallet-sdk/coins/ethereum"
	"github.com/okx/go-wallet-sdk/crypto/go-bip32"
	"github.com/okx/go-wallet-sdk/crypto/go-bip39"
	"github.com/okx/go-wallet-sdk/util"
)


func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	return mnemonic, err
}

func GetDerivedPath(index int) string {
	return fmt.Sprintf(`m/44'/60'/0'/0/%d`, index)
}

func GetDerivedPrivateKey(mnemonic string, hdPath string) (string, error) {
	seed := bip39.NewSeed(mnemonic, "")
	rp, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", err
	}
	c, err := rp.NewChildKeyByPathString(hdPath)
	if err != nil {
		return "", err
	}
	childPrivateKey := hex.EncodeToString(c.Key)
	return childPrivateKey, nil
}

func GetNewAddress(prvHex string) string {
	prvBytes, err := hex.DecodeString(prvHex)
	if err != nil {
		return ""
	}
	prv, pub := btcec.PrivKeyFromBytes(prvBytes)
	if prv == nil {
		return ""
	}
	return ethereum.GetNewAddress(pub)
}

func ValidAddress(address string) bool {
	return ethereum.ValidateAddress(address)
}

func SignTransaction(txJson, prvHex string) (string, error) {
	if len(txJson) == 0 {
		return "", errors.New("invalid txJson")
	}
	if len(prvHex) == 0 {
		return "", errors.New("invalid prvHex")
	}
	var err error
	var s SignParams
	if err := json.Unmarshal([]byte(txJson), &s); err != nil {
		return "", err
	}
	chainId := util.ConvertToBigInt(s.ChainId)
	var to *common.Address
	if len(s.To) != 0 {
		addr := common.HexToAddress(s.To)
		to = &addr
	}
	var data []byte
	if len(s.Data) != 0 {
		if data, err = util.DecodeHexStringErr(s.Data); err != nil {
			return "", err
		}
	}
	prvBytes, err := hex.DecodeString(prvHex)
	if err != nil {
		return "", errors.New("invalid prvHex")
	}
	prv, _ := btcec.PrivKeyFromBytes(prvBytes)

	var evmTx *ethereum.EVMTx

	if s.Type == ethereum.DynamicFeeTxType {
		tx := ethereum.NewEip1559Transaction(
			chainId,
			util.ConvertToBigInt(s.Nonce).Uint64(),
			util.ConvertToBigInt(s.MaxPriorityFeePerGas),
			util.ConvertToBigInt(s.MaxFeePerGas),
			util.ConvertToBigInt(s.GasLimit).Uint64(),
			to,
			util.ConvertToBigInt(s.Value),
			data,
		)
		evmTx = &ethereum.EVMTx{
			TxType:  s.Type,
			ChainId: chainId,
			Tx1559:  tx,
		}
	} else {
		tx := ethereum.NewEthTransaction(
			util.ConvertToBigInt(s.Nonce),
			util.ConvertToBigInt(s.GasLimit),
			util.ConvertToBigInt(s.GasPrice),
			util.ConvertToBigInt(s.Value),
			to.Hex(),
			util.EncodeHex(data),
		)
		evmTx = &ethereum.EVMTx{
			TxType:  s.Type,
			ChainId: chainId,
			Tx:      tx,
		}
	}

	return ethereum.SignTx(evmTx, prv)
}

func CalTxHash(rawTx string) (string, error) {
	hash, err := ethereum.CalTxHash(rawTx)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}
