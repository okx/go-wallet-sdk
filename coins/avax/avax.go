package avax

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

func NewAddress(chainId string, hrp string, publicKey *btcec.PublicKey) (string, error) {
	hash := btcutil.Hash160(publicKey.SerializeCompressed())
	addr, err := bech32.EncodeFromBase256(hrp, hash)
	if err != nil {
		return "", err
	}
	return chainId + "-" + addr, nil
}

// ParseAddress
// returned value: (address, chainId, hrp, error)
func ParseAddress(address string) ([]byte, string, string, error) {
	arr := strings.Split(address, "-")
	if len(arr) != 2 {
		return nil, "", "", errors.New("invalid address")
	}
	hrp, data, err := bech32.Decode(arr[1])
	if err != nil {
		return nil, "", "", err
	}
	add, err := bech32.ConvertBits(data, 5, 8, false)
	return add, arr[0], hrp, err
}

func NewTransferTransaction(netWorkId uint32, blockchainId string, inputs *[]TransferInput, outputs *[]TransferOutPut) (string, error) {
	sl := NewSerializer()
	blockchainIDBytes, err := CheckDecodeWithCheckSumLast(blockchainId)
	if err != nil {
		return "", err
	}
	trans := Transaction{Codecid: 0, TypeId: BASETX, NetworkID: netWorkId, BlockchainID: blockchainIDBytes}
	for _, output := range *outputs {
		if err := trans.AddOutput(output.Address, output.Value, output.AssetId); err != nil {
			return "", err
		}
	}

	for _, input := range *inputs {
		if err := trans.AddInput(input.TxId, input.Index, input.Amount, input.AssetId, input.PrivateKey); err != nil {
			return "", err
		}
	}

	// sort input and output
	trans.sortInputAndOutput()
	trans.SerializeToBytes(sl)
	hash := sha256.Sum256(sl.Payload())

	// sign
	for _, input := range trans.Ins {
		pk, err := hex.DecodeString(input.PrivateKey)
		if err != nil {
			return "", err
		}
		privateKey, _ := btcec.PrivKeyFromBytes(pk)
		sig2 := ecdsa.SignCompact(privateKey, hash[:], false)
		sig := make([]byte, len(sig2))
		copy(sig, sig2[1:])
		sig[64] = sig2[0] - 27
		trans.Credentials = append(trans.Credentials, Credential{sig, SECPCREDENTIAL})
	}

	// Re-serialize after signing
	sl = NewSerializer()
	trans.SerializeToBytes(sl)
	return C58Encode(sl.Payload()), nil
}

func C58Encode(input []byte) string {
	b := make([]byte, 0, len(input)+4)
	b = append(b, input...)
	h := sha256.Sum256(input)
	b = append(b, h[len(h)-4:]...)
	return base58.Encode(b)
}

// CheckDecodeWithCheckSumLast decodes a string that was encoded with CheckEncode and verifies the checksum.
func CheckDecodeWithCheckSumLast(input string) (result []byte, err error) {
	decoded := base58.Decode(input)
	if len(decoded) < 5 {
		return nil, errors.New("invalid format")
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])

	var cksum2 [4]byte
	h := sha256.Sum256(decoded[:len(decoded)-4])
	copy(cksum2[:], h[len(h)-4:])

	if cksum2 != cksum {
		return nil, errors.New("invalid check sum")
	}
	payload := decoded[0 : len(decoded)-4]
	result = append(result, payload...)
	return
}
