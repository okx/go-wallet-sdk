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

package token

import (
	"errors"
	"fmt"

	"github.com/emresenyuva/go-wallet-sdk/coins/solana/base"
)

// Mints new tokens to an account.  The native mint does not support minting.
//
// This instruction differs from MintTo in that the decimals value is
// checked by the caller.  This may be useful when creating transactions
// offline or within a hardware wallet.
type MintToChecked struct {
	// The amount of new tokens to mint.
	Amount *uint64

	// Expected number of base 10 digits to the right of the decimal place.
	Decimals *uint8

	// [0] = [WRITE] mint
	// ··········· The mint.
	//
	// [1] = [WRITE] destination
	// ··········· The account to mint tokens to.
	//
	// [2] = [] authority
	// ··········· The mint's minting authority.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *MintToChecked) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice MintToChecked) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewMintToCheckedInstructionBuilder creates a new `MintToChecked` instruction builder.
func NewMintToCheckedInstructionBuilder() *MintToChecked {
	nd := &MintToChecked{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of new tokens to mint.
func (inst *MintToChecked) SetAmount(amount uint64) *MintToChecked {
	inst.Amount = &amount
	return inst
}

// SetDecimals sets the "decimals" parameter.
// Expected number of base 10 digits to the right of the decimal place.
func (inst *MintToChecked) SetDecimals(decimals uint8) *MintToChecked {
	inst.Decimals = &decimals
	return inst
}

// SetMintAccount sets the "mint" account.
// The mint.
func (inst *MintToChecked) SetMintAccount(mint base.PublicKey) *MintToChecked {
	inst.Accounts[0] = base.Meta(mint).WRITE()
	return inst
}

// GetMintAccount gets the "mint" account.
// The mint.
func (inst *MintToChecked) GetMintAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetDestinationAccount sets the "destination" account.
// The account to mint tokens to.
func (inst *MintToChecked) SetDestinationAccount(destination base.PublicKey) *MintToChecked {
	inst.Accounts[1] = base.Meta(destination).WRITE()
	return inst
}

// GetDestinationAccount gets the "destination" account.
// The account to mint tokens to.
func (inst *MintToChecked) GetDestinationAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetAuthorityAccount sets the "authority" account.
// The mint's minting authority.
func (inst *MintToChecked) SetAuthorityAccount(authority base.PublicKey, multisigSigners ...base.PublicKey) *MintToChecked {
	inst.Accounts[2] = base.Meta(authority)
	if len(multisigSigners) == 0 {
		inst.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		inst.Signers = append(inst.Signers, base.Meta(signer).SIGNER())
	}
	return inst
}

// GetAuthorityAccount gets the "authority" account.
// The mint's minting authority.
func (inst *MintToChecked) GetAuthorityAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst MintToChecked) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_MintToChecked),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst MintToChecked) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *MintToChecked) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
		if inst.Decimals == nil {
			return errors.New("Decimals parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.Accounts[1] == nil {
			return errors.New("accounts.Destination is not set")
		}
		if inst.Accounts[2] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if !inst.Accounts[2].IsSigner && len(inst.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(inst.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(inst.Signers))
		}
	}
	return nil
}

func (obj MintToChecked) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(obj.Decimals)
	if err != nil {
		return err
	}
	return nil
}

// NewMintToCheckedInstruction declares a new MintToChecked instruction with the provided parameters and accounts.
func NewMintToCheckedInstruction(
	// Parameters:
	amount uint64,
	decimals uint8,
	// Accounts:
	mint base.PublicKey,
	destination base.PublicKey,
	authority base.PublicKey,
	multisigSigners []base.PublicKey,
) *MintToChecked {
	return NewMintToCheckedInstructionBuilder().
		SetAmount(amount).
		SetDecimals(decimals).
		SetMintAccount(mint).
		SetDestinationAccount(destination).
		SetAuthorityAccount(authority, multisigSigners...)
}
