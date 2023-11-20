package bitcoin

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts"
)

type TxInput struct {
	TxId              string
	VOut              uint32
	Sequence          uint32
	Amount            int64
	Address           string
	PrivateKey        string
	NonWitnessUtxo    string
	MasterFingerprint uint32
	DerivationPath    string
	PublicKey         string
}

type TxOutput struct {
	Address           string
	Amount            int64
	IsChange          bool
	MasterFingerprint uint32
	DerivationPath    string
	PublicKey         string
}

const SellerSignatureIndex = 2

func GenerateSignedListingPSBTBase64(in *TxInput, out *TxOutput, network *chaincfg.Params) (string, error) {
	txHash, err := chainhash.NewHashFromStr(in.TxId)
	if err != nil {
		return "", err
	}
	prevOut := wire.NewOutPoint(txHash, in.VOut)
	inputs := []*wire.OutPoint{{Index: 0}, {Index: 1}, prevOut}

	pkScript, err := AddrToPkScript(out.Address, network)
	if err != nil {
		return "", err
	}
	// placeholder
	dummyPkScript, err := AddrToPkScript("bc1pcyj5mt2q4t4py8jnur8vpxvxxchke4pzy7tdr9yvj3u3kdfgrj6sw3rzmr", network)
	if err != nil {
		return "", err
	}
	outputs := []*wire.TxOut{{PkScript: dummyPkScript}, {PkScript: dummyPkScript}, wire.NewTxOut(out.Amount, pkScript)}

	nSequences := []uint32{wire.MaxTxInSequenceNum, wire.MaxTxInSequenceNum, wire.MaxTxInSequenceNum}
	p, err := psbt.New(inputs, outputs, int32(2), uint32(0), nSequences)
	if err != nil {
		return "", err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return "", err
	}

	dummyWitnessUtxo := wire.NewTxOut(0, dummyPkScript)
	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 0)
	if err != nil {
		return "", err
	}
	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 1)
	if err != nil {
		return "", err
	}

	prevPkScript, err := AddrToPkScript(in.Address, network)
	if err != nil {
		return "", err
	}
	witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
	prevOuts := map[wire.OutPoint]*wire.TxOut{
		wire.OutPoint{Index: 0}: dummyWitnessUtxo,
		wire.OutPoint{Index: 1}: dummyWitnessUtxo,
		*prevOut:                witnessUtxo,
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	err = signInput(updater, SellerSignatureIndex, in, prevOutputFetcher, txscript.SigHashSingle|txscript.SigHashAnyOneCanPay, network)
	if err != nil {
		return "", err
	}

	return p.B64Encode()
}

