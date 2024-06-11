package ethereum

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okx/go-wallet-sdk/crypto"
	"github.com/okx/go-wallet-sdk/util"
	"golang.org/x/crypto/sha3"
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
)

type EthTransaction struct {
	Nonce    *big.Int `json:"nonce"`
	GasPrice *big.Int `json:"gasPrice"`
	GasLimit *big.Int `json:"gas"`
	To       []byte   `json:"to"`
	Value    *big.Int `json:"value"`
	Data     []byte   `json:"data"`
	// Signature values
	V *big.Int `json:"v"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

func (tx *EthTransaction) SignTransaction(chainId *big.Int, prvKey *btcec.PrivateKey) (string, error) {
	tx.V = chainId
	rawTransaction, err := rlp.EncodeToBytes([]interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
		chainId, uint(0), uint(0),
	})
	if err != nil {
		return "", err
	}
	sig, err := SignMessage(rawTransaction, prvKey)
	if err != nil {
		return "", err
	}
	tx.V = big.NewInt(chainId.Int64()*2 + sig.V.Int64() + 8)
	tx.R = sig.R
	tx.S = sig.S
	value, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return "", err
	}
	return util.EncodeHexWith0x(value), nil
}

func (tx *EthTransaction) UnSignedTx(chainId *big.Int) (string, error) {
	tx.V = chainId
	rawTransaction, err := rlp.EncodeToBytes([]interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
		chainId, uint(0), uint(0),
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawTransaction), nil
}

func (tx *EthTransaction) GetSigningHash(chainId *big.Int) (string, string, error) {
	unSignedTx, err := tx.UnSignedTx(chainId)
	if err != nil {
		return "", "", err
	}
	raw, err := hex.DecodeString(unSignedTx)
	if err != nil {
		return "", "", err
	}
	h := sha3.NewLegacyKeccak256()
	h.Write(raw)
	msgHash := h.Sum(nil)
	return hex.EncodeToString(msgHash), unSignedTx, nil
}

func (tx *EthTransaction) SignedTx(chainId *big.Int, sig *SignatureData) (string, error) {
	tx.V = big.NewInt(chainId.Int64()*2 + sig.V.Int64() + 8)
	tx.R = sig.R
	tx.S = sig.S
	value, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return "", err
	}
	return util.EncodeHexWith0x(value), nil
}

func SignMessage(message []byte, prvKey *btcec.PrivateKey) (*SignatureData, error) {
	hash256 := sha3.NewLegacyKeccak256()
	hash256.Write(message)
	messageHash := hash256.Sum(nil)
	return SignAsRecoverable(messageHash, prvKey)
}

func SignEthTypeMessage(message string, prvKey *btcec.PrivateKey, addPrefix bool) (string, error) {
	// support hex message and non-hex message
	msg := OnlyRemovePrefix(message)
	msgData, err := hex.DecodeString(msg)
	if err != nil {
		msgData = []byte(msg)
	}
	res, err := SignAsRecoverable(signHash(msgData, addPrefix), prvKey)
	if err != nil {
		return "", err
	}
	minV := big.NewInt(27)
	if res.V.Cmp(minV) == -1 {
		res.V.Add(res.V, minV)
	}
	r, err := hex.DecodeString(hex.EncodeToString(res.ByteR) + hex.EncodeToString(res.ByteS) + hex.EncodeToString(res.V.Bytes()))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(r), nil
}

func signHash(data []byte, addPrefix bool) []byte {
	if addPrefix {
		msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
		s := sha3.NewLegacyKeccak256()
		s.Write([]byte(msg))
		return s.Sum(nil)
	}
	return data
}

func NewEthTransaction(nonce, gasLimit, gasPrice, value *big.Int, to, data string) *EthTransaction {
	toBytes := util.RemoveZeroHex(to)
	dataBytes := util.RemoveZeroHex(data)
	return &EthTransaction{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		To:       toBytes,
		Value:    value,
		Data:     dataBytes,
	}
}

func NewUnsignedTx(nonce, gasLimit, gasPrice, value, chainId *big.Int, to, data string) (*UnsignedTx, error) {
	tx := NewEthTransaction(nonce, gasLimit, gasPrice, value, to, data)
	data, hash, err := tx.GetSigningHash(chainId)
	if err != nil {
		return nil, err
	}
	return &UnsignedTx{Tx: data, Hash: hash}, nil
}

func NewTransactionFromRaw(raw string) (*EthTransaction, error) {
	bytes := util.RemoveZeroHex(raw)
	t := new(EthTransaction)
	err := rlp.DecodeBytes(bytes, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func SignAsRecoverable(value []byte, prvKey *btcec.PrivateKey) (*SignatureData, error) {
	sig, err := ecdsa.SignCompact(prvKey, value, false)
	if err != nil {
		return nil, err
	}
	V := sig[0]
	R := sig[1:33]
	S := sig[33:65]
	return &SignatureData{
		V:     new(big.Int).SetBytes([]byte{V}),
		R:     new(big.Int).SetBytes(R),
		S:     new(big.Int).SetBytes(S),
		ByteV: V,
		ByteR: R,
		ByteS: S,
	}, nil
}

func VerifySignMsg(signature, message, address string, addPrefix bool) error {
	addr, err := EcRecover(signature, message, addPrefix)
	if err != nil {
		return err
	}
	if addr == address {
		return nil
	}
	return errors.New("invali sign")
}

func EcRecover(signature, message string, addPrefix bool) (string, error) {
	publicKey, err := EcRecoverPubKey(signature, message, addPrefix)
	if publicKey == nil {
		return "", err
	}
	return GetNewAddress(publicKey), nil
}

func GetEthGroupAddress(prefix string, pubKey *btcec.PublicKey) string {
	addressByte := GetEthGroupPubHash(pubKey)
	return prefix + hex.EncodeToString(addressByte[12:])
}

func GetEthGroupPubHash(pubKey *btcec.PublicKey) []byte {
	pubBytes := pubKey.SerializeUncompressed()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:])
	addressByte := hash.Sum(nil)
	return addressByte
}

func EcRecoverPubKey(signature, message string, addPrefix bool) (*btcec.PublicKey, error) {
	signatureData := util.RemoveZeroHex(signature)
	R := signatureData[:33]
	S := signatureData[33:64]
	V := signatureData[64:65]
	realData, err := hex.DecodeString(hex.EncodeToString(V) + hex.EncodeToString(R) + hex.EncodeToString(S))
	if err != nil {
		return nil, err
	}
	// support hex message or non-hex message
	msg := OnlyRemovePrefix(message)
	msgData, err := hex.DecodeString(msg)
	if err != nil {
		msgData = []byte(msg)
	}
	hash := signHash(msgData, addPrefix)
	publicKey, _, err := ecdsa.RecoverCompact(realData, hash)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

type SignatureData struct {
	V *big.Int
	R *big.Int
	S *big.Int

	ByteV byte
	ByteR []byte
	ByteS []byte
}

func NewSignatureData(msgHash []byte, publicKey string, r, s *big.Int) (*SignatureData, error) {
	// Calculate v, r, and s
	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}

	pubKey, err := btcec.ParsePubKey(pubBytes)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.SignCompact(btcec.S256(), r, s, *pubKey, msgHash, false)
	if err != nil {
		return nil, err
	}

	V := sig[0]
	R := sig[1:33]
	S := sig[33:65]
	return &SignatureData{
		V:     new(big.Int).SetBytes([]byte{V}),
		R:     new(big.Int).SetBytes(R),
		S:     new(big.Int).SetBytes(S),
		ByteV: V,
		ByteR: R,
		ByteS: S,
	}, nil
}

func (sd *SignatureData) ToHex() string {
	return hex.EncodeToString(sd.ToBytes())
}

func (sd SignatureData) ToBytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, sd.ByteR...)
	bytes = append(bytes, sd.ByteS...)
	bytes = append(bytes, sd.ByteV)
	return bytes
}

func GetNewAddressBytes(pubKey *btcec.PublicKey) []byte {
	pubBytes := pubKey.SerializeUncompressed()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:])
	addressByte := hash.Sum(nil)
	return addressByte[12:]
}

func GetNewAddress(pubKey *btcec.PublicKey) string {
	return "0x" + hex.EncodeToString(GetNewAddressBytes(pubKey))
}

func GetEthereumMessagePrefix(message string) string {
	return fmt.Sprintf(MessagePrefixTmp, len(message))
}

func PubKeyToAddr(publicKey []byte) (string, error) {
	pubKey, err := btcec.ParsePubKey(publicKey)
	if err != nil {
		return "", err
	}
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKey.SerializeUncompressed()[1:])
	addressByte := hash.Sum(nil)
	return "0x" + hex.EncodeToString(addressByte[12:]), nil
}
