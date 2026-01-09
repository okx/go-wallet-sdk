package cosmos

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okx/go-wallet-sdk/coins/cosmos/tx"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types/ethsecp256k1"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types/ibc"
	"golang.org/x/crypto/sha3"
)

const (
	ATOM_HRP = "cosmos"
)

type CommonParam struct {
	ChainId       string
	Sequence      uint64
	AccountNumber uint64
	FeeDemon      string
	FeeAmount     string
	GasLimit      uint64
	Memo          string
	TimeoutHeight uint64
}

type TransferParam struct {
	CommonParam
	FromAddress string
	ToAddress   string
	Demon       string
	Amount      string
}

type IbcTransferParam struct {
	CommonParam
	FromAddress      string
	ToAddress        string
	Demon            string
	Amount           string
	SourcePort       string
	SourceChannel    string
	TimeOutHeight    ibc.Height
	TimeOutInSeconds uint64
}

func NewAddress(privateKey string, hrp string, followETH bool) (string, error) {
	pkBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	_, pb := btcec.PrivKeyFromBytes(pkBytes)
	if followETH {
		pubBytes := pb.SerializeUncompressed()
		hash := sha3.NewLegacyKeccak256()
		hash.Write(pubBytes[1:])
		addressByte := hash.Sum(nil)
		address, err := bech32.EncodeFromBase256(hrp, addressByte[12:])
		if err != nil {
			return "", err
		}
		return address, nil
	}
	bytes := btcutil.Hash160(pb.SerializeCompressed())
	address, err := bech32.EncodeFromBase256(hrp, bytes)
	if err != nil {
		return "", err
	}
	return address, nil
}

func PubHex2AnyHex(publicKeyHex string, compress bool) (string, error) {
	bb, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", err
	}
	publicKey, err := btcec.ParsePubKey(bb)
	if err != nil {
		return "", err
	}
	var pk []byte
	if compress {
		pk = publicKey.SerializeCompressed()
	} else {
		pk = publicKey.SerializeUncompressed()
	}
	pubkey := types.PubKey{Key: pk}
	anyPubkey, err := types.NewAnyWithValue(&pubkey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(anyPubkey.Value), nil
}

func Convert2AnyPubKey(pubKeyOrHex string, compress, eth bool) (string, error) {
	pk, err := DecodeAnyOrPubKey(pubKeyOrHex)
	if err != nil {
		return "", err
	}
	pubKey, err := btcec.ParsePubKey(pk)
	if err != nil {
		return "", err
	}
	if eth {
		if compress {
			return hex.EncodeToString(pubKey.SerializeCompressed()), nil
		}
		return hex.EncodeToString(pubKey.SerializeUncompressed()), nil
	}
	var b []byte
	if compress {
		b = pubKey.SerializeCompressed()
	} else {
		b = pubKey.SerializeUncompressed()
	}
	pubkey := types.PubKey{Key: b}
	anyPubkey, err := types.NewAnyWithValue(&pubkey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(anyPubkey.Value), nil
}

func DecodeAnyOrPubKey(pubKeyHex string) ([]byte, error) {
	b, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, err
	}
	if len(b) == 33 || len(b) == 65 {
		return b, nil
	}

	pubKey := types.PubKey{}
	if err := pubKey.Unmarshal(b); err != nil {
		return nil, err
	}
	//anyPubKey, _ := types.NewAnyWithValue(&pubKey)
	return pubKey.GetKey(), nil
}

func GetAddressByPublicKey(pubKeyHex string, HRP string) (string, error) {
	pb, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return "", err
	}
	bytes := btcutil.Hash160(pb)
	address, err := bech32.EncodeFromBase256(HRP, bytes)
	if err != nil {
		return "", err
	}
	return address, nil
}
func ValidateAddress(address string, hrp string) bool {
	h, _, err := bech32.DecodeToBase256(address)
	return err == nil && h == hrp
}

