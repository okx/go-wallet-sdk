package v3

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/okx/go-wallet-sdk/coins/starknet"
	"github.com/okx/go-wallet-sdk/coins/starknet/cairo1"
	"github.com/okx/go-wallet-sdk/coins/starknet/juno_core/felt"
)

var (
	PREFIX_TRANSACTION      = new(felt.Felt).SetBytes([]byte("invoke"))
	PREFIX_DECLARE          = new(felt.Felt).SetBytes([]byte("declare"))
	PREFIX_CONTRACT_ADDRESS = new(felt.Felt).SetBytes([]byte("STARKNET_CONTRACT_ADDRESS"))
	PREFIX_DEPLOY_ACCOUNT   = new(felt.Felt).SetBytes([]byte("deploy_account"))
	CairoCalldataLength     = 2
)

func CreateDeployAccountTxV3(starkPub string, accountClass, proxyAccountClass string, nonce *big.Int, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit string, chainId string) (*DeployAccountTxnV3, error) {
	pub, err := starknet.HexToBN(starkPub)
	if err != nil {
		return nil, err
	}

	var classHash *big.Int
	var constructorCallData []*big.Int
	var contractAddress *big.Int

	// cairo0
	if proxyAccountClass != "" {
		classHash, err = starknet.HexToBN(proxyAccountClass)
		if err != nil {
			return nil, err
		}
		accountClassHash, err := starknet.HexToBN(accountClass)
		if err != nil {
			return nil, err
		}

		constructorCallData = []*big.Int{accountClassHash, starknet.GetSelectorFromName("initialize")}
		calldate := []*big.Int{big.NewInt(2), pub, big.NewInt(0)}
		constructorCallData = append(constructorCallData, calldate...)
		// calculate address
		contractAddress, err = starknet.CalculateContractAddressFromHash(starkPub, accountClass, proxyAccountClass)
		if err != nil {
			return nil, err
		}
	} else { // cairo1
		classHash, err = starknet.HexToBN(accountClass)
		if err != nil {
			return nil, err
		}
		constructorCallData = []*big.Int{pub, big.NewInt(0)}
		contractAddress, err = cairo1.CalculateContractAddressFromHash(starkPub, accountClass)
		if err != nil {
			return nil, err
		}
	}

	txn := &DeployAccountTxnV3{
		Type:                starknet.TransactionTypeDeployAccount,
		Version:             "0x3",
		Nonce:               BigIntToFelt(nonce),
		ContractAddressSalt: BigIntToFelt(pub),
		ConstructorCalldata: BigIntToFelts(constructorCallData),
		ClassHash:           BigIntToFelt(classHash),
		ResourceBounds: ResourceBoundsMapping{
			L1Gas: ResourceBounds{
				MaxAmount:       U64(l1GasMaxAmount),
				MaxPricePerUnit: U128(l1GasMaxPricePerUnit),
			},
			L2Gas: ResourceBounds{
				MaxAmount:       U64(l2GasMaxAmount),
				MaxPricePerUnit: U128(l2GasMaxPricePerUnit),
			},
			L1DataGas: ResourceBounds{
				MaxAmount:       U64(l1DataGasMaxAmount),
				MaxPricePerUnit: U128(l1DataMaxPricePerUnit),
			},
		},
		Tip:             U64(tip),
		PayMasterData:   []*felt.Felt{},
		NonceDataMode:   DAModeL1,
		FeeMode:         DAModeL1,
		TransactionHash: nil,
	}
	txn.TransactionHash, err = TransactionHashDeployAccountV3(txn, BigIntToFelt(contractAddress), new(felt.Felt).SetBytes([]byte(chainId)))
	if err != nil {
		return nil, err
	}

	return txn, nil
}

func CreateSignedDeployAccountTxV3(curve starknet.StarkCurve, starkPub string, accountClass, proxyAccountClass string, nonce *big.Int, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit string, chainId string, privKeyHex string) (*DeployAccountTxnV3, error) {
	tx, err := CreateDeployAccountTxV3(starkPub, accountClass, proxyAccountClass, nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, chainId)
	if err != nil {
		return nil, err
	}
	if err := starknet.SignTx(curve, tx, privKeyHex); err != nil {
		return nil, err
	}
	return tx, nil
}

