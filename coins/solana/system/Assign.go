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

// Assign account to a program
type Assign struct {
	// Owner program account
	Owner *base.PublicKey

	// [0] = [WRITE, SIGNER] AssignedAccount
	// ··········· Assigned account public key
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAssignInstructionBuilder creates a new `Assign` instruction builder.
func NewAssignInstructionBuilder() *Assign {
	nd := &Assign{
		AccountMetaSlice: make(base.AccountMetaSlice, 1),
	}
	return nd
}

// Owner program account
func (inst *Assign) SetOwner(owner base.PublicKey) *Assign {
	inst.Owner = &owner
	return inst
}

// Assigned account public key
func (inst *Assign) SetAssignedAccount(assignedAccount base.PublicKey) *Assign {
	inst.AccountMetaSlice[0] = base.Meta(assignedAccount).WRITE().SIGNER()
	return inst
}

func (inst *Assign) GetAssignedAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

func (inst Assign) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_Assign, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Assign) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Assign) Validate() error {
	// Check whether all (required) parameters are set:
	{
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

func (inst Assign) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Owner` param:
	{
		err := encoder.Encode(*inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAssignInstruction declares a new Assign instruction with the provided parameters and accounts.
func NewAssignInstruction(
	// Parameters:
	owner base.PublicKey,
	// Accounts:
	assignedAccount base.PublicKey) *Assign {
	return NewAssignInstructionBuilder().
		SetOwner(owner).
		SetAssignedAccount(assignedAccount)
}
