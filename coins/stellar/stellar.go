package stellar

import (
	"errors"
	"github.com/okx/go-wallet-sdk/coins/stellar/keypair"
	"github.com/okx/go-wallet-sdk/coins/stellar/strkey"
	"github.com/okx/go-wallet-sdk/coins/stellar/txnbuild"
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
	"strings"
)

func NewRandSecret() (string, error) {
	for i := 0; i < 5; i++ {
		kp, err := keypair.Random()
		if err != nil {
			continue
		}
		return kp.Seed(), nil
	}
	return "", errors.New("failed to generate random secret")
}

func GetNewAddress(secret string) (string, error) {
	full, err := keypair.ParseFull(secret)
	if err != nil {
		return "", err
	}
	return full.Address(), nil
}

func PubKeyToAddr(publicKey []byte) (string, error) {
	address, err := strkey.Encode(strkey.VersionByteAccountID, publicKey)
	if err != nil {
		return "", err
	}
	return address, nil
}

func GetMuxedAddress(accountID string, id uint64) (string, error) {
	muxedAccount, err := xdr.MuxedAccountFromAccountId(accountID, id)
	if err != nil {
		return "", err
	}
	muxedAddress := muxedAccount.Address()
	if err != nil {
		return "", err
	}
	return muxedAddress, nil
}

func ValidateAddress(address string) error {
	if strings.HasPrefix(address, "G") {
		_, err := strkey.Decode(strkey.VersionByteAccountID, address)
		return err
	} else if strings.HasPrefix(address, "M") {
		_, err := xdr.AddressToMuxedAccount(address)
		return err
	}
	return errors.New("invalid address")
}

func SignTransaction(tx *txnbuild.Transaction, network string, secret ...string) (string, error) {
	kps := make([]*keypair.Full, 0)
	for _, item := range secret {
		full, err := keypair.ParseFull(item)
		if err != nil {
			return "", err
		}
		kps = append(kps, full)
	}
	tx, err := tx.Sign(network, kps...)
	if err != nil {
		return "", err
	}
	return tx.Base64()
}

func TransferAssetTx(sourceAccount txnbuild.Account, toAddress, amount string, asset txnbuild.Asset, baseFee, timeout int64,
	memo txnbuild.Memo) (*txnbuild.Transaction, error) {
	transferAssetOp := txnbuild.Payment{
		Destination: toAddress,
		Amount:      amount,
		Asset:       asset,
	}
	// 创建交易
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		IncrementSequenceNum: true,
		SourceAccount:        sourceAccount, // 源账户的私钥
		Operations:           []txnbuild.Operation{&transferAssetOp},
		Preconditions: txnbuild.Preconditions{
			TimeBounds: getTimeBounds(timeout),
		},
		Memo:    memo,
		BaseFee: baseFee,
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// timeout: 0表示没时间限制
func CreateTrustLineTx(sourceAccount txnbuild.Account, limit string, asset txnbuild.Asset, baseFee, timeout int64, memo txnbuild.Memo) (*txnbuild.Transaction, error) {
	line := txnbuild.ChangeTrustAssetWrapper{
		Asset: asset,
	}
	// 创建信任线操作
	trustLineOp := txnbuild.ChangeTrust{
		Line:  line,
		Limit: limit, // 信任资产的数量上限
	}
	// 创建交易
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		IncrementSequenceNum: true,
		SourceAccount:        sourceAccount, // 源账户的私钥
		Operations:           []txnbuild.Operation{&trustLineOp},
		Preconditions: txnbuild.Preconditions{
			TimeBounds: getTimeBounds(timeout),
		},
		Memo:    memo,
		BaseFee: baseFee,
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func getTimeBounds(timeout int64) txnbuild.TimeBounds {
	if timeout == 0 {
		return txnbuild.NewInfiniteTimeout()
	} else {
		return txnbuild.NewTimeout(timeout)
	}
}