func CreateTransferTxV3(contractAddr, from, to string, amount, nonce *big.Int, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit string, chainId string, cairoVersion int) (*InvokeTxnV3, error) {
	contractBn, err := starknet.HexToBN(contractAddr)
	if err != nil {
		return nil, err
	}
	toBn, err := starknet.HexToBN(to)
	if err != nil {
		return nil, err
	}
	transaction := starknet.Transaction{
		ContractAddress:    contractBn,
		EntryPointSelector: starknet.GetSelectorFromName("transfer"),
		Calldata:           []*big.Int{toBn, amount, big.NewInt(0)},
	}

	txs := []starknet.Transaction{transaction}

	var calldata []*big.Int
	if cairoVersion == 1 {
		calldata = starknet.FmtExecuteCalldataCairo1(txs)
	} else if cairoVersion == 0 {
		calldata = starknet.FmtExecuteCalldata(txs)
	}

	txn := &InvokeTxnV3{
		Type:          starknet.TransactionTypeInvoke,
		SenderAddress: felt.HexToFelt(from),
		Calldata:      BigIntToFelts(calldata),
		Version:       "0x3",
		Signature:     nil,
		Nonce:         BigIntToFelt(nonce),
		ResourceBounds: ResourceBoundsMapping{
			L1Gas: ResourceBounds{
				MaxAmount:       U64(l1GasMaxAmount),
				MaxPricePerUnit: U128(l1GasMaxPricePerUnit),
			},
			L2Gas: ResourceBounds{
				MaxAmount:       U64(l2GasMaxAmount),
				MaxPricePerUnit: U128(l2GasMaxPricePerUnit),
			},
			L1DataGas: ResourceBounds{
				MaxAmount:       U64(l1DataGasMaxAmount),
				MaxPricePerUnit: U128(l1DataMaxPricePerUnit),
			},
		},
		Tip:                   U64(tip),
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         DAModeL1,
		FeeMode:               DAModeL1,
	}

	txn.TransactionHash, err = TransactionHashInvokeV3(txn, new(felt.Felt).SetBytes([]byte(chainId)))
	if err != nil {
		return nil, err
	}

	return txn, nil
}

func CreateSignedTransferTxV3(curve starknet.StarkCurve, contractAddr, from, to string, amount, nonce *big.Int, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit string, chainId string, cairoVersion int, privKeyHex string) (*InvokeTxnV3, error) {
	tx, err := CreateTransferTxV3(contractAddr, from, to, amount, nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, chainId, cairoVersion)
	if err != nil {
		return nil, err
	}
	if err := starknet.SignTx(curve, tx, privKeyHex); err != nil {
		return nil, err
	}
	return tx, nil
}

