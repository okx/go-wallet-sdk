// This file is referenced from https://github.com/decred/dcrdex/blob/281e6bc/dex/networks/zec/tx.go
// We updated the Tx implementation to read consensus parameters (such as branch ID) from the transaction itself
// instead of hardcoding them, as these values may change during chain hard forks.

package zec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/okx/go-wallet-sdk/crypto/blake2b"
	"io"
)

const (
	VersionPreOverwinter = 2
	VersionOverwinter    = 3
	VersionSapling       = 4
	VersionNU5           = 5
	MaxExpiryHeight      = 499999999 // https://zips.z.cash/zip-0203

	versionOverwinterGroupID = 0x03C48270
	versionSaplingGroupID    = 0x892f2085
	versionNU5GroupID        = 0x26A7270A

	overwinterMask = ^uint32(1 << 31)
	pver           = 0

	overwinterJoinSplitSize = 1802
	saplingJoinSplitSize    = 1698

	pkTransparentDigest = "ZTxIdTranspaHash"
	pkPrevOutsV5        = "ZTxIdPrevoutHash"
	pkAmounts           = "ZTxTrAmountsHash"
	pkPrevScripts       = "ZTxTrScriptsHash"
	pkSequenceV5        = "ZTxIdSequencHash"
	pkOutputsV5         = "ZTxIdOutputsHash"
	pkTxIn              = "Zcash___TxInHash"
	pkV5TxDigest        = "ZcashTxHash_"
	pkHeader            = "ZTxIdHeadersHash"
	pkPrevOutsV4        = "ZcashPrevoutHash"
	pkSequenceV4        = "ZcashSequencHash"
	pkOutputsV4         = "ZcashOutputsHash"
)

var (
	// Zclassic only
	ConsensusBranchButtercup = [4]byte{0x0d, 0x54, 0x0b, 0x93}

	emptySaplingDigest = [32]byte{0x6f, 0x2f, 0xc8, 0xf9, 0x8f, 0xea, 0xfd, 0x94,
		0xe7, 0x4a, 0x0d, 0xf4, 0xbe, 0xd7, 0x43, 0x91, 0xee, 0x0b, 0x5a, 0x69,
		0x94, 0x5e, 0x4c, 0xed, 0x8c, 0xa8, 0xa0, 0x95, 0x20, 0x6f, 0x00, 0xae}

	emptyOrchardDigest = [32]byte{0x9f, 0xbe, 0x4e, 0xd1, 0x3b, 0x0c, 0x08, 0xe6,
		0x71, 0xc1, 0x1a, 0x34, 0x07, 0xd8, 0x4e, 0x11, 0x17, 0xcd, 0x45, 0x02,
		0x8a, 0x2e, 0xee, 0x1b, 0x9f, 0xea, 0xe7, 0x8b, 0x48, 0xa6, 0xe2, 0xc1}
)

// JoinSplit is only the new and old fields of a vJoinSplit.
type JoinSplit struct {
	Old, New uint64
}

// Tx is a Zcash-adapted MsgTx. Tx will decode any version transaction, but will
// not save most data for shielded transactions.
// Tx can only produce tx hashes for unshielded transactions. Tx can only create
// signature hashes for unshielded version 5 transactions.
type Tx struct {
	*wire.MsgTx
	ConsensusBranchId   uint32
	ExpiryHeight        uint32
	NSpendsSapling      uint64
	NOutputsSapling     uint64
	ValueBalanceSapling int64
	NActionsOrchard     uint64
	SizeProofsOrchard   uint64
	NJoinSplit          uint64
	VJoinSplit          []*JoinSplit
	ValueBalanceOrchard int64
}

// NewTxFromMsgTx creates a Tx embedding the MsgTx, and adding Zcash-specific
// fields.
func NewTxFromMsgTx(tx *wire.MsgTx, expiryHeight uint32) *Tx {
	zecTx := &Tx{
		MsgTx:        tx,
		ExpiryHeight: expiryHeight,
	}
	return zecTx
}

// TxHash generates the Hash for the transaction.
func (tx *Tx) TxHash() chainhash.Hash {
	if tx.Version == 5 {
		txHash, err := tx.txHashV5()
		if err != nil {
			return chainhash.Hash{}
		}
		return txHash
	}
	b, _ := tx.Bytes()
	return chainhash.DoubleHashH(b)
}

func (tx *Tx) txHashV5() (_ chainhash.Hash, err error) {
	td, err := tx.transparentDigestV5()
	if err != nil {
		return
	}
	return tx.txDigestV5(td)
}