func MakeTransactionWithMessage(p CommonParam, publicKeyCompressed string, messages []*types.Any) (string, error) {
	// body = messages(message array) + memo + timeoutHeight
	body := tx.TxBody{Messages: messages, Memo: p.Memo, TimeoutHeight: p.TimeoutHeight}

	// Public key 33bytes compressed format
	publicKeyBytes, err := hex.DecodeString(publicKeyCompressed)
	if err != nil {
		return "", err
	}

	pubkey := types.PubKey{Key: publicKeyBytes}
	anyPubkey, err := types.NewAnyWithValue(&pubkey)
	if err != nil {
		return "", err
	}

	single := tx.ModeInfo_Single{Mode: types.SignMode_SIGN_MODE_DIRECT}
	single_ := tx.ModeInfo_Single_{Single: &single}
	modeInfo := tx.ModeInfo{Sum: &single_}
	signerInfo := make([]*tx.SignerInfo, 0)
	signerInfo = append(signerInfo, &tx.SignerInfo{PublicKey: anyPubkey, ModeInfo: &modeInfo, Sequence: p.Sequence})

	// authInfo = signerInfo(publicKey + modeInfo + sequence) + fee(amount + gasLimit)
	feeAmount, ok := types.NewIntFromString(p.FeeAmount)
	if !ok {
		return "", errors.New("invalid fee amount")
	}

	feeCoin := types.NewCoin(p.FeeDemon, feeAmount)
	feeCoins := types.NewCoins(feeCoin)
	fee := tx.Fee{Amount: feeCoins, GasLimit: p.GasLimit}
	authInfo := tx.AuthInfo{SignerInfos: signerInfo, Fee: &fee}

	// signDoc = bodyBytes + authInfoBytes + ChainId + AccountNumber
	bodyBytes, err := body.Marshal()
	if err != nil {
		return "", err
	}

	authInfoBytes, err := authInfo.Marshal()
	if err != nil {
		return "", err
	}

	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: p.ChainId, AccountNumber: p.AccountNumber}
	signDocBtyes, err := signDoc.Marshal()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signDocBtyes), nil
}

func MakeTransactionWithSignDoc(body string, auth string, ChainId string, AccountNumber uint64) (string, error) {
	bodyBytes, err := hexutil.Decode(body)
	if err != nil {
		return "", err
	}

	authInfoBytes, err := hexutil.Decode(auth)
	if err != nil {
		return "", err
	}

	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: ChainId, AccountNumber: AccountNumber}
	signDocBtyes, err := signDoc.Marshal()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signDocBtyes), nil
}

// GetRawTransaction Gets the string to be signed
// param - （TransferParam | IbcTransferParam）
// publicKeyCompressed - （hex 33 bytes）
// return  signDoc
func GetRawTransaction(param interface{}, publicKeyCompressed string) (string, error) {
	switch param.(type) {
	case TransferParam:
		{
			p, _ := param.(TransferParam)
			amount, ok := types.NewIntFromString(p.Amount)
			if !ok {
				return "", errors.New("invalid  amount")
			}
			coin := types.NewCoin(p.Demon, amount)
			coins := types.NewCoins(coin)
			sendMsg := types.MsgSend{FromAddress: p.FromAddress, ToAddress: p.ToAddress, Amount: coins}

			messages := make([]*types.Any, 0)
			anySend, err := types.NewAnyWithValue(&sendMsg)
			if err != nil {
				return "", err
			}
			messages = append(messages, anySend)
			return MakeTransactionWithMessage(p.CommonParam, publicKeyCompressed, messages)
		}
	case IbcTransferParam:
		p, _ := param.(IbcTransferParam)
		amount, ok := types.NewIntFromString(p.Amount)
		if !ok {
			return "", errors.New("invalid  amount")
		}
		coin := types.NewCoin(p.Demon, amount)
		sendMsg := ibc.MsgTransfer{
			SourcePort:       p.SourcePort,
			SourceChannel:    p.SourceChannel,
			Token:            coin,
			Sender:           p.FromAddress,
			Receiver:         p.ToAddress,
			TimeoutHeight:    p.TimeOutHeight,
			TimeoutTimestamp: p.TimeOutInSeconds * 1_000_000_000,
		}
		messages := make([]*types.Any, 0)
		anySend, err := types.NewAnyWithValue(&sendMsg)
		if err != nil {
			return "", err
		}
		messages = append(messages, anySend)
		return MakeTransactionWithMessage(p.CommonParam, publicKeyCompressed, messages)
	default:
		return "", errors.New("unsupported param type")
	}
}

