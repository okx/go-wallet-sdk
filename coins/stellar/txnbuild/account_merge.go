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
	"github.com/okx/go-wallet-sdk/coins/stellar/support/errors"
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
)

// AccountMerge represents the Stellar merge account operation. See
// https://developers.stellar.org/docs/start/list-of-operations/
type AccountMerge struct {
	Destination   string
	SourceAccount string
}

// BuildXDR for AccountMerge returns a fully configured XDR Operation.
func (am *AccountMerge) BuildXDR() (xdr.Operation, error) {
	var xdrOp xdr.MuxedAccount
	err := xdrOp.SetAddress(am.Destination)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set destination address")
	}

	opType := xdr.OperationTypeAccountMerge
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	SetOpSourceAccount(&op, am.SourceAccount)
	return op, nil
}

// FromXDR for AccountMerge initialises the txnbuild struct from the corresponding xdr Operation.
func (am *AccountMerge) FromXDR(xdrOp xdr.Operation) error {
	if xdrOp.Body.Type != xdr.OperationTypeAccountMerge {
		return errors.New("error parsing account_merge operation from xdr")
	}

	am.SourceAccount = accountFromXDR(xdrOp.SourceAccount)
	if xdrOp.Body.Destination != nil {
		am.Destination = xdrOp.Body.Destination.Address()
	}

	return nil
}

// Validate for AccountMerge validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (am *AccountMerge) Validate() error {
	var err error
	_, err = xdr.AddressToMuxedAccount(am.Destination)
	if err != nil {
		return NewValidationError("Destination", err.Error())
	}
	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (am *AccountMerge) GetSourceAccount() string {
	return am.SourceAccount
}
