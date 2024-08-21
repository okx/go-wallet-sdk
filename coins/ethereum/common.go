package ethereum

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
)

const AddressLength = 20

func OnlyRemovePrefix(s string) string {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			return s[2:]
		}
	}
	return s
}

func encodeRSV(r, s, v *big.Int) []byte {
	sig := make([]byte, 65)
	copy(sig[0:32], r.Bytes())
	copy(sig[32:64], s.Bytes())
	sig[64] = byte(v.Uint64() - 27)
	return sig
}

func getEthGroupPubHash(pubKey *btcec.PublicKey) []byte {
	pubBytes := pubKey.SerializeUncompressed()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:])
	addressByte := hash.Sum(nil)
	return addressByte
}

func generateEIP1559Tx(rlpStr string) (*types.Transaction, error) {
	tx := new(Eip1559Transaction)
	if len(rlpStr) > 1 {
		if rlpStr[0:2] == "02" {
			rlpStr = rlpStr[2:]
		}
	}
	unsignedByte, err := hex.DecodeString(rlpStr)
	if err != nil {
		return nil, err
	}
	if err = rlp.DecodeBytes(unsignedByte, &tx); err != nil {
		return nil, err
	}
	return NewEip1559Transaction(tx.ChainId, tx.Nonce, tx.GasTipCap, tx.GasFeeCap, tx.Gas, tx.To, tx.Value, tx.Data), nil
}

func calMessageHash(data string) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	s := sha3.NewLegacyKeccak256()
	s.Write([]byte(msg))
	return s.Sum(nil)
}
