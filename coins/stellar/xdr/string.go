package xdr

func String32Ptr(val string) *String32 {
	pval := String32(val)
	return &pval
}
