package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type UnsignedTx struct {
	Tx   string `json:"tx"`
	Hash string `json:"hash"`
}

type Eip1559Transaction struct {
	ChainId    *big.Int         `json:"chainId"`
	Nonce      uint64           `json:"nonce"`
	GasTipCap  *big.Int         `json:"gasTipCap"`
	GasFeeCap  *big.Int         `json:"gasFeeCap"`
	Gas        uint64           `json:"gas"`
	To         *common.Address  `json:"to"`
	Value      *big.Int         `json:"value"`
	Data       []byte           `json:"data"`
	AccessList types.AccessList `json:"accessList"`
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

type Eip1559Token struct {
	EIP1559Transaction
	ContractAddress string `json:"contract_address"`
	Amount          string `json:"amount"`
}

type Transaction struct {
	Nonce    string `json:"nonce"`
	GasLimit string `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	Fee      string `json:"fee"`
	ChainId  string `json:"chainId"`
}

type EIP1559Transaction struct {
	Transaction
	TxType               int    `json:"txType"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
}
