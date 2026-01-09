package ton

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"github.com/okx/go-wallet-sdk/coins/ton/tlb"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
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

func (b *TonTransferBuilder) BuildMultiTransferMessages(request *MultiRequest) ([]*wallet.Message, error) {
	messages := make([]*wallet.Message, len(request.Messages))
	for k, v := range request.Messages {
		to, err := address.ParseAddr(v.Address)
		if err != nil {
			return nil, err
		}
		toAmount, ok := new(big.Int).SetString(v.Amount, 10)
		if !ok {
			return nil, err
		}

		if v.ExtraFlags == "" {
			v.ExtraFlags = "0"
		}
		extraFlags, err := tlb.FromDecimal(v.ExtraFlags, 0)
		if err != nil {
			return nil, err
		}
		message, err := b.wallet.BuildTransferByBody(to, tlb.FromNanoTON(toAmount), v.Payload, v.StateInit, extraFlags)
		if err != nil {
			return nil, err
		}
		messages[k] = message
	}

	b.messages = messages
	return messages, nil
}
func SignMultiTransfer(seed, pub []byte, seqno uint32, request *MultiRequest, simulate bool, version wallet.Version) (*SignedTx, error) {
	params := &TransferParams{
		Seed:     seed,
		PubKey:   pub,
		Seqno:    seqno,
		Simulate: simulate,
		Version:  version,
		ExpireAt: request.ValidUntil,
	}
	b, err := NewTonTransferBuilder(params)
	if err != nil {
		return nil, err
	}
	_, err = b.BuildMultiTransferMessages(request)
	if err != nil {
		return nil, err
	}
	_, err = b.BuildTransferDirect()
	if err != nil {
		return nil, err
	}
	return b.BuildSignedTx(true)
}

func SignProof(addr string, seed []byte, payload *ProofData) (string, error) {
	if payload == nil || len(seed) != ed25519.SeedSize {
		return "", errors.New("invalid param")
	}
	msg, err := CreateMessage(addr, payload)
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
	msg, err := CreateMessage(addr, payload)
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

func GetWalletInformation(seed, pubKey []byte, version wallet.Version) (*AccontInfo, error) {
	w, err := NewWallet(seed, pubKey, version)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, errors.New("invalid wallet")
	}
	state, err := wallet.GetStateInit(w.PublicKey(), w.GetVersionConfig(), w.GetSubwalletID())
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

func CreateMessage(addr string, message *ProofData) ([]byte, error) {
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
