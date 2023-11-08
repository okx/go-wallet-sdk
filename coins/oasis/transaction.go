package oasis

import (
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/bech32"
	"math/big"
)

// Transfer is a stake transfer.
type Transfer struct {
	To     [21]byte `json:"to"`
	Amount []byte   `json:"amount"`
}

// Transaction is an unsigned consensus transaction.
type Transaction struct {
	// Nonce is a nonce to prevent replay.
	Nonce uint64 `json:"nonce"`
	// Fee is an optional fee that the sender commits to pay to execute this
	// transaction.
	Fee *Fee `json:"fee,omitempty"`

	// Method is the method that should be called.
	Method string `json:"method"`
	// Body is the method call body.
	Body interface{} `json:"body,omitempty"`
}

// Fee is the consensus transaction fee the sender wishes to pay for
// operations which require a fee to be paid to validators.
type Fee struct {
	// Amount is the fee amount to be paid.
	Amount []byte `json:"amount"`
	// Gas is the maximum gas that a transaction can use.
	Gas uint64 `json:"gas"`
}

// SignedTransaction is a signed transaction.
type SignedTransaction struct {
	Signed
}

// Signed is a signed blob.
type Signed struct {
	// Blob is the signed blob.
	Blob []byte `json:"untrusted_raw_value"`

	// Signature is the signature over blob.
	Signature Signature `json:"signature"`
}

// Signature is a signature, bundled with the signing public key.
type Signature struct {
	// PublicKey is the public key that produced the signature.
	PublicKey []byte `json:"public_key"`

	// Signature is the actual raw signature.
	Signature []byte `json:"signature"`
}

func NewTx(nonce, gas uint64, feeAmount *big.Int, body interface{}) *Transaction {
	tx := &Transaction{
		Nonce: nonce,
		Fee: &Fee{
			Amount: feeAmount.Bytes(),
			Gas:    gas,
		},
		Method: "staking.Transfer",
		Body:   body,
	}
	return tx
}

func NewTransferTx(nonce, gas uint64, feeAmount *big.Int, toAddr string, amount *big.Int) *Transaction {
	_, toBytes, _ := bech32.DecodeToBase256(toAddr)
	to := [21]byte{}
	copy(to[:], toBytes)

	transfer := Transfer{
		To:     to,
		Amount: amount.Bytes(),
	}

	return NewTx(nonce, gas, feeAmount, transfer)
}
