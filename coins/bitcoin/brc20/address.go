package brc20

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func NewTapRootAddress(privateKey *btcec.PrivateKey, params *chaincfg.Params) (string, error) {
	rootKey := txscript.ComputeTaprootKeyNoScript(privateKey.PubKey())
	address, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(rootKey), params)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func NewTapRootAddressWithScript(privateKey *btcec.PrivateKey, script []byte, params *chaincfg.Params) (string, error) {
	proof := &txscript.TapscriptProof{
		TapLeaf:  txscript.NewBaseTapLeaf(schnorr.SerializePubKey(privateKey.PubKey())),
		RootNode: txscript.NewBaseTapLeaf(script),
	}
	tapHash := proof.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(privateKey.PubKey(), tapHash[:])
	address, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(outputKey), params)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func NewTapRootAddressWithScriptWithPubKey(serializedPubKey []byte, script []byte, params *chaincfg.Params) string {
	proof := &txscript.TapscriptProof{
		TapLeaf:  txscript.NewBaseTapLeaf(serializedPubKey),
		RootNode: txscript.NewBaseTapLeaf(script),
	}
	tapHash := proof.RootNode.TapHash()
	pubKey, err := schnorr.ParsePubKey(serializedPubKey)
	if err != nil {
		return ""
	}
	outputKey := txscript.ComputeTaprootOutputKey(pubKey, tapHash[:])
	address, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(outputKey), params)
	if err != nil {
		return ""
	}
	return address.String()
}

func CreateInscriptionScript(privateKey *btcec.PrivateKey, contentType string, body []byte) ([]byte, error) {
	inscriptionBuilder := txscript.NewScriptBuilder().
		AddData(schnorr.SerializePubKey(privateKey.PubKey())).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_FALSE).
		AddOp(txscript.OP_IF).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		// text/plain;charset=utf-8
		AddData([]byte(contentType)).
		AddOp(txscript.OP_0)

	maxChunkSize := 520
	bodySize := len(body)
	for i := 0; i < bodySize; i += maxChunkSize {
		end := i + maxChunkSize
		if end > bodySize {
			end = bodySize
		}
		inscriptionBuilder.AddFullData(body[i:end])
	}
	inscriptionScript, err := inscriptionBuilder.Script()
	if err != nil {
		return nil, err
	}
	// to skip txscript.MaxScriptSize 10000
	inscriptionScript = append(inscriptionScript, txscript.OP_ENDIF)
	return inscriptionScript, nil
}

func CreateInscriptionScriptWithPubKey(publicKey []byte, contentType string, body []byte) ([]byte, error) {
	inscriptionBuilder := txscript.NewScriptBuilder().
		AddData(publicKey).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_FALSE).
		AddOp(txscript.OP_IF).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		// text/plain;charset=utf-8
		AddData([]byte(contentType)).
		AddOp(txscript.OP_0)

	maxChunkSize := 520
	bodySize := len(body)
	for i := 0; i < bodySize; i += maxChunkSize {
		end := i + maxChunkSize
		if end > bodySize {
			end = bodySize
		}
		inscriptionBuilder.AddFullData(body[i:end])
	}
	inscriptionScript, err := inscriptionBuilder.Script()
	if err != nil {
		return nil, err
	}
	// to skip txscript.MaxScriptSize 10000
	inscriptionScript = append(inscriptionScript, txscript.OP_ENDIF)
	return inscriptionScript, nil
}
