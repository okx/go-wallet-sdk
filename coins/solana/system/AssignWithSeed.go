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

// Assign account to a program based on a seed
type AssignWithSeed struct {
	// Base public key
	Base *base.PublicKey

	// String of ASCII chars, no longer than pubkey::MAX_SEED_LEN
	Seed *string

	// Owner program account
	Owner *base.PublicKey

	// [0] = [WRITE] AssignedAccount
	// ··········· Assigned account
	//
	// [1] = [SIGNER] BaseAccount
	// ··········· Base account
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAssignWithSeedInstructionBuilder creates a new `AssignWithSeed` instruction builder.
func NewAssignWithSeedInstructionBuilder() *AssignWithSeed {
	nd := &AssignWithSeed{
		AccountMetaSlice: make(base.AccountMetaSlice, 2),
	}
	return nd
}

// Base public key
func (inst *AssignWithSeed) SetBase(base base.PublicKey) *AssignWithSeed {
	inst.Base = &base
	return inst
}

// String of ASCII chars, no longer than pubkey::MAX_SEED_LEN
func (inst *AssignWithSeed) SetSeed(seed string) *AssignWithSeed {
	inst.Seed = &seed
	return inst
}

// Owner program account
func (inst *AssignWithSeed) SetOwner(owner base.PublicKey) *AssignWithSeed {
	inst.Owner = &owner
	return inst
}

// Assigned account
func (inst *AssignWithSeed) SetAssignedAccount(assignedAccount base.PublicKey) *AssignWithSeed {
	inst.AccountMetaSlice[0] = base.Meta(assignedAccount).WRITE()
	return inst
}

func (inst *AssignWithSeed) GetAssignedAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Base account
func (inst *AssignWithSeed) SetBaseAccount(baseAccount base.PublicKey) *AssignWithSeed {
	inst.AccountMetaSlice[1] = base.Meta(baseAccount).SIGNER()
	return inst
}

func (inst *AssignWithSeed) GetBaseAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst AssignWithSeed) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_AssignWithSeed, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst AssignWithSeed) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AssignWithSeed) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Base == nil {
			return errors.New("Base parameter is not set")
		}
		if inst.Seed == nil {
			return errors.New("Seed parameter is not set")
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

func (inst AssignWithSeed) MarshalWithEncoder(encoder *base.Encoder) error {
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
	// Serialize `Owner` param:
	{
		err := encoder.Encode(*inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAssignWithSeedInstruction declares a new AssignWithSeed instruction with the provided parameters and accounts.
func NewAssignWithSeedInstruction(
	// Parameters:
	base base.PublicKey,
	seed string,
	owner base.PublicKey,
	// Accounts:
	assignedAccount base.PublicKey,
	baseAccount base.PublicKey) *AssignWithSeed {
	return NewAssignWithSeedInstructionBuilder().
		SetBase(base).
		SetSeed(seed).
		SetOwner(owner).
		SetAssignedAccount(assignedAccount).
		SetBaseAccount(baseAccount)
}
