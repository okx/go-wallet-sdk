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
