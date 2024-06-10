package bitcoin

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/okx/go-wallet-sdk/coins/bitcoin/doginals"
	"io"
	"reflect"
)

var (
	ErrInvalidPsbtHex = errors.New("invalid psbt hex")
)
var (
	emptyHash = chainhash.Hash{}
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
type TxInputs []*TxInput

func (inputs TxInputs) UtxoViewpoint(net *chaincfg.Params) (UtxoViewpoint, error) {
	view := make(UtxoViewpoint, len(inputs))
	for _, v := range inputs {
		h, err := chainhash.NewHashFromStr(v.TxId)
		if err != nil {
			return nil, err
		}
		changePkScript, err := AddrToPkScript(v.Address, net)
		if err != nil {
			return nil, err
		}
		view[wire.OutPoint{Index: v.VOut, Hash: *h}] = changePkScript
	}
	return view, nil
}

type TxOutput struct {
	Address           string
	Amount            int64
	IsChange          bool
	MasterFingerprint uint32
	DerivationPath    string
	PublicKey         string
}

type ToSignInput struct {
	Index              int    `json:"index"`
	Address            string `json:"address"`
	PublicKey          string `json:"publicKey"`
	SigHashTypes       []int  `json:"sighashTypes"`
	DisableTweakSigner bool   `json:"disableTweakSigner"`
}

type SignPsbtOption struct {
	AutoFinalized bool           `json:"autoFinalized"`
	ToSignInputs  []*ToSignInput `json:"toSignInputs"`
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

func CalcFee(ins TxInputs, outs []*TxOutput, sellerPsbt string, feeRate int64, network *chaincfg.Params) (int64, error) {
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

	view, _ := ins.UtxoViewpoint(network)
	return GetTxVirtualSizeByView(btcutil.NewTx(tx), view) * feeRate, nil
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

type OutPoint struct {
	TxId string `json:"txId"`
	VOut uint32 `json:"vOut"`
}

func DecodeFromSignedPSBT(psbtHex string) ([]*OutPoint, error) {
	ps := make([]*OutPoint, 0)
	psbtBytes, err := hex.DecodeString(psbtHex)
	if err != nil {
		return ps, ErrInvalidPsbtHex
	}
	p, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
	if err != nil {
		return ps, ErrInvalidPsbtHex
	}
	if p == nil {
		return ps, ErrInvalidPsbtHex
	}
	if p.UnsignedTx != nil {
		for _, v := range p.UnsignedTx.TxIn {
			if v.PreviousOutPoint.Hash == emptyHash {
				continue
			}
			ps = append(ps, &OutPoint{TxId: v.PreviousOutPoint.Hash.String(), VOut: v.PreviousOutPoint.Index})
		}
	}
	return ps, nil
}

type PsbtInput struct {
	TxId   string `json:"txId"`
	Amount int64  `json:"amount"`
	VOut   uint32 `json:"vOut"`
}

type PsbtOutput struct {
	VOut     uint32 `json:"vOut"`
	Amount   int64  `json:"amount"`
	Address  string `json:"address"`
	PkScript string `json:"pkScript"`
}

type PsbtInputOutputs struct {
	UnSignedTx string        `json:"un_signed_tx"`
	Input      []*PsbtInput  `json:"input"`
	Output     []*PsbtOutput `json:"output"`
}

func parsePsbt(psbtHex string) (*psbt.Packet, error) {
	psbtBytes, err := hex.DecodeString(psbtHex)
	if err == nil {
		if r, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false); err == nil {
			return r, nil
		}
	}
	psbtBytes, err = base64.StdEncoding.DecodeString(psbtHex)
	if err != nil {
		return nil, err
	}
	return psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
}

func DecodePSBTsInputOutputs(psbtHexs []string, params *chaincfg.Params) ([]*PsbtInputOutputs, error) {
	r := make([]*PsbtInputOutputs, len(psbtHexs))
	for k, psbtHex := range psbtHexs {
		p, err := parsePsbt(psbtHex)
		if err != nil || p == nil {
			return nil, ErrInvalidPsbtHex
		}
		is := make([]*PsbtInput, 0)
		outs := make([]*PsbtOutput, 0)
		if p.UnsignedTx != nil {
			for k1, v1 := range p.UnsignedTx.TxIn {
				if v1.PreviousOutPoint.Hash == emptyHash {
					continue
				}
				var value int64
				if p.Inputs[k1].NonWitnessUtxo != nil {
					index := p.UnsignedTx.TxIn[k1].PreviousOutPoint.Index
					value = p.Inputs[k1].NonWitnessUtxo.TxOut[index].Value
				}
				if p.Inputs[k1].WitnessUtxo != nil {
					value = p.Inputs[k1].WitnessUtxo.Value
				}
				is = append(is, &PsbtInput{TxId: v1.PreviousOutPoint.Hash.String(), Amount: value, VOut: v1.PreviousOutPoint.Index})
			}
			for k1, v1 := range p.UnsignedTx.TxOut {
				_, addrs, _, err := txscript.ExtractPkScriptAddrs(v1.PkScript, params)
				if err != nil {
					continue
				}
				var addr string
				if len(addrs) > 0 {
					addr = addrs[0].EncodeAddress()
				}
				outs = append(outs, &PsbtOutput{VOut: uint32(k1), Amount: v1.Value, Address: addr})
			}
			r[k] = &PsbtInputOutputs{Input: is, Output: outs}
		}
	}
	return r, nil
}

func DecodePSBTInputOutputs(psbtHex string, params *chaincfg.Params) (*PsbtInputOutputs, error) {
	p, err := parsePsbt(psbtHex)
	if err != nil || p == nil {
		return nil, ErrInvalidPsbtHex
	}
	is := make([]*PsbtInput, 0)
	outs := make([]*PsbtOutput, 0)
	var unSignedTx string
	if p.UnsignedTx != nil {
		var buf bytes.Buffer
		if err := p.UnsignedTx.Serialize(&buf); err == nil {
			unSignedTx = hex.EncodeToString(buf.Bytes())
		}
		for k1, v1 := range p.UnsignedTx.TxIn {
			if v1.PreviousOutPoint.Hash == emptyHash {
				continue
			}
			var value int64
			if p.Inputs[k1].NonWitnessUtxo != nil {
				index := p.UnsignedTx.TxIn[k1].PreviousOutPoint.Index
				value = p.Inputs[k1].NonWitnessUtxo.TxOut[index].Value
			}
			if p.Inputs[k1].WitnessUtxo != nil {
				value = p.Inputs[k1].WitnessUtxo.Value
			}
			is = append(is, &PsbtInput{TxId: v1.PreviousOutPoint.Hash.String(), Amount: value, VOut: v1.PreviousOutPoint.Index})
		}
		for k1, v1 := range p.UnsignedTx.TxOut {
			_, addrs, _, err := txscript.ExtractPkScriptAddrs(v1.PkScript, params)
			if err != nil {
				continue
			}
			var addr string
			if len(addrs) > 0 {
				addr = addrs[0].EncodeAddress()
			}
			outs = append(outs, &PsbtOutput{VOut: uint32(k1), Amount: v1.Value, PkScript: hex.EncodeToString(v1.PkScript), Address: addr})
		}
	}

	return &PsbtInputOutputs{Input: is, UnSignedTx: unSignedTx, Output: outs}, nil
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
func SignRawPSBTTransaction(psbtHex string, privKey string) (string, error) {
	psbtBytes, err := hex.DecodeString(psbtHex)
	if err != nil {
		return "", err
	}
	p, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
	if err != nil {
		return "", err
	}
	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return "", err
	}
	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for i, in := range p.UnsignedTx.TxIn {
		prevOut := &in.PreviousOutPoint
		if p.Inputs[i].NonWitnessUtxo != nil {
			prevOuts[*prevOut] = p.Inputs[i].NonWitnessUtxo.TxOut[prevOut.Index]
		}
		if p.Inputs[i].WitnessUtxo != nil {
			prevOuts[*prevOut] = p.Inputs[i].WitnessUtxo
		}
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)
	for i, pIn := range p.Inputs {
		err = signPSBTPacket(updater, privKey, i, p, prevOutputFetcher, pIn.SighashType)
		if err != nil {
			//return "", err
		}
	}

	var b bytes.Buffer
	if err := p.Serialize(&b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b.Bytes()), nil
}

func signPSBTPacket(updater *psbt.Updater, priv string, i int, packet *psbt.Packet, prevOutFetcher *txscript.MultiPrevOutFetcher, hashType txscript.SigHashType) error {
	wif, err := btcutil.DecodeWIF(priv)
	if err != nil {
		return err
	}
	privKey := wif.PrivKey

	var prevPkScript []byte
	var value int64
	if packet.Inputs[i].NonWitnessUtxo != nil {
		index := packet.UnsignedTx.TxIn[i].PreviousOutPoint.Index
		prevPkScript = packet.Inputs[i].NonWitnessUtxo.TxOut[index].PkScript
		value = packet.Inputs[i].NonWitnessUtxo.TxOut[index].Value
	}
	if packet.Inputs[i].WitnessUtxo != nil {
		prevPkScript = packet.Inputs[i].WitnessUtxo.PkScript
		value = packet.Inputs[i].WitnessUtxo.Value
	}

	if txscript.IsPayToTaproot(prevPkScript) {
		internalPubKey := schnorr.SerializePubKey(privKey.PubKey())
		updater.Upsbt.Inputs[i].TaprootInternalKey = internalPubKey

		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)
		if hashType == txscript.SigHashAll {
			hashType = txscript.SigHashDefault
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return err
		}
		witness, err := txscript.TaprootWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes,
			i, packet.Inputs[i].WitnessUtxo.Value, prevPkScript, hashType, privKey)
		if err != nil {
			return err
		}
		updater.Upsbt.Inputs[i].TaprootKeySpendSig = witness[0]
	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		if hashType == txscript.SigHashDefault {
			hashType = txscript.SigHashAll
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return err
		}
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
		if hashType == txscript.SigHashDefault {
			hashType = txscript.SigHashAll
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return err
		}
		pubKeyBytes := privKey.PubKey().SerializeCompressed()
		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)

		script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
		if err != nil {
			return err
		}

		signature, err := txscript.RawTxInWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes, i, value, script, hashType, privKey)
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

