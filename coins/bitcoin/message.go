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
	"reflect"
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

var (
	ErrNonSupportedAddrType = errors.New("non-supported address type")
	ErrInvalidSignature     = errors.New("invalid signature")
	ErrInvalidPubKey        = errors.New("invalid public key")
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

const (
	freeListMaxItems        = 125
	binaryFreeListMaxItems  = 1024
	maxWitnessItemSize      = 4_000_000
	maxWitnessItemsPerInput = 4_000_000
)

type binaryFreeList chan []byte

// Borrow returns a byte slice from the free list with a length of 8.  A new
// buffer is allocated if there are not any available on the free list.
func (l binaryFreeList) Borrow() []byte {
	var buf []byte
	select {
	case buf = <-l:
	default:
		buf = make([]byte, 8)
	}
	return buf[:8]
}

// Return puts the provided byte slice back on the free list.  The buffer MUST
// have been obtained via the Borrow function and therefore have a cap of 8.
func (l binaryFreeList) Return(buf []byte) {
	select {
	case l <- buf:
	default:
		// Let it go to the garbage collector.
	}
}

type scriptFreeList chan *scriptSlab

const scriptSlabSize = 1 << 22

type scriptSlab [scriptSlabSize]byte

var binarySerializer binaryFreeList = make(chan []byte, binaryFreeListMaxItems)
var scriptPool = make(scriptFreeList, freeListMaxItems)

// ignored and allowed to go the garbage collector.
func (c scriptFreeList) Borrow() *scriptSlab {
	var buf *scriptSlab
	select {
	case buf = <-c:
	default:
		buf = new(scriptSlab)
	}
	return buf
}

// Return puts the provided byte slice back on the free list when it has a cap
// of the expected length.  The buffer is expected to have been obtained via
// the Borrow function.  Any slices that are not of the appropriate size, such
// as those whose size is greater than the largest allowed free list item size
// are simply ignored so they can go to the garbage collector.
func (c scriptFreeList) Return(buf *scriptSlab) {
	// Return the buffer to the free list when it's not full.  Otherwise let
	// it be garbage collected.
	select {
	case c <- buf:
	default:
		// Let it go to the garbage collector.
	}
}

func BtcDecodeWitnessForBip0322(r io.Reader, pver uint32, enc wire.MessageEncoding, msg *wire.MsgTx) error {
	if enc == wire.WitnessEncoding {
		buf := binarySerializer.Borrow()
		defer binarySerializer.Return(buf)

		sbufP := scriptPool.Borrow()
		defer scriptPool.Return(sbufP)
		sbuf := sbufP[:]
		for _, ti := range msg.TxIn {
			witCount, err := wire.ReadVarInt(r, pver)
			if err != nil {
				return err
			}
			if witCount > maxWitnessItemsPerInput {
				str := fmt.Sprintf("too many witness items to fit "+
					"into max message size [count %d, max %d]",
					witCount, maxWitnessItemsPerInput)
				return fmt.Errorf("MsgTx.BtcDecode %s", str)
			}
			ti.Witness = make([][]byte, witCount)
			for j := uint64(0); j < witCount; j++ {
				ti.Witness[j], err = readScriptBuf(
					r, pver, buf, sbuf, maxWitnessItemSize,
					"script witness item",
				)
				if err != nil {
					return err
				}
				sbuf = sbuf[len(ti.Witness[j]):]
			}
		}
	} else {
		return errors.New("do not support nonwithness decode")
	}
	return nil
}

func readScriptBuf(r io.Reader, pver uint32, buf, s []byte,
	maxAllowed uint32, fieldName string) ([]byte, error) {

	count, err := ReadVarIntBuf(r, pver, buf)
	if err != nil {
		return nil, err
	}

	// Prevent byte array larger than the max message size.  It would
	// be possible to cause memory exhaustion and panics without a sane
	// upper bound on this count.
	if count > uint64(maxAllowed) {
		str := fmt.Sprintf("%s is larger than the max allowed size "+
			"[count %d, max %d]", fieldName, count, maxAllowed)
		return nil, fmt.Errorf("readScript %s", str)
	}

	_, err = io.ReadFull(r, s[:count])
	if err != nil {
		return nil, err
	}
	return s[:count], nil
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

func SignMessage(wif string, prefix, message string) (string, error) {
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, prepareMessage(prefix))
	wire.WriteVarString(&buf, 0, message)
	messageHash := chainhash.DoubleHashB(buf.Bytes())
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return "", err
	}
	sig, err := ecdsa.SignCompact(w.PrivKey, messageHash, true)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func prepareMessage(prefix string) string {
	first := "Bitcoin Signed Message:\n"
	if len(prefix) > 0 {
		first = prefix + " Signed Message:\n"
	}
	return first
}

func MPCUnsignedMessage(prefix string, message string) string {
	var buf bytes.Buffer
	err := wire.WriteVarString(&buf, 0, prepareMessage(prefix))
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

func MPCSignedMessageCompat(prefix, message string, signature string, publicKeyHex string, network *chaincfg.Params) (string, error) {
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	msgHash := MPCUnsignedMessage(prefix, message)
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

func VerifyMessage(signatureStr, prefix, message, publicKeyHex, address, signType string, network *chaincfg.Params) error {
	if signType == "bip322-simple" {
		return VerifySimpleForBip0322(message, address, signatureStr, publicKeyHex, network)
	}
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	var signature []byte
	var err error
	if IsHexString(signatureStr) {
		signature, err = hex.DecodeString(signatureStr)
	} else {
		signature, err = base64.StdEncoding.DecodeString(signatureStr)
	}
	if err != nil {
		return err
	}
	messageHash := MPCUnsignedMessage(prefix, message)

	h, err := hex.DecodeString(messageHash)
	if err != nil {
		return err
	}
	p, ok, err := ecdsa.RecoverCompact(signature, h)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("invalid signature")
	}
	if err := checkPublicKeyHex(p, publicKeyHex); err != nil {
		return err
	}
	if err := checkAddr(p, address, network); err != nil {
		return err
	}
	return nil
}

func NewPossibleAddrs(pub *btcec.PublicKey, network *chaincfg.Params) ([]string, error) {
	addrs := make([]string, 0)
	newAddr, err := btcutil.NewAddressPubKey(pub.SerializeUncompressed(), network)
	if err != nil {
		return addrs, err
	}
	addrs = append(addrs, newAddr.EncodeAddress())
	newAddr, err = btcutil.NewAddressPubKey(pub.SerializeCompressed(), network)
	if err != nil {
		return addrs, err
	}
	addrs = append(addrs, newAddr.EncodeAddress())

	pkHash := btcutil.Hash160(pub.SerializeUncompressed())
	newAddr2, err := btcutil.NewAddressPubKeyHash(pkHash, network)
	if err != nil {
		return addrs, err
	}

	addrs = append(addrs, newAddr2.EncodeAddress())
	pkHash = btcutil.Hash160(pub.SerializeCompressed())
	newAddr3, err := btcutil.NewAddressPubKeyHash(pkHash, network)
	if err != nil {
		return addrs, err
	}

	addrs = append(addrs, newAddr3.EncodeAddress())
	pkHash2 := btcutil.Hash160(pub.SerializeUncompressed())
	newAddr4, err := btcutil.NewAddressWitnessPubKeyHash(pkHash2, network)
	if err != nil {
		return addrs, err
	}

	addrs = append(addrs, newAddr4.EncodeAddress())
	pkHash = btcutil.Hash160(pub.SerializeCompressed())
	newAddr5, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, network)
	if err != nil {
		return addrs, err
	}

	addrs = append(addrs, newAddr5.EncodeAddress())

	rootKey := txscript.ComputeTaprootKeyNoScript(pub)
	newAddr6, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(rootKey), network)
	if err != nil {
		return addrs, err
	}
	addrs = append(addrs, newAddr6.EncodeAddress())
	return addrs, nil
}