// SignatureDigest produces a hash of tx data suitable for signing.
// SignatureDigest only works correctly for unshielded version 5 transactions.
func (tx *Tx) SignatureDigest(
	vin int, hashType txscript.SigHashType, script []byte, vals []int64, prevScripts [][]byte,
) (_ [32]byte, err error) {

	if tx.Version == 4 {
		return tx.txDigestV4(hashType, vin, vals, script)
	}
	td, err := tx.transparentSigDigestV5(vin, hashType, vals, prevScripts)
	if err != nil {
		return
	}
	return tx.txDigestV5(td)
}

// txDigestV5 produces hashes of transaction data in accordance with ZIP-244.
func (tx *Tx) txDigestV5(transparentPart [32]byte) (_ chainhash.Hash, err error) {
	hd, err := tx.headerDigestV5()
	if err != nil {
		return
	}
	b := make([]byte, 128)

	copy(b[:32], hd[:])
	copy(b[32:64], transparentPart[:])
	copy(b[64:96], emptySaplingDigest[:])
	copy(b[96:], emptyOrchardDigest[:])
	h, err := blake2bHash(b, append([]byte(pkV5TxDigest), uint32Bytes(tx.ConsensusBranchId)[:]...))
	if err != nil {
		return
	}
	var txHash chainhash.Hash
	copy(txHash[:], h[:])
	return txHash, nil
}

func (tx *Tx) headerDigestV5() ([32]byte, error) {
	b := make([]byte, 20)
	copy(b[:4], uint32Bytes(uint32(tx.Version)|(1<<31)))
	copy(b[4:8], uint32Bytes(versionNU5GroupID))
	copy(b[8:12], uint32Bytes(tx.ConsensusBranchId))
	copy(b[12:16], uint32Bytes(tx.LockTime))
	copy(b[16:], uint32Bytes(tx.ExpiryHeight))
	return blake2bHash(b, []byte(pkHeader))
}

// Zclassic only. Based on ZIP-0243, but uses ConsensusBranchButtercup from
// Zclassic.
// isSimnet is hack for https://github.com/ZclassicCommunity/zclassic/issues/83
func (tx *Tx) txDigestV4(
	hashType txscript.SigHashType, vin int, vals []int64, script []byte,
) (_ chainhash.Hash, err error) {

	b, err := tx.sighashPreimageV4(hashType, vin, vals, script)
	if err != nil {
		return
	}
	consensusBranchID := ConsensusBranchButtercup
	h, err := blake2bHash(b, append([]byte("ZcashSigHash"), consensusBranchID[:]...))
	if err != nil {
		return
	}

	var txHash chainhash.Hash
	copy(txHash[:], h[:])
	return txHash, nil
}

func (tx *Tx) sighashPreimageV4(hashType txscript.SigHashType, vin int, vals []int64, script []byte) (_ []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 270+len(script)))
	buf.Write(uint32Bytes(uint32(tx.Version) | (1 << 31))) // 4 bytes
	buf.Write(uint32Bytes(versionSaplingGroupID))          // + 4 = 8
	prevoutsDigest, err := tx.calcHashPrevOuts(pkPrevOutsV4)
	if err != nil {
		return
	}
	buf.Write(prevoutsDigest[:]) // + 32 = 40

	seqDigest, err := tx.hashSequence(pkSequenceV4)
	if err != nil {
		return
	}
	buf.Write(seqDigest[:]) // + 32 = 72

	outputsDigest, err := tx.hashOutputs(pkOutputsV4)
	if err != nil {
		return
	}
	buf.Write(outputsDigest[:]) // + 32 = 104
	// The following three fields are all zero hashes for transparent txs.
	// hashJoinSplits, hashShieldedSpends, hashShieldedOutputs [32]byte
	buf.Write(make([]byte, 96))             // + 96 = 200
	buf.Write(uint32Bytes(tx.LockTime))     // + 4 = 204
	buf.Write(uint32Bytes(tx.ExpiryHeight)) // + 4 = 208
	// valueBalance
	buf.Write(uint64Bytes(0))                // + 8 = 216
	buf.Write(uint32Bytes(uint32(hashType))) // + 4 = 220

	txInsDigest, err := tx.preimageTxInSig(vin, vals[vin], script) // + 50 + len(prevScript) = 270 + len(prevScript)
	if err != nil {
		return
	}
	buf.Write(txInsDigest[:])
	return buf.Bytes(), nil
}

func (tx *Tx) transparentDigestV5() (h [32]byte, err error) {
	prevoutsDigest, err := tx.calcHashPrevOuts(pkPrevOutsV5)
	if err != nil {
		return
	}

	seqDigest, err := tx.hashSequence(pkSequenceV5)
	if err != nil {
		return
	}

	outputsDigest, err := tx.hashOutputs(pkOutputsV5)
	if err != nil {
		return
	}

	b := make([]byte, 96)
	copy(b[:32], prevoutsDigest[:])
	copy(b[32:64], seqDigest[:])
	copy(b[64:], outputsDigest[:])
	return blake2bHash(b, []byte(pkTransparentDigest))
}

