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

// IsAuthorized returns true if issuer has authorized account to perform
// transactions with its credit
func (e TrustLineFlags) IsAuthorized() bool {
	return (e & TrustLineFlagsAuthorizedFlag) != 0
}

// IsAuthorizedToMaintainLiabilitiesFlag returns true if the issuer has authorized
// the account to maintain and reduce liabilities for its credit
func (e TrustLineFlags) IsAuthorizedToMaintainLiabilitiesFlag() bool {
	return (e & TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag) != 0
}

// IsClawbackEnabledFlag returns true if the issuer has authorized
// the account to claw assets back
func (e TrustLineFlags) IsClawbackEnabledFlag() bool {
	return (e & TrustLineFlagsTrustlineClawbackEnabledFlag) != 0
}
