package ton

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"

	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"github.com/okx/go-wallet-sdk/coins/ton/tlb"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/jetton"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
)

type TransferParams struct {
	Seed     []byte
	PubKey   []byte
	Seqno    uint32
	ExpireAt int64
	Simulate bool
	To       string
	Amount   string
	Comment  string
	Mode     uint8

	// Jetton only parameters
	From                  string
	Decimals              int
	MessageAttachedTons   string
	InvokeNotificationFee string
	CustomPayload         string
	StateInit             string
	Rnd                   uint64

	IsToken bool
	Version wallet.Version
}

type TonTransferBuilder struct {
	wallet          *wallet.Wallet
	params          *TransferParams
	payload         *cell.Builder
	messages        []*wallet.Message
	externalMessage *tlb.ExternalMessage
}

func NewTonTransferBuilder(params *TransferParams) (*TonTransferBuilder, error) {
	b := &TonTransferBuilder{
		params:   params,
		messages: make([]*wallet.Message, 0),
	}
	w, err := NewWallet(b.params.Seed, b.params.PubKey, b.params.Version)
	if err != nil {
		return nil, err
	}
	spec := w.GetSpec().(wallet.SpecRegularSetter)
	spec.SetCustomSeqnoFetcher(func() uint32 {
		return b.params.Seqno
	})
	if b.params.ExpireAt < 0 {
		return nil, errors.New("invalid expiration")
	}
	spec.SetExpireAt(b.params.ExpireAt)
	b.wallet = w
	return b, nil
}

func (b *TonTransferBuilder) BuildSingleMessage() (*wallet.Message, error) {
	if b.params.IsToken {
		return b.BuildJettonMessage()
	}
	return b.BuildTonMessage()
}

func (b *TonTransferBuilder) BuildTonMessage() (*wallet.Message, error) {
	toAddr, err := address.ParseAddr(b.params.To)
	if err != nil {
		return nil, err
	}
	toAmount, ok := new(big.Int).SetString(b.params.Amount, 10)
	if !ok {
		return nil, err
	}
	msg, err := b.wallet.BuildTransfer(toAddr, tlb.FromNanoTON(toAmount), false, b.params.Comment, b.params.Mode)
	if err != nil {
		return nil, err
	}
	b.messages = append(b.messages, msg)
	return msg, nil
}

func (b *TonTransferBuilder) BuildJettonMessage() (*wallet.Message, error) {
	fromAddr, err := address.ParseAddr(b.params.From)
	if err != nil {
		return nil, err
	}
	toAddr, err := address.ParseAddr(b.params.To)
	if err != nil {
		return nil, err
	}
	responseToAddress := b.wallet.Address()
	toAmount, ok := new(big.Int).SetString(b.params.Amount, 10)
	if !ok {
		return nil, err
	}
	amountForwardTON := tlb.ZeroCoins
	if b.params.InvokeNotificationFee != "" {
		invokenFee, ok := new(big.Int).SetString(b.params.InvokeNotificationFee, 10)
		if !ok {
			return nil, err
		}
		amountForwardTON = tlb.FromNanoTON(invokenFee)
	}
	var payloadForward *cell.Cell
	if b.params.Comment != "" {
		payloadForward = cell.BeginCell().MustStoreUInt(0, 32).MustStoreStringSnake(b.params.Comment).EndCell()
	}
	jw := &jetton.WalletClient{}
	toAmountCoins, err := tlb.FromNano(toAmount, b.params.Decimals)
	if err != nil {
		return nil, err
	}
	var customPayloadCell *cell.Cell
	if len(b.params.CustomPayload) > 0 {
		customPayloadCell, err = wallet.TryParseCell(b.params.CustomPayload)
		if err != nil {
			return nil, err
		}
	}
	transferPayload, err := jw.BuildTransferPayloadV2(toAddr, responseToAddress, toAmountCoins, amountForwardTON, payloadForward, customPayloadCell, b.params.Rnd)
	if err != nil {
		return nil, err
	}

	messageAttachedVal := "50000000"
	if b.params.MessageAttachedTons != "" {
		messageAttachedVal = b.params.MessageAttachedTons
	}
	messageAttachedValBig, ok := new(big.Int).SetString(messageAttachedVal, 10)
	if !ok {
		return nil, err
	}
	message := wallet.SimpleMessage(fromAddr, tlb.FromNanoTON(messageAttachedValBig), transferPayload)

	if len(b.params.StateInit) > 0 {
		b, err := wallet.TryParseBase64(b.params.StateInit)
		if err != nil {
			return nil, err
		}
		bd2, err := cell.FromBOC(b)
		if err != nil {
			return nil, err
		}
		r1, err := bd2.PeekRef(0)
		if err != nil {
			return nil, err
		}
		r2, err := bd2.PeekRef(1)
		if err != nil {
			return nil, err
		}
		in := &tlb.StateInit{
			Code: r1,
			Data: r2,
		}
		message.InternalMessage.StateInit = in
	}
	b.messages = append(b.messages, message)
	return message, nil
}