func GenerateBatchBuyingTx(ins []*TxInput, outs []*TxOutput, sellerPSBTList []string, network *chaincfg.Params) (string, error) {
	sellerIndex := len(sellerPSBTList) + 1
	var spList []*psbt.Packet
	for _, sellerPSBT := range sellerPSBTList {
		sp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPSBT)), true)
		if err != nil {
			return "", err
		}
		spList = append(spList, sp)
	}

	var inputs []*wire.OutPoint
	var nSequences []uint32
	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for i, in := range ins {
		var prevOut *wire.OutPoint
		var sequence uint32
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			nftIn := spList[i-sellerIndex].UnsignedTx.TxIn[SellerSignatureIndex]
			prevOut = &nftIn.PreviousOutPoint
			sequence = nftIn.Sequence
		} else {
			txHash, err := chainhash.NewHashFromStr(in.TxId)
			if err != nil {
				return "", err
			}
			prevOut = wire.NewOutPoint(txHash, in.VOut)

			sequence = wire.MaxTxInSequenceNum
			if in.Sequence > 0 {
				sequence = in.Sequence | wire.SequenceLockTimeDisabled
			}
		}
		inputs = append(inputs, prevOut)

		prevPkScript, err := AddrToPkScript(in.Address, network)
		if err != nil {
			return "", err
		}
		witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
		prevOuts[*prevOut] = witnessUtxo

		nSequences = append(nSequences, sequence)
	}

	var outputs []*wire.TxOut
	for i, out := range outs {
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			outputs = append(outputs, spList[i-sellerIndex].UnsignedTx.TxOut[SellerSignatureIndex])
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
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			continue
		}

		err = signInput(updater, i, in, prevOutputFetcher, txscript.SigHashAll, network)
		if err != nil {
			return "", err
		}

		err = psbt.Finalize(bp, i)
		if err != nil {
			return "", err
		}
	}

	for i, sp := range spList {
		bp.Inputs[sellerIndex+i] = sp.Inputs[SellerSignatureIndex]
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

func GenerateBatchBuyingTxPsbt(ins []*TxInput, outs []*TxOutput, sellerPSBTList []string, network *chaincfg.Params) (string, string, error) {
	sellerIndex := len(sellerPSBTList) + 1
	var spList []*psbt.Packet
	for _, sellerPSBT := range sellerPSBTList {
		sp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPSBT)), true)
		if err != nil {
			return "", "", err
		}
		spList = append(spList, sp)
	}

	var inputs []*wire.OutPoint
	var nSequences []uint32
	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for i, in := range ins {
		var prevOut *wire.OutPoint
		var sequence uint32
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			nftIn := spList[i-sellerIndex].UnsignedTx.TxIn[SellerSignatureIndex]
			prevOut = &nftIn.PreviousOutPoint
			sequence = nftIn.Sequence
		} else {
			txHash, err := chainhash.NewHashFromStr(in.TxId)
			if err != nil {
				return "", "", err
			}
			prevOut = wire.NewOutPoint(txHash, in.VOut)

			sequence = wire.MaxTxInSequenceNum
			if in.Sequence > 0 {
				sequence = in.Sequence | wire.SequenceLockTimeDisabled
			}
		}
		inputs = append(inputs, prevOut)

		prevPkScript, err := AddrToPkScript(in.Address, network)
		if err != nil {
			return "", "", err
		}
		witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
		prevOuts[*prevOut] = witnessUtxo

		nSequences = append(nSequences, sequence)
	}

	var outputs []*wire.TxOut
	for i, out := range outs {
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			outputs = append(outputs, spList[i-sellerIndex].UnsignedTx.TxOut[SellerSignatureIndex])
		} else {
			pkScript, err := AddrToPkScript(out.Address, network)
			if err != nil {
				return "", "", err
			}
			outputs = append(outputs, wire.NewTxOut(out.Amount, pkScript))
		}
	}

	bp, err := psbt.New(inputs, outputs, int32(2), uint32(0), nSequences)
	if err != nil {
		return "", "", err
	}

	updater, err := psbt.NewUpdater(bp)
	if err != nil {
		return "", "", err
	}

	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	for i, in := range ins {
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			continue
		}

		err = signInput(updater, i, in, prevOutputFetcher, txscript.SigHashAll, network)
		if err != nil {
			return "", "", err
		}

		err = psbt.Finalize(bp, i)
		if err != nil {
			return "", "", err
		}
	}

	for i, sp := range spList {
		bp.Inputs[sellerIndex+i] = sp.Inputs[SellerSignatureIndex]
	}

	err = psbt.MaybeFinalizeAll(bp)
	if err != nil {
		return "", "", err
	}

	buyerSignedTx, err := psbt.Extract(bp)
	if err != nil {
		return "", "", err
	}

	var buf bytes.Buffer
	if err := buyerSignedTx.Serialize(&buf); err != nil {
		return "", "", err
	}

	return hex.EncodeToString(buf.Bytes()), hex.EncodeToString(buf.Bytes()), nil
}