func checkPublicKeyHex(pub *btcec.PublicKey, publicKeyHex string) error {
	if ph := hex.EncodeToString(pub.SerializeCompressed()); ph == publicKeyHex {
		return nil
	}
	if ph := hex.EncodeToString(pub.SerializeUncompressed()); ph == publicKeyHex {
		return nil
	}
	return ErrInvalidPubKey
}

func NewAddressPubKeyHashFromWif(wif string, network *chaincfg.Params) (string, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return "", err
	}
	return NewAddressPubKeyHash(w.PrivKey.PubKey(), network)
}

func NewAddressPubKeyHash(pub *btcec.PublicKey, network *chaincfg.Params) (string, error) {
	if pub == nil {
		return "", ErrInvalidPubKey
	}
	if network != nil && network.Net == zecNet {
		return NewZECAddr(pub.SerializeCompressed()), nil
	}
	pkHash := btcutil.Hash160(pub.SerializeUncompressed())
	newAddr, err := btcutil.NewAddressPubKeyHash(pkHash, network)
	if err != nil {
		return "", err
	}
	return newAddr.EncodeAddress(), nil
}

func checkAddr(pub *btcec.PublicKey, addr string, network *chaincfg.Params) error {
	if len(addr) == 0 {
		return ErrNonSupportedAddrType
	}
	//zec address is different from other address types.
	if network != nil && network.Net == zecNet {
		if addr == NewZECAddr(pub.SerializeCompressed()) {
			return nil
		}
		if addr == NewZECAddr(pub.SerializeUncompressed()) {
			return nil
		}
		return ErrInvalidSignature
	}
	a, err := btcutil.DecodeAddress(addr, network)
	if err != nil {
		return err
	}
	switch a.(type) {
	case *btcutil.AddressPubKey:
		newAddr, err := btcutil.NewAddressPubKey(pub.SerializeUncompressed(), network)
		if err != nil {
			return err
		}
		if newAddr.EncodeAddress() == addr {
			return nil
		}
		newAddr, err = btcutil.NewAddressPubKey(pub.SerializeCompressed(), network)
		if err != nil {
			return err
		}
		if newAddr.EncodeAddress() == addr {
			return nil
		}
		return ErrInvalidSignature

	case *btcutil.AddressPubKeyHash:
		pkHash := btcutil.Hash160(pub.SerializeUncompressed())
		newAddr, err := btcutil.NewAddressPubKeyHash(pkHash, network)
		if err != nil {
			return err
		}
		if newAddr.EncodeAddress() == addr {
			return nil
		}
		pkHash = btcutil.Hash160(pub.SerializeCompressed())
		newAddr, err = btcutil.NewAddressPubKeyHash(pkHash, network)
		if err != nil {
			return err
		}
		if newAddr.EncodeAddress() == addr {
			return nil
		}
		return ErrInvalidSignature
	case *btcutil.AddressWitnessPubKeyHash:
		pkHash := btcutil.Hash160(pub.SerializeUncompressed())
		newAddr, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, network)
		if err != nil {
			return err
		}
		if newAddr.EncodeAddress() == addr {
			return nil
		}
		pkHash = btcutil.Hash160(pub.SerializeCompressed())
		newAddr, err = btcutil.NewAddressWitnessPubKeyHash(pkHash, network)
		if err != nil {
			return err
		}
		if newAddr.EncodeAddress() == addr {
			return nil
		}
		return ErrInvalidSignature
	case *btcutil.AddressTaproot:
		rootKey := txscript.ComputeTaprootKeyNoScript(pub)
		newAddr, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(rootKey), network)
		if err != nil {
			return err
		}
		if newAddr.EncodeAddress() == addr {
			return nil
		}
		return ErrInvalidSignature
	default:
		return ErrNonSupportedAddrType
	}
}

