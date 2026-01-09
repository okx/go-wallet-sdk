package v3

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/okx/go-wallet-sdk/coins/starknet/juno_core/felt"
)

type Resource string

const (
	ResourceL1Gas     Resource = "L1_GAS"
	ResourceL2Gas     Resource = "L2_GAS"
	ResourceL1DataGas Resource = "L1_DATA"
)

type InvokeTxnV3 struct {
	Type           string                `json:"type"`
	SenderAddress  *felt.Felt            `json:"sender_address"`
	Calldata       []*felt.Felt          `json:"calldata"`
	Version        string                `json:"version"`
	Signature      []*felt.Felt          `json:"signature"`
	Nonce          *felt.Felt            `json:"nonce"`
	ResourceBounds ResourceBoundsMapping `json:"resource_bounds"`
	Tip            U64                   `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []*felt.Felt `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode         DataAvailabilityMode `json:"fee_data_availability_mode"`
	TransactionHash *felt.Felt           `json:"-"`
}

func (tx *InvokeTxnV3) SetSignature(x *big.Int, y *big.Int) {
	tx.Signature = []*felt.Felt{BigIntToFelt(x), BigIntToFelt(y)}
}

func (tx *InvokeTxnV3) GetTxHash() *big.Int {
	return felt.FeltToBigInt(tx.TransactionHash)
}

func (tx *InvokeTxnV3) GetTxRequest() *InvokeTxnV3Request {
	calldata := make([]string, len(tx.Calldata))
	for i, data := range tx.Calldata {
		calldata[i] = data.String()
	}

	payMasterData := make([]string, len(tx.PayMasterData))
	for i, data := range tx.PayMasterData {
		payMasterData[i] = data.String()
	}

	accountDeploymentData := make([]string, len(tx.AccountDeploymentData))
	for i, data := range tx.AccountDeploymentData {
		accountDeploymentData[i] = data.String()
	}

	signature := []string{}
	if len(tx.Signature) == 2 {
		signature = append(signature, tx.Signature[0].String())
		signature = append(signature, tx.Signature[1].String())
	}

	nonceDataMode, err := tx.NonceDataMode.UInt64()
	if err != nil {
		return nil
	}
	feeMode, err := tx.FeeMode.UInt64()
	if err != nil {
		return nil
	}
	return &InvokeTxnV3Request{
		Type:                  tx.Type,
		SenderAddress:         tx.SenderAddress.String(),
		Calldata:              calldata,
		Version:               tx.Version,
		Signature:             signature,
		Nonce:                 tx.Nonce.String(),
		ResourceBounds:        tx.ResourceBounds,
		Tip:                   tx.Tip,
		PayMasterData:         payMasterData,
		AccountDeploymentData: accountDeploymentData,
		NonceDataMode:         nonceDataMode,
		FeeMode:               feeMode,
	}
}

func (tx *InvokeTxnV3) GetTxRequestJson() string {
	request := tx.GetTxRequest()
	jstr, err := json.Marshal(request)
	if err != nil {
		return ""
	}
	return string(jstr)
}

type InvokeTxnV3Request struct {
	Type           string                `json:"type"`
	SenderAddress  string                `json:"sender_address"`
	Calldata       []string              `json:"calldata"`
	Version        string                `json:"version"`
	Signature      []string              `json:"signature"`
	Nonce          string                `json:"nonce"`
	ResourceBounds ResourceBoundsMapping `json:"resource_bounds"`
	Tip            U64                   `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []string `json:"paymaster_data"`
	// The data needed to deploy the account contract from which this tx will be initiated
	AccountDeploymentData []string `json:"account_deployment_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode uint64 `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode uint64 `json:"fee_data_availability_mode"`
}