// GetSignedTransaction - Getting signed transactions
// signDoc - Transaction string to be signed（hex）
// signature - Transaction signature (hex)
func GetSignedTransaction(signDoc string, signature string) (string, error) {
	signDocBytes, err := hex.DecodeString(signDoc)
	if err != nil {
		return "", err
	}

	doc := tx.SignDoc{}
	err = doc.Unmarshal(signDocBytes)
	if err != nil {
		return "", err
	}

	signBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}

	signatures := make([][]byte, 0)
	signatures = append(signatures, signBytes)

	trans := tx.TxRaw{BodyBytes: doc.BodyBytes, AuthInfoBytes: doc.AuthInfoBytes, Signatures: signatures}
	transBytes, err := trans.Marshal()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(transBytes), nil
}

func SignRawTransaction(signDoc string, privateKey *btcec.PrivateKey) (string, error) {
	signDocBtyes, err := hex.DecodeString(signDoc)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(signDocBtyes)
	signature := ecdsa.SignCompact(privateKey, hash[:], false)
	return hex.EncodeToString(signature[1:]), nil
}

func Transfer(param TransferParam, privateKeyHex string) (string, error) {
	return TransferAction(param, privateKeyHex, false)
}

func TransferAction(param TransferParam, privateKeyHex string, useEthSecp256k1 bool) (string, error) {
	amount, ok := types.NewIntFromString(param.Amount)
	if !ok {
		return "", errors.New("invalid  amount")
	}
	coin := types.NewCoin(param.Demon, amount)
	coins := types.NewCoins(coin)
	sendMsg := types.MsgSend{FromAddress: param.FromAddress, ToAddress: param.ToAddress, Amount: coins}

	messages := make([]*types.Any, 0)
	anySend, err := types.NewAnyWithValue(&sendMsg)
	if err != nil {
		return "", err
	}
	messages = append(messages, anySend)
	return BuildTxAction(param.CommonParam, messages, privateKeyHex, useEthSecp256k1)
}

func IbcTransfer(param IbcTransferParam, privateKeyHex string) (string, error) {
	return IbcTransferAction(param, privateKeyHex, false)
}

func IbcTransferAction(param IbcTransferParam, privateKeyHex string, useEthSecp256k1 bool) (string, error) {
	amount, ok := types.NewIntFromString(param.Amount)
	if !ok {
		return "", errors.New("invalid  amount")
	}
	coin := types.NewCoin(param.Demon, amount)
	sendMsg := ibc.MsgTransfer{
		SourcePort:       param.SourcePort,
		SourceChannel:    param.SourceChannel,
		Token:            coin,
		Sender:           param.FromAddress,
		Receiver:         param.ToAddress,
		TimeoutHeight:    param.TimeOutHeight,
		TimeoutTimestamp: param.TimeOutInSeconds * 1_000_000_000,
	}

	messages := make([]*types.Any, 0)
	anySend, err := types.NewAnyWithValue(&sendMsg)
	if err != nil {
		return "", err
	}
	messages = append(messages, anySend)
	return BuildTxAction(param.CommonParam, messages, privateKeyHex, useEthSecp256k1)
}

func BuildTx(param CommonParam, messages []*types.Any, privateKeyHex string) (string, error) {
	return BuildTxAction(param, messages, privateKeyHex, false)
}

func getPublicKey(privateKeyHex string, useEthSecp256k1 bool) (*types.Any, error) {
	pkBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	_, publicKey := btcec.PrivKeyFromBytes(pkBytes)
	if useEthSecp256k1 {
		pubkey := ethsecp256k1.PubKey{Key: publicKey.SerializeCompressed()}
		anyPubkey, err := types.NewAnyWithValue(&pubkey)
		if err != nil {
			return nil, err
		}
		return anyPubkey, nil
	} else {
		pubkey := types.PubKey{Key: publicKey.SerializeCompressed()}
		anyPubkey, err := types.NewAnyWithValue(&pubkey)
		if err != nil {
			return nil, err
		}
		return anyPubkey, nil
	}
}

