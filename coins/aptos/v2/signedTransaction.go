package v2

import (
	"errors"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
)

// SignedTransaction a raw transaction plus its authenticator for a fully verifiable message
type SignedTransaction struct {
	Transaction   RawTransaction
	Authenticator crypto.Authenticator
}

func (txn *SignedTransaction) MarshalBCS(bcs *bcs.Serializer) {
	txn.Transaction.MarshalBCS(bcs)
	txn.Authenticator.MarshalBCS(bcs)
}
func (txn *SignedTransaction) UnmarshalBCS(bcs *bcs.Deserializer) {
	txn.Transaction.UnmarshalBCS(bcs)
	txn.Authenticator.UnmarshalBCS(bcs)
}

// Verify checks a signed transaction's signature
func (txn *SignedTransaction) Verify() error {
	bytes, err := txn.Transaction.SigningMessage()
	if err != nil {
		return err
	}
	if txn.Authenticator.Verify(bytes) {
		return nil
	}
	return errors.New("signature is invalid")
}

func (txn *SignedTransaction) Hash() ([]byte, error) {
	if TransactionPrefix == nil {
		hash := Sha3256Hash([][]byte{[]byte("APTOS::Transaction")})
		TransactionPrefix = &hash
	}

	txnBytes, err := bcs.Serialize(txn)
	if err != nil {
		return nil, err
	}

	// Transaction signature is defined as, the domain separated prefix based on struct (Transaction)
	// Then followed by the type of the transaction for the enum, UserTransaction is 0
	// Then followed by BCS encoded bytes of the signed transaction
	return Sha3256Hash([][]byte{*TransactionPrefix, {0}, txnBytes}), nil
}

var TransactionPrefix *[]byte
