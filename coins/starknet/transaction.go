package starknet

import (
	"encoding/json"
	"math/big"
	"strings"
)

const (
	TransactionTypeDeploy        = "DEPLOY"
	TransactionTypeInvoke        = "INVOKE_FUNCTION"
	TransactionTypeDeployAccount = "DEPLOY_ACCOUNT"

	GOERLI_ID  = "SN_GOERLI"
	MAINNET_ID = "SN_MAIN"

	EXECUTE_SELECTOR = "__execute__"

	TRANSACTION_PREFIX  string = "invoke"
	DEPLOY_ACCOUNT      string = "deploy_account"
	TRANSACTION_VERSION int64  = 1
)

type signTx interface {
	SetSignature(*big.Int, *big.Int)
	GetTxHash() *big.Int
}

type DeployTransaction struct {
	Type                string     `json:"type"`
	ContractAddressSalt *big.Int   `json:"contract_address_salt"`
	ConstructorCalldata []*big.Int `json:"constructor_calldata"`
	ClassHash           *big.Int   `json:"class_hash"`
	MaxFee              *big.Int   `json:"max_fee"`
	Version             *big.Int   `json:"version"`
	Nonce               *big.Int   `json:"nonce"`
	Signature           []*big.Int `json:"signature"`
	TransactionHash     *big.Int   `json:"-"`
}

func (d *DeployTransaction) SetSignature(x, y *big.Int) {
	d.Signature = []*big.Int{x, y}
}

func (d *DeployTransaction) GetTxHash() *big.Int {
	return d.TransactionHash
}

type DeployTransactionReq struct {
	Type                string   `json:"type"`
	ContractAddressSalt string   `json:"contract_address_salt"`
	ConstructorCalldata []string `json:"constructor_calldata"`
	ClassHash           string   `json:"class_hash"`
	MaxFee              string   `json:"max_fee"`
	Version             string   `json:"version"`
	Nonce               string   `json:"nonce"`
	Signature           []string `json:"signature"`
}

func (d *DeployTransaction) GetDeployAccountReq() DeployTransactionReq {
	calldata := make([]string, len(d.ConstructorCalldata))
	for i, data := range d.ConstructorCalldata {
		calldata[i] = data.String()
	}

	return DeployTransactionReq{
		Type:                d.Type,
		ContractAddressSalt: BigToHex(d.ContractAddressSalt),
		ConstructorCalldata: calldata,
		ClassHash:           BigToHex(d.ClassHash),
		MaxFee:              BigToHex(d.MaxFee),
		Version:             BigToHex(d.Version),
		Nonce:               BigToHex(d.Nonce),
		Signature:           []string{d.Signature[0].String(), d.Signature[1].String()},
	}
}

type Transaction struct {
	Type               string     `json:"type"`
	ContractAddress    *big.Int   `json:"contract_address"`
	Calldata           []*big.Int `json:"calldata"`
	EntryPointSelector *big.Int   `json:"entry_point_selector"`
	Nonce              *big.Int   `json:"nonce,omitempty"`
	TransactionHash    *big.Int   `json:"transaction_hash"`
	MaxFee             *big.Int   `json:"max_fee"`
	Signature          []*big.Int `json:"signature"`
	SenderAddress      string     `json:"sender_address"`
}

func (tx *Transaction) SetSignature(x, y *big.Int) {
	tx.Signature = []*big.Int{x, y}
}

func (tx *Transaction) GetTxHash() *big.Int {
	return tx.TransactionHash
}

type TransactionRequest struct {
	Type          string   `json:"type"`
	SenderAddress string   `json:"sender_address"`
	Calldata      []string `json:"calldata"`
	MaxFee        string   `json:"max_fee"`
	Signature     []string `json:"signature,omitempty"`
	Version       string   `json:"version"`
	Nonce         string   `json:"nonce"`
}

type FunctionInvocation struct {
	ContractAddress    string   `json:"contract_address"`
	Calldata           []string `json:"calldata"`
	EntryPointSelector string   `json:"entry_point_selector"`
}

