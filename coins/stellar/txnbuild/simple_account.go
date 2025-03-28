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

// SimpleAccount is a minimal implementation of an Account.
type SimpleAccount struct {
	AccountID string
	Sequence  int64
}

// GetAccountID returns the Account ID.
func (sa *SimpleAccount) GetAccountID() string {
	return sa.AccountID
}

// IncrementSequenceNumber increments the internal record of the
// account's sequence number by 1.
func (sa *SimpleAccount) IncrementSequenceNumber() (int64, error) {
	sa.Sequence++
	return sa.Sequence, nil
}

// GetSequenceNumber returns the sequence number of the account.
func (sa *SimpleAccount) GetSequenceNumber() (int64, error) {
	return sa.Sequence, nil
}

// NewSimpleAccount is a factory method that creates a SimpleAccount from "accountID" and "sequence".
func NewSimpleAccount(accountID string, sequence int64) SimpleAccount {
	return SimpleAccount{accountID, sequence}
}

// ensure that SimpleAccount implements Account interface.
var _ Account = &SimpleAccount{}
