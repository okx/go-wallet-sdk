package tendermint

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/amino"
)

// RegisterAmino registers all crypto related types in the given (amino) codec.
func RegisterCodec(cdc *amino.Codec) {
	// These are all written here instead of
	cdc.RegisterInterface((*PubKey)(nil), nil)
	cdc.RegisterConcrete(PubKeySecp256k1{}, PubKeyAminoName, nil)

	cdc.RegisterInterface((*PrivKey)(nil), nil)
	cdc.RegisterConcrete(PrivKeySecp256k1{}, PrivKeyAminoName, nil)
}

func init() {
	RegisterCodec(amino.GCodec)
}
