package xdr

// NewDecoratedSignature constructs a decorated signature structure directly
// from the given signature and hint. Note that the hint should
// correspond to the signer that created the signature, but this helper cannot
// ensure that.
func NewDecoratedSignature(sig []byte, hint [4]byte) DecoratedSignature {
	return DecoratedSignature{
		Hint:      SignatureHint(hint),
		Signature: Signature(sig),
	}
}

// NewDecoratedSignatureForPayload creates a decorated signature with a hint
// that uses the key hint, the last four bytes of signature, and the last four
// bytes of the input that got signed. Note that the signature should be the
// signature of the payload via the key being hinted, but this construction
// method cannot ensure that.
func NewDecoratedSignatureForPayload(
	sig []byte, keyHint [4]byte, payload []byte,
) DecoratedSignature {
	hint := [4]byte{}
	// copy the last four bytes of the payload into the hint
	if len(payload) >= len(hint) {
		copy(hint[:], payload[len(payload)-len(hint):])
	} else {
		copy(hint[:], payload[:])
	}

	for i := 0; i < len(keyHint); i++ {
		hint[i] ^= keyHint[i]
	}

	return DecoratedSignature{
		Hint:      SignatureHint(hint),
		Signature: Signature(sig),
	}
}