func (req *InvokeTxnV3Request) ToTxn(chainId string) (*InvokeTxnV3, error) {
	caldatas := make([]*felt.Felt, len(req.Calldata))
	for i, s := range req.Calldata {
		f, err := new(felt.Felt).SetString(s)
		if err != nil {
			return nil, err
		}
		caldatas[i] = f
	}

	signature := make([]*felt.Felt, len(req.Signature))
	for i, s := range req.Signature {
		f, err := new(felt.Felt).SetString(s)
		if err != nil {
			return nil, err
		}
		signature[i] = f
	}

	payMasterData := make([]*felt.Felt, len(req.PayMasterData))
	for i, s := range req.PayMasterData {
		f, err := new(felt.Felt).SetString(s)
		if err != nil {
			return nil, err
		}
		payMasterData[i] = f
	}

	accountDeploymentData := make([]*felt.Felt, len(req.AccountDeploymentData))
	for i, s := range req.AccountDeploymentData {
		f, err := new(felt.Felt).SetString(s)
		if err != nil {
			return nil, err
		}
		accountDeploymentData[i] = f
	}

	nonce, err := new(felt.Felt).SetString(req.Nonce)
	if err != nil {
		return nil, err
	}

	sender, err := new(felt.Felt).SetString(req.SenderAddress)
	if err != nil {
		return nil, err
	}

	nonceDataMode, err := getDataAvailabilityMode(req.NonceDataMode)
	if err != nil {
		return nil, err
	}

	feeMode, err := getDataAvailabilityMode(req.FeeMode)
	if err != nil {
		return nil, err
	}

	txn := &InvokeTxnV3{
		Type:                  req.Type,
		SenderAddress:         sender,
		Calldata:              caldatas,
		Version:               req.Version,
		Signature:             signature,
		Nonce:                 nonce,
		ResourceBounds:        req.ResourceBounds,
		Tip:                   req.Tip,
		PayMasterData:         payMasterData,
		AccountDeploymentData: accountDeploymentData,
		NonceDataMode:         nonceDataMode,
		FeeMode:               feeMode,
	}
	txn.TransactionHash, err = TransactionHashInvokeV3(txn, new(felt.Felt).SetBytes([]byte(chainId)))
	if err != nil {
		return nil, err
	}
	return txn, nil
}

type DeployAccountTxnV3 struct {
	Type                string                `json:"type"`
	Version             string                `json:"version"`
	Signature           []*felt.Felt          `json:"signature"`
	Nonce               *felt.Felt            `json:"nonce"`
	ContractAddressSalt *felt.Felt            `json:"contract_address_salt"`
	ConstructorCalldata []*felt.Felt          `json:"constructor_calldata"`
	ClassHash           *felt.Felt            `json:"class_hash"`
	ResourceBounds      ResourceBoundsMapping `json:"resource_bounds"`
	Tip                 U64                   `json:"tip"`
	// The data needed to allow the paymaster to pay for the transaction in native tokens
	PayMasterData []*felt.Felt `json:"paymaster_data"`
	// The storage domain of the account's nonce (an account has a nonce per DA mode)
	NonceDataMode DataAvailabilityMode `json:"nonce_data_availability_mode"`
	// The storage domain of the account's balance from which fee will be charged
	FeeMode         DataAvailabilityMode `json:"fee_data_availability_mode"`
	TransactionHash *felt.Felt           `json:"-"`
}

func (tx *DeployAccountTxnV3) SetSignature(x *big.Int, y *big.Int) {
	tx.Signature = []*felt.Felt{BigIntToFelt(x), BigIntToFelt(y)}
}

func (tx *DeployAccountTxnV3) GetTxHash() *big.Int {
	return felt.FeltToBigInt(tx.TransactionHash)
}

func (tx *DeployAccountTxnV3) GetTxRequest() *DeployAccountTxnV3Request {
	calldata := make([]string, len(tx.ConstructorCalldata))
	for i, data := range tx.ConstructorCalldata {
		calldata[i] = data.String()
	}

	payMasterData := make([]string, len(tx.PayMasterData))
	for i, data := range tx.PayMasterData {
		payMasterData[i] = data.String()
	}

	signature := []string{}
	if len(tx.Signature) == 2 {
		signature = append(signature, tx.Signature[0].String())
		signature = append(signature, tx.Signature[1].String())
	}

	nonceDataMode, err := tx.NonceDataMode.UInt64()
	if err != nil {
		return nil
	}
	feeMode, err := tx.FeeMode.UInt64()
	if err != nil {
		return nil
	}
	return &DeployAccountTxnV3Request{
		Type:                tx.Type,
		Version:             tx.Version,
		Signature:           signature,
		Nonce:               tx.Nonce.String(),
		ContractAddressSalt: tx.ContractAddressSalt.String(),
		ConstructorCalldata: calldata,
		ClassHash:           tx.ClassHash.String(),
		ResourceBounds:      tx.ResourceBounds,
		Tip:                 tx.Tip,
		PayMasterData:       payMasterData,
		NonceDataMode:       nonceDataMode,
		FeeMode:             feeMode,
	}
}
func (tx *DeployAccountTxnV3) GetTxRequestJson() string {
	request := tx.GetTxRequest()
	jstr, err := json.Marshal(request)
	if err != nil {
		return ""
	}
	return string(jstr)
}

type DeployAccountTxnV3Request struct {
	Type                string                `json:"type"`
	Version             string                `json:"version"`
	Signature           []string              `json:"signature"`
	Nonce               string                `json:"nonce"`
	ContractAddressSalt string                `json:"contract_address_salt"`
	ConstructorCalldata []string              `json:"constructor_calldata"`
	ClassHash           string                `json:"class_hash"`
	ResourceBounds      ResourceBoundsMapping `json:"resource_bounds"`
	Tip                 U64                   `json:"tip"`
	PayMasterData       []string              `json:"paymaster_data"`
	//AccountDeploymentData []string              `json:"account_deployment_data"`
	NonceDataMode uint64 `json:"nonce_data_availability_mode"`
	FeeMode       uint64 `json:"fee_data_availability_mode"`
}

