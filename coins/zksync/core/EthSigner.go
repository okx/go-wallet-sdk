package core

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/coins/ethereum"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"
)

type EthSigner interface {
	GetAddress() string
	SignMessage([]byte) ([]byte, error)
	SignHash(msg []byte) ([]byte, error)
	SignAuth(txData *ChangePubKey) (*ChangePubKeyECDSA, error)
	SignTransaction(tx ZksTransaction, nonce uint32, token *Token, fee *big.Int) (*EthSignature, error)
	SignBatch(txs []ZksTransaction, nonce uint32, token *Token, fee *big.Int) (*EthSignature, error)
	SignOrder(order *Order, sell, buy *Token) (*EthSignature, error)
}

type EthSignatureType string

const (
	EthSignatureTypeEth     EthSignatureType = "EthereumSignature"
	EthSignatureTypeEIP1271 EthSignatureType = "EIP1271Signature"
)

type EthSignature struct {
	Type      EthSignatureType `json:"type"`
	Signature string           `json:"signature"`
}

type OkEthSigner struct {
	privKey []byte
}

func NewOkEthSignerFromPrivBytes(privKeyBytes []byte) (*OkEthSigner, error) {
	return &OkEthSigner{privKeyBytes}, nil
}

func (s *OkEthSigner) GetAddress() string {
	ethPrivKey, _ := btcec.PrivKeyFromBytes(s.privKey)
	pubKey := ethPrivKey.PubKey()
	pubBytes := pubKey.SerializeUncompressed()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:])
	addressByte := hash.Sum(nil)
	return "0x" + hex.EncodeToString(addressByte)[24:]

}

func (s *OkEthSigner) SignMessage(msg []byte) ([]byte, error) {
	prvKey, _ := btcec.PrivKeyFromBytes(s.privKey)
	prefix := fmt.Sprintf(ethereum.MessagePrefixTmp, len(msg))
	toSign := append([]byte(prefix), msg...)
	message, err := ethereum.SignMessage(toSign, prvKey)
	if err != nil {
		return nil, err
	}
	return message.ToBytes(), nil
}

func (s *OkEthSigner) SignHash(msg []byte) ([]byte, error) {
	prvKey, _ := btcec.PrivKeyFromBytes(s.privKey)
	message, err := ethereum.SignMessage(msg, prvKey)
	if err != nil {
		return nil, err
	}
	return message.ToBytes(), nil
}

func (s *OkEthSigner) SignAuth(txData *ChangePubKey) (*ChangePubKeyECDSA, error) {
	auth := &ChangePubKeyECDSA{
		Type:         ChangePubKeyAuthTypeECDSA,
		EthSignature: "",
		BatchHash:    "0x" + hex.EncodeToString(make([]byte, 32)),
	}
	txData.EthAuthData = auth
	msg, err := getChangePubKeyData(txData)
	if err != nil {
		return nil, errors.New("failed to get ChangePubKey data for sign")
	}
	sig, err := s.SignMessage(msg)
	if err != nil {
		return nil, errors.New("failed to sign ChangePubKeyECDSA msg")
	}
	auth.EthSignature = "0x" + hex.EncodeToString(sig)
	return auth, nil
}

func (s *OkEthSigner) SignTransaction(tx ZksTransaction, nonce uint32, token *Token, fee *big.Int) (*EthSignature, error) {
	switch tx.getType() {
	case "ChangePubKey":
		if txData, ok := tx.(*ChangePubKey); ok {
			msg, err := getChangePubKeyData(txData)
			if err != nil {
				return nil, errors.New("failed to get ChangePubKey data for sign")
			}
			sig, err := s.SignMessage(msg)
			if err != nil {
				return nil, errors.New("failed to sign ChangePubKey tx")
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: "0x" + hex.EncodeToString(sig),
			}, nil
		}
	case "Transfer":
		if txData, ok := tx.(*Transfer); ok {
			var tokenToUse *Token
			if txData.Token != nil {
				tokenToUse = txData.Token
			} else {
				tokenToUse = token
			}
			fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
			if !ok {
				return nil, errors.New("failed to convert string fee to big.Int")
			}
			msg, err := getTransferMessagePart(txData.To, txData.Amount, fee, tokenToUse)
			if err != nil {
				return nil, errors.New("failed to get Transfer message part")
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, errors.New("failed to sign Transfer tx")
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: "0x" + hex.EncodeToString(sig),
			}, nil
		}
	case "Withdraw":
		if txData, ok := tx.(*Withdraw); ok {
			msg, err := getWithdrawMessagePart(txData.To, txData.Amount, fee, token)
			if err != nil {
				return nil, errors.New("failed to get Withdraw message part")
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, errors.New("failed to sign Withdraw tx")
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: "0x" + hex.EncodeToString(sig),
			}, nil
		}
	case "ForcedExit":
		if txData, ok := tx.(*ForcedExit); ok {
			msg, err := getForcedExitMessagePart(txData.Target, fee, token)
			if err != nil {
				return nil, errors.New("failed to get ForcedExit message part")
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, errors.New("failed to sign ForcedExit tx")
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: "0x" + hex.EncodeToString(sig),
			}, nil
		}
	case "MintNFT":
		if txData, ok := tx.(*MintNFT); ok {
			msg, err := getMintNFTMessagePart(txData.ContentHash, txData.Recipient, fee, token)
			if err != nil {
				return nil, errors.New("failed to get MintNFT message part")
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, errors.New("failed to sign MintNFT tx")
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: "0x" + hex.EncodeToString(sig),
			}, nil
		}
	case "WithdrawNFT":
		if txData, ok := tx.(*WithdrawNFT); ok {
			msg, err := getWithdrawNFTMessagePart(txData.To, txData.Token, fee, token)
			if err != nil {
				return nil, errors.New("failed to get WithdrawNFT message part")
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, errors.New("failed to sign WithdrawNFT tx")
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: "0x" + hex.EncodeToString(sig),
			}, nil
		}
	case "Swap":
		msg := getSwapMessagePart(token, fee)
		msg += "\n" + getNonceMessagePart(nonce)
		sig, err := s.SignMessage([]byte(msg))
		if err != nil {
			return nil, errors.New("failed to sign Swap tx")
		}
		return &EthSignature{
			Type:      EthSignatureTypeEth,
			Signature: "0x" + hex.EncodeToString(sig),
		}, nil
	}
	return nil, errors.New("unknown tx type")
}

