package system

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

const (
	// Create a new account
	Instruction_CreateAccount uint32 = iota

	// Assign account to a program
	Instruction_Assign

	// Transfer lamports
	Instruction_Transfer

	// Create a new account at an address derived from a base pubkey and a seed
	Instruction_CreateAccountWithSeed

	// Consumes a stored nonce, replacing it with a successor
	Instruction_AdvanceNonceAccount

	// Withdraw funds from a nonce account
	Instruction_WithdrawNonceAccount

	// Drive state of Uninitialized nonce account to Initialized, setting the nonce value
	Instruction_InitializeNonceAccount

	// Change the entity authorized to execute nonce instructions on the account
	Instruction_AuthorizeNonceAccount

	// Allocate space in a (possibly new) account without funding
	Instruction_Allocate

	// Allocate space for and assign an account at an address derived from a base public key and a seed
	Instruction_AllocateWithSeed

	// Assign account to a program based on a seed
	Instruction_AssignWithSeed

	// Transfer lamports from a derived address
	Instruction_TransferWithSeed
)

type Instruction struct {
	base.BaseVariant
}

func (inst *Instruction) ProgramID() base.PublicKey {
	return base.SystemProgramID
}

func (inst *Instruction) Accounts() (out []*base.AccountMeta) {
	return inst.Impl.(base.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := base.NewBinEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) MarshalWithEncoder(encoder *base.Encoder) error {
	err := encoder.WriteUint32(inst.TypeID.Uint32(), binary.LittleEndian)
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}