func (tx *Tx) transparentSigDigestV5(vin int, hashType txscript.SigHashType, vals []int64, prevScripts [][]byte) (h [32]byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 193))

	buf.Write([]byte{byte(hashType)})

	anyoneCanPay := hashType&txscript.SigHashAnyOneCanPay > 0

	prevoutsDigest, err := tx.hashPrevOutsSigV5(anyoneCanPay)
	if err != nil {
		return
	}
	buf.Write(prevoutsDigest[:])

	amtsDigest, err := tx.hashAmountsSig(anyoneCanPay, vals)
	if err != nil {
		return
	}
	buf.Write(amtsDigest[:])

	prevScriptsDigest, err := tx.hashPrevScriptsSig(anyoneCanPay, prevScripts)
	if err != nil {
		return
	}
	buf.Write(prevScriptsDigest[:])

	seqsDigest, err := tx.hashSequenceSigV5(anyoneCanPay)
	if err != nil {
		return
	}
	buf.Write(seqsDigest[:])

	outputsDigest, err := tx.hashOutputsSigV5(anyoneCanPay)
	if err != nil {
		return
	}
	buf.Write(outputsDigest[:])

	txInsDigest, err := tx.hashTxInSig(vin, vals[vin], prevScripts[vin])
	if err != nil {
		return
	}
	buf.Write(txInsDigest[:])

	return blake2bHash(buf.Bytes(), []byte(pkTransparentDigest))
}

func (tx *Tx) calcHashPrevOuts(pk string) ([32]byte, error) {
	var buf bytes.Buffer
	for _, in := range tx.TxIn {
		buf.Write(in.PreviousOutPoint.Hash[:])
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], in.PreviousOutPoint.Index)
		buf.Write(b[:])
	}
	return blake2bHash(buf.Bytes(), []byte(pk))
}

func (tx *Tx) hashPrevOutsSigV5(anyoneCanPay bool) ([32]byte, error) {
	if anyoneCanPay {
		return blake2bHash([]byte{}, []byte(pkPrevOutsV5))
	}
	return tx.calcHashPrevOuts(pkPrevOutsV5)
}

func (tx *Tx) hashAmountsSig(anyoneCanPay bool, vals []int64) ([32]byte, error) {
	if anyoneCanPay {
		return blake2bHash([]byte{}, []byte(pkAmounts))
	}
	b := make([]byte, 0, 8*len(vals))
	for _, v := range vals {
		b = append(b, int64Bytes(v)...)
	}
	return blake2bHash(b, []byte(pkAmounts))
}

func (tx *Tx) hashPrevScriptsSig(anyoneCanPay bool, prevScripts [][]byte) (_ [32]byte, err error) {
	if anyoneCanPay {
		return blake2bHash([]byte{}, []byte(pkPrevScripts))
	}
	buf := new(bytes.Buffer)
	for _, s := range prevScripts {
		if err = wire.WriteVarBytes(buf, pver, s); err != nil {
			return
		}
	}

	return blake2bHash(buf.Bytes(), []byte(pkPrevScripts))
}

func (tx *Tx) hashSequence(pk string) ([32]byte, error) {
	var b bytes.Buffer
	for _, in := range tx.TxIn {
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], in.Sequence)
		b.Write(buf[:])
	}
	return blake2bHash(b.Bytes(), []byte(pk))
}

func (tx *Tx) hashSequenceSigV5(anyoneCanPay bool) (h [32]byte, err error) {
	if anyoneCanPay {
		return blake2bHash([]byte{}, []byte(pkSequenceV5))
	}
	return tx.hashSequence(pkSequenceV5)
}

func (tx *Tx) hashOutputs(pk string) (_ [32]byte, err error) {
	var b bytes.Buffer
	for _, out := range tx.TxOut {
		if err = wire.WriteTxOut(&b, 0, 0, out); err != nil {
			return chainhash.Hash{}, err
		}
	}
	return blake2bHash(b.Bytes(), []byte(pk))
}

func (tx *Tx) hashOutputsSigV5(anyoneCanPay bool) (_ [32]byte, err error) {
	if anyoneCanPay {
		return blake2bHash([]byte{}, []byte(pkOutputsV5))
	}
	return tx.hashOutputs(pkOutputsV5)
}

func (tx *Tx) hashTxInSig(idx int, prevVal int64, prevScript []byte) (h [32]byte, err error) {
	b, err := tx.preimageTxInSig(idx, prevVal, prevScript)
	if err != nil {
		return
	}
	return blake2bHash(b, []byte(pkTxIn))
}