func (s *OkEthSigner) SignBatch(txs []ZksTransaction, nonce uint32, token *Token, fee *big.Int) (*EthSignature, error) {
	batchMsgs := make([]string, 0, len(txs))
	for _, tx := range txs {

		switch tx.getType() {
		case "Transfer":
			if txData, ok := tx.(*Transfer); ok {
				var tokenToUse *Token
				if txData.Token != nil {
					tokenToUse = txData.Token
				} else {
					tokenToUse = token
				}
				fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
				if !ok {
					return nil, errors.New("failed to convert string fee to big.Int")
				}
				msg, err := getTransferMessagePart(txData.To, txData.Amount, fee, tokenToUse)
				if err != nil {
					return nil, errors.New("failed to get Transfer message part")
				}
				batchMsgs = append(batchMsgs, msg)
			}
		case "Withdraw":
			if txData, ok := tx.(*Withdraw); ok {
				msg, err := getWithdrawMessagePart(txData.To, txData.Amount, fee, token)
				if err != nil {
					return nil, errors.New("failed to get Withdraw message part")
				}
				batchMsgs = append(batchMsgs, msg)
			}
		case "ForcedExit":
			if txData, ok := tx.(*ForcedExit); ok {
				msg, err := getForcedExitMessagePart(txData.Target, fee, token)
				if err != nil {
					return nil, errors.New("failed to get ForcedExit message part")
				}
				batchMsgs = append(batchMsgs, msg)
			}
		case "MintNFT":
			if txData, ok := tx.(*MintNFT); ok {
				msg, err := getMintNFTMessagePart(txData.ContentHash, txData.Recipient, fee, token)
				if err != nil {
					return nil, errors.New("failed to get MintNFT message part")
				}
				batchMsgs = append(batchMsgs, msg)
			}
		case "WithdrawNFT":
			if txData, ok := tx.(*WithdrawNFT); ok {
				msg, err := getWithdrawNFTMessagePart(txData.To, txData.Token, fee, token)
				if err != nil {
					return nil, errors.New("failed to get WithdrawNFT message part")
				}
				batchMsgs = append(batchMsgs, msg)
			}
		default:
			return nil, errors.New("unknown tx type")
		}
	}
	batchMsg := strings.Join(batchMsgs, "\n")
	batchMsg += "\n" + getNonceMessagePart(nonce)
	sig, err := s.SignMessage([]byte(batchMsg))
	if err != nil {
		return nil, errors.New("failed to sign batch of txs")
	}
	return &EthSignature{
		Type:      EthSignatureTypeEth,
		Signature: "0x" + hex.EncodeToString(sig),
	}, nil
}

func (s *OkEthSigner) SignOrder(order *Order, sell, buy *Token) (*EthSignature, error) {
	msg, err := getOrderMessagePart(order.RecipientAddress, order.Amount, sell, buy, order.Ratio)
	if err != nil {
		return nil, errors.New("failed to get Order message part")
	}
	msg += "\n" + getNonceMessagePart(order.Nonce)
	sig, err := s.SignMessage([]byte(msg))
	if err != nil {
		return nil, errors.New("failed to sign Order")
	}
	return &EthSignature{
		Type:      EthSignatureTypeEth,
		Signature: "0x" + hex.EncodeToString(sig),
	}, nil
}
