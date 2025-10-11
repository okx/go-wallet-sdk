package zkspace

import (
	"encoding/hex"
	"github.com/emresenyuva/go-wallet-sdk/coins/zksync/core"
	"math/big"
)

func CreateChangePubKeyTx(accountId uint32, from string, newPkHash string, nonce uint32, ethSignature []byte) ChangePubKey {
	tx := ChangePubKey{
		Type:         "ChangePubKey",
		AccountId:    accountId,
		Account:      from,
		NewPkHash:    newPkHash,
		Nonce:        nonce,
		EthSignature: "0x" + hex.EncodeToString(ethSignature),
	}
	return tx
}

func CreateTransferTx(accountId uint32, from string, to string, tokenId uint16, amount *big.Int, feeTokenId uint8, fee *big.Int,
	chainId uint8, nonce uint32) Transfer {
	transferTx := Transfer{
		Type:       "Transfer",
		AccountId:  accountId,
		From:       from,
		To:         to,
		TokenId:    tokenId,
		Amount:     amount,
		FeeTokenId: feeTokenId,
		Fee:        fee,
		ChainId:    chainId,
		Nonce:      nonce,
	}
	return transferTx
}

type SignedTransaction struct {
	Transaction       SignTransfer       `json:"tx"`
	EthereumSignature *core.EthSignature `json:"signature"`
}

func CreateSignTransferTx(tx *Transfer) SignTransfer {
	transferTx := SignTransfer{
		Type:       "Transfer",
		AccountId:  tx.AccountId,
		From:       tx.From,
		To:         tx.To,
		TokenId:    tx.TokenId,
		Amount:     tx.Amount.String(),
		FeeTokenId: tx.FeeTokenId,
		Fee:        tx.Fee.String(),
		ChainId:    tx.ChainId,
		Nonce:      tx.Nonce,
		Signature:  tx.Signature,
	}
	return transferTx
}
