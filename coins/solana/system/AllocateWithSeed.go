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

	"github.com/emresenyuva/go-wallet-sdk/coins/solana/base"
)

// Allocate space for and assign an account at an address derived from a base public key and a seed
type AllocateWithSeed struct {
	// Base public key
	Base *base.PublicKey

	// String of ASCII chars, no longer than pubkey::MAX_SEED_LEN
	Seed *string

	// Number of bytes of memory to allocate
	Space *uint64

	// Owner program account address
	Owner *base.PublicKey

	// [0] = [WRITE] AllocatedAccount
	// ··········· Allocated account
	//
	// [1] = [SIGNER] BaseAccount
	// ··········· Base account
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAllocateWithSeedInstructionBuilder creates a new `AllocateWithSeed` instruction builder.
func NewAllocateWithSeedInstructionBuilder() *AllocateWithSeed {
	nd := &AllocateWithSeed{
		AccountMetaSlice: make(base.AccountMetaSlice, 2),
	}
	return nd
}

// Base public key
func (inst *AllocateWithSeed) SetBase(base base.PublicKey) *AllocateWithSeed {
	inst.Base = &base
	return inst
}

// String of ASCII chars, no longer than pubkey::MAX_SEED_LEN
func (inst *AllocateWithSeed) SetSeed(seed string) *AllocateWithSeed {
	inst.Seed = &seed
	return inst
}

// Number of bytes of memory to allocate
func (inst *AllocateWithSeed) SetSpace(space uint64) *AllocateWithSeed {
	inst.Space = &space
	return inst
}

// Owner program account address
func (inst *AllocateWithSeed) SetOwner(owner base.PublicKey) *AllocateWithSeed {
	inst.Owner = &owner
	return inst
}

// Allocated account
func (inst *AllocateWithSeed) SetAllocatedAccount(allocatedAccount base.PublicKey) *AllocateWithSeed {
	inst.AccountMetaSlice[0] = base.Meta(allocatedAccount).WRITE()
	return inst
}

func (inst *AllocateWithSeed) GetAllocatedAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Base account
func (inst *AllocateWithSeed) SetBaseAccount(baseAccount base.PublicKey) *AllocateWithSeed {
	inst.AccountMetaSlice[1] = base.Meta(baseAccount).SIGNER()
	return inst
}

func (inst *AllocateWithSeed) GetBaseAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst AllocateWithSeed) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_AllocateWithSeed, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst AllocateWithSeed) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AllocateWithSeed) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Base == nil {
			return errors.New("Base parameter is not set")
		}
		if inst.Seed == nil {
			return errors.New("Seed parameter is not set")
		}
		if inst.Space == nil {
			return errors.New("Space parameter is not set")
		}
		if inst.Owner == nil {
			return errors.New("Owner parameter is not set")
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

func (inst AllocateWithSeed) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Base` param:
	{
		err := encoder.Encode(*inst.Base)
		if err != nil {
			return err
		}
	}
	// Serialize `Seed` param:
	{
		err := encoder.WriteRustString(*inst.Seed)
		if err != nil {
			return err
		}
	}
	// Serialize `Space` param:
	{
		err := encoder.Encode(*inst.Space)
		if err != nil {
			return err
		}
	}
	// Serialize `Owner` param:
	{
		err := encoder.Encode(*inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAllocateWithSeedInstruction declares a new AllocateWithSeed instruction with the provided parameters and accounts.
func NewAllocateWithSeedInstruction(
	// Parameters:
	base base.PublicKey,
	seed string,
	space uint64,
	owner base.PublicKey,
	// Accounts:
	allocatedAccount base.PublicKey,
	baseAccount base.PublicKey) *AllocateWithSeed {
	return NewAllocateWithSeedInstructionBuilder().
		SetBase(base).
		SetSeed(seed).
		SetSpace(space).
		SetOwner(owner).
		SetAllocatedAccount(allocatedAccount).
		SetBaseAccount(baseAccount)
}
