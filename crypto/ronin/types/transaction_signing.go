// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/go-ethereum/common"
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/go-ethereum/crypto"
)

var (
	ErrInvalidChainId    = errors.New("invalid chain id for signer")
	errMissingPayerField = errors.New("transaction has no payer field")
)

// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache struct {
	signer Signer
	from   common.Address
}

// SignTx signs the transaction using the given signer and private key.
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

// SignNewTx creates a transaction and signs it.
func SignNewTx(prv *ecdsa.PrivateKey, s Signer, txdata TxData) (*Transaction, error) {
	tx := NewTx(txdata)
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

func PayerSign(prv *ecdsa.PrivateKey, signer Signer, sender common.Address, txdata TxData) (r, s, v *big.Int, err error) {
	payerHash := rlpHash([]interface{}{
		signer.ChainID(),
		sender,
		txdata.nonce(),
		txdata.gasTipCap(),
		txdata.gasFeeCap(),
		txdata.gas(),
		txdata.to(),
		txdata.value(),
		txdata.data(),
		txdata.expiredTime(),
	})

	sig, err := crypto.Sign(payerHash[:], prv)
	if err != nil {
		return nil, nil, nil, err
	}

	r, s, _ = decodeSignature(sig)
	v = big.NewInt(int64(sig[64]))
	return r, s, v, nil
}

// MustSignNewTx creates a transaction and signs it.
// This panics if the transaction cannot be signed.
func MustSignNewTx(prv *ecdsa.PrivateKey, s Signer, txdata TxData) *Transaction {
	tx, err := SignNewTx(prv, s, txdata)
	if err != nil {
		panic(err)
	}
	return tx
}

// Sender returns the address derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// Sender may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func Sender(signer Signer, tx *Transaction) (common.Address, error) {
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.from.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

// Payer returns the address derived from payer's signature in sponsored
// transaction or nil in other transaction types.
//
// Payer may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func Payer(signer Signer, tx *Transaction) (common.Address, error) {
	if tx.Type() != SponsoredTxType {
		return common.Address{}, errMissingPayerField
	}

	if sc := tx.payer.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Payer(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.payer.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

// Signer encapsulates transaction signature handling. The name of this type is slightly
// misleading because Signers don't actually sign, they're just for validating and
// processing of signatures.
//
// Note that this interface is not a stable API and may change at any time to accommodate
// new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction) (common.Address, error)

	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	ChainID() *big.Int

	// Hash returns 'signature hash', i.e. the transaction hash that is signed by the
	// private key. This hash does not uniquely identify the transaction.
	Hash(tx *Transaction) common.Hash

	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool

	// Payer returns the payer address of sponsored transaction
	Payer(tx *Transaction) (common.Address, error)
}

type MikoSigner struct {
	EIP155Signer
}

func NewMikoSigner(chainId *big.Int) MikoSigner {
	return MikoSigner{NewEIP155Signer(chainId)}
}

func (s MikoSigner) Equal(s2 Signer) bool {
	miko, ok := s2.(MikoSigner)
	return ok && miko.chainId.Cmp(s.chainId) == 0
}

func (s MikoSigner) Sender(tx *Transaction) (common.Address, error) {
	switch tx.Type() {
	case LegacyTxType:
		return s.EIP155Signer.Sender(tx)
	case SponsoredTxType:
		if tx.ChainId().Cmp(s.chainId) != 0 {
			return common.Address{}, ErrInvalidChainId
		}
		// V in sponsored signature is {0, 1}, but the recoverPlain expects
		// {0, 1} + 27, so we need to add 27 to V
		V, R, S := tx.RawSignatureValues()
		V = new(big.Int).Add(V, big.NewInt(27))
		return recoverPlain(s.Hash(tx), R, S, V, true)
	default:
		return common.Address{}, ErrTxTypeNotSupported
	}
}

func (s MikoSigner) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	switch tx.Type() {
	case LegacyTxType:
		return s.EIP155Signer.SignatureValues(tx, sig)
	case SponsoredTxType:
		// V in sponsored signature is {0, 1}, get it directly from raw signature
		// because decodeSignature returns {0, 1} + 27
		R, S, _ := decodeSignature(sig)
		V := big.NewInt(int64(sig[64]))
		return R, S, V, nil
	default:
		return nil, nil, nil, ErrTxTypeNotSupported
	}
}

