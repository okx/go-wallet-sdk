package nostrassets

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/okx/go-wallet-sdk/util"
	"strings"
)

const (
	// pubkeyCompressed is the header byte for a compressed secp256k1 pubkey.
	PUBKEY_COMPRESSED byte = 0x2 // y_bit + x coord

	HRP  = "npub"
	NSec = "nsec"
)

func GetPublicKey(prvBech string) (pub string, err error) {
	defer func() {
		if r := recover(); r != nil {
			pub, err = "", errors.New("invalid public key")
			return
		}
	}()
	prv, err := DecodeBech32(NSec, prvBech)
	if err != nil {
		return "", err
	}
	pk, _ := btcec.PrivKeyFromBytes(prv)
	return util.EncodeHex(pk.PubKey().SerializeCompressed()[1:]), nil
}

func NpubEncode(pubHex string) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			res, err = "", errors.New("invalid public key")
			return
		}
	}()
	pub, err := util.DecodeHexString(pubHex)
	if err != nil {
		return "", err
	}
	return bech32.EncodeFromBase256(HRP, pub)
}
func NsecEncode(prvHex string) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			res, err = "", errors.New("invalid private key")
			return
		}
	}()
	pub, err := util.DecodeHexString(prvHex)
	if err != nil {
		return "", err
	}
	return bech32.EncodeFromBase256(NSec, pub)
}

func DecodeBech32(prefix, data string) ([]byte, error) {
	hrp, res, err := bech32.DecodeToBase256(data)
	if err != nil {
		return nil, err
	}
	if hrp != prefix {
		return nil, errors.New("invalid bech32 data")
	}
	return res, nil
}

func ValidateAddress(address string) bool {
	hrp, _, err := bech32.DecodeToBase256(address)
	if err != nil {
		return false
	}
	return hrp == HRP
}

func AddressFromPrvKey(prvBech string) (addr string, err error) {
	defer func() {
		if r := recover(); r != nil {
			addr, err = "", errors.New("invalid public key")
			return
		}
	}()
	pv, err := DecodeBech32(NSec, prvBech)
	if err != nil {
		return "", err
	}
	pk, _ := btcec.PrivKeyFromBytes(pv)
	return bech32.EncodeFromBase256(HRP, pk.PubKey().SerializeCompressed()[1:])
}

