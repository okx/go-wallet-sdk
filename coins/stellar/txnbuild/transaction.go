/*
Package txnbuild implements transactions and operations on the Stellar network.
This library provides an interface to the Stellar transaction model. It supports the building of Go applications on
top of the Stellar network (https://www.stellar.org/). Transactions constructed by this library may be submitted
to any Horizon instance for processing onto the ledger, using any Stellar SDK client. The recommended client for Go
programmers is horizonclient (https://github.com/stellar/go/tree/master/clients/horizonclient). Together, these two
libraries provide a complete Stellar SDK.
For more information and further examples, see https://github.com/stellar/go/blob/master/docs/reference/readme.md
*/
package txnbuild

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/stellar/keypair"
	"github.com/okx/go-wallet-sdk/coins/stellar/network"
	"github.com/okx/go-wallet-sdk/coins/stellar/support/collections/set"
	"github.com/okx/go-wallet-sdk/coins/stellar/support/errors"
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
	"math"
	"math/bits"
)

// MinBaseFee is the minimum transaction fee for the Stellar network of 100 stroops (0.00001 XLM).
const MinBaseFee = 100

// Account represents the aspects of a Stellar account necessary to construct transactions. See
// https://developers.stellar.org/docs/glossary/accounts/
type Account interface {
	GetAccountID() string
	IncrementSequenceNumber() (int64, error)
	GetSequenceNumber() (int64, error)
}

func hashHex(e xdr.TransactionEnvelope, networkStr string) (string, error) {
	h, err := network.HashTransactionInEnvelope(e, networkStr)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h[:]), nil
}

func concatSignatures(
	e xdr.TransactionEnvelope,
	networkStr string,
	signatures []xdr.DecoratedSignature,
	kps ...*keypair.Full,
) ([]xdr.DecoratedSignature, error) {
	// Hash the transaction
	h, err := network.HashTransactionInEnvelope(e, networkStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash transaction")
	}

	extended := make(
		[]xdr.DecoratedSignature,
		len(signatures),
		len(signatures)+len(kps),
	)
	copy(extended, signatures)
	// Sign the hash
	for _, kp := range kps {
		sig, err := kp.SignDecorated(h[:])
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign transaction")
		}
		extended = append(extended, sig)
	}
	return extended, nil
}

func concatSignatureDecorated(e xdr.TransactionEnvelope, signatures []xdr.DecoratedSignature, newSignatures []xdr.DecoratedSignature) ([]xdr.DecoratedSignature, error) {
	extended := make([]xdr.DecoratedSignature, len(signatures)+len(newSignatures))
	copy(extended, signatures)
	copy(extended[len(signatures):], newSignatures)
	return extended, nil
}

func concatSignatureBase64(e xdr.TransactionEnvelope, signatures []xdr.DecoratedSignature, networkStr, publicKey, signature string) ([]xdr.DecoratedSignature, error) {
	if signature == "" {
		return nil, errors.New("signature not presented")
	}

	kp, err := keypair.ParseAddress(publicKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse the public key %s", publicKey)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to base64-decode the signature %s", signature)
	}

	h, err := network.HashTransactionInEnvelope(e, networkStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash transaction")
	}

	err = kp.Verify(h[:], sigBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify the signature")
	}

	extended := make([]xdr.DecoratedSignature, len(signatures), len(signatures)+1)
	copy(extended, signatures)
	extended = append(extended, xdr.DecoratedSignature{
		Hint:      xdr.SignatureHint(kp.Hint()),
		Signature: xdr.Signature(sigBytes),
	})

	return extended, nil
}

func stringsToKP(keys ...string) ([]*keypair.Full, error) {
	var signers []*keypair.Full
	for _, k := range keys {
		kp, err := keypair.Parse(k)
		if err != nil {
			return nil, errors.Wrapf(err, "provided string %s is not a valid Stellar key", k)
		}
		kpf, ok := kp.(*keypair.Full)
		if !ok {
			return nil, errors.New("provided string %s is not a valid Stellar secret key")
		}
		signers = append(signers, kpf)
	}

	return signers, nil
}

