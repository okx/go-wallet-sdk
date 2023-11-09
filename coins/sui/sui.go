package sui

import (
	crypto_ed25519 "crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/dchest/blake2b"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/okx/go-wallet-sdk/crypto/ed25519"
	"regexp"
	"strings"
)

const (
	SignatureScheme    = "ED25519"
	PUBLIC_KEY_SIZE    = 32
	SUI_ADDRESS_LENGTH = 32
)

var (
	ErrInvalidSuiRequest = errors.New("invalid sui request")
	ErrInvalidSuiSeedHex = errors.New("invalid sui seed hex")
	ErrInvalidSuiParam   = errors.New("invalid sui param")
	ErrUnknownSuiRequest = errors.New("unknown sui request")
)

func ValidateAddress(address string) bool {
	if strings.HasPrefix(address, "0x") {
		re1, err := regexp.Compile("^0x[\\dA-Fa-f]{64}$")
		if err != nil {
			panic(err)
		}
		return re1.Match([]byte(address))
	}
	re2, err := regexp.Compile("^[\\dA-Fa-f]{64}$")
	if err != nil {
		panic(err)
	}
	return re2.Match([]byte(address))
}

func NormalizeSuiAddress(value string) string {
	v := strings.ToLower(value)
	if strings.HasPrefix(v, "0x") {
		if len(v) == SUI_ADDRESS_LENGTH*2+2 {
			return v
		}
		return "0x" + strings.Repeat("0", SUI_ADDRESS_LENGTH*2+2-len(v)) + v[2:]
	}
	return "0x" + strings.Repeat("0", SUI_ADDRESS_LENGTH*2-len(v)) + v
}

func Hash(txBytes string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(txBytes)
	if err != nil {
		return "", err
	}
	data := make([]byte, len("TransactionData::")+len(b))
	copy(data, "TransactionData::")
	copy(data[len("TransactionData::"):], b)
	hash := blake2b.New256()
	hash.Write(data)
	result := hash.Sum(nil)
	return base58.Encode(result), nil
}

func GenerateKey() (crypto_ed25519.PrivateKey, error) {
	return ed25519.GenerateKey()
}

func NewAddress(seedHex string) string {
	publicKey, err := ed25519.PublicKeyFromSeed(seedHex)
	if err != nil {
		return ""
	}
	if len(publicKey) != PUBLIC_KEY_SIZE {
		return ""
	}
	k := make([]byte, PUBLIC_KEY_SIZE+1)
	copy(k[1:], publicKey)
	hash := blake2b.New256()
	hash.Write(k)
	h := hash.Sum(nil)
	address := "0x" + hex.EncodeToString(h)[0:64]
	return address
}

func GetAddressByPubKey(publicKeyHex string) (string, error) {
	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", err
	}
	if len(publicKey) != PUBLIC_KEY_SIZE {
		return "", nil
	}
	k := make([]byte, PUBLIC_KEY_SIZE+1)
	copy(k[1:], publicKey)
	hash := blake2b.New256()
	hash.Write(k)
	h := hash.Sum(nil)
	address := "0x" + hex.EncodeToString(h)[0:64]
	return address, nil
}

type Request struct {
	Data string `json:"data"`
	Type Type   `json:"type"`
}

func GetTxHash(r *Request, to string, gasBudget uint64, gasPrice uint64, addr string) (string, error) {
	txBytes, err := PrepareTx(r, to, gasBudget, gasPrice, addr)
	if err != nil {
		return "", err
	}

	if len(txBytes) == 0 {
		return "", errors.New("err txBytes")
	}

	b, err := GetRawTx(txBytes)
	if err != nil {
		return "", errors.New("get raw tx err")
	}
	return b, nil
}

func GetRawTx(txBytes string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(txBytes)
	if err != nil {
		return "", errors.New("decode txBytes error")
	}
	signDo := make([]byte, len(data)+3)
	copy(signDo[3:], data)
	hash := blake2b.New256()
	hash.Write(signDo)
	b := hash.Sum(nil)
	return hex.EncodeToString(b), err
}