func HashMessage(p []byte) []byte {
	hf := sha3.NewLegacyKeccak256()
	hf.Reset()
	hf.Write(p)
	return hf.Sum(nil)
}

func BuildTxAction(param CommonParam, messages []*types.Any, privateKeyHex string, useEthSecp256k1 bool) (string, error) {
	// body = messages(message array) + memo + timeoutHeight
	body := tx.TxBody{Messages: messages, Memo: param.Memo, TimeoutHeight: param.TimeoutHeight}

	// Public key 33bytes compressed format
	pkBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(pkBytes)

	pubKey, err := getPublicKey(privateKeyHex, useEthSecp256k1)
	if err != nil {
		return "", err
	}
	single := tx.ModeInfo_Single{Mode: types.SignMode_SIGN_MODE_DIRECT}
	single_ := tx.ModeInfo_Single_{Single: &single}
	modeInfo := tx.ModeInfo{Sum: &single_}
	signerInfo := make([]*tx.SignerInfo, 0)
	signerInfo = append(signerInfo, &tx.SignerInfo{PublicKey: pubKey, ModeInfo: &modeInfo, Sequence: param.Sequence})

	// authInfo = signerInfo(publicKey + modeInfo + sequence) + fee(amount + gasLimit)
	feeAmount, ok := types.NewIntFromString(param.FeeAmount)
	if !ok {
		return "", errors.New("invalid  fee amount")
	}
	feeCoin := types.NewCoin(param.FeeDemon, feeAmount)
	feeCoins := types.NewCoins(feeCoin)
	fee := tx.Fee{Amount: feeCoins, GasLimit: param.GasLimit}
	authInfo := tx.AuthInfo{SignerInfos: signerInfo, Fee: &fee}

	// signDoc = bodyBytes + authInfoBytes + ChainId + AccountNumber
	bodyBytes, err := body.Marshal()
	if err != nil {
		return "", err
	}
	authInfoBytes, err := authInfo.Marshal()
	if err != nil {
		return "", err
	}
	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: param.ChainId, AccountNumber: param.AccountNumber}
	signDocBtyes, err := signDoc.Marshal()
	if err != nil {
		return "", err
	}

	var signBytes []byte
	if useEthSecp256k1 {
		m := HashMessage(signDocBtyes)
		result := ecdsa.SignCompact(privateKey, m, false)
		V := result[0]
		R := result[1:33]
		S := result[33:65]
		signBytes = make([]byte, 0)
		signBytes = append(signBytes, R...)
		signBytes = append(signBytes, S...)
		signBytes = append(signBytes, V-27)
	} else {
		hash := sha256.Sum256(signDocBtyes)
		signBytes = ecdsa.SignCompact(privateKey, hash[:], false)
		signBytes = signBytes[1:]
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

func BuildTxActionForSignMessage(param CommonParam, messages []*types.Any, privateKeyHex string, useEthSecp256k1 bool) (string, string, error) {
	// body = messages(message array) + memo + timeoutHeight
	body := tx.TxBody{Messages: messages, Memo: param.Memo, TimeoutHeight: param.TimeoutHeight}

	// Public key 33bytes compressed format
	pkBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", "", err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(pkBytes)

	pubKey, err := getPublicKey(privateKeyHex, useEthSecp256k1)
	if err != nil {
		return "", "", err
	}

	single := tx.ModeInfo_Single{Mode: types.SignMode_SIGN_MODE_DIRECT}
	single_ := tx.ModeInfo_Single_{Single: &single}
	modeInfo := tx.ModeInfo{Sum: &single_}
	signerInfo := make([]*tx.SignerInfo, 0)
	signerInfo = append(signerInfo, &tx.SignerInfo{PublicKey: pubKey, ModeInfo: &modeInfo, Sequence: param.Sequence})

	// authInfo = signerInfo(publicKey + modeInfo + sequence) + fee(amount + gasLimit)
	feeAmount, _ := types.NewIntFromString(param.FeeAmount)
	feeCoin := types.NewCoin(param.FeeDemon, feeAmount)
	feeCoins := types.NewCoins(feeCoin)
	fee := tx.Fee{Amount: feeCoins, GasLimit: param.GasLimit}
	authInfo := tx.AuthInfo{SignerInfos: signerInfo, Fee: &fee}

	// signDoc = bodyBytes + authInfoBytes + ChainId + AccountNumber
	bodyBytes, _ := body.Marshal()
	authInfoBytes, _ := authInfo.Marshal()
	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: param.ChainId, AccountNumber: param.AccountNumber}
	signDocBtyes, _ := signDoc.Marshal()

	// signature
	var signBytes []byte
	if useEthSecp256k1 {
		m := HashMessage(signDocBtyes)
		result := ecdsa.SignCompact(privateKey, m, false)
		V := result[0]
		R := result[1:33]
		S := result[33:65]
		signBytes = make([]byte, 0)
		signBytes = append(signBytes, R...)
		signBytes = append(signBytes, S...)
		signBytes = append(signBytes, V-27)
	} else {
		hash := sha256.Sum256(signDocBtyes)
		signBytes = ecdsa.SignCompact(privateKey, hash[:], false)
		signBytes = signBytes[1:]
	}

	signatures := make([][]byte, 0)
	signatures = append(signatures, signBytes)

	trans := tx.TxRaw{BodyBytes: signDoc.BodyBytes, AuthInfoBytes: signDoc.AuthInfoBytes, Signatures: signatures}
	transBytes, err := trans.Marshal()
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(transBytes), base64.StdEncoding.EncodeToString(signBytes), nil
}

