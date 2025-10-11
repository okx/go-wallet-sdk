package zksync

import (
	"github.com/emresenyuva/go-wallet-sdk/coins/zksync/core"
	"math/big"
	"strconv"
)

const (
	WITHDRAW       = "Withdraw"
	TRANSFER       = "Transfer"
	CHANGE_PUB_KEY = "ChangePubKey"
	MINT_NFT       = "MintNFT"
	WITHDRAW_NFT   = "WithdrawNFT"
)

type BatchTransaction struct {
	Txs       []core.SignedTransaction `json:"txs"`
	Signature *core.EthSignature       `json:"signature"`
}

func CreateWithdrawTx(accountId uint32, address string, amount *big.Int, fee *big.Int, token *core.Token, nonce uint32, validFrom, validUntil uint64) *core.Withdraw {
	tx := &core.Withdraw{
		Type:      WITHDRAW,
		AccountId: accountId,
		From:      address,
		To:        address,
		TokenId:   token.Id,
		Amount:    amount,
		Nonce:     nonce,
		Fee:       fee.String(),
		TimeRange: &core.TimeRange{
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
	}

	return tx
}

func CreateTransferTx(accountId uint32, fromAddress, toAddress string, amount *big.Int, fee *big.Int, token *core.Token, nonce uint32, validFrom, validUntil uint64) *core.Transfer {
	tx := &core.Transfer{
		Type:      TRANSFER,
		AccountId: accountId,
		From:      fromAddress,
		To:        toAddress,
		Token:     token,
		TokenId:   token.Id,
		Amount:    amount,
		Nonce:     nonce,
		Fee:       fee.String(),
		TimeRange: &core.TimeRange{
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
	}

	return tx
}

func CreateTransferWithFeeTokenTx(accountId uint32, fromAddress, toAddress string, amount *big.Int, token *core.Token, fee *big.Int, feeToken *core.Token, nonce uint32, validFrom, validUntil uint64) []*core.Transfer {
	tx := &core.Transfer{
		Type:      TRANSFER,
		AccountId: accountId,
		From:      fromAddress,
		To:        toAddress,
		Token:     token,
		TokenId:   token.Id,
		Amount:    amount,
		Fee:       big.NewInt(0).String(),
		Nonce:     nonce,
		TimeRange: &core.TimeRange{
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
	}

	feeTx := &core.Transfer{
		Type:      TRANSFER,
		AccountId: accountId,
		From:      fromAddress,
		To:        fromAddress,
		Token:     feeToken,
		TokenId:   feeToken.Id,
		Amount:    big.NewInt(0),
		Nonce:     nonce + 1,
		Fee:       fee.String(),
		TimeRange: &core.TimeRange{
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
	}

	return []*core.Transfer{tx, feeTx}
}

func CreateChangePubKeyTx(accountId uint32, fromAddress, pubKeyHash string, feeToken *core.Token, fee *big.Int, nonce uint32, validFrom, validUntil uint64) *core.ChangePubKey {
	tx := &core.ChangePubKey{
		Type:        CHANGE_PUB_KEY,
		AccountId:   accountId,
		Account:     fromAddress,
		NewPkHash:   pubKeyHash,
		FeeToken:    feeToken.Id,
		Fee:         fee.String(),
		Nonce:       nonce,
		EthAuthData: &core.ChangePubKeyOnchain{Type: core.ChangePubKeyAuthTypeOnchain},
		TimeRange: &core.TimeRange{
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
	}
	return tx
}

func CreateMintNFTTx(accountId uint32, creator, recipient, contentHash string, feeToken *core.Token, fee *big.Int, nonce uint32) *core.MintNFT {
	hash := core.HexToHash(contentHash)

	tx := &core.MintNFT{
		Type:           MINT_NFT,
		CreatorId:      accountId,
		CreatorAddress: creator,
		ContentHash:    hash,
		Recipient:      recipient,
		Nonce:          nonce,
		Fee:            fee.String(),
		FeeToken:       feeToken.Id,
	}
	return tx
}

func CreateTransferNFTTx(accountId uint32, fromAddress, toAddress, nftSymbol string, feeToken *core.Token, fee *big.Int, nonce uint32, validFrom, validUntil uint64) ([]*core.Transfer, error) {
	nftIdStr := nftSymbol[4:]
	nftId, err := strconv.Atoi(nftIdStr)
	if err != nil {
		return nil, err
	}
	nft := core.NFT{
		Id:     uint32(nftId),
		Symbol: nftSymbol,
	}
	nftTx := &core.Transfer{
		Type:      TRANSFER,
		AccountId: accountId,
		From:      fromAddress,
		To:        toAddress,
		Token:     nft.ToToken(),
		TokenId:   uint32(nftId),
		Amount:    big.NewInt(1),
		Nonce:     nonce,
		Fee:       big.NewInt(0).String(),
		TimeRange: &core.TimeRange{
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
	}
	feeTx := &core.Transfer{
		Type:      TRANSFER,
		AccountId: accountId,
		From:      fromAddress,
		To:        fromAddress,
		Token:     feeToken,
		TokenId:   feeToken.Id,
		Amount:    big.NewInt(0),
		Nonce:     nonce + 1,
		Fee:       fee.String(),
		TimeRange: &core.TimeRange{
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
	}
	return []*core.Transfer{nftTx, feeTx}, nil
}

func CreateWithdrawNFTTx(accountId uint32, addr string, nftId uint32, feeToken *core.Token, fee *big.Int, nonce uint32, validFrom, validUntil uint64) *core.WithdrawNFT {
	tx := &core.WithdrawNFT{
		Type:      WITHDRAW_NFT,
		AccountId: accountId,
		From:      addr,
		To:        addr,
		Token:     nftId,
		Nonce:     nonce,
		Fee:       fee.String(),
		FeeToken:  feeToken.Id,
		TimeRange: &core.TimeRange{
			ValidFrom:  validFrom,
			ValidUntil: validUntil,
		},
	}
	return tx
}
