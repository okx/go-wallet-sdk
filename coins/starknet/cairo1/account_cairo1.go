package cairo1

import (
	"github.com/okx/go-wallet-sdk/coins/starknet"
	"math/big"
)

func CalculateContractAddressFromHash(starkPub string, accountClass string) (hash *big.Int, err error) {
	salt, err := starknet.HexToBN(starkPub)
	if err != nil {
		return nil, err
	}
	classHash, err := starknet.HexToBN(accountClass)
	if err != nil {
		return nil, err
	}

	deployerAddress := big.NewInt(0)

	calldate := []*big.Int{salt, big.NewInt(0)}

	constructorCallData := []*big.Int{}
	constructorCallData = append(constructorCallData, calldate...)

	constructorCalldataHash, err := starknet.ComputeHashOnElements(constructorCallData)
	if err != nil {
		return nil, err
	}
	ContractAddressPrefix, err := starknet.HexToBN("0x535441524b4e45545f434f4e54524143545f41444452455353")
	if err != nil {
		return nil, err
	}

	ele := []*big.Int{
		ContractAddressPrefix,
		deployerAddress,
		salt,
		classHash,
		constructorCalldataHash,
	}
	return starknet.ComputeHashOnElements(ele)
}

func CreateDeployAccountTx(starkPub string, accountClass string, nonce, maxFee *big.Int, chainId string) (*starknet.DeployTransaction, error) {
	pub, err := starknet.HexToBN(starkPub)
	if err != nil {
		return nil, err
	}
	version := big.NewInt(starknet.TRANSACTION_VERSION)
	classHash, err := starknet.HexToBN(accountClass)
	if err != nil {
		return nil, err
	}

	calldate := []*big.Int{pub, big.NewInt(0)}

	constructorCallData := []*big.Int{}
	constructorCallData = append(constructorCallData, calldate...)

	// calculate address
	contractAddress, err := CalculateContractAddressFromHash(starkPub, accountClass)
	if err != nil {
		return nil, err
	}

	// calculate tx hash
	txHash, err := starknet.CalculateDeployAccountTransactionHashCairo1(contractAddress, classHash, constructorCallData, pub, version, starknet.UTF8StrToBig(chainId), nonce, maxFee)
	if err != nil {
		return nil, err
	}

	return &starknet.DeployTransaction{
		Type:                starknet.TransactionTypeDeployAccount,
		ContractAddressSalt: pub,
		ConstructorCalldata: constructorCallData,
		ClassHash:           classHash,
		MaxFee:              maxFee,
		Version:             version,
		Nonce:               nonce,
		Signature:           nil,
		TransactionHash:     txHash,
	}, nil
}

func CreateSignedDeployAccountTx(curve starknet.StarkCurve, starkPub string, accountClass string, nonce, maxFee *big.Int, chainId string, privateKey string) (*starknet.DeployTransaction, error) {
	tx, err := CreateDeployAccountTx(starkPub, accountClass, nonce, maxFee, chainId)
	if err != nil {
		return nil, err
	}

	if err = starknet.SignTx(curve, tx, privateKey); err != nil {
		return nil, err
	}
	return tx, nil
}