type MessageData struct {
	ChainId       string `json:"chain_id,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	Sequence      string `json:"sequence,omitempty"`
	Fee           struct {
		Gas    string       `json:"gas,omitempty"`
		Amount []types.Coin `json:"amount,omitempty"`
	} `json:"fee,omitempty"`
	Msgs []struct {
		T string      `json:"type,omitempty"`
		V interface{} `json:"value,omitempty"`
	} `json:"msgs,omitempty"`
	Memo string `json:"memo,omitempty"`
}

func SignDoc(body string, auth string, privateKeyHex string, ChainId string, AccountNumber uint64) (string, string, error) {
	bodyBytes, err := hexutil.Decode(body)
	if err != nil {
		return "", "", err
	}

	authInfoBytes, err := hexutil.Decode(auth)
	if err != nil {
		return "", "", err
	}

	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: ChainId, AccountNumber: AccountNumber}
	signDocBtyes, err := signDoc.Marshal()
	if err != nil {
		return "", "", err
	}

	pkBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", "", err
	}

	privateKey, _ := btcec.PrivKeyFromBytes(pkBytes)

	var signBytes []byte
	hash := sha256.Sum256(signDocBtyes)
	signature := ecdsa.SignCompact(privateKey, hash[:], false)
	signBytes = signature[1:]

	signatures := make([][]byte, 0)
	signatures = append(signatures, signBytes)

	trans := tx.TxRaw{BodyBytes: signDoc.BodyBytes, AuthInfoBytes: signDoc.AuthInfoBytes, Signatures: signatures}
	transBytes, err := trans.Marshal()
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(transBytes), base64.StdEncoding.EncodeToString(signBytes), nil
}

func sortedObject(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		result := make(map[string]interface{})
		for _, key := range keys {
			result[key] = sortedObject(v[key])
		}
		return result
	case []interface{}:
		for i, item := range v {
			v[i] = sortedObject(item)
		}
		return v
	default:
		return obj
	}
}

func SignAminoMessage(data string, privateKeyHex string) (string, error) {
	var msg map[string]interface{}
	json.Unmarshal([]byte(data), &msg)

	sortedMsg := sortedObject(msg)
	msgBytes, err := json.Marshal(sortedMsg)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(msgBytes)
	pkBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}

	privateKey, _ := btcec.PrivKeyFromBytes(pkBytes)
	signature := ecdsa.SignCompact(privateKey, hash[:], false)
	return base64.StdEncoding.EncodeToString(signature[1:]), nil
}

// SignMessage Sign the message, using JSON format
func SignMessage(data string, privateKeyHex string) (string, string, error) {
	tx, _, err := SignMessageAction(data, privateKeyHex, false)
	if err != nil {
		return "", "", err
	}

	sig, err := SignAminoMessage(data, privateKeyHex)
	if err != nil {
		return "", "", err
	}
	return tx, sig, nil
}

func SignMessageAction(data string, privateKeyHex string, useEthSecp256k1 bool) (string, string, error) {
	messageData := MessageData{}
	_ = json.Unmarshal([]byte(data), &messageData)

	messages := make([]*types.Any, 0)
	for _, m := range messageData.Msgs {
		fn := types.GetMessageConverter(m.T)
		if fn != nil {
			b, err := json.Marshal(m.V)
			if err != nil {
				return "", "", err
			}
			messages = append(messages, fn(string(b)))
		}
	}

	param := CommonParam{}
	param.AccountNumber = string2Uint64(messageData.AccountNumber)
	param.Sequence = string2Uint64(messageData.Sequence)
	param.Memo = messageData.Memo
	param.TimeoutHeight = 0
	param.FeeDemon = messageData.Fee.Amount[0].Denom
	param.FeeAmount = messageData.Fee.Amount[0].Amount.String()
	param.GasLimit = string2Uint64(messageData.Fee.Gas)
	param.ChainId = messageData.ChainId
	return BuildTxActionForSignMessage(param, messages, privateKeyHex, useEthSecp256k1)
}

func string2Uint64(intStr string) uint64 {
	i := big.Int{}
	_, _ = i.SetString(intStr, 10)
	return i.Uint64()
}

func getJsonSignDoc(p *CommonParam, msg *types.StdAny) (*types.StdSignDoc, error) {
	signDoc := types.StdSignDoc{}
	signDoc.AccountNumber = strconv.FormatUint(p.AccountNumber, 10)
	signDoc.Sequence = strconv.FormatUint(p.Sequence, 10)
	signDoc.ChainID = p.ChainId
	signDoc.Memo = p.Memo
	if p.TimeoutHeight != 0 {
		signDoc.TimeoutHeight = strconv.FormatUint(p.TimeoutHeight, 10)
	}

	signDoc.Fee = types.StdFee{}
	signDoc.Fee.Gas = strconv.FormatUint(p.GasLimit, 10)
	feeAmount, ok := types.NewIntFromString(p.FeeAmount)
	if !ok {
		return nil, errors.New("invalid fee amount")
	}
	feeCoin := types.NewCoin(p.FeeDemon, feeAmount)
	feeCoins := types.NewCoins(feeCoin)
	signDoc.Fee.Amount = feeCoins

	signDoc.Msgs = make([]types.StdAny, 0)
	signDoc.Msgs = append(signDoc.Msgs, *msg)
	return &signDoc, nil
}

func GetRawJsonTransaction(param interface{}) (string, error) {
	switch param.(type) {
	case TransferParam:
		{
			p, _ := param.(TransferParam)
			amount, ok := types.NewIntFromString(p.Amount)
			if !ok {
				return "", errors.New("invalid  amount")
			}
			coin := types.NewCoin(p.Demon, amount)
			coins := types.NewCoins(coin)
			sendMsg := types.MsgSend{FromAddress: p.FromAddress, ToAddress: p.ToAddress, Amount: coins}
			signDoc, err := getJsonSignDoc(&p.CommonParam, &types.StdAny{T: "cosmos-sdk/MsgSend", V: sendMsg})
			if err != nil {
				return "", err
			}
			bytes, err := json.Marshal(signDoc)
			if err != nil {
				return "", err
			}
			return string(types.MustSortJSON(bytes)), nil
		}
	case IbcTransferParam:
		{
			p, _ := param.(IbcTransferParam)
			amount, ok := types.NewIntFromString(p.Amount)
			if !ok {
				return "", errors.New("invalid  amount")
			}
			coin := types.NewCoin(p.Demon, amount)
			sendMsg := ibc.MsgTransfer{
				SourcePort:       p.SourcePort,
				SourceChannel:    p.SourceChannel,
				Token:            coin,
				Sender:           p.FromAddress,
				Receiver:         p.ToAddress,
				TimeoutHeight:    p.TimeOutHeight,
				TimeoutTimestamp: p.TimeOutInSeconds * 1_000_000_000,
			}
			signDoc, err := getJsonSignDoc(&p.CommonParam, &types.StdAny{T: "cosmos-sdk/MsgTransfer", V: sendMsg})
			if err != nil {
				return "", err
			}
			bytes, err := json.Marshal(signDoc)
			if err != nil {
				return "", err
			}
			return string(types.MustSortJSON(bytes)), nil
		}
	default:
		return "", fmt.Errorf("unspport param type")
	}
}

func GetSignedJsonTransaction(signDoc string, publicKey string, signature string) (string, error) {
	stdDoc := types.StdSignDoc{}
	err := json.Unmarshal([]byte(signDoc), &stdDoc)
	if err != nil {
		return "", err
	}

	messages := make([]*types.Any, 0)
	for _, msg := range stdDoc.Msgs {
		if msg.T == "cosmos-sdk/MsgSend" {
			mBytes, err := json.Marshal(msg.V)
			if err != nil {
				return "", err
			}
			sendMsg := types.MsgSend{}
			err = json.Unmarshal(mBytes, &sendMsg)
			if err != nil {
				return "", err
			}
			anySend, err := types.NewAnyWithValue(&sendMsg)
			if err != nil {
				return "", err
			}
			messages = append(messages, anySend)
		} else {
			mBytes, err := json.Marshal(msg.V)
			if err != nil {
				return "", err
			}
			sendMsg := ibc.MsgTransfer{}
			err = json.Unmarshal(mBytes, &sendMsg)
			if err != nil {
				return "", err
			}
			anySend, err := types.NewAnyWithValue(&sendMsg)
			if err != nil {
				return "", err
			}
			messages = append(messages, anySend)
		}
	}

	compressedBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	pubKey := types.PubKey{Key: compressedBytes}
	anyPubkey, err := types.NewAnyWithValue(&pubKey)
	if err != nil {
		return "", err
	}

	single := tx.ModeInfo_Single{Mode: types.SignMode_SIGN_MODE_LEGACY_AMINO_JSON}
	single_ := tx.ModeInfo_Single_{Single: &single}
	modeInfo := tx.ModeInfo{Sum: &single_}
	signerInfo := make([]*tx.SignerInfo, 0)

	signerInfo = append(signerInfo, &tx.SignerInfo{PublicKey: anyPubkey, ModeInfo: &modeInfo, Sequence: string2Uint64(stdDoc.Sequence)})

	feeCoins := make([]types.Coin, 0)
	for _, coin := range stdDoc.Fee.Amount {
		demon := coin.Denom
		amount := coin.Amount
		feeCoin := types.NewCoin(demon, amount)
		feeCoins = append(feeCoins, feeCoin)
	}
	fee := tx.Fee{Amount: feeCoins, GasLimit: string2Uint64(stdDoc.Fee.Gas)}

	body := tx.TxBody{Messages: messages, Memo: stdDoc.Memo, TimeoutHeight: string2Uint64(stdDoc.TimeoutHeight)}
	bodyBytes, err := body.Marshal()
	if err != nil {
		return "", err
	}

	authInfo := tx.AuthInfo{SignerInfos: signerInfo, Fee: &fee}
	authInfoBytes, err := authInfo.Marshal()
	if err != nil {
		return "", err
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}
	signatures := make([][]byte, 0)
	signatures = append(signatures, signatureBytes)

	trans := tx.TxRaw{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, Signatures: signatures}
	transBytes, err := trans.Marshal()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(transBytes), nil
}
func GetSigningHash(rawTxByte string) (string, error) {
	txHashByte := sha256.Sum256([]byte(rawTxByte))
	txHashHex := hex.EncodeToString(txHashByte[:])
	return txHashHex, nil
}
func SignRawJsonTransaction(signDoc string, privateKey *btcec.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(signDoc))
	signature := ecdsa.SignCompact(privateKey, hash[:], false)
	return hex.EncodeToString(signature[1:]), nil
}
