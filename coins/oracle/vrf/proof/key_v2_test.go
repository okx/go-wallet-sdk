package proof

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/go-ethereum/common"
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/vrf/secp256k1"
	"math/big"
	"testing"

	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/curve"
)

func TestVRFKeys_KeyV2_Raw(t *testing.T) {
	privK, err := ecdsa.GenerateKey(curve.S256(), rand.Reader)
	if err != nil {
		t.Error("generate vrf key failed")
	}
	t.Logf("VRF Private Key: %s", hex.EncodeToString(privK.D.Bytes()))

	r := Raw(privK.D.Bytes())
	k, err := r.Key()
	if err != nil {
		t.Errorf("init private key failed, %v", err)
	}
	t.Logf("VRF RAW Private Key: %s", hex.EncodeToString(k.Raw()))
}

func TestVRFKeys_KeyV2(t *testing.T) {
	/*k, err := NewV2()
	if err != nil {
		t.Errorf("NewV2 failed, %v", err)
	}*/

	privKeyBytes, _ := hex.DecodeString("36778dbc3a61764ed00aa1d38ed1ece4eaa830ab675d30483035de06cb1e65b9")
	r := Raw(privKeyBytes)
	k, err := r.Key()
	if err != nil {
		t.Errorf("init private key failed, %v", err)
	}

	t.Logf("VRF RAW Private Key: %s", hex.EncodeToString(k.Raw()))
	t.Logf("VRF Public Key: %s", k.PublicKey.String())
	pkPoint, _ := k.PublicKey.Point()
	pkX, pkY := secp256k1.Coordinates(pkPoint)
	t.Logf("VRF Public Key Coordinates: %v", [2]*big.Int{pkX, pkY})
	uncompressedPK, _ := k.PublicKey.StringUncompressed()
	t.Logf("VRF Uncompressed Public Key: %s", uncompressedPK)
	t.Logf("VRF Address: %v", k.PublicKey.Address())

	t.Run("VRF Public Key SetCoordinates", func(t *testing.T) {
		x, _ := big.NewInt(0).SetString("51571074400993387374180102297480811906841540904676382877334458261957238918398", 10)
		y, _ := big.NewInt(0).SetString("44067233396471738740462061844825239534566793919002732591025918586599711673876", 10)
		rv := secp256k1.EthereumAddress(secp256k1.SetCoordinates(x, y))
		t.Logf("VRF Address: %v", common.BytesToAddress(rv[:]))
	})

	t.Run("generates proof", func(t *testing.T) {
		p, err := k.GenerateProof(big.NewInt(1))
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("PublicKey: %v", hex.EncodeToString(secp256k1.LongMarshal(p.PublicKey)))
			t.Logf("Gamma: %s", hex.EncodeToString(secp256k1.LongMarshal(p.Gamma)))
			t.Logf("C: %s", p.C.String())
			t.Logf("S: %s", p.S.String())
			t.Logf("Seed: %s", p.Seed.String())
			t.Logf("Output: %s", p.Output.String())
		}
	})
}
