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

package xdr

import (
	"sort"
)

// SortSignersByKey returns a new []Signer array sorted by signer key.
func SortSignersByKey(signers []Signer) []Signer {
	keys := make([]string, 0, len(signers))
	keysMap := make(map[string]Signer)
	newSigners := make([]Signer, 0, len(signers))

	for _, signer := range signers {
		key := signer.Key.Address()
		keys = append(keys, key)
		keysMap[key] = signer
	}

	sort.Strings(keys)

	for _, key := range keys {
		newSigners = append(newSigners, keysMap[key])
	}

	return newSigners
}
