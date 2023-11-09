package zksync

import (
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/zksync/core"
	"math/big"
)

type signer struct {
	ethSigner *core.OkEthSigner
	zkSigner  *core.ZkSigner
}

func newSigner(ethPrivKeyHex string, chainId int) (*signer, error) {
	privKeyBytes, err := hex.DecodeString(ethPrivKeyHex)
	if err != nil {
		return nil, err
	}
	ethSigner, err := core.NewOkEthSignerFromPrivBytes(privKeyBytes)
	if err != nil {
		return nil, err
	}

	zkSigner, err := core.NewZkSignerFromEthSigner(ethSigner, core.ChainId(chainId))
	if err != nil {
		return nil, err
	}
	return &signer{ethSigner: ethSigner, zkSigner: zkSigner}, nil
}

func (s *signer) getPublicKeyHash() string {
	return s.zkSigner.GetPublicKeyHash()
}

func (s signer) signChangePubKey(txData *core.ChangePubKey) (*core.Signature, error) {
	return s.zkSigner.SignChangePubKey(txData)
}

func (s *signer) signWithdraw(txData *core.Withdraw) (*core.Signature, error) {
	return s.zkSigner.SignWithdraw(txData)
}

func (s *signer) signTransfer(transfer *core.Transfer) (*core.Signature, error) {
	return s.zkSigner.SignTransfer(transfer)
}

func (s *signer) signMintNFT(txData *core.MintNFT) (*core.Signature, error) {
	return s.zkSigner.SignMintNFT(txData)
}

func (s *signer) signWithdrawNFT(txData *core.WithdrawNFT) (*core.Signature, error) {
	return s.zkSigner.SignWithdrawNFT(txData)
}

func (s *signer) signAuth(txData *core.ChangePubKey) (*core.ChangePubKeyECDSA, error) {
	return s.ethSigner.SignAuth(txData)
}

func (s *signer) signTransaction(tx core.ZksTransaction, nonce uint32, token *core.Token, fee *big.Int) (*core.EthSignature, error) {
	return s.ethSigner.SignTransaction(tx, nonce, token, fee)
}

func (s *signer) signBatch(txs []core.ZksTransaction, nonce uint32, token *core.Token, fee *big.Int) (*core.EthSignature, error) {
	return s.ethSigner.SignBatch(txs, nonce, token, fee)
}
