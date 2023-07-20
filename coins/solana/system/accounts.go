// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"encoding/binary"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

type NonceAccount struct {
	Version          uint32
	State            uint32
	AuthorizedPubkey base.PublicKey
	Nonce            base.PublicKey
	FeeCalculator    FeeCalculator
}

type FeeCalculator struct {
	LamportsPerSignature uint64
}

func (obj NonceAccount) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	err = encoder.WriteUint32(obj.Version, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint32(obj.State, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteBytes(obj.AuthorizedPubkey[:], false)
	if err != nil {
		return err
	}
	err = encoder.WriteBytes(obj.Nonce[:], false)
	if err != nil {
		return err
	}
	return obj.FeeCalculator.MarshalWithEncoder(encoder)
}

func (obj FeeCalculator) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	return encoder.WriteUint64(obj.LamportsPerSignature, binary.LittleEndian)
}
