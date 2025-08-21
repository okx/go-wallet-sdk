package v2

import (
	"encoding/json"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
	"golang.org/x/crypto/sha3"
)

//region RawTransaction

var rawTransactionPrehash []byte
var rawTransactionPrehashEndless []byte

const rawTransactionPrehashStr = "APTOS::RawTransaction"
const rawTransactionPrehashStrEndless = "ENDLESS::RawTransaction"

// RawTransactionPrehash Return the sha3-256 prehash for RawTransaction
// Do not write to the []byte returned
func RawTransactionPrehash() []byte {
	if rawTransactionPrehash == nil {
		rawTransactionPrehash = rawTransactionPrehashWithPrefix(rawTransactionPrehashStr)
	}
	return rawTransactionPrehash
}

func RawTransactionPrehashEndless() []byte {
	if rawTransactionPrehashEndless == nil {
		rawTransactionPrehashEndless = rawTransactionPrehashWithPrefix(rawTransactionPrehashStrEndless)
	}
	return rawTransactionPrehashEndless
}

func rawTransactionPrehashWithPrefix(prefix string) []byte {
	b32 := sha3.Sum256([]byte(prefix))
	out := make([]byte, len(b32))
	copy(out, b32[:])
	return out
}

type RawTransactionImpl interface {
	bcs.Struct

	// SigningMessage creates a raw signing message for the transaction
	// Note that this should only be used externally if signing transactions outside the SDK.  Otherwise, use Sign.
	SigningMessage() (message []byte, err error)

	// Sign signs a transaction and returns the associated AccountAuthenticator, it will underneath sign the SigningMessage
	Sign(signer crypto.Signer) (*crypto.AccountAuthenticator, error)
}

// RawTransaction representation of a transaction's parts prior to signing
// Implements crypto.MessageSigner, crypto.Signer, bcs.Struct
type RawTransaction struct {
	Sender         AccountAddress
	SequenceNumber uint64
	Payload        TransactionPayload
	MaxGasAmount   uint64
	GasUnitPrice   uint64

	// ExpirationTimestampSeconds is seconds since Unix epoch
	ExpirationTimestampSeconds uint64

	ChainId uint8
}

func (txn *RawTransaction) SignedTransaction(sender crypto.Signer) (*SignedTransaction, error) {
	auth, err := txn.Sign(sender)
	if err != nil {
		return nil, err
	}
	return txn.SignedTransactionWithAuthenticator(auth)
}

// SignedTransactionWithAuthenticator signs the sender only signed transaction
func (txn *RawTransaction) SignedTransactionWithAuthenticator(auth *crypto.AccountAuthenticator) (*SignedTransaction, error) {
	txnAuth, err := NewTransactionAuthenticator(auth)
	if err != nil {
		return nil, err
	}
	return &SignedTransaction{
		Transaction:   txn,
		Authenticator: txnAuth,
	}, nil
}

//region RawTransaction bcs.Struct

func (txn *RawTransaction) MarshalBCS(ser *bcs.Serializer) {
	txn.Sender.MarshalBCS(ser)
	ser.U64(txn.SequenceNumber)
	txn.Payload.MarshalBCS(ser)
	ser.U64(txn.MaxGasAmount)
	ser.U64(txn.GasUnitPrice)
	ser.U64(txn.ExpirationTimestampSeconds)
	ser.U8(txn.ChainId)
}

func (txn *RawTransaction) UnmarshalBCS(des *bcs.Deserializer) {
	txn.Sender.UnmarshalBCS(des)
	txn.SequenceNumber = des.U64()
	txn.Payload.UnmarshalBCS(des)
	txn.MaxGasAmount = des.U64()
	txn.GasUnitPrice = des.U64()
	txn.ExpirationTimestampSeconds = des.U64()
	txn.ChainId = des.U8()
}

//endregion

//region RawTransaction MessageSigner

// SigningMessage generates the bytes needed to be signed by a signer
func (txn *RawTransaction) SigningMessage() (message []byte, err error) {
	txnBytes, err := bcs.Serialize(txn)
	if err != nil {
		return
	}
	var prehash []byte
	prehash = RawTransactionPrehash()
	//if common.IsEndless(txn.ChainId) {
	//	prehash = RawTransactionPrehashEndless()
	//} else {
	//	prehash = RawTransactionPrehash()
	//}
	message = make([]byte, len(prehash)+len(txnBytes))
	copy(message, prehash)
	copy(message[len(prehash):], txnBytes)
	return message, nil
}

//endregion

// String returns a JSON formatted string representation of the RawTransaction
func (txn *RawTransaction) String() string {
	jsonBytes, err := json.MarshalIndent(txn, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling RawTransaction: %v", err)
	}
	return string(jsonBytes)
}

//region RawTransaction Signer

func (txn *RawTransaction) Sign(signer crypto.Signer) (authenticator *crypto.AccountAuthenticator, err error) {
	message, err := txn.SigningMessage()
	if err != nil {
		return
	}
	return signer.Sign(message)
}