func concatHashX(signatures []xdr.DecoratedSignature, preimage []byte) ([]xdr.DecoratedSignature, error) {
	if maxSize := xdr.Signature(preimage).XDRMaxSize(); len(preimage) > maxSize {
		return nil, errors.Errorf(
			"preimage cannot be more than %d bytes", maxSize,
		)
	}
	extended := make(
		[]xdr.DecoratedSignature,
		len(signatures),
		len(signatures)+1,
	)
	copy(extended, signatures)

	preimageHash := sha256.Sum256(preimage)
	var hint [4]byte
	// copy the last 4-bytes of the signer public key to be used as hint
	copy(hint[:], preimageHash[28:])

	sig := xdr.DecoratedSignature{
		Hint:      xdr.SignatureHint(hint),
		Signature: xdr.Signature(preimage),
	}
	return append(extended, sig), nil
}

func marshallBinary(e xdr.TransactionEnvelope, signatures []xdr.DecoratedSignature) ([]byte, error) {
	switch e.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		e.V1.Signatures = signatures
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		e.V0.Signatures = signatures
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		e.FeeBump.Signatures = signatures
	default:
		panic("invalid transaction type: " + e.Type.String())
	}

	var txBytes bytes.Buffer
	_, err := xdr.Marshal(&txBytes, e)
	if err != nil {
		return nil, err
	}
	return txBytes.Bytes(), nil
}

func marshallBase64(e xdr.TransactionEnvelope, signatures []xdr.DecoratedSignature) (string, error) {
	binary, err := marshallBinary(e, signatures)
	if err != nil {
		return "", errors.Wrap(err, "failed to get XDR bytestring")
	}

	return base64.StdEncoding.EncodeToString(binary), nil
}

func marshallBase64Bytes(e xdr.TransactionEnvelope, signatures []xdr.DecoratedSignature) ([]byte, error) {
	binary, err := marshallBinary(e, signatures)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get XDR bytestring")
	}

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(binary)))
	base64.StdEncoding.Encode(encoded, binary)
	return encoded, nil
}

// Transaction represents a Stellar transaction. See
// https://developers.stellar.org/docs/glossary/transactions/
// A Transaction may be wrapped by a FeeBumpTransaction in which case
// the account authorizing the FeeBumpTransaction will pay for the transaction fees
// instead of the Transaction's source account.
type Transaction struct {
	envelope      xdr.TransactionEnvelope
	baseFee       int64
	maxFee        int64
	sourceAccount SimpleAccount
	operations    []Operation
	memo          Memo
	preconditions Preconditions
}

// BaseFee returns the per operation fee for this transaction.
func (t *Transaction) BaseFee() int64 {
	return t.baseFee
}

// MaxFee returns the total fees which can be spent to submit this transaction.
func (t *Transaction) MaxFee() int64 {
	return t.maxFee
}

// SourceAccount returns the account which is originating this account.
func (t *Transaction) SourceAccount() SimpleAccount {
	return t.sourceAccount
}

// SequenceNumber returns the sequence number of the transaction.
func (t *Transaction) SequenceNumber() int64 {
	return t.sourceAccount.Sequence
}

// Memo returns the memo configured for this transaction.
func (t *Transaction) Memo() Memo {
	return t.memo
}

// Timebounds returns the Timebounds configured for this transaction.
func (t *Transaction) Timebounds() TimeBounds {
	return t.preconditions.TimeBounds
}

// Operations returns the list of operations included in this transaction.
// The contents of the returned slice should not be modified.
func (t *Transaction) Operations() []Operation {
	return t.operations
}

// Signatures returns the list of signatures attached to this transaction.
// The contents of the returned slice should not be modified.
func (t *Transaction) Signatures() []xdr.DecoratedSignature {
	return t.envelope.Signatures()
}

// Hash returns the network specific hash of this transaction
// encoded as a byte array.
func (t *Transaction) Hash(networkStr string) ([32]byte, error) {
	return network.HashTransactionInEnvelope(t.envelope, networkStr)
}

// HashHex returns the network specific hash of this transaction
// encoded as a hexadecimal string.
func (t *Transaction) HashHex(network string) (string, error) {
	return hashHex(t.envelope, network)
}

func (t *Transaction) clone(signatures []xdr.DecoratedSignature) *Transaction {
	newTx := new(Transaction)
	*newTx = *t
	newTx.envelope = t.envelope

	switch newTx.envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		newTx.envelope.V1 = new(xdr.TransactionV1Envelope)
		*newTx.envelope.V1 = *t.envelope.V1
		newTx.envelope.V1.Signatures = signatures
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		newTx.envelope.V0 = new(xdr.TransactionV0Envelope)
		*newTx.envelope.V0 = *t.envelope.V0
		newTx.envelope.V0.Signatures = signatures
	default:
		panic("invalid transaction type: " + newTx.envelope.Type.String())
	}

	return newTx
}

