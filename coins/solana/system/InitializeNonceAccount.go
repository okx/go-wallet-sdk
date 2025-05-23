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

// Drive state of Uninitialized nonce account to Initialized, setting the nonce value
type InitializeNonceAccount struct {
	// The Pubkey parameter specifies the entity authorized to execute nonce instruction on the account.
	// No signatures are required to execute this instruction, enabling derived nonce account addresses.
	Authorized *base.PublicKey

	// [0] = [WRITE] NonceAccount
	// ··········· Nonce account
	//
	// [1] = [] $(SysVarRecentBlockHashesPubkey)
	// ··········· RecentBlockhashes sysvar
	//
	// [2] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewInitializeNonceAccountInstructionBuilder creates a new `InitializeNonceAccount` instruction builder.
func NewInitializeNonceAccountInstructionBuilder() *InitializeNonceAccount {
	nd := &InitializeNonceAccount{
		AccountMetaSlice: make(base.AccountMetaSlice, 3),
	}
	nd.AccountMetaSlice[1] = base.Meta(base.SysVarRecentBlockHashesPubkey)
	nd.AccountMetaSlice[2] = base.Meta(base.SysVarRentPubkey)
	return nd
}

// The Pubkey parameter specifies the entity authorized to execute nonce instruction on the account.
// No signatures are required to execute this instruction, enabling derived nonce account addresses.
func (inst *InitializeNonceAccount) SetAuthorized(authorized base.PublicKey) *InitializeNonceAccount {
	inst.Authorized = &authorized
	return inst
}

// Nonce account
func (inst *InitializeNonceAccount) SetNonceAccount(nonceAccount base.PublicKey) *InitializeNonceAccount {
	inst.AccountMetaSlice[0] = base.Meta(nonceAccount).WRITE()
	return inst
}

func (inst *InitializeNonceAccount) GetNonceAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// RecentBlockhashes sysvar
func (inst *InitializeNonceAccount) SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey base.PublicKey) *InitializeNonceAccount {
	inst.AccountMetaSlice[1] = base.Meta(SysVarRecentBlockHashesPubkey)
	return inst
}

func (inst *InitializeNonceAccount) GetSysVarRecentBlockHashesPubkeyAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// Rent sysvar
func (inst *InitializeNonceAccount) SetSysVarRentPubkeyAccount(SysVarRentPubkey base.PublicKey) *InitializeNonceAccount {
	inst.AccountMetaSlice[2] = base.Meta(SysVarRentPubkey)
	return inst
}

func (inst *InitializeNonceAccount) GetSysVarRentPubkeyAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst InitializeNonceAccount) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_InitializeNonceAccount, binary.LittleEndian),
	}}
}

func (inst InitializeNonceAccount) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Authorized` param:
	{
		err := encoder.Encode(*inst.Authorized)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewInitializeNonceAccountInstruction declares a new InitializeNonceAccount instruction with the provided parameters and accounts.
func NewInitializeNonceAccountInstruction(
	// Parameters:
	authorized base.PublicKey,
	// Accounts:
	nonceAccount base.PublicKey) *InitializeNonceAccount {
	return NewInitializeNonceAccountInstructionBuilder().
		SetAuthorized(authorized).
		SetNonceAccount(nonceAccount).
		SetSysVarRecentBlockHashesPubkeyAccount(base.SysVarRecentBlockHashesPubkey).
		SetSysVarRentPubkeyAccount(base.SysVarRentPubkey)
}