func GenerateSignedBuyingTx(ins []*TxInput, outs []*TxOutput, sellerPsbt string, network *chaincfg.Params) (string, error) {
	sp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPsbt)), true)
	if err != nil {
		return "", err
	}

	var inputs []*wire.OutPoint
	var nSequences []uint32
	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for i, in := range ins {
		var prevOut *wire.OutPoint
		if i == SellerSignatureIndex {
			prevOut = &sp.UnsignedTx.TxIn[i].PreviousOutPoint
		} else {
			txHash, err := chainhash.NewHashFromStr(in.TxId)
			if err != nil {
				return "", err
			}
			prevOut = wire.NewOutPoint(txHash, in.VOut)
		}
		inputs = append(inputs, prevOut)

		prevPkScript, err := AddrToPkScript(in.Address, network)
		if err != nil {
			return "", err
		}
		witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
		prevOuts[*prevOut] = witnessUtxo

		nSequences = append(nSequences, wire.MaxTxInSequenceNum)
	}

	var outputs []*wire.TxOut
	for i, out := range outs {
		if i == SellerSignatureIndex {
			outputs = append(outputs, sp.UnsignedTx.TxOut[i])
		} else {
			pkScript, err := AddrToPkScript(out.Address, network)
			if err != nil {
				return "", err
			}
			outputs = append(outputs, wire.NewTxOut(out.Amount, pkScript))
		}
	}

	bp, err := psbt.New(inputs, outputs, int32(2), uint32(0), nSequences)
	if err != nil {
		return "", err
	}

	updater, err := psbt.NewUpdater(bp)
	if err != nil {
		return "", err
	}

	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	for i, in := range ins {
		if i == SellerSignatureIndex {
			continue
		}

		if err = signInput(updater, i, in, prevOutputFetcher, txscript.SigHashAll, network); err != nil {
			return "", err
		}

		if err = psbt.Finalize(bp, i); err != nil {
			return "", err
		}
	}

	// bp.UnsignedTx.TxIn[SellerSignatureIndex] = sp.UnsignedTx.TxIn[SellerSignatureIndex]
	bp.Inputs[SellerSignatureIndex] = sp.Inputs[SellerSignatureIndex]

	if err = psbt.MaybeFinalizeAll(bp); err != nil {
		return "", err
	}

	buyerSignedTx, err := psbt.Extract(bp)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := buyerSignedTx.Serialize(&buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

func signInput(updater *psbt.Updater, i int, in *TxInput, prevOutFetcher *txscript.MultiPrevOutFetcher, hashType txscript.SigHashType, network *chaincfg.Params) error {
	wif, err := btcutil.DecodeWIF(in.PrivateKey)
	if err != nil {
		return err
	}
	privKey := wif.PrivKey

	prevPkScript, err := AddrToPkScript(in.Address, network)
	if err != nil {
		return err
	}
	if txscript.IsPayToPubKeyHash(prevPkScript) {
		prevTx := wire.NewMsgTx(2)
		txBytes, err := hex.DecodeString(in.NonWitnessUtxo)
		if err != nil {
			return err
		}
		if err = prevTx.Deserialize(bytes.NewReader(txBytes)); err != nil {
			return err
		}
		if err = updater.AddInNonWitnessUtxo(prevTx, i); err != nil {
			return err
		}
	} else {
		witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
		if err = updater.AddInWitnessUtxo(witnessUtxo, i); err != nil {
			return err
		}
	}

	if err = updater.AddInSighashType(hashType, i); err != nil {
		return err
	}

	if txscript.IsPayToTaproot(prevPkScript) {
		internalPubKey := schnorr.SerializePubKey(privKey.PubKey())
		updater.Upsbt.Inputs[i].TaprootInternalKey = internalPubKey

		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)
		if hashType == txscript.SigHashAll {
			hashType = txscript.SigHashDefault
		}
		witness, err := txscript.TaprootWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes,
			i, in.Amount, prevPkScript, hashType, privKey)
		if err != nil {
			return err
		}

		updater.Upsbt.Inputs[i].TaprootKeySpendSig = witness[0]
	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		signature, err := txscript.RawTxInSignature(updater.Upsbt.UnsignedTx, i, prevPkScript, hashType, privKey)
		if err != nil {
			return err
		}

		if _, err := updater.Sign(i, signature, privKey.PubKey().SerializeCompressed(), nil, nil); err != nil {
			return err
		}
	} else {
		pubKeyBytes := privKey.PubKey().SerializeCompressed()
		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)

		script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
		if err != nil {
			return err
		}
		signature, err := txscript.RawTxInWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes, i, in.Amount, script, hashType, privKey)
		if err != nil {
			return err
		}

		if txscript.IsPayToScriptHash(prevPkScript) {
			redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return err
			}
			err = updater.AddInRedeemScript(redeemScript, i)
			if err != nil {
				return err
			}
		}

		if _, err := updater.Sign(i, signature, pubKeyBytes, nil, nil); err != nil {
			return err
		}
	}
	return nil
}