func (tx *Tx) preimageTxInSig(idx int, prevVal int64, script []byte) ([]byte, error) {
	if len(tx.TxIn) <= idx {
		return nil, fmt.Errorf("no input at index %d", idx)
	}
	txIn := tx.TxIn[idx]

	// prev_hash 32 + prev_index 4 + prev_value 8 + script_pub_key var_int+L (size) + nSequence 4
	b := bytes.NewBuffer(make([]byte, 0, 50+len(script)))

	b.Write(txIn.PreviousOutPoint.Hash[:])
	b.Write(uint32Bytes(txIn.PreviousOutPoint.Index))
	if tx.Version >= 5 {
		b.Write(int64Bytes(prevVal))
	}
	if err := wire.WriteVarBytes(b, pver, script); err != nil {
		return nil, err
	}
	if tx.Version < 5 {
		b.Write(int64Bytes(prevVal))
	}
	b.Write(uint32Bytes(txIn.Sequence))

	return b.Bytes(), nil
}

// Bytes encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
// See Serialize for encoding transactions to be stored to disk, such as in a
// database, as opposed to encoding transactions for the wire.
// msg.Version must be 4 or 5.
func (tx *Tx) Bytes() (_ []byte, err error) {
	w := new(bytes.Buffer)
	header := uint32(tx.Version)
	if tx.Version >= VersionOverwinter {
		header |= 1 << 31
	}

	if err = putUint32(w, header); err != nil {
		return nil, fmt.Errorf("error writing version: %w", err)
	}

	if tx.Version >= VersionOverwinter {
		var groupID uint32 = versionOverwinterGroupID
		switch tx.Version {
		case VersionSapling:
			groupID = versionSaplingGroupID
		case VersionNU5:
			groupID = versionNU5GroupID
		}

		// nVersionGroupId
		if err = putUint32(w, groupID); err != nil {
			return nil, fmt.Errorf("error writing nVersionGroupId: %w", err)
		}
	}

	if tx.Version == VersionNU5 {
		// nConsensusBranchId
		if err = putUint32(w, tx.ConsensusBranchId); err != nil {
			return nil, fmt.Errorf("error writing nConsensusBranchId: %w", err)
		}

		// lock_time
		if err = putUint32(w, tx.LockTime); err != nil {
			return nil, fmt.Errorf("error writing lock_time: %w", err)
		}

		// nExpiryHeight
		if err = putUint32(w, tx.ExpiryHeight); err != nil {
			return nil, fmt.Errorf("error writing nExpiryHeight: %w", err)
		}
	}

	// tx_in_count
	if err = wire.WriteVarInt(w, pver, uint64(len(tx.MsgTx.TxIn))); err != nil {
		return nil, fmt.Errorf("error writing tx_in_count: %w", err)
	}

	// tx_in
	for vin, ti := range tx.TxIn {
		if err = writeTxIn(w, ti); err != nil {
			return nil, fmt.Errorf("error writing tx_in %d: %w", vin, err)
		}
	}

	// tx_out_count
	if err = wire.WriteVarInt(w, pver, uint64(len(tx.TxOut))); err != nil {
		return nil, fmt.Errorf("error writing tx_out_count: %w", err)
	}

	// tx_out
	for vout, to := range tx.TxOut {
		if err = wire.WriteTxOut(w, pver, tx.Version, to); err != nil {
			return nil, fmt.Errorf("error writing tx_out %d: %w", vout, err)
		}
	}

	if tx.Version <= VersionSapling {
		// lock_time
		if err = putUint32(w, tx.LockTime); err != nil {
			return nil, fmt.Errorf("error writing lock_time: %w", err)
		}

		if tx.Version >= VersionOverwinter {
			// nExpiryHeight
			if err = putUint32(w, tx.ExpiryHeight); err != nil {
				return nil, fmt.Errorf("error writing nExpiryHeight: %w", err)
			}
		}
	}

	if tx.Version == VersionSapling {
		// valueBalanceSapling
		if err = putUint64(w, 0); err != nil {
			return nil, fmt.Errorf("error writing valueBalanceSapling: %w", err)
		}
	}

	if tx.Version >= VersionSapling {
		// nSpendsSapling
		if err = wire.WriteVarInt(w, pver, 0); err != nil {
			return nil, fmt.Errorf("error writing nSpendsSapling: %w", err)
		}

		// nOutputsSapling
		if err = wire.WriteVarInt(w, pver, 0); err != nil {
			return nil, fmt.Errorf("error writing nOutputsSapling: %w", err)
		}
	}

	if tx.Version >= VersionPreOverwinter && tx.Version <= VersionSapling {
		// nJoinSplit
		if err = wire.WriteVarInt(w, pver, 0); err != nil {
			return nil, fmt.Errorf("error writing nJoinSplit: %w", err)
		}
		return w.Bytes(), nil
	}

	// NU 5

	// no anchorSapling, because nSpendsSapling = 0
	// no bindingSigSapling or valueBalanceSapling, because nSpendsSapling + nOutputsSapling = 0

	if tx.Version == VersionNU5 {
		// nActionsOrchard
		if err = wire.WriteVarInt(w, pver, 0); err != nil {
			return nil, fmt.Errorf("error writing nActionsOrchard: %w", err)
		}
	}

	// vActionsOrchard, flagsOrchard, valueBalanceOrchard, anchorOrchard,
	// sizeProofsOrchard, proofsOrchard, vSpendAuthSigsOrchard, and
	// bindingSigOrchard are all empty, because nActionsOrchard = 0.

	return w.Bytes(), nil
}