//endregion
//endregion

//region RawTransactionWithData

var rawTransactionWithDataPrehash []byte
var rawTransactionWithDataPrehashEndless []byte

const rawTransactionWithDataPrehashStr = "APTOS::RawTransactionWithData"
const rawTransactionWithDataPrehashStrEndless = "ENDLESS::RawTransactionWithData"

// RawTransactionWithDataPrehash Return the sha3-256 prehash for RawTransactionWithData
// Do not write to the []byte returned
func RawTransactionWithDataPrehash() []byte {
	// Cache the prehash
	if rawTransactionWithDataPrehash == nil {
		rawTransactionWithDataPrehash = rawTransactionPrehashWithPrefix(rawTransactionWithDataPrehashStr)
	}
	return rawTransactionWithDataPrehash
}

func RawTransactionWithDataPrehashEndless() []byte {
	// Cache the prehash
	if rawTransactionWithDataPrehashEndless == nil {
		rawTransactionWithDataPrehashEndless = rawTransactionPrehashWithPrefix(rawTransactionWithDataPrehashStrEndless)
	}
	return rawTransactionWithDataPrehashEndless
}

type RawTransactionWithDataVariant uint32

const (
	MultiAgentRawTransactionWithDataVariant             RawTransactionWithDataVariant = 0
	MultiAgentWithFeePayerRawTransactionWithDataVariant RawTransactionWithDataVariant = 1
)

type RawTransactionWithDataImpl interface {
	bcs.Struct
}

// TODO: make a function to make creating this easier

type RawTransactionWithData struct {
	Variant RawTransactionWithDataVariant
	Inner   RawTransactionWithDataImpl
}

func (txn *RawTransactionWithData) SetFeePayer(
	feePayer AccountAddress,
) bool {
	if inner, ok := txn.Inner.(*MultiAgentWithFeePayerRawTransactionWithData); ok {
		inner.FeePayer = &feePayer
		return true
	}
	return false
}

func (txn *RawTransactionWithData) ToMultiAgentSignedTransaction(
	sender *crypto.AccountAuthenticator,
	additionalSigners []crypto.AccountAuthenticator,
) (*SignedTransaction, bool) {
	switch multiAgent := txn.Inner.(type) {
	case *MultiAgentRawTransactionWithData:
		return &SignedTransaction{
			Transaction: multiAgent.RawTxn,
			Authenticator: &TransactionAuthenticator{
				Variant: TransactionAuthenticatorMultiAgent,
				Auth: &MultiAgentTransactionAuthenticator{
					Sender:                   sender,
					SecondarySignerAddresses: multiAgent.SecondarySigners,
					SecondarySigners:         additionalSigners,
				},
			},
		}, true
	default:
		return nil, false
	}
}

func (txn *RawTransactionWithData) ToFeePayerSignedTransaction(
	sender *crypto.AccountAuthenticator,
	feePayerAuthenticator *crypto.AccountAuthenticator,
	additionalSigners []crypto.AccountAuthenticator,
) (*SignedTransaction, bool) {
	switch feePayerTxn := txn.Inner.(type) {
	case *MultiAgentWithFeePayerRawTransactionWithData:
		return &SignedTransaction{
			Transaction: feePayerTxn.RawTxn,
			Authenticator: &TransactionAuthenticator{
				Variant: TransactionAuthenticatorFeePayer,
				Auth: &FeePayerTransactionAuthenticator{
					Sender:                   sender,
					SecondarySignerAddresses: feePayerTxn.SecondarySigners,
					SecondarySigners:         additionalSigners,
					FeePayer:                 feePayerTxn.FeePayer,
					FeePayerAuthenticator:    feePayerAuthenticator,
				},
			},
		}, true
	default:
		return nil, false
	}
}

// MarshalTypeScriptBCS converts to RawTransactionWithData to the TypeScript type MultiAgentTransaction
func (txn *RawTransactionWithData) MarshalTypeScriptBCS(ser *bcs.Serializer) {
	switch inner := (txn.Inner).(type) {
	case *MultiAgentRawTransactionWithData:
		ser.Struct(inner.RawTxn)
		bcs.SerializeSequence(inner.SecondarySigners, ser)
		ser.Bool(false)
	case *MultiAgentWithFeePayerRawTransactionWithData:
		ser.Struct(inner.RawTxn)
		bcs.SerializeSequence(inner.SecondarySigners, ser)
		ser.Bool(true)
		ser.Struct(inner.FeePayer)
	}
}

