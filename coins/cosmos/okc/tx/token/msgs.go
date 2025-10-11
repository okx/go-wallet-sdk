package token

import (
	"github.com/emresenyuva/go-wallet-sdk/coins/cosmos/okc/tx/amino"
	"github.com/emresenyuva/go-wallet-sdk/coins/cosmos/okc/tx/common"
	"github.com/emresenyuva/go-wallet-sdk/coins/cosmos/okc/tx/common/types"
)

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	FromAddress types.AccAddress `json:"from_address"`
	ToAddress   types.AccAddress `json:"to_address"`
	Amount      types.SysCoins   `json:"amount"`
}

func NewMsgTokenSend(from, to types.AccAddress, coins types.SysCoins) MsgSend {
	return MsgSend{
		FromAddress: from,
		ToAddress:   to,
		Amount:      coins,
	}
}

func (msg MsgSend) Route() string { return "token" }

func (msg MsgSend) Type() string { return "send" }

func (msg MsgSend) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return common.ErrAddressIsRequired
	}
	if msg.ToAddress.Empty() {
		return common.ErrAddressIsRequired
	}
	if !msg.Amount.IsValid() {
		return common.ErrInvalidCoins
	}
	if !msg.Amount.IsAllPositive() {
		return common.ErrInsufficientCoins
	}
	return nil
}

func (msg MsgSend) GetSignBytes() []byte {
	bz := amino.GCodec.MustMarshalJSON(msg)
	return common.MustSortJSON(bz)
}

func (msg MsgSend) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.FromAddress}
}
