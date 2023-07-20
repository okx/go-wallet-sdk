package types

import (
	"encoding/json"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/amino"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/common"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/common/types"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/tendermint"
)

var (
	_ types.Tx = (*StdTx)(nil)

	maxGasWanted = uint64((1 << 63) - 1)
)

// StdTx is a standard way to wrap a Msg with Fee and Signatures.
// NOTE: the first signature is the fee payer (Signatures must not be nil).
type StdTx struct {
	Msgs          []types.Msg    `json:"msg" yaml:"msg"`
	Fee           StdFee         `json:"fee" yaml:"fee"`
	Signatures    []StdSignature `json:"signatures" yaml:"signatures"`
	Memo          string         `json:"memo" yaml:"memo"`
	TimeoutHeight uint64         `json:"timeout_height" yaml:"timeout_height"`

	types.BaseTx `json:"-" yaml:"-"`
}

func NewStdTx(msgs []types.Msg, fee StdFee, sigs []StdSignature, memo string) *StdTx {
	return &StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}

type StdFee struct {
	Amount types.Coins `json:"amount" yaml:"amount"`
	Gas    uint64      `json:"gas" yaml:"gas"`
}

// NewStdFee returns a new instance of StdFee
func NewStdFee(gas uint64, amount types.Coins) StdFee {
	return StdFee{
		Amount: amount,
		Gas:    gas,
	}
}

// Bytes for signing later
func (fee *StdFee) Bytes() []byte {
	// normalize. XXX
	// this is a sign of something ugly
	// (in the lcd_test, client side its null,
	// server side its [])
	if len(fee.Amount) == 0 {
		fee.Amount = types.NewCoins()
	}
	bz, err := amino.GCodec.MarshalJSON(fee) // TODO
	if err != nil {
		panic(err)
	}
	return bz
}

//__________________________________________________________

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDoc struct {
	AccountNumber uint64            `json:"account_number" yaml:"account_number"`
	ChainID       string            `json:"chain_id" yaml:"chain_id"`
	Fee           json.RawMessage   `json:"fee" yaml:"fee"`
	Memo          string            `json:"memo" yaml:"memo"`
	Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
	Sequence      uint64            `json:"sequence" yaml:"sequence"`
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, accnum uint64, sequence uint64, fee StdFee, msgs []types.Msg, memo string) []byte {
	msgsBytes := make([]json.RawMessage, 0, len(msgs))
	for _, msg := range msgs {
		msgsBytes = append(msgsBytes, msg.GetSignBytes())
	}
	bz, err := amino.GCodec.MarshalJSON(StdSignDoc{
		AccountNumber: accnum,
		ChainID:       chainID,
		Fee:           fee.Bytes(),
		Memo:          memo,
		Msgs:          msgsBytes,
		Sequence:      sequence,
	})
	if err != nil {
		panic(err)
	}
	return common.MustSortJSON(bz)
}

// StdSignature represents a sig
type StdSignature struct {
	tendermint.PubKey `json:"pub_key" yaml:"pub_key"` // optional
	Signature         []byte                          `json:"signature" yaml:"signature"`
}