func CalcFeeForBatchBuy(ins TxInputs, outs []*TxOutput, sellerPSBTList []string, feeRate int64, network *chaincfg.Params) (int64, error) {
	txHex, err := GenerateBatchBuyingTx(ins, outs, sellerPSBTList, network)

	tx := wire.NewMsgTx(2)
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return 0, err
	}
	err = tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return 0, err
	}
	if network != nil && network.PubKeyHashAddrID == doginals.PubKeyHashAddrID && network.ScriptHashAddrID == doginals.ScriptHashAddrID {
		return doginals.GetTxVirtualSize(btcutil.NewTx(tx)) * feeRate, nil
	}
	view, _ := ins.UtxoViewpoint(network)
	return GetTxVirtualSizeByView(btcutil.NewTx(tx), view) * feeRate, nil
}

func CalcFeeForBatchBuyWithMPC(ins TxInputs, outs []*TxOutput, sellerPSBTList []string, feeRate int64, network *chaincfg.Params, ops ...string) (int64, error) {
	var txHex string
	var err error
	if len(ops) > 0 {
		txHex = ops[0]
	} else {
		txHex, err = GenerateBatchBuyingTx(ins, outs, sellerPSBTList, network)
		if err != nil {
			return 0, err
		}
	}

	tx := wire.NewMsgTx(2)
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return 0, err
	}
	err = tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return 0, err
	}
	if network != nil && network.PubKeyHashAddrID == doginals.PubKeyHashAddrID && network.ScriptHashAddrID == doginals.ScriptHashAddrID {
		return doginals.GetTxVirtualSize(btcutil.NewTx(tx)) * feeRate, nil
	}
	view, _ := ins.UtxoViewpoint(network)
	return GetTxVirtualSizeByView(btcutil.NewTx(tx), view) * feeRate, nil
}

