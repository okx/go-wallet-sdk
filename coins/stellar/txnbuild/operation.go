package txnbuild

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
)

// Operation represents the operation types of the Stellar network.
type Operation interface {
	BuildXDR() (xdr.Operation, error)
	FromXDR(xdrOp xdr.Operation) error
	Validate() error
	GetSourceAccount() string
}

// SetOpSourceAccount sets the source account ID on an Operation, allowing M-strkeys (as defined in SEP23).
func SetOpSourceAccount(op *xdr.Operation, sourceAccount string) {
	if sourceAccount == "" {
		return
	}
	var opSourceAccountID xdr.MuxedAccount
	opSourceAccountID.SetAddress(sourceAccount)
	op.SourceAccount = &opSourceAccountID
}

// operationFromXDR returns a txnbuild Operation from its corresponding XDR operation
func operationFromXDR(xdrOp xdr.Operation) (Operation, error) {
	var newOp Operation
	switch xdrOp.Body.Type {
	case xdr.OperationTypePayment:
		newOp = &Payment{}
	case xdr.OperationTypeChangeTrust:
		newOp = &ChangeTrust{}
	case xdr.OperationTypeAllowTrust:
		newOp = &AllowTrust{}
	case xdr.OperationTypeAccountMerge:
		newOp = &AccountMerge{}
	case xdr.OperationTypeInflation:
		newOp = &Inflation{}
	default:
		return nil, fmt.Errorf("unknown operation type: %d", xdrOp.Body.Type)
	}

	err := newOp.FromXDR(xdrOp)
	return newOp, err
}

func accountFromXDR(account *xdr.MuxedAccount) string {
	if account != nil {
		return account.Address()
	}
	return ""
}

// SorobanOperation represents a smart contract operation on the Stellar network.
type SorobanOperation interface {
	BuildTransactionExt() (xdr.TransactionExt, error)
}