// Sign returns a new Transaction instance which extends the current instance
// with additional signatures derived from the given list of keypair instances.
func (t *Transaction) Sign(network string, kps ...*keypair.Full) (*Transaction, error) {
	extendedSignatures, err := concatSignatures(t.envelope, network, t.Signatures(), kps...)
	if err != nil {
		return nil, err
	}

	return t.clone(extendedSignatures), nil
}

// SignWithKeyString returns a new Transaction instance which extends the current instance
// with additional signatures derived from the given list of private key strings.
func (t *Transaction) SignWithKeyString(network string, keys ...string) (*Transaction, error) {
	kps, err := stringsToKP(keys...)
	if err != nil {
		return nil, err
	}
	return t.Sign(network, kps...)
}

// SignHashX returns a new Transaction instance which extends the current instance
// with HashX signature type.
// See description here: https://developers.stellar.org/docs/glossary/multisig/#hashx
func (t *Transaction) SignHashX(preimage []byte) (*Transaction, error) {
	extendedSignatures, err := concatHashX(t.Signatures(), preimage)
	if err != nil {
		return nil, err
	}

	return t.clone(extendedSignatures), nil
}

// ClearSignatures returns a new Transaction instance which extends the current instance
// with signatures removed.
func (t *Transaction) ClearSignatures() (*Transaction, error) {
	return t.clone(nil), nil
}

// AddSignatureDecorated returns a new Transaction instance which extends the current instance
// with an additional decorated signature(s).
func (t *Transaction) AddSignatureDecorated(signature ...xdr.DecoratedSignature) (*Transaction, error) {
	extendedSignatures, err := concatSignatureDecorated(t.envelope, t.Signatures(), signature)
	if err != nil {
		return nil, err
	}

	return t.clone(extendedSignatures), nil
}

// AddSignatureBase64 returns a new Transaction instance which extends the current instance
// with an additional signature derived from the given base64-encoded signature.
func (t *Transaction) AddSignatureBase64(network, publicKey, signature string) (*Transaction, error) {
	extendedSignatures, err := concatSignatureBase64(t.envelope, t.Signatures(), network, publicKey, signature)
	if err != nil {
		return nil, err
	}

	return t.clone(extendedSignatures), nil
}

// ToXDR returns the a xdr.TransactionEnvelope which is equivalent to this transaction.
// The envelope should not be modified because any changes applied may
// affect the internals of the Transaction instance.
func (t *Transaction) ToXDR() xdr.TransactionEnvelope {
	return t.envelope
}

// MarshalBinary returns the binary XDR representation of the transaction envelope.
func (t *Transaction) MarshalBinary() ([]byte, error) {
	return marshallBinary(t.envelope, t.Signatures())
}

// MarshalText returns the base64 XDR representation of the transaction envelope.
func (t *Transaction) MarshalText() ([]byte, error) {
	return marshallBase64Bytes(t.envelope, t.Signatures())
}

// UnmarshalText consumes into the value the base64 XDR representation of the
// transaction envelope.
func (t *Transaction) UnmarshalText(b []byte) error {
	gtx, err := TransactionFromXDR(string(b))
	if err != nil {
		return err
	}
	tx, ok := gtx.Transaction()
	if !ok {
		return errors.New("transaction envelope unmarshaled into FeeBumpTransaction is not a fee bump transaction")
	}
	*t = *tx
	return nil
}

// Base64 returns the base 64 XDR representation of the transaction envelope.
func (t *Transaction) Base64() (string, error) {
	return marshallBase64(t.envelope, t.Signatures())
}

// ToGenericTransaction creates a GenericTransaction containing the Transaction.
func (t *Transaction) ToGenericTransaction() *GenericTransaction {
	return &GenericTransaction{simple: t}
}

// FeeBumpTransaction represents a CAP 15 fee bump transaction.
// Fee bump transactions allow an arbitrary account to pay the fee for a transaction.
type FeeBumpTransaction struct {
	envelope   xdr.TransactionEnvelope
	baseFee    int64
	maxFee     int64
	feeAccount string
	inner      *Transaction
}

// BaseFee returns the per operation fee for this transaction.
func (t *FeeBumpTransaction) BaseFee() int64 {
	return t.baseFee
}