type GenerateMPCPSbtTxRes struct {
	Psbt         string   `json:"psbt"`
	PsbtTx       string   `json:"psbtTx"`
	SignHashList []string `json:"signHashList"`
}

func GenerateMPCUnsignedListingPSBT(in *TxInput, out *TxOutput, network *chaincfg.Params) (*GenerateMPCPSbtTxRes, error) {
	txHash, err := chainhash.NewHashFromStr(in.TxId)
	if err != nil {
		return nil, err
	}
	prevOut := wire.NewOutPoint(txHash, in.VOut)
	inputs := []*wire.OutPoint{{Index: 0}, {Index: 1}, prevOut}

	pkScript, err := AddrToPkScript(out.Address, network)
	if err != nil {
		return nil, err
	}

	dummyPkScript, err := AddrToPkScript("bc1pcyj5mt2q4t4py8jnur8vpxvxxchke4pzy7tdr9yvj3u3kdfgrj6sw3rzmr", network)
	if err != nil {
		return nil, err
	}
	outputs := []*wire.TxOut{{PkScript: dummyPkScript}, {PkScript: dummyPkScript}, wire.NewTxOut(out.Amount, pkScript)}

	nSequences := []uint32{wire.MaxTxInSequenceNum, wire.MaxTxInSequenceNum, wire.MaxTxInSequenceNum}
	p, err := psbt.New(inputs, outputs, int32(2), uint32(0), nSequences)
	if err != nil {
		return nil, err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return nil, err
	}

	dummyWitnessUtxo := wire.NewTxOut(0, dummyPkScript)
	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 0)
	if err != nil {
		return nil, err
	}
	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 1)
	if err != nil {
		return nil, err
	}

	prevPkScript, err := AddrToPkScript(in.Address, network)
	if err != nil {
		return nil, err
	}
	witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
	prevOuts := map[wire.OutPoint]*wire.TxOut{
		wire.OutPoint{Index: 0}: dummyWitnessUtxo,
		wire.OutPoint{Index: 1}: dummyWitnessUtxo,
		*prevOut:                witnessUtxo,
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	sigHash, err := calcInputSigHash(updater, SellerSignatureIndex, in, prevOutputFetcher, txscript.SigHashSingle|txscript.SigHashAnyOneCanPay, network)
	if err != nil {
		return nil, err
	}
	psbtBase64, err := p.B64Encode()
	if err != nil {
		return nil, err
	}
	res := &GenerateMPCPSbtTxRes{
		PsbtTx:       psbtBase64,
		SignHashList: []string{sigHash},
	}
	return res, nil
}

