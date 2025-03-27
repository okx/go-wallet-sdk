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
