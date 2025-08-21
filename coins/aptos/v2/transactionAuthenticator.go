package v2

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
)

//region TransactionAuthenticator

type TransactionAuthenticatorVariant uint8

const (
	TransactionAuthenticatorEd25519      TransactionAuthenticatorVariant = 0
	TransactionAuthenticatorMultiEd25519 TransactionAuthenticatorVariant = 1
	TransactionAuthenticatorMultiAgent   TransactionAuthenticatorVariant = 2
	TransactionAuthenticatorFeePayer     TransactionAuthenticatorVariant = 3
	TransactionAuthenticatorSingleSender TransactionAuthenticatorVariant = 4
)

type TransactionAuthenticatorImpl interface {
	bcs.Struct
	// Verify Return true if this AccountAuthenticator approves
	Verify(data []byte) bool
}

// TransactionAuthenticator is used for authorizing a transaction.  This differs from crypto.AccountAuthenticator because it handles
// constructs like FeePayer and MultiAgent.  Some keys can't stand on their own as TransactionAuthenticators.
// Implements TransactionAuthenticatorImpl, bcs.Struct
type TransactionAuthenticator struct {
	Variant TransactionAuthenticatorVariant
	Auth    TransactionAuthenticatorImpl
}

func NewTransactionAuthenticator(auth *crypto.AccountAuthenticator) (*TransactionAuthenticator, error) {
	txnAuth := &TransactionAuthenticator{}
	switch auth.Variant {
	case crypto.AccountAuthenticatorEd25519:
		txnAuth.Variant = TransactionAuthenticatorEd25519
		txnAuth.Auth = &Ed25519TransactionAuthenticator{
			Sender: auth,
		}
	case crypto.AccountAuthenticatorMultiEd25519:
		txnAuth.Variant = TransactionAuthenticatorMultiEd25519
		txnAuth.Auth = &MultiEd25519TransactionAuthenticator{
			Sender: auth,
		}
	case crypto.AccountAuthenticatorSingleSender:
		txnAuth.Variant = TransactionAuthenticatorSingleSender
		txnAuth.Auth = &SingleSenderTransactionAuthenticator{
			Sender: auth,
		}
	case crypto.AccountAuthenticatorMultiKey:
		txnAuth.Variant = TransactionAuthenticatorSingleSender
		txnAuth.Auth = &SingleSenderTransactionAuthenticator{
			Sender: auth,
		}
	default:
		return nil, fmt.Errorf("unknown authenticator type %d", auth.Variant)
	}
	return txnAuth, nil
}

//region TransactionAuthenticator TransactionAuthenticatorImpl

func (ea *TransactionAuthenticator) Verify(msg []byte) bool {
	return ea.Auth.Verify(msg)
}

//endregion

//region TransactionAuthenticator bcs.Struct

func (ea *TransactionAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(ea.Variant))
	ea.Auth.MarshalBCS(ser)
}

func (ea *TransactionAuthenticator) UnmarshalBCS(des *bcs.Deserializer) {
	kindNum := des.Uleb128()
	if des.Error() != nil {
		return
	}
	ea.Variant = TransactionAuthenticatorVariant(kindNum)
	switch ea.Variant {
	case TransactionAuthenticatorEd25519:
		ea.Auth = &Ed25519TransactionAuthenticator{}
	case TransactionAuthenticatorMultiEd25519:
		ea.Auth = &MultiEd25519TransactionAuthenticator{}
	case TransactionAuthenticatorMultiAgent:
		ea.Auth = &MultiAgentTransactionAuthenticator{}
	case TransactionAuthenticatorFeePayer:
		ea.Auth = &FeePayerTransactionAuthenticator{}
	case TransactionAuthenticatorSingleSender:
		ea.Auth = &SingleSenderTransactionAuthenticator{}
	default:
		des.SetError(fmt.Errorf("unknown TransactionAuthenticator kind: %d", kindNum))
		return
	}
	ea.Auth.UnmarshalBCS(des)
}

//endregion
//endregion

//region Ed25519TransactionAuthenticator

// Ed25519TransactionAuthenticator for legacy ED25519 accounts
// Implements TransactionAuthenticatorImpl, bcs.Struct
type Ed25519TransactionAuthenticator struct {
	Sender *crypto.AccountAuthenticator
}

//region Ed25519TransactionAuthenticator TransactionAuthenticatorImpl

func (ea *Ed25519TransactionAuthenticator) Verify(msg []byte) bool {
	return ea.Sender.Verify(msg)
}

//endregion

//region Ed25519TransactionAuthenticator bcs.Struct

func (ea *Ed25519TransactionAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	ea.Sender.Auth.MarshalBCS(ser)
}

func (ea *Ed25519TransactionAuthenticator) UnmarshalBCS(des *bcs.Deserializer) {
	ea.Sender = &crypto.AccountAuthenticator{}
	ea.Sender.Variant = crypto.AccountAuthenticatorEd25519
	ea.Sender.Auth = &crypto.Ed25519Authenticator{}
	des.Struct(ea.Sender.Auth)
}

//endregion
//endregion

