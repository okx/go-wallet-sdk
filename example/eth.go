package example

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	btcec2 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/okx/go-wallet-sdk/coins/ethereum"
	"github.com/okx/go-wallet-sdk/crypto/bip32"
	"github.com/okx/go-wallet-sdk/util"
	"github.com/tyler-smith/go-bip39"
	"math/big"
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
	childPrivateKey := hex.EncodeToString(c.Key.Key)
	return childPrivateKey, nil
}

func GetNewAddress(prvHex string) string {
	prvBytes, err := hex.DecodeString(prvHex)
	if err != nil {
		return ""
	}
	prv, pub := btcec.PrivKeyFromBytes(btcec.S256(), prvBytes)
	if prv == nil {
		return ""
	}
	return ethereum.GetAddress(hex.EncodeToString(pub.SerializeCompressed()))
}

func ValidAddress(address string) bool {
	return ethereum.ValidateAddress(address)
}

type SignParams struct {
	Type                 int    `json:"type"`
	ChainId              string `json:"chainId"`
	Nonce                string `json:"nonce"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	GasLimit             string `json:"gasLimit"`
	To                   string `json:"to"`
	Value                string `json:"value"`
	Data                 string `json:"data"`
	isToken              bool   `json:"isToken"`
}

func SignEip1559Transaction(chainId *big.Int, tx *types.Transaction, prvKey *ecdsa.PrivateKey) ([]byte, string, error) {
	signer := types.NewLondonSigner(chainId)
	signedTx, err := types.SignTx(tx, signer, prvKey)
	if err != nil {
		return nil, "", err
	}
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, "", err
	}
	return rawTx, signedTx.Hash().Hex(), nil
}

func toJosn(r interface{}) string {
	res, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(res)
}

type SignedTx struct {
	Hash string `json:"hash"`
	Hex  string `json:"hex"`
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
		if data, err = util.DecodeHexString(s.Data); err != nil {
			return "", err
		}
	}
	prvBytes, err := hex.DecodeString(prvHex)
	if err != nil {
		return "", errors.New("invalid prvHex")
	}
	var jsonTx ethereum.Eip1559Token
	if err := json.Unmarshal([]byte(txJson), &jsonTx); err != nil {
		return "", err
	}
	if jsonTx.TxType == types.DynamicFeeTxType { // EIP1559 sign
		prv, _ := btcec.PrivKeyFromBytes(btcec.S256(), prvBytes)
		tx := ethereum.NewEip1559Transaction(
			chainId,
			util.ConvertToBigInt(jsonTx.Nonce).Uint64(),
			util.ConvertToBigInt(jsonTx.MaxPriorityFeePerGas),
			util.ConvertToBigInt(jsonTx.MaxFeePerGas),
			util.ConvertToBigInt(jsonTx.GasLimit).Uint64(),
			to,
			util.ConvertToBigInt(jsonTx.Value),
			data,
		)
		res, hash, err := SignEip1559Transaction(chainId, tx, (*ecdsa.PrivateKey)(prv))
		if err != nil {
			return "", err
		}
		return toJosn(SignedTx{Hash: hash, Hex: util.EncodeHexWith0x(res)}), nil
	} else {
		prv, _ := btcec2.PrivKeyFromBytes(prvBytes)
		// Token processing
		var tx *ethereum.EthTransaction
		if s.isToken {
			tx = ethereum.NewEthTransaction(util.ConvertToBigInt(jsonTx.Nonce), util.ConvertToBigInt(jsonTx.GasLimit), util.ConvertToBigInt(jsonTx.GasPrice), big.NewInt(0), jsonTx.ContractAddress, util.EncodeHexWith0x(data))
		} else {
			tx = ethereum.NewEthTransaction(util.ConvertToBigInt(jsonTx.Nonce), util.ConvertToBigInt(jsonTx.GasLimit), util.ConvertToBigInt(jsonTx.GasPrice), util.ConvertToBigInt(jsonTx.Value), jsonTx.To, util.EncodeHexWith0x(data))
		}
		res, err := tx.SignTransaction(chainId, (*secp256k1.PrivateKey)(prv))
		if err != nil {
			return "", err
		}
		return toJosn(SignedTx{Hash: ethereum.CalTxHash(res), Hex: res}), nil
	}
}

func MessageHash(data string) string {
	return ethereum.MessageHash(data)
}

func GetEthereumMessagePrefix(message string) string {
	return ethereum.GetEthereumMessagePrefix(message)
}

func CalTxHash(rawTx string) string {
	return ethereum.CalTxHash(rawTx)
}

func GenerateRawTransactionWithSignature(txType int, chainId, unsignedRawTx, r, s, v string) (string, error) {
	return ethereum.GenerateRawTransactionWithSignature(txType, chainId, unsignedRawTx, r, s, v)
}
func DecodeTx(rawTx string) (string, error) {
	return ethereum.DecodeTx(rawTx)
}

func EcRecover(signature, message string, addPrefix bool) string {
	return ethereum.EcRecover(signature, message, addPrefix)
}
