package zil

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/emresenyuva/go-wallet-sdk/coins/zil/keytools"
	"github.com/emresenyuva/go-wallet-sdk/coins/zil/protobuf"
	"github.com/emresenyuva/go-wallet-sdk/coins/zil/util"
	"google.golang.org/protobuf/proto"
	"math/big"
	"strconv"
	"strings"
)

var bintZero = big.NewInt(0)

type Transaction struct {
	ID              string
	Version         string
	Nonce           string
	Amount          string
	GasPrice        string
	GasLimit        string
	Signature       string
	SenderPubKey    string
	ToAddr          string
	Code            string
	Data            interface{}
	ContractAddress string
	Priority        bool
}

type txParams struct {
	ID           string `json:"ID"`
	Version      string `json:"version"`
	Nonce        string `json:"nonce"`
	Amount       string `json:"amount"`
	GasPrice     string `json:"gasPrice"`
	GasLimit     string `json:"gasLimit"`
	Signature    string `json:"signature"`
	SenderPubKey string `json:"senderPubKey"`
	ToAddr       string `json:"toAddr"`
	Code         string `json:"code"`
	Data         string `json:"data"`
}

type TransactionPayload struct {
	Version   int    `json:"version"`
	Nonce     int    `json:"nonce"`
	ToAddr    string `json:"toAddr"`
	Amount    string `json:"amount"`
	PubKey    string `json:"pubKey"`
	GasPrice  string `json:"gasPrice"`
	GasLimit  string `json:"gasLimit"`
	Code      string `json:"code"`
	Data      string `json:"data"`
	Signature string `json:"signature"`
	Priority  bool   `json:"priority"`
}

type Init struct {
	Version   int           `json:"version"`
	Nonce     int           `json:"nonce"`
	ToAddr    string        `json:"toAddr"`
	Amount    int64         `json:"amount"`
	PubKey    string        `json:"pubKey"`
	GasPrice  int64         `json:"gasPrice"`
	GasLimit  int64         `json:"gasLimit"`
	Code      string        `json:"code"`
	Data      []interface{} `json:"data"`
	Signature string        `json:"signature"`
}

func (t *Transaction) toTransactionParam() txParams {
	data, _ := json.Marshal(t.Data)
	param := txParams{
		ID:           t.ID,
		Version:      t.Version,
		Nonce:        t.Nonce,
		Amount:       t.Amount,
		GasPrice:     t.GasPrice,
		GasLimit:     t.GasLimit,
		Signature:    t.Signature,
		SenderPubKey: t.SenderPubKey,
		Code:         t.Code,
		Data:         string(data),
	}

	if t.ToAddr == "" {
		param.ToAddr = "0000000000000000000000000000000000000000"
	} else {
		param.ToAddr = t.ToAddr
	}
	return param
}

func EncodeTransactionProto(txParams txParams) ([]byte, error) {
	amount, ok := new(big.Int).SetString(txParams.Amount, 10)
	if !ok {
		return nil, errors.New("amount error")
	}

	gasPrice, ok2 := new(big.Int).SetString(txParams.GasPrice, 10)
	if !ok2 {
		return nil, errors.New("gas price error")
	}

	v, err := strconv.ParseUint(txParams.Version, 10, 32)
	if err != nil {
		return nil, err
	}
	version := uint32(v)

	nonce, err2 := strconv.ParseUint(txParams.Nonce, 10, 64)
	if err2 != nil {
		return nil, err2
	}

	pubKeyBytes, _ := hex.DecodeString(txParams.SenderPubKey)
	senderpubkey := protobuf.ByteArray{
		Data: pubKeyBytes,
	}

	amountArray := protobuf.ByteArray{
		Data: bigIntToPaddedBytes(amount, 32),
	}

	gasPriceArray := protobuf.ByteArray{
		Data: bigIntToPaddedBytes(gasPrice, 32),
	}

	gasLimit, err3 := strconv.ParseUint(txParams.GasLimit, 10, 64)
	if err3 != nil {
		return nil, err3
	}

	addrBytes, _ := hex.DecodeString(txParams.ToAddr)
	protoTransactionCoreInfo := protobuf.ProtoTransactionCoreInfo{
		Version:      &version,
		Nonce:        &nonce,
		Toaddr:       addrBytes,
		Senderpubkey: &senderpubkey,
		Amount:       &amountArray,
		Gasprice:     &gasPriceArray,
		Gaslimit:     &gasLimit,
	}

	if txParams.Data == "\"\"" {
		txParams.Data = ""
	}

	if txParams.Data != "" {
		protoTransactionCoreInfo.Data = []byte(txParams.Data)
	}

	if txParams.Code != "" {
		protoTransactionCoreInfo.Code = []byte(txParams.Code)
	}

	bytes, err4 := proto.Marshal(&protoTransactionCoreInfo)
	if err4 != nil {
		return nil, err4
	}
	return bytes, nil

}

func bigIntToPaddedBytes(i *big.Int, paddedSize int32) []byte {
	bytes := i.Bytes()
	padded, _ := hex.DecodeString(fmt.Sprintf("%0*x", paddedSize, bytes))
	return padded
}

