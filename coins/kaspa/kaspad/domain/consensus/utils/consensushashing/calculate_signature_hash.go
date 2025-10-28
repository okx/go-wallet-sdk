package consensushashing

import (
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/hashes"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/serialization"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/subnetworks"
	"github.com/pkg/errors"
)

// SigHashType represents hash type bits at the end of a signature.
type SigHashType uint8

// Hash type bits from the end of a signature.
const (
	SigHashAll          SigHashType = 0b00000001
	SigHashNone         SigHashType = 0b00000010
	SigHashSingle       SigHashType = 0b00000100
	SigHashAnyOneCanPay SigHashType = 0b10000000

	// SigHashMask defines the number of bits of the hash type which is used
	// to identify which outputs are signed.
	SigHashMask = 0b00000111
)

// IsStandardSigHashType returns true if sht represents a standard SigHashType
func (sht SigHashType) IsStandardSigHashType() bool {
	switch sht {
	case SigHashAll, SigHashNone, SigHashSingle,
		SigHashAll | SigHashAnyOneCanPay, SigHashNone | SigHashAnyOneCanPay, SigHashSingle | SigHashAnyOneCanPay:
		return true
	default:
		return false
	}
}

func (sht SigHashType) isSigHashAll() bool {
	return sht&SigHashMask == SigHashAll
}
func (sht SigHashType) isSigHashNone() bool {
	return sht&SigHashMask == SigHashNone
}
func (sht SigHashType) isSigHashSingle() bool {
	return sht&SigHashMask == SigHashSingle
}
func (sht SigHashType) isSigHashAnyOneCanPay() bool {
	return sht&SigHashAnyOneCanPay == SigHashAnyOneCanPay
}

// SighashReusedValues holds all fields used in the calculation of a transaction's sigHash, that are
// the same for all transaction inputs.
// Reuse of such values prevents the quadratic hashing problem.
type SighashReusedValues struct {
	previousOutputsHash *externalapi.DomainHash
	sequencesHash       *externalapi.DomainHash
	sigOpCountsHash     *externalapi.DomainHash
	outputsHash         *externalapi.DomainHash
	payloadHash         *externalapi.DomainHash
}

// CalculateSignatureHashSchnorr will, given a script and hash type calculate the signature hash
// to be used for signing and verification for Schnorr.
// This returns error only if one of the provided parameters are consensus-invalid.
func CalculateSignatureHashSchnorr(tx *externalapi.DomainTransaction, inputIndex int, hashType SigHashType,
	reusedValues *SighashReusedValues) (*externalapi.DomainHash, error) {

	if !hashType.IsStandardSigHashType() {
		return nil, errors.Errorf("SigHashType %d is not a valid SigHash type", hashType)
	}

	txIn := tx.Inputs[inputIndex]
	prevScriptPublicKey := txIn.UTXOEntry.ScriptPublicKey()
	return calculateSignatureHash(tx, inputIndex, txIn, prevScriptPublicKey, hashType, reusedValues)
}

// CalculateSignatureHashECDSA will, given a script and hash type calculate the signature hash
// to be used for signing and verification for ECDSA.
// This returns error only if one of the provided parameters are consensus-invalid.
func CalculateSignatureHashECDSA(tx *externalapi.DomainTransaction, inputIndex int, hashType SigHashType,
	reusedValues *SighashReusedValues) (*externalapi.DomainHash, error) {

	hash, err := CalculateSignatureHashSchnorr(tx, inputIndex, hashType, reusedValues)
	if err != nil {
		return nil, err
	}

	hashWriter := hashes.NewTransactionSigningHashECDSAWriter()
	hashWriter.InfallibleWrite(hash.ByteSlice())

	return hashWriter.Finalize(), nil
}

func calculateSignatureHash(tx *externalapi.DomainTransaction, inputIndex int, txIn *externalapi.DomainTransactionInput,
	prevScriptPublicKey *externalapi.ScriptPublicKey, hashType SigHashType, reusedValues *SighashReusedValues) (
	*externalapi.DomainHash, error) {

	hashWriter := hashes.NewTransactionSigningHashWriter()
	infallibleWriteElement(hashWriter, tx.Version)

	previousOutputsHash := getPreviousOutputsHash(tx, hashType, reusedValues)
	infallibleWriteElement(hashWriter, previousOutputsHash)

	sequencesHash := getSequencesHash(tx, hashType, reusedValues)
	infallibleWriteElement(hashWriter, sequencesHash)

	sigOpCountsHash := getSigOpCountsHash(tx, hashType, reusedValues)
	infallibleWriteElement(hashWriter, sigOpCountsHash)

	hashOutpoint(hashWriter, txIn.PreviousOutpoint)

	infallibleWriteElement(hashWriter, prevScriptPublicKey.Version)
	infallibleWriteElement(hashWriter, prevScriptPublicKey.Script)

	infallibleWriteElement(hashWriter, txIn.UTXOEntry.Amount())

	infallibleWriteElement(hashWriter, txIn.Sequence)

	infallibleWriteElement(hashWriter, txIn.SigOpCount)

	outputsHash := getOutputsHash(tx, inputIndex, hashType, reusedValues)
	infallibleWriteElement(hashWriter, outputsHash)

	infallibleWriteElement(hashWriter, tx.LockTime)

	infallibleWriteElement(hashWriter, tx.SubnetworkID)
	infallibleWriteElement(hashWriter, tx.Gas)

	payloadHash := getPayloadHash(tx, reusedValues)
	infallibleWriteElement(hashWriter, payloadHash)

	infallibleWriteElement(hashWriter, uint8(hashType))

	return hashWriter.Finalize(), nil
}

