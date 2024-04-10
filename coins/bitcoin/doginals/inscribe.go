package doginals

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"strconv"
)

type InscriptionData struct {
	ContentType string `json:"contentType"`
	Body        []byte `json:"body"`
	RevealAddr  string `json:"revealAddr"`
}

type PrevOutput struct {
	TxId       string `json:"txId"`
	VOut       uint32 `json:"vOut"`
	Amount     int64  `json:"amount"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

const errInsufficientBalance = "Insufficient Balance"

type InscriptionRequest struct {
	CommitTxPrevOutputList []*PrevOutput    `json:"commitTxPrevOutputList"`
	CommitFeeRate          int64            `json:"commitFeeRate"`
	RevealFeeRate          int64            `json:"revealFeeRate"`
	RevealOutValue         int64            `json:"revealOutValue"`
	InscriptionData        *InscriptionData `json:"inscriptionData"`
	Address                string           `json:"address"`
	DustSize               int64            `json:"dustSize"`
}

type Inscription struct {
	P   string `json:"p"`
	Op  string `json:"op"`
	Amt string `json:"amt"`
}

func (i *Inscription) Amount() (int64, error) {
	return strconv.ParseInt(i.Amt, 10, 64)
}

type inscriptionTxCtxData struct {
	PrivateKey              *btcec.PrivateKey
	InscriptionScript       []byte
	RedeemScript            []byte
	CommitTxAddress         btcutil.Address
	CommitTxAddressPkScript []byte
	Hash                    []byte
	RevealPkScript          []byte
	RevealTxPrevOutput      *wire.TxOut
}

type InscriptionTool struct {
	CommitTxPrevOutputFetcher *txscript.MultiPrevOutFetcher
	CommitTxPrivateKeyList    []*btcec.PrivateKey
	InscriptionTxCtxData      []*inscriptionTxCtxData
	RevealTxPrevOutputFetcher *txscript.MultiPrevOutFetcher
	CommitTxPrevOutputList    []*PrevOutput
	RevealTxs                 []*wire.MsgTx
	CommitTx                  *wire.MsgTx
	MustCommitTxFee           int64
	MustRevealTxFees          []int64
	CommitAddrs               []string
	FromAddr                  btcutil.Address
	RevealAddr                btcutil.Address
}

type InscribeTxs struct {
	CommitTx     string   `json:"commitTx"`
	RevealTxs    []string `json:"revealTxs"`
	CommitTxFee  int64    `json:"commitTxFee"`
	RevealTxFees []int64  `json:"revealTxFees"`
	CommitAddrs  []string `json:"commitAddrs"`
}

const (
	DefaultTxVersion      = 2
	DefaultSequenceNum    = 0xfffffffd
	DefaultRevealOutValue = int64(100000)
	DefaultMinChangeValue = int64(100000)

	WitnessScaleFactor = 4

	ChangeOutputMaxSize = int64(20 + 4 + 34 + 4)
	MaxChunkLen         = 240
	MaxPayloadLen       = 1500
)

type Chunk struct {
	Buf       []byte
	Len       int
	OpcodeNum int
}

func bufferToBuffer(b []byte) []byte {
	c := bufferToChunk(b)
	buf := make([]byte, 0)
	buf = append(buf, byte(c.OpcodeNum))
	if len(c.Buf) > 0 {
		if c.OpcodeNum < txscript.OP_PUSHDATA1 {
		} else if c.OpcodeNum == txscript.OP_PUSHDATA1 {
			buf = append(buf, byte(c.Len))
		} else if c.OpcodeNum == txscript.OP_PUSHDATA2 {
			buf = binary.LittleEndian.AppendUint64(buf, uint64(c.Len))
			//bw.writeUInt64(c.len)
		} else if c.OpcodeNum == txscript.OP_PUSHDATA4 {
			buf = binary.LittleEndian.AppendUint32(buf, uint32(c.Len))
			//bw.writeUInt32(c.len)
		}
		buf = append(buf, c.Buf...)
	}
	return buf
}

type DogScript struct {
	chunks []*Chunk
}

func (d *DogScript) push(c *Chunk) {
	if c == nil {
		return
	}
	d.chunks = append(d.chunks, c)
}

func (d *DogScript) total() int {
	if len(d.chunks) == 0 {
		return 0
	}
	size := 0
	for _, chunk := range d.chunks {
		size += wire.VarIntSerializeSize(uint64(chunk.OpcodeNum))
		var opcodenum = chunk.OpcodeNum
		if len(chunk.Buf) > 0 {
			if opcodenum < txscript.OP_PUSHDATA1 {
			} else if opcodenum == txscript.OP_PUSHDATA1 {
				size += wire.VarIntSerializeSize(uint64(chunk.Len))
			} else if opcodenum == txscript.OP_PUSHDATA2 {
				size += wire.VarIntSerializeSize(uint64(chunk.Len))
			} else if opcodenum == txscript.OP_PUSHDATA4 {
				size += wire.VarIntSerializeSize(uint64(chunk.Len))
			}
			size += len(chunk.Buf)
		}
	}
	return size
}

func (d *DogScript) toBuffer() []byte {
	buf := make([]byte, 0, d.total())
	for _, chunk := range d.chunks {
		var opcodenum = chunk.OpcodeNum
		buf = append(buf, byte(chunk.OpcodeNum))
		if len(chunk.Buf) > 0 {
			if opcodenum < txscript.OP_PUSHDATA1 {
			} else if opcodenum == txscript.OP_PUSHDATA1 {
				buf = append(buf, byte(chunk.Len))
			} else if opcodenum == txscript.OP_PUSHDATA2 {
				buf = binary.LittleEndian.AppendUint64(buf, uint64(chunk.Len))
			} else if opcodenum == txscript.OP_PUSHDATA4 {
				buf = binary.LittleEndian.AppendUint32(buf, uint32(chunk.Len))
			}
			buf = append(buf, chunk.Buf...)
		}
	}
	return buf
}

func numberToChunk(n int) *Chunk {
	var buf []byte
	var length int
	var op int
	if n <= 16 {
		length, op = 0, 80+n
		if n == 0 {
			op = 0
		}
	} else if n < 128 {
		buf, length, op = []byte{byte(n)}, 1, 1
	} else {
		buf, length, op = []byte{byte(n % 256), byte(n / 256)}, 2, 2
	}
	return &Chunk{
		// @ts-ignore
		Buf:       buf,    //n <= 16 ? undefined : n < 128 ? Buffer.from([n]) : Buffer.from([n % 256, n / 256]),
		Len:       length, //n <= 16 ? 0 : n < 128 ? 1 : 2,
		OpcodeNum: op,     //n == 0 ? 0 : n <= 16 ? 80 + n : n < 128 ? 1 : 2
	}
}

func bufferToChunk(b []byte) *Chunk {
	var buf []byte
	var op int
	if len(b) > 0 {
		buf = b
	}
	if len(b) <= 75 {
		op = len(b)
	} else if len(b) <= 255 {
		op = 76
	} else {
		op = 77
	}
	return &Chunk{
		Buf:       buf,
		Len:       len(b),
		OpcodeNum: op, //b.length <= 75 ? b.length : b.length <= 255 ? 76 : 77
	}
}

func opcodeToChunk(op int) *Chunk {
	return &Chunk{OpcodeNum: op}
}

func NewInscriptionTool(request *InscriptionRequest) (*InscriptionTool, error) {
	var commitTxPrivateKeyList []*btcec.PrivateKey
	for _, prevOutput := range request.CommitTxPrevOutputList {
		privateKeyWif, err := btcutil.DecodeWIF(prevOutput.PrivateKey)
		if err != nil {
			return nil, err
		}
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, privateKeyWif.PrivKey)
	}
	tool := &InscriptionTool{
		CommitTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		CommitTxPrivateKeyList:    commitTxPrivateKeyList,
		RevealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		CommitTxPrevOutputList:    request.CommitTxPrevOutputList,
	}
	return tool, tool._initTool(request)
}

func (tool *InscriptionTool) _initTool(request *InscriptionRequest) error {
	minChangeValue := DefaultMinChangeValue
	if request.DustSize > 0 {
		minChangeValue = request.DustSize
	}
	fromAddr, inscriptionTxCtxData, err := createInscriptionTxCtxData(request)
	if err != nil {
		return err
	}
	revealAddr, err := btcutil.DecodeAddress(request.InscriptionData.RevealAddr, &DogeMainNetParams)
	if err != nil {
		return err
	}
	tool.FromAddr, tool.InscriptionTxCtxData, tool.RevealAddr = fromAddr, inscriptionTxCtxData, revealAddr
	totalRevealPrevOutputValue, err := tool.buildEmptyRevealTx(DefaultRevealOutValue, request.RevealFeeRate)
	if err != nil {
		return err
	}
	err = tool.buildCommitTx(request.CommitTxPrevOutputList, request.Address, totalRevealPrevOutputValue, request.CommitFeeRate, DefaultRevealOutValue, minChangeValue)
	if err != nil {
		return err
	}
	err = tool.signCommitTx()
	if err != nil {
		return errors.New("sign commit tx error")
	}
	err = tool.completeRevealTx()
	if err != nil {
		return err
	}
	return err
}

func createInscriptionTxCtxData(inscriptionRequest *InscriptionRequest) (btcutil.Address, []*inscriptionTxCtxData, error) {
	privateKeyWif, err := btcutil.DecodeWIF(inscriptionRequest.CommitTxPrevOutputList[0].PrivateKey)
	if err != nil {
		return nil, nil, err
	}
	privateKey := privateKeyWif.PrivKey
	pubKey := privateKeyWif.PrivKey.PubKey().SerializeCompressed()
	pkHash := btcutil.Hash160(pubKey)
	addr, err := btcutil.NewAddressPubKeyHash(pkHash, &DogeMainNetParams)
	if err != nil {
		return nil, nil, err
	}
	parts, data := make([][]byte, 0), inscriptionRequest.InscriptionData.Body
	for i := 0; i < len(data); {
		end := i + MaxChunkLen
		if end > len(data) {
			end = len(data)
		}
		parts = append(parts, data[i:end])
		i = end
	}
	inscription := DogScript{}
	inscription.push(bufferToChunk([]byte("ord")))
	inscription.push(numberToChunk(len(parts)))
	inscription.push(bufferToChunk([]byte(inscriptionRequest.InscriptionData.ContentType)))
	for n, part := range parts {
		inscription.push(numberToChunk(len(parts) - n - 1))
		inscription.push(bufferToChunk(part))
	}
	ctxDatas := make([]*inscriptionTxCtxData, 0)
	i := 0
	for i < len(inscription.chunks) {
		partial := DogScript{}
		if len(ctxDatas) == 0 {
			// @ts-ignore
			partial.push(inscription.chunks[i])
			i++
		}
		for partial.total() <= MaxPayloadLen && len(inscription.chunks) > i {
			partial.push(inscription.chunks[i])
			i++
			partial.push(inscription.chunks[i])
			i++
		}

		if partial.total() > MaxPayloadLen {
			// @ts-ignore
			partial.chunks = partial.chunks[0 : len(partial.chunks)-2]
			i -= 2
		}
		lock := DogScript{}
		lock.push(bufferToChunk(pubKey))
		lock.push(opcodeToChunk(txscript.OP_CHECKSIGVERIFY))
		for i := 0; i < len(partial.chunks); i++ {
			lock.push(opcodeToChunk(txscript.OP_DROP))
		}
		lock.push(opcodeToChunk(txscript.OP_TRUE))
		address, err := btcutil.NewAddressScriptHash(lock.toBuffer(), &DogeMainNetParams)
		if err != nil {
			return nil, nil, err
		}
		script, err := txscript.PayToAddrScript(address)
		if err != nil {
			return nil, nil, err
		}
		revealAddr, err := btcutil.DecodeAddress(inscriptionRequest.InscriptionData.RevealAddr, &DogeMainNetParams)
		if err != nil {
			return nil, nil, err
		}
		revealPkScript, err := txscript.PayToAddrScript(revealAddr)
		if err != nil {
			return nil, nil, err
		}

		ctx := &inscriptionTxCtxData{
			PrivateKey:              privateKey,
			InscriptionScript:       partial.toBuffer(),
			RedeemScript:            lock.toBuffer(),
			CommitTxAddress:         address,
			CommitTxAddressPkScript: script,
			RevealTxPrevOutput: &wire.TxOut{
				PkScript: nil,
				Value:    100000,
			},
			RevealPkScript: revealPkScript,
		}
		ctxDatas = append(ctxDatas, ctx)
	}
	return addr, ctxDatas, nil
}

func AddrToPkScript(addr string) ([]byte, error) {
	address, err := btcutil.DecodeAddress(addr, &DogeMainNetParams)
	if err != nil {
		return nil, err
	}

	return txscript.PayToAddrScript(address)
}

func (tool *InscriptionTool) buildEmptyRevealTx(revealOutValue, revealFeeRate int64) (int64, error) {
	totalPrevOutputValue := int64(0)
	total := len(tool.InscriptionTxCtxData)
	revealTxs := make([]*wire.MsgTx, total)
	mustRevealTxFees := make([]int64, total)
	commitAddrs := make([]string, total)
	left := int64(0)
	prevOutputValue := int64(0)
	for i := len(tool.InscriptionTxCtxData) - 1; i > -1; i-- {
		ctx := tool.InscriptionTxCtxData[i]
		tx := wire.NewMsgTx(DefaultTxVersion)
		in := wire.NewTxIn(&wire.OutPoint{Index: uint32(0)}, nil, nil)
		in.Sequence = DefaultSequenceNum
		tx.AddTxIn(in)
		in1 := wire.NewTxIn(&wire.OutPoint{Index: uint32(1)}, nil, nil)
		in1.Sequence = DefaultSequenceNum
		tx.AddTxIn(in1)

		scriptPubKey := ctx.CommitTxAddressPkScript
		if i == len(tool.InscriptionTxCtxData)-1 {
			scriptPubKey = ctx.RevealPkScript
		}
		out := wire.NewTxOut(revealOutValue, scriptPubKey)
		tx.AddTxOut(out)

		emptySignature := bufferToBuffer(make([]byte, 72))
		redeemScript := bufferToBuffer(ctx.RedeemScript)
		scrip0 := make([]byte, 0, len(ctx.InscriptionScript)+len(emptySignature)+len(redeemScript))
		scrip0 = append(scrip0, ctx.InscriptionScript...)
		scrip0 = append(scrip0, emptySignature...)
		scrip0 = append(scrip0, redeemScript...)
		tx.TxIn[0].SignatureScript = scrip0
		if i != len(tool.InscriptionTxCtxData)-1 {
			script1, err := txscript.PayToAddrScript(tool.FromAddr)
			if err != nil {
				return 0, err
			}
			tx.AddTxOut(wire.NewTxOut(left, script1))
		}

		tx.TxIn[1].SignatureScript = make([]byte, 106)
		fee := (int64(DogeByteLength(tx)) + ChangeOutputMaxSize) * revealFeeRate
		left += fee
		tool.InscriptionTxCtxData[i].RevealTxPrevOutput = &wire.TxOut{
			PkScript: tool.InscriptionTxCtxData[i].CommitTxAddressPkScript,
			Value:    fee,
		}
		totalPrevOutputValue += fee
		totalPrevOutputValue += prevOutputValue
		revealTxs[i] = tx
		mustRevealTxFees[i] = fee
		commitAddrs[i] = tool.InscriptionTxCtxData[i].CommitTxAddress.EncodeAddress()
	}
	tool.RevealTxs = revealTxs
	tool.MustRevealTxFees = mustRevealTxFees
	tool.CommitAddrs = commitAddrs
	totalPrevOutputValue += revealOutValue
	return totalPrevOutputValue, nil
}

func (tool *InscriptionTool) buildCommitTx(commitTxPrevOutputList []*PrevOutput, changeAddress string, totalRevealPrevOutputValue, commitFeeRate int64, revealOutValue int64, minChangeValue int64) error {
	totalSenderAmount := btcutil.Amount(0)
	tx := wire.NewMsgTx(DefaultTxVersion)
	for _, prevOutput := range commitTxPrevOutputList {
		txHash, err := chainhash.NewHashFromStr(prevOutput.TxId)
		if err != nil {
			return err
		}
		outPoint := wire.NewOutPoint(txHash, prevOutput.VOut)
		pkScript, err := AddrToPkScript(prevOutput.Address)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(prevOutput.Amount, pkScript)
		tool.CommitTxPrevOutputFetcher.AddPrevOut(*outPoint, txOut)

		in := wire.NewTxIn(outPoint, nil, nil)
		in.Sequence = DefaultSequenceNum
		tx.AddTxIn(in)

		totalSenderAmount += btcutil.Amount(prevOutput.Amount)
	}
	tx.AddTxOut(wire.NewTxOut(revealOutValue, tool.InscriptionTxCtxData[0].CommitTxAddressPkScript))
	pkScript, err := txscript.PayToAddrScript(tool.FromAddr)
	if err != nil {
		return err
	}
	tx.AddTxOut(wire.NewTxOut(totalRevealPrevOutputValue, pkScript))
	if len(changeAddress) > 0 {
		changePkScript, err := AddrToPkScript(changeAddress)
		if err != nil {
			return err
		}
		tx.AddTxOut(wire.NewTxOut(0, changePkScript))
	}

	txForEstimate := wire.NewMsgTx(DefaultTxVersion)
	txForEstimate.TxIn = tx.TxIn
	txForEstimate.TxOut = tx.TxOut
	if err := Sign(txForEstimate, tool.CommitTxPrivateKeyList, tool.CommitTxPrevOutputFetcher); err != nil {
		return err
	}

	fee := btcutil.Amount(DogeByteLength(txForEstimate)+ChangeOutputMaxSize) * btcutil.Amount(commitFeeRate)
	if int64(totalSenderAmount) < totalRevealPrevOutputValue+int64(fee) {
		return errors.New("insufficient amount")
	}
	changeAmount := totalSenderAmount - btcutil.Amount(totalRevealPrevOutputValue) - fee
	if int64(changeAmount) >= minChangeValue {
		tx.TxOut[len(tx.TxOut)-1].Value = int64(changeAmount)
	} else {
		tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
		txForEstimate.TxOut = txForEstimate.TxOut[:len(txForEstimate.TxOut)-1]
		feeWithoutChange := btcutil.Amount(DogeByteLength(txForEstimate)+ChangeOutputMaxSize) * btcutil.Amount(commitFeeRate)
		if totalSenderAmount < btcutil.Amount(totalRevealPrevOutputValue)+feeWithoutChange {
			tool.MustCommitTxFee = int64(fee)
			return errors.New(errInsufficientBalance)
		}
	}
	tool.CommitTx = tx
	return nil
}

func (tool *InscriptionTool) completeRevealTx() error {
	for i, r := range tool.RevealTxs {
		if i == 0 {
			tool.RevealTxPrevOutputFetcher.AddPrevOut(wire.OutPoint{
				Hash:  tool.CommitTx.TxHash(),
				Index: uint32(1),
			}, tool.InscriptionTxCtxData[i].RevealTxPrevOutput)
			r.TxIn[0].PreviousOutPoint.Hash = tool.CommitTx.TxHash()
			r.TxIn[1].PreviousOutPoint.Hash = tool.CommitTx.TxHash()
		} else {
			tool.RevealTxPrevOutputFetcher.AddPrevOut(wire.OutPoint{
				Hash:  tool.RevealTxs[i-1].TxHash(),
				Index: uint32(1),
			}, tool.InscriptionTxCtxData[i].RevealTxPrevOutput)
			r.TxIn[0].PreviousOutPoint.Hash = tool.RevealTxs[i-1].TxHash()
			r.TxIn[1].PreviousOutPoint.Hash = tool.RevealTxs[i-1].TxHash()
		}

		//tool.revealTxPrevOutputFetcher.push(this.inscriptionTxCtxDataList[i].revealTxPrevOutput!.value);

		sig, err := txscript.RawTxInSignature(r, 0, tool.InscriptionTxCtxData[i].RedeemScript, txscript.SigHashAll, tool.InscriptionTxCtxData[i].PrivateKey)
		if err != nil {
			return err
		}
		sigScript := bufferToBuffer(sig)
		redeemScript := bufferToBuffer(tool.InscriptionTxCtxData[i].RedeemScript)
		signatureScript := make([]byte, 0, len(tool.InscriptionTxCtxData[i].InscriptionScript)+len(sigScript)+len(redeemScript))
		signatureScript = append(signatureScript, tool.InscriptionTxCtxData[i].InscriptionScript...)
		signatureScript = append(signatureScript, sigScript...)
		signatureScript = append(signatureScript, redeemScript...)

		if err != nil {
			return err
		}
		tool.RevealTxs[i].TxIn[0].SignatureScript = signatureScript
		prevScript2, err := txscript.PayToAddrScript(tool.FromAddr)
		if err != nil {
			return err
		}
		sigScript1, err := txscript.SignatureScript(r, 1, prevScript2, txscript.SigHashAll, tool.InscriptionTxCtxData[i].PrivateKey, true)
		if err != nil {
			return err
		}
		tool.RevealTxs[i].TxIn[1].SignatureScript = sigScript1
	}
	return nil
}

func (tool *InscriptionTool) signCommitTx() error {
	return Sign(tool.CommitTx, tool.CommitTxPrivateKeyList, tool.CommitTxPrevOutputFetcher)
}

func Sign(tx *wire.MsgTx, privateKeys []*btcec.PrivateKey, prevOutFetcher *txscript.MultiPrevOutFetcher) error {
	for i, in := range tx.TxIn {
		prevOut := prevOutFetcher.FetchPrevOutput(in.PreviousOutPoint)
		privKey := privateKeys[i]
		if !(txscript.IsPayToScriptHash(prevOut.PkScript) || txscript.IsPayToPubKeyHash(prevOut.PkScript) || txscript.IsPayToPubKey(prevOut.PkScript)) {
			return errors.New("non-supported address type")
		}
		sigScript, err := txscript.SignatureScript(tx, i, prevOut.PkScript, txscript.SigHashAll, privKey, true)
		if err != nil {
			return err
		}
		in.SignatureScript = sigScript
	}

	return nil
}

func GetTxHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (tool *InscriptionTool) GetCommitTxHex() (string, error) {
	return GetTxHex(tool.CommitTx)
}

func (tool *InscriptionTool) GetRevealTxHexList() ([]string, error) {
	txHexList := make([]string, len(tool.RevealTxs))
	for i := range tool.RevealTxs {
		txHex, err := GetTxHex(tool.RevealTxs[i])
		if err != nil {
			return nil, err
		}
		txHexList[i] = txHex
	}
	return txHexList, nil
}

func (tool *InscriptionTool) CalculateFee() (int64, []int64) {
	commitTxFee := int64(0)
	for _, in := range tool.CommitTx.TxIn {
		commitTxFee += tool.CommitTxPrevOutputFetcher.FetchPrevOutput(in.PreviousOutPoint).Value
	}
	for _, out := range tool.CommitTx.TxOut {
		commitTxFee -= out.Value
	}
	revealTxFees := make([]int64, 0)
	for _, tx := range tool.RevealTxs {
		revealTxFee := tool.RevealTxPrevOutputFetcher.FetchPrevOutput(tx.TxIn[1].PreviousOutPoint).Value
		revealTxFees = append(revealTxFees, revealTxFee)
	}
	return commitTxFee, revealTxFees
}

func Inscribe(request *InscriptionRequest) (*InscribeTxs, error) {
	tool, err := NewInscriptionTool(request)
	if err != nil && err.Error() == errInsufficientBalance {
		return &InscribeTxs{
			CommitTx:     "",
			RevealTxs:    []string{},
			CommitTxFee:  tool.MustCommitTxFee,
			RevealTxFees: tool.MustRevealTxFees,
			CommitAddrs:  tool.CommitAddrs,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	commitTx, err := tool.GetCommitTxHex()
	if err != nil {
		return nil, err
	}
	revealTxs, err := tool.GetRevealTxHexList()
	if err != nil {
		return nil, err
	}

	commitTxFee, revealTxFees := tool.CalculateFee()

	return &InscribeTxs{
		CommitTx:     commitTx,
		RevealTxs:    revealTxs,
		CommitTxFee:  commitTxFee,
		RevealTxFees: revealTxFees,
		CommitAddrs:  tool.CommitAddrs,
	}, nil
}

// GetTransactionWeight computes the value of the weight metric for a given
// transaction. Currently the weight metric is simply the sum of the
// transactions's serialized size without any witness data scaled
// proportionally by the WitnessScaleFactor, and the transaction's serialized
// size including any witness data.
func GetTransactionWeight(tx *btcutil.Tx) int64 {
	msgTx := tx.MsgTx()

	baseSize := msgTx.SerializeSizeStripped()
	totalSize := msgTx.SerializeSize()

	// (baseSize * 3) + totalSize
	return int64((baseSize * (WitnessScaleFactor - 1)) + totalSize)
}

// GetTxVirtualSize computes the virtual size of a given transaction. A
// transaction's virtual size is based off its weight, creating a discount for
// any witness data it contains, proportional to the current
// blockchain.WitnessScaleFactor value.
func GetTxVirtualSize(tx *btcutil.Tx) int64 {
	if tx == nil {
		return 0
	}
	return DogeByteLength(tx.MsgTx())
}

func DogeByteLength(tx *wire.MsgTx) int64 {
	if tx == nil {
		return 0
	}
	result := 4 + 9 + 9 + 4
	for _, in := range tx.TxIn {
		l := 32 + 4 + 4 + wire.VarIntSerializeSize(uint64(len(in.SignatureScript))) + len(in.SignatureScript)
		result = result + 32 + 4 + l
	}
	for _, out := range tx.TxOut {
		result = result + 9 + len(out.PkScript)
	}
	return int64(result)
}
