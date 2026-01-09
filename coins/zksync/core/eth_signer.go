/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
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

var (
	ErrUnknownTxTYPE = errors.New("unknown tx type")
	ErrConvertBigInt = errors.New("failed to convert string fee to big.Int")
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
	return HEX_PREFIX + hex.EncodeToString(addressByte)[24:]

}

func (s *OkEthSigner) SignMessage(msg []byte) ([]byte, error) {
	prvKey, _ := btcec.PrivKeyFromBytes(s.privKey)
	prefix := fmt.Sprintf("\u0019Ethereum Signed Message:\n%d", len(msg))
	toSign := append([]byte(prefix), msg...)
	message := ethereum.SignMessage(toSign, prvKey)
	return message.ToBytes(), nil
}

func (s *OkEthSigner) SignHash(msg []byte) ([]byte, error) {
	prvKey, _ := btcec.PrivKeyFromBytes(s.privKey)
	message := ethereum.SignMessage(msg, prvKey)
	return message.ToBytes(), nil
}

func (s *OkEthSigner) SignAuth(txData *ChangePubKey) (*ChangePubKeyECDSA, error) {
	auth := &ChangePubKeyECDSA{
		Type:         ChangePubKeyAuthTypeECDSA,
		EthSignature: "",
		BatchHash:    HEX_PREFIX + hex.EncodeToString(make([]byte, 32)),
	}
	txData.EthAuthData = auth
	msg, err := getChangePubKeyData(txData)
	if err != nil {
		return nil, err
	}
	sig, err := s.SignMessage(msg)
	if err != nil {
		return nil, err
	}
	auth.EthSignature = HEX_PREFIX + hex.EncodeToString(sig)
	return auth, nil
}

func (s *OkEthSigner) SignTransaction(tx ZksTransaction, nonce uint32, token *Token, fee *big.Int) (*EthSignature, error) {
	switch tx.getType() {
	// case "ChangePubKey":
	case TransactionTypeChangePubKey_:
		if txData, ok := tx.(*ChangePubKey); ok {
			msg, err := getChangePubKeyData(txData)
			if err != nil {
				return nil, err
			}
			sig, err := s.SignMessage(msg)
			if err != nil {
				return nil, err
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: HEX_PREFIX + hex.EncodeToString(sig),
			}, nil
		}
	case TransactionTypeTransfer:
		if txData, ok := tx.(*Transfer); ok {
			var tokenToUse *Token
			if txData.Token != nil {
				tokenToUse = txData.Token
			} else {
				tokenToUse = token
			}
			fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
			if !ok {
				return nil, ErrConvertBigInt
			}
			msg, err := getTransferMessagePart(txData.To, txData.Amount, fee, tokenToUse)
			if err != nil {
				return nil, err
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, err
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: HEX_PREFIX + hex.EncodeToString(sig),
			}, nil
		}
	case TransactionTypeWithdraw:
		if txData, ok := tx.(*Withdraw); ok {
			msg, err := getWithdrawMessagePart(txData.To, txData.Amount, fee, token)
			if err != nil {
				return nil, err
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, err
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: HEX_PREFIX + hex.EncodeToString(sig),
			}, nil
		}
	case TransactionTypeForcedExit:
		if txData, ok := tx.(*ForcedExit); ok {
			msg, err := getForcedExitMessagePart(txData.Target, fee, token)
			if err != nil {
				return nil, err
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, err
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: HEX_PREFIX + hex.EncodeToString(sig),
			}, nil
		}
	case TransactionTypeMintNFT:
		if txData, ok := tx.(*MintNFT); ok {
			msg, err := getMintNFTMessagePart(txData.ContentHash, txData.Recipient, fee, token)
			if err != nil {
				return nil, err
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, err
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: HEX_PREFIX + hex.EncodeToString(sig),
			}, nil
		}
	case TransactionTypeWithdrawNFT:
		if txData, ok := tx.(*WithdrawNFT); ok {
			msg, err := getWithdrawNFTMessagePart(txData.To, txData.Token, fee, token)
			if err != nil {
				return nil, err
			}
			msg += "\n" + getNonceMessagePart(nonce)
			sig, err := s.SignMessage([]byte(msg))
			if err != nil {
				return nil, err
			}
			return &EthSignature{
				Type:      EthSignatureTypeEth,
				Signature: HEX_PREFIX + hex.EncodeToString(sig),
			}, nil
		}
	case TransactionTypeSwap:
		msg := getSwapMessagePart(token, fee)
		msg += "\n" + getNonceMessagePart(nonce)
		sig, err := s.SignMessage([]byte(msg))
		if err != nil {
			return nil, err
		}
		return &EthSignature{
			Type:      EthSignatureTypeEth,
			Signature: HEX_PREFIX + hex.EncodeToString(sig),
		}, nil
	}
	return nil, ErrUnknownTxTYPE
}