type Event struct {
	Kind      uint64     `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	CreatedAt uint64     `json:"created_at"`
	Pubkey    string     `json:"pubkey"`
	Id        string     `json:"id"`
	Sig       string     `json:"sig"`
}

func (e *Event) Copy() *Event {
	ee := &Event{Kind: e.Kind, Content: e.Content, CreatedAt: e.CreatedAt, Pubkey: e.Pubkey, Id: e.Id, Sig: e.Sig}
	ee.Tags = make([][]string, len(e.Tags))
	for k, v := range e.Tags {
		vv := make([]string, len(v))
		copy(vv, v)
		ee.Tags[k] = vv
	}
	return ee
}

func SignEvent(prvBech string, evt *Event) (res *Event, err error) {
	defer func() {
		if r := recover(); r != nil {
			res, err = nil, r.(error)
			return
		}
	}()
	e := evt.Copy()
	pub, err := GetPublicKey(prvBech)
	if err != nil {
		return nil, err
	}
	e.Pubkey = pub
	body, err := json.Marshal([]interface{}{0, e.Pubkey, e.CreatedAt, e.Kind, e.Tags, e.Content})
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	h.Write(body)
	hash := h.Sum(nil)
	eventHash := util.EncodeHex(hash)
	e.Id = eventHash
	pv, err := DecodeBech32(NSec, prvBech)
	if err != nil {
		return nil, err
	}
	prv, _ := btcec.PrivKeyFromBytes(pv)
	s, err := schnorr.Sign(prv, hash)
	if err != nil {
		return nil, err
	}
	e.Sig = hex.EncodeToString(s.Serialize())
	return e, nil
}

func GetEventHash(e *Event) (string, error) {
	body, err := json.Marshal([]interface{}{0, e.Pubkey, e.CreatedAt, e.Kind, e.Tags, e.Content})
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(body)
	hash := h.Sum(nil)
	eventHash := util.EncodeHex(hash)
	return eventHash, nil
}

func VerifyEvent(e *Event) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			return
		}
	}()
	if len(e.Pubkey) == 0 || len(e.Id) == 0 || len(e.Sig) == 0 {
		return false
	}
	body, err := json.Marshal([]interface{}{0, e.Pubkey, e.CreatedAt, e.Kind, e.Tags, e.Content})
	if err != nil {
		return false
	}
	h := sha256.New()
	h.Write(body)
	hash := h.Sum(nil)
	eventHash := util.EncodeHex(hash)
	if e.Id != eventHash {
		return false
	}

	sigBytes, err := util.DecodeHexString(e.Sig)
	if err != nil {
		return false
	}
	s, err := schnorr.ParseSignature(sigBytes)
	if err != nil {
		return false
	}
	pubKeyBytes, err := util.DecodeHexString(e.Pubkey)
	if err != nil {
		return false
	}
	p, err := schnorr.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false
	}
	return s.Verify(hash, p)
}

func TryGetPubkey(pubkeyHex string) []byte {
	pub, err := util.DecodeHexString(pubkeyHex)
	if err != nil {
		return make([]byte, 0)
	}
	full := make([]byte, len(pub)+1)
	full[0] = secp256k1.PubKeyFormatCompressedOdd
	format := full[0]
	format &= ^byte(0x1)
	copy(full[1:], pub)
	if format != PUBKEY_COMPRESSED {
		full[0] = secp256k1.PubKeyFormatCompressedEven
	}
	return full
}

func Encrypt(prvBech string, pubkey string, text string, fakeIv string) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			res, err = "", r.(error)
			return
		}
	}()
	iv := make([]byte, 16)
	if len(fakeIv) > 0 {
		rr, err := base64.StdEncoding.DecodeString("RK8MdoXLJLe0R4DbYAd2oQ==")
		if err != nil {
			return "", err
		}
		copy(iv[:], rr)
	} else {

		n, err := rand.Read(iv)
		if n != 16 {
			return "", errors.New("invalid rand")
		}
		if err != nil {
			return "", err
		}
	}
	pv, err := DecodeBech32(NSec, prvBech)
	if err != nil {
		return "", err
	}
	prv, _ := btcec.PrivKeyFromBytes(pv)
	pub, err := btcec.ParsePubKey(TryGetPubkey(pubkey))
	if err != nil {
		return "", err
	}
	secret := secp256k1.GenerateSharedSecret(prv, pub)
	clip, err := AseEncrypt(text, secret, iv, aes.BlockSize)
	if err != nil {
		return "", err
	}
	res = fmt.Sprintf("%s?iv=%s", base64.StdEncoding.EncodeToString(clip), base64.StdEncoding.EncodeToString(iv))
	return
}

func Decrypt(prvBech string, pubkey string, ciphertext string) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			res, err = "", r.(error)
			return
		}
	}()
	pv, err := DecodeBech32(NSec, prvBech)
	if err != nil {
		return "", err
	}
	params := strings.Split(ciphertext, "?iv=")
	if len(params) != 2 {
		return "", errors.New("invalid ciphertext")
	}
	ctb, err := base64.StdEncoding.DecodeString(params[0])
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(params[1])
	if err != nil {
		return "", err
	}
	prv, _ := btcec.PrivKeyFromBytes(pv)
	pub, err := btcec.ParsePubKey(TryGetPubkey(pubkey))
	if err != nil {
		return "", err
	}
	secret := secp256k1.GenerateSharedSecret(prv, pub)
	r, err := AseDecrypt(ctb, secret, iv)
	if err != nil {
		return "", err
	}
	res = string(r)
	return
}

func AseEncrypt(plaintext string, bKey []byte, bIV []byte, blockSize int) ([]byte, error) {
	bPlaintext := PKCS5Padding([]byte(plaintext), blockSize, len(plaintext))
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return ciphertext, nil
}

func AseDecrypt(ciphertext []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	ecb := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	ecb.CryptBlocks(decrypted, ciphertext)
	return PKCS5Trimming(decrypted), nil
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
