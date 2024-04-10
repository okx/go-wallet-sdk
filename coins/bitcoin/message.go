package bitcoin

import (
	"bytes"
	"crypto/sha256"
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
	"github.com/okx/go-wallet-sdk/crypto"
	"io"
	"math/big"
)

const (
	compactSigSize = 65

	compactSigMagicOffset = 27

	compactSigCompPubKey = 4

	ScriptSigPrefix     = "0020"
	Bip0322Tag          = "BIP0322-signed-message"
	Bip0322Opt          = "bip0322-simple"
	SignedMessagePrefix = "Bitcoin Signed Message:\n"
)

func Bip0322Hash(message string) string {
	tagHash := sha256.Sum256([]byte(Bip0322Tag))
	result := sha256.Sum256(append(tagHash[:], append(tagHash[:], []byte(message)...)...))
	return hex.EncodeToString(result[:])
}

func SignBip0322(message string, address string, privateKey string) (string, error) {
	network := &chaincfg.MainNetParams
	txId, err := BuildToSpend(message, address, network)
	if err != nil {
		return "", err
	}

	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return "", err
	}

	prevOut := wire.NewOutPoint(txHash, 0)
	inputs := []*wire.OutPoint{prevOut}
	script, err := hex.DecodeString("6a")
	if err != nil {
		return "", err
	}
	outputs := []*wire.TxOut{wire.NewTxOut(0, script)}
	nSequences := []uint32{uint32(0)}
	p, err := NewPsbt(inputs, outputs, int32(0), uint32(0), nSequences, Bip0322Opt)
	if err != nil {
		return "", err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return "", err
	}

	dummyPkScript, err := AddrToPkScript(address, network)
	if err != nil {
		return "", err
	}
	dummyWitnessUtxo := wire.NewTxOut(0, dummyPkScript)

	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 0)
	if err != nil {
		return "", err
	}

	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", err
	}
	privKey := wif.PrivKey

	prevPkScript, err := AddrToPkScript(address, network)
	witnessUtxo := wire.NewTxOut(0, prevPkScript)
	prevOuts := map[wire.OutPoint]*wire.TxOut{
		*prevOut: witnessUtxo,
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	if txscript.IsPayToTaproot(prevPkScript) {
		internalPubKey := schnorr.SerializePubKey(privKey.PubKey())
		updater.Upsbt.Inputs[0].TaprootInternalKey = internalPubKey

		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutputFetcher)
		hashType := txscript.SigHashDefault
		witness, err := txscript.TaprootWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes,
			0, 0, prevPkScript, hashType, privKey)
		if err != nil {
			return "", err
		}
		updater.Upsbt.Inputs[0].TaprootKeySpendSig = witness[0]
	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		signature, err := txscript.RawTxInSignature(updater.Upsbt.UnsignedTx, 0, prevPkScript, txscript.SigHashDefault, privKey)
		if err != nil {
			return "", err
		}
		signOutcome, err := updater.Sign(0, signature, privKey.PubKey().SerializeCompressed(), nil, nil)
		if err != nil {
			return "", err
		}
		if signOutcome != psbt.SignSuccesful {
			return "", err
		}
	} else {
		pubKeyBytes := privKey.PubKey().SerializeCompressed()
		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutputFetcher)

		script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
		if err != nil {
			return "", err
		}
		signature, err := txscript.RawTxInWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes, 0, 0, script, txscript.SigHashAll, privKey)
		if err != nil {
			return "", err
		}

		if txscript.IsPayToScriptHash(prevPkScript) {
			redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return "", err
			}
			err = updater.AddInRedeemScript(redeemScript, 0)
			if err != nil {
				return "", err
			}
		}

		signOutcome, err := updater.Sign(0, signature, pubKeyBytes, nil, nil)
		if err != nil {
			return "", err
		}
		if signOutcome != psbt.SignSuccesful {
			return "", err
		}
	}
	var byts bytes.Buffer
	if err := p.Serialize(&byts); err != nil {
		return "", err
	}
	txHex, err := ExtractTxFromSignedPSBTBIP322(hex.EncodeToString(byts.Bytes()))
	if err != nil {
		return "", err
	}

	txHexBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(txHexBytes), nil
}

func BuildToSpend(message string, address string, network *chaincfg.Params) (string, error) {
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	outputScript, err := AddrToPkScript(address, network)
	if err != nil {
		return "", nil
	}

	tx := wire.NewMsgTx(0)

	sequence := uint32(0)
	scriptSig, err := hex.DecodeString(ScriptSigPrefix + Bip0322Hash(message))
	if err != nil {
		return "", err
	}

	in := wire.NewTxIn(&wire.OutPoint{Index: wire.MaxPrevOutIndex}, scriptSig, nil)
	in.Sequence = sequence
	tx.AddTxIn(in)
	tx.AddTxOut(wire.NewTxOut(0, outputScript))

	return tx.TxHash().String(), nil
}