// MaxFee returns the total fees which can be spent to submit this transaction.
func (t *FeeBumpTransaction) MaxFee() int64 {
	return t.maxFee
}

// FeeAccount returns the address of the account which will be paying for the inner transaction.
func (t *FeeBumpTransaction) FeeAccount() string {
	return t.feeAccount
}

// Signatures returns the list of signatures attached to this transaction.
// The contents of the returned slice should not be modified.
func (t *FeeBumpTransaction) Signatures() []xdr.DecoratedSignature {
	return t.envelope.FeeBumpSignatures()
}

// Hash returns the network specific hash of this transaction
// encoded as a byte array.
func (t *FeeBumpTransaction) Hash(networkStr string) ([32]byte, error) {
	return network.HashTransactionInEnvelope(t.envelope, networkStr)
}

// HashHex returns the network specific hash of this transaction
// encoded as a hexadecimal string.
func (t *FeeBumpTransaction) HashHex(network string) (string, error) {
	return hashHex(t.envelope, network)
}

func (t *FeeBumpTransaction) clone(signatures []xdr.DecoratedSignature) *FeeBumpTransaction {
	newTx := new(FeeBumpTransaction)
	*newTx = *t
	newTx.envelope.FeeBump = new(xdr.FeeBumpTransactionEnvelope)
	*newTx.envelope.FeeBump = *t.envelope.FeeBump
	newTx.envelope.FeeBump.Signatures = signatures
	return newTx
}

// Sign returns a new FeeBumpTransaction instance which extends the current instance
// with additional signatures derived from the given list of keypair instances.
func (t *FeeBumpTransaction) Sign(network string, kps ...*keypair.Full) (*FeeBumpTransaction, error) {
	extendedSignatures, err := concatSignatures(t.envelope, network, t.Signatures(), kps...)
	if err != nil {
		return nil, err
	}

	return t.clone(extendedSignatures), nil
}

// SignWithKeyString returns a new FeeBumpTransaction instance which extends the current instance
// with additional signatures derived from the given list of private key strings.
func (t *FeeBumpTransaction) SignWithKeyString(network string, keys ...string) (*FeeBumpTransaction, error) {
	kps, err := stringsToKP(keys...)
	if err != nil {
		return nil, err
	}
	return t.Sign(network, kps...)
}

// SignHashX returns a new FeeBumpTransaction instance which extends the current instance
// with HashX signature type.
// See description here: https://developers.stellar.org/docs/glossary/multisig/#hashx
func (t *FeeBumpTransaction) SignHashX(preimage []byte) (*FeeBumpTransaction, error) {
	extendedSignatures, err := concatHashX(t.Signatures(), preimage)
	if err != nil {
		return nil, err
	}

	return t.clone(extendedSignatures), nil
}

// ClearSignatures returns a new Transaction instance which extends the current instance
// with signatures removed.
func (t *FeeBumpTransaction) ClearSignatures() (*FeeBumpTransaction, error) {
	return t.clone(nil), nil
}

// AddSignatureDecorated returns a new FeeBumpTransaction instance which extends the current instance
// with an additional decorated signature(s).
func (t *FeeBumpTransaction) AddSignatureDecorated(signature ...xdr.DecoratedSignature) (*FeeBumpTransaction, error) {
	extendedSignatures, err := concatSignatureDecorated(t.envelope, t.Signatures(), signature)
	if err != nil {
		return nil, err
	}

	return t.clone(extendedSignatures), nil
}

// AddSignatureBase64 returns a new FeeBumpTransaction instance which extends the current instance
// with an additional signature derived from the given base64-encoded signature.
func (t *FeeBumpTransaction) AddSignatureBase64(network, publicKey, signature string) (*FeeBumpTransaction, error) {
	extendedSignatures, err := concatSignatureBase64(t.envelope, t.Signatures(), network, publicKey, signature)
	if err != nil {
		return nil, err
	}

	return t.clone(extendedSignatures), nil
}

// ToXDR returns the a xdr.TransactionEnvelope which is equivalent to this transaction.
// The envelope should not be modified because any changes applied may
// affect the internals of the FeeBumpTransaction instance.
func (t *FeeBumpTransaction) ToXDR() xdr.TransactionEnvelope {
	return t.envelope
}

// MarshalBinary returns the binary XDR representation of the transaction envelope.
func (t *FeeBumpTransaction) MarshalBinary() ([]byte, error) {
	return marshallBinary(t.envelope, t.Signatures())
}

