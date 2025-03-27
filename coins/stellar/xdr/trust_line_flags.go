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