func CreateMultiContractTxV3(from string, calls []starknet.Calls, nonce *big.Int, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit string, chainId string, cairoVersion int) (*InvokeTxnV3, error) {
	var txs []starknet.Transaction
	for _, call := range calls {
		var callDatas []*big.Int
		for _, v := range call.Calldata {
			if strings.HasPrefix(v, "0x") {
				bigV, err := starknet.HexToBN(v)
				if err != nil {
					return nil, err
				}
				callDatas = append(callDatas, bigV)
			} else {
				bigCall, ok := new(big.Int).SetString(v, 10)
				if !ok {
					return nil, errors.New("calldata error")
				}
				callDatas = append(callDatas, bigCall)
			}
		}

		var hashEntryPoint *big.Int
		if strings.HasPrefix(call.Entrypoint, "0x") {
			hashEntryPointBn, err := starknet.HexToBN(call.Entrypoint)
			if err != nil {
				return nil, err
			}
			hashEntryPoint = hashEntryPointBn
		} else {
			hashEntryPoint = starknet.GetSelectorFromName(call.Entrypoint)
		}
		contracBn, err := starknet.HexToBN(call.ContractAddress)
		if err != nil {
			return nil, err
		}
		transaction := starknet.Transaction{
			ContractAddress:    contracBn,
			EntryPointSelector: hashEntryPoint,
			Calldata:           callDatas,
		}
		txs = append(txs, transaction)
	}

	var calldata []*big.Int
	if cairoVersion == 1 {
		calldata = starknet.FmtExecuteCalldataCairo1(txs)
	} else if cairoVersion == 0 {
		calldata = starknet.FmtExecuteCalldata(txs)
	}

	txn := &InvokeTxnV3{
		Type:          starknet.TransactionTypeInvoke,
		SenderAddress: felt.HexToFelt(from),
		Calldata:      BigIntToFelts(calldata),
		Version:       "0x3",
		Signature:     nil,
		Nonce:         BigIntToFelt(nonce),
		ResourceBounds: ResourceBoundsMapping{
			L1Gas: ResourceBounds{
				MaxAmount:       U64(l1GasMaxAmount),
				MaxPricePerUnit: U128(l1GasMaxPricePerUnit),
			},
			L2Gas: ResourceBounds{
				MaxAmount:       U64(l2GasMaxAmount),
				MaxPricePerUnit: U128(l2GasMaxPricePerUnit),
			},
			L1DataGas: ResourceBounds{
				MaxAmount:       U64(l1DataGasMaxAmount),
				MaxPricePerUnit: U128(l1DataMaxPricePerUnit),
			},
		},
		Tip:                   U64(tip),
		PayMasterData:         []*felt.Felt{},
		AccountDeploymentData: []*felt.Felt{},
		NonceDataMode:         DAModeL1,
		FeeMode:               DAModeL1,
	}

	hash, err := TransactionHashInvokeV3(txn, new(felt.Felt).SetBytes([]byte(chainId)))
	if err != nil {
		return nil, err
	}

	txn.TransactionHash = hash
	return txn, nil
}

func CreateSignedMultiContractTxV3(curve starknet.StarkCurve, from string, calls []starknet.Calls, nonce *big.Int, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit string, chainId string, cairoVersion int, privKeyHex string) (*InvokeTxnV3, error) {
	tx, err := CreateMultiContractTxV3(from, calls, nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, chainId, cairoVersion)
	if err != nil {
		return nil, err
	}
	if err := starknet.SignTx(curve, tx, privKeyHex); err != nil {
		return nil, err
	}
	return tx, nil
}

func GetTxHashV3(jsonTx string, chainId string) (string, error) {
	var txType struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal([]byte(jsonTx), &txType)
	if err != nil {
		return "", err
	}
	if txType.Type == "INVOKE_FUNCTION" {
		var request InvokeTxnV3Request
		err = json.Unmarshal([]byte(jsonTx), &request)
		if err != nil {
			return "", err
		}
		tx, err := request.ToTxn(chainId)
		if err != nil {
			return "", err
		}
		return starknet.BigToHex(tx.GetTxHash()), nil
	} else if txType.Type == "DEPLOY_ACCOUNT" {
		var request DeployAccountTxnV3Request
		err = json.Unmarshal([]byte(jsonTx), &request)
		if err != nil {
			return "", err
		}

		var contractAddress *big.Int

		if len(request.ConstructorCalldata) == CairoCalldataLength {
			contractAddress, err = cairo1.CalculateContractAddressFromHash(request.ContractAddressSalt, request.ClassHash)
		} else {
			contractAddress, err = starknet.CalculateContractAddressFromHash(request.ContractAddressSalt, request.ConstructorCalldata[0], request.ClassHash)
		}

		tx, err := request.ToTxn(starknet.BigToHex(contractAddress), chainId)
		if err != nil {
			return "", err
		}
		return starknet.BigToHex(tx.GetTxHash()), nil
	}

	return "", errors.New("unsupported transaction type: " + txType.Type)

}
