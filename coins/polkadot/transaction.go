package polkadot

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"golang.org/x/crypto/blake2b"

	"github.com/okx/go-wallet-sdk/util"
)

type TxStruct struct {
	ModuleMethod string `json:"module_method"`
	Version      string `json:"version"`
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       uint64 `json:"amount"`
	Nonce        uint64 `json:"nonce"`
	Tip          uint64 `json:"tip"`
	BlockHeight  uint64 `json:"block_height"`
	BlockHash    string `json:"block_hash"`
	GenesisHash  string `json:"genesis_hash"`
	SpecVersion  uint32 `json:"spec_version"`
	TxVersion    uint32 `json:"tx_version"`
	KeepAlive    string `json:"keep_alive"`
	EraHeight    uint64 `json:"era_height"`
}

type TxStruct2 struct {
	ModuleMethod string `json:"module_method"`
	Version      string `json:"version"`
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       uint64 `json:"amount"`
	Nonce        uint64 `json:"nonce"`
	Tip          uint64 `json:"tip"`
	BlockHash    string `json:"block_hash"`
	GenesisHash  string `json:"genesis_hash"`
	SpecVersion  string `json:"spec_version"`
	TxVersion    string `json:"tx_version"`
	Era          string `json:"era"`
}

type UnSignedTx struct {
	Method      []byte
	Era         []byte
	Nonce       []byte
	Tip         []byte
	SpecVersion []byte
	GenesisHash []byte
	BlockHash   []byte
	TxVersion   []byte
}

func (t *UnSignedTx) ToBytesString() string {
	payload := make([]byte, 0)
	payload = append(payload, t.Method...)
	payload = append(payload, t.Era...)
	payload = append(payload, t.Nonce...)
	payload = append(payload, t.Tip...)
	payload = append(payload, t.SpecVersion...)
	payload = append(payload, t.TxVersion...)
	payload = append(payload, t.GenesisHash...)
	payload = append(payload, t.BlockHash...)
	return hex.EncodeToString(payload)
}

func UnSignedTxFromTxStruct(tx TxStruct, txType int32) (UnSignedTx, error) {
	var tp UnSignedTx
	var err error
	if txType == Transfer {
		tp.Method, err = BalanceTransfer(tx.ModuleMethod, tx.To, tx.Amount)
		if err != nil {
			return UnSignedTx{}, err
		}
	} else if txType == TransferAll {
		tp.Method, err = BalanceTransferAll(tx.ModuleMethod, tx.To, tx.KeepAlive)
		if err != nil {
			return UnSignedTx{}, err
		}
	} else {
		return UnSignedTx{}, nil
	}
	tp.Era = GetEra(tx.BlockHeight, tx.EraHeight)
	if tx.Nonce == 0 {
		tp.Nonce = []byte{0}
	} else {
		nonce := Encode(uint64(tx.Nonce))
		tp.Nonce, err = hex.DecodeString(nonce)
		if err != nil {
			return UnSignedTx{}, err
		}
	}
	if tx.Tip == 0 {
		tp.Tip = []byte{0}
	} else {
		fee := Encode(uint64(tx.Tip))
		tp.Tip, err = hex.DecodeString(fee)
		if err != nil {
			return UnSignedTx{}, err
		}
	}

	specv := make([]byte, 4)
	binary.LittleEndian.PutUint32(specv, tx.SpecVersion)
	tp.SpecVersion = specv

	txv := make([]byte, 4)
	binary.LittleEndian.PutUint32(txv, tx.TxVersion)
	tp.TxVersion = txv

	genesis := util.RemoveZeroHex(tx.GenesisHash)
	tp.GenesisHash = genesis
	block := util.RemoveZeroHex(tx.BlockHash)
	tp.BlockHash = block
	return tp, nil
}

