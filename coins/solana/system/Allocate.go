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
	"errors"
	"fmt"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

// Allocate space in a (possibly new) account without funding
type Allocate struct {
	// Number of bytes of memory to allocate
	Space *uint64

	// [0] = [WRITE, SIGNER] NewAccount
	// ··········· New account
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAllocateInstructionBuilder creates a new `Allocate` instruction builder.
func NewAllocateInstructionBuilder() *Allocate {
	nd := &Allocate{
		AccountMetaSlice: make(base.AccountMetaSlice, 1),
	}
	return nd
}

// Number of bytes of memory to allocate
func (inst *Allocate) SetSpace(space uint64) *Allocate {
	inst.Space = &space
	return inst
}

// New account
func (inst *Allocate) SetNewAccount(newAccount base.PublicKey) *Allocate {
	inst.AccountMetaSlice[0] = base.Meta(newAccount).WRITE().SIGNER()
	return inst
}

func (inst *Allocate) GetNewAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

func (inst Allocate) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_Allocate, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Allocate) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Allocate) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Space == nil {
			return errors.New("Space parameter is not set")
		}
	}

	// Check whether all accounts are set:
	for accIndex, acc := range inst.AccountMetaSlice {
		if acc == nil {
			return fmt.Errorf("ins.AccountMetaSlice[%v] is not set", accIndex)
		}
	}
	return nil
}

func (inst Allocate) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Space` param:
	{
		err := encoder.Encode(*inst.Space)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAllocateInstruction declares a new Allocate instruction with the provided parameters and accounts.
func NewAllocateInstruction(
	// Parameters:
	space uint64,
	// Accounts:
	newAccount base.PublicKey) *Allocate {
	return NewAllocateInstructionBuilder().
		SetSpace(space).
		SetNewAccount(newAccount)
}