// UnmarshalTypeScriptBCS converts to RawTransactionWithData from the TypeScript type MultiAgentTransaction
func (txn *RawTransactionWithData) UnmarshalTypeScriptBCS(des *bcs.Deserializer) {
	rawTxn := &RawTransaction{}
	des.Struct(rawTxn)
	secondarySigners := bcs.DeserializeSequence[AccountAddress](des)
	feePayer := bcs.DeserializeOption(des, func(des *bcs.Deserializer, out *AccountAddress) {
		des.Struct(out)
	})
	if des.Error() != nil {
		return
	}

	if feePayer == nil {
		txn.Variant = MultiAgentRawTransactionWithDataVariant
		txn.Inner = &MultiAgentRawTransactionWithData{
			RawTxn:           rawTxn,
			SecondarySigners: secondarySigners,
		}
	} else {
		txn.Variant = MultiAgentWithFeePayerRawTransactionWithDataVariant
		txn.Inner = &MultiAgentWithFeePayerRawTransactionWithData{
			RawTxn:           rawTxn,
			SecondarySigners: secondarySigners,
			FeePayer:         feePayer,
		}
	}
}

//region RawTransactionWithData Signer

func (txn *RawTransactionWithData) Sign(signer crypto.Signer) (authenticator *crypto.AccountAuthenticator, err error) {
	message, err := txn.SigningMessage()
	if err != nil {
		return
	}
	return signer.Sign(message)
}

//endregion

//region RawTransactionWithData MessageSigner

func (txn *RawTransactionWithData) SigningMessage() (message []byte, err error) {
	txnBytes, err := bcs.Serialize(txn)
	if err != nil {
		return
	}
	//var chainId uint8
	//switch txn.Variant {
	//case MultiAgentRawTransactionWithDataVariant:
	//	rawTx := txn.Inner.(*MultiAgentRawTransactionWithData)
	//	//chainId = rawTx.RawTxn.ChainId
	//case MultiAgentWithFeePayerRawTransactionWithDataVariant:
	//	rawTx := txn.Inner.(*MultiAgentWithFeePayerRawTransactionWithData)
	//	//chainId = rawTx.RawTxn.ChainId
	//}
	var prehash []byte
	prehash = RawTransactionWithDataPrehash()
	//if common.IsEndless(chainId) {
	//	prehash = RawTransactionWithDataPrehashEndless()
	//} else {
	//	prehash = RawTransactionWithDataPrehash()
	//}
	message = make([]byte, len(prehash)+len(txnBytes))
	copy(message, prehash)
	copy(message[len(prehash):], txnBytes)
	return message, nil
}

//endregion

//region RawTransactionWithData bcs.Struct

// String returns a JSON formatted string representation of the RawTransactionWithData
func (txn *RawTransactionWithData) String() string {
	jsonBytes, err := json.MarshalIndent(txn, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling RawTransactionWithData: %v", err)
	}
	return string(jsonBytes)
}

func (txn *RawTransactionWithData) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(txn.Variant))
	ser.Struct(txn.Inner)
}

func (txn *RawTransactionWithData) UnmarshalBCS(des *bcs.Deserializer) {
	txn.Variant = RawTransactionWithDataVariant(des.Uleb128())
	switch txn.Variant {
	case MultiAgentRawTransactionWithDataVariant:
		txn.Inner = &MultiAgentRawTransactionWithData{}
	case MultiAgentWithFeePayerRawTransactionWithDataVariant:
		txn.Inner = &MultiAgentWithFeePayerRawTransactionWithData{}
	default:
		des.SetError(fmt.Errorf("unknown RawTransactionWithData variant %d", txn.Variant))
		return
	}
	des.Struct(txn.Inner)
}

//endregion
//endregion

//region MultiAgentRawTransactionWithData

type MultiAgentRawTransactionWithData struct {
	RawTxn           *RawTransaction
	SecondarySigners []AccountAddress
}

//region MultiAgentRawTransactionWithData bcs.Struct

func (txn *MultiAgentRawTransactionWithData) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(txn.RawTxn)
	bcs.SerializeSequence(txn.SecondarySigners, ser)
}

func (txn *MultiAgentRawTransactionWithData) UnmarshalBCS(des *bcs.Deserializer) {
	txn.RawTxn = &RawTransaction{}
	des.Struct(txn.RawTxn)
	txn.SecondarySigners = bcs.DeserializeSequence[AccountAddress](des)
}

//endregion
//endregion

//region MultiAgentWithFeePayerRawTransactionWithData

type MultiAgentWithFeePayerRawTransactionWithData struct {
	RawTxn           *RawTransaction
	SecondarySigners []AccountAddress
	FeePayer         *AccountAddress
}

//region MultiAgentWithFeePayerRawTransactionWithData bcs.Struct

func (txn *MultiAgentWithFeePayerRawTransactionWithData) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(txn.RawTxn)
	bcs.SerializeSequence(txn.SecondarySigners, ser)
	ser.Struct(txn.FeePayer)
}

func (txn *MultiAgentWithFeePayerRawTransactionWithData) UnmarshalBCS(des *bcs.Deserializer) {
	txn.RawTxn = &RawTransaction{}
	des.Struct(txn.RawTxn)
	txn.SecondarySigners = bcs.DeserializeSequence[AccountAddress](des)
	txn.FeePayer = &AccountAddress{}
	des.Struct(txn.FeePayer)
}

//endregion
//endregion
