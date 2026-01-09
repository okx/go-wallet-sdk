package example

type SignParams struct {
	Type                 int    `json:"txType"`
	ChainId              string `json:"chainId"`
	Nonce                string `json:"nonce"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	GasLimit             string `json:"gasLimit"`
	GasPrice             string `json:"gasPrice"`
	To                   string `json:"to"`
	Value                string `json:"value"`
	Data                 string `json:"data"`
}
