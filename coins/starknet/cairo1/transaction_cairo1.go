package cairo1

import (
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/starknet"
	"math/big"
	"strings"
)

func CreateTransferTx(curve starknet.StarkCurve, contractAddr, from, to string, amount, nonce, maxFee *big.Int, chainId string) (*starknet.Transaction, error) {
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

	fromBn, err := starknet.HexToBN(from)
	if err != nil {
		return nil, err
	}
	hash, err := curve.HashMulticallCairo1(fromBn, nonce, maxFee, starknet.UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}
	tx := &starknet.Transaction{
		Type:               starknet.TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    fromBn,
		Calldata:           starknet.FmtExecuteCalldataCairo1(txs),
		EntryPointSelector: starknet.GetSelectorFromName(starknet.EXECUTE_SELECTOR),
		TransactionHash:    hash,
		SenderAddress:      from,
		Nonce:              nonce,
	}

	return tx, nil
}

func CreateSignedTransferTx(curve starknet.StarkCurve, contractAddr, from, to string, amount, nonce, maxFee *big.Int, chainId string, privKeyHex string) (*starknet.Transaction, error) {
	tx, err := CreateTransferTx(curve, contractAddr, from, to, amount, nonce, maxFee, chainId)
	if err != nil {
		return nil, err
	}
	if err := starknet.SignTx(curve, tx, privKeyHex); err != nil {
		return nil, err
	}
	return tx, nil
}

func CreateContractTx(curve starknet.StarkCurve, contractAddr, from, functionName string, callData []string, nonce, maxFee *big.Int, chainId string) (*starknet.Transaction, error) {
	var callDatas []*big.Int
	for _, v := range callData {
		bigV, err := starknet.HexToBN(v)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, bigV)
	}
	callDatas = append(callDatas, big.NewInt(0))
	contractAddrBn, err := starknet.HexToBN(contractAddr)
	if err != nil {
		return nil, err
	}
	transaction := starknet.Transaction{
		ContractAddress:    contractAddrBn,
		EntryPointSelector: starknet.GetSelectorFromName(functionName),
		Calldata:           callDatas,
	}

	txs := []starknet.Transaction{transaction}

	fromBn, err := starknet.HexToBN(from)
	if err != nil {
		return nil, err
	}
	hash, err := curve.HashMulticallCairo1(fromBn, nonce, maxFee, starknet.UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}
	tx := &starknet.Transaction{
		Type:               starknet.TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    fromBn,
		Calldata:           starknet.FmtExecuteCalldataCairo1(txs),
		EntryPointSelector: starknet.GetSelectorFromName(starknet.EXECUTE_SELECTOR),
		TransactionHash:    hash,
		SenderAddress:      from,
		Nonce:              nonce,
	}

	return tx, nil
}

func CreateSignedMultiContractTx(curve starknet.StarkCurve, from string, calls []starknet.Calls, nonce, maxFee *big.Int, chainId string, privateKey string) (*starknet.Transaction, error) {
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
	fromBn, err := starknet.HexToBN(from)
	if err != nil {
		return nil, err
	}
	hash, err := curve.HashMulticallCairo1(fromBn, nonce, maxFee, starknet.UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}
	tx := &starknet.Transaction{
		Type:               starknet.TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    fromBn,
		Calldata:           starknet.FmtExecuteCalldataCairo1(txs),
		EntryPointSelector: starknet.GetSelectorFromName(starknet.EXECUTE_SELECTOR),
		TransactionHash:    hash,
		SenderAddress:      from,
		Nonce:              nonce,
	}

	if err = starknet.SignTx(curve, tx, privateKey); err != nil {
		return nil, err
	}

	return tx, nil
}

func SignMessageV1(curve starknet.StarkCurve, address string, msg string, chainId string, privKey string) (string, error) {
	jsonMsg := `{
    "accountAddress" : "%s",
    "typedData" : {
          "types" : {
              "StarkNetDomain": [
                  { "name" : "name", "type" : "felt" },
                  { "name" : "version", "type" : "felt" },
                  { "name" : "chainId", "type" : "felt" }
              ],
              "SignatureMessage": [{ "name": "message", "type": "felt" }]
          },
          "primaryType" : "SignatureMessage",
          "domain" : {
              "name" : "account.signature_message",
              "version" : "1",
              "chainId" : "%s"
          },
          "message" : {
              "message": "%s"
          }
  }
}`
	jsonMsg = fmt.Sprintf(jsonMsg, address, starknet.BigToHex(starknet.UTF8StrToBig(chainId)), msg)
	hash, err := starknet.GetMessageHashWithJson(jsonMsg)
	if err != nil {
		return "", err
	}

	return starknet.SignMsg(curve, hash, privKey)
}
