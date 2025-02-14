package ton

import (
	"bytes"
	"crypto/ed25519"
	"errors"
	"github.com/okx/go-wallet-sdk/util"
	"strings"

	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
)

var (
	ErrInvalidMnemonic = errors.New("invalid mnemonic")
	ErrInvalidVersion  = errors.New("invalid wallet version")
)

func GetVersionConfigAndSubWallet(version wallet.Version) (wallet.VersionConfig, uint32, error) {
	var versionConfig wallet.VersionConfig
	var subWallet uint32
	if version == wallet.V4R2 {
		versionConfig = wallet.V4R2
		subWallet = wallet.DefaultSubwallet
	} else if version == wallet.V5R1Final {
		versionConfig = wallet.ConfigV5R1Final{NetworkGlobalID: wallet.MainnetGlobalID}
		subWallet = 0
	} else {
		return versionConfig, subWallet, errors.New("invalid version")
	}
	return versionConfig, subWallet, nil
}

func NewWallet(seed, pubKey []byte, version wallet.Version) (*wallet.Wallet, error) {
	versionConfig, _, err := GetVersionConfigAndSubWallet(version)
	if err != nil {
		return nil, err
	}

	if len(pubKey) == ed25519.PublicKeySize && len(seed) == ed25519.SeedSize {
		if bytes.Equal(pubKey, ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)) {
			return wallet.FromPrivateKey(ed25519.NewKeyFromSeed(seed), versionConfig)
		}
	}
	if len(pubKey) > 0 {
		return wallet.FakeFromPublicKey(ed25519.PublicKey(pubKey), versionConfig)
	}
	return wallet.FromPrivateKey(ed25519.NewKeyFromSeed(seed), versionConfig)
}

func NewAddress(seed []byte, version wallet.Version) (string, error) {
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	versionConfig, subWallet, err := GetVersionConfigAndSubWallet(version)
	if err != nil {
		return "", err
	}

	addr, err := wallet.AddressFromPubKey(pubKey, versionConfig, subWallet)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func CheckPubKeyAddress(pubKeyHex string, addr string) error {
	b, err := util.DecodeHexString(pubKeyHex)
	if err != nil {
		return err
	}
	if len(b) != ed25519.PublicKeySize {
		return errors.New("invalid public key size")
	}
	pubKey := ed25519.PublicKey(b)
	a, err := wallet.AddressFromPubKey(pubKey, wallet.V4R2, wallet.DefaultSubwallet)
	if err != nil {
		return err
	}
	if a.Equal(addr) {
		return nil
	}
	a2, err := wallet.AddressFromPubKey(pubKey, wallet.V5R1Final, wallet.DefaultSubwallet)
	if a2.Equal(addr) {
		return nil
	}
	return errors.New("invalid address")
}

func NewPubKeyAddress(pubKeyHex string, version string) (string, error) {
	b, err := util.DecodeHexString(pubKeyHex)
	if err != nil {
		return "", err
	}
	if len(b) != ed25519.PublicKeySize {
		return "", errors.New("invalid public key size")
	}
	pubKey := ed25519.PublicKey(b)
	switch version {
	case "v4r2":
		fallthrough
	case "":
		a, err := wallet.AddressFromPubKey(pubKey, wallet.V4R2, wallet.DefaultSubwallet)
		if err != nil {
			return "", err
		}
		return a.String(), nil
	case "w5":
		a, err := wallet.AddressFromPubKey(pubKey, wallet.V5R1Final, wallet.DefaultSubwallet)
		if err != nil {
			return "", err
		}
		return a.String(), nil
	default:
		return "", ErrInvalidVersion
	}
	return "", ErrInvalidVersion
}
func AddressStrings(a string) ([]*address.AddrWithType, error) {
	addr, err := address.ParseAddr(a)
	if err != nil {
		return nil, err
	}
	return addr.Strings(), nil
}

func FromSeedV4R2(mnemonic, password string) (ed25519.PrivateKey, error) {
	if len(mnemonic) == 0 {
		return nil, ErrInvalidMnemonic
	}
	words := strings.Split(mnemonic, " ")
	if len(words) != 24 {
		return nil, ErrInvalidMnemonic
	}
	w, err := wallet.FromSeedWithPassword(words, password, wallet.V4R2)
	if err != nil {
		return nil, err
	}
	return w.PrivateKey(), nil
}

func ValidateAddress(addr string) bool {
	_, err := address.ParseAddr(addr)
	if err != nil {
		return false
	}
	return true
}

func VenomNewAddress(seed []byte) (string, error) {
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	addr, err := wallet.AddressFromPubKey(pubKey, wallet.VenomV3, wallet.VenomDefaultSubwallet)
	if err != nil {
		return "", err
	}

	return addr.RawString(), nil
}
