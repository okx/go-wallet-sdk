package zksync

import (
	"encoding/hex"
	"github.com/emresenyuva/go-wallet-sdk/coins/zksync/core"
	"math/big"
)

func GetPubKeyHash(ethPrivKeyHex string, chainId int) (string, error) {
	signer, err := newSigner(ethPrivKeyHex, chainId)
	if err != nil {
		return "", err
	}

	return signer.getPublicKeyHash(), nil
}

func GetAddress(ethPrivKeyHex string) (string, error) {
	privKeyBytes, err := hex.DecodeString(ethPrivKeyHex)
	if err != nil {
		return "", err
	}
	ethSigner, err := core.NewOkEthSignerFromPrivBytes(privKeyBytes)
	if err != nil {
		return "", err
	}

	return ethSigner.GetAddress(), nil
}

func SignWithdraw(withdraw *core.Withdraw, feeToken *core.Token, ethPrivKeyHex string, chainId int) (*core.SignedTransaction, error) {
	signer, err := newSigner(ethPrivKeyHex, chainId)
	if err != nil {
		return nil, err
	}

	sign, _ := signer.zkSigner.SignWithdraw(withdraw)
	withdraw.Signature = sign

	fee, _ := new(big.Int).SetString(withdraw.Fee, 10)

	ethSignature, err := signer.signTransaction(withdraw, withdraw.Nonce, feeToken, fee)
	if err != nil {
		return nil, err
	}

	signedTransaction := &core.SignedTransaction{
		Transaction:       withdraw,
		EthereumSignature: ethSignature,
	}

	return signedTransaction, nil
}

func SignTransfer(transfer *core.Transfer, ethPrivKeyHex string, chainId int) (*core.SignedTransaction, error) {
	signer, err := newSigner(ethPrivKeyHex, chainId)
	if err != nil {
		return nil, err
	}

	sign, _ := signer.signTransfer(transfer)
	transfer.Signature = sign

	fee, _ := new(big.Int).SetString(transfer.Fee, 10)

	ethSignature, err := signer.signTransaction(transfer, transfer.Nonce, nil, fee)
	if err != nil {
		return nil, err
	}

	signedTransaction := &core.SignedTransaction{
		Transaction:       transfer,
		EthereumSignature: ethSignature,
	}

	return signedTransaction, nil

}

func SignBatchTransfer(transfers []*core.Transfer, ethPrivKeyHex string, chainId int) (*BatchTransaction, error) {
	signer, err := newSigner(ethPrivKeyHex, chainId)
	if err != nil {
		return nil, err
	}

	feeStr := "0"
	nonce := UINT32_MAX

	signedTxs := []core.SignedTransaction{}
	zksTxs := []core.ZksTransaction{}
	for _, transfer := range transfers {
		sign, _ := signer.signTransfer(transfer)
		transfer.Signature = sign

		zksTx := core.ZksTransaction(transfer)
		signedTxs = append(signedTxs, core.SignedTransaction{
			Transaction: zksTx,
		})

		zksTxs = append(zksTxs, zksTx)

		if transfer.Fee != "0" {
			feeStr = transfer.Fee
		}

		if transfer.Nonce < nonce {
			nonce = transfer.Nonce
		}
	}

	fee, _ := new(big.Int).SetString(feeStr, 10)

	ethSignature, err := signer.signBatch(zksTxs, nonce, nil, fee)
	if err != nil {
		return nil, err
	}

	signedTransaction := &BatchTransaction{
		Txs:       signedTxs,
		Signature: ethSignature,
	}

	return signedTransaction, nil
}

func SignChangePubKey(changePubKey *core.ChangePubKey, ethPrivKeyHex string, chainId int) (*core.SignedTransaction, error) {
	signer, err := newSigner(ethPrivKeyHex, chainId)
	if err != nil {
		return nil, err
	}

	sign, _ := signer.signChangePubKey(changePubKey)
	changePubKey.Signature = sign

	_, err = signer.signAuth(changePubKey)
	if err != nil {
		return nil, err
	}

	signedTransaction := &core.SignedTransaction{
		Transaction: changePubKey,
	}

	return signedTransaction, nil
}

func SignMintNFT(mintNFT *core.MintNFT, feeToken *core.Token, ethPrivKeyHex string, chainId int) (*core.SignedTransaction, error) {
	signer, err := newSigner(ethPrivKeyHex, chainId)
	if err != nil {
		return nil, err
	}

	sign, _ := signer.signMintNFT(mintNFT)
	mintNFT.Signature = sign

	fee, _ := new(big.Int).SetString(mintNFT.Fee, 10)

	ethSignature, err := signer.signTransaction(mintNFT, mintNFT.Nonce, feeToken, fee)
	if err != nil {
		return nil, err
	}

	signedTransaction := &core.SignedTransaction{
		Transaction:       mintNFT,
		EthereumSignature: ethSignature,
	}

	return signedTransaction, nil
}

func SignWithdrawNFT(withdrawNft *core.WithdrawNFT, feeToken *core.Token, ethPrivKeyHex string, chainId int) (*core.SignedTransaction, error) {
	signer, err := newSigner(ethPrivKeyHex, chainId)
	if err != nil {
		return nil, err
	}

	sign, _ := signer.signWithdrawNFT(withdrawNft)
	withdrawNft.Signature = sign

	fee, _ := new(big.Int).SetString(withdrawNft.Fee, 10)

	ethSignature, err := signer.signTransaction(withdrawNft, withdrawNft.Nonce, feeToken, fee)
	if err != nil {
		return nil, err
	}

	signedTransaction := &core.SignedTransaction{
		Transaction:       withdrawNft,
		EthereumSignature: ethSignature,
	}

	return signedTransaction, nil
}