// MarshalText returns the base64 XDR representation of the transaction
// envelope.
func (t *FeeBumpTransaction) MarshalText() ([]byte, error) {
	return marshallBase64Bytes(t.envelope, t.Signatures())
}

// UnmarshalText consumes into the value the base64 XDR representation of the
// transaction envelope.
func (t *FeeBumpTransaction) UnmarshalText(b []byte) error {
	gtx, err := TransactionFromXDR(string(b))
	if err != nil {
		return err
	}
	fbtx, ok := gtx.FeeBump()
	if !ok {
		return errors.New("transaction envelope unmarshaled into Transaction is not a transaction")
	}
	*t = *fbtx
	return nil
}

// Base64 returns the base 64 XDR representation of the transaction envelope.
func (t *FeeBumpTransaction) Base64() (string, error) {
	return marshallBase64(t.envelope, t.Signatures())
}

// ToGenericTransaction creates a GenericTransaction containing the
// FeeBumpTransaction.
func (t *FeeBumpTransaction) ToGenericTransaction() *GenericTransaction {
	return &GenericTransaction{feeBump: t}
}

// InnerTransaction returns the Transaction which is wrapped by
// this FeeBumpTransaction instance.
func (t *FeeBumpTransaction) InnerTransaction() *Transaction {
	innerCopy := new(Transaction)
	*innerCopy = *t.inner
	return innerCopy
}

// GenericTransaction represents a parsed transaction envelope returned by TransactionFromXDR.
// A GenericTransaction can be either a Transaction or a FeeBumpTransaction.
type GenericTransaction struct {
	simple  *Transaction
	feeBump *FeeBumpTransaction
}

// NewGenericTransactionWithTransaction creates a GenericTransaction containing
// a Transaction.
func NewGenericTransactionWithTransaction(tx *Transaction) *GenericTransaction {
	return &GenericTransaction{simple: tx}
}

// NewGenericTransactionWithFeeBumpTransaction creates a GenericTransaction
// containing a FeeBumpTransaction.
func NewGenericTransactionWithFeeBumpTransaction(feeBumpTx *FeeBumpTransaction) *GenericTransaction {
	return &GenericTransaction{feeBump: feeBumpTx}
}

// Transaction unpacks the GenericTransaction instance into a Transaction.
// The function also returns a boolean which is true if the GenericTransaction can be
// unpacked into a Transaction.
func (t GenericTransaction) Transaction() (*Transaction, bool) {
	return t.simple, t.simple != nil
}

// FeeBump unpacks the GenericTransaction instance into a FeeBumpTransaction.
// The function also returns a boolean which is true if the GenericTransaction
// can be unpacked into a FeeBumpTransaction.
func (t GenericTransaction) FeeBump() (*FeeBumpTransaction, bool) {
	return t.feeBump, t.feeBump != nil
}

// ToXDR returns the a xdr.TransactionEnvelope which is equivalent to this
// transaction. The envelope should not be modified because any changes applied
// may affect the internals of the GenericTransaction.
func (t *GenericTransaction) ToXDR() (xdr.TransactionEnvelope, error) {
	if tx, ok := t.Transaction(); ok {
		return tx.envelope, nil
	}
	if fbtx, ok := t.FeeBump(); ok {
		return fbtx.envelope, nil
	}
	return xdr.TransactionEnvelope{}, fmt.Errorf("unable to get xdr of empty GenericTransaction")
}

// Hash returns the network specific hash of this transaction
// encoded as a byte array.
func (t GenericTransaction) Hash(networkStr string) ([32]byte, error) {
	if tx, ok := t.Transaction(); ok {
		return tx.Hash(networkStr)
	}
	if fbtx, ok := t.FeeBump(); ok {
		return fbtx.Hash(networkStr)
	}
	return [32]byte{}, fmt.Errorf("unable to get hash of empty GenericTransaction")
}

// HashHex returns the network specific hash of this transaction
// encoded as a hexadecimal string.
func (t GenericTransaction) HashHex(network string) (string, error) {
	if tx, ok := t.Transaction(); ok {
		return tx.HashHex(network)
	}
	if fbtx, ok := t.FeeBump(); ok {
		return fbtx.HashHex(network)
	}
	return "", fmt.Errorf("unable to get hash of empty GenericTransaction")
}

