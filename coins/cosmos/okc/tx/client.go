package tx

import (
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/amino"
	authtypes "github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/auth/types"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/common/ethsecp256k1"
	"github.com/okx/go-wallet-sdk/coins/cosmos/okc/tx/common/types"
	"golang.org/x/crypto/sha3"
)

func hashMessage(p []byte) []byte {
	hf := sha3.NewLegacyKeccak256()
	hf.Reset()
	hf.Write(p)
	return hf.Sum(nil)
}

// MakeSignature completes the signature
func MakeSignature(privateKeyHex string, msg authtypes.StdSignMsg) (sig authtypes.StdSignature, err error) {
	m := hashMessage(msg.Bytes())

	pkBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return authtypes.StdSignature{}, err
	}
	ecPriv, ecPub := btcec.PrivKeyFromBytes(pkBytes)

	result := ecdsa.SignCompact(ecPriv, m, false)

	V := result[0]
	R := result[1:33]
	S := result[33:65]
	sigBytes := make([]byte, 0)
	sigBytes = append(sigBytes, R...)
	sigBytes = append(sigBytes, S...)
	sigBytes = append(sigBytes, V-27)

	return authtypes.StdSignature{
		PubKey:    ethsecp256k1.PubKey(ecPub.SerializeCompressed()),
		Signature: sigBytes,
	}, nil
}

func BuildStdTx(privateKey, chainId, memo string, msgs []types.Msg, feeCoins types.Coins, gas, accNumber, seqNumber uint64) (
	stdTx *authtypes.StdTx, err error) {
	if len(chainId) == 0 {
		return stdTx, errors.New("failed. empty chain ID")
	}

	stdFee := authtypes.NewStdFee(gas, feeCoins)

	signMsg := authtypes.StdSignMsg{
		ChainID:       chainId,
		AccountNumber: accNumber,
		Sequence:      seqNumber,
		Memo:          memo,
		Msgs:          msgs,
		Fee:           stdFee,
	}

	sigBytes, err := MakeSignature(privateKey, signMsg)
	if err != nil {
		return
	}

	return authtypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authtypes.StdSignature{sigBytes}, signMsg.Memo), err
}

func MarshalStdTx(stdTx *authtypes.StdTx) (string, error) {
	bytes, err := amino.GCodec.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
