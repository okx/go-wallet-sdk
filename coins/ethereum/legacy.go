package ethereum

import (
	"github.com/okx/go-wallet-sdk/util"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

type EthTransaction struct {
	Nonce    *big.Int `json:"nonce"`
	GasPrice *big.Int `json:"gasPrice"`
	GasLimit *big.Int `json:"gas"`
	To       []byte   `json:"to"`
	Value    *big.Int `json:"value"`
	Data     []byte   `json:"data"`
	// Signature values
	V *big.Int `json:"v"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

func NewEthTransaction(nonce, gasLimit, gasPrice, value *big.Int, to, data string) *EthTransaction {
	toBytes := util.DecodeHexStringPad(to)
	dataBytes := util.DecodeHexStringPad(data)
	return &EthTransaction{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		To:       toBytes,
		Value:    value,
		Data:     dataBytes,
	}
}

func NewTransactionFromRaw(raw string) (*EthTransaction, error) {
	bytes := util.DecodeHexStringPad(raw)
	t := new(EthTransaction)
	err := rlp.DecodeBytes(bytes, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (tx *EthTransaction) UnSignedTx(chainId *big.Int) string {
	tx.V = chainId
	rawTransaction, _ := rlp.EncodeToBytes([]interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
		chainId, uint(0), uint(0),
	})
	return util.EncodeHex(rawTransaction)
}

func (tx *EthTransaction) GetSigningHash(chainId *big.Int) (string, string, error) {
	unSignedTx := tx.UnSignedTx(chainId)
	raw, err := util.DecodeHexStringErr(unSignedTx)
	if err != nil {
		return "", "", err
	}
	h := sha3.NewLegacyKeccak256()
	h.Write(raw)
	msgHash := h.Sum(nil)
	return util.EncodeHex(msgHash), unSignedTx, nil
}

func (tx *EthTransaction) Hash() (string, error) {
	sha := sha3.NewLegacyKeccak256()
	sha.Reset()
	value, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return "", err
	}
	_, err = sha.Write(value)
	if err != nil {
		return "", err
	}

	hash := sha.Sum(nil)
	return util.EncodeHex(hash[:]), nil
}

func (tx *EthTransaction) GenUnsignedTx(chainId *big.Int) string {
	tx.V = chainId
	rawData, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return ""
	}
	return util.EncodeHex(rawData)
}

func (tx *EthTransaction) SignTransaction(chainId *big.Int, prvKey *btcec.PrivateKey) string {
	tx.V = chainId
	rawTransaction, err := rlp.EncodeToBytes([]interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
		chainId, uint(0), uint(0),
	})
	if err != nil {
		return ""
	}
	sig := SignMessage(rawTransaction, prvKey)
	tx.V = big.NewInt(chainId.Int64()*2 + sig.V.Int64() + 8)
	tx.R = sig.R
	tx.S = sig.S
	value, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return ""
	}
	return util.EncodeHexWithPrefix(value)
}

func (tx *EthTransaction) SignedTx(chainId *big.Int, sig *SignatureData) string {
	tx.V = big.NewInt(chainId.Int64()*2 + sig.V.Int64() + 8)
	tx.R = sig.R
	tx.S = sig.S
	value, _ := rlp.EncodeToBytes(tx)
	return util.EncodeHexWithPrefix(value)
}
func SignLegacyTx(tx *EthTransaction, chainId *big.Int, prvKey *btcec.PrivateKey) ([]byte, error) {
	rawTransaction, _ := rlp.EncodeToBytes([]interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
		chainId, uint(0), uint(0),
	})
	sig := SignMessage(rawTransaction, prvKey)
	tx.V = big.NewInt(chainId.Int64()*2 + sig.V.Int64() + 8)
	tx.R = sig.R
	tx.S = sig.S
	return rlp.EncodeToBytes(tx)
}

func GenLegacyTxWithSig(unsignedRawTx string, chainID, R, S, V *big.Int) (string, error) {
	unsignedRawTxByte, err := util.DecodeHexStringPadErr(unsignedRawTx)
	if err != nil {
		return "", err
	}
	var tx EthTransaction
	err = rlp.DecodeBytes(unsignedRawTxByte, &tx)
	if err != nil {
		return "", err
	}
	tx.V = V
	if tx.V.Int64() == 0 || tx.V.Int64() == 1 {
		tx.V = tx.V.Add(tx.V, chainID.Mul(chainID, big.NewInt(2))).Add(tx.V, big.NewInt(35))
	}
	tx.R = R
	tx.S = S
	rawTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return "", err
	}
	return util.EncodeHexWithPrefix(rawTx), nil
}
