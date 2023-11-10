package vrf

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	vrfProof "github.com/okx/go-wallet-sdk/coins/oracle/vrf/proof"
	"github.com/shopspring/decimal"
	"github.com/vordev/VOR/core/services/signatures/secp256k1"
	"math/big"
)

type VRFResponse struct {
	Proof      VRFProof             `json:"VRFProof"`
	Commitment VRFRequestCommitment `json:"VRFRequestCommitment"`
}

type VRFRequestCommitment struct {
	BlockNum         uint64         `json:"BlockNum"`
	SubId            uint64         `json:"SubId"`
	CallbackGasLimit uint32         `json:"CallbackGasLimit"`
	NumWords         uint32         `json:"NumWords"`
	Sender           common.Address `json:"Sender"`
}

type VRFProof struct {
	Pk            [2]*big.Int    `json:"PublicKey"`
	Gamma         [2]*big.Int    `json:"Gamma"`
	C             *big.Int       `json:"C"`
	S             *big.Int       `json:"S"`
	Seed          *big.Int       `json:"Seed"`
	UWitness      common.Address `json:"UWitness"`
	CGammaWitness [2]*big.Int    `json:"CGammaWitness"`
	SHashWitness  [2]*big.Int    `json:"SHashWitness"`
	ZInv          *big.Int       `json:"ZInv"`
}

func (v *VRFResponse) GenerateProofResponseFromProof(p vrfProof.Proof, s vrfProof.PreSeedData) error {
	solidityProof, err := vrfProof.SolidityPrecalculations(&p)
	if err != nil {
		return fmt.Errorf("SolidityPrecalculations failed, %v", err)
	}
	solidityProof.P.Seed = common.BytesToHash(s.PreSeed[:]).Big()
	x, y := secp256k1.Coordinates(solidityProof.P.PublicKey)
	gx, gy := secp256k1.Coordinates(solidityProof.P.Gamma)
	cgx, cgy := secp256k1.Coordinates(solidityProof.CGammaWitness)
	shx, shy := secp256k1.Coordinates(solidityProof.SHashWitness)
	v.Proof = VRFProof{
		Pk:            [2]*big.Int{x, y},
		Gamma:         [2]*big.Int{gx, gy},
		C:             solidityProof.P.C,
		S:             solidityProof.P.S,
		Seed:          common.BytesToHash(s.PreSeed[:]).Big(),
		UWitness:      solidityProof.UWitness,
		CGammaWitness: [2]*big.Int{cgx, cgy},
		SHashWitness:  [2]*big.Int{shx, shy},
		ZInv:          solidityProof.ZInv,
	}
	v.Commitment = VRFRequestCommitment{
		BlockNum:         s.BlockNum,
		SubId:            s.SubId,
		CallbackGasLimit: s.CallbackGasLimit,
		NumWords:         s.NumWords,
		Sender:           s.Sender,
	}

	return nil
}

func InitPreSeedData(preSeed, blockHash, sender string, blockNum, subID uint64, cbGasLimit, numWords uint32) (vrfProof.PreSeedData, error) {
	ps, err := vrfProof.BigToSeed(decimal.RequireFromString(preSeed).BigInt())
	if err != nil {
		return vrfProof.PreSeedData{}, fmt.Errorf("init preSeed big to seed failed, %v", err)
	}
	preSeedData := vrfProof.PreSeedData{
		PreSeed:          ps,
		BlockHash:        common.HexToHash(blockHash),
		BlockNum:         blockNum,
		SubId:            subID,
		CallbackGasLimit: cbGasLimit,
		NumWords:         numWords,
		Sender:           common.HexToAddress(sender),
	}
	return preSeedData, nil
}

func GenerateVRFProofResponse(privateKeyHex string, preSeed vrfProof.PreSeedData) (*VRFResponse, error) {
	// init private key
	privateBytes, _ := hex.DecodeString(privateKeyHex)
	key, err := vrfProof.Raw(privateBytes).Key()
	if err != nil {
		return &VRFResponse{}, fmt.Errorf("init private key failed, %v", err)
	}
	// generate final seed
	finalSeed := vrfProof.FinalSeed(preSeed)
	// generate proof
	proof, err := key.GenerateProof(finalSeed)
	if err != nil {
		return &VRFResponse{}, fmt.Errorf("generate proof failed, %v", err)
	}

	vrfResponse := new(VRFResponse)
	// generate proof response
	if err = vrfResponse.GenerateProofResponseFromProof(proof, preSeed); err != nil {
		return &VRFResponse{}, fmt.Errorf("generate proof response from proof failed, %v", err)
	}

	return vrfResponse, nil
}

func MakeVRFProofMarshalResponse(privateKeyHex string, preSeed vrfProof.PreSeedData) (string, error) {
	vrfResp, err := GenerateVRFProofResponse(privateKeyHex, preSeed)
	if err != nil {
		return "", err
	}
	vrfRespBytes, _ := json.Marshal(vrfResp)

	return string(vrfRespBytes), nil
}
