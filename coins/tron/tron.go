package tron

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/sha3"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/okx/go-wallet-sdk/coins/tron/pb"
	"github.com/okx/go-wallet-sdk/coins/tron/token"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

func GetAddress(publicKey *btcec.PublicKey) string {
	pubKey := publicKey.SerializeUncompressed()
	h := sha3.NewLegacyKeccak256()
	h.Write(pubKey[1:])
	hash := h.Sum(nil)[12:]
	return base58.CheckEncode(hash, GetNetWork()[0])
}

func GetAddressByPublicKey(pubKey string) (string, error) {
	pubKeyByte, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key")
	}

	pk, err := btcec.ParsePubKey(pubKeyByte)
	uncompressedPubKey := pk.SerializeUncompressed()
	if err != nil {
		return "", fmt.Errorf("pubKey encoding err ")
	}

	h := sha3.NewLegacyKeccak256()
	h.Write(uncompressedPubKey[1:])
	hash := h.Sum(nil)[12:]
	return base58.CheckEncode(hash, GetNetWork()[0]), nil
}

func ValidateAddress(address string) bool {
	_, _, err := base58.CheckDecode(address)
	return err == nil
}

func GetAddressHash(address string) ([]byte, error) {
	to, v, err := base58.CheckDecode(address)
	if err != nil {
		return nil, err
	}
	var bs []byte
	bs = append(bs, v)
	bs = append(bs, to...)
	return bs, nil
}

func GetNetWork() []byte {
	return []byte{0x41}
}

func ParseTxStr(txStr string) (pb.Transaction, error) {
	bytes, err := hex.DecodeString(txStr)
	var trans pb.Transaction
	err = proto.Unmarshal(bytes, &trans)
	if err != nil {
		return pb.Transaction{}, err
	}
	return trans, nil
}

func SignStart(txStr string) (string, error) {
	trans, err := ParseTxStr(txStr)
	if err != nil {
		return "", err
	}
	rawData, err := proto.Marshal(trans.GetRawData())
	if err != nil {
		return "", err
	}
	s256 := sha256.New()
	s256.Write(rawData)
	hash := s256.Sum(nil)
	return hex.EncodeToString(hash), nil
}

func Sign(data string, privateKey *btcec.PrivateKey) (string, error) {
	hash, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	signature, err := ecdsa.SignCompact(privateKey, hash, false)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

func SignEnd(txStr string, txSignStr string) (string, error) {
	trans, err := ParseTxStr(txStr)
	if err != nil {
		return "", err
	}
	contractList := trans.GetRawData().GetContract()
	for range contractList {
		signature, err := hex.DecodeString(txSignStr)
		if err != nil {
			return "", err
		}
		var sig []byte
		sig = append(sig, signature[1:33]...)
		sig = append(sig, signature[33:65]...)
		sig = append(sig, signature[0]-27)
		trans.Signature = append(trans.Signature, sig)
	}
	bytes, err := proto.Marshal(&trans)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

//	  create a TRC transfer transaction
//		 refBlockBytes - transaction reference block height (take 6-7 2 bytes)
//		 refBlockHash - block hash referenced by the transaction (take 8-15 8 bytes)
//		 expiration - transaction expiration time, beyond this time the transaction will not be packed
//		 timestamp - transaction creation time
//		 fee_limit - the maximum energy consumption of smart contract transactions, only need to be set when deploying or calling smart contracts
func NewTransfer(fromAddress string, toAddress string, amount int64, refBlockBytes string, refBlockHash string, expiration int64, timestamp int64) (string, error) {
	owner, err := GetAddressHash(fromAddress)
	if err != nil {
		return "", err
	}
	to, err := GetAddressHash(toAddress)
	if err != nil {
		return "", err
	}
	transferContract := &pb.TransferContract{OwnerAddress: owner, ToAddress: to, Amount: amount}
	param, err := ptypes.MarshalAny(transferContract)
	if err != nil {
		return "", err
	}
	contract := &pb.Transaction_Contract{Type: pb.Transaction_Contract_TransferContract, Parameter: param}
	raw := new(pb.TransactionRaw)
	refBytes, err := hex.DecodeString(refBlockBytes)
	if err != nil {
		return "", err
	}
	raw.RefBlockBytes = refBytes
	refHash, err := hex.DecodeString(refBlockHash)
	if err != nil {
		return "", err
	}
	raw.RefBlockHash = refHash
	raw.Expiration = expiration
	raw.Timestamp = timestamp
	raw.Contract = []*pb.Transaction_Contract{contract}
	trans := pb.Transaction{RawData: raw}
	data, err := proto.Marshal(&trans)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// create a TRC20 transfer transaction
func NewTRC20TokenTransfer(fromAddress string, toAddress string, contractAddress string, amount *big.Int, feeLimit int64, refBlockBytes string, refBlockHash string, expiration int64, timestamp int64) (string, error) {
	raw := new(pb.TransactionRaw)
	refBytes, err := hex.DecodeString(refBlockBytes)
	if err != nil {
		return "", err
	}
	raw.RefBlockBytes = refBytes
	refHash, err := hex.DecodeString(refBlockHash)
	if err != nil {
		return "", err
	}
	raw.RefBlockHash = refHash
	raw.Expiration = expiration
	raw.Timestamp = timestamp

	fromAddressHash, err := GetAddressHash(fromAddress)
	if err != nil {
		return "", err
	}
	toAddressHash, err := GetAddressHash(toAddress)
	if err != nil {
		return "", err
	}
	contractAddressHash, err := GetAddressHash(contractAddress)
	if err != nil {
		return "", err
	}
	input, err := token.Transfer(hex.EncodeToString(toAddressHash), amount)
	if err != nil {
		return "", err
	}
	transferContract := &pb.TriggerSmartContract{OwnerAddress: fromAddressHash, ContractAddress: contractAddressHash, CallValue: 0, CallTokenValue: 0, Data: input}
	param, err := ptypes.MarshalAny(transferContract)
	if err != nil {
		return "", err
	}

	contract := &pb.Transaction_Contract{Type: pb.Transaction_Contract_TriggerSmartContract, Parameter: param}
	raw.FeeLimit = feeLimit

	raw.Contract = []*pb.Transaction_Contract{contract}
	trans := pb.Transaction{RawData: raw}
	data, err := proto.Marshal(&trans)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}