func UnSignedTxFromTxStruct2(tx TxStruct2) (UnSignedTx, error) {
	var tp UnSignedTx
	var err error
	tp.Method, err = BalanceTransfer(tx.ModuleMethod, tx.To, tx.Amount)
	if err != nil {
		return UnSignedTx{}, err
	}
	tp.Era, err = hex.DecodeString(tx.Era)
	if err != nil {
		return UnSignedTx{}, err
	}
	if tx.Nonce == 0 {
		tp.Nonce = []byte{0}
	} else {
		nonce := Encode(tx.Nonce)
		tp.Nonce, err = hex.DecodeString(nonce)
		if err != nil {
			return UnSignedTx{}, err
		}
	}
	if tx.Tip == 0 {
		tp.Tip = []byte{0}
	} else {
		fee := Encode(tx.Tip)
		tp.Tip, err = hex.DecodeString(fee)
		if err != nil {
			return UnSignedTx{}, err
		}
	}

	tp.SpecVersion, err = hex.DecodeString(tx.SpecVersion)
	if err != nil {
		return UnSignedTx{}, err
	}
	tp.TxVersion, err = hex.DecodeString(tx.TxVersion)
	if err != nil {
		return UnSignedTx{}, err
	}

	genesis := util.RemoveZeroHex(tx.GenesisHash)
	tp.GenesisHash = genesis
	block := util.RemoveZeroHex(tx.BlockHash)
	tp.BlockHash = block
	return tp, nil
}

func SignTx(tx TxStruct, txType int32, privateKey string) (string, error) {
	unSignedTx, err := UnSignedTxFromTxStruct(tx, txType)
	if err != nil {
		return "", err
	}
	message := unSignedTx.ToBytesString()
	payload, err := hex.DecodeString(message)
	if err != nil {
		return "", err
	}
	prikey, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	key := ed25519.NewKeyFromSeed(prikey)
	signature, err := key.Sign(rand.Reader, payload, crypto.Hash(0))
	if err != nil {
		return "", err
	}
	signed := make([]byte, 0)
	version, err := hex.DecodeString(tx.Version)
	if err != nil {
		return "", err
	}
	signed = append(signed, version...)
	signed = append(signed, 0x00)
	pubKey, err := AddressToPublicKey(tx.From)
	if err != nil {
		return "", err
	}
	from, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", err
	}
	signed = append(signed, from...)
	signed = append(signed, 0x00) // signature type 00:ed25519  01:sr25519  02:ecdsa
	signed = append(signed, signature...)
	signed = append(signed, unSignedTx.Era...)
	signed = append(signed, unSignedTx.Nonce...)
	signed = append(signed, unSignedTx.Tip...)
	signed = append(signed, unSignedTx.Method...)
	lengthBytes, err := hex.DecodeString(Encode(uint64(len(signed))))
	if err != nil {
		return "", err
	}
	return util.EncodeHexWith0x(lengthBytes) + hex.EncodeToString(signed), nil
}

func SignTx2(tx TxStruct2, signature []byte) (string, error) {
	unSignedTx, err := UnSignedTxFromTxStruct2(tx)
	if err != nil {
		return "", err
	}

	signed := make([]byte, 0)
	version, err := hex.DecodeString(tx.Version)
	if err != nil {
		return "", err
	}
	signed = append(signed, version...)
	signed = append(signed, 0x00)
	pubKey, err := AddressToPublicKey(tx.From)
	if err != nil {
		return "", err
	}
	from, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", err
	}
	signed = append(signed, from...)
	signed = append(signed, 0x00) // signature type 00:ed25519  01:sr25519  02:ecdsa
	signed = append(signed, signature...)
	signed = append(signed, unSignedTx.Era...)
	signed = append(signed, unSignedTx.Nonce...)
	signed = append(signed, unSignedTx.Tip...)
	signed = append(signed, unSignedTx.Method...)
	lengthBytes, err := hex.DecodeString(Encode(uint64(len(signed))))
	if err != nil {
		return "", err
	}
	return util.EncodeHexWith0x(lengthBytes) + hex.EncodeToString(signed), nil
}

func CalTxHash(signedTx string) (string, error) {
	if signedTx[:2] == "0x" || signedTx[:2] == "0X" {
		signedTx = signedTx[2:]
	}
	b, err := hex.DecodeString(signedTx)
	if err != nil {
		return "", err
	}
	h, err := blake2b.New256(nil)
	if err != nil {
		return "", err
	}
	_, err = h.Write(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
