package core

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/zksync/zkscrypto"
	"math/big"
	"strings"
)

const (
	Message                 = "Access zkSync account.\n\nOnly sign this message for a trusted client!"
	TransactionVersion byte = 0x01
)

func NewZkSignerFromSeed(seed []byte) (*ZkSigner, error) {
	privateKey, err := zkscrypto.NewPrivateKey(seed)
	if err != nil {
		return nil, errors.New("failed to create private key")
	}
	return newZkSignerFromPrivateKey(privateKey)
}

func NewZkSignerFromEthSigner(es EthSigner, cid ChainId) (*ZkSigner, error) {
	signMsg := Message
	if cid != ChainIdMainnet {
		signMsg = fmt.Sprintf("%s\nChain ID: %d.", Message, cid)
	}
	sig, err := es.SignMessage([]byte(signMsg))
	if err != nil {
		return nil, errors.New("failed to sign special message")
	}
	return NewZkSignerFromSeed(sig)
}

func newZkSignerFromPrivateKey(privateKey *zkscrypto.PrivateKey) (*ZkSigner, error) {
	publicKey, err := privateKey.PublicKey()
	if err != nil {
		return nil, errors.New("failed to create public key")
	}
	publicKeyHash, err := publicKey.Hash()
	if err != nil {
		return nil, errors.New("failed to get public key hash")
	}
	return &ZkSigner{
		privateKey:    privateKey,
		publicKey:     publicKey,
		publicKeyHash: publicKeyHash,
	}, nil
}

type ZkSigner struct {
	privateKey    *zkscrypto.PrivateKey
	publicKey     *zkscrypto.PublicKey
	publicKeyHash *zkscrypto.PublicKeyHash
}

func (s *ZkSigner) Sign(message []byte) (*zkscrypto.Signature, error) {
	signature, err := s.privateKey.Sign(message)
	if err != nil {
		return nil, errors.New("failed to sign message")
	}
	return signature, nil
}

func (s *ZkSigner) SignChangePubKey(txData *ChangePubKey) (*Signature, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xff - 0x07)
	buf.WriteByte(TransactionVersion)
	buf.Write(Uint32ToBytes(txData.AccountId))
	buf.Write(ParseAddress(txData.Account))
	pkhBytes, err := pkhToBytes(txData.NewPkHash)
	if err != nil {
		return nil, errors.New("failed to get pkh bytes")
	}
	buf.Write(pkhBytes)
	buf.Write(Uint32ToBytes(txData.FeeToken))
	fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
	if !ok {
		return nil, errors.New("failed to convert string fee to big.Int")
	}
	packedFee, err := packFee(fee)
	if err != nil {
		return nil, errors.New("failed to pack fee")
	}
	buf.Write(packedFee)
	buf.Write(Uint32ToBytes(txData.Nonce))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidFrom))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidUntil))
	sig, err := s.Sign(buf.Bytes())
	if err != nil {
		return nil, errors.New("failed to sign ChangePubKey tx data")
	}
	res := &Signature{
		PubKey:    s.GetPublicKey(),
		Signature: sig.HexString(),
	}
	return res, nil
}

func (s *ZkSigner) SignTransfer(txData *Transfer) (*Signature, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xff - 0x05)
	buf.WriteByte(TransactionVersion)
	buf.Write(Uint32ToBytes(txData.AccountId))
	buf.Write(ParseAddress(txData.From))
	buf.Write(ParseAddress(txData.To))
	buf.Write(Uint32ToBytes(txData.Token.Id))
	packedAmount, err := packAmount(txData.Amount)
	if err != nil {
		return nil, errors.New("failed to pack amount")
	}
	buf.Write(packedAmount)
	fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
	if !ok {
		return nil, errors.New("failed to convert string fee to big.Int")
	}
	packedFee, err := packFee(fee)
	if err != nil {
		return nil, errors.New("failed to pack fee")
	}
	buf.Write(packedFee)
	buf.Write(Uint32ToBytes(txData.Nonce))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidFrom))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidUntil))
	sig, err := s.Sign(buf.Bytes())
	if err != nil {
		return nil, errors.New("failed to sign Transfer tx data")
	}
	res := &Signature{
		PubKey:    s.GetPublicKey(),
		Signature: sig.HexString(),
	}
	return res, nil
}

