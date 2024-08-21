package ton

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"
	"math/big"

	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"github.com/okx/go-wallet-sdk/coins/ton/tlb"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/jetton"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
)

func buildTx(w *wallet.Wallet, withInit bool) (*SignedTx, error) {
	if w == nil {
		return nil, errors.New("invalid wallet")
	}
	if !withInit {
		return &SignedTx{
			Address: w.WalletAddress().String(),
		}, nil
	}
	stateInit, err := wallet.GetStateInit(w.PublicKey(), wallet.V4R2, wallet.DefaultSubwallet)
	if err != nil {
		return nil, err
	}

	return &SignedTx{
		Code:    base64.StdEncoding.EncodeToString(stateInit.Code.ToBOC()),
		Data:    base64.StdEncoding.EncodeToString(stateInit.Data.ToBOC()),
		Address: w.WalletAddress().String(),
	}, nil
}

func newWallet(seed, pubKey []byte) (*wallet.Wallet, error) {
	if len(pubKey) == ed25519.PublicKeySize && len(seed) == ed25519.SeedSize {
		if bytes.Equal(pubKey, ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)) {
			return wallet.FromPrivateKey(ed25519.NewKeyFromSeed(seed), wallet.V4R2)
		}
	}
	if len(pubKey) > 0 {
		return wallet.FakeFromPublicKey(ed25519.PublicKey(pubKey), wallet.V4R2)
	}
	return wallet.FromPrivateKey(ed25519.NewKeyFromSeed(seed), wallet.V4R2)
}

func Transfer(seed, pubKey []byte, to, amount, comment string, seqno uint32, expireAt int64, mode uint8, simulate bool) (*SignedTx, error) {
	w, err := newWallet(seed, pubKey)
	if err != nil {
		return nil, err
	}
	specV4R2 := w.GetSpec().(*wallet.SpecV4R2)
	specV4R2.SetCustomSeqnoFetcher(func() uint32 {
		return seqno
	})
	specV4R2.SetExpireAt(expireAt)
	toAddr, err := address.ParseAddr(to)

	if err != nil {
		return nil, err
	}
	toAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, err
	}
	message, err := w.BuildTransfer(toAddr, tlb.FromNanoTON(toAmount), false, comment, mode)
	if err != nil {
		return nil, err
	}
	initialized := false
	if seqno > 0 {
		initialized = true
	}

	externalMessage, err := w.BuildExternalMessage(context.Background(), message, initialized)
	if err != nil {
		return nil, err
	}

	if simulate {
		signedTx, err := buildTx(w, seqno == 0)
		if err != nil {
			return nil, err
		}
		signedTx.FillTx(base64.StdEncoding.EncodeToString(externalMessage.Body.ToBOC()))
		return signedTx, nil
	}
	emCell, err := tlb.ToCell(externalMessage)
	if err != nil {
		return nil, err
	}
	signedTx, err := buildTx(w, seqno == 0)
	if err != nil {
		return nil, err
	}
	signedTx.FillTx(base64.StdEncoding.EncodeToString(emCell.ToBOC()))
	return signedTx, nil
}

func TransferJetton(seed, pubKey []byte, from, to, amount string, decimals int, seqno uint32, messageAttachedTons string, invokeNotificationFee string, comment string, expireAt int64, rnd uint64, simulate bool) (*SignedTx, error) {
	w, err := newWallet(seed, pubKey)
	if err != nil {
		return nil, err
	}
	fromAddr, err := address.ParseAddr(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := address.ParseAddr(to)
	if err != nil {
		return nil, err
	}
	specV4R2 := w.GetSpec().(*wallet.SpecV4R2)
	specV4R2.SetCustomSeqnoFetcher(func() uint32 {
		return seqno
	})
	specV4R2.SetExpireAt(expireAt)
	responseToAddress := w.Address()
	toAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, err
	}
	amountForwardTON := tlb.ZeroCoins
	if invokeNotificationFee != "" {
		invokenFee, ok := new(big.Int).SetString(invokeNotificationFee, 10)
		if !ok {
			return nil, err
		}
		amountForwardTON = tlb.FromNanoTON(invokenFee)
	}
	var payloadForward *cell.Cell
	if comment != "" {
		payloadForward = cell.BeginCell().MustStoreUInt(0, 32).MustStoreStringSnake(comment).EndCell()
	}
	jw := &jetton.WalletClient{}
	toAmountCoins, err := tlb.FromNano(toAmount, decimals)
	if err != nil {
		return nil, err
	}
	transferPayload, err := jw.BuildTransferPayloadV2(toAddr, responseToAddress, toAmountCoins, amountForwardTON, payloadForward, nil, rnd)
	if err != nil {
		return nil, err
	}

	messageAttachedVal := "50000000"
	if messageAttachedTons != "" {
		messageAttachedVal = messageAttachedTons
	}
	messageAttachedValBig, ok := new(big.Int).SetString(messageAttachedVal, 10)
	if !ok {
		return nil, err
	}
	message := wallet.SimpleMessage(fromAddr, tlb.FromNanoTON(messageAttachedValBig), transferPayload)

	initialized := false
	if seqno > 0 {
		initialized = true
	}

	externalMessage, err := w.BuildExternalMessage(context.Background(), message, initialized)
	if err != nil {
		return nil, err
	}
	if simulate {
		signedTx, err := buildTx(w, seqno == 0)
		if err != nil {
			return nil, err
		}
		signedTx.FillTx(base64.StdEncoding.EncodeToString(externalMessage.Body.ToBOCWithFlags(false)))
		return signedTx, nil
	}
	emCell, err := tlb.ToCell(externalMessage)
	if err != nil {
		return nil, err
	}
	signedTx, err := buildTx(w, seqno == 0)
	if err != nil {
		return nil, err
	}
	signedTx.FillTx(base64.StdEncoding.EncodeToString(emCell.ToBOCWithFlags(false)))
	return signedTx, nil
}

// VenomTransfer venom chain use v3 til now
func VenomTransfer(seed []byte, to, amount, comment string, seqno uint32, bounce bool, globalID uint32, expireAt int64, mode uint8) (string, error) {
	w, err := wallet.FromPrivateKeyVenom(ed25519.NewKeyFromSeed(seed), wallet.VenomV3)
	specVenomV3 := w.GetSpec().(*wallet.SpecVenomV3)
	specVenomV3.SetCustomSeqnoFetcher(func() uint32 {
		return seqno
	})
	specVenomV3.SetExpireAt(expireAt)
	toAddr, err := address.ParseRawAddr(to)
	if err != nil {
		return "", err
	}
	toAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return "", err
	}
	message, err := w.BuildTransfer(toAddr, tlb.FromNanoTON(toAmount), bounce, comment, mode)
	if err != nil {
		return "", err
	}
	// specVenomV3.SetGlobalID(globalID)

	initialized := false
	if seqno > 0 {
		initialized = true
	}

	externalMessage, err := w.BuildExternalMessage(context.Background(), message, initialized)
	if err != nil {
		return "", err
	}
	emCell, err := tlb.ToCell(externalMessage)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(struct {
		Id   string `json:"id"`
		Body string `json:"body"`
	}{
		base64.StdEncoding.EncodeToString(emCell.Hash()),
		base64.StdEncoding.EncodeToString(emCell.ToBOCWithFlags(false)),
	})
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func CalTxHash(boc string) (string, error) {
	emCellBytes, err := base64.StdEncoding.DecodeString(boc)
	if err != nil {
		return "", err
	}

	emCell, err := cell.FromBOC(emCellBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(emCell.Hash()), nil
}
