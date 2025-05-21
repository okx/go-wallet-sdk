package xdr

import (
	"fmt"
	"strconv"
	"strings"
)

// GoString prints Uint32 as decimal instead of hexadecimal numbers.
func (u Uint32) GoString() string {
	return strconv.FormatInt(int64(u), 10)
}

// GoString implements fmt.GoStringer.
func (e TransactionEnvelope) GoString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("xdr.TransactionEnvelope{Type: xdr.%s,", envelopeTypeMap[int32(e.Type)]))
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTxV0:
		sb.WriteString(fmt.Sprintf("V0: &%#v", *e.V0))
	case EnvelopeTypeEnvelopeTypeTx:
		sb.WriteString(fmt.Sprintf("V1: &%#v", *e.V1))
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		sb.WriteString(fmt.Sprintf("FeeBump: &%#v", *e.FeeBump))
	default:
		panic("Unknown type")
	}
	sb.WriteString("}")
	return sb.String()
}

// GoString implements fmt.GoStringer.
func (e FeeBumpTransactionInnerTx) GoString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("xdr.FeeBumpTransactionInnerTx{Type: xdr.%s,", envelopeTypeMap[int32(e.Type)]))
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTx:
		sb.WriteString(fmt.Sprintf("V1: &%#v", *e.V1))
	default:
		panic("Unknown type")
	}
	sb.WriteString("}")
	return sb.String()
}

// GoString implements fmt.GoStringer.
func (a AccountId) GoString() string {
	return fmt.Sprintf("xdr.MustAddress(%#v)", a.Address())
}

// GoString implements fmt.GoStringer.
func (a Asset) GoString() string {
	if a.Type == AssetTypeAssetTypeNative {
		return "xdr.MustNewNativeAsset()"
	}

	var typ, code, issuer string
	a.MustExtract(&typ, &code, &issuer)
	return fmt.Sprintf("xdr.MustNewCreditAsset(%#v, %#v)", code, issuer)
}

// GoString implements fmt.GoStringer.
func (m Memo) GoString() string {
	switch m.Type {
	case MemoTypeMemoNone:
		return fmt.Sprintf("xdr.Memo{Type: xdr.%s}", memoTypeMap[int32(m.Type)])
	case MemoTypeMemoText:
		return fmt.Sprintf(`xdr.MemoText(%#v)`, *m.Text)
	case MemoTypeMemoId:
		return fmt.Sprintf(`xdr.MemoID(%d)`, *m.Id)
	case MemoTypeMemoHash:
		return fmt.Sprintf(`xdr.MemoHash(%#v)`, *m.Hash)
	case MemoTypeMemoReturn:
		return fmt.Sprintf(`xdr.MemoRetHash(%#v)`, *m.RetHash)
	default:
		panic("Unknown type")
	}
}

// GoString implements fmt.GoStringer.
func (m MuxedAccount) GoString() string {
	switch m.Type {
	case CryptoKeyTypeKeyTypeEd25519:
		accountID := m.ToAccountId()
		return fmt.Sprintf("xdr.MustMuxedAddress(%#v)", accountID.Address())
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		var sb strings.Builder
		sb.WriteString("xdr.MuxedAccount{Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,")
		sb.WriteString(fmt.Sprintf("Med25519: &%#v", *m.Med25519))
		sb.WriteString("}")
		return sb.String()
	default:
		panic("Unknown type")
	}
}

// GoString implements fmt.GoStringer.
func (o Operation) GoString() string {
	var sb strings.Builder
	sb.WriteString("xdr.Operation{")
	if o.SourceAccount != nil {
		if o.SourceAccount.Type == CryptoKeyTypeKeyTypeEd25519 {
			accountID := o.SourceAccount.ToAccountId()
			sb.WriteString(fmt.Sprintf("SourceAccount: xdr.MustMuxedAddressPtr(%#v),", accountID.Address()))
		} else {
			sb.WriteString(fmt.Sprintf("SourceAccount: &%#v,", *o.SourceAccount))
		}
	}
	sb.WriteString(fmt.Sprintf("Body: %#v", o.Body))
	sb.WriteString("}")
	return sb.String()
}

// GoString implements fmt.GoStringer.
func (o OperationBody) GoString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("xdr.OperationBody{Type: xdr.%s,", operationTypeMap[int32(o.Type)]))
	switch {
	case o.PaymentOp != nil:
		sb.WriteString(fmt.Sprintf("PaymentOp: &%#v", *o.PaymentOp))
	case o.ChangeTrustOp != nil:
		sb.WriteString(fmt.Sprintf("ChangeTrustOp: &%#v", *o.ChangeTrustOp))
	case o.AllowTrustOp != nil:
		sb.WriteString(fmt.Sprintf("AllowTrustOp: &%#v", *o.AllowTrustOp))
	case o.Destination != nil:
		sb.WriteString(fmt.Sprintf("Destination: %#v", *o.Destination))
	default:
		panic("Unknown type")
	}
	sb.WriteString("}")
	return sb.String()
}

// GoString implements fmt.GoStringer.
func (s AssetCode) GoString() string {
	var code string
	switch s.Type {
	case AssetTypeAssetTypeCreditAlphanum4:
		code = string(s.AssetCode4[:])
	case AssetTypeAssetTypeCreditAlphanum12:
		code = string(s.AssetCode12[:])
	default:
		panic("Unknown type")
	}
	return fmt.Sprintf("xdr.MustNewAssetCodeFromString(%#v)", strings.TrimRight(code, string([]byte{0})))
}

// GoString implements fmt.GoStringer.
func (s Signer) GoString() string {
	var sb strings.Builder
	sb.WriteString("xdr.Signer{")
	sb.WriteString(fmt.Sprintf("Key: xdr.MustSigner(%#v),", s.Key.Address()))
	sb.WriteString(fmt.Sprintf("Weight: %#v", s.Weight))
	sb.WriteString("}")
	return sb.String()
}

// GoString implements fmt.GoStringer.
func (t *TimeBounds) GoString() string {
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("&xdr.TimeBounds{MinTime: xdr.TimePoint(%d), MaxTime: xdr.TimePoint(%d)}", t.MinTime, t.MaxTime)
}

// GoString implements fmt.GoStringer.
func (pt PreconditionType) GoString() string {
	return "xdr." + preconditionTypeMap[int32(pt)]
}

// GoString implements fmt.GoStringer.
func (p Preconditions) GoString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("xdr.Preconditions{Type: %s, ", p.Type.GoString()))
	switch p.Type {
	case PreconditionTypePrecondNone:
		sb.WriteString("nil")
	case PreconditionTypePrecondTime:
		sb.WriteString(fmt.Sprintf("TimeBounds: %s", p.TimeBounds.GoString()))
	case PreconditionTypePrecondV2:
		sb.WriteString(fmt.Sprintf("V2: {%#v}", p.V2))
	default:
		sb.WriteString("(unknown)")
	}

	sb.WriteString("}")
	return sb.String()
}
