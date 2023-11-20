package elrond

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"golang.org/x/crypto/sha3"
)

type TxBuilder struct {
	privateKey *ed25519.PrivateKey
}

// ArgCreateTransaction will hold the transaction fields
type ArgCreateTransaction struct {
	Nonce    uint64
	Value    string
	RcvAddr  string
	GasPrice uint64
	GasLimit uint64
	// Arbitrary information about the transaction, base64-encoded.
	Data []byte
	// The chain identifier.
	ChainID string
	Version uint32
	Options uint32
}

type Transaction struct {
	Nonce    uint64 `json:"nonce"`
	Value    string `json:"value"`
	RcvAddr  string `json:"receiver"`
	SndAddr  string `json:"sender"`
	GasPrice uint64 `json:"gasPrice,omitempty"`
	GasLimit uint64 `json:"gasLimit,omitempty"`
	Data     []byte `json:"data,omitempty"`
	ChainID  string `json:"chainID"`
	Version  uint32 `json:"version"`
	Options  uint32 `json:"options,omitempty"`
	// The digital signature consisting of 128 hex-characters (thus 64 bytes in a raw representation)
	Signature string `json:"signature,omitempty"`
}

// NewTxBuilder will create a new transaction builder able to build and correctly sign a transaction
func NewTxBuilder(key *ed25519.PrivateKey) *TxBuilder {
	return &TxBuilder{
		privateKey: key,
	}
}

// createTransaction assembles a transaction from the provided arguments
func (builder *TxBuilder) createTransaction(arg ArgCreateTransaction, sender string, signature string) *Transaction {
	return &Transaction{
		Nonce:     arg.Nonce,
		Value:     arg.Value,
		SndAddr:   sender,
		RcvAddr:   arg.RcvAddr,
		GasPrice:  arg.GasPrice,
		GasLimit:  arg.GasLimit,
		Data:      arg.Data,
		ChainID:   arg.ChainID,
		Version:   arg.Version,
		Options:   arg.Options,
		Signature: signature,
	}
}

func (builder *TxBuilder) createUnsignedMessage(arg ArgCreateTransaction, sender string) ([]byte, error) {
	tx := builder.createTransaction(arg, sender, "")
	return json.Marshal(tx)
}

// ApplySignatureAndGenerateTx will apply the corresponding sender and compute the signature field and
// generate the transaction instance
func (builder *TxBuilder) build(arg ArgCreateTransaction) (*Transaction, error) {
	pkBytes := builder.privateKey.Public().(ed25519.PublicKey)
	sender, _ := bech32.EncodeFromBase256(HRP, pkBytes)
	unsignedMessage, err := builder.createUnsignedMessage(arg, sender)
	if err != nil {
		return nil, err
	}
	shouldSignOnTxHash := arg.Version >= 2 && arg.Options&1 > 0
	if shouldSignOnTxHash {
		hasher := sha3.NewLegacyKeccak256()
		hasher.Write(unsignedMessage[:])
		unsignedMessage = hasher.Sum(nil)
	}
	signature, err := builder.privateKey.Sign(crypto_rand.Reader, unsignedMessage, crypto.Hash(0))
	if err != nil {
		return nil, err
	}
	return builder.createTransaction(arg, sender, hex.EncodeToString(signature)), nil
}
