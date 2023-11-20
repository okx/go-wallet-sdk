package stacks

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
)

// Stack Deprecated
func Stack(contractAddress, contractName, functionName, poxAddress string, amount, height, cycle *big.Int, senderKey string, nonce, fee *big.Int) (string, error) {
	postConditionAmount := amount
	address, err := GetPoxAddress(poxAddress)
	if err != nil {
		return "", fmt.Errorf("err")
	}
	functionArgs := []ClarityValue{
		NewUintCV(postConditionAmount),
		address,
		NewUintCV(height),
		NewUintCV(cycle),
	}

	signedContractCallOptions := &SignedContractCallOptions{
		ContractAddress: contractAddress,
		ContractName:    contractName,
		FunctionName:    functionName,
		FunctionArgs:    functionArgs,
		SendKey:         senderKey,
		ValidateWithAbi: true,
		Fee:             *fee,
		Nonce:           *nonce,
	}

	stacksTransaction, err := MakeContractCall(signedContractCallOptions)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	return hex.EncodeToString(Serialize(*stacksTransaction)), nil
}

func ContractCall(senderKey string, from string, to string, memo string, amount *big.Int, contractAddress, contractName, tokenName, functionName string, nonce, fee *big.Int) (string, error) {
	functionArgs := []ClarityValue{
		NewUintCV(amount),
		*NewStandardPrincipalCV(from),
		*NewStandardPrincipalCV(to),
		&SomeCV{OptionalSome, &BufferCV{
			Type_:  Buffer,
			Buffer: []byte(memo),
		}},
	}

	postConditions := []FungiblePostCondition{
		makeStandardFungiblePostCondition(from, 01, amount, AssetInfo{
			type_:        ASSETINFO,
			address:      *createAddress(contractAddress),
			contractName: *createLPString(contractName, nil, nil),
			assetName:    *createLPString(tokenName, nil, nil),
		}),
	}
	anchorMode := 3
	if fee == nil {
		fee = big.NewInt(int64(400000))
	}

	var interfaceSlice []PostConditionInterface
	for _, condition := range postConditions {
		interfaceSlice = append(interfaceSlice, condition)
	}
	signedContractCallOptions := &SignedContractCallOptions{
		ContractAddress:   contractAddress,
		ContractName:      contractName,
		FunctionName:      functionName,
		FunctionArgs:      functionArgs,
		SendKey:           senderKey,
		Fee:               *fee,
		Nonce:             *nonce,
		AnchorMode:        anchorMode,
		PostConditionMode: PostConditionModeDeny,
		PostConditions:    interfaceSlice,
	}
	stacksTransaction, err := MakeContractCall(signedContractCallOptions)
	if err != nil {
		return "", err
	}
	txSerialize := hex.EncodeToString(Serialize(*stacksTransaction))
	txId := Txid(*stacksTransaction)

	transactionRes := TransactionRes{
		TxId:        txId,
		TxSerialize: txSerialize,
	}
	res, err := json.Marshal(transactionRes)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func makeStandardFungiblePostCondition(address string, fungibleConditionCode int, amount *big.Int, assetInfo AssetInfo) FungiblePostCondition {
	standardPrincipal := createStandardPrincipal(address)
	return createFungiblePostCondition(standardPrincipal, fungibleConditionCode, amount, assetInfo)
}

func createFungiblePostCondition(standardPrincipal PostConditionPrincipal, fungibleConditionCode int, amount *big.Int, assetInfo AssetInfo) FungiblePostCondition {
	return FungiblePostCondition{
		PostCondition: PostCondition{
			StacksMessage: StacksMessage{
				Type: POSTCONDITION,
			},
			ConditionType: 1,
			Principal:     standardPrincipal,
			ConditionCode: fungibleConditionCode,
		},
		assetInfo: assetInfo,
		amount:    amount,
	}
}

func createStandardPrincipal(address string) PostConditionPrincipal {
	addObj := createAddress(address)
	return PostConditionPrincipal{
		Type:    PRINCIPAL,
		Prefix:  2,
		Address: *addObj,
	}
}

func MakeContractCall(txOptions *SignedContractCallOptions) (*StacksTransaction, error) {
	sendKey := txOptions.SendKey
	publicKey, err := GetPublicKey(sendKey)
	if err != nil {
		return nil, err
	}
	stacksTransaction, err := makeUnsignedContractCall(publicKey, txOptions)
	if err != nil {
		return nil, err
	}
	privKey, err := createStacksPrivateKey(txOptions.SendKey)
	if err != nil {
		return nil, err
	}

	stacks := StandardAuthorization{
		stacksTransaction.Auth.AuthType,
		&SingleSigSpendingCondition{
			HashMode:    stacksTransaction.Auth.SpendingCondition.HashMode,
			Signer:      stacksTransaction.Auth.SpendingCondition.Signer,
			Nonce:       *big.NewInt(int64(0)),
			Fee:         *big.NewInt(int64(0)),
			KeyEncoding: stacksTransaction.Auth.SpendingCondition.KeyEncoding,
			Signature:   MessageSignature{9, "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		},
		stacksTransaction.Auth.SponsorSpendingCondition,
	}
	tmpAuth := stacksTransaction.Auth
	stacksTransaction.Auth = stacks
	tx := Txid(*stacksTransaction)
	stacksTransaction.Auth = tmpAuth

	signer := &TransactionSigner{
		transaction: stacksTransaction,
		//6fbc7ff90a741471572bf406ba58381bf39743fe8fc749e5dbfc6a88dfbc340e
		sigHash:       tx,
		originDone:    false,
		checkOversign: true,
		checkOverlap:  true,
	}

	if err := signer.signOrigin(privKey); err != nil {
		return nil, err
	}
	stacksTransaction.Auth.SpendingCondition.Signature = signer.transaction.Auth.SpendingCondition.Signature

	return stacksTransaction, nil
}

func makeUnsignedContractCall(publicKey string, txOptions *SignedContractCallOptions) (*StacksTransaction, error) {
	payload := CreateContractCallPayload(txOptions.ContractAddress, txOptions.ContractName, txOptions.FunctionName, txOptions.FunctionArgs)
	spendingCondition, err := createSingleSigSpendingCondition(SerializeP2PKH, publicKey, txOptions.Nonce, txOptions.Fee)
	if err != nil {
		return nil, err
	}
	authorization := StandardAuthorization{4, spendingCondition, nil}

	newPost := make([]PostConditionInterface, len(txOptions.PostConditions))
	if txOptions.PostConditions != nil && len(txOptions.PostConditions) > 0 {
		copy(newPost, txOptions.PostConditions)
	}

	lpPostConditions := &LPList{7, 4, newPost}
	var chainId = new(int64)
	*chainId = 1
	anchorMode := txOptions.AnchorMode
	if anchorMode == 0 {
		anchorMode = 3
	}
	// stacksTransaction := newStacksTransaction(0, chainId, authorization, anchorMode, payload, 2, lpPostConditions)
	stacksTransaction := newStacksTransaction(0, chainId, authorization, &anchorMode, payload, txOptions.PostConditionMode, lpPostConditions)
	return stacksTransaction, nil
}

func makeUnsignedContractCallWithSerializePostCondition(publicKey string, txOptions *SignedContractCallOptions) (*StacksTransaction, error) {
	payload := CreateContractCallPayload(txOptions.ContractAddress, txOptions.ContractName, txOptions.FunctionName, txOptions.FunctionArgs)
	spendingCondition, err := createSingleSigSpendingCondition(SerializeP2PKH, publicKey, txOptions.Nonce, txOptions.Fee)
	if err != nil {
		return nil, err
	}
	authorization := StandardAuthorization{4, spendingCondition, nil}

	var newpost []PostConditionInterface
	if len(txOptions.SerializePostConditions) > 0 {
		for _, ps := range txOptions.SerializePostConditions {
			newpost = append(newpost, DeserializePostCondition(ps))
		}
	}

	lpPostConditions := &LPList{7, 4, newpost}
	var chainId = new(int64)
	*chainId = 1
	anchorMode := txOptions.AnchorMode
	if anchorMode == 0 {
		anchorMode = 3
	}
	stacksTransaction := newStacksTransaction(0, chainId, authorization, &anchorMode, payload, txOptions.PostConditionMode, lpPostConditions)
	return stacksTransaction, nil
}

func CreateContractCallPayload(contractAddress, contractName, functionName string, functionArgs []ClarityValue) *ContractCallPayload {
	address := createAddress(contractAddress)
	lpContractName := createLPString(contractName, nil, nil)
	lpFunctionName := createLPString(functionName, nil, nil)

	return &ContractCallPayload{
		type_:           8,
		payloadType:     2,
		contractAddress: address,
		contractName:    lpContractName,
		functionName:    lpFunctionName,
		functionArgs:    functionArgs,
	}
}