func GenerateMPCSignedListingPSBT(psbtBase64 string, signature string, pubKey string) (*GenerateMPCPSbtTxRes, error) {
	p, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(psbtBase64)), true)
	if err != nil {
		return nil, err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return nil, err
	}

	err = addInputSignature(updater, SellerSignatureIndex, signature[:64], signature[64:128],
		pubKey, txscript.SigHashSingle|txscript.SigHashAnyOneCanPay)
	if err != nil {
		return nil, err
	}
	pB64, err := p.B64Encode()
	if err != nil {
		return nil, err
	}
	res := &GenerateMPCPSbtTxRes{
		PsbtTx:       pB64,
		SignHashList: nil,
	}
	return res, nil
}

func GenerateMPCUnsignedBuyingPSBT(ins []*TxInput, outs []*TxOutput, sellerPSBTList []string, network *chaincfg.Params) (*GenerateMPCPSbtTxRes, error) {
	sellerIndex := len(sellerPSBTList) + 1
	var spList []*psbt.Packet
	for _, sellerPSBT := range sellerPSBTList {
		sp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPSBT)), true)
		if err != nil {
			return nil, err
		}
		spList = append(spList, sp)
	}

	var inputs []*wire.OutPoint
	var nSequences []uint32
	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for i, in := range ins {
		var prevOut *wire.OutPoint
		var sequence uint32
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			nftIn := spList[i-sellerIndex].UnsignedTx.TxIn[SellerSignatureIndex]
			prevOut = &nftIn.PreviousOutPoint
			sequence = nftIn.Sequence
		} else {
			txHash, err := chainhash.NewHashFromStr(in.TxId)
			if err != nil {
				return nil, err
			}
			prevOut = wire.NewOutPoint(txHash, in.VOut)

			sequence = wire.MaxTxInSequenceNum
			if in.Sequence > 0 {
				sequence = in.Sequence | wire.SequenceLockTimeDisabled
			}
		}
		inputs = append(inputs, prevOut)

		prevPkScript, err := AddrToPkScript(in.Address, network)
		if err != nil {
			return nil, err
		}
		witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
		prevOuts[*prevOut] = witnessUtxo

		nSequences = append(nSequences, sequence)
	}

	var outputs []*wire.TxOut
	for i, out := range outs {
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			outputs = append(outputs, spList[i-sellerIndex].UnsignedTx.TxOut[SellerSignatureIndex])
		} else {
			pkScript, err := AddrToPkScript(out.Address, network)
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, wire.NewTxOut(out.Amount, pkScript))
		}
	}

	bp, err := psbt.New(inputs, outputs, int32(2), uint32(0), nSequences)
	if err != nil {
		return nil, err
	}

	updater, err := psbt.NewUpdater(bp)
	if err != nil {
		return nil, err
	}

	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	var sigHashList []string
	for i, in := range ins {
		if sellerIndex <= i && i < sellerIndex+len(spList) {
			continue
		}

		sigHash, err := calcInputSigHash(updater, i, in, prevOutputFetcher, txscript.SigHashAll, network)
		if err != nil {
			return nil, err
		}

		sigHashList = append(sigHashList, sigHash)
	}

	for i, sp := range spList {
		bp.Inputs[sellerIndex+i] = sp.Inputs[SellerSignatureIndex]
	}

	psbtBase64, err := bp.B64Encode()
	if err != nil {
		return nil, err
	}
	res := &GenerateMPCPSbtTxRes{
		PsbtTx:       psbtBase64,
		SignHashList: sigHashList,
	}
	return res, nil
}

func GenerateMPCSignedBuyingTx(psbtBase64 string, signatures []string, pubKey string, batchSize int) (*GenerateMPCPSbtTxRes, error) {
	p, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(psbtBase64)), true)
	if err != nil {
		return nil, err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return nil, err
	}

	sellerIndex := batchSize + 1
	signatureIndex := 0
	for i := range p.UnsignedTx.TxIn {
		if sellerIndex <= i && i < sellerIndex+batchSize {
			continue
		}

		err = addInputSignature(updater, i, signatures[signatureIndex][:64], signatures[signatureIndex][64:128], pubKey, txscript.SigHashAll)
		if err != nil {
			return nil, err
		}
		signatureIndex++
	}
	err = psbt.MaybeFinalizeAll(p)
	if err != nil {
		return nil, err
	}
	buyerSignedTx, err := psbt.Extract(p)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := buyerSignedTx.Serialize(&buf); err != nil {
		return nil, err
	}
	b := bytes.Buffer{}
	err = p.Serialize(&b)
	if err != nil {
		return nil, err
	}
	res := &GenerateMPCPSbtTxRes{
		Psbt:         hex.EncodeToString(b.Bytes()),
		PsbtTx:       hex.EncodeToString(buf.Bytes()),
		SignHashList: nil,
	}
	return res, nil
}

