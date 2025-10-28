package types

import (
	"bytes"
	"math/big"

	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/go-ethereum/common"
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/rlp"
)

type SponsoredTx struct {
	ChainID                *big.Int        // destination chain ID
	Nonce                  uint64          // nonce of sender account
	GasTipCap              *big.Int        // maximum tip to the miner
	GasFeeCap              *big.Int        // maximum gas fee want to pay
	Gas                    uint64          // gas limit
	To                     *common.Address `rlp:"nil"` // nil means contract creation
	Value                  *big.Int        // wei amount
	Data                   []byte          // contract invocation input data
	ExpiredTime            uint64          // the expired time of payer's signature
	PayerV, PayerR, PayerS *big.Int        // payer's signature values
	V, R, S                *big.Int        // sender's signature values
}

func (tx *SponsoredTx) copy() TxData {
	cpy := &SponsoredTx{
		Nonce:       tx.Nonce,
		To:          copyAddressPtr(tx.To),
		Data:        common.CopyBytes(tx.Data),
		Gas:         tx.Gas,
		ExpiredTime: tx.ExpiredTime,
		// These are initialized below.
		ChainID:   new(big.Int),
		Value:     new(big.Int),
		GasTipCap: new(big.Int),
		GasFeeCap: new(big.Int),
		PayerV:    new(big.Int),
		PayerR:    new(big.Int),
		PayerS:    new(big.Int),
		V:         new(big.Int),
		R:         new(big.Int),
		S:         new(big.Int),
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.GasTipCap != nil {
		cpy.GasTipCap.Set(tx.GasTipCap)
	}
	if tx.GasFeeCap != nil {
		cpy.GasFeeCap.Set(tx.GasFeeCap)
	}
	if tx.PayerV != nil {
		cpy.PayerV.Set(tx.PayerV)
	}
	if tx.PayerR != nil {
		cpy.PayerR.Set(tx.PayerR)
	}
	if tx.PayerS != nil {
		cpy.PayerS.Set(tx.PayerS)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *SponsoredTx) txType() byte           { return SponsoredTxType }
func (tx *SponsoredTx) chainID() *big.Int      { return tx.ChainID }
func (tx *SponsoredTx) accessList() AccessList { return nil }
func (tx *SponsoredTx) data() []byte           { return tx.Data }
func (tx *SponsoredTx) gas() uint64            { return tx.Gas }
func (tx *SponsoredTx) gasPrice() *big.Int     { return tx.GasFeeCap }
func (tx *SponsoredTx) gasTipCap() *big.Int    { return tx.GasTipCap }
func (tx *SponsoredTx) gasFeeCap() *big.Int    { return tx.GasFeeCap }
func (tx *SponsoredTx) value() *big.Int        { return tx.Value }
func (tx *SponsoredTx) nonce() uint64          { return tx.Nonce }
func (tx *SponsoredTx) to() *common.Address    { return tx.To }
func (tx *SponsoredTx) expiredTime() uint64    { return tx.ExpiredTime }

func (tx *SponsoredTx) rawPayerSignatureValues() (v, r, s *big.Int) {
	return tx.PayerV, tx.PayerR, tx.PayerS
}

func (tx *SponsoredTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *SponsoredTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}

func (tx *SponsoredTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *SponsoredTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
