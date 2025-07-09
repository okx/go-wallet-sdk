package ethereum

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okx/go-wallet-sdk/util"
)

const (
	AuthorizationTxType      = 0x04
	AuthorizationTxTypeMagic = 0x05
)

type EthAuthorization struct {
	ChainId *big.Int `json:"chainId"`
	Address []byte   `json:"address"`
	Nonce   *big.Int `json:"nonce"`
	YParity *big.Int `json:"yParit"`
	R       *big.Int `json:"r"`
	S       *big.Int `json:"s"`
}

type Eip7702Transaction struct {
	ChainId           *big.Int            `json:"chainId"`
	Nonce             uint64              `json:"nonce"`
	GasTipCap         *big.Int            `json:"gasTipCap"`
	GasFeeCap         *big.Int            `json:"gasFeeCap"`
	Gas               uint64              `json:"gas"`
	To                *common.Address     `json:"to" rlp:"nil"`
	Value             *big.Int            `json:"value"`
	Data              []byte              `json:"data"`
	AccessList        types.AccessList    `json:"accessList"`
	AuthorizationList []*EthAuthorization `json:"authorizationList"`
}

type Eip7702TransactionVRS struct {
	ChainId           *big.Int            `json:"chainId"`
	Nonce             uint64              `json:"nonce"`
	GasTipCap         *big.Int            `json:"gasTipCap"`
	GasFeeCap         *big.Int            `json:"gasFeeCap"`
	Gas               uint64              `json:"gas"`
	To                *common.Address     `json:"to" rlp:"nil"`
	Value             *big.Int            `json:"value"`
	Data              []byte              `json:"data"`
	AccessList        types.AccessList    `json:"accessList"`
	AuthorizationList []*EthAuthorization `json:"authorizationList"`
	V                 *big.Int
	R                 *big.Int
	S                 *big.Int
}

func GenUnsignedEip7702Tx(tx *types.Transaction, authList []*EthAuthorization, chainId *big.Int) (string, error) {
	err := CheckAuthList(authList)
	if err != nil {
		return "", err
	}
	transaction := Eip7702Transaction{
		ChainId:           chainId,
		Nonce:             tx.Nonce(),
		GasTipCap:         tx.GasTipCap(),
		GasFeeCap:         tx.GasFeeCap(),
		Gas:               tx.Gas(),
		To:                tx.To(),
		Value:             tx.Value(),
		Data:              tx.Data(),
		AccessList:        tx.AccessList(),
		AuthorizationList: authList,
	}
	baseRawTransaction, err := rlp.EncodeToBytes(transaction)
	if err != nil {
		return "", err
	}
	return util.EncodeHex(appendAuthTxType(baseRawTransaction)), nil
}

func SignEip7702Tx(tx *types.Transaction, authList []*EthAuthorization, chainId *big.Int, prvKey *btcec.PrivateKey) ([]byte, error) {
	err := CheckAuthList(authList)
	if err != nil {
		return nil, err
	}
	transaction := Eip7702Transaction{
		ChainId:           chainId,
		Nonce:             tx.Nonce(),
		GasTipCap:         tx.GasTipCap(),
		GasFeeCap:         tx.GasFeeCap(),
		Gas:               tx.Gas(),
		To:                tx.To(),
		Value:             tx.Value(),
		Data:              tx.Data(),
		AccessList:        tx.AccessList(),
		AuthorizationList: authList,
	}
	baseRawTransaction, err := rlp.EncodeToBytes(transaction)
	if err != nil {
		return nil, err
	}
	rawTransaction := appendAuthTxType(baseRawTransaction)

	sig := SignMessage(rawTransaction, prvKey)

	extendedTx := Eip7702TransactionVRS{
		ChainId:           chainId,
		Nonce:             tx.Nonce(),
		GasTipCap:         tx.GasTipCap(),
		GasFeeCap:         tx.GasFeeCap(),
		Gas:               tx.Gas(),
		To:                tx.To(),
		Value:             tx.Value(),
		Data:              tx.Data(),
		AccessList:        tx.AccessList(),
		AuthorizationList: authList,
		V:                 big.NewInt(sig.V.Int64() - 27),
		R:                 sig.R,
		S:                 sig.S,
	}

	signed, err := rlp.EncodeToBytes(extendedTx)
	if err != nil {
		return nil, err
	}
	return appendAuthTxType(signed), err
}