func calcInputSigHash(updater *psbt.Updater, i int, in *TxInput, prevOutFetcher *txscript.MultiPrevOutFetcher, hashType txscript.SigHashType, network *chaincfg.Params) (string, error) {
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
		err = prevTx.Deserialize(bytes.NewReader(txBytes))
		if err != nil {
			return "", err
		}
		err = updater.AddInNonWitnessUtxo(prevTx, i)
		if err != nil {
			return "", err
		}
	} else {
		witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
		err = updater.AddInWitnessUtxo(witnessUtxo, i)
		if err != nil {
			return "", err
		}
	}

	err = updater.AddInSighashType(hashType, i)
	if err != nil {
		return "", err
	}

	tx := updater.Upsbt.UnsignedTx
	txSigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)
	pubKeyBytes, err := hex.DecodeString(in.PublicKey)
	var sigHash []byte
	if txscript.IsPayToTaproot(prevPkScript) {
		if hashType == txscript.SigHashAll {
			hashType = txscript.SigHashDefault
		}
		sigHash, err = txscript.CalcTaprootSignatureHash(txSigHashes, hashType, tx, i, prevOutFetcher)
		if err != nil {
			return "", err
		}

		updater.Upsbt.Inputs[i].TaprootInternalKey = pubKeyBytes[1:]
	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		sigHash, err = txscript.CalcSignatureHash(prevPkScript, hashType, tx, i)
		if err != nil {
			return "", err
		}
	} else {
		script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
		if err != nil {
			return "", err
		}
		sigHash, err = txscript.CalcWitnessSigHash(script, txSigHashes, hashType, tx, i, in.Amount)
		if err != nil {
			return "", err
		}

		if txscript.IsPayToScriptHash(prevPkScript) {
			redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return "", err
			}
			err = updater.AddInRedeemScript(redeemScript, i)
			if err != nil {
				return "", err
			}
		}
	}

	return hex.EncodeToString(sigHash), nil
}

func addInputSignature(updater *psbt.Updater, i int, rHex string, sHex string, pubKeyHex string, hashType txscript.SigHashType) error {
	rBytes, err := hex.DecodeString(rHex)
	if err != nil {
		return err
	}
	sBytes, err := hex.DecodeString(sHex)
	if err != nil {
		return err
	}

	r := new(btcec.ModNScalar)
	r.SetByteSlice(rBytes)
	s := new(btcec.ModNScalar)
	s.SetByteSlice(sBytes)
	signature := append(ecdsa.NewSignature(r, s).Serialize(), byte(hashType))
	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return err
	}
	_, err = updater.Sign(i, signature, pubKey, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func CalcInputSigHashForUnsignedPSBT(updater *psbt.Updater, i int, packet *psbt.Packet, prevOutFetcher *txscript.MultiPrevOutFetcher, hashType txscript.SigHashType, publicKey string) (string, error) {

	var prevPkScript []byte
	var value int64
	if packet.Inputs[i].NonWitnessUtxo != nil {
		index := packet.UnsignedTx.TxIn[i].PreviousOutPoint.Index
		prevPkScript = packet.Inputs[i].NonWitnessUtxo.TxOut[index].PkScript
		value = packet.Inputs[i].NonWitnessUtxo.TxOut[index].Value
	}
	if packet.Inputs[i].WitnessUtxo != nil {
		prevPkScript = packet.Inputs[i].WitnessUtxo.PkScript
		value = packet.Inputs[i].WitnessUtxo.Value
	}

	tx := updater.Upsbt.UnsignedTx
	txSigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	var sigHash []byte
	if txscript.IsPayToTaproot(prevPkScript) {
		if hashType == txscript.SigHashAll {
			hashType = txscript.SigHashDefault
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return "", err
		}
		sigHash, err = txscript.CalcTaprootSignatureHash(txSigHashes, hashType, tx, i, prevOutFetcher)
		if err != nil {
			return "", err
		}

		updater.Upsbt.Inputs[i].TaprootInternalKey = pubKeyBytes[1:]
	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		if hashType == txscript.SigHashDefault {
			hashType = txscript.SigHashAll
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return "", err
		}
		sigHash, err = txscript.CalcSignatureHash(prevPkScript, hashType, tx, i)
		if err != nil {
			return "", err
		}
	} else {
		if hashType == txscript.SigHashDefault {
			hashType = txscript.SigHashAll
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return "", err
		}
		script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
		if err != nil {
			return "", err
		}
		sigHash, err = txscript.CalcWitnessSigHash(script, txSigHashes, hashType, tx, i, value)
		if err != nil {
			return "", err
		}
		if txscript.IsPayToScriptHash(prevPkScript) {
			redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return "", err
			}
			err = updater.AddInRedeemScript(redeemScript, i)
			if err != nil {
				return "", err
			}
		}
	}

	return hex.EncodeToString(sigHash), nil
}

