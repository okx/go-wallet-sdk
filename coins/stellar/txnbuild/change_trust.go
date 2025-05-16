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
	"github.com/okx/go-wallet-sdk/coins/stellar/amount"
	"github.com/okx/go-wallet-sdk/coins/stellar/support/errors"
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
	"math"
)

// ChangeTrust represents the Stellar change trust operation. See
// https://developers.stellar.org/docs/start/list-of-operations/
// If Limit is omitted, it defaults to txnbuild.MaxTrustlineLimit.
type ChangeTrust struct {
	Line          ChangeTrustAsset
	Limit         string
	SourceAccount string
}

// MaxTrustlineLimit represents the maximum value that can be set as a trustline limit.
var MaxTrustlineLimit = amount.StringFromInt64(math.MaxInt64)

// RemoveTrustlineOp returns a ChangeTrust operation to remove the trustline of the described asset,
// by setting the limit to "0".
func RemoveTrustlineOp(issuedAsset ChangeTrustAsset) ChangeTrust {
	return ChangeTrust{
		Line:  issuedAsset,
		Limit: "0",
	}
}

// BuildXDR for ChangeTrust returns a fully configured XDR Operation.
func (ct *ChangeTrust) BuildXDR() (xdr.Operation, error) {
	if ct.Line.IsNative() {
		return xdr.Operation{}, errors.New("trustline cannot be extended to a native (XLM) asset")
	}
	xdrLine, err := ct.Line.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "can't convert trustline asset to XDR")
	}

	if ct.Limit == "" {
		ct.Limit = MaxTrustlineLimit
	}

	xdrLimit, err := amount.Parse(ct.Limit)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse limit amount")
	}

	opType := xdr.OperationTypeChangeTrust
	xdrOp := xdr.ChangeTrustOp{
		Line:  xdrLine,
		Limit: xdrLimit,
	}
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	SetOpSourceAccount(&op, ct.SourceAccount)
	return op, nil
}

// FromXDR for ChangeTrust initialises the txnbuild struct from the corresponding xdr Operation.
func (ct *ChangeTrust) FromXDR(xdrOp xdr.Operation) error {
	result, ok := xdrOp.Body.GetChangeTrustOp()
	if !ok {
		return errors.New("error parsing change_trust operation from xdr")
	}

	ct.SourceAccount = accountFromXDR(xdrOp.SourceAccount)
	ct.Limit = amount.String(result.Limit)
	asset, err := assetFromChangeTrustAssetXDR(result.Line)
	if err != nil {
		return errors.Wrap(err, "error parsing asset in change_trust operation")
	}
	ct.Line = asset
	return nil
}

// Validate for ChangeTrust validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (ct *ChangeTrust) Validate() error {
	// only validate limit if it has a value. Empty limit is set to the max trustline limit.
	if ct.Limit != "" {
		err := validateAmount(ct.Limit)
		if err != nil {
			return NewValidationError("Limit", err.Error())
		}
	}

	err := validateChangeTrustAsset(ct.Line)
	if err != nil {
		return NewValidationError("Line", err.Error())
	}
	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (ct *ChangeTrust) GetSourceAccount() string {
	return ct.SourceAccount
}