// see https://zips.z.cash/protocol/protocol.pdf section 7.1
func DeserializeTx(b []byte) (*Tx, error) {
	tx := &Tx{MsgTx: new(wire.MsgTx)}
	r := bytes.NewReader(b)
	if err := tx.ZecDecode(r); err != nil {
		return nil, err
	}

	remains := r.Len()
	if remains > 0 {
		return nil, fmt.Errorf("incomplete deserialization. %d bytes remaining", remains)
	}

	return tx, nil
}

// ZecDecode reads the serialized transaction from the reader and populates the
// *Tx's fields.
func (tx *Tx) ZecDecode(r io.Reader) (err error) {
	ver, err := readUint32(r)
	if err != nil {
		return fmt.Errorf("error reading version: %w", err)
	}

	// overWintered := (ver & (1 << 31)) > 0
	ver &= overwinterMask // Clear the overwinter bit
	tx.Version = int32(ver)

	if tx.Version > VersionNU5 {
		return fmt.Errorf("unsupported tx version %d > 4", ver)
	}

	if ver >= VersionOverwinter {
		// nVersionGroupId uint32
		if err = discardBytes(r, 4); err != nil {
			return fmt.Errorf("error reading nVersionGroupId: %w", err)
		}
	}

	if ver == VersionNU5 {
		// nConsensusBranchId uint32
		if tx.ConsensusBranchId, err = readUint32(r); err != nil {
			return fmt.Errorf("error reading nConsensusBranchId: %w", err)
		}
		// lock_time
		if tx.LockTime, err = readUint32(r); err != nil {
			return fmt.Errorf("error reading lock_time: %w", err)
		}
		// nExpiryHeight
		if tx.ExpiryHeight, err = readUint32(r); err != nil {
			return fmt.Errorf("error reading nExpiryHeight: %w", err)
		}
	}

	txInCount, err := wire.ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	tx.TxIn = make([]*wire.TxIn, 0, txInCount)
	for i := 0; i < int(txInCount); i++ {
		ti := new(wire.TxIn)
		if err = readTxIn(r, ti); err != nil {
			return err
		}
		tx.TxIn = append(tx.TxIn, ti)
	}

	txOutCount, err := wire.ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	tx.TxOut = make([]*wire.TxOut, 0, txOutCount)
	for i := 0; i < int(txOutCount); i++ {
		to := new(wire.TxOut)
		if err = readTxOut(r, to); err != nil {
			return err
		}
		tx.TxOut = append(tx.TxOut, to)
	}

	if ver < VersionNU5 {
		// lock_time
		if tx.LockTime, err = readUint32(r); err != nil {
			return fmt.Errorf("error reading lock_time: %w", err)
		}
	}

	if ver == VersionOverwinter || ver == VersionSapling {
		// nExpiryHeight
		if tx.ExpiryHeight, err = readUint32(r); err != nil {
			return fmt.Errorf("error reading nExpiryHeight: %w", err)
		}
	}

	// That's it for pre-overwinter.
	if ver < VersionPreOverwinter {
		return nil
	}

	var bindingSigRequired bool
	if ver == VersionSapling {
		// valueBalanceSapling uint64
		if tx.ValueBalanceSapling, err = readInt64(r); err != nil {
			return fmt.Errorf("error reading valueBalanceSapling: %w", err)
		}

		if tx.NSpendsSapling, err = wire.ReadVarInt(r, pver); err != nil {
			return fmt.Errorf("error reading nSpendsSapling: %w", err)
		} else if tx.NSpendsSapling > 0 {
			// vSpendsSapling - discard
			bindingSigRequired = true
			if err = discardBytes(r, int64(tx.NSpendsSapling*384)); err != nil {
				return fmt.Errorf("error reading vSpendsSapling: %w", err)
			}
		}

		if tx.NOutputsSapling, err = wire.ReadVarInt(r, pver); err != nil {
			return fmt.Errorf("error reading nOutputsSapling: %w", err)
		} else if tx.NOutputsSapling > 0 {
			// vOutputsSapling - discard
			bindingSigRequired = true
			if err = discardBytes(r, int64(tx.NOutputsSapling*948)); err != nil {
				return fmt.Errorf("error reading vOutputsSapling: %w", err)
			}
		}
	}

	if ver <= VersionSapling && ver >= VersionPreOverwinter {
		if tx.NJoinSplit, err = wire.ReadVarInt(r, pver); err != nil {
			return fmt.Errorf("error reading nJoinSplit: %w", err)
		} else if tx.NJoinSplit > 0 {
			// vJoinSplit - discard
			tx.VJoinSplit = make([]*JoinSplit, 0, tx.NJoinSplit)
			for i := uint64(0); i < tx.NJoinSplit; i++ {
				sz := overwinterJoinSplitSize
				if ver == 4 {
					sz = saplingJoinSplitSize
				}
				old, err := readUint64(r)
				if err != nil {
					return fmt.Errorf("error reading joinsplit old: %w", err)
				}
				new, err := readUint64(r)
				if err != nil {
					return fmt.Errorf("error reading joinsplit new: %w", err)
				}
				tx.VJoinSplit = append(tx.VJoinSplit, &JoinSplit{
					Old: old,
					New: new,
				})
				if err = discardBytes(r, int64(sz-16)); err != nil {
					return fmt.Errorf("error reading vJoinSplit: %w", err)
				}
			}
			// joinSplitPubKey
			if err = discardBytes(r, 32); err != nil {
				return fmt.Errorf("error reading joinSplitPubKey: %w", err)
			}

			// joinSplitSig
			if err = discardBytes(r, 64); err != nil {
				return fmt.Errorf("error reading joinSplitSig: %w", err)
			}
		}
	} else { // NU5
		// nSpendsSapling
		tx.NSpendsSapling, err = wire.ReadVarInt(r, pver)
		if err != nil {
			return fmt.Errorf("error reading nSpendsSapling: %w", err)
		} else if tx.NSpendsSapling > 0 {
			// vSpendsSapling - discard
			bindingSigRequired = true
			if err = discardBytes(r, int64(tx.NSpendsSapling*96)); err != nil {
				return fmt.Errorf("error reading vSpendsSapling: %w", err)
			}
		}

		// nOutputsSapling
		tx.NOutputsSapling, err = wire.ReadVarInt(r, pver)
		if err != nil {
			return fmt.Errorf("error reading nSpendsSapling: %w", err)
		} else if tx.NOutputsSapling > 0 {
			// vOutputsSapling - discard
			bindingSigRequired = true
			if err = discardBytes(r, int64(tx.NOutputsSapling*756)); err != nil {
				return fmt.Errorf("error reading vOutputsSapling: %w", err)
			}
		}

		if tx.NOutputsSapling+tx.NSpendsSapling > 0 {
			// valueBalanceSpending uint64
			if err = discardBytes(r, 8); err != nil {
				return fmt.Errorf("error reading valueBalanceSpending: %w", err)
			}
		}

		if tx.NSpendsSapling > 0 {
			// anchorSapling
			if err = discardBytes(r, 32); err != nil {
				return fmt.Errorf("error reading anchorSapling: %w", err)
			}
			// vSpendProofsSapling
			if err = discardBytes(r, int64(tx.NSpendsSapling*192)); err != nil {
				return fmt.Errorf("error reading vSpendProofsSapling: %w", err)
			}
			// vSpendAuthSigsSapling
			if err = discardBytes(r, int64(tx.NSpendsSapling*64)); err != nil {
				return fmt.Errorf("error reading vSpendAuthSigsSapling: %w", err)
			}
		}

		if tx.NOutputsSapling > 0 {
			// vOutputProofsSapling
			if err = discardBytes(r, int64(tx.NOutputsSapling*192)); err != nil {
				return fmt.Errorf("error reading vOutputProofsSapling: %w", err)
			}
		}
	}

	if bindingSigRequired {
		// bindingSigSapling
		if err = discardBytes(r, 64); err != nil {
			return fmt.Errorf("error reading bindingSigSapling: %w", err)
		}
	}

	// pre-NU5 is done now.
	if ver < VersionNU5 {
		return nil
	}

	// NU5-only fields below.

	// nActionsOrchard
	tx.NActionsOrchard, err = wire.ReadVarInt(r, pver)
	if err != nil {
		return fmt.Errorf("error reading bindingSigSapling: %w", err)
	}

	if tx.NActionsOrchard == 0 {
		return nil
	}

	// vActionsOrchard
	if err = discardBytes(r, int64(tx.NActionsOrchard*820)); err != nil {
		return fmt.Errorf("error reading vActionsOrchard: %w", err)
	}

	// flagsOrchard
	if err = discardBytes(r, 1); err != nil {
		return fmt.Errorf("error reading flagsOrchard: %w", err)
	}

	// valueBalanceOrchard uint64
	if tx.ValueBalanceOrchard, err = readInt64(r); err != nil {
		return fmt.Errorf("error reading valueBalanceOrchard: %w", err)
	}

	// anchorOrchard
	if err = discardBytes(r, 32); err != nil {
		return fmt.Errorf("error reading anchorOrchard: %w", err)
	}

	// sizeProofsOrchard
	tx.SizeProofsOrchard, err = wire.ReadVarInt(r, pver)
	if err != nil {
		return fmt.Errorf("error reading sizeProofsOrchard: %w", err)
	}

	// proofsOrchard
	if err = discardBytes(r, int64(tx.SizeProofsOrchard)); err != nil {
		return fmt.Errorf("error reading proofsOrchard: %w", err)
	}

	// vSpendAuthSigsOrchard
	if err = discardBytes(r, int64(tx.NActionsOrchard*64)); err != nil {
		return fmt.Errorf("error reading vSpendAuthSigsOrchard: %w", err)
	}

	// bindingSigOrchard
	if err = discardBytes(r, 64); err != nil {
		return fmt.Errorf("error reading bindingSigOrchard: %w", err)
	}

	return nil
}

