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

// BumpSequence represents the Stellar bump sequence operation. See
// https://developers.stellar.org/docs/start/list-of-operations/
type BumpSequence struct {
	BumpTo        int64
	SourceAccount string
}

// Validate for BumpSequence validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (bs *BumpSequence) Validate() error {
	err := validateAmount(bs.BumpTo)
	if err != nil {
		return NewValidationError("BumpTo", err.Error())
	}
	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (bs *BumpSequence) GetSourceAccount() string {
	return bs.SourceAccount
}
