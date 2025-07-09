package ethereum

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okx/go-wallet-sdk/util"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type AccessList []AccessTuple

type Hash [32]byte

// AccessTuple is the element type of an access list.
type AccessTuple struct {
	Address     []byte `json:"address"     gencodec:"required"`
	StorageKeys []Hash `json:"storageKeys" gencodec:"required"`
}

// DynamicFeeTx represents an EIP-1559 transaction.
type DynamicFeeTx struct {
	ChainID    *big.Int `json:"chainId"`
	Nonce      uint64   `json:"nonce"`
	GasTipCap  *big.Int `json:"gasTipCap"` // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int `json:"gasFeeCap"` // a.k.a. maxFeePerGas
	Gas        uint64   `json:"gas"`
	To         *Address `json:"to"` // nil means contract creation
	Value      *big.Int `json:"value"`
	Data       []byte   `json:"data"`
	AccessList AccessList

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

func NewEthDynamicFeeTx(chainId *big.Int, nonce uint64, gasTipCap, gasFeeCap *big.Int, gas uint64, value *big.Int, to,
	data string, accessList AccessList) *DynamicFeeTx {
	toBytes := HexToAddress(to)
	dataBytes := util.DecodeHexStringPad(data)
	return &DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      nonce,
		GasTipCap:  gasTipCap,
		GasFeeCap:  gasFeeCap,
		Gas:        gas,
		To:         &toBytes,
		Value:      value,
		Data:       dataBytes,
		AccessList: accessList,
	}
}

func (tx *DynamicFeeTx) SignTransaction(prvKey *btcec.PrivateKey) (string, error) {
	rawTransaction, _ := rlp.EncodeToBytes([]interface{}{
		tx.ChainID,
		tx.Nonce,
		tx.GasTipCap,
		tx.GasFeeCap,
		tx.Gas,
		tx.To,
		tx.Value,
		tx.Data,
		tx.AccessList,
	})
	sig, err := SignMessageEIP1559(rawTransaction, prvKey)
	if err != nil {
		return "", err
	}
	tx.V = sig.V
	tx.R = sig.R
	tx.S = sig.S
	var buf bytes.Buffer
	buf.Write([]byte{DynamicFeeTxType})
	err = rlp.Encode(&buf, tx)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(buf.Bytes()), nil
}

func (tx *DynamicFeeTx) Hash() (string, error) {
	hash256 := sha3.NewLegacyKeccak256()
	hash256.Write([]byte{DynamicFeeTxType})
	err := rlp.Encode(hash256, tx)
	if err != nil {
		return "", err
	}
	messageHash := hash256.Sum(nil)
	return hex.EncodeToString(messageHash), nil
}

// accessors for innerTx.
func (tx *DynamicFeeTx) txType() byte           { return DynamicFeeTxType }
func (tx *DynamicFeeTx) chainID() *big.Int      { return tx.ChainID }
func (tx *DynamicFeeTx) accessList() AccessList { return tx.AccessList }
func (tx *DynamicFeeTx) data() []byte           { return tx.Data }
func (tx *DynamicFeeTx) gas() uint64            { return tx.Gas }
func (tx *DynamicFeeTx) gasFeeCap() *big.Int    { return tx.GasFeeCap }
func (tx *DynamicFeeTx) gasTipCap() *big.Int    { return tx.GasTipCap }
func (tx *DynamicFeeTx) gasPrice() *big.Int     { return tx.GasFeeCap }
func (tx *DynamicFeeTx) value() *big.Int        { return tx.Value }
func (tx *DynamicFeeTx) nonce() uint64          { return tx.Nonce }
func (tx *DynamicFeeTx) to() *Address           { return tx.To }

func (tx *DynamicFeeTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	if baseFee == nil {
		return dst.Set(tx.GasFeeCap)
	}
	tip := dst.Sub(tx.GasFeeCap, baseFee)
	if tip.Cmp(tx.GasTipCap) > 0 {
		tip.Set(tx.GasTipCap)
	}
	return tip.Add(tip, baseFee)
}

func (tx *DynamicFeeTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *DynamicFeeTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}

func (tx *DynamicFeeTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *DynamicFeeTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
