package common

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/sha3"
	"regexp"
)

const CHAIN_ID_MAINNET_ENDLESS = 220
const CHAIN_ID_TESTNET_ENDLESS = 221

func IsEndless(chainId uint8) bool {
	if chainId == CHAIN_ID_MAINNET_ENDLESS || chainId == CHAIN_ID_TESTNET_ENDLESS {
		return true
	} else {
		return false
	}
}

func ComputeTransactionHash(prefix []byte, hexStr string) string {
	bcsBytes, _ := hex.DecodeString(hexStr)
	// Serialize Transaction structure, so need to add sequence number
	message := make([]byte, 0)
	message = append(message, prefix...)
	message = append(message, 0x0)
	message = append(message, bcsBytes...)
	return "0x" + hex.EncodeToString(Sha256Hash(message))
}

func Sha256Hash(bytes []byte) []byte {
	sha256 := sha3.New256()
	sha256.Write(bytes)
	return sha256.Sum(nil)
}

func IsHexString(s string) bool {
	res, err := regexp.MatchString("0x[0-9a-fA-F]+", s)
	if err != nil {
		return false
	}
	return res
}

func ConvertInterfaceToByte(arg []interface{}) (res []byte, err error) {
	res = make([]byte, 0)
	for _, v := range arg {
		val, ok := v.(byte)
		if !ok {
			err = errors.New("value is not byte")
			return
		}
		res = append(res, val)
	}
	return
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
