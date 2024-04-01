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

package computebudget

import (
	"errors"
	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

type SetComputeUnitPrice struct {
	MicroLamports uint64
}

func (obj *SetComputeUnitPrice) SetAccounts(accounts []*base.AccountMeta) error {
	return nil
}

func (slice SetComputeUnitPrice) GetAccounts() (accounts []*base.AccountMeta) {
	return
}

// NewSetComputeUnitPriceInstructionBuilder creates a new `SetComputeUnitPrice` instruction builder.
func NewSetComputeUnitPriceInstructionBuilder() *SetComputeUnitPrice {
	nd := &SetComputeUnitPrice{}
	return nd
}

func (inst *SetComputeUnitPrice) SetMicroLamports(microLamports uint64) *SetComputeUnitPrice {
	inst.MicroLamports = microLamports
	return inst
}

func (inst SetComputeUnitPrice) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_SetComputeUnitPrice),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SetComputeUnitPrice) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetComputeUnitPrice) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.MicroLamports == 0 {
			return errors.New("MicroLamports parameter is not set")
		}
	}
	return nil
}

func (obj SetComputeUnitPrice) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `MicroLamports` param:
	err = encoder.Encode(obj.MicroLamports)
	if err != nil {
		return err
	}
	return nil
}

// NewSetComputeUnitPriceInstruction declares a new SetComputeUnitPrice instruction with the provided parameters and accounts.
func NewSetComputeUnitPriceInstruction(
	// Parameters:
	microLamports uint64,
) *SetComputeUnitPrice {
	return NewSetComputeUnitPriceInstructionBuilder().SetMicroLamports(microLamports)
}