func (t *Transaction) Bytes() ([]byte, error) {
	txParams := t.toTransactionParam()
	bytes, err := EncodeTransactionProto(txParams)

	if err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func pack(a int, b int) int {
	return a<<16 + b
}

func CreateTransferTransaction(to, gasPrice string, amount, gasLimit *big.Int, nonce int, chainId int) *Transaction {
	tx := &Transaction{
		Version:  strconv.FormatInt(int64(pack(chainId, 1)), 10),
		ToAddr:   to,
		Amount:   amount.String(),
		GasPrice: gasPrice,
		GasLimit: gasLimit.String(),
		Nonce:    strconv.Itoa(nonce),
		Code:     "",
		Data:     "",
		Priority: false,
	}
	return tx
}

func SignTransaction(privateKeyhex string, tx *Transaction) error {

	if !IsBech32(tx.ToAddr) {
		return errors.New("not bech32")
	}

	address, err := FromBech32Addr(tx.ToAddr)
	if err != nil {
		return err
	}
	tx.ToAddr = address

	publicKeyHex, err := GetPublicKeyFromPrivateKey(privateKeyhex)
	if err != nil {
		return err
	}
	tx.SenderPubKey = publicKeyHex

	message, err := tx.Bytes()
	if err != nil {
		return err
	}

	privBytes, err := hex.DecodeString(privateKeyhex)
	if err != nil {
		return err
	}

	pubBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return err
	}

	rb, err := keytools.GenerateRandomBytes(btcec.S256().N.BitLen() / 8)

	if err != nil {
		return err
	}

	r, s, err := trySign(privBytes, pubBytes, message, rb)
	if err != nil {
		return err
	}
	sig := fmt.Sprintf("%064s%064s", util.EncodeHex(r), util.EncodeHex(s))
	tx.Signature = sig

	return nil
}

func trySign(privateKey []byte, publicKey []byte, message []byte, k []byte) ([]byte, []byte, error) {
	priKey := new(big.Int).SetBytes(privateKey)
	bintK := new(big.Int).SetBytes(k)

	curve := btcec.S256()

	// 1a. check if private key is 0
	if priKey.Cmp(new(big.Int).SetInt64(0)) <= 0 {
		return nil, nil, errors.New("private key must be > 0")
	}

	// 1b. check if private key is less than curve order, i.e., within [1...n-1]
	if priKey.Cmp(curve.N) >= 0 {
		return nil, nil, errors.New("private key cannot be greater than curve order")
	}

	if bintK.Cmp(bintZero) == 0 {
		return nil, nil, errors.New("k cannot be zero")
	}

	if bintK.Cmp(curve.N) > 0 {
		return nil, nil, errors.New("k cannot be greater than order of secp256k1")
	}

	// 2. Compute commitment Q = kG, where G is the base point
	Qx, Qy := curve.ScalarBaseMult(k)

	Q := util.Compress(curve, Qx, Qy, true)

	// 3. Compute the challenge r = H(Q || pubKey || msg)
	// mod reduce r by the order of secp256k1, n
	r := new(big.Int).SetBytes(hash(Q, publicKey, message[:]))
	r = r.Mod(r, curve.N)

	if r.Cmp(bintZero) == 0 {
		return nil, nil, errors.New("invalid r")
	}

	//4. Compute s = k - r * prv
	// 4a. Compute r * prv
	_r := *r
	s := new(big.Int).Mod(_r.Mul(&_r, priKey), curve.N)
	s = new(big.Int).Mod(new(big.Int).Sub(bintK, s), curve.N)

	if s.Cmp(big.NewInt(0)) == 0 {
		return nil, nil, errors.New("invalid s")
	}

	return r.Bytes(), s.Bytes(), nil
}

func hash(Q []byte, pubKey []byte, msg []byte) []byte {
	var buffer bytes.Buffer
	buffer.Write(Q)
	buffer.Write(pubKey[:33])
	buffer.Write(msg)
	return util.Sha256(buffer.Bytes())
}

func (t *Transaction) ToTransactionPayload() TransactionPayload {
	version, _ := strconv.ParseInt(t.Version, 10, 32)
	nonce, _ := strconv.ParseInt(t.Nonce, 10, 32)
	data, _ := json.Marshal(t.Data)

	p := TransactionPayload{
		Version:   int(version),
		Nonce:     int(nonce),
		ToAddr:    util.ToCheckSumAddress(t.ToAddr)[2:],
		Amount:    t.Amount,
		PubKey:    strings.ToLower(t.SenderPubKey),
		GasPrice:  t.GasPrice,
		GasLimit:  t.GasLimit,
		Code:      t.Code,
		Signature: strings.ToLower(t.Signature),
		Priority:  t.Priority,
	}

	if string(data) != "\"\"" {
		p.Data = string(data)
	}

	if p.ToAddr == "0000000000000000000000000000000000000000" {
		p.ToAddr = "0x0000000000000000000000000000000000000000"
	}
	return p
}
