package atomical

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
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
	PkScript          string
	Amount            int64
	IsChange          bool
	MasterFingerprint uint32
	DerivationPath    string
	PublicKey         string
}

const SellerSignatureIndex = 1

func GenerateAtomicalSignedListingPSBTBase64(in *TxInput, out *TxOutput, network *chaincfg.Params) (string, error) {
	txHash, err := chainhash.NewHashFromStr(in.TxId)
	if err != nil {
		return "", err
	}
	prevOut := wire.NewOutPoint(txHash, in.VOut)
	inputs := []*wire.OutPoint{{Index: 0}, prevOut}
	var pkScript []byte
	if len(out.PkScript) > 0 {
		pkScript, err = hex.DecodeString(out.PkScript)
		if err != nil {
			return "", err
		}
	} else {
		pkScript, err = AddrToPkScript(out.Address, network)
		if err != nil {
			return "", err
		}
	}

	dummyPkScript, _ := AddrToPkScript("bc1p9f7eu4devpp6d84lv8hwdmccmdfqhmnvvtx48hnyleuje3m6769qy4evm8", network)
	outputs := []*wire.TxOut{{PkScript: dummyPkScript}, wire.NewTxOut(out.Amount, pkScript)}

	nSequences := []uint32{wire.MaxTxInSequenceNum, wire.MaxTxInSequenceNum}
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

	prevPkScript, err := AddrToPkScript(in.Address, network)
	if err != nil {
		return "", err
	}
	witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
	prevOuts := map[wire.OutPoint]*wire.TxOut{
		wire.OutPoint{Index: 0}: dummyWitnessUtxo,
		*prevOut:                witnessUtxo,
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	err = signInput(updater, SellerSignatureIndex, in, prevOutputFetcher, txscript.SigHashSingle|txscript.SigHashAnyOneCanPay, network)
	if err != nil {
		return "", err
	}

	return p.B64Encode()
}

func GenerateAtomicalSignedBuyingTx(ins []*TxInput, outs []*TxOutput, dustSize, feePerB int64, sellerPsbt string, network *chaincfg.Params) (int64, string, error) {
	// Include change and calculate whether the handling fee is sufficient
	totalInput, totalOutput, vsize, err := calFee(ins, outs, sellerPsbt, network)
	if err != nil {
		return 0, "", err
	}
	fee := vsize * feePerB
	if totalInput-totalOutput > fee && totalInput-totalOutput-fee >= dustSize {
		outs[len(outs)-1].Amount = totalInput - totalOutput - fee
	} else {
		outs = outs[0 : len(outs)-1]
		totalInput, totalOutput, vsize, err = calFee(ins, outs, sellerPsbt, network)
		feeWithoutChange := vsize * feePerB
		if totalInput-totalOutput < feeWithoutChange {
			return totalOutput + fee, "", errors.New(ErInsufficientBalance)
		}
		fee = feeWithoutChange
	}
	tx, err := generateBuyPsbt(ins, outs, sellerPsbt, network, false)
	if err != nil {
		return 0, "", err
	}
	return fee, tx, err
}

func calFee(ins []*TxInput, outs []*TxOutput, sellerPsbt string, network *chaincfg.Params) (int64, int64, int64, error) {
	sp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPsbt)), true)
	if err != nil {
		return 0, 0, 0, err
	}
	dummyPrivKey := "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22"

	txBuild := bitcoin.NewTxBuild(2, network)
	txBuild.AddInput2(ins[0].TxId, ins[0].VOut, dummyPrivKey, ins[0].Address, ins[0].Amount)
	txBuild.AddOutput2(outs[0].Address, outs[0].PkScript, outs[0].Amount)

	txBuild.AddInput2(sp.UnsignedTx.TxIn[SellerSignatureIndex].PreviousOutPoint.Hash.String(), sp.UnsignedTx.TxIn[SellerSignatureIndex].PreviousOutPoint.Index, dummyPrivKey, ins[1].Address, ins[1].Amount)
	txBuild.AddOutput2("", hex.EncodeToString(sp.UnsignedTx.TxOut[SellerSignatureIndex].PkScript), sp.UnsignedTx.TxOut[SellerSignatureIndex].Value)

	for i := 2; i < len(ins); i++ {
	}

	for i := 2; i < len(outs); i++ {
		txBuild.AddOutput2(outs[i].Address, outs[i].PkScript, outs[i].Amount)
	}
	tx, err := txBuild.Build()
	if err != nil {
		return 0, 0, 0, err
	}
	vsize := bitcoin.GetTxVirtualSize(btcutil.NewTx(tx))
	return txBuild.TotalInputAmount(), txBuild.TotalOutputAmount(), vsize, nil
}

func generateBuyPsbt(ins []*TxInput, outs []*TxOutput, sellerPsbt string, network *chaincfg.Params, finalize bool) (string, error) {
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
			var pkScript []byte
			if len(out.PkScript) > 0 {
				pkScript, err = hex.DecodeString(out.PkScript)
				if err != nil {
					return "", err
				}
			} else {
				pkScript, err = AddrToPkScript(out.Address, network)
				if err != nil {
					return "", err
				}
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

		err = signInput(updater, i, in, prevOutputFetcher, txscript.SigHashAll, network)
		if err != nil {
			return "", err
		}

		if finalize {
			err = psbt.Finalize(bp, i)
			if err != nil {
				return "", err
			}
		}
	}
	bp.Inputs[SellerSignatureIndex] = sp.Inputs[SellerSignatureIndex]
	if !finalize {
		r, err := bp.B64Encode()
		return r, err
	}

	err = psbt.MaybeFinalizeAll(bp)
	if err != nil {
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
		err = prevTx.Deserialize(bytes.NewReader(txBytes))
		if err != nil {
			return err
		}
		err = updater.AddInNonWitnessUtxo(prevTx, i)
		if err != nil {
			return err
		}
	} else {
		witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
		err = updater.AddInWitnessUtxo(witnessUtxo, i)
		if err != nil {
			return err
		}
	}

	err = updater.AddInSighashType(hashType, i)
	if err != nil {
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
		signOutcome, err := updater.Sign(i, signature, privKey.PubKey().SerializeCompressed(), nil, nil)
		if err != nil {
			return err
		}
		if signOutcome != psbt.SignSuccesful {
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

		signOutcome, err := updater.Sign(i, signature, pubKeyBytes, nil, nil)
		if err != nil {
			return err
		}
		if signOutcome != psbt.SignSuccesful {
			return err
		}
	}
	return nil
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
