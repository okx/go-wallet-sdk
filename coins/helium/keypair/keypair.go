/**
Author： https://github.com/hecodev007/block_sign
*/

package keypair

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	cre "github.com/okx/go-wallet-sdk/coins/helium/crypto"
	"github.com/okx/go-wallet-sdk/coins/helium/utils"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

const (
	NISTP256Version = iota
	Ed25519Version
	WIFVersion = 0x80
)

type Keypair struct {
	curve cre.Curves

	privateKey []byte
	version    int
}

func New(version int) *Keypair {
	c := cre.NewCurve(version)
	ks := new(Keypair)
	ks.curve = c
	ks.version = version
	return ks
}

func NewKeypairFromHex(version int, privHex string) *Keypair {
	kp := New(version)
	data, err := hex.DecodeString(privHex)
	if err != nil {
		return nil
	}
	kp.privateKey = data
	return kp
}
func (kp *Keypair) CreatePublicFromPrivate(private string) string {
	seedBytes, _ := hex.DecodeString(private)
	privateKey := ed25519.NewKeyFromSeed(seedBytes)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	address := kp.CreateAddress(publicKey)

	return address
}
func (kp *Keypair) GenerateKey() ([]byte, []byte) {
	return kp.curve.GenerateKey()
}

func (kp *Keypair) CreateAddressable() *Addressable {
	if kp.privateKey == nil {
		return nil
	}
	priv := ed25519.NewKeyFromSeed(kp.privateKey)
	pub := make([]byte, 32)
	copy(pub, priv[32:])
	address := kp.CreateAddress(pub)
	var bin []byte
	v := kp.curve.GetVersion()
	bin = append(bin, v...)
	bin = append(bin, pub...)
	aa := new(Addressable)
	aa.base58 = address
	aa.bin = bin
	aa.publicKey = pub
	return aa
}

// CreateAddress first byte 0, second byte is network | keyType, next is public key and then base64 for all
func (kp *Keypair) CreateAddress(publicKey []byte) string {
	var (
		payload  []byte
		vpayload []byte
	)
	v := kp.curve.GetVersion()           // 1->ed25519 0-> NIST p256
	payload = append(v, publicKey[:]...) //[0,1,pub][1,pub..]
	version := []byte{0}                 //mainNet 0, testNet 0x10
	vpayload = append(version, payload...)
	//vpayload[1] = 17 //testNet
	//double sha256
	checksum := utils.DoubleSha256(vpayload)[:4]
	vpayload = append(vpayload, checksum...)

	return base58.Encode(vpayload)
}
func (kp *Keypair) SetPrivateKey(privateKey []byte) {
	kp.privateKey = privateKey
	return
}

func (kp *Keypair) Sign(message []byte) ([]byte, error) {
	if kp.privateKey == nil {
		return nil, errors.New("private key is null")
	}
	var data []byte
	if kp.version == 1 {
		privKey := ed25519.NewKeyFromSeed(kp.privateKey)
		data = ed25519.Sign(privKey, message)
	} else {
		//todo
		//privKey:=nist_p256.NewNISTP256PrivateBySeed(kp.privateKey)
		//ecdsa.Sign(rand.Reader,privKey,message)
	}
	return data, nil
}