//region MultiEd25519TransactionAuthenticator

type MultiEd25519TransactionAuthenticator struct {
	Sender *crypto.AccountAuthenticator
}

//region Ed25519TransactionAuthenticator TransactionAuthenticatorImpl

func (ea *MultiEd25519TransactionAuthenticator) Verify(msg []byte) bool {
	return ea.Sender.Verify(msg)
}

//endregion

//region MultiEd25519TransactionAuthenticator bcs.Struct

func (ea *MultiEd25519TransactionAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	ea.Sender.MarshalBCS(ser)
}

func (ea *MultiEd25519TransactionAuthenticator) UnmarshalBCS(des *bcs.Deserializer) {
	//ea.Sender = &crypto.AccountAuthenticator{}
	//ea.Sender.Variant = crypto.AccountAuthenticatorMultiEd25519
	//ea.Sender.Auth = &crypto.MultiEd25519Authenticator{}
	//des.Struct(ea.Sender.Auth)
}

//endregion
//endregion

//region MultiAgentTransactionAuthenticator

type MultiAgentTransactionAuthenticator struct {
	Sender                   *crypto.AccountAuthenticator
	SecondarySignerAddresses []AccountAddress
	SecondarySigners         []crypto.AccountAuthenticator
}

//region MultiAgentTransactionAuthenticator TransactionAuthenticatorImpl

func (ea *MultiAgentTransactionAuthenticator) Verify(msg []byte) bool {
	sender := ea.Sender.Verify(msg)
	if !sender {
		return false
	}
	for _, sa := range ea.SecondarySigners {
		verified := sa.Verify(msg)
		if !verified {
			return false
		}
	}
	return true
}

//endregion

//region MultiAgentTransactionAuthenticator bcs.Struct

func (ea *MultiAgentTransactionAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	ea.Sender.MarshalBCS(ser)
	bcs.SerializeSequence(ea.SecondarySignerAddresses, ser)
	bcs.SerializeSequence(ea.SecondarySigners, ser)
}

func (ea *MultiAgentTransactionAuthenticator) UnmarshalBCS(des *bcs.Deserializer) {
	ea.Sender = &crypto.AccountAuthenticator{}
	des.Struct(ea.Sender)
	ea.SecondarySignerAddresses = bcs.DeserializeSequence[AccountAddress](des)
	ea.SecondarySigners = bcs.DeserializeSequence[crypto.AccountAuthenticator](des)
}

//endregion
//endregion

//region FeePayerTransactionAuthenticator

type FeePayerTransactionAuthenticator struct {
	Sender                   *crypto.AccountAuthenticator
	SecondarySignerAddresses []AccountAddress
	SecondarySigners         []crypto.AccountAuthenticator
	FeePayer                 *AccountAddress
	FeePayerAuthenticator    *crypto.AccountAuthenticator
}

//region FeePayerTransactionAuthenticator bcs.Struct

func (ea *FeePayerTransactionAuthenticator) Verify(msg []byte) bool {
	sender := ea.Sender.Verify(msg)
	if !sender {
		return false
	}
	for _, sa := range ea.SecondarySigners {
		verified := sa.Verify(msg)
		if !verified {
			return false
		}
	}
	return ea.FeePayerAuthenticator.Verify(msg)
}

//endregion

//region FeePayerTransactionAuthenticator bcs.Struct

func (ea *FeePayerTransactionAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	ea.Sender.MarshalBCS(ser)
	bcs.SerializeSequence(ea.SecondarySignerAddresses, ser)
	bcs.SerializeSequence(ea.SecondarySigners, ser)
	ser.Struct(ea.FeePayer)
	ser.Struct(ea.FeePayerAuthenticator)
}

func (ea *FeePayerTransactionAuthenticator) UnmarshalBCS(des *bcs.Deserializer) {
	ea.Sender = &crypto.AccountAuthenticator{}
	des.Struct(ea.Sender)
	ea.SecondarySignerAddresses = bcs.DeserializeSequence[AccountAddress](des)
	ea.SecondarySigners = bcs.DeserializeSequence[crypto.AccountAuthenticator](des)

	ea.FeePayer = &AccountAddress{}
	des.Struct(ea.FeePayer)
	ea.FeePayerAuthenticator = &crypto.AccountAuthenticator{}
	des.Struct(ea.FeePayerAuthenticator)
}

//endregion
//endregion

//region SingleSenderTransactionAuthenticator

type SingleSenderTransactionAuthenticator struct {
	Sender *crypto.AccountAuthenticator
}

//region SingleSenderTransactionAuthenticator TransactionAuthenticatorImpl

func (ea *SingleSenderTransactionAuthenticator) Verify(msg []byte) bool {
	return ea.Sender.Verify(msg)
}

//endregion

//region SingleSenderTransactionAuthenticator bcs.Struct

func (ea *SingleSenderTransactionAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(ea.Sender)
}

func (ea *SingleSenderTransactionAuthenticator) UnmarshalBCS(des *bcs.Deserializer) {
	ea.Sender = &crypto.AccountAuthenticator{}
	des.Struct(ea.Sender)
}

//endregion
//endregion
