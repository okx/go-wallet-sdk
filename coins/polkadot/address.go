package polkadot

import (
	"encoding/hex"

	"github.com/okx/go-wallet-sdk/crypto/ss58"
)

var (
	SubstratePrefix  = []byte{0x2a}
	PolkadotPrefix   = []byte{0x00}
	KsmPrefix        = []byte{0x02}
	DarwiniaPrefix   = []byte{0x12}
	EdgewarePrefix   = []byte{0x07}
	CentrifugePrefix = []byte{0x24}
	PlasmPrefix      = []byte{0x05}
	StafiPrefix      = []byte{0x14}
	KulupuPrefix     = []byte{0x10}
	BifrostPrefix    = []byte{0x06}
	KaruraPrefix     = []byte{0x08}
	ReynoldsPrefix   = []byte{0x09}
	AcalaPrefix      = []byte{0x0a}
	LaminarPrefix    = []byte{0x0b}
	PolymathPrefix   = []byte{0x0c}
	RobonomicsPrefix = []byte{0x20}
	ChainxPrefix     = []byte{0x2c}
	SSPrefix         = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
)

func PubKeyToAddress(publicKey []byte, prefix []byte) (address string, err error) {
	address, err = ss58.Encode(publicKey, prefix)
	if err != nil {
		return "", err
	}
	return address, nil
}

func AddressToPublicKey(address string) string {
	pub, _ := ss58.DecodeToPub(address)
	pubHex := hex.EncodeToString(pub)
	return pubHex
}

func ValidateAddress(addr string) bool {
	publicKey, err := ss58.DecodeToPub(addr)
	if err != nil {
		return false
	}
	if len(publicKey) != 32 {
		return false
	}
	return true
}