func ExtractTxFromSignedPSBTBIP322(psbtHex string) (string, error) {
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

	return GetTxHexBIP322(tx)
}
func GetTxHexBIP322(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := BtcEncodeBip322(&buf, 0, wire.WitnessEncoding, tx); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func BtcEncodeBip322(w io.Writer, pver uint32, enc wire.MessageEncoding, msg *wire.MsgTx) error {
	doWitness := enc == wire.WitnessEncoding && msg.HasWitness()
	if doWitness {
		for _, ti := range msg.TxIn {
			err := wire.WriteVarInt(w, pver, uint64(len(ti.Witness)))
			if err != nil {
				return err
			}
			for _, item := range ti.Witness {
				err = wire.WriteVarBytes(w, pver, item)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func MPCUnsignedBip0322(message string, address string, publicKey string, network *chaincfg.Params) (*GenerateMPCPSbtTxRes, error) {
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	txId, err := BuildToSpend(message, address, network)
	if err != nil {
		return nil, err
	}

	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return nil, err
	}

	prevOut := wire.NewOutPoint(txHash, 0)
	inputs := []*wire.OutPoint{prevOut}
	script, err := hex.DecodeString("6a")
	if err != nil {
		return nil, err
	}
	outputs := []*wire.TxOut{wire.NewTxOut(0, script)}
	nSequences := []uint32{uint32(0)}
	p, err := NewPsbt(inputs, outputs, int32(0), uint32(0), nSequences, Bip0322Opt)
	if err != nil {
		return nil, err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return nil, err
	}

	dummyPkScript, err := AddrToPkScript(address, network)
	if err != nil {
		return nil, err
	}
	dummyWitnessUtxo := wire.NewTxOut(0, dummyPkScript)

	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 0)
	if err != nil {
		return nil, err
	}

	prevPkScript, err := AddrToPkScript(address, network)
	if err != nil {
		return nil, err
	}
	witnessUtxo := wire.NewTxOut(0, prevPkScript)
	prevOuts := map[wire.OutPoint]*wire.TxOut{
		*prevOut: witnessUtxo,
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)
	publicKeyBs, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	pub, err := btcec.ParsePubKey(publicKeyBs)
	if err != nil {
		return nil, err
	}
	var signHash []byte
	if txscript.IsPayToTaproot(prevPkScript) {
		internalPubKey := schnorr.SerializePubKey(pub)
		updater.Upsbt.Inputs[0].TaprootInternalKey = internalPubKey

		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutputFetcher)
		hashType := txscript.SigHashDefault

		err = updater.AddInSighashType(hashType, 0)
		if err != nil {
			return nil, err
		}
		signHash, err = txscript.CalcTaprootSignatureHash(sigHashes, hashType, updater.Upsbt.UnsignedTx, 0, prevOutputFetcher)
		if err != nil {
			return nil, err
		}
	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		hashType := txscript.SigHashAll
		err = updater.AddInSighashType(hashType, 0)
		if err != nil {
			return nil, err
		}
		signHash, err = txscript.CalcSignatureHash(prevPkScript, hashType, updater.Upsbt.UnsignedTx, 0)
		if err != nil {
			return nil, err
		}
	} else {
		hashType := txscript.SigHashAll
		err = updater.AddInSighashType(hashType, 0)
		if err != nil {
			return nil, err
		}
		pubKeyBytes, err := hex.DecodeString(publicKey)
		if err != nil {
			return nil, err
		}
		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutputFetcher)

		script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
		if err != nil {
			return nil, err
		}
		signHash, err = txscript.CalcWitnessSigHash(script, sigHashes, hashType, updater.Upsbt.UnsignedTx, 0, 0)
		if err != nil {
			return nil, err
		}
		if txscript.IsPayToScriptHash(prevPkScript) {
			redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return nil, err
			}
			err = updater.AddInRedeemScript(redeemScript, 0)
			if err != nil {
				return nil, err
			}
		}
	}
	var byts bytes.Buffer
	if err := p.Serialize(&byts); err != nil {
		return nil, err
	}
	res := &GenerateMPCPSbtTxRes{
		PsbtTx:       hex.EncodeToString(byts.Bytes()),
		SignHashList: []string{hex.EncodeToString(signHash)},
	}
	return res, nil
}

