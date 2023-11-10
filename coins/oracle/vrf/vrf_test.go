package vrf

import (
	"encoding/hex"
	vrfProof "github.com/okx/go-wallet-sdk/coins/oracle/vrf/proof"
	"github.com/vordev/VOR/core/services/signatures/secp256k1"
	"math/big"
	"testing"
)

func TestInitPreSeedData(t *testing.T) {
	preSeed := "18656632679933127775015571000634247866258735772678922953932585824482961251492"
	blockHash := "0x5f05030c72c506d463c198ccd1cb48f470f61b5eae6520386ec5d9fce596a535"
	sender := "0x690b9a9e9aa1c9db991c7721a92d351db4fac990"
	blockNum := uint64(16116783)
	subID := uint64(1)
	cbGasLimit := uint32(1000)
	numWords := uint32(2)
	t.Run("init preSeed data", func(t *testing.T) {
		p, err := InitPreSeedData(preSeed, blockHash, sender, blockNum, subID, cbGasLimit, numWords)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("PreSeedData: %v", p)
		}
	})
}

func TestVRFResponse_GenerateProofResponseFromProof(t *testing.T) {
	/*k, err := vrfProof.NewV2()
	if err != nil {
		t.Errorf("NewV2 failed, %v", err)
	}
	*/

	privKeyBytes, _ := hex.DecodeString("36778dbc3a61764ed00aa1d38ed1ece4eaa830ab675d30483035de06cb1e65b9")
	r := vrfProof.Raw(privKeyBytes)
	k, err := r.Key()
	if err != nil {
		t.Errorf("init private key failed, %v", err)
	}
	privateKeyHex := hex.EncodeToString(k.Raw())

	t.Logf("VRF RAW Private Key: %s", privateKeyHex)
	t.Logf("VRF Public Key: %s", k.PublicKey.String())
	pkPoint, _ := k.PublicKey.Point()
	pkX, pkY := secp256k1.Coordinates(pkPoint)
	t.Logf("VRF Public Key Coordinates: %v", [2]*big.Int{pkX, pkY})
	uncompressedPK, _ := k.PublicKey.StringUncompressed()
	t.Logf("VRF Uncompressed Public Key: %s", uncompressedPK)
	t.Logf("VRF Address: %v", k.PublicKey.Address())

	preSeed := "18656632679933127775015571000634247866258735772678922953932585824482961251492"
	blockHash := "0x5f05030c72c506d463c198ccd1cb48f470f61b5eae6520386ec5d9fce596a535"
	sender := "0x690b9a9e9aa1c9db991c7721a92d351db4fac990"
	blockNum := uint64(16116783)
	subID := uint64(1)
	cbGasLimit := uint32(1000)
	numWords := uint32(2)

	psData, err := InitPreSeedData(preSeed, blockHash, sender, blockNum, subID, cbGasLimit, numWords)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("PreSeedData: %v", psData)
	}

	t.Run("generates vrf proof response", func(t *testing.T) {
		vrfResp, err := GenerateVRFProofResponse(privateKeyHex, psData)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("VRF Proof: %v", vrfResp.Proof)
			t.Logf("VRF RequestCommitment: %v", vrfResp.Commitment)
		}
	})

	t.Run("make vrf proof marshal response", func(t *testing.T) {
		vrfRespStr, err := MakeVRFProofMarshalResponse(privateKeyHex, psData)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("VRF Marshal Response: %s", vrfRespStr)
		}
	})
}
