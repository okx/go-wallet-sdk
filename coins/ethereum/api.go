package ethereum

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okx/go-wallet-sdk/coins/ethereum/token"
	"github.com/okx/go-wallet-sdk/util"
	"golang.org/x/crypto/sha3"
	"math/big"
)

// generate tx with json param
func GenerateTxWithJSON(message string, chainId *big.Int, isToken bool) (*UnsignedTx, error) {
	var jsonTx Eip1559Token
	err := json.Unmarshal([]byte(message), &jsonTx)
	if err != nil {
		return nil, err
	}
	// read chainId, Use the incoming  chain id first
	if len(jsonTx.ChainId) > 0 {
		newChainId, ok := new(big.Int).SetString(jsonTx.ChainId, 10)
		if ok {
			chainId = newChainId
		}
	}
	// Generate transaction object
	// token logic
	var data []byte
	var toAddress common.Address
	if isToken {
		data, err = token.Transfer(jsonTx.To, util.ConvertToBigInt(jsonTx.Amount))
		if err != nil {
			return nil, err
		}
		toAddress = common.HexToAddress(jsonTx.ContractAddress)
	} else {
		data = util.RemoveZeroHex(jsonTx.Data)
		toAddress = common.HexToAddress(jsonTx.To)
	}
	if jsonTx.TxType == types.DynamicFeeTxType { // EIP1559 sign
		tx := NewEip1559Transaction(
			chainId,
			util.ConvertToUint64(jsonTx.Nonce),
			util.ConvertToBigInt(jsonTx.MaxPriorityFeePerGas),
			util.ConvertToBigInt(jsonTx.MaxFeePerGas),
			util.ConvertToUint64(jsonTx.GasLimit),
			&toAddress,
			util.ConvertToBigInt(jsonTx.Value),
			data,
		)
		res, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		hash := tx.Hash()
		return &UnsignedTx{Hash: hash.Hex(), Tx: util.EncodeHexWith0x(res)}, nil
	} else {
		// Token processing
		var tx *EthTransaction
		if isToken {
			tx = NewEthTransaction(util.ConvertToBigInt(jsonTx.Nonce), util.ConvertToBigInt(jsonTx.GasLimit), util.ConvertToBigInt(jsonTx.GasPrice), big.NewInt(0), jsonTx.ContractAddress, util.EncodeHexWith0x(data))
		} else {
			tx = NewEthTransaction(util.ConvertToBigInt(jsonTx.Nonce), util.ConvertToBigInt(jsonTx.GasLimit), util.ConvertToBigInt(jsonTx.GasPrice), util.ConvertToBigInt(jsonTx.Value), jsonTx.To, util.EncodeHexWith0x(data))
		}
		hash, res, err := tx.GetSigningHash(chainId)
		if err != nil {
			return nil, err
		}
		return &UnsignedTx{Tx: res, Hash: hash}, nil
	}
}

// GenerateRawTransactionWithSignature Generate the transaction to be broadcast based on the unsigned transaction and the signature result
func GenerateRawTransactionWithSignature(txType int, chainId, unsignedRawTx, r, s, v string) (string, error) {
	unsignedRawTxByte := util.RemoveZeroHex(unsignedRawTx)
	chainID, ok := new(big.Int).SetString(chainId, 10)
	if !ok {
		return "", ErrInvalidParam
	}
	R, ok := new(big.Int).SetString(r, 16)
	if !ok {
		return "", ErrInvalidParam
	}
	S, ok := new(big.Int).SetString(s, 16)
	if !ok {
		return "", ErrInvalidParam
	}
	V, ok := new(big.Int).SetString(v, 16)
	if !ok {
		return "", ErrInvalidParam
	}

	if txType == types.DynamicFeeTxType { // EIP1559 sign
		tx, err := generateEIP1559Tx(unsignedRawTx)
		if err != nil {
			return "", err
		}
		signer := types.NewLondonSigner(chainID)
		signedTx, err := tx.WithSignature(signer, encodeRSV(R, S, V))
		if err != nil {
			return "", err
		}
		rawTx, err := signedTx.MarshalBinary()
		if err != nil {
			return "", err
		}
		return util.EncodeHexWith0x(rawTx), err
	} else { // legacy sign
		var tx EthTransaction
		if err := rlp.DecodeBytes(unsignedRawTxByte, &tx); err != nil {
			return "", err
		}
		tx.V = V
		tx.R = R
		tx.S = S
		value, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return "", err
		}
		return util.EncodeHexWith0x(value), err
	}
}

