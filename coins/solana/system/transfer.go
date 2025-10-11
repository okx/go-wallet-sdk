package system

import (
	"encoding/binary"

	"github.com/emresenyuva/go-wallet-sdk/coins/solana/base"
)

// Transfer lamports
type Transfer struct {
	// Number of lamports to transfer to the new account
	Lamports *uint64

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE] RecipientAccount
	// ··········· Recipient account
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewTransferInstructionBuilder creates a new `Transfer` instruction builder.
func NewTransferInstructionBuilder() *Transfer {
	nd := &Transfer{
		AccountMetaSlice: make(base.AccountMetaSlice, 2),
	}
	return nd
}

// Number of lamports to transfer to the new account
func (inst *Transfer) SetLamports(lamports uint64) *Transfer {
	inst.Lamports = &lamports
	return inst
}

// Funding account
func (inst *Transfer) SetFundingAccount(fundingAccount base.PublicKey) *Transfer {
	inst.AccountMetaSlice[0] = base.Meta(fundingAccount).WRITE().SIGNER()
	return inst
}

// Recipient account
func (inst *Transfer) SetRecipientAccount(recipientAccount base.PublicKey) *Transfer {
	inst.AccountMetaSlice[1] = base.Meta(recipientAccount).WRITE()
	return inst
}

func (inst Transfer) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*inst.Lamports)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst Transfer) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_Transfer, binary.LittleEndian),
	}}
}

// NewTransferInstruction declares a new Transfer instruction with the provided parameters and accounts.
func NewTransferInstruction(
	// Parameters:
	lamports uint64,
	// Accounts:
	fundingAccount base.PublicKey,
	recipientAccount base.PublicKey) *Transfer {
	return NewTransferInstructionBuilder().
		SetLamports(lamports).
		SetFundingAccount(fundingAccount).
		SetRecipientAccount(recipientAccount)
}