func GenerateMPCUnsignedPSBT(psbtStr string, pubKeyHex string) (*GenerateMPCPSbtTxRes, error) {
	p, err := GetPsbtFromString(psbtStr)
	if err != nil {
		return nil, err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return nil, err
	}
	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for i, in := range p.UnsignedTx.TxIn {
		prevOut := &in.PreviousOutPoint
		if p.Inputs[i].NonWitnessUtxo != nil {
			prevOuts[*prevOut] = p.Inputs[i].NonWitnessUtxo.TxOut[prevOut.Index]
		}
		if p.Inputs[i].WitnessUtxo != nil {
			prevOuts[*prevOut] = p.Inputs[i].WitnessUtxo
		}
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	var sigHashList []string
	for i, pIn := range p.Inputs {
		sigHash, err := CalcInputSigHashForUnsignedPSBT(updater, i, p, prevOutputFetcher, pIn.SighashType, pubKeyHex)
		if err != nil {
			hash, err := GetRandomHash()
			if err != nil {
				sigHash = pubKeyHex
			}
			sigHash = hash
		}
		sigHashList = append(sigHashList, sigHash)
	}

	psbtBase64, err := p.B64Encode()
	if err != nil {
		return nil, err
	}
	res := &GenerateMPCPSbtTxRes{
		PsbtTx:       psbtBase64,
		SignHashList: sigHashList,
	}
	return res, nil
}

func GetRandomHash() (string, error) {
	s := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, s); err != nil {
		return "", err
	}
	hashBs := chainhash.HashB(s)
	return "ffffffff" + hex.EncodeToString(hashBs[:28]), nil
}

func GenerateMPCSignedPSBT(psbtStr string, pubKeyHex string, signatureList []string) (*GenerateMPCPSbtTxRes, error) {
	unsignedPsbtInfor, err := GenerateMPCUnsignedPSBT(psbtStr, pubKeyHex)
	rawPsbtStr := unsignedPsbtInfor.PsbtTx
	p, err := GetPsbtFromString(rawPsbtStr)
	if err != nil {
		return nil, err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return nil, err
	}
	for i := range p.UnsignedTx.TxIn {
		signHashType := p.Inputs[i].SighashType
		if signHashType == txscript.SigHashDefault {
			signHashType = txscript.SigHashAll
		}
		err = addInputSignature(updater, i, signatureList[i][:64], signatureList[i][64:128], pubKeyHex, signHashType)
		if err != nil {
			//return nil, err
		}
	}
	buf := &bytes.Buffer{}
	err = p.Serialize(buf)
	if err != nil {
		return nil, err
	}
	res := &GenerateMPCPSbtTxRes{
		PsbtTx:       hex.EncodeToString(buf.Bytes()),
		SignHashList: nil,
	}
	return res, nil
}

func SignPsbtWithKeyPathAndScriptPath(psbtHex string, privKey string, network *chaincfg.Params, option *SignPsbtOption) (string, error) {
	psbtBytes, err := hex.DecodeString(psbtHex)
	if err != nil {
		return "", err
	}
	p, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
	if err != nil {
		return "", err
	}
	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return "", err
	}
	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for i, in := range p.UnsignedTx.TxIn {
		prevOut := &in.PreviousOutPoint
		if p.Inputs[i].NonWitnessUtxo != nil {
			prevOuts[*prevOut] = p.Inputs[i].NonWitnessUtxo.TxOut[prevOut.Index]
		}
		if p.Inputs[i].WitnessUtxo != nil {
			prevOuts[*prevOut] = p.Inputs[i].WitnessUtxo
		}
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)
	m := make(map[int]*ToSignInput)
	if option != nil {
		for _, v := range option.ToSignInputs {
			m[v.Index] = v
		}
	}
	for i, pIn := range p.Inputs {
		toSignInput, ok := m[i]
		if len(m) > 0 && !ok {
			continue
		}
		err = signPsbtWithKeyPathAndScriptPath(updater, privKey, i, p, prevOutputFetcher, pIn.SighashType, toSignInput, network)
		if err != nil {
			continue
		}
		if option != nil && !option.AutoFinalized {
			continue
		}
		err = psbt.Finalize(p, i)
		if err != nil {
		}
	}

	var b bytes.Buffer
	if err := p.Serialize(&b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b.Bytes()), nil
}

