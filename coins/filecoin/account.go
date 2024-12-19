package filecoin

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	ecdsa2 "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/dchest/blake2b"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/okx/go-wallet-sdk/util"
)

var (
	payloadHashConfig  = &blake2b.Config{Size: 20}
	checksumHashConfig = &blake2b.Config{Size: 4}
	AddressEncoding    = base32.NewEncoding(encodeStd)
)

const (
	MainnetPrefix = "f"
	TestnetPrefix = "t"
	encodeStd     = "abcdefghijklmnopqrstuvwxyz234567"
)

const (
	ID byte = iota
	SECP256K1
	Actor
	BLS
)

func NewPrivateKey() string {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return ""
	}
	privkey := make([]byte, 32)
	blob := key.D.Bytes()
	copy(privkey[32-len(blob):], blob)
	return util.EncodeHexWith0x(privkey)
}

func GetPublicKey(privateKeyHex string) (string, error) {
	privateKeyBytes, err := util.DecodeHexString(privateKeyHex)
	if err != nil {
		return "", err
	}
	x, y := secp256k1.S256().ScalarBaseMult(privateKeyBytes)
	publicKeyBytes := elliptic.Marshal(secp256k1.S256(), x, y)
	return util.EncodeHexWith0x(publicKeyBytes), nil
}

func GetAddressByPublicKey(publicKeyHex string, chainId string) (string, error) {
	if strings.HasPrefix(publicKeyHex, "0x") || strings.HasPrefix(publicKeyHex, "0X") {
		publicKeyHex = publicKeyHex[2:]
	}
	pubKeyByte, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", err
	}
	pk, err := btcec.ParsePubKey(pubKeyByte)
	if err != nil {
		return "", err
	}
	uncompressedPubKey := pk.SerializeUncompressed()
	pubKeyHash, err := hash(uncompressedPubKey, payloadHashConfig)
	if err != nil {
		return "", err
	}

	explen := 1 + len(pubKeyHash)
	buf := make([]byte, explen)
	var protocol byte = 1
	buf[0] = protocol
	copy(buf[1:], pubKeyHash)

	cksm, err := hash(buf, checksumHashConfig)
	if err != nil {
		return "", err
	}
	address := chainId + fmt.Sprintf("%d", protocol) + AddressEncoding.WithPadding(-1).EncodeToString(append(pubKeyHash, cksm[:]...))

	return address, nil
}

func GetAddressByPrivateKey(privateKeyHex string, chainId string) (string, error) {
	publicKeyHex, err := GetPublicKey(privateKeyHex)
	if err != nil {
		return "", err
	}
	return GetAddressByPublicKey(publicKeyHex, chainId)
}

func AddressToBytes(addr string) []byte {
	if len(addr) == 0 {
		return nil
	}

	if string(addr[0]) != MainnetPrefix && string(addr[0]) != TestnetPrefix {
		return nil
	}

	var protocol byte
	switch addr[1] {
	case '0':
		protocol = ID
	case '1':
		protocol = SECP256K1
	case '2':
		protocol = Actor
	case '3':
		protocol = BLS
	default:
		return nil
	}

	raw := addr[2:]
	if protocol == ID {
		if len(raw) > 20 {
			return nil
		}
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil
		}
		return toBytes(protocol, toUvarint(id))
	}

	payloadcksm, err := AddressEncoding.WithPadding(-1).DecodeString(raw)
	if err != nil {
		return nil
	}
	payload := payloadcksm[:len(payloadcksm)-4]
	cksm := payloadcksm[len(payloadcksm)-4:]

	if protocol == SECP256K1 || protocol == Actor {
		if len(payload) != 20 {
			return nil
		}
	}

	if !validateChecksum(append([]byte{protocol}, payload...), cksm) {
		return nil
	}

	return toBytes(protocol, payload)
}

func SignTx(message *Message, privateKeyHex string) (*SignedMessage, error) {
	privKeyBytes, err := util.DecodeHexString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	sig := ecdsa2.SignCompact(privateKey, message.Hash(), false)
	V := sig[0]
	R := sig[1:33]
	S := sig[33:65]
	signature := append(R, S...)
	signature = append(signature, V-27)

	return &SignedMessage{
		Message: message,
		Signature: Signature{
			Type: SECP256K1,
			Data: signature,
		},
	}, nil
}

func toBytes(protocol byte, payload []byte) []byte {
	switch protocol {
	case ID:
		_, n, err := fromUvarint(payload)
		if err != nil {
			return nil
		}
		if n != len(payload) {
			return nil
		}
	case SECP256K1, Actor:
		if len(payload) != 20 {
			return nil
		}
	case BLS:
		if len(payload) != 48 {
			return nil
		}
	default:
		return nil
	}
	explen := 1 + len(payload)
	buf := make([]byte, explen)

	buf[0] = protocol
	copy(buf[1:], payload)

	return buf
}

// fromUvarint reads an unsigned varint from the beginning of buf, returns the
// varint, and the number of bytes read.
func fromUvarint(buf []byte) (uint64, int, error) {
	// Modified from the go standard library. Copyright the Go Authors and
	// released under the BSD License.
	var x uint64
	var s uint
	for i, b := range buf {
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, 0, errors.New("varints larger than uint64 not yet supported")
			} else if b == 0 && s > 0 {
				return 0, 0, errors.New("varint not minimally encoded")
			}
			return x | uint64(b)<<s, i + 1, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0, errors.New("varints malformed, could not reach the end")
}

func toUvarint(num uint64) []byte {
	buf := make([]byte, uvarintSize(num))
	n := binary.PutUvarint(buf, uint64(num))
	return buf[:n]
}

func uvarintSize(num uint64) int {
	bits := bits.Len64(num)
	q, r := bits/7, bits%7
	size := q
	if r > 0 || size == 0 {
		size++
	}
	return size
}

func validateChecksum(ingest, expect []byte) bool {
	digest, err := hash(ingest, checksumHashConfig)
	if err != nil {
		return false
	}
	return bytes.Equal(digest, expect)
}

func hash(ingest []byte, cfg *blake2b.Config) ([]byte, error) {
	hasher, err := blake2b.New(cfg)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid address hash configuration: %v", err))
	}
	if _, err := hasher.Write(ingest); err != nil {
		return nil, errors.New(fmt.Sprintf("blake2b is unable to process hashes: %v", err))
	}
	return hasher.Sum(nil), nil
}