func (s *OkEthSigner) SignBatch(txs []ZksTransaction, nonce uint32, token *Token, fee *big.Int) (*EthSignature, error) {
	batchMsgs := make([]string, 0, len(txs))
	for _, tx := range txs {

		switch tx.getType() {
		case TransactionTypeTransfer:
			if txData, ok := tx.(*Transfer); ok {
				var tokenToUse *Token
				if txData.Token != nil {
					tokenToUse = txData.Token
				} else {
					tokenToUse = token
				}
				fee, ok := big.NewInt(0).SetString(txData.Fee, 10)
				if !ok {
					return nil, ErrConvertBigInt
				}
				msg, err := getTransferMessagePart(txData.To, txData.Amount, fee, tokenToUse)
				if err != nil {
					return nil, err
				}
				batchMsgs = append(batchMsgs, msg)
			}
		case TransactionTypeWithdraw:
			if txData, ok := tx.(*Withdraw); ok {
				msg, err := getWithdrawMessagePart(txData.To, txData.Amount, fee, token)
				if err != nil {
					return nil, err
				}
				batchMsgs = append(batchMsgs, msg)
			}
		case TransactionTypeForcedExit:
			if txData, ok := tx.(*ForcedExit); ok {
				msg, err := getForcedExitMessagePart(txData.Target, fee, token)
				if err != nil {
					return nil, err
				}
				batchMsgs = append(batchMsgs, msg)
			}
		case TransactionTypeMintNFT:
			if txData, ok := tx.(*MintNFT); ok {
				msg, err := getMintNFTMessagePart(txData.ContentHash, txData.Recipient, fee, token)
				if err != nil {
					return nil, err
				}
				batchMsgs = append(batchMsgs, msg)
			}
		case TransactionTypeWithdrawNFT:
			if txData, ok := tx.(*WithdrawNFT); ok {
				msg, err := getWithdrawNFTMessagePart(txData.To, txData.Token, fee, token)
				if err != nil {
					return nil, err
				}
				batchMsgs = append(batchMsgs, msg)
			}
		default:
			return nil, ErrUnknownTxTYPE
		}
	}
	batchMsg := strings.Join(batchMsgs, "\n")
	batchMsg += "\n" + getNonceMessagePart(nonce)
	sig, err := s.SignMessage([]byte(batchMsg))
	if err != nil {
		return nil, err
	}
	return &EthSignature{
		Type:      EthSignatureTypeEth,
		Signature: HEX_PREFIX + hex.EncodeToString(sig),
	}, nil
}

func (s *OkEthSigner) SignOrder(order *Order, sell, buy *Token) (*EthSignature, error) {
	msg, err := getOrderMessagePart(order.RecipientAddress, order.Amount, sell, buy, order.Ratio)
	if err != nil {
		return nil, err
	}
	msg += "\n" + getNonceMessagePart(order.Nonce)
	sig, err := s.SignMessage([]byte(msg))
	if err != nil {
		return nil, err
	}
	return &EthSignature{
		Type:      EthSignatureTypeEth,
		Signature: HEX_PREFIX + hex.EncodeToString(sig),
	}, nil
}
