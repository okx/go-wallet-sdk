/*
 * Copyright 2016 Stellar Development Foundation and contributors.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file includes portions of third-party code from [https://github.com/stellar/go].
 * The original code is licensed under the Apache License 2.0.
 */

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
