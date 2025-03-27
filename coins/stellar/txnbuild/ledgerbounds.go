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

import "github.com/okx/go-wallet-sdk/coins/stellar/support/errors"

// LedgerBounds represent a transaction precondition that controls the ledger
// range for which a transaction is valid. Setting MaxLedger = 0 indicates there
// is no maximum ledger.
type LedgerBounds struct {
	MinLedger uint32
	MaxLedger uint32
}

func (lb *LedgerBounds) Validate() error {
	if lb == nil {
		return nil
	}

	if lb.MaxLedger > 0 && lb.MaxLedger < lb.MinLedger {
		return errors.New("invalid ledgerbound: max ledger < min ledger")
	}

	return nil
}
