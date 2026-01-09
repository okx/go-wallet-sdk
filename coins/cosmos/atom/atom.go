package atom

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/okx/go-wallet-sdk/coins/cosmos/tx"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types"
)

const (
	HRP = "cosmos"
)

func NewAddress(privateKey *btcec.PrivateKey) (string, error) {
	pb := privateKey.PubKey().SerializeCompressed()
	bytes := btcutil.Hash160(pb)
	address, err := bech32.EncodeFromBase256(HRP, bytes)
	if err != nil {
		return "", err
	}
	return address, nil
}

func ValidateAddress(address string) bool {
	hrp, _, err := bech32.DecodeToBase256(address)
	return err == nil && hrp == HRP
}

func SignStart(chainId string, from string, to string, demon string, memo string,
	amount *big.Int, timeoutHeight uint64, sequence uint64, accountNumber uint64, feeAmount *big.Int, gasLimit uint64, privateKey *btcec.PrivateKey) (string, error) {
	coin := types.NewCoin(demon, types.NewIntFromBigInt(amount))
	coins := types.NewCoins(coin)
	sendMsg := types.MsgSend{FromAddress: from, ToAddress: to, Amount: coins}

	messages := make([]*types.Any, 0)
	anySend, err := types.NewAnyWithValue(&sendMsg)
	if err != nil {
		return "", err
	}
	messages = append(messages, anySend)

	body := tx.TxBody{Messages: messages, Memo: memo, TimeoutHeight: timeoutHeight}

	// Public key 33bytes compressed format
	publickey := privateKey.PubKey().SerializeCompressed()

	pubkey := types.PubKey{Key: publickey}
	anyPubkey, err := types.NewAnyWithValue(&pubkey)
	if err != nil {
		return "", err
	}

	single := tx.ModeInfo_Single{Mode: types.SignMode_SIGN_MODE_DIRECT}
	single_ := tx.ModeInfo_Single_{Single: &single}
	modeInfo := tx.ModeInfo{Sum: &single_}
	signerInfo := make([]*tx.SignerInfo, 0)
	signerInfo = append(signerInfo, &tx.SignerInfo{PublicKey: anyPubkey, ModeInfo: &modeInfo, Sequence: sequence})

	feeCoin := types.NewCoin(demon, types.NewIntFromBigInt(feeAmount))
	feeCoins := types.NewCoins(feeCoin)
	fee := tx.Fee{Amount: feeCoins, GasLimit: gasLimit}
	authInfo := tx.AuthInfo{SignerInfos: signerInfo, Fee: &fee}

	bodyBytes, err := body.Marshal()
	if err != nil {
		return "", err
	}
	authInfoBytes, err := authInfo.Marshal()
	if err != nil {
		return "", err
	}
	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: chainId, AccountNumber: accountNumber}
	signDocBtyes, err := signDoc.Marshal()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signDocBtyes), nil
}

func Sign(rawHex string, privateKey *btcec.PrivateKey) (string, error) {
	signDocBytes, err := hex.DecodeString(rawHex)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(signDocBytes)
	signature := ecdsa.SignCompact(privateKey, hash[:], false)
	return hex.EncodeToString(signature[1:]), nil
}

func SignEnd(rawHex string, signHex string) (string, error) {
	signDocBytes, err := hex.DecodeString(rawHex)
	if err != nil {
		return "", err
	}
	var signDoc tx.SignDoc
	signDoc.Unmarshal(signDocBytes)

	signBytes, err := hex.DecodeString(signHex)
	if err != nil {
		return "", err
	}
	signatures := make([][]byte, 0)
	signatures = append(signatures, signBytes)

	trans := tx.TxRaw{BodyBytes: signDoc.BodyBytes, AuthInfoBytes: signDoc.AuthInfoBytes, Signatures: signatures}
	transBytes, err := trans.Marshal()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(transBytes), nil
}
