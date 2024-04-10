package bitcoin

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"regexp"
)

func GetPsbtFromString(psbtStr string) (*psbt.Packet, error) {
	isHex := IsHexString(psbtStr)
	var bs []byte
	var err error
	if isHex {
		bs, err = hex.DecodeString(psbtStr)
	} else {
		bs = []byte(psbtStr)
	}
	p, err := psbt.NewFromRawBytes(bytes.NewReader(bs), !isHex)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func IsHexString(s string) bool {
	if len(s) <= 1 {
		return false
	}
	if s[:2] != "0x" {
		s = "0x" + s
	}
	res, err := regexp.MatchString("^0x[0-9a-fA-F]+$", s)
	if err != nil {
		return false
	}
	return res
}

// NewPsbt compatible for bip0322
func NewPsbt(inputs []*wire.OutPoint,
	outputs []*wire.TxOut, version int32, nLockTime uint32,
	nSequences []uint32, opts ...string) (*psbt.Packet, error) {

	// Create the new struct; the input and output lists will be empty, the
	// unsignedTx object must be constructed and serialized, and that
	// serialization should be entered as the only entry for the
	// globalKVPairs list.
	//
	// Ensure that the version of the transaction is greater then our
	// minimum allowed transaction version. There must be one sequence
	// number per input.
	if version < psbt.MinTxVersion || len(nSequences) != len(inputs) {
		if version < psbt.MinTxVersion {
			if len(opts) == 0 || opts[0] != Bip0322Opt {
				return nil, psbt.ErrInvalidPsbtFormat
			}
		} else {
			return nil, psbt.ErrInvalidPsbtFormat
		}
	}

	unsignedTx := wire.NewMsgTx(version)
	unsignedTx.LockTime = nLockTime
	for i, in := range inputs {
		unsignedTx.AddTxIn(&wire.TxIn{
			PreviousOutPoint: *in,
			Sequence:         nSequences[i],
		})
	}
	for _, out := range outputs {
		unsignedTx.AddTxOut(out)
	}

	// The input and output lists are empty, but there is a list of those
	// two lists, and each one must be of length matching the unsigned
	// transaction; the unknown list can be nil.
	pInputs := make([]psbt.PInput, len(unsignedTx.TxIn))
	pOutputs := make([]psbt.POutput, len(unsignedTx.TxOut))

	// This new Psbt is "raw" and contains no key-value fields, so sanity
	// checking with c.Cpsbt.SanityCheck() is not required.
	return &psbt.Packet{
		UnsignedTx: unsignedTx,
		Inputs:     pInputs,
		Outputs:    pOutputs,
		Unknowns:   nil,
	}, nil
}
