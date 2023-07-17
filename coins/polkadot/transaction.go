package polkadot

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"

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

func UnSignedTxFromTxStruct(tx TxStruct, txType int32) UnSignedTx {
	var tp UnSignedTx
	if txType == Transfer {
		tp.Method, _ = BalanceTransfer(tx.ModuleMethod, tx.To, tx.Amount)
	} else if txType == TransferAll {
		tp.Method, _ = BalanceTransferAll(tx.ModuleMethod, tx.To, tx.KeepAlive)
	} else {
		return UnSignedTx{}
	}
	tp.Era = GetEra(tx.BlockHeight, tx.EraHeight)
	if tx.Nonce == 0 {
		tp.Nonce = []byte{0}
	} else {
		nonce := Encode(uint64(tx.Nonce))
		tp.Nonce, _ = hex.DecodeString(nonce)
	}
	if tx.Tip == 0 {
		tp.Tip = []byte{0}
	} else {
		fee := Encode(uint64(tx.Tip))
		tp.Tip, _ = hex.DecodeString(fee)
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
	return tp
}

func UnSignedTxFromTxStruct2(tx TxStruct2) UnSignedTx {
	var tp UnSignedTx
	tp.Method, _ = BalanceTransfer(tx.ModuleMethod, tx.To, tx.Amount)
	tp.Era, _ = hex.DecodeString(tx.Era)
	if tx.Nonce == 0 {
		tp.Nonce = []byte{0}
	} else {
		nonce := Encode(tx.Nonce)
		tp.Nonce, _ = hex.DecodeString(nonce)
	}
	if tx.Tip == 0 {
		tp.Tip = []byte{0}
	} else {
		fee := Encode(tx.Tip)
		tp.Tip, _ = hex.DecodeString(fee)
	}

	tp.SpecVersion, _ = hex.DecodeString(tx.SpecVersion)
	tp.TxVersion, _ = hex.DecodeString(tx.TxVersion)

	genesis := util.RemoveZeroHex(tx.GenesisHash)
	tp.GenesisHash = genesis
	block := util.RemoveZeroHex(tx.BlockHash)
	tp.BlockHash = block
	return tp
}

func SignTx(tx TxStruct, txType int32, privateKey string) string {
	unSignedTx := UnSignedTxFromTxStruct(tx, txType)
	message := unSignedTx.ToBytesString()
	payload, _ := hex.DecodeString(message)
	prikey, _ := hex.DecodeString(privateKey)
	key := ed25519.NewKeyFromSeed(prikey)
	signature, _ := key.Sign(rand.Reader, payload, crypto.Hash(0))

	signed := make([]byte, 0)
	version, _ := hex.DecodeString(tx.Version)
	signed = append(signed, version...)
	signed = append(signed, 0x00)
	from, _ := hex.DecodeString(AddressToPublicKey(tx.From))
	signed = append(signed, from...)
	signed = append(signed, 0x00) // signature type 00:ed25519  01:sr25519  02:ecdsa
	signed = append(signed, signature...)
	signed = append(signed, unSignedTx.Era...)
	signed = append(signed, unSignedTx.Nonce...)
	signed = append(signed, unSignedTx.Tip...)
	signed = append(signed, unSignedTx.Method...)
	lengthBytes, _ := hex.DecodeString(Encode(uint64(len(signed))))
	return "0x" + hex.EncodeToString(lengthBytes) + hex.EncodeToString(signed)
}

func SignTx2(tx TxStruct2, signature []byte) string {
	unSignedTx := UnSignedTxFromTxStruct2(tx)

	signed := make([]byte, 0)
	version, _ := hex.DecodeString(tx.Version)
	signed = append(signed, version...)
	signed = append(signed, 0x00)
	from, _ := hex.DecodeString(AddressToPublicKey(tx.From))
	signed = append(signed, from...)
	signed = append(signed, 0x00) // signature type 00:ed25519  01:sr25519  02:ecdsa
	signed = append(signed, signature...)
	signed = append(signed, unSignedTx.Era...)
	signed = append(signed, unSignedTx.Nonce...)
	signed = append(signed, unSignedTx.Tip...)
	signed = append(signed, unSignedTx.Method...)
	lengthBytes, _ := hex.DecodeString(Encode(uint64(len(signed))))
	return "0x" + hex.EncodeToString(lengthBytes) + hex.EncodeToString(signed)
}
