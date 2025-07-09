package ethereum

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	btcecEcdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okx/go-wallet-sdk/util"
	"golang.org/x/crypto/sha3"
)

const (
	DynamicFeeTxType = 0x02
)

type Eip1559Transaction struct {
	ChainId    *big.Int         `json:"chainId"`
	Nonce      uint64           `json:"nonce"`
	GasTipCap  *big.Int         `json:"gasTipCap"`
	GasFeeCap  *big.Int         `json:"gasFeeCap"`
	Gas        uint64           `json:"gas"`
	To         *common.Address  `json:"to" rlp:"nil"` // nil for contract creation
	Value      *big.Int         `json:"value"`
	Data       []byte           `json:"data"`
	AccessList types.AccessList `json:"accessList"`
}

type Eip1559TransactionVRS struct {
	ChainId    *big.Int         `json:"chainId"`
	Nonce      uint64           `json:"nonce"`
	GasTipCap  *big.Int         `json:"gasTipCap"`
	GasFeeCap  *big.Int         `json:"gasFeeCap"`
	Gas        uint64           `json:"gas"`
	To         *common.Address  `json:"to" rlp:"nil"` // nil for contract creation
	Value      *big.Int         `json:"value"`
	Data       []byte           `json:"data"`
	AccessList types.AccessList `json:"accessList"`
	V          *big.Int
	R          *big.Int
	S          *big.Int
}

func NewEip1559Transaction(
	chainId *big.Int,
	nonce uint64,
	maxPriorityFeePerGas *big.Int,
	maxFeePerGas *big.Int,
	gasLimit uint64,
	to *common.Address,
	value *big.Int,
	data []byte) *types.Transaction {
	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Nonce:     nonce,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Gas:       gasLimit,
		To:        to,
		Value:     value,
		Data:      data,
	})
}

func GenUnsignedEip1559Tx(tx *types.Transaction, chainId *big.Int) (string, error) {
	transaction := Eip1559Transaction{
		ChainId:    chainId,
		Nonce:      tx.Nonce(),
		GasTipCap:  tx.GasTipCap(),
		GasFeeCap:  tx.GasFeeCap(),
		Gas:        tx.Gas(),
		To:         tx.To(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	}
	txBytes, err := rlp.EncodeToBytes(transaction)
	if err != nil {
		return "", err
	}
	return util.EncodeHex(append([]byte{DynamicFeeTxType}, txBytes...)), nil
}

func GenerateEIP1559Tx(rlpStr string) (*types.Transaction, error) {
	tx := new(Eip1559Transaction)
	if len(rlpStr) > 1 {
		if rlpStr[0:2] == "02" {
			rlpStr = rlpStr[2:]
		}
	}
	unsignedByte, err := util.DecodeHexStringErr(rlpStr)
	if err != nil {
		return nil, err
	}
	err = rlp.DecodeBytes(unsignedByte, &tx)
	if err != nil {
		return nil, err
	}
	return NewEip1559Transaction(tx.ChainId, tx.Nonce, tx.GasTipCap, tx.GasFeeCap, tx.Gas, tx.To, tx.Value, tx.Data), nil
}

func SignMessageEIP1559(message []byte, prvKey *btcec.PrivateKey) (*SignatureData, error) {
	hash256 := sha3.NewLegacyKeccak256()
	hash256.Write([]byte{DynamicFeeTxType})
	hash256.Write(message)
	messageHash := hash256.Sum(nil)
	sig, err := btcecEcdsa.SignCompact(prvKey, messageHash, false)
	if err != nil {
		return nil, err
	}
	V := sig[0] - 27
	R := sig[1:33]
	S := sig[33:65]
	return &SignatureData{
		V:     new(big.Int).SetBytes([]byte{V}),
		R:     new(big.Int).SetBytes(R),
		S:     new(big.Int).SetBytes(S),
		ByteV: V,
		ByteR: R,
		ByteS: S,
	}, nil
}

func SignEip1559Tx(chainId *big.Int, tx *types.Transaction, prvKey *ecdsa.PrivateKey) ([]byte, error) {
	signer := types.NewLondonSigner(chainId)
	signedTx, err := types.SignTx(tx, signer, prvKey)
	if err != nil {
		return nil, err
	}
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return rawTx, nil
}

func GenEip1559TxWithSig(unsignedRawTx string, chainID, R, S, V *big.Int) (string, error) {
	unsignedRawTxByte, err := util.DecodeHexStringErr(unsignedRawTx)
	if err != nil {
		return "", err
	}
	var tx Eip1559Transaction
	err = rlp.DecodeBytes(unsignedRawTxByte[1:], &tx)
	if err != nil {
		return "", err
	}

	signedTx := Eip1559TransactionVRS{
		ChainId:    tx.ChainId,
		Nonce:      tx.Nonce,
		GasTipCap:  tx.GasTipCap,
		GasFeeCap:  tx.GasFeeCap,
		Gas:        tx.Gas,
		To:         tx.To,
		Value:      tx.Value,
		Data:       tx.Data,
		AccessList: tx.AccessList,
		V:          V,
		R:          R,
		S:          S,
	}

	signedTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", err
	}
	return util.EncodeHexWithPrefix(append([]byte{DynamicFeeTxType}, signedTxBytes...)), nil

}