// SerializeSize is the size of the transaction when serialized.
func (tx *Tx) SerializeSize() uint64 {
	var sz uint64 = 4 // header
	ver := tx.Version
	sz += uint64(wire.VarIntSerializeSize(uint64(len(tx.TxIn)))) // tx_in_count
	for _, txIn := range tx.TxIn {                               // tx_in
		sz += 32 /* prev hash */ + 4 /* prev index */ + 4 /* sequence */
		sz += uint64(wire.VarIntSerializeSize(uint64(len(txIn.SignatureScript))) + len(txIn.SignatureScript))
	}
	sz += uint64(wire.VarIntSerializeSize(uint64(len(tx.TxOut)))) // tx_out_count
	for _, txOut := range tx.TxOut {                              // tx_out
		sz += 8 /* value */
		sz += uint64(wire.VarIntSerializeSize(uint64(len(txOut.PkScript))) + len(txOut.PkScript))
	}
	sz += 4 // lockTime

	// join-splits are only versions 2 to 4.
	if ver >= VersionPreOverwinter && ver < VersionNU5 {
		sz += uint64(wire.VarIntSerializeSize(tx.NJoinSplit))
		if tx.NJoinSplit > 0 {
			if ver < VersionSapling {
				sz += tx.NJoinSplit * overwinterJoinSplitSize
			} else {
				sz += tx.NJoinSplit * saplingJoinSplitSize
			}
			sz += 32 // joinSplitPubKey
			sz += 64 // joinSplitSig
		}
	}

	if ver >= VersionOverwinter {
		sz += 4 // nExpiryHeight
		sz += 4 // nVersionGroupId
	}

	if ver >= VersionSapling {
		sz += 8                                                    // valueBalanceSapling
		sz += uint64(wire.VarIntSerializeSize(tx.NSpendsSapling))  // nSpendsSapling
		sz += 384 * tx.NSpendsSapling                              // vSpendsSapling
		sz += uint64(wire.VarIntSerializeSize(tx.NOutputsSapling)) // nOutputsSapling
		sz += 948 * tx.NOutputsSapling                             // vOutputsSapling
		if tx.NSpendsSapling+tx.NOutputsSapling > 0 {
			sz += 64 // bindingSigSapling
		}
	}

	if ver == VersionNU5 {
		// With nSpendsSapling = 0 and nOutputsSapling = 0
		sz += 4                                                    // nConsensusBranchId
		sz += uint64(wire.VarIntSerializeSize(tx.NActionsOrchard)) // nActionsOrchard
		if tx.NActionsOrchard > 0 {
			sz += tx.NActionsOrchard * 820                               // vActionsOrchard
			sz++                                                         // flagsOrchard
			sz += 8                                                      // valueBalanceOrchard
			sz += 32                                                     // anchorOrchard
			sz += uint64(wire.VarIntSerializeSize(tx.SizeProofsOrchard)) // sizeProofsOrchard
			sz += tx.SizeProofsOrchard                                   // proofsOrchard
			sz += 64 * tx.NActionsOrchard                                // vSpendAuthSigsOrchard
			sz += 64                                                     // bindingSigOrchard

		}
	}
	return sz
}