func Wif2PubKeyHex(wif string) (string, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(w.PrivKey.PubKey().SerializeCompressed()), nil
}

func VerifySimpleForBip0322(message, address, signature, publicKey string, network *chaincfg.Params) error {
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	txId, err := BuildToSpend(message, address, network)
	if err != nil {
		return err
	}

	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return err
	}

	prevOut := wire.NewOutPoint(txHash, 0)
	inputs := []*wire.OutPoint{prevOut}
	script, _ := hex.DecodeString("6a")
	outputs := []*wire.TxOut{wire.NewTxOut(0, script)}
	nSequences := []uint32{uint32(0)}
	p, err := NewPsbt(inputs, outputs, int32(0), uint32(0), nSequences, "bip0322-simple")
	if err != nil {
		return err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return err
	}

	dummyPkScript, _ := AddrToPkScript(address, network)
	dummyWitnessUtxo := wire.NewTxOut(0, dummyPkScript)

	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 0)
	if err != nil {
		return err
	}

	prevPkScript, err := AddrToPkScript(address, network)
	witnessUtxo := wire.NewTxOut(0, prevPkScript)
	prevOuts := map[wire.OutPoint]*wire.TxOut{
		*prevOut: witnessUtxo,
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)
	publicKeyBs, err := hex.DecodeString(publicKey)
	if err != nil {
		return err
	}
	pub, err := btcec.ParsePubKey(publicKeyBs)
	if err != nil {
		return err
	}
	tx := updater.Upsbt.UnsignedTx
	sigB64, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	r := bytes.NewReader(sigB64)
	sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutputFetcher)
	if txscript.IsPayToTaproot(prevPkScript) {
		internalPubKey := schnorr.SerializePubKey(pub)
		updater.Upsbt.Inputs[0].TaprootInternalKey = internalPubKey

		sigHash, err := txscript.CalcTaprootSignatureHash(sigHashes, txscript.SigHashDefault, tx, 0, prevOutputFetcher)
		if err != nil {
			return err
		}
		err = BtcDecodeWitnessForBip0322(r, 0, wire.WitnessEncoding, tx)
		if err != nil {
			return err
		}
		sig, err := schnorr.ParseSignature(tx.TxIn[0].Witness[0])
		if err != nil {
			return err
		}
		tweakedPublicKey := txscript.ComputeTaprootKeyNoScript(pub)
		if sig.Verify(sigHash, tweakedPublicKey) {
			return nil
		}
		return ErrInvalidSignature
	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		// todo
		return nil
	} else {
		script, err := PayToPubKeyHashScript(btcutil.Hash160(publicKeyBs))
		if err != nil {
			return err
		}
		sigHash, err := txscript.CalcWitnessSigHash(script, sigHashes, txscript.SigHashAll, tx, 0, 0)
		if err != nil {
			return err
		}
		err = BtcDecodeWitnessForBip0322(r, 0, wire.WitnessEncoding, tx)
		if err != nil {
			return err
		}
		sig, err := ecdsa.ParseSignature(tx.TxIn[0].Witness[0])
		if err != nil {
			return err
		}
		pubInWitnessStack := tx.TxIn[0].Witness[1]
		if !reflect.DeepEqual(publicKeyBs, pubInWitnessStack) {
			return fmt.Errorf("pubInWitnessStack is wrong %s", publicKey)
		}
		if sig.Verify(sigHash, pub) {
			return nil
		}
		return ErrInvalidSignature
	}
}