// MarshalBinary returns the binary XDR representation of the transaction
// envelope.
func (t *GenericTransaction) MarshalBinary() ([]byte, error) {
	if tx, ok := t.Transaction(); ok {
		return tx.MarshalBinary()
	}
	if fbtx, ok := t.FeeBump(); ok {
		return fbtx.MarshalBinary()
	}
	return nil, errors.New("unable to marshal empty GenericTransaction")
}

// MarshalText returns the base64 XDR representation of the transaction
// envelope.
func (t *GenericTransaction) MarshalText() ([]byte, error) {
	if tx, ok := t.Transaction(); ok {
		return tx.MarshalText()
	}
	if fbtx, ok := t.FeeBump(); ok {
		return fbtx.MarshalText()
	}
	return nil, errors.New("unable to marshal empty GenericTransaction")
}

// UnmarshalText consumes into the value the base64 XDR representation of the
// transaction envelope.
func (t *GenericTransaction) UnmarshalText(b []byte) error {
	gtx, err := TransactionFromXDR(string(b))
	if err != nil {
		return err
	}
	*t = *gtx
	return nil
}

// TransactionFromXDR parses the supplied transaction envelope in base64 XDR
// and returns a GenericTransaction instance.
func TransactionFromXDR(txeB64 string) (*GenericTransaction, error) {
	var xdrEnv xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(txeB64, &xdrEnv)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal transaction envelope")
	}
	return transactionFromParsedXDR(xdrEnv)
}

func transactionFromParsedXDR(xdrEnv xdr.TransactionEnvelope) (*GenericTransaction, error) {
	var err error
	newTx := &GenericTransaction{}

	if xdrEnv.IsFeeBump() {
		var innerTx *GenericTransaction
		innerTx, err = transactionFromParsedXDR(xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1:   xdrEnv.FeeBump.Tx.InnerTx.V1,
		})
		if err != nil {
			return newTx, errors.New("could not parse inner transaction")
		}
		feeBumpAccount := xdrEnv.FeeBumpAccount()
		feeAccount := feeBumpAccount.Address()

		newTx.feeBump = &FeeBumpTransaction{
			envelope: xdrEnv,
			// A fee-bump transaction has an effective number of operations equal to one plus the
			// number of operations in the inner transaction. Correspondingly, the minimum fee for
			// the fee-bump transaction is one base fee more than the minimum fee for the inner
			// transaction.
			baseFee:    xdrEnv.FeeBumpFee() / int64(len(innerTx.simple.operations)+1),
			maxFee:     xdrEnv.FeeBumpFee(),
			inner:      innerTx.simple,
			feeAccount: feeAccount,
		}

		return newTx, nil
	}
	sourceAccount := xdrEnv.SourceAccount()
	accountID := sourceAccount.Address()

	totalFee := int64(xdrEnv.Fee())
	baseFee := totalFee
	if count := int64(len(xdrEnv.Operations())); count > 0 {
		baseFee = baseFee / count
	}

	newTx.simple = &Transaction{
		envelope: xdrEnv,
		baseFee:  baseFee,
		maxFee:   totalFee,
		sourceAccount: SimpleAccount{
			AccountID: accountID,
			Sequence:  xdrEnv.SeqNum(),
		},
		operations: nil,
		memo:       nil,
	}

	newTx.simple.preconditions.FromXDR(xdrEnv.Preconditions())
	newTx.simple.memo, err = memoFromXDR(xdrEnv.Memo())
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse memo")
	}

	operations := xdrEnv.Operations()
	for _, op := range operations {
		newOp, err := operationFromXDR(op)
		if err != nil {
			return nil, err
		}
		// if it's a soroban transaction, and we found a InvokeHostFunction operation.
		if xdrEnv.V1 != nil && xdrEnv.V1.Tx.Ext.V != 0 {
			//if invoke, ok := newOp.(*InvokeHostFunction); ok {
			//	invoke.Ext = xdrEnv.V1.Tx.Ext
			//}
		}
		newTx.simple.operations = append(newTx.simple.operations, newOp)
	}

	return newTx, nil
}

// TransactionParams is a container for parameters which are used to construct
// new Transaction instances
type TransactionParams struct {
	SourceAccount        Account
	IncrementSequenceNum bool
	Operations           []Operation
	BaseFee              int64
	Memo                 Memo
	Preconditions        Preconditions
}

