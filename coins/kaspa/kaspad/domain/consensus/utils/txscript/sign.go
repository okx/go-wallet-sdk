// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/consensushashing"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/dagconfig"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/util"
	"github.com/pkg/errors"
)

// RawTxInSignature returns the serialized Schnorr signature for the input idx of
// the given transaction, with hashType appended to it.
func RawTxInSignature(tx *externalapi.DomainTransaction, idx int, hashType consensushashing.SigHashType,
	key *btcec.PrivateKey, sighashReusedValues *consensushashing.SighashReusedValues) ([]byte, error) {

	hash, err := consensushashing.CalculateSignatureHashSchnorr(tx, idx, hashType, sighashReusedValues)
	if err != nil {
		return nil, err
	}
	signature, err := schnorr.Sign(key, (*hash.ByteArray())[:])
	if err != nil {
		return nil, errors.Errorf("cannot sign tx input: %s", err)
	}

	return append(signature.Serialize(), byte(hashType)), nil
}

// RawTxInSignatureECDSA returns the serialized ECDSA signature for the input idx of
// the given transaction, with hashType appended to it.
func RawTxInSignatureECDSA(tx *externalapi.DomainTransaction, idx int, hashType consensushashing.SigHashType,
	key *btcec.PrivateKey, sighashReusedValues *consensushashing.SighashReusedValues) ([]byte, error) {

	hash, err := consensushashing.CalculateSignatureHashECDSA(tx, idx, hashType, sighashReusedValues)
	if err != nil {
		return nil, err
	}
	signature, err := ecdsa.SignCompact(key, (*hash.ByteArray())[:], true)
	if err != nil {
		return nil, errors.Errorf("cannot sign tx input: %s", err)
	}

	return append(signature[1:], byte(hashType)), nil
}

// SignatureScript creates an input signature script for tx to spend KAS sent
// from a previous output to the owner of a Schnorr private key. tx must include all
// transaction inputs and outputs, however txin scripts are allowed to be filled
// or empty. The returned script is calculated to be used as the idx'th txin
// sigscript for tx. script is the ScriptPublicKey of the previous output being used
// as the idx'th input. privKey is serialized in either a compressed or
// uncompressed format based on compress. This format must match the same format
// used to generate the payment address, or the script validation will fail.
func SignatureScript(tx *externalapi.DomainTransaction, idx int, hashType consensushashing.SigHashType,
	privKey *btcec.PrivateKey, sighashReusedValues *consensushashing.SighashReusedValues) ([]byte, error) {

	sig, err := RawTxInSignature(tx, idx, hashType, privKey, sighashReusedValues)
	if err != nil {
		return nil, err
	}

	return NewScriptBuilder().AddData(sig).Script()
}

// SignatureScriptECDSA creates an input signature script for tx to spend KAS sent
// from a previous output to the owner of an ECDSA private key. tx must include all
// transaction inputs and outputs, however txin scripts are allowed to be filled
// or empty. The returned script is calculated to be used as the idx'th txin
// sigscript for tx. script is the ScriptPublicKey of the previous output being used
// as the idx'th input. privKey is serialized in either a compressed or
// uncompressed format based on compress. This format must match the same format
// used to generate the payment address, or the script validation will fail.
func SignatureScriptECDSA(tx *externalapi.DomainTransaction, idx int, hashType consensushashing.SigHashType,
	privKey *btcec.PrivateKey, sighashReusedValues *consensushashing.SighashReusedValues) ([]byte, error) {

	sig, err := RawTxInSignatureECDSA(tx, idx, hashType, privKey, sighashReusedValues)
	if err != nil {
		return nil, err
	}

	return NewScriptBuilder().AddData(sig).Script()
}

func sign(dagParams *dagconfig.Params, tx *externalapi.DomainTransaction, idx int,
	script *externalapi.ScriptPublicKey, hashType consensushashing.SigHashType,
	sighashReusedValues *consensushashing.SighashReusedValues, kdb KeyDB, sdb ScriptDB) (
	[]byte, ScriptClass, util.Address, error) {

	class, address, err := ExtractScriptPubKeyAddress(script, dagParams)
	if err != nil {
		return nil, NonStandardTy, nil, err
	}

	switch class {
	case PubKeyTy:
		// look up key for address
		key, err := kdb.GetKey(address)
		if err != nil {
			return nil, class, nil, err
		}

		signedScript, err := SignatureScript(tx, idx, hashType, key, sighashReusedValues)
		if err != nil {
			return nil, class, nil, err
		}

		return signedScript, class, address, nil
	case ScriptHashTy:
		script, err := sdb.GetScript(address)
		if err != nil {
			return nil, class, nil, err
		}

		return script, class, address, nil
	default:
		return nil, class, nil, errors.New("can't sign unknown transactions")
	}
}

