package starknet

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (d *DeployTransaction) UnmarshalJSON(data []byte) error {
	temp := &struct {
		Type                string   `json:"type"`
		ContractAddressSalt string   `json:"contract_address_salt"`
		ConstructorCalldata []string `json:"constructor_calldata"`
		ClassHash           string   `json:"class_hash"`
		MaxFee              string   `json:"max_fee"`
		Version             string   `json:"version"`
		Nonce               string   `json:"nonce"`
		Signature           []string `json:"signature"`
		TransactionHash     string   `json:"-"`
	}{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	d.Type = temp.Type

	d.ContractAddressSalt = new(big.Int)
	d.ContractAddressSalt.SetString(temp.ContractAddressSalt, 16)

	for _, call := range temp.ConstructorCalldata {
		c, ok := new(big.Int).SetString(call, 10)
		if !ok {
			return fmt.Errorf("unmarshal json error")
		}
		d.ConstructorCalldata = append(d.ConstructorCalldata, c)
	}

	d.ClassHash = new(big.Int)
	d.ClassHash.SetString(temp.ClassHash[2:], 16)

	d.MaxFee = new(big.Int)
	d.MaxFee.SetString(temp.MaxFee[2:], 16)

	d.Version = new(big.Int)
	d.Version.SetString(temp.Version[2:], 16)

	d.Nonce = new(big.Int)
	d.Nonce.SetString(temp.Nonce[2:], 16)

	for _, sig := range temp.Signature {
		s, ok := new(big.Int).SetString(sig, 10)
		if !ok {
			return fmt.Errorf("unmarshal json error")
		}
		d.Signature = append(d.Signature, s)
	}

	return nil
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

func (d *DeployTransaction) GetDeployAccountReqWithOutSign() DeployTransactionReq {
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
	}
}

type Calls struct {
	ContractAddress string   `json:"contract_address"`
	Entrypoint      string   `json:"entry_point_selector"`
	Calldata        []string `json:"calldata"`
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

func (t *Transaction) UnmarshalJSON(data []byte) error {
	temp := &struct {
		Type               string   `json:"type"`
		ContractAddress    string   `json:"contract_address"`
		Calldata           []string `json:"calldata"`
		EntryPointSelector string   `json:"entry_point_selector"`
		Nonce              string   `json:"nonce,omitempty"`
		TransactionHash    string   `json:"transaction_hash"`
		MaxFee             string   `json:"max_fee"`
		Signature          []string `json:"signature"`
		SenderAddress      string   `json:"sender_address"`
	}{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	t.Type = temp.Type

	t.ContractAddress = new(big.Int)
	t.ContractAddress.SetString(temp.ContractAddress[2:], 16)

	for _, call := range temp.Calldata {
		c, ok := new(big.Int).SetString(call, 10)
		if !ok {
			return fmt.Errorf("unmarshal json error")
		}
		t.Calldata = append(t.Calldata, c)
	}

	t.EntryPointSelector = new(big.Int)
	t.EntryPointSelector.SetString(temp.EntryPointSelector[2:], 16)

	t.Nonce = new(big.Int)
	t.Nonce.SetString(temp.Nonce[2:], 16)

	t.TransactionHash = new(big.Int)
	t.TransactionHash.SetString(temp.TransactionHash[2:], 16)

	t.MaxFee = new(big.Int)
	t.MaxFee.SetString(temp.MaxFee[2:], 16)

	for _, sig := range temp.Signature {
		s, ok := new(big.Int).SetString(sig, 10)
		if !ok {
			return fmt.Errorf("unmarshal json error")
		}
		t.Signature = append(t.Signature, s)
	}

	t.SenderAddress = temp.SenderAddress
	return nil
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

type OldTransactionRequest struct {
	Type               string   `json:"type"`
	ContractAddress    string   `json:"contract_address"`
	Calldata           []string `json:"calldata"`
	EntryPointSelector string   `json:"entry_point_selector"`
	MaxFee             string   `json:"max_fee"`
	Signature          []string `json:"signature,omitempty"`
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

func (tx *Transaction) GetTxRequestWithOutSign() *TransactionRequest {
	calldata := make([]string, len(tx.Calldata))

	for i, data := range tx.Calldata {
		calldata[i] = data.String()
	}

	return &TransactionRequest{
		Type:          tx.Type,
		SenderAddress: strings.ToLower(tx.SenderAddress),
		Calldata:      calldata,
		MaxFee:        BigToHex(tx.MaxFee),
		Version:       "0x1",
		Nonce:         BigToHex(tx.Nonce),
	}
}

func (tx Transaction) GetOldTxRequest() *OldTransactionRequest {
	calldata := make([]string, len(tx.Calldata))

	for i, data := range tx.Calldata {
		calldata[i] = data.String()
	}

	return &OldTransactionRequest{
		Type:               tx.Type,
		ContractAddress:    BigToHex(tx.ContractAddress),
		Calldata:           calldata,
		EntryPointSelector: BigToHex(tx.EntryPointSelector),
		MaxFee:             BigToHex(tx.MaxFee),
		Signature:          []string{tx.Signature[0].String(), tx.Signature[1].String()},
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

	bytes, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CreateDeployAccountTx(starkPub string, nonce, maxFee *big.Int, chainId string) (*DeployTransaction, error) {
	pub, err := HexToBN(starkPub)
	if err != nil {
		return nil, err
	}
	version := big.NewInt(TRANSACTION_VERSION)
	accountClassHash, err := HexToBN(AccountClassHash)
	if err != nil {
		return nil, err
	}
	classHash, err := HexToBN(ProxyAccountClassHash)
	if err != nil {
		return nil, err
	}

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
	contractBn, err := HexToBN(contractAddr)
	if err != nil {
		return nil, err
	}
	toBn, err := HexToBN(to)
	if err != nil {
		return nil, err
	}
	transaction := Transaction{
		ContractAddress:    contractBn,
		EntryPointSelector: GetSelectorFromName("transfer"),
		Calldata:           []*big.Int{toBn, amount, big.NewInt(0)},
	}

	txs := []Transaction{transaction}

	fromBn, err := HexToBN(from)
	if err != nil {
		return nil, err
	}
	hash, err := curve.HashMulticall(fromBn, nonce, maxFee, UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}
	tx := &Transaction{
		Type:               TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    fromBn,
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
	priBn, err := HexToBN(privKeyHex)
	if err != nil {
		return nil
	}
	x, y, err := curve.Sign(tx.GetTxHash(), priBn)
	if err != nil {
		return err
	}
	tx.SetSignature(x, y)
	return nil
}

func SignMsg(curve StarkCurve, msg string, privKey string) (string, error) {
	priBn, err := HexToBN(privKey)
	if err != nil {
		return "", err
	}
	x, y, err := curve.PrivateToPoint(priBn)
	if err != nil {
		return "", err
	}
	msgBn, err := HexToBN(msg)
	if err != nil {
		return "", err
	}
	r, s, err := curve.Sign(msgBn, priBn)
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

func CreateSignedMultiContractTx(curve StarkCurve, from string, calls []Calls, nonce, maxFee *big.Int, chainId string, privateKey string) (*Transaction, error) {
	var txs []Transaction
	for _, call := range calls {
		var callDatas []*big.Int
		for _, v := range call.Calldata {
			if strings.HasPrefix(v, "0x") {
				bigV, err := HexToBN(v)
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
			hashEntryPointBn, err := HexToBN(call.Entrypoint)
			if err != nil {
				return nil, err
			}
			hashEntryPoint = hashEntryPointBn
		} else {
			hashEntryPoint = GetSelectorFromName(call.Entrypoint)
		}
		contracBn, err := HexToBN(call.ContractAddress)
		if err != nil {
			return nil, err
		}
		transaction := Transaction{
			ContractAddress:    contracBn,
			EntryPointSelector: hashEntryPoint,
			Calldata:           callDatas,
		}
		txs = append(txs, transaction)
	}
	fromBn, err := HexToBN(from)
	if err != nil {
		return nil, err
	}
	hash, err := curve.HashMulticall(fromBn, nonce, maxFee, UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}
	tx := &Transaction{
		Type:               TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    fromBn,
		Calldata:           FmtExecuteCalldata(txs),
		EntryPointSelector: GetSelectorFromName(EXECUTE_SELECTOR),
		TransactionHash:    hash,
		SenderAddress:      from,
		Nonce:              nonce,
	}

	if err = SignTx(curve, tx, privateKey); err != nil {
		return nil, err
	}

	return tx, nil
}

func CreateContractTx(curve StarkCurve, contractAddr, from, functionName string, callData []string, nonce, maxFee *big.Int, chainId string) (*Transaction, error) {
	var callDatas []*big.Int
	for _, v := range callData {
		bigV, err := HexToBN(v)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, bigV)
	}
	callDatas = append(callDatas, big.NewInt(0))
	contractAddrBn, err := HexToBN(contractAddr)
	if err != nil {
		return nil, err
	}
	transaction := Transaction{
		ContractAddress:    contractAddrBn,
		EntryPointSelector: GetSelectorFromName(functionName),
		Calldata:           callDatas,
	}

	txs := []Transaction{transaction}

	fromBn, err := HexToBN(from)
	if err != nil {
		return nil, err
	}
	hash, err := curve.HashMulticall(fromBn, nonce, maxFee, UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}
	tx := &Transaction{
		Type:               TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    fromBn,
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
	fromBn, err := HexToBN(from)
	if err != nil {
		return nil, err
	}
	callBn, err := HexToBN("0x33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2")
	if err != nil {
		return nil, err
	}
	transaction := Transaction{
		ContractAddress:    fromBn,
		EntryPointSelector: GetSelectorFromName("upgrade"),
		Calldata:           []*big.Int{callBn},
	}

	txs := []Transaction{transaction}

	hash, err := curve.OldHashMulticall(fromBn, nonce, maxFee, UTF8StrToBig(chainId), txs)
	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		Type:               TransactionTypeInvoke,
		MaxFee:             maxFee,
		ContractAddress:    fromBn,
		Calldata:           OldFmtExecuteCalldata(nonce, txs),
		EntryPointSelector: GetSelectorFromName(EXECUTE_SELECTOR),
		TransactionHash:    hash,
	}
	if err = SignTx(curve, tx, privateKey); err != nil {
		return nil, err
	}
	return tx, nil
}

func GetTxHash(txStr string) (string, error) {
	curve := SC()

	var data map[string]interface{}
	err := json.Unmarshal([]byte(txStr), &data)
	if err != nil {
		return "", err
	}

	var ok bool
	ty, ok := data["type"].(string)
	if !ok {
		return "", fmt.Errorf("json parse error")
	}

	n, ok := data["nonce"].(string)
	if !ok {
		return "", fmt.Errorf("json parse error")
	}
	nonce, err := HexToBN(n)
	if err != nil {
		return "", err
	}

	m, ok := data["max_fee"].(string)
	if !ok {
		return "", fmt.Errorf("json parse error")
	}
	maxFee, err := HexToBN(m)
	if err != nil {
		return "", err
	}

	switch ty {
	case "INVOKE_FUNCTION":
		from, ok := data["sender_address"].(string)
		if !ok {
			return "", fmt.Errorf("json parse error")
		}

		calldata, ok := data["calldata"].([]interface{})
		if !ok {
			return "", fmt.Errorf("json parse error")
		}
		amount, _ := new(big.Int).SetString(calldata[7].(string), 10)
		to, _ := new(big.Int).SetString(calldata[6].(string), 10)
		selector, _ := new(big.Int).SetString(calldata[2].(string), 10)
		contractAddress, _ := new(big.Int).SetString(calldata[1].(string), 10)
		transaction := Transaction{
			ContractAddress:    contractAddress,
			EntryPointSelector: selector,
			Calldata:           []*big.Int{to, amount, big.NewInt(0)},
		}
		txs := []Transaction{transaction}
		fromBn, err := HexToBN(from)
		if err != nil {
			return "", err
		}
		hash, err := curve.HashMulticall(fromBn, nonce, maxFee, UTF8StrToBig(MAINNET_ID), txs)
		if err != nil {
			return "", err
		}
		return BigToHex(hash), nil
	case "DEPLOY_ACCOUNT":
		version := big.NewInt(TRANSACTION_VERSION)
		accountClassHash, err := HexToBN(AccountClassHash)
		if err != nil {
			return "", err
		}
		classHash, err := HexToBN(ProxyAccountClassHash)
		if err != nil {
			return "", err
		}

		pub, ok := data["contract_address_salt"].(string)
		if !ok {
			return "", fmt.Errorf("json parse error")
		}
		// calculate address
		contractAddress, err := CalculateContractAddressFromHash(pub)
		if err != nil {
			return "", err
		}

		// build callData
		pubBn, err := HexToBN(pub)
		if err != nil {
			return "", err
		}
		constructorCallData := []*big.Int{accountClassHash, GetSelectorFromName("initialize")}
		calldate := []*big.Int{big.NewInt(2), pubBn, big.NewInt(0)}
		constructorCallData = append(constructorCallData, calldate...)

		txHash, err := calculateDeployAccountTransactionHash(contractAddress, classHash, constructorCallData, pubBn, version, UTF8StrToBig(MAINNET_ID), nonce, maxFee)
		if err != nil {
			return "", err
		}
		return BigToHex(txHash), err
	}

	return "", err
}