func signPsbtWithKeyPathAndScriptPath(updater *psbt.Updater, priv string, i int, packet *psbt.Packet, prevOutFetcher *txscript.MultiPrevOutFetcher, hashType txscript.SigHashType, toSignInput *ToSignInput, network *chaincfg.Params) error {
	wif, err := btcutil.DecodeWIF(priv)
	if err != nil {
		return err
	}
	privKey := wif.PrivKey
	if toSignInput != nil && toSignInput.PublicKey != "" {
		pub := hex.EncodeToString(privKey.PubKey().SerializeCompressed())
		if pub != toSignInput.PublicKey {
			return fmt.Errorf("invlid public key %s", toSignInput.PublicKey)
		}
	}

	var prevPkScript []byte
	var value int64
	if packet.Inputs[i].NonWitnessUtxo != nil {
		index := packet.UnsignedTx.TxIn[i].PreviousOutPoint.Index
		prevPkScript = packet.Inputs[i].NonWitnessUtxo.TxOut[index].PkScript
		value = packet.Inputs[i].NonWitnessUtxo.TxOut[index].Value
	}
	if packet.Inputs[i].WitnessUtxo != nil {
		prevPkScript = packet.Inputs[i].WitnessUtxo.PkScript
		value = packet.Inputs[i].WitnessUtxo.Value
	}
	if toSignInput != nil && len(toSignInput.SigHashTypes) > 0 {
		hashType = txscript.SigHashType(toSignInput.SigHashTypes[0])
	}
	if toSignInput != nil && toSignInput.Address != "" {
		pks, err := AddrToPkScript(toSignInput.Address, network)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(pks, prevPkScript) {
			return fmt.Errorf("invalid address %s", toSignInput.Address)
		}
	}
	if txscript.IsPayToTaproot(prevPkScript) {
		// ket path only
		internalPubKey := schnorr.SerializePubKey(privKey.PubKey())
		updater.Upsbt.Inputs[i].TaprootInternalKey = internalPubKey

		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)
		if hashType == txscript.SigHashAll {
			hashType = txscript.SigHashDefault
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return err
		}
		witness, err := txscript.TaprootWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes,
			i, packet.Inputs[i].WitnessUtxo.Value, prevPkScript, hashType, privKey)
		if err != nil {
			return err
		}
		updater.Upsbt.Inputs[i].TaprootKeySpendSig = witness[0]

		// script path but key path spend
		rootHash := updater.Upsbt.Inputs[i].TaprootMerkleRoot
		if rootHash != nil {
			if toSignInput != nil {
				if toSignInput.DisableTweakSigner {
					// fake root and it is invalid
					rootHash = []byte{}
				}
			}
			sig, err := txscript.RawTxInTaprootSignature(updater.Upsbt.UnsignedTx, sigHashes,
				i, packet.Inputs[i].WitnessUtxo.Value, prevPkScript, rootHash, hashType, privKey)
			if err != nil {
				return err
			}
			updater.Upsbt.Inputs[i].TaprootKeySpendSig = sig
		} else {
			if len(updater.Upsbt.Inputs[i].TaprootLeafScript) > 0 {
				// btcd only support one leaf till now
				tapLeaves := updater.Upsbt.Inputs[i].TaprootLeafScript
				taprootScriptSpendSignatures := make([]*psbt.TaprootScriptSpendSig, 0)
				for _, leaf := range tapLeaves {
					tapLeaf := txscript.TapLeaf{
						LeafVersion: leaf.LeafVersion,
						Script:      leaf.Script,
					}
					sig, err := txscript.RawTxInTapscriptSignature(updater.Upsbt.UnsignedTx, sigHashes,
						i, packet.Inputs[i].WitnessUtxo.Value, prevPkScript, tapLeaf, hashType, privKey)
					if err != nil {
						return err
					}
					tapHash := tapLeaf.TapHash()
					tapLeafSignature := &psbt.TaprootScriptSpendSig{
						XOnlyPubKey: internalPubKey,
						LeafHash:    tapHash.CloneBytes(),
						Signature:   sig,
						SigHash:     hashType,
					}
					taprootScriptSpendSignatures = append(taprootScriptSpendSignatures, tapLeafSignature)
				}
				updater.Upsbt.Inputs[i].TaprootInternalKey = nil
				updater.Upsbt.Inputs[i].TaprootKeySpendSig = nil
				// remove duplicate
				updater.Upsbt.Inputs[i].TaprootScriptSpendSig = append(updater.Upsbt.Inputs[i].TaprootScriptSpendSig, taprootScriptSpendSignatures...)
				CheckDuplicateOfUpdater(updater, i)
			}
		}

	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		if hashType == txscript.SigHashDefault {
			hashType = txscript.SigHashAll
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return err
		}
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
		if hashType == txscript.SigHashDefault {
			hashType = txscript.SigHashAll
		}
		err = updater.AddInSighashType(hashType, i)
		if err != nil {
			return err
		}
		pubKeyBytes := privKey.PubKey().SerializeCompressed()
		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)

		script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
		if err != nil {
			return err
		}

		signature, err := txscript.RawTxInWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes, i, value, script, hashType, privKey)
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

func CheckDuplicateOfUpdater(updater *psbt.Updater, index int) {
	signatures := updater.Upsbt.Inputs[index].TaprootScriptSpendSig
	m := map[string]*psbt.TaprootScriptSpendSig{}
	signs := make([]*psbt.TaprootScriptSpendSig, 0)
	for _, v := range signatures {
		key := append(v.XOnlyPubKey, v.LeafHash...)
		keyHex := hex.EncodeToString(key)
		_, ok := m[keyHex]
		if !ok {
			m[keyHex] = v
			signs = append(signs, v)
		}
	}
	updater.Upsbt.Inputs[index].TaprootScriptSpendSig = signs
}