func (s MikoSigner) Hash(tx *Transaction) common.Hash {
	switch tx.Type() {
	case LegacyTxType:
		return s.EIP155Signer.Hash(tx)
	case SponsoredTxType:
		payerV, payerR, payerS := tx.RawPayerSignatureValues()
		return prefixedRlpHash(
			tx.Type(),
			[]interface{}{
				s.chainId,
				tx.Nonce(),
				tx.GasTipCap(),
				tx.GasFeeCap(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				tx.ExpiredTime(),
				payerV, payerR, payerS,
			},
		)
	default:
		return common.Hash{}
	}
}

func payerInternal(s Signer, tx *Transaction) (common.Address, error) {
	if tx.Type() != SponsoredTxType {
		return common.Address{}, ErrInvalidTxType
	}

	sender, err := Sender(s, tx)
	if err != nil {
		return common.Address{}, err
	}

	payerV, payerR, payerS := tx.RawPayerSignatureValues()
	payerHash := rlpHash([]interface{}{
		tx.ChainId(), // The chainId is checked in Sender already
		sender,
		tx.Nonce(),
		tx.GasTipCap(),
		tx.GasFeeCap(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		tx.ExpiredTime(),
	})

	// V in payer signature is {0, 1}, but the recoverPlain expects
	// {0, 1} + 27, so we need to add 27 to V
	payerV = new(big.Int).Add(payerV, big.NewInt(27))
	return recoverPlain(payerHash, payerR, payerS, payerV, true)
}

func (s MikoSigner) Payer(tx *Transaction) (common.Address, error) {
	return payerInternal(s, tx)
}

// EIP155Signer implements Signer using the EIP-155 rules. This accepts transactions which
// are replay-protected as well as unprotected homestead transactions.
type EIP155Signer struct {
	chainId, chainIdMul *big.Int
}

func NewEIP155Signer(chainId *big.Int) EIP155Signer {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return EIP155Signer{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}

func (s EIP155Signer) ChainID() *big.Int {
	return s.chainId
}

func (s EIP155Signer) Equal(s2 Signer) bool {
	eip155, ok := s2.(EIP155Signer)
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}

var big8 = big.NewInt(8)

func (s EIP155Signer) Sender(tx *Transaction) (common.Address, error) {
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	if !tx.Protected() {
		return HomesteadSigner{}.Sender(tx)
	}
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}
	V, R, S := tx.RawSignatureValues()
	V = new(big.Int).Sub(V, s.chainIdMul)
	V.Sub(V, big8)
	return recoverPlain(s.Hash(tx), R, S, V, true)
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s EIP155Signer) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	if tx.Type() != LegacyTxType {
		return nil, nil, nil, ErrTxTypeNotSupported
	}
	R, S, V = decodeSignature(sig)
	if s.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, s.chainIdMul)
	}
	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s EIP155Signer) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		s.chainId, uint(0), uint(0),
	})
}

func (s EIP155Signer) Payer(tx *Transaction) (common.Address, error) {
	return common.Address{}, ErrInvalidTxType
}

// HomesteadTransaction implements TransactionInterface using the
// homestead rules.
type HomesteadSigner struct{ FrontierSigner }

func (s HomesteadSigner) ChainID() *big.Int {
	return nil
}

func (s HomesteadSigner) Equal(s2 Signer) bool {
	_, ok := s2.(HomesteadSigner)
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (hs HomesteadSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	return hs.FrontierSigner.SignatureValues(tx, sig)
}

func (hs HomesteadSigner) Sender(tx *Transaction) (common.Address, error) {
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	v, r, s := tx.RawSignatureValues()
	return recoverPlain(hs.Hash(tx), r, s, v, true)
}

type FrontierSigner struct{}

func (s FrontierSigner) ChainID() *big.Int {
	return nil
}

func (s FrontierSigner) Equal(s2 Signer) bool {
	_, ok := s2.(FrontierSigner)
	return ok
}

func (fs FrontierSigner) Sender(tx *Transaction) (common.Address, error) {
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	v, r, s := tx.RawSignatureValues()
	return recoverPlain(fs.Hash(tx), r, s, v, false)
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (fs FrontierSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	if tx.Type() != LegacyTxType {
		return nil, nil, nil, ErrTxTypeNotSupported
	}
	r, s, v = decodeSignature(sig)
	return r, s, v, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (fs FrontierSigner) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
	})
}

func (fs FrontierSigner) Payer(tx *Transaction) (common.Address, error) {
	return common.Address{}, ErrInvalidTxType
}

func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		return common.Address{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return common.Address{}, ErrInvalidSig
	}
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
