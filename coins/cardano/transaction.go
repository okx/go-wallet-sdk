package cardano

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func buildMultiAsset(multiAsset []*MultiAssets) (*MultiAsset, error) {
	ma := NewMultiAsset()

	for _, a := range multiAsset {
		policyId, err := NewPolicyIDFromHex(a.PolicyId)
		if err != nil {
			return nil, err
		}
		for _, asset := range a.Assets {
			assetName, err := NewAssetNameFromHex(asset.AssetName)
			if err != nil {
				return nil, err
			}
			if ma.Get(policyId) != nil {
				ma.Get(policyId).Set(assetName, BigNum(asset.Amount))
			} else {
				assets := NewAssets()
				assets.Set(assetName, BigNum(asset.Amount))

				ma.Set(policyId, assets)
			}
		}
	}

	return ma, nil
}

func CreateTxBuilder(txData *TxData) (*TxBuilder, error) {
	txBuilder := NewTxBuilder(protocolParams)

	for _, input := range txData.Inputs {
		txHash, err := hex.DecodeString(input.TxId)
		if err != nil {
			return nil, err
		}
		assets, err := buildMultiAsset(input.MultiAsset)
		if err != nil {
			return nil, err
		}
		amount := NewValueWithAssets(Coin(input.Amount), assets)
		txBuilder.AddInputs(&TxInput{
			TxHash: txHash,
			Index:  input.Index,
			Amount: amount,
		})
	}

	if txData.ToAddress == "" {
		return nil, fmt.Errorf("toAddress is required")
	}

	toAddress, err := NewAddress(txData.ToAddress)
	if err != nil {
		return nil, err
	}

	assets, err := buildMultiAsset(txData.MultiAsset)
	if err != nil {
		return nil, err
	}
	if !txData.Max {
		minAda, err := MinAda(txData.ToAddress, txData.MultiAsset)
		if err != nil {
			return nil, err
		}
		if txData.Amount < minAda {
			return nil, fmt.Errorf("amount is less than minAda required for output: %d", minAda)
		}
	}

	toValue := NewValueWithAssets(Coin(txData.Amount), assets)
	txBuilder.AddOutputs(&TxOutput{
		Address: toAddress,
		Amount:  toValue,
	})

	if txData.Max {
		txBuilder.SetMax()
	}

	changeAddress, err := NewAddress(txData.ChangeAddress)
	if err != nil {
		return nil, err
	}
	txBuilder.AddChangeIfNeeded(changeAddress)

	if txData.TTL > 0 {
		txBuilder.SetTTL(txData.TTL)
	}

	return txBuilder, nil
}

func Transfer(txData *TxData) (string, error) {
	txBuilder, err := CreateTxBuilder(txData)
	if err != nil {
		return "", err
	}

	prvKey, err := hex.DecodeString(txData.PrvKey)
	if err != nil {
		return "", err
	}

	tx, err := txBuilder.Build(prvKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(tx.Bytes()), nil
}

func MinAda(addrStr string, multiAsset []*MultiAssets) (uint64, error) {
	address, err := NewAddress(addrStr)
	if err != nil {
		return 0, err
	}
	assets, err := buildMultiAsset(multiAsset)
	if err != nil {
		return 0, err
	}
	return uint64(calcMinAda(address, assets)), nil
}

func MinFee(txData *TxData) (*MinFeeData, error) {
	txBuilder, err := CreateTxBuilder(txData)
	if err != nil {
		return nil, err
	}

	valid, err := txBuilder.BuildUnsigned()
	if err != nil {
		return nil, err
	}

	fee, err := txBuilder.Fee()

	if err != nil {
		return nil, err
	}

	change := txBuilder.GetChange()

	return &MinFeeData{
		Valid:  valid,
		Fee:    uint64(fee),
		Change: change,
	}, nil
}

func GetTxHash(transaction string) (string, error) {
	var tx Tx
	b, err := base64.StdEncoding.DecodeString(transaction)
	if err != nil {
		return "", err
	}

	err = tx.UnmarshalCBOR(b)
	if err != nil {
		return "", err
	}

	h, err := tx.Hash()
	if err != nil {
		return "", err
	}

	return h.String(), nil
}