func PrepareTx(r *Request, to string, gasBudget uint64, gasPrice uint64, addr string) (string, error) {
	var s []byte
	var err error
	switch {
	case r.Type == Empty || r.Type == Transfer:
		var req PaySuiRequest
		if err = json.Unmarshal([]byte(r.Data), &req); err != nil {
			return "", ErrInvalidSuiRequest
		}
		if len(req.Coins) == 0 {
			return "", errors.New("invalid sui request")
		}
		s, err = BuildTx(addr, to, req.Coins, req.Amount, req.Epoch, gasBudget, gasPrice)
		if err != nil {
			return "", err
		}
	case r.Type == Split:
		var req SplitSuiRequest
		if err = json.Unmarshal([]byte(r.Data), &req); err != nil {
			return "", ErrInvalidSuiRequest
		}
		if len(req.Coins) == 0 {
			return "", errors.New("invalid sui request")
		}
		s, err = BuildSplitTx(addr, to, req.Coins, req.Amounts, req.Epoch, gasBudget, gasPrice)
		if err != nil {
			return "", err
		}
	case r.Type == Stake:
		var req StakeSuiRequest
		if err = json.Unmarshal([]byte(r.Data), &req); err != nil {
			return "", ErrInvalidSuiRequest
		}
		if len(req.Coins) == 0 {
			return "", errors.New("invalid sui request")
		}
		s, err = BuildStakeTx(addr, to, req.Coins, req.Amount, req.Epoch, gasBudget, gasPrice)
		if err != nil {
			return "", err
		}
	case r.Type == WithdrawStake:
		var req WithdrawStakSuiRequest
		if err = json.Unmarshal([]byte(r.Data), &req); err != nil {
			return "", ErrInvalidSuiRequest
		}
		if len(req.Coins) == 0 {
			return "", errors.New("invalid sui request")
		}
		s, err = BuildWithdrawStakeTx(addr, req.Coins, req.StakeCoin, req.Epoch, gasBudget, gasPrice)
		if err != nil {
			return "", err
		}
	case r.Type == Merge:
		var req MergeSuiRequest
		if err = json.Unmarshal([]byte(r.Data), &req); err != nil {
			return "", ErrInvalidSuiRequest
		}
		if len(req.Coins) == 0 || len(req.Objects) == 0 {
			return "", errors.New("invalid sui request")
		}
		s, err = BuildMergeTx(addr, req.Coins, req.Objects, req.Epoch, gasBudget, gasPrice)
		if err != nil {
			return "", err
		}
	}

	if len(s) == 0 {
		return "", ErrInvalidSuiParam
	}
	return base64.StdEncoding.EncodeToString(s), nil
}

func Execute(r *Request, from, to string, gasBudget uint64, gasPrice uint64, seedHex string) (string, error) {
	addr := NewAddress(seedHex)
	if len(addr) == 0 || from != addr {
		return "", ErrInvalidSuiSeedHex
	}
	if r == nil {
		return "", errors.New("invalid sui request")
	}
	if to == "" {
		return "", errors.New("invalid sui to")
	}
	if gasBudget == 0 {
		return "", errors.New("invalid sui gasBudget")
	}
	if gasPrice == 0 {
		return "", errors.New("invalid sui gasPrice")
	}
	raw, err := PrepareTx(r, to, gasBudget, gasPrice, addr)
	if err != nil {
		return "", err
	}
	tx, err := SignTransaction(raw, seedHex)
	if err != nil {
		return "", err
	}
	if tx == nil {
		return "", ErrInvalidSuiSeedHex
	}
	b, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func SignMessage(input string, seedHex string) (string, error) {
	signDo := make([]byte, len(input)+3)
	signDo[0] = 3
	copy(signDo[3:], []byte(input))
	hash := blake2b.New256()
	hash.Write(signDo)
	b := hash.Sum(nil)
	signature, err := ed25519.Sign(seedHex, b)
	if err != nil {
		return "", err
	}
	publicKey, err := ed25519.PublicKeyFromSeed(seedHex)
	if err != nil {
		return "", err
	}
	sign := make([]byte, 1+len(publicKey)+len(signature))
	copy(sign[1:], signature)
	copy(sign[1+len(signature):], publicKey)
	return base64.StdEncoding.EncodeToString(sign), nil
}

func SignTransaction(txBytes string, seedHex string) (*SignedTransaction, error) {
	if len(txBytes) == 0 || len(seedHex) == 0 {
		return nil, nil
	}
	publicKey, err := ed25519.PublicKeyFromSeed(seedHex)
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(txBytes)
	if err != nil {
		return nil, err
	}
	signDo := make([]byte, len(data)+3)
	copy(signDo[3:], data)
	hash := blake2b.New256()
	hash.Write(signDo)
	b := hash.Sum(nil)
	signature, err := ed25519.Sign(seedHex, b)
	if err != nil {
		return nil, err
	}
	sign := make([]byte, 1+len(publicKey)+len(signature))
	copy(sign[1:], signature)
	copy(sign[1+len(signature):], publicKey)
	return &SignedTransaction{TxBytes: txBytes, Signature: base64.StdEncoding.EncodeToString(sign)}, nil
}

type PaySuiRequest struct {
	Coins  []*SuiObjectRef `json:"coins"`
	Amount uint64          `json:"amount"`
	Epoch  uint64          `json:"epoch"`
}

type SplitSuiRequest struct {
	Coins   []*SuiObjectRef `json:"coins"`
	Amounts []uint64        `json:"amounts"`
	Epoch   uint64          `json:"epoch"`
}

type StakeSuiRequest struct {
	Coins  []*SuiObjectRef `json:"coins"`
	Amount uint64          `json:"amount"`
	Epoch  uint64          `json:"epoch"`
}

type WithdrawStakSuiRequest struct {
	Coins     []*SuiObjectRef `json:"coins"`
	StakeCoin *SuiObjectRef   `json:"stake_coin"`
	Epoch     uint64          `json:"epoch"`
}

type MergeSuiRequest struct {
	Coins   []*SuiObjectRef `json:"coins"`
	Objects []*SuiObjectRef `json:"objects"`
	Epoch   uint64          `json:"epoch"`
}
