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

const MAX_COMPUTE_UNIT_LIMIT = 1400000

type SetComputeUnitLimit struct {
	Units uint32
}

func (inst *SetComputeUnitLimit) SetAccounts(accounts []*base.AccountMeta) error {
	return nil
}

func (inst SetComputeUnitLimit) GetAccounts() (accounts []*base.AccountMeta) {
	return
}

// NewSetComputeUnitLimitInstructionBuilder creates a new `SetComputeUnitLimit` instruction builder.
func NewSetComputeUnitLimitInstructionBuilder() *SetComputeUnitLimit {
	nd := &SetComputeUnitLimit{}
	return nd
}

// Unit limit
func (inst *SetComputeUnitLimit) SetUnits(units uint32) *SetComputeUnitLimit {
	inst.Units = units
	return inst
}

func (inst SetComputeUnitLimit) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_SetComputeUnitLimit),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SetComputeUnitLimit) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetComputeUnitLimit) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Units == 0 {
			return errors.New("Units parameter is not set")
		}
		if inst.Units > MAX_COMPUTE_UNIT_LIMIT {
			return errors.New("Units parameter exceeds the maximum compute unit")
		}
	}
	return nil
}

func (obj SetComputeUnitLimit) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `Units` param:
	err = encoder.Encode(obj.Units)
	if err != nil {
		return err
	}
	return nil
}

// NewSetComputeUnitLimitInstruction declares a new SetComputeUnitLimit instruction with the provided parameters and accounts.
func NewSetComputeUnitLimitInstruction(
	// Parameters:
	units uint32,
) *SetComputeUnitLimit {
	return NewSetComputeUnitLimitInstructionBuilder().SetUnits(units)
}
