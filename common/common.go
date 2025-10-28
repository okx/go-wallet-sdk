package common

import "encoding/json"

type SignTxResp struct {
	TxHash string `json:"txHash"`
	Tx     string `json:"tx"`
}

func (b SignTxResp) ToJson() (string, error) {
	res, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func FormatSignTxResp(txHash, tx string) (string, error) {
	r := SignTxResp{TxHash: txHash, Tx: tx}
	return r.ToJson()
}

type SignMultiAgentTxResp struct {
	RawTxn           string `json:"rawTxn"`
	AccAuthenticator string `json:"accAuthenticator"`
}

func FormatSignMultiAgentTxResp(rawTxn, accAuthenticator string) (string, error) {
	resp := SignMultiAgentTxResp{
		RawTxn:           rawTxn,
		AccAuthenticator: accAuthenticator,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type SignSimpleTxResp struct {
	TxHash           string `json:"txHash"`
	RawTxn           string `json:"rawTxn"`
	AccAuthenticator string `json:"accAuthenticator"`
}

func FormatSimpleTxRespTxResp(txHash, rawTxn, accAuthenticator string) (string, error) {
	resp := SignSimpleTxResp{
		TxHash:           txHash,
		RawTxn:           rawTxn,
		AccAuthenticator: accAuthenticator,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