func CalTxHash(rawTx string) string {
	bytes := util.RemoveZeroHex(rawTx)
	s256 := sha3.NewLegacyKeccak256()
	s256.Write(bytes)
	txBytes := s256.Sum(nil)
	return util.EncodeHexWith0x(txBytes)
}

func DecodeTx(rawTx string) (string, error) {
	rawTxBytes := util.RemoveZeroHex(rawTx)
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(rawTxBytes); err != nil {
		return "", err
	}
	txData := make(map[string]interface{})
	txData["nonce"] = big.NewInt(int64(tx.Nonce()))
	txData["gasLimit"] = float64(tx.Gas())
	txData["to"] = tx.To().String()
	txData["value"] = float64(tx.Value().Int64())
	txData["chainId"] = tx.ChainId()
	txData["txType"] = int(tx.Type())
	if tx.Type() == types.DynamicFeeTxType {
		txData["maxFeePerGas"] = float64(tx.GasFeeCap().Int64())
		txData["maxPriorityFeePerGas"] = float64(tx.GasTipCap().Int64())
	} else {
		txData["gasPrice"] = float64(tx.GasPrice().Int64())
	}
	signer := types.NewLondonSigner(tx.ChainId())
	if addr, err := types.Sender(signer, tx); err == nil {
		txData["from"] = addr.String()
	}
	v, r, s := tx.RawSignatureValues()
	txData["v"] = v
	txData["r"] = r
	txData["s"] = s
	if len(tx.Data()) != 0 {
		data := hex.EncodeToString(tx.Data())
		txData["inputData"] = "0x" + data
		inputData := make(map[string]interface{})
		if len(data) >= 72 {
			methodId := data[:8]
			approveMethod := "approve(address,uint256)"
			transferMethod := "transfer(address,uint256)"
			approveMethodId := hex.EncodeToString(crypto.Keccak256([]byte("approve(address,uint256)")))[:8]
			transferMethodId := hex.EncodeToString(crypto.Keccak256([]byte("transfer(address,uint256)")))[:8]
			address := "0x" + data[32:72]
			amount, ok := new(big.Int).SetString(data[72:], 16)
			if !ok {
				return "", ErrInvalidParam
			}
			if methodId == approveMethodId {
				inputData["method"] = approveMethod
				inputData["address"] = address
				inputData["amount"] = amount
			} else if methodId == transferMethodId {
				inputData["method"] = transferMethod
				inputData["address"] = address
				inputData["amount"] = amount
			}
		}
		txData["decodedData"] = inputData
	}
	b, err := json.Marshal(txData)
	return string(b), err
}

func GetAddress(pubkeyHex string) string {
	p, err := util.DecodeHexString(pubkeyHex)
	if err != nil {
		return ""
	}
	pubKey, err := btcec.ParsePubKey(p)
	if err != nil {
		return ""
	}
	return util.EncodeHexWith0x(getEthGroupPubHash(pubKey)[12:])
}

func ValidateAddress(address string) bool {
	if util.HasHexPrefix(address) {
		address = address[2:]
	}
	return len(address) == 2*AddressLength && util.IsHex(address)
}

func MessageHash(data string) string {
	return util.EncodeHexWith0x(calMessageHash(data))
}