// mergeScripts merges sigScript and prevScript assuming they are both
// partial solutions for scriptPubKey spending output idx of tx. class, addresses
// and nrequired are the result of extracting the addresses from scriptPubKey.
// The return value is the best effort merging of the two scripts. Calling this
// function with addresses, class and nrequired that do not match scriptPubKey is
// an error and results in undefined behaviour.
func mergeScripts(dagParams *dagconfig.Params, tx *externalapi.DomainTransaction, idx int,
	class ScriptClass, sigScript []byte, prevScript *externalapi.ScriptPublicKey) ([]byte, error) {

	switch class {
	case ScriptHashTy:
		// Remove the last push in the script and then recurse.
		// this could be a lot less inefficient.
		sigPops, err := parseScript(sigScript)
		if err != nil || len(sigPops) == 0 {
			return prevScript.Script, nil
		}
		prevPops, err := parseScript(prevScript.Script)
		if err != nil || len(prevPops) == 0 {
			return sigScript, nil
		}

		// assume that script in sigPops is the correct one, we just
		// made it.
		script := sigPops[len(sigPops)-1].data
		scriptPubKey := &externalapi.ScriptPublicKey{
			Script:  script,
			Version: prevScript.Version,
		}
		// We already know this information somewhere up the stack.
		class, _, _ :=
			ExtractScriptPubKeyAddress(scriptPubKey, dagParams)

		// regenerate scripts.
		sigScript, _ := unparseScript(sigPops)
		prevScriptByte, _ := unparseScript(prevPops)
		prevScript = &externalapi.ScriptPublicKey{
			Script:  prevScriptByte,
			Version: prevScript.Version,
		}
		// Merge
		mergedScript, err := mergeScripts(dagParams, tx, idx, class, sigScript, prevScript)
		if err != nil {
			return nil, err
		}

		// Reappend the script and return the result.
		builder := NewScriptBuilder()
		builder.AddOps(mergedScript)
		builder.AddData(script)
		return builder.Script()

	// It doesn't actually make sense to merge anything other than multiig
	// and scripthash (because it could contain multisig). Everything else
	// has either zero signature, can't be spent, or has a single signature
	// which is either present or not. The other two cases are handled
	// above. In the conflict case here we just assume the longest is
	// correct (this matches behaviour of the reference implementation).
	default:
		if len(sigScript) > len(prevScript.Script) {
			return sigScript, nil
		}
		return prevScript.Script, nil
	}
}

// KeyDB is an interface type provided to SignTxOutput, it encapsulates
// any user state required to get the private keys for an address.
type KeyDB interface {
	GetKey(util.Address) (*btcec.PrivateKey, error)
}

// KeyClosure implements KeyDB with a closure.
type KeyClosure func(util.Address) (*btcec.PrivateKey, error)

// GetKey implements KeyDB by returning the result of calling the closure.
func (kc KeyClosure) GetKey(address util.Address) (*btcec.PrivateKey, error) {
	return kc(address)
}

// ScriptDB is an interface type provided to SignTxOutput, it encapsulates any
// user state required to get the scripts for an pay-to-script-hash address.
type ScriptDB interface {
	GetScript(util.Address) ([]byte, error)
}

// ScriptClosure implements ScriptDB with a closure.
type ScriptClosure func(util.Address) ([]byte, error)

// GetScript implements ScriptDB by returning the result of calling the closure.
func (sc ScriptClosure) GetScript(address util.Address) ([]byte, error) {
	return sc(address)
}

// SignTxOutput signs output idx of the given tx to resolve the script given in
// scriptPublicKey with a signature type of hashType. Any keys required will be
// looked up by calling getKey() with the string of the given address.
// Any pay-to-script-hash signatures will be similarly looked up by calling
// getScript. If previousScript is provided then the results in previousScript
// will be merged in a type-dependent manner with the newly generated.
// signature script.
func SignTxOutput(dagParams *dagconfig.Params, tx *externalapi.DomainTransaction, idx int,
	scriptPublicKey *externalapi.ScriptPublicKey, hashType consensushashing.SigHashType,
	sighashReusedValues *consensushashing.SighashReusedValues, kdb KeyDB, sdb ScriptDB,
	previousScript *externalapi.ScriptPublicKey) ([]byte, error) {

	sigScript, class, _, err := sign(dagParams, tx,
		idx, scriptPublicKey, hashType, sighashReusedValues, kdb, sdb)
	if err != nil {
		return nil, err
	}

	if class == ScriptHashTy {
		scriptHashPreimageScriptPublicKey := &externalapi.ScriptPublicKey{
			Script:  sigScript,
			Version: scriptPublicKey.Version,
		}

		realSigScript, _, _, err := sign(dagParams, tx, idx,
			scriptHashPreimageScriptPublicKey, hashType, sighashReusedValues, kdb, sdb)
		if err != nil {
			return nil, err
		}

		// Append the p2sh script as the last push in the script.
		builder := NewScriptBuilder()
		builder.AddOps(realSigScript)
		builder.AddData(sigScript)

		sigScript, _ = builder.Script()
	}

	// Merge scripts. with any previous data, if any.
	return mergeScripts(dagParams, tx, idx, class, sigScript, previousScript)
}
