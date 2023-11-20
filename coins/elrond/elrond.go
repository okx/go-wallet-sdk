package elrond

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"golang.org/x/crypto/sha3"
)

const (
	HRP        = "erd"
	zeroString = "0"
)

func NewAddress(privateKey string) (string, error) {
	pkBytes, _ := hex.DecodeString(privateKey)
	if len(pkBytes) != 64 {
		return "", errors.New("length of private key must 64 bytes")
	}
	key := ed25519.PrivateKey(pkBytes)
	pb := key.Public().(ed25519.PublicKey)
	address, err := bech32.EncodeFromBase256(HRP, pb)
	if err != nil {
		return "", err
	}
	return address, nil
}

func ValidateAddress(address string) bool {
	h, _, err := bech32.DecodeToBase256(address)
	return err == nil && h == HRP
}

func AddressFromSeed(seed string) (string, error) {
	seedBytes, _ := hex.DecodeString(seed)
	if len(seedBytes) != 32 {
		return "", errors.New("length of private key must 32 bytes")
	}
	privateKey := ed25519.NewKeyFromSeed(seedBytes)
	pb := privateKey.Public().(ed25519.PublicKey)
	address, _ := bech32.EncodeFromBase256(HRP, pb)
	return address, nil
}

func Transfer(args ArgCreateTransaction, privateKeyHex string) (string, error) {
	pk, _ := hex.DecodeString(privateKeyHex)
	privateKey := ed25519.NewKeyFromSeed(pk)
	builder := NewTxBuilder(&privateKey)
	tran, err := builder.build(args)
	if err != nil {
		return "", err
	}
	ss, err := json.Marshal(tran)
	if err != nil {
		return "", err
	}
	return string(ss), nil
}

func GetNewAddressByPub(pubkey string) (string, error) {
	pb, err := hex.DecodeString(pubkey)
	if err != nil {
		return "", err
	}
	address, err := bech32.EncodeFromBase256(HRP, pb)
	if err != nil {
		return "", err
	}
	return address, nil
}

func UnsignedTx(arg ArgCreateTransaction, sndAddr string, signature string) (string, error) {
	tx := Transaction{
		Nonce:    arg.Nonce,
		Value:    arg.Value,
		RcvAddr:  arg.RcvAddr,
		SndAddr:  sndAddr,
		GasPrice: arg.GasPrice,
		GasLimit: arg.GasLimit,
		Data:     arg.Data,
		ChainID:  arg.ChainID,
		Version:  arg.Version,
		Options:  arg.Options,
		// The digital signature consisting of 128 hex-characters (thus 64 bytes in a raw representation)
		Signature: signature,
	}
	unsignedMessage, _ := json.Marshal(tx)
	shouldSignOnTxHash := arg.Version >= 2 && arg.Options&1 > 0
	if shouldSignOnTxHash {
		hasher := sha3.NewLegacyKeccak256()
		hasher.Write(unsignedMessage[:])
		unsignedMessage = hasher.Sum(nil)
	}
	return hex.EncodeToString(unsignedMessage), nil
}

func SignedTx(arg ArgCreateTransaction, sender string, signature string) string {
	tx := Transaction{
		Nonce:    arg.Nonce,
		Value:    arg.Value,
		RcvAddr:  arg.RcvAddr,
		SndAddr:  sender,
		GasPrice: arg.GasPrice,
		GasLimit: arg.GasLimit,
		Data:     arg.Data,
		ChainID:  arg.ChainID,
		Version:  arg.Version,
		Options:  arg.Options,
		// The digital signature consisting of 128 hex-characters (thus 64 bytes in a raw representation)
		Signature: signature,
	}
	ss, _ := json.Marshal(tx)
	return string(ss)
}