func NewEthAuthorization(nonce, chainid, yParity, r, s *big.Int, address []byte) *EthAuthorization {
	return &EthAuthorization{
		Nonce:   nonce,
		Address: address,
		ChainId: chainid,
		YParity: yParity,
		R:       r,
		S:       s,
	}
}

func SignAuthorization(tx EthAuthorization, prvKey *btcec.PrivateKey) (EthAuthorization, error) {
	rawTransaction, _ := rlp.EncodeToBytes([]interface{}{
		tx.ChainId,
		tx.Address,
		tx.Nonce,
	})
	addMagic := append([]byte{AuthorizationTxTypeMagic}, rawTransaction...)
	msgHash := crypto.Keccak256(addMagic)

	sig := SignAsRecoverable(msgHash, prvKey)
	tx.YParity = big.NewInt(sig.V.Int64() - 27)
	tx.R = sig.R
	tx.S = sig.S
	return tx, nil
}

func GenerateSignAuthHash(address []byte, nonce *big.Int, chainId *big.Int) ([]byte, error) {
	rawTransaction, _ := rlp.EncodeToBytes([]interface{}{
		chainId,
		address,
		nonce,
	})
	return CalTxHash(util.EncodeHex(append([]byte{AuthorizationTxTypeMagic}, rawTransaction...)))
}

func CheckAuthList(authList []*EthAuthorization) error {
	for i, auth := range authList {
		if auth == nil {
			return fmt.Errorf("authorization at index %d is nil", i)
		}
		if auth.Nonce == nil {
			return fmt.Errorf("missing Nonce at index %d", i)
		}
		if len(auth.Address) == 0 {
			return fmt.Errorf("missing Address at index %d", i)
		}
		if auth.ChainId == nil {
			return fmt.Errorf("missing ChainId at index %d", i)
		}
		if auth.YParity == nil {
			return fmt.Errorf("missing YParity at index %d", i)
		}
		if auth.R == nil || auth.R.Sign() == 0 {
			return fmt.Errorf("missing R at index %d", i)
		}
		if auth.S == nil || auth.S.Sign() == 0 {
			return fmt.Errorf("missing S at index %d", i)
		}
	}
	return nil
}

func EcRecoverAuthorization(tx EthAuthorization) (string, error) {
	rawTransaction, _ := rlp.EncodeToBytes([]interface{}{
		tx.ChainId,
		tx.Address,
		tx.Nonce,
	})
	addMagic := append([]byte{AuthorizationTxTypeMagic}, rawTransaction...)
	msgHash := crypto.Keccak256(addMagic)
	//sig := SignAsRecoverable(msgHash, prvKey)
	buffR := make([]byte, 32)
	buffS := make([]byte, 32)

	curRBytes := tx.R.FillBytes(buffR)
	curSBytes := tx.S.FillBytes(buffS)
	curVBytes := byte(tx.YParity.Int64())

	sigByte := append(curRBytes, curSBytes...)
	sigByte = append(sigByte, curVBytes)

	pubByte, err := crypto.Ecrecover(msgHash, sigByte)
	if err != nil {
		return "", err
	}
	addr := common.BytesToAddress(crypto.Keccak256(pubByte[1:])[12:])

	return addr.String(), nil
}

func appendAuthTxType(rawTransaction []byte) []byte {
	if len(rawTransaction) == 0 {
		return nil
	}
	if rawTransaction[0] == AuthorizationTxType {
		return rawTransaction
	}
	return append([]byte{AuthorizationTxType}, rawTransaction...)
}

func GenEip7702TxWithSig(unsignedRawTx string, chainID, R, S, V *big.Int) (string, error) {
	unsignedRawTxByte, err := util.DecodeHexStringErr(unsignedRawTx)
	if err != nil {
		return "", err
	}
	var tx Eip7702Transaction
	err = rlp.DecodeBytes(unsignedRawTxByte[1:], &tx)
	if err != nil {
		return "", err
	}

	signedTx := Eip7702TransactionVRS{
		ChainId:           tx.ChainId,
		Nonce:             tx.Nonce,
		GasTipCap:         tx.GasTipCap,
		GasFeeCap:         tx.GasFeeCap,
		Gas:               tx.Gas,
		To:                tx.To,
		Value:             tx.Value,
		Data:              tx.Data,
		AccessList:        tx.AccessList,
		AuthorizationList: tx.AuthorizationList,
		V:                 V,
		R:                 R,
		S:                 S,
	}

	signedTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", err
	}
	return util.EncodeHex(append([]byte{AuthorizationTxType}, signedTxBytes...)), nil
}