// // RequiredTxFeesZIP317 calculates the minimum tx fees according to ZIP-0317.
// func (tx *Tx) RequiredTxFeesZIP317() uint64 {
// 	txInsSize := uint64(wire.VarIntSerializeSize(uint64(len(tx.TxIn))))
// 	for _, txIn := range tx.TxIn {
// 		txInsSize += uint64(txIn.SerializeSize())
// 	}
//
// 	txOutsSize := uint64(wire.VarIntSerializeSize(uint64(len(tx.TxOut))))
// 	for _, txOut := range tx.TxOut {
// 		txOutsSize += uint64(txOut.SerializeSize())
// 	}
//
// 	return TxFeesZIP317(txInsSize, txOutsSize, tx.NSpendsSapling, tx.NOutputsSapling, tx.NJoinSplit, tx.NActionsOrchard)
// }

// writeTxIn encodes ti to the bitcoin protocol encoding for a transaction
// input (TxIn) to w.
func writeTxIn(w io.Writer, ti *wire.TxIn) error {
	err := writeOutPoint(w, &ti.PreviousOutPoint)
	if err != nil {
		return err
	}

	err = wire.WriteVarBytes(w, pver, ti.SignatureScript)
	if err != nil {
		return err
	}

	return putUint32(w, ti.Sequence)
}