func (req *DeployAccountTxnV3Request) ToTxn(address, chainId string) (*DeployAccountTxnV3, error) {
	calldata := make([]*felt.Felt, len(req.ConstructorCalldata))
	for i, s := range req.ConstructorCalldata {
		f, err := new(felt.Felt).SetString(s)
		if err != nil {
			return nil, err
		}
		calldata[i] = f
	}

	signature := make([]*felt.Felt, len(req.Signature))
	for i, s := range req.Signature {
		f, err := new(felt.Felt).SetString(s)
		if err != nil {
			return nil, err
		}
		signature[i] = f
	}

	payMasterData := make([]*felt.Felt, len(req.PayMasterData))
	for i, s := range req.PayMasterData {
		f, err := new(felt.Felt).SetString(s)
		if err != nil {
			return nil, err
		}
		payMasterData[i] = f
	}

	nonce, err := new(felt.Felt).SetString(req.Nonce)
	if err != nil {
		return nil, err
	}

	salt, err := new(felt.Felt).SetString(req.ContractAddressSalt)
	if err != nil {
		return nil, err
	}

	classHash, err := new(felt.Felt).SetString(req.ClassHash)
	if err != nil {
		return nil, err
	}

	nonceDataMode, err := getDataAvailabilityMode(req.NonceDataMode)
	if err != nil {
		return nil, err
	}
	feeMode, err := getDataAvailabilityMode(req.FeeMode)
	if err != nil {
		return nil, err
	}

	txn := &DeployAccountTxnV3{
		Type:                req.Type,
		Version:             req.Version,
		Signature:           signature,
		Nonce:               nonce,
		ContractAddressSalt: salt,
		ConstructorCalldata: calldata,
		ClassHash:           classHash,
		ResourceBounds:      req.ResourceBounds,
		Tip:                 req.Tip,
		PayMasterData:       payMasterData,
		NonceDataMode:       nonceDataMode,
		FeeMode:             feeMode,
	}

	txn.TransactionHash, err = TransactionHashDeployAccountV3(txn, felt.HexToFelt(address), new(felt.Felt).SetBytes([]byte(chainId)))
	if err != nil {
		return nil, err
	}

	return txn, nil
}

type ResourceBoundsMapping struct {
	// The max amount and max price per unit of L1 gas used in this tx
	L1Gas ResourceBounds `json:"L1_GAS"`
	// The max amount and max price per unit of L1 blob gas used in this tx
	L1DataGas ResourceBounds `json:"L1_DATA_GAS"`
	// The max amount and max price per unit of L2 gas used in this tx
	L2Gas ResourceBounds `json:"L2_GAS"`
}

type ResourceBounds struct {
	// The max amount of the resource that can be used in the tx
	MaxAmount U64 `json:"max_amount"`
	// The max price per unit of this resource for this tx
	MaxPricePerUnit U128 `json:"max_price_per_unit"`
}

func (rb ResourceBounds) Bytes(resource Resource) ([]byte, error) {
	const eight = 8
	maxAmountBytes := make([]byte, eight)
	maxAmountUint64, err := rb.MaxAmount.ToUint64()
	if err != nil {
		return nil, err
	}
	binary.BigEndian.PutUint64(maxAmountBytes, maxAmountUint64)
	maxPricePerUnitFelt, err := new(felt.Felt).SetString(string(rb.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}
	maxPriceBytes := maxPricePerUnitFelt.Bytes()
	return Flatten(
		[]byte{0},
		[]byte(resource),
		maxAmountBytes,
		maxPriceBytes[16:], // uint128.
	), nil
}

func Flatten[T any](sl ...[]T) []T {
	var result []T
	for _, slice := range sl {
		result = append(result, slice...)
	}

	return result
}

type DataAvailabilityMode string

const (
	DAModeL1 DataAvailabilityMode = "L1"
	DAModeL2 DataAvailabilityMode = "L2"
)

func (da *DataAvailabilityMode) UInt64() (uint64, error) {
	switch *da {
	case DAModeL1:
		return uint64(0), nil
	case DAModeL2:
		return uint64(1), nil
	}
	return 0, errors.New("Unknown DAMode")
}

func getDataAvailabilityMode(mode uint64) (DataAvailabilityMode, error) {
	switch mode {
	case 0:
		return DAModeL1, nil
	case 1:
		return DAModeL2, nil
	}
	return "", errors.New("Unknown DAMode")
}