func CalcFee(ins []*TxInput, outs []*TxOutput, sellerPsbt string, feeRate int64, network *chaincfg.Params) (int64, error) {
	txHex, err := GenerateSignedBuyingTx(ins, outs, sellerPsbt, network)

	tx := wire.NewMsgTx(2)
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return 0, err
	}
	err = tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return 0, err
	}

	return GetTxVirtualSize(btcutil.NewTx(tx)) * feeRate, nil
}

func AddrToPkScript(addr string, network *chaincfg.Params) ([]byte, error) {
	address, err := btcutil.DecodeAddress(addr, network)
	if err != nil {
		return nil, err
	}

	return txscript.PayToAddrScript(address)
}

func PayToPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(pubKeyHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()
}

func PayToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

func GenerateUnsignedPSBTHex(ins []*TxInput, outs []*TxOutput, network *chaincfg.Params) (string, error) {
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	var inputs []*wire.OutPoint
	var nSequences []uint32
	for _, in := range ins {
		txHash, err := chainhash.NewHashFromStr(in.TxId)
		if err != nil {
			return "", err
		}
		inputs = append(inputs, wire.NewOutPoint(txHash, in.VOut))

		nSequences = append(nSequences, in.Sequence|wire.SequenceLockTimeDisabled)
	}

	var outputs []*wire.TxOut
	for _, out := range outs {
		pkScript, err := AddrToPkScript(out.Address, network)
		if err != nil {
			return "", err
		}
		outputs = append(outputs, wire.NewTxOut(out.Amount, pkScript))
	}

	p, err := psbt.New(inputs, outputs, int32(2), uint32(0), nSequences)
	if err != nil {
		return "", err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return "", err
	}

	for i, in := range ins {
		publicKeyBytes, err := hex.DecodeString(in.PublicKey)
		if err != nil {
			return "", err
		}
		prevPkScript, err := AddrToPkScript(in.Address, network)
		if err != nil {
			return "", err
		}
		if txscript.IsPayToPubKeyHash(prevPkScript) {
			prevTx := wire.NewMsgTx(2)
			txBytes, err := hex.DecodeString(in.NonWitnessUtxo)
			if err != nil {
				return "", err
			}
			if err := prevTx.Deserialize(bytes.NewReader(txBytes)); err != nil {
				return "", err
			}
			if err := updater.AddInNonWitnessUtxo(prevTx, i); err != nil {
				return "", err
			}
		} else {
			witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
			if err := updater.AddInWitnessUtxo(witnessUtxo, i); err != nil {
				return "", err
			}
			if txscript.IsPayToScriptHash(prevPkScript) {
				redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(publicKeyBytes))
				if err != nil {
					return "", err
				}
				if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
					return "", err
				}
			}
		}

		derivationPath, err := accounts.ParseDerivationPath(in.DerivationPath)
		if err != nil {
			return "", err
		}
		if err := updater.AddInBip32Derivation(in.MasterFingerprint, derivationPath, publicKeyBytes, i); err != nil {
			return "", err
		}
	}

	for i, out := range outs {
		if out.IsChange {
			derivationPath, err := accounts.ParseDerivationPath(out.DerivationPath)
			if err != nil {
				return "", err
			}
			publicKeyBytes, err := hex.DecodeString(out.PublicKey)
			if err != nil {
				return "", err
			}
			if err := updater.AddOutBip32Derivation(out.MasterFingerprint, derivationPath, publicKeyBytes, i); err != nil {
				return "", err
			}
		}
	}

	var b bytes.Buffer
	if err := p.Serialize(&b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b.Bytes()), nil
}

func ExtractTxFromSignedPSBT(psbtHex string) (string, error) {
	psbtBytes, err := hex.DecodeString(psbtHex)
	if err != nil {
		return "", err
	}
	p, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
	if err != nil {
		return "", err
	}

	if err = psbt.MaybeFinalizeAll(p); err != nil {
		return "", err
	}

	tx, err := psbt.Extract(p)

	return GetTxHex(tx)
}