func (b *TonTransferBuilder) BuildTransferSigningHash() ([]byte, error) {
	if len(b.messages) == 0 {
		return nil, errors.New("no messages")
	}
	payload, err := b.wallet.BuildMessageUnsigned(context.Background(), b.messages, b.params.Seqno > 0)
	if err != nil {
		return nil, err
	}
	b.payload = payload
	return payload.EndCell().Hash(), nil
}

func (b *TonTransferBuilder) BuildTransferWithSignature(signature []byte) (*tlb.ExternalMessage, error) {
	if b.payload == nil {
		return nil, errors.New("payload is not set")
	}
	msg, err := b.wallet.BuildMessageWithSignature(context.Background(), b.payload, signature, b.params.Seqno > 0)
	if err != nil {
		return nil, err
	}
	b.externalMessage = msg
	return msg, nil
}

func (b *TonTransferBuilder) BuildTransferDirect() (*tlb.ExternalMessage, error) {
	if len(b.messages) == 0 {
		return nil, errors.New("no messages")
	}
	initialized := b.params.Seqno > 0
	msg, err := b.wallet.BuildExternalMessageForMany(context.Background(), b.messages, initialized)
	if err != nil {
		return nil, err
	}
	b.externalMessage = msg
	return msg, nil
}

func (b *TonTransferBuilder) BuildSignedTx(useBOCWithFlags bool) (*SignedTx, error) {
	if b.externalMessage == nil {
		return nil, errors.New("external message is not set")
	}

	cell := b.externalMessage.Body

	if !b.params.Simulate {
		emCell, err := tlb.ToCell(b.externalMessage)
		if err != nil {
			return nil, err
		}
		cell = emCell
	}

	signedTx := NewSignedTx(b.wallet.WalletAddress().String())

	err := signedTx.FillInit(b.wallet, b.params.Seqno)
	if err != nil {
		return nil, err
	}
	err = signedTx.FillTx(base64.StdEncoding.EncodeToString(cell.ToBOCWithFlags(useBOCWithFlags)), !b.params.Simulate)
	if err != nil {
		return nil, err
	}
	return signedTx, err
}

// For backward compatibility
func Transfer(seed, pubKey []byte, to, amount, comment string, seqno uint32, expireAt int64, mode uint8, simulate bool, version wallet.Version) (*SignedTx, error) {
	params := &TransferParams{
		Seed:     seed,
		PubKey:   pubKey,
		Seqno:    seqno,
		ExpireAt: expireAt,
		Simulate: simulate,
		Version:  version,
		To:       to,
		Amount:   amount,
		Comment:  comment,
		Mode:     mode,
		IsToken:  false,
	}
	b, err := NewTonTransferBuilder(params)
	if err != nil {
		return nil, err
	}
	_, err = b.BuildTonMessage()
	if err != nil {
		return nil, err
	}
	_, err = b.BuildTransferDirect()
	if err != nil {
		return nil, err
	}
	return b.BuildSignedTx(true)
}

func TransferJetton(seed, pubKey []byte, from, to, amount string, decimals int, seqno uint32, messageAttachedTons string, invokeNotificationFee string, customPayload, stateInit, comment string, expireAt int64, rnd uint64, simulate bool, version wallet.Version) (*SignedTx, error) {
	params := &TransferParams{
		Seed:                  seed,
		PubKey:                pubKey,
		Seqno:                 seqno,
		ExpireAt:              expireAt,
		Simulate:              simulate,
		Version:               version,
		From:                  from,
		To:                    to,
		Amount:                amount,
		Decimals:              decimals,
		MessageAttachedTons:   messageAttachedTons,
		InvokeNotificationFee: invokeNotificationFee,
		CustomPayload:         customPayload,
		StateInit:             stateInit,
		Comment:               comment,
		Rnd:                   rnd,
		IsToken:               true,
	}
	b, err := NewTonTransferBuilder(params)
	if err != nil {
		return nil, err
	}
	_, err = b.BuildJettonMessage()
	if err != nil {
		return nil, err
	}
	_, err = b.BuildTransferDirect()
	if err != nil {
		return nil, err
	}
	return b.BuildSignedTx(false)
}

func getExternalMessageCell(boc string) (*cell.Cell, error) {
	emCellBytes, err := base64.StdEncoding.DecodeString(boc)
	if err != nil {
		return nil, err
	}

	return cell.FromBOC(emCellBytes)
}

func CalTxHash(boc string) (string, error) {
	emCell, err := getExternalMessageCell(boc)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(emCell.Hash()), nil
}

func CalNormMsgHash(boc string) (string, error) {
	emCell, err := getExternalMessageCell(boc)
	if err != nil {
		return "", err
	}

	var message tlb.Message
	err = tlb.LoadFromCell(&message, emCell.BeginParse())
	if err != nil {
		return "", err
	}

	externalMessage, ok := message.Msg.(*tlb.ExternalMessage)
	if !ok {
		return "", errors.New("invalid external message")
	}
	return hex.EncodeToString(externalMessage.NormalizedHash()), nil
}