// NewTransaction returns a new Transaction instance
func NewTransaction(params TransactionParams) (*Transaction, error) {
	if params.SourceAccount == nil {
		return nil, errors.New("transaction has no source account")
	}

	var sequence int64
	var err error
	if params.IncrementSequenceNum {
		sequence, err = params.SourceAccount.IncrementSequenceNumber()
	} else {
		sequence, err = params.SourceAccount.GetSequenceNumber()
	}
	if err != nil {
		return nil, errors.Wrap(err, "could not obtain account sequence")
	}

	tx := &Transaction{
		baseFee: params.BaseFee,
		sourceAccount: SimpleAccount{
			AccountID: params.SourceAccount.GetAccountID(),
			Sequence:  sequence,
		},
		operations:    params.Operations,
		memo:          params.Memo,
		preconditions: params.Preconditions,
	}
	var sourceAccount xdr.MuxedAccount
	if err = sourceAccount.SetAddress(tx.sourceAccount.AccountID); err != nil {
		return nil, errors.Wrap(err, "account id is not valid")
	}
	if tx.baseFee < 0 {
		return nil, errors.Errorf("base fee cannot be negative")
	}

	if len(tx.operations) == 0 {
		return nil, errors.New("transaction has no operations")
	}

	// check if maxFee fits in a uint32
	// 64 bit fees are only available in fee bump transactions
	// if maxFee is negative then there must have been an int overflow
	hi, lo := bits.Mul64(uint64(params.BaseFee), uint64(len(params.Operations)))
	if hi > 0 || lo > math.MaxUint32 {
		return nil, errors.Errorf(
			"base fee %d results in an overflow of max fee", params.BaseFee)
	}
	tx.maxFee = int64(lo)

	// Check that all preconditions are valid
	if err = tx.preconditions.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid preconditions")
	}

	precondXdr, err := tx.preconditions.BuildXDR()
	if err != nil {
		return nil, errors.Wrap(err, "invalid preconditions")
	}

	envelope := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: xdr.Transaction{
				SourceAccount: sourceAccount,
				Fee:           xdr.Uint32(tx.maxFee),
				SeqNum:        xdr.SequenceNumber(sequence),
				Cond:          precondXdr,
			},
			Signatures: nil,
		},
	}

	// Handle the memo, if one is present
	if tx.memo != nil {
		var xdrMemo xdr.Memo
		xdrMemo, err = tx.memo.ToXDR()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't build memo XDR")
		}
		envelope.V1.Tx.Memo = xdrMemo
	}

	var sorobanOp SorobanOperation

	for _, op := range tx.operations {
		if verr := op.Validate(); verr != nil {
			return nil, errors.Wrap(verr, fmt.Sprintf("validation failed for %T operation", op))
		}
		xdrOperation, err2 := op.BuildXDR()
		if err2 != nil {
			return nil, errors.Wrap(err2, fmt.Sprintf("failed to build operation %T", op))
		}
		envelope.V1.Tx.Operations = append(envelope.V1.Tx.Operations, xdrOperation)

		if scOp, ok := op.(SorobanOperation); ok {
			// this is a smart contract operation.
			// smart contract operations are limited to 1 operation / transaction.
			sorobanOp = scOp
		}
	}

	// In case it's a smart contract transaction, we need to include the Ext field within the envelope.
	if sorobanOp != nil {
		envelope.V1.Tx.Ext, err = sorobanOp.BuildTransactionExt()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to build operation %T", sorobanOp))
		}
	}

	tx.envelope = envelope
	return tx, nil
}

// FeeBumpTransactionParams is a container for parameters
// which are used to construct new FeeBumpTransaction instances
type FeeBumpTransactionParams struct {
	Inner      *Transaction
	FeeAccount string
	BaseFee    int64
}

func convertToV1(tx *Transaction) (*Transaction, error) {
	sourceAccount := tx.SourceAccount()
	signatures := tx.Signatures()
	tx, err := NewTransaction(TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: false,
		Operations:           tx.Operations(),
		BaseFee:              tx.BaseFee(),
		Memo:                 tx.Memo(),
		Preconditions:        Preconditions{TimeBounds: tx.Timebounds()},
	})
	if err != nil {
		return tx, err
	}
	tx.envelope.V1.Signatures = signatures
	return tx, nil
}

