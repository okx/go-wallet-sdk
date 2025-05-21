package near

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/near/serialize"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"math/big"
)

const (
	SignatureLength = 65
)

type Transaction struct {
	SignerId   serialize.String
	PublicKey  serialize.PublicKey
	Nonce      serialize.U64
	ReceiverId serialize.String
	BlockHash  serialize.BlockHash
	Actions    []serialize.IAction
}

type SignatureData struct {
	V *big.Int
	R *big.Int
	S *big.Int

	ByteV byte
	ByteR []byte
	ByteS []byte
}
type ActionTransfer struct {
	Transfer Transfer `json:"transfer"`
}

type Transfer struct {
	Deposit string `json:"deposit"` //amount
}

func CreateTransaction(from, to, publicKeyHex, blockHash string, nonce int64) (*Transaction, error) {
	bh := base58.Decode(blockHash)
	if len(bh) == 0 {
		return nil, fmt.Errorf("base58  decode blockhash error ,BlockHash=%s", blockHash)
	}

	pubBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("public key hex decode error,public key hex=%s", publicKeyHex)
	}

	if len(pubBytes) != 32 {
		return nil, fmt.Errorf("public key len error,public key=%s", publicKeyHex)
	}

	tx := Transaction{
		SignerId:   serialize.String{Value: from},
		PublicKey:  serialize.PublicKey{KeyType: 0, Value: pubBytes},
		Nonce:      serialize.U64{Value: uint64(nonce)},
		ReceiverId: serialize.String{Value: to},
		BlockHash:  serialize.BlockHash{Value: bh},
		Actions:    nil,
	}

	return &tx, nil
}

func (tx *Transaction) SetAction(action ...serialize.IAction) {
	tx.Actions = append(tx.Actions, action...)
}

func (tx *Transaction) Serialize() ([]byte, error) {
	var (
		data []byte
	)
	ss, err := tx.SignerId.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: signerId error,Err=%v", err)
	}
	data = append(data, ss...)
	ps, err := tx.PublicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: publickey error,Err=%v", err)
	}
	data = append(data, ps...)
	ns, err := tx.Nonce.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: nonce error,Err=%v", err)
	}
	data = append(data, ns...)
	rs, err := tx.ReceiverId.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: ReceiverId error,Err=%v", err)
	}
	data = append(data, rs...)
	bs, err := tx.BlockHash.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: blockhash error,Err=%v", err)
	}
	data = append(data, bs...)
	//serialize action
	al := len(tx.Actions)
	uAL := serialize.U32{
		Value: uint32(al),
	}
	uALData, err := uAL.Serialize()
	if err != nil {
		return nil, fmt.Errorf("tx serialize: action length error,Err=%v", err)
	}
	data = append(data, uALData...)
	for _, action := range tx.Actions {
		as, err := action.Serialize()
		if err != nil {
			return nil, fmt.Errorf("tx serialize: action error,Err=%v", err)
		}
		data = append(data, as...)
	}
	return data, nil
}

func SignTransaction(txBase58 string, privateKey string) (string, error) {
	pkBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	key := ed25519.PrivateKey(pkBytes)

	data := base58.Decode(txBase58)

	txHash := sha256.Sum256(data)
	sig := ed25519.Sign(key, txHash[:])
	if len(sig) != 64 {
		return "", fmt.Errorf("sign error,length is not equal 64,length=%d", len(sig))
	}
	return base58.Encode(sig), nil
}

type SignatureTransaction struct {
	Sig serialize.Signature
	Tx  *Transaction
}

func CreateSignedTransaction(tx *Transaction, sig string) (*SignatureTransaction, error) {

	signature := base58.Decode(sig)
	if len(signature) == 0 {
		return nil, fmt.Errorf("base58 decode sig error,sig=%s", sig)
	}

	stx := SignatureTransaction{
		Sig: serialize.Signature{KeyType: tx.PublicKey.KeyType, Value: signature},
		Tx:  tx,
	}

	return &stx, nil
}

func (stx *SignatureTransaction) Serialize() ([]byte, error) {
	data, err := stx.Tx.Serialize()
	if err != nil {
		return nil, fmt.Errorf("sign serialize: tx serialize error,Err=%v", err)
	}
	ss, err := stx.Sig.Serialize()
	if err != nil {
		return nil, fmt.Errorf("sign serialize: sig serialize error,Err=%v", err)
	}
	data = append(data, ss...)
	return data, nil
}

func (tx *Transaction) GetSigningHash() (string, error) {
	signedTx, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	txHash := sha256.Sum256(signedTx[:])
	return hex.EncodeToString(txHash[:]), nil
}

func CalTxHash(tx string, signed bool) (string, error) {
	data, err := base64.StdEncoding.DecodeString(tx)
	if err != nil {
		return "", err
	}
	if signed {
		if len(data) < SignatureLength {
			return "", err
		}
		data = data[:len(data)-SignatureLength] // Remove signature portion of tx
	}
	txHash := sha256.Sum256(data)
	return base58.Encode(txHash[:]), nil
}
