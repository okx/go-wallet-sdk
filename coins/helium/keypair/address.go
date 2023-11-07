package keypair

import (
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

type Addressable struct {
	bin       []byte
	base58    string
	publicKey []byte
}

func NewAddressable(address string) *Addressable {
	data := base58.Decode(address)
	//err:=utils.ValidHeliumAddress(address)
	//if err != nil {
	//	return nil
	//}
	bin := data[1 : len(data)-4]
	publicKey := bin[1:]
	aa := new(Addressable)
	aa.base58 = address
	aa.bin = bin
	aa.publicKey = publicKey
	return aa
}

func (aa *Addressable) GetAddress() string {
	if aa == nil {
		return ""
	}
	return aa.base58
}

func (aa *Addressable) GetBin() []byte {
	if aa == nil {
		return nil
	}
	return aa.bin
}

func (aa *Addressable) GetPublicKey() []byte {
	if aa == nil {
		return nil
	}
	return aa.publicKey
}
