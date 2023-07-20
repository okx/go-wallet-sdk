package token

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/amino"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *amino.Codec) {
	//cdc.RegisterConcrete(MsgTokenIssue{}, "okexchain/token/MsgIssue", nil)
	//cdc.RegisterConcrete(MsgTokenBurn{}, "okexchain/token/MsgBurn", nil)
	//cdc.RegisterConcrete(MsgTokenMint{}, "okexchain/token/MsgMint", nil)
	//cdc.RegisterConcrete(MsgMultiSend{}, "okexchain/token/MsgMultiTransfer", nil)
	cdc.RegisterConcrete(MsgSend{}, "okexchain/token/MsgTransfer", nil)
	//cdc.RegisterConcrete(MsgTransferOwnership{}, "okexchain/token/MsgTransferOwnership", nil)
	//cdc.RegisterConcrete(MsgConfirmOwnership{}, "okexchain/token/MsgConfirmOwnership", nil)
	//cdc.RegisterConcrete(MsgTokenModify{}, "okexchain/token/MsgModify", nil)
}

func init() {
	RegisterCodec(amino.GCodec)
}
