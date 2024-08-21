package ton

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"github.com/okx/go-wallet-sdk/coins/ton/tlb"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
	"math/big"
)

const (
	tonProofPrefix   = "ton-proof-item-v2/"
	tonConnectPrefix = "ton-connect"
)

var (
	ErrInvalidParam = errors.New("invalid param")
	ErrInvalidSign  = errors.New("invalid sign")
)

type ProofData struct {
	Timestamp uint64 `json:"timestamp"`
	Domain    string `json:"domain"`
	Payload   string `json:"payload"`
}

func SignMultiTransfer(seed, pub []byte, seqno uint32, request *MultiRequest, simulate bool) (*SignedTx, error) {
	w, err := newWallet(seed, pub)
	if err != nil {
		return nil, err
	}
	specV4R2 := w.GetSpec().(*wallet.SpecV4R2)
	specV4R2.SetCustomSeqnoFetcher(func() uint32 {
		return seqno
	})
	expireAt := request.ValidUntil
	if expireAt < 0 {
		return nil, ErrInvalidMultiRequest
	}
	specV4R2.SetExpireAt(expireAt)
	initialized := true
	if seqno == 0 {
		initialized = false
	}

	msgs := make([]*wallet.Message, len(request.Messages))
	for k, v := range request.Messages {
		to, err := address.ParseAddr(v.Address)
		if err != nil {
			return nil, err
		}
		toAmount, ok := new(big.Int).SetString(v.Amount, 10)
		if !ok {
			return nil, err
		}
		vv, err := w.BuildTransferByBody(to, tlb.FromNanoTON(toAmount), v.Payload, v.StateInit)
		if err != nil {
			return nil, err
		}
		msgs[k] = vv
	}
	externalMessage, err := w.BuildExternalMessageForMany(context.Background(), msgs, initialized)
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

func SignProof(addr string, seed []byte, payload *ProofData) (string, error) {
	if payload == nil || len(seed) != ed25519.SeedSize {
		return "", errors.New("invalid param")
	}
	msg, err := createMessage(addr, payload)
	if err != nil {
		return "", err
	}
	prv := ed25519.NewKeyFromSeed(seed)
	return base64.StdEncoding.EncodeToString(ed25519.Sign(prv, msg)), nil
}

func VerifySignProofStr(addr string, pubHex, signBase64 string, payload *ProofData) error {
	if len(pubHex) == 0 || len(signBase64) == 0 {
		return ErrInvalidParam
	}
	pub, err := hex.DecodeString(pubHex)
	if err != nil {
		return ErrInvalidParam
	}
	sign, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return ErrInvalidParam
	}
	return VerifySignProof(addr, pub, sign, payload)
}

func VerifySignProof(addr string, pub, sign []byte, payload *ProofData) error {
	if payload == nil || len(pub) != ed25519.PublicKeySize || len(sign) == 0 {
		return ErrInvalidParam
	}
	msg, err := createMessage(addr, payload)
	if err != nil {
		return err
	}

	if !ed25519.Verify(ed25519.PublicKey(pub), msg, sign) {
		return ErrInvalidSign
	}
	return nil
}

type AccontInfo struct {
	InitCode        string `json:"initCode"`
	InitData        string `json:"initData"`
	WalletStateInit string `json:"walletStateInit"`
	WalletAddress   string `json:"walletAddress"`
}

func (s *AccontInfo) Str() (string, error) {
	if s == nil {
		return "", errors.New("invalid account info")
	}
	r, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(r), err
}

func GetWalletInformation(seed, pubKey []byte) (*AccontInfo, error) {
	w, err := newWallet(seed, pubKey)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, errors.New("invalid wallet")
	}
	state, err := wallet.GetStateInit(w.PublicKey(), wallet.V4R2, wallet.DefaultSubwallet)
	if err != nil {
		return nil, err
	}

	stateCell, err := tlb.ToCell(state)

	if err != nil {
		return nil, fmt.Errorf("failed to get state cell: %w", err)
	}
	return &AccontInfo{
		InitCode:        base64.StdEncoding.EncodeToString(state.Code.ToBOC()),
		InitData:        base64.StdEncoding.EncodeToString(state.Data.ToBOC()),
		WalletStateInit: base64.StdEncoding.EncodeToString(stateCell.ToBOC()),
		WalletAddress:   w.WalletAddress().String(),
	}, nil
}

func createMessage(addr string, message *ProofData) ([]byte, error) {
	addr2, err := address.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	wc := make([]byte, 4)
	binary.BigEndian.PutUint32(wc, uint32(addr2.Workchain()))

	ts := make([]byte, 8)
	binary.LittleEndian.PutUint64(ts, uint64(message.Timestamp))

	dl := make([]byte, 4)
	binary.LittleEndian.PutUint32(dl, uint32(len(message.Domain)))

	m := []byte(tonProofPrefix)
	m = append(m, wc...)
	m = append(m, addr2.Data()...)
	m = append(m, dl...)
	m = append(m, []byte(message.Domain)...)
	m = append(m, ts...)
	m = append(m, []byte(message.Payload)...)

	messageHash := sha256.Sum256(m)
	fullMes := []byte{0xff, 0xff}
	fullMes = append(fullMes, []byte(tonConnectPrefix)...)
	fullMes = append(fullMes, messageHash[:]...)

	res := sha256.Sum256(fullMes)
	return res[:], nil
}