func getPreviousOutputsHash(tx *externalapi.DomainTransaction, hashType SigHashType, reusedValues *SighashReusedValues) *externalapi.DomainHash {
	if hashType.isSigHashAnyOneCanPay() {
		return externalapi.NewZeroHash()
	}

	if reusedValues.previousOutputsHash == nil {
		hashWriter := hashes.NewTransactionSigningHashWriter()
		for _, txIn := range tx.Inputs {
			hashOutpoint(hashWriter, txIn.PreviousOutpoint)
		}
		reusedValues.previousOutputsHash = hashWriter.Finalize()
	}

	return reusedValues.previousOutputsHash
}

func getSequencesHash(tx *externalapi.DomainTransaction, hashType SigHashType, reusedValues *SighashReusedValues) *externalapi.DomainHash {
	if hashType.isSigHashSingle() || hashType.isSigHashAnyOneCanPay() || hashType.isSigHashNone() {
		return externalapi.NewZeroHash()
	}

	if reusedValues.sequencesHash == nil {
		hashWriter := hashes.NewTransactionSigningHashWriter()
		for _, txIn := range tx.Inputs {
			infallibleWriteElement(hashWriter, txIn.Sequence)
		}
		reusedValues.sequencesHash = hashWriter.Finalize()
	}

	return reusedValues.sequencesHash
}

func getSigOpCountsHash(tx *externalapi.DomainTransaction, hashType SigHashType, reusedValues *SighashReusedValues) *externalapi.DomainHash {
	if hashType.isSigHashAnyOneCanPay() {
		return externalapi.NewZeroHash()
	}

	if reusedValues.sigOpCountsHash == nil {
		hashWriter := hashes.NewTransactionSigningHashWriter()
		for _, txIn := range tx.Inputs {
			infallibleWriteElement(hashWriter, txIn.SigOpCount)
		}
		reusedValues.sigOpCountsHash = hashWriter.Finalize()
	}

	return reusedValues.sigOpCountsHash
}

func getOutputsHash(tx *externalapi.DomainTransaction, inputIndex int, hashType SigHashType, reusedValues *SighashReusedValues) *externalapi.DomainHash {
	// SigHashNone: return zero-hash
	if hashType.isSigHashNone() {
		return externalapi.NewZeroHash()
	}

	// SigHashSingle: If the relevant output exists - return it's hash, otherwise return zero-hash
	if hashType.isSigHashSingle() {
		if inputIndex >= len(tx.Outputs) {
			return externalapi.NewZeroHash()
		}
		hashWriter := hashes.NewTransactionSigningHashWriter()
		hashTxOut(hashWriter, tx.Outputs[inputIndex])
		return hashWriter.Finalize()
	}

	// SigHashAll: Return hash of all outputs. Re-use hash if available.
	if reusedValues.outputsHash == nil {
		hashWriter := hashes.NewTransactionSigningHashWriter()
		for _, txOut := range tx.Outputs {
			hashTxOut(hashWriter, txOut)
		}
		reusedValues.outputsHash = hashWriter.Finalize()
	}

	return reusedValues.outputsHash
}

func getPayloadHash(tx *externalapi.DomainTransaction, reusedValues *SighashReusedValues) *externalapi.DomainHash {
	if tx.SubnetworkID.Equal(&subnetworks.SubnetworkIDNative) {
		return externalapi.NewZeroHash()
	}

	if reusedValues.payloadHash == nil {
		hashWriter := hashes.NewTransactionSigningHashWriter()
		infallibleWriteElement(hashWriter, tx.Payload)
		reusedValues.payloadHash = hashWriter.Finalize()
	}
	return reusedValues.payloadHash
}

func hashTxOut(hashWriter hashes.HashWriter, txOut *externalapi.DomainTransactionOutput) {
	infallibleWriteElement(hashWriter, txOut.Value)
	infallibleWriteElement(hashWriter, txOut.ScriptPublicKey.Version)
	infallibleWriteElement(hashWriter, txOut.ScriptPublicKey.Script)
}

func hashOutpoint(hashWriter hashes.HashWriter, outpoint externalapi.DomainOutpoint) {
	infallibleWriteElement(hashWriter, outpoint.TransactionID)
	infallibleWriteElement(hashWriter, outpoint.Index)
}

func infallibleWriteElement(hashWriter hashes.HashWriter, element interface{}) {
	err := serialization.WriteElement(hashWriter, element)
	if err != nil {
		// It seems like this could only happen if the writer returned an error.
		// and this writer should never return an error (no allocations or possible failures)
		// the only non-writer error path here is unknown types in `WriteElement`
		panic(errors.Wrap(err, "TransactionHashForSigning() failed. this should never fail for structurally-valid transactions"))
	}
}
