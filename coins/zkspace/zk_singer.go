package zkspace

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/zksync/core"
	"github.com/okx/go-wallet-sdk/coins/zksync/zkscrypto"
)

const (
	Message                 = "Access ZKSwap account.\n\nOnly sign this message for a trusted client!"
	TransactionVersion byte = 0x01
)

func NewZkSignerFromSeed(seed []byte) (*ZkSigner, error) {
	privateKey, err := zkscrypto.NewPrivateKey(seed)
	if err != nil {
		return nil, err
	}
	return newZkSignerFromPrivateKey(privateKey)
}

func NewZkSignerFromRawPrivateKey(rawPk []byte) (*ZkSigner, error) {
	privateKey, err := zkscrypto.NewPrivateKeyRaw(rawPk)
	if err != nil {
		return nil, err
	}
	return newZkSignerFromPrivateKey(privateKey)
}

func NewZkSignerFromEthSigner(es core.EthSigner, cid core.ChainId) (*ZkSigner, error) {
	signMsg := Message
	if cid != core.ChainIdMainnet {
		signMsg = fmt.Sprintf("%s\nChain ID: %d.", Message, cid)
	}
	sig, err := es.SignMessage([]byte(signMsg))
	if err != nil {
		return nil, err
	}
	return NewZkSignerFromSeed(sig)
}

func newZkSignerFromPrivateKey(privateKey *zkscrypto.PrivateKey) (*ZkSigner, error) {
	publicKey, err := privateKey.PublicKey()
	if err != nil {
		return nil, err
	}
	publicKeyHash, err := publicKey.Hash()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return signature, nil
}

func (s *ZkSigner) SignTransfer(txData *Transfer) (*Signature, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(0x05)
	buf.Write(core.Uint32ToBytes(txData.AccountId))
	fromBytes, err := hex.DecodeString(txData.From[2:])
	if err != nil {
		return nil, err
	}
	buf.Write(fromBytes)
	toBytes, err := hex.DecodeString(txData.To[2:])
	if err != nil {
		return nil, err
	}
	buf.Write(toBytes)
	buf.Write(core.Uint16ToBytes(txData.TokenId))
	packedAmount, err := core.PackAmount(txData.Amount)
	if err != nil {
		return nil, err
	}
	buf.Write(packedAmount)
	buf.WriteByte(txData.FeeTokenId)
	packedFee, err := core.PackFee(txData.Fee)
	if err != nil {
		return nil, err
	}
	buf.Write(packedFee)
	buf.WriteByte(txData.ChainId)
	buf.Write(core.Uint32ToBytes(txData.Nonce))
	sig, err := s.Sign(buf.Bytes())
	if err != nil {
		return nil, err
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

func (s *ZkSigner) GetPrivateKey() string {
	return s.privateKey.HexString()
}
