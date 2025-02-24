package bitcoin

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

type TaprootInfo struct {
	ControlBlockWitness string `json:"controlBlockWitness"`
	TaprootAddress      string `json:"taprootAddress"`
}

func NewTaprootAddress(script string, net *chaincfg.Params, pubHex string) (*TaprootInfo, error) {
	if net == nil {
		net = &chaincfg.MainNetParams
	}
	inscriptionScript, err := hex.DecodeString(script)
	if err != nil {
		return nil, err
	}
	pubBytes, err := hex.DecodeString(pubHex)
	if err != nil {
		return nil, err
	}
	pubKey, err := btcec.ParsePubKey(pubBytes)
	if err != nil {
		return nil, err
	}
	leafNode := txscript.NewBaseTapLeaf(inscriptionScript)
	proof := &txscript.TapscriptProof{
		TapLeaf:  leafNode,
		RootNode: leafNode,
	}
	controlBlock := proof.ToControlBlock(pubKey)
	controlBlockWitness, err := controlBlock.ToBytes()
	if err != nil {
		return nil, err
	}
	tapHash := proof.RootNode.TapHash()
	commitTxAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(pubKey, tapHash[:])), net)
	if err != nil {
		return nil, err

	}
	return &TaprootInfo{
		ControlBlockWitness: hex.EncodeToString(controlBlockWitness),
		TaprootAddress:      commitTxAddress.String(),
	}, nil
}
