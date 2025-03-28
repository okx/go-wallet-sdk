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

package strkey

import (
	"bytes"
	"github.com/okx/go-wallet-sdk/coins/stellar/support/errors"
	xdr "github.com/okx/go-wallet-sdk/coins/stellar/xdr3"
)

type SignedPayload struct {
	signer  string
	payload []byte
}

const maxPayloadLen = 64

// NewSignedPayload creates a signed payload from an account ID (G... address)
// and a payload. The payload buffer is copied directly into the structure, so
// it should not be modified after construction.
func NewSignedPayload(signerPublicKey string, payload []byte) (*SignedPayload, error) {
	if len(payload) > maxPayloadLen {
		return nil, errors.Errorf("payload length %d exceeds max %d",
			len(payload), maxPayloadLen)
	}

	return &SignedPayload{signer: signerPublicKey, payload: payload}, nil
}

// Encode turns a signed payload structure into its StrKey equivalent.
func (sp *SignedPayload) Encode() (string, error) {
	signerBytes, err := Decode(VersionByteAccountID, sp.Signer())
	if err != nil {
		return "", errors.Wrap(err, "failed to decode signed payload signer")
	}

	b := new(bytes.Buffer)
	b.Write(signerBytes)
	xdr.Marshal(b, sp.Payload())

	strkey, err := Encode(VersionByteSignedPayload, b.Bytes())
	if err != nil {
		return "", errors.Wrap(err, "failed to encode signed payload")
	}
	return strkey, nil
}

func (sp *SignedPayload) Signer() string {
	return sp.signer
}

func (sp *SignedPayload) Payload() []byte {
	return sp.payload
}

// DecodeSignedPayload transforms a P... signer into a `SignedPayload` instance.
func DecodeSignedPayload(address string) (*SignedPayload, error) {
	raw, err := Decode(VersionByteSignedPayload, address)
	if err != nil {
		return nil, errors.New("invalid signed payload")
	}

	const signerLen = 32
	rawSigner, raw := raw[:signerLen], raw[signerLen:]
	signer, err := Encode(VersionByteAccountID, rawSigner)
	if err != nil {
		return nil, errors.Wrap(err, "invalid signed payload signer")
	}

	payload := []byte{}
	reader := bytes.NewBuffer(raw)
	readBytes, err := xdr.Unmarshal(reader, &payload)
	if err != nil {
		return nil, errors.Wrap(err, "invalid signed payload")
	}

	if len(raw) != readBytes || reader.Len() > 0 {
		return nil, errors.New("invalid signed payload padding")
	}

	return NewSignedPayload(signer, payload)
}