type Params struct {
	FunctionInvocation FunctionInvocation `json:"function_invocation"`
	Signature          []string           `json:"signature"`
	MaxFee             string             `json:"max_fee"`
	Version            string             `json:"version"`
}

type TransactionJsonRpc struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		FunctionInvocation FunctionInvocation `json:"function_invocation"`
		Signature          []string           `json:"signature"`
		MaxFee             string             `json:"max_fee"`
		Version            string             `json:"version"`
	} `json:"params"`
	Id int `json:"id"`
}

func (tx *Transaction) GetTxRequest() *TransactionRequest {
	calldata := make([]string, len(tx.Calldata))

	for i, data := range tx.Calldata {
		calldata[i] = data.String()
	}

	return &TransactionRequest{
		Type:          tx.Type,
		SenderAddress: strings.ToLower(tx.SenderAddress),
		Calldata:      calldata,
		MaxFee:        BigToHex(tx.MaxFee),
		Signature:     []string{tx.Signature[0].String(), tx.Signature[1].String()},
		Version:       "0x1",
		Nonce:         BigToHex(tx.Nonce),
	}
}

func (tx *Transaction) ToJsonRpcParams() string {
	calldata := make([]string, len(tx.Calldata))

	for i, data := range tx.Calldata {
		calldata[i] = BigToHex(data)
	}

	functionInvocation := FunctionInvocation{
		ContractAddress:    BigToHex(tx.ContractAddress),
		Calldata:           calldata,
		EntryPointSelector: BigToHex(tx.EntryPointSelector),
	}

	signature := make([]string, len(tx.Signature))
	for i, data := range tx.Signature {
		signature[i] = data.String()
	}

	params := Params{
		FunctionInvocation: functionInvocation,
		Signature:          signature,
		MaxFee:             BigToHex(tx.MaxFee),
		Version:            "0x0",
	}

	bytes, _ := json.Marshal(params)
	return string(bytes)
}