// NewFeeBumpTransaction returns a new FeeBumpTransaction instance
func NewFeeBumpTransaction(params FeeBumpTransactionParams) (*FeeBumpTransaction, error) {
	inner := params.Inner
	if inner == nil {
		return nil, errors.New("inner transaction is missing")
	}
	switch inner.envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx, xdr.EnvelopeTypeEnvelopeTypeTxV0:
	default:
		return nil, errors.Errorf("%s transactions cannot be fee bumped", inner.envelope.Type)
	}

	innerEnv := inner.ToXDR()
	if innerEnv.Type == xdr.EnvelopeTypeEnvelopeTypeTxV0 {
		var err error
		inner, err = convertToV1(inner)
		if err != nil {
			return nil, errors.Wrap(err, "could not upgrade transaction from v0 to v1")
		}
	} else if innerEnv.Type != xdr.EnvelopeTypeEnvelopeTypeTx {
		return nil, errors.Errorf("%v transactions cannot be fee bumped", innerEnv.Type.String())
	}

	tx := &FeeBumpTransaction{
		baseFee: params.BaseFee,
		// A fee-bump transaction has an effective number of operations equal to one plus the
		// number of operations in the inner transaction. Correspondingly, the minimum fee for
		// the fee-bump transaction is one base fee more than the minimum fee for the inner
		// transaction.
		maxFee:     params.BaseFee * int64(len(inner.operations)+1),
		feeAccount: params.FeeAccount,
		inner:      new(Transaction),
	}
	*tx.inner = *inner

	hi, lo := bits.Mul64(uint64(params.BaseFee), uint64(len(inner.operations)+1))
	if hi > 0 || lo > math.MaxInt64 {
		return nil, errors.Errorf("base fee %d results in an overflow of max fee", params.BaseFee)
	}
	tx.maxFee = int64(lo)

	if tx.baseFee < tx.inner.baseFee {
		return tx, errors.New("base fee cannot be lower than provided inner transaction fee")
	}
	if tx.baseFee < MinBaseFee {
		return tx, errors.Errorf(
			"base fee cannot be lower than network minimum of %d", MinBaseFee,
		)
	}

	var feeSource xdr.MuxedAccount
	if err := feeSource.SetAddress(tx.feeAccount); err != nil {
		return tx, errors.Wrap(err, "fee account is not a valid address")
	}

	tx.envelope = xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTxFeeBump,
		FeeBump: &xdr.FeeBumpTransactionEnvelope{
			Tx: xdr.FeeBumpTransaction{
				FeeSource: feeSource,
				Fee:       xdr.Int64(tx.maxFee),
				InnerTx: xdr.FeeBumpTransactionInnerTx{
					Type: xdr.EnvelopeTypeEnvelopeTypeTx,
					V1:   innerEnv.V1,
				},
			},
		},
	}

	return tx, nil
}

// generateRandomNonce creates a cryptographically secure random slice of `n` bytes.
func generateRandomNonce(n int) ([]byte, error) {
	binary := make([]byte, n)
	_, err := rand.Read(binary)

	if err != nil {
		return []byte{}, err
	}

	return binary, err
}

// verifyTxSignature checks if a transaction has been signed by the provided Stellar account.
func verifyTxSignature(tx *Transaction, network string, signer string) error {
	signersFound, err := verifyTxSignatures(tx, network, signer)
	if len(signersFound) == 0 {
		return errors.Errorf("transaction not signed by %s", signer)
	}
	return err
}

// verifyTxSignature checks if a transaction has been signed by one or more of
// the signers, returning a list of signers that were found to have signed the
// transaction.
func verifyTxSignatures(tx *Transaction, network string, signers ...string) ([]string, error) {
	txHash, err := tx.Hash(network)
	if err != nil {
		return nil, err
	}

	// find and verify signatures
	signatureUsed := map[int]bool{}
	signersFound := set.Set[string]{}
	for _, signer := range signers {
		kp, err := keypair.ParseAddress(signer)
		if err != nil {
			return nil, errors.Wrap(err, "signer not address")
		}

		for i, decSig := range tx.Signatures() {
			if signatureUsed[i] {
				continue
			}
			if decSig.Hint != kp.Hint() {
				continue
			}
			err := kp.Verify(txHash[:], decSig.Signature)
			if err == nil {
				signatureUsed[i] = true
				signersFound.Add(signer)
				break
			}
		}
	}

	signersFoundList := make([]string, 0, len(signersFound))
	for _, signer := range signers {
		if signersFound.Contains(signer) {
			signersFoundList = append(signersFoundList, signer)
			delete(signersFound, signer)
		}
	}
	return signersFoundList, nil
}
