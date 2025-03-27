package xdr

// TimeBounds extracts the timebounds (if any) from the transaction's
// Preconditions.
func (tx *Transaction) TimeBounds() *TimeBounds {
	switch tx.Cond.Type {
	case PreconditionTypePrecondNone:
		return nil
	case PreconditionTypePrecondTime:
		return tx.Cond.TimeBounds
	case PreconditionTypePrecondV2:
		return tx.Cond.V2.TimeBounds
	default:
		panic("unsupported precondition type: " + tx.Cond.Type.String())
	}
}

// LedgerBounds extracts the ledgerbounds (if any) from the transaction's
// Preconditions.
func (tx *Transaction) LedgerBounds() *LedgerBounds {
	switch tx.Cond.Type {
	case PreconditionTypePrecondNone, PreconditionTypePrecondTime:
		return nil
	case PreconditionTypePrecondV2:
		return tx.Cond.V2.LedgerBounds
	default:
		panic("unsupported precondition type: " + tx.Cond.Type.String())
	}
}

// MinSeqNum extracts the min seq number (if any) from the transaction's
// Preconditions.
func (tx *Transaction) MinSeqNum() *SequenceNumber {
	switch tx.Cond.Type {
	case PreconditionTypePrecondNone, PreconditionTypePrecondTime:
		return nil
	case PreconditionTypePrecondV2:
		return tx.Cond.V2.MinSeqNum
	default:
		panic("unsupported precondition type: " + tx.Cond.Type.String())
	}
}

// MinSeqAge extracts the min seq age (if any) from the transaction's
// Preconditions.
func (tx *Transaction) MinSeqAge() *Duration {
	switch tx.Cond.Type {
	case PreconditionTypePrecondNone, PreconditionTypePrecondTime:
		return nil
	case PreconditionTypePrecondV2:
		return &tx.Cond.V2.MinSeqAge
	default:
		panic("unsupported precondition type: " + tx.Cond.Type.String())
	}
}

// MinSeqLedgerGap extracts the min seq ledger gap (if any) from the transaction's
// Preconditions.
func (tx *Transaction) MinSeqLedgerGap() *Uint32 {
	switch tx.Cond.Type {
	case PreconditionTypePrecondNone, PreconditionTypePrecondTime:
		return nil
	case PreconditionTypePrecondV2:
		return &tx.Cond.V2.MinSeqLedgerGap
	default:
		panic("unsupported precondition type: " + tx.Cond.Type.String())
	}
}

// ExtraSigners extracts the extra signers (if any) from the transaction's
// Preconditions.
func (tx *Transaction) ExtraSigners() []SignerKey {
	switch tx.Cond.Type {
	case PreconditionTypePrecondNone, PreconditionTypePrecondTime:
		return nil
	case PreconditionTypePrecondV2:
		return tx.Cond.V2.ExtraSigners
	default:
		panic("unsupported precondition type: " + tx.Cond.Type.String())
	}
}