// writeOutPoint encodes op to the bitcoin protocol encoding for an OutPoint
// to w.
func writeOutPoint(w io.Writer, op *wire.OutPoint) error {
	_, err := w.Write(op.Hash[:])
	if err != nil {
		return err
	}
	return putUint32(w, op.Index)
}

func uint32Bytes(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

// putUint32 writes a little-endian encoded uint32 to the Writer.
func putUint32(w io.Writer, v uint32) error {
	_, err := w.Write(uint32Bytes(v))
	return err
}

func uint64Bytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	return b
}

func int64Bytes(v int64) []byte {
	return uint64Bytes(uint64(v))
}

// putUint64 writes a little-endian encoded uint64 to the Writer.
func putUint64(w io.Writer, v uint64) error {
	_, err := w.Write(uint64Bytes(v))
	return err
}

// readUint32 reads a little-endian encoded uint32 from the Reader.
func readUint32(r io.Reader) (uint32, error) {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}

// readUint64 reads a little-endian encoded uint64 from the Reader.
func readUint64(r io.Reader) (uint64, error) {
	b := make([]byte, 8)
	if _, err := io.ReadFull(r, b); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b), nil
}

func readInt64(r io.Reader) (int64, error) {
	u, err := readUint64(r)
	if err != nil {
		return 0, err
	}
	return int64(u), nil
}

// readTxIn reads the next sequence of bytes from r as a transaction input.
func readTxIn(r io.Reader, ti *wire.TxIn) error {
	err := readOutPoint(r, &ti.PreviousOutPoint)
	if err != nil {
		return err
	}

	ti.SignatureScript, err = readScript(r)
	if err != nil {
		return err
	}

	ti.Sequence, err = readUint32(r)
	return err
}

// readTxOut reads the next sequence of bytes from r as a transaction output.
func readTxOut(r io.Reader, to *wire.TxOut) error {
	v, err := readUint64(r)
	if err != nil {
		return err
	}
	to.Value = int64(v)

	to.PkScript, err = readScript(r)
	return err
}

// readOutPoint reads the next sequence of bytes from r as an OutPoint.
func readOutPoint(r io.Reader, op *wire.OutPoint) error {
	_, err := io.ReadFull(r, op.Hash[:])
	if err != nil {
		return err
	}

	op.Index, err = readUint32(r)
	return err
}

// readScript reads a variable length byte array. Copy of unexported
// btcd/wire.readScript.
func readScript(r io.Reader) ([]byte, error) {
	count, err := wire.ReadVarInt(r, pver)
	if err != nil {
		return nil, err
	}
	if count > uint64(wire.MaxMessagePayload) {
		return nil, fmt.Errorf("larger than the max allowed size "+
			"[count %d, max %d]", count, wire.MaxMessagePayload)
	}
	b := make([]byte, count)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func discardBytes(r io.Reader, n int64) error {
	m, err := io.CopyN(io.Discard, r, n)
	if err != nil {
		return err
	}
	if m != n {
		return fmt.Errorf("only discarded %d of %d bytes", m, n)
	}
	return nil
}

// blake2bHash is a BLAKE-2B hash of the data with the specified personalization
// key.
func blake2bHash(data, personalizationKey []byte) (_ [32]byte, err error) {
	bHash, err := blake2b.New(&blake2b.Config{Size: 32, Person: personalizationKey})
	if err != nil {
		return
	}

	if _, err = bHash.Write(data); err != nil {
		return
	}

	var h [32]byte
	copy(h[:], bHash.Sum(nil))
	return h, err
}

// CalcTxSize calculates the size of a Zcash transparent transaction. CalcTxSize
// won't return accurate results for shielded or blended transactions.
func CalcTxSize(tx *wire.MsgTx) uint64 {
	return (&Tx{MsgTx: tx}).SerializeSize()
}
