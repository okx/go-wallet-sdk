package xdr

// Liabilities returns TrustLineEntry's liabilities
func (trustLine *TrustLineEntry) Liabilities() Liabilities {
	var liabilities Liabilities
	if trustLine.Ext.V1 != nil {
		liabilities = trustLine.Ext.V1.Liabilities
	}
	return liabilities
}