func MPCSignedBip0322(message string, address string, publicKey string, signatureList []string, network *chaincfg.Params) (*GenerateMPCPSbtTxRes, error) {
	unsignedPsbtInfor, err := MPCUnsignedBip0322(message, address, publicKey, network)
	if err != nil {
		return nil, err
	}
	signedPsbtInfor, err := GenerateMPCSignedPSBT(unsignedPsbtInfor.PsbtTx, publicKey, signatureList)
	if err != nil {
		return nil, err
	}
	txHex, err := ExtractTxFromSignedPSBTBIP322(signedPsbtInfor.PsbtTx)
	if err != nil {
		return nil, err
	}

	txHexBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	res := &GenerateMPCPSbtTxRes{
		PsbtTx:       base64.StdEncoding.EncodeToString(txHexBytes),
		SignHashList: nil,
	}
	return res, nil
}

func SignMessage(wif string, message string) (string, error) {
	var buf bytes.Buffer
	err := wire.WriteVarString(&buf, 0, SignedMessagePrefix)
	if err != nil {
		return "", err
	}
	err = wire.WriteVarString(&buf, 0, message)
	if err != nil {
		return "", err
	}
	messageHash := chainhash.DoubleHashB(buf.Bytes())
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return "", err
	}
	sig, err := ecdsa.SignCompact(w.PrivKey, messageHash, w.CompressPubKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func MPCUnsignedMessage(message string) string {
	var buf bytes.Buffer
	err := wire.WriteVarString(&buf, 0, SignedMessagePrefix)
	if err != nil {
		return ""
	}
	err = wire.WriteVarString(&buf, 0, message)
	if err != nil {
		return ""
	}
	messageHash := chainhash.DoubleHashB(buf.Bytes())
	return hex.EncodeToString(messageHash)
}

func MPCSignedMessageCompat(message string, signature string, publicKeyHex string, network *chaincfg.Params) (string, error) {
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	msgHash := MPCUnsignedMessage(message)
	r, ok := new(big.Int).SetString(signature[0:64], 16)
	if !ok {
		return "", errors.New("parse r failed")
	}
	s, ok := new(big.Int).SetString(signature[64:128], 16)
	if !ok {
		return "", errors.New("parse s failed")
	}
	publicKeyBs, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", err
	}
	pub, err := btcec.ParsePubKey(publicKeyBs)
	if err != nil {
		return "", err
	}
	msgHashBs, err := hex.DecodeString(msgHash)
	if err != nil {
		return "", err
	}
	addressPubKey, err := btcutil.NewAddressPubKey(publicKeyBs, network)
	if err != nil {
		return "", err
	}
	isCompressedKey := true
	if addressPubKey.Format() == btcutil.PKFUncompressed {
		isCompressedKey = false
	}
	sig, err := crypto.SignCompact(btcec.S256(), r, s, *pub, msgHashBs, isCompressedKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func MPCSignedMessage(signature string, publicKeyHex string, network *chaincfg.Params) (string, error) {
	if network == nil {
		network = &chaincfg.MainNetParams
	}

	sigBs, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}
	if len(sigBs) < compactSigSize {
		return "", errors.New("invalid length of signature")
	}
	// v r s
	compactSigRecoveryCode := sigBs[64]
	r := new(btcec.ModNScalar)
	r.SetByteSlice(sigBs[:32])
	s := new(btcec.ModNScalar)
	s.SetByteSlice(sigBs[32:64])

	var b [compactSigSize]byte
	// v r s
	b[0] = compactSigRecoveryCode
	r.PutBytesUnchecked(b[1:33])
	s.PutBytesUnchecked(b[33:65])
	return base64.StdEncoding.EncodeToString(b[:]), nil
}

func VerifyMessage(signatureStr, message, publicKeyHex, signType string) (bool, error) {
	if signType == "bip322-simple" {
		return false, nil
	}
	var signature []byte
	var err error
	if IsHexString(signatureStr) {
		signature, err = hex.DecodeString(signatureStr)
	} else {
		signature, err = base64.StdEncoding.DecodeString(signatureStr)
	}
	if err != nil {
		return false, err
	}
	messageHash := MPCUnsignedMessage(message)

	h, err := hex.DecodeString(messageHash)
	if err != nil {
		return false, err
	}
	p, ok, err := ecdsa.RecoverCompact(signature, h)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("invalid signature")
	}
	ph := hex.EncodeToString(p.SerializeCompressed())
	if ph != publicKeyHex {
		return false, fmt.Errorf("invalid public Key %s, given publicKeyHex %s", ph, publicKeyHex)
	}
	return ok, nil
}