func CreateDeployAccountTx(starkPub string, nonce, maxFee *big.Int, chainId string) (*DeployTransaction, error) {
	pub := HexToBN(starkPub)
	version := big.NewInt(TRANSACTION_VERSION)
	accountClassHash := HexToBN(AccountClassHash)
	classHash := HexToBN(ProxyAccountClassHash)

	// build callData
	constructorCallData := []*big.Int{accountClassHash, GetSelectorFromName("initialize")}
	calldate := []*big.Int{big.NewInt(2), pub, big.NewInt(0)}
	constructorCallData = append(constructorCallData, calldate...)

	// calculate address
	contractAddress, err := CalculateContractAddressFromHash(starkPub)
	if err != nil {
		return nil, err
	}

	// calculate tx hash
	txHash, err := calculateDeployAccountTransactionHash(contractAddress, classHash, constructorCallData, pub, version, UTF8StrToBig(chainId), nonce, maxFee)
	if err != nil {
		return nil, err
	}

	return &DeployTransaction{
		Type:                TransactionTypeDeployAccount,
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

func CreateTransferTx(curve StarkCurve, contractAddr, from, to string, amount, nonce, maxFee *big.Int, chainId string) (*Transaction, error) {
	transaction := Transaction{
		ContractAddress:    HexToBN(contractAddr),
		EntryPointSelector: GetSelectorFromName("transfer"),
		Calldata:           []*big.Int{HexToBN(to), amount, big.NewInt(0)},
	}

	txs := []Transaction{transaction}

	hash, err := curve.HashMulticall(HexToBN(from), nonce, maxFee, UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}
	tx := &Transaction{
		Type:               TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    HexToBN(from),
		Calldata:           FmtExecuteCalldata(txs),
		EntryPointSelector: GetSelectorFromName(EXECUTE_SELECTOR),
		TransactionHash:    hash,
		SenderAddress:      from,
		Nonce:              nonce,
	}

	return tx, nil
}

func CreateSignedTransferTx(curve StarkCurve, contractAddr, from, to string, amount, nonce, maxFee *big.Int, chainId string, privKeyHex string) (*Transaction, error) {
	tx, err := CreateTransferTx(curve, contractAddr, from, to, amount, nonce, maxFee, chainId)
	if err != nil {
		return nil, err
	}
	if err := SignTx(curve, tx, privKeyHex); err != nil {
		return nil, err
	}
	return tx, nil
}

func CreateSignedDeployAccountTx(curve StarkCurve, starkPub string, nonce, maxFee *big.Int, chainId string, privateKey string) (*DeployTransaction, error) {
	tx, err := CreateDeployAccountTx(starkPub, nonce, maxFee, chainId)
	if err != nil {
		return nil, err
	}

	if err = SignTx(curve, tx, privateKey); err != nil {
		return nil, err
	}
	return tx, nil
}

func SignTx(curve StarkCurve, tx signTx, privKeyHex string) error {
	x, y, err := curve.Sign(tx.GetTxHash(), HexToBN(privKeyHex))
	if err != nil {
		return err
	}
	tx.SetSignature(x, y)
	return nil
}

func SignMsg(curve StarkCurve, msg string, privKey string) (string, error) {
	x, y, err := curve.PrivateToPoint(HexToBN(privKey))
	if err != nil {
		return "", err
	}
	r, s, err := curve.Sign(HexToBN(msg), HexToBN(privKey))
	if err != nil {
		return "", err
	}
	sig, err := json.Marshal(struct {
		X string `json:"publicKey"`
		Y string `json:"publicKeyY"`
		R string `json:"signedDataR"`
		S string `json:"signedDataS"`
	}{BigToHexWithPadding(x), BigToHexWithPadding(y), BigToHexWithPadding(r), BigToHexWithPadding(s)})

	return string(sig), err
}

func CreateContractTx(curve StarkCurve, contractAddr, from, functionName string, callData []string, nonce, maxFee *big.Int, chainId string) (*Transaction, error) {
	var callDatas []*big.Int
	for _, v := range callData {
		callDatas = append(callDatas, HexToBN(v))
	}
	callDatas = append(callDatas, big.NewInt(0))
	transaction := Transaction{
		ContractAddress:    HexToBN(contractAddr),
		EntryPointSelector: GetSelectorFromName(functionName),
		Calldata:           callDatas,
	}

	txs := []Transaction{transaction}

	hash, err := curve.HashMulticall(HexToBN(from), nonce, maxFee, UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}
	tx := &Transaction{
		Type:               TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    HexToBN(from),
		Calldata:           FmtExecuteCalldata(txs),
		EntryPointSelector: GetSelectorFromName(EXECUTE_SELECTOR),
		TransactionHash:    hash,
		SenderAddress:      from,
		Nonce:              nonce,
	}

	return tx, nil
}

func CreateSignedContractTx(curve StarkCurve, contractAddr, from, functionName string, callData []string, nonce, maxFee *big.Int, chainId string, privateKey string) (*Transaction, error) {
	tx, err := CreateContractTx(curve, contractAddr, from, functionName, callData, nonce, maxFee, chainId)
	if err != nil {
		return nil, err
	}
	if err = SignTx(curve, tx, privateKey); err != nil {
		return nil, err
	}

	return tx, nil
}

func CreateSignedUpgradeTx(curve StarkCurve, from string, nonce, maxFee *big.Int, chainId string, privateKey string) (*Transaction, error) {
	transaction := Transaction{
		ContractAddress:    HexToBN(from),
		EntryPointSelector: GetSelectorFromName("upgrade"),
		Calldata:           []*big.Int{HexToBN("0x33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2")},
	}

	txs := []Transaction{transaction}

	hash, err := curve.OldHashMulticall(HexToBN(from), nonce, maxFee, UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		Type:               TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    HexToBN(from),
		Calldata:           OldFmtExecuteCalldata(nonce, txs),
		EntryPointSelector: GetSelectorFromName(EXECUTE_SELECTOR),
		TransactionHash:    hash,
	}
	if err = SignTx(curve, tx, privateKey); err != nil {
		return nil, err
	}
	return tx, nil
}
