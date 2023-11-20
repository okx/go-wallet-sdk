/**
Author： https://github.com/hecodev007/block_sign
*/

package transactions

import (
	"github.com/okx/go-wallet-sdk/coins/helium/keypair"
)

type PaymentV1Tx struct {
	Payer  *keypair.Addressable
	Payee  *keypair.Addressable
	Amount uint64
	Fee    uint64
	Nonce  uint64
	Sig    []byte
}

func NewPaymentV1Tx(from, to *keypair.Addressable, amount, fee, nonce uint64, sig []byte) *PaymentV1Tx {
	return &PaymentV1Tx{
		Payer:  from,
		Payee:  to,
		Amount: amount,
		Fee:    fee,
		Nonce:  nonce,
		Sig:    sig,
	}
}

func (v1 *PaymentV1Tx) SetFee(fee uint64) {
	v1.Fee = fee
}

func (v1 *PaymentV1Tx) SetSignature(sig []byte) {
	v1.Sig = sig
	return
}