func (s *ZkSigner) SignWithdraw(txData *Withdraw) (*Signature, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xff - 0x03)
	buf.WriteByte(TransactionVersion)
	buf.Write(Uint32ToBytes(txData.AccountId))
	buf.Write(ParseAddress(txData.From))
	buf.Write(ParseAddress(txData.To))
	buf.Write(Uint32ToBytes(txData.TokenId))
	amountBytes := txData.Amount.Bytes()
	buf.Write(make([]byte, 16-len(amountBytes))) // total amount slot is 16 bytes BE
	buf.Write(amountBytes)
	fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
	if !ok {
		return nil, errors.New("failed to convert string fee to big.Int")
	}
	packedFee, err := packFee(fee)
	if err != nil {
		return nil, errors.New("failed to pack fee")
	}
	buf.Write(packedFee)
	buf.Write(Uint32ToBytes(txData.Nonce))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidFrom))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidUntil))
	sig, err := s.Sign(buf.Bytes())
	if err != nil {
		return nil, errors.New("failed to sign Withdraw tx data")
	}
	res := &Signature{
		PubKey:    s.GetPublicKey(),
		Signature: sig.HexString(),
	}
	return res, nil
}

func (s *ZkSigner) SignForcedExit(txData *ForcedExit) (*Signature, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xff - 0x08)
	buf.WriteByte(TransactionVersion)
	buf.Write(Uint32ToBytes(txData.AccountId))
	buf.Write(ParseAddress(txData.Target))
	buf.Write(Uint32ToBytes(txData.TokenId))
	fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
	if !ok {
		return nil, errors.New("failed to convert string fee to big.Int")
	}
	packedFee, err := packFee(fee)
	if err != nil {
		return nil, errors.New("failed to pack fee")
	}
	buf.Write(packedFee)
	buf.Write(Uint32ToBytes(txData.Nonce))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidFrom))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidUntil))
	sig, err := s.Sign(buf.Bytes())
	if err != nil {
		return nil, errors.New("failed to sign ForcedExit tx data")
	}
	res := &Signature{
		PubKey:    s.GetPublicKey(),
		Signature: sig.HexString(),
	}
	return res, nil
}

func (s *ZkSigner) SignMintNFT(txData *MintNFT) (*Signature, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xff - 0x09)
	buf.WriteByte(TransactionVersion)
	buf.Write(Uint32ToBytes(txData.CreatorId))
	buf.Write(ParseAddress(txData.CreatorAddress))
	buf.Write(txData.ContentHash.Bytes())
	buf.Write(ParseAddress(txData.Recipient))
	buf.Write(Uint32ToBytes(txData.FeeToken))
	fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
	if !ok {
		return nil, errors.New("failed to convert string fee to big.Int")
	}
	packedFee, err := packFee(fee)
	if err != nil {
		return nil, errors.New("failed to pack fee")
	}
	buf.Write(packedFee)
	buf.Write(Uint32ToBytes(txData.Nonce))
	sig, err := s.Sign(buf.Bytes())
	if err != nil {
		return nil, errors.New("failed to sign MintNFT tx data")
	}
	res := &Signature{
		PubKey:    s.GetPublicKey(),
		Signature: sig.HexString(),
	}
	return res, nil
}

func (s *ZkSigner) SignWithdrawNFT(txData *WithdrawNFT) (*Signature, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xff - 0x0a)
	buf.WriteByte(TransactionVersion)
	buf.Write(Uint32ToBytes(txData.AccountId))
	buf.Write(ParseAddress(txData.From))
	buf.Write(ParseAddress(txData.To))
	buf.Write(Uint32ToBytes(txData.Token))
	buf.Write(Uint32ToBytes(txData.FeeToken))
	fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
	if !ok {
		return nil, errors.New("failed to convert string fee to big.Int")
	}
	packedFee, err := packFee(fee)
	if err != nil {
		return nil, errors.New("failed to pack fee")
	}
	buf.Write(packedFee)
	buf.Write(Uint32ToBytes(txData.Nonce))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidFrom))
	buf.Write(Uint64ToBytes(txData.TimeRange.ValidUntil))
	sig, err := s.Sign(buf.Bytes())
	if err != nil {
		return nil, errors.New("failed to sign WithdrawNFT tx data")
	}
	res := &Signature{
		PubKey:    s.GetPublicKey(),
		Signature: sig.HexString(),
	}
	return res, nil
}

func (s *ZkSigner) GetPublicKeyHash() string {
	return "sync:" + s.publicKeyHash.HexString()
}

func (s *ZkSigner) GetPublicKey() string {
	return s.publicKey.HexString()
}

func ParseAddress(address string) []byte {
	var minAddress = address
	if strings.HasPrefix(address, "0x") {
		minAddress = minAddress[2:]
	}

	addrBytes, _ := hex.DecodeString(minAddress)
	return addrBytes
}
