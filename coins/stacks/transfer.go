package stacks

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
)

func Transfer(privateKey string, to string, memo string, amount *big.Int, nonce *big.Int, fee *big.Int) (string, error) {
	publicKeyHexStr, err := GetPublicKey(privateKey)
	if err != nil {
		return "", err
	}
	signedtokentransferoptions := createSignedTokenTransferOptions(to, *amount, *fee, *nonce, memo, privateKey)
	stacksTransaction, err := makeUnsignedSTXTokenTransfer(signedtokentransferoptions, publicKeyHexStr)
	if err != nil {
		return "", err
	}
	stacksPrivateKey, err := createStacksPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	stacks := StandardAuthorization{
		stacksTransaction.Auth.AuthType,
		&SingleSigSpendingCondition{
			HashMode:    stacksTransaction.Auth.SpendingCondition.HashMode,
			Signer:      stacksTransaction.Auth.SpendingCondition.Signer,
			Nonce:       *big.NewInt(int64(0)),
			Fee:         *big.NewInt(int64(0)),
			KeyEncoding: stacksTransaction.Auth.SpendingCondition.KeyEncoding,
			Signature:   stacksTransaction.Auth.SpendingCondition.Signature,
		},
		stacksTransaction.Auth.SponsorSpendingCondition,
	}
	tmpAuth := stacksTransaction.Auth
	stacksTransaction.Auth = stacks
	tx := Txid(*stacksTransaction)
	stacksTransaction.Auth = tmpAuth

	signer := &TransactionSigner{
		transaction:   stacksTransaction,
		sigHash:       tx,
		originDone:    false,
		checkOversign: true,
		checkOverlap:  true,
	}

	signer.signOrigin(stacksPrivateKey)
	stacksTransaction.Auth.SpendingCondition.Signature = signer.transaction.Auth.SpendingCondition.Signature

	buf := bytes.NewBuffer(make([]byte, 0))
	buf.Write(getBytes(int64(stacksTransaction.Version), 0))
	chainIdBuffer := bytes.NewBuffer(make([]byte, 0, 4))

	chainIdBuffer.Write(getBytesByLength(stacksTransaction.ChainId, 8))
	buf.Write(sliceByteBuffer(chainIdBuffer))
	buf.Write(serializeAuth(&stacksTransaction.Auth))
	buf.Write(getBytes(int64(stacksTransaction.AnchorMode), 0))
	buf.Write(getBytes(int64(stacksTransaction.PostConditionMode), 0))
	buffer2 := serializeLPList(stacksTransaction.PostConditions)
	buf.Write(buffer2)

	buffer3 := SerializePayload(stacksTransaction.Payload)
	buf.Write(buffer3)

	txSerialize := hex.EncodeToString(sliceByteBuffer(buf))
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

func Txid(stacksTransaction StacksTransaction) string {
	serialized := Serialize(stacksTransaction)
	return txidFromData(serialized)
}

func makeUnsignedSTXTokenTransfer(s *SignedTokenTransferOptions, pubKey string) (*StacksTransaction, error) {
	tokenTransferPayload, err := CreateTokenTransferPayload(s.Recipient, s.Amount, s.Memo)
	if err != nil {
		return nil, err
	}
	addressHashMode := uint64(0)
	spendingCondition, err := createSingleSigSpendingCondition(addressHashMode, pubKey, s.Nonce, s.Fee)
	authorization := &StandardAuthorization{}
	authorization.AuthType = 4
	authorization.SpendingCondition = spendingCondition
	lpPostConditions := createLPList(nil, []PostConditionInterface{})
	var chainId = new(int64)
	*chainId = 1
	var anchorMode = new(int)
	*anchorMode = 3
	transactions := newStacksTransaction(0, chainId, *authorization, anchorMode, tokenTransferPayload, 0, lpPostConditions)
	return transactions, nil
}

func newStacksTransaction(version int, chainID *int64, auth StandardAuthorization, anchorMode *int,
	payload Payload, postConditionMode int, postConditions *LPList) *StacksTransaction {
	chainId := int64(0)
	if chainID != nil {
		chainId = *chainID
	} else {
		chainId = 1
	}

	if postConditionMode == 0 {
		postConditionMode = 2
	}
	if anchorMode == nil {
		payloadType := payload.getPayloadType()
		switch payloadType {
		case 4, 3:
			*anchorMode = 1
		case 2, 1, 0:
			*anchorMode = 3
		}
	}
	if postConditions == nil {
		postConditions = createLPList(nil, []PostConditionInterface{})
	}
	return &StacksTransaction{
		Version:           version,
		ChainId:           chainId,
		Auth:              auth,
		AnchorMode:        *anchorMode,
		Payload:           payload,
		PostConditionMode: postConditionMode,
		PostConditions:    *postConditions,
	}
}
