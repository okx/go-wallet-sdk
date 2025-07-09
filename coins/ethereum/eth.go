package ethereum

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/okx/go-wallet-sdk/util"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
)

type EVMTx struct {
	TxType            int
	ChainId           *big.Int
	Tx                *EthTransaction
	Tx1559            *types.Transaction
	AuthorizationList []*EthAuthorization
}

func SignTx(evmTx *EVMTx, privateKey *btcec.PrivateKey) (string, error) {
	var signedTxByte []byte
	var err error
	if evmTx.TxType == DynamicFeeTxType {
		signedTxByte, err = SignEip1559Tx(evmTx.ChainId, evmTx.Tx1559, privateKey.ToECDSA())
	} else if evmTx.TxType == AuthorizationTxType {
		signedTxByte, err = SignEip7702Tx(evmTx.Tx1559, evmTx.AuthorizationList, evmTx.ChainId, privateKey)
	} else {
		signedTxByte, err = SignLegacyTx(evmTx.Tx, evmTx.ChainId, privateKey)
	}
	if err != nil {
		return "", err
	}
	return util.EncodeHexWithPrefix(signedTxByte), nil
}

func GenUnsignedTx(evmTx *EVMTx) (string, error) {
	if evmTx.TxType == DynamicFeeTxType {
		return GenUnsignedEip1559Tx(evmTx.Tx1559, evmTx.ChainId)
	} else if evmTx.TxType == AuthorizationTxType {
		return GenUnsignedEip7702Tx(evmTx.Tx1559, evmTx.AuthorizationList, evmTx.ChainId)
	}
	return evmTx.Tx.GenUnsignedTx(evmTx.ChainId), nil
}

func GenTxWithSig(txType int, chainId, unsignedRawTx, r, s, v string) (string, error) {
	chainID, _ := new(big.Int).SetString(chainId, 10)
	R, _ := new(big.Int).SetString(r, 16)
	S, _ := new(big.Int).SetString(s, 16)
	V, _ := new(big.Int).SetString(v, 16)

	if txType == DynamicFeeTxType {
		return GenEip1559TxWithSig(unsignedRawTx, chainID, R, S, V)
	} else if txType == AuthorizationTxType {
		return GenEip7702TxWithSig(unsignedRawTx, chainID, R, S, V)
	} else {
		return GenLegacyTxWithSig(unsignedRawTx, chainID, R, S, V)
	}
}

func CalcSignHash(data []byte, addPrefix bool) []byte {
	if addPrefix {
		msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
		s := sha3.NewLegacyKeccak256()
		s.Write([]byte(msg))
		return s.Sum(nil)
	}
	return data
}

func CalTxHash(tx string) ([]byte, error) {
	bytes, err := util.DecodeHexStringErr(tx)
	if err != nil {
		return nil, err
	}

	hash256 := sha3.NewLegacyKeccak256()
	hash256.Write(bytes)
	return hash256.Sum(nil), nil
}

func SignMessage(message []byte, prvKey *btcec.PrivateKey) *SignatureData {
	hash256 := sha3.NewLegacyKeccak256()
	hash256.Write(message)
	messageHash := hash256.Sum(nil)
	return SignAsRecoverable(messageHash, prvKey)
}

func SignEthTypeMessage(message string, prvKey *btcec.PrivateKey, addPrefix bool) (string, error) {
	msg := util.RemoveHexPrefix(message)
	msgData, err := hex.DecodeString(msg)
	if err != nil {
		msgData = []byte(msg)
	}
	res := SignAsRecoverable(CalcSignHash(msgData, addPrefix), prvKey)
	minV := big.NewInt(27)
	if res.V.Cmp(minV) == -1 {
		res.V.Add(res.V, minV)
	}
	r := append(append(res.ByteR, res.ByteS...), res.V.Bytes()...)
	return util.EncodeHex(r), nil
}

func SignAsRecoverable(value []byte, prvKey *btcec.PrivateKey) *SignatureData {
	sig, _ := ecdsa.SignCompact(prvKey, value, false)

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
	}
}

func VerifySignMsg(signature, message, address string, addPrefix bool) error {
	addr, err := EcRecover(signature, message, addPrefix)
	if err != nil {
		return err
	}
	if addr == address {
		return nil
	}
	return errors.New("invalid sign")
}

func EcRecover(signature, message string, addPrefix bool) (string, error) {
	publicKey, err := EcRecoverPubKey(signature, message, addPrefix)
	if publicKey == nil {
		return "", err
	}
	return GetNewAddress(publicKey), nil
}

func EcRecoverPubKey(signature, message string, addPrefix bool) (*btcec.PublicKey, error) {
	signatureData := util.DecodeHexStringPad(signature)
	if len(signatureData) < 65 {
		return nil, errors.New("signature too short")
	}
	R := signatureData[:33]
	S := signatureData[33:64]
	V := signatureData[64:65]
	realData := append(append(V, R...), S...)
	msg := util.RemoveHexPrefix(message)
	msgData, err := hex.DecodeString(msg)
	if err != nil {
		msgData = []byte(msg)
	}
	hash := CalcSignHash(msgData, addPrefix)
	publicKey, _, err := ecdsa.RecoverCompact(realData, hash)
	if publicKey == nil {
		return nil, err
	}
	return publicKey, nil
}

func EcRecoverBytes(signature, message []byte, addPrefix bool) (string, error) {
	publicKeyBytes, err := EcRecoverPubKeyBytes(signature, message, addPrefix)
	if err != nil {
		return "", err
	}
	publicKey, err := btcec.ParsePubKey(publicKeyBytes)
	if err != nil {
		return "", err
	}
	return GetNewAddress(publicKey), nil
}

func EcRecoverPubKeyBytes(signature, message []byte, addPrefix bool) ([]byte, error) {
	if len(signature) < 65 {
		return nil, errors.New("signature too short")
	}
	R := signature[:33]
	S := signature[33:64]
	V := signature[64:65]
	realData := append(append(V, R...), S...)
	hash := CalcSignHash(message, addPrefix)
	publicKey, _, err := ecdsa.RecoverCompact(realData, hash)
	if err != nil {
		return nil, err
	}
	return publicKey.SerializeUncompressed(), nil
}
