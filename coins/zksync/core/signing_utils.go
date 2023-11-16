/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
package core

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	AmountExponentBitWidth int64 = 5
	AmountMantissaBitWidth int64 = 35
	FeeExponentBitWidth    int64 = 5
	FeeMantissaBitWidth    int64 = 11
)

func Uint16ToBytes(v uint16) []byte {
	res := make([]byte, 2)
	binary.BigEndian.PutUint16(res, v)
	return res
}

func Uint32ToBytes(v uint32) []byte {
	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, v)
	return res
}

func Uint64ToBytes(v uint64) []byte {
	res := make([]byte, 8)
	binary.BigEndian.PutUint64(res, v)
	return res
}

func BigIntToBytesBE(v *big.Int, numBytes int) []byte {
	val := v.Bytes()
	res := make([]byte, numBytes-len(val)) // left padded with 0 bytes to target length
	return append(res, val...)
}

func pkhToBytes(pkh string) ([]byte, error) {
	if pkh[:5] != "sync:" {
		return nil, errors.New("pubKeyHash must start with 'sync:'")
	}
	res, err := hex.DecodeString(pkh[5:])
	if err != nil {
		return nil, err
	}
	if len(res) != 20 {
		return nil, errors.New("pkh must be 20 bytes long")
	}
	return res, nil
}

func PackFee(fee *big.Int) ([]byte, error) {
	return packFee(fee)
}

func PackAmount(fee *big.Int) ([]byte, error) {
	return packAmount(fee)
}

func packFee(fee *big.Int) ([]byte, error) {
	packedFee, err := integerToDecimalByteArray(fee, FeeExponentBitWidth, FeeMantissaBitWidth, 10)
	if err != nil {
		return nil, err
	}
	// check that unpacked fee still has same value
	if unpackedFee, err := decimalByteArrayToInteger(packedFee, FeeExponentBitWidth, FeeMantissaBitWidth, 10); err != nil {
		return nil, err
	} else if unpackedFee.Cmp(fee) != 0 {
		return nil, errors.New("fee amount is not packable")
	}
	return packedFee, nil
}

func packAmount(amount *big.Int) ([]byte, error) {
	packedAmount, err := integerToDecimalByteArray(amount, AmountExponentBitWidth, AmountMantissaBitWidth, 10)
	if err != nil {
		return nil, err
	}
	// check that unpacked amount still has same value
	if unpackedFee, err := decimalByteArrayToInteger(packedAmount, AmountExponentBitWidth, AmountMantissaBitWidth, 10); err != nil {
		return nil, err
	} else if unpackedFee.Cmp(amount) != 0 {
		return nil, errors.New("amount amount is not packable")
	}
	return packedAmount, nil
}

func integerToDecimalByteArray(value *big.Int, expBits, mantissaBits, expBase int64) ([]byte, error) {
	bigExpBase := big.NewInt(expBase)
	// maxExponent = expBase ^ ((2 ^ expBits) - 1)
	maxExpPow := big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(expBits), nil), big.NewInt(1))
	maxExponent := big.NewInt(0).Exp(bigExpBase, maxExpPow, nil)
	// maxMantissa = (2 ^ mantissaBits) - 1
	maxMantissa := big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(mantissaBits), nil), big.NewInt(1))
	// check for max possible value
	if value.Cmp(big.NewInt(0).Mul(maxMantissa, maxExponent)) > 0 {
		return nil, errors.New("integer is too big")
	}
	exponent := uint64(0)
	mantissa := big.NewInt(0).Set(value)
	for mantissa.Cmp(maxMantissa) > 0 {
		mantissa.Div(mantissa, bigExpBase)
		exponent++
	}

	exponentData := uint64ToBitsLE(exponent, uint(expBits))
	mantissaData := uint64ToBitsLE(mantissa.Uint64(), uint(mantissaBits))
	combined := exponentData.Clone().Append(mantissaData)
	reversed := combined.Reverse()
	bytes, err := reversed.ToBytesBE()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func decimalByteArrayToInteger(value []byte, expBits, mantissaBits, expBase int64) (*big.Int, error) {
	if int64(len(value)*8) != expBits+mantissaBits {
		return nil, errors.New("decimal unpacking, incorrect input length")
	}
	bits := NewBits(uint(expBits + mantissaBits))
	bits.FromBytesBE(value).Reverse()
	exponent := big.NewInt(0)
	expPow2 := big.NewInt(1)
	for i := uint(0); i < uint(expBits); i++ {
		if bits.GetBit(i) {
			exponent.Add(exponent, expPow2)
		}
		expPow2.Mul(expPow2, big.NewInt(2))
	}
	exponent.Exp(big.NewInt(expBase), exponent, nil)

	mantissa := big.NewInt(0)
	mantissaPow2 := big.NewInt(1)
	for i := uint(expBits); i < uint(expBits+mantissaBits); i++ {
		if bits.GetBit(i) {
			mantissa.Add(mantissa, mantissaPow2)
		}
		mantissaPow2.Mul(mantissaPow2, big.NewInt(2))
	}
	return exponent.Mul(exponent, mantissa), nil
}

func uint64ToBitsLE(v uint64, size uint) *Bits {
	res := NewBits(size)
	for i := uint(0); i < size; i++ {
		res.SetBit(i, v&1 == 1)
		v /= 2
	}
	return res
}

func getChangePubKeyData(txData *ChangePubKey) ([]byte, error) {
	buf := bytes.Buffer{}
	pkhBytes, err := pkhToBytes(txData.NewPkHash)
	if err != nil {
		return nil, err
	}
	buf.Write(pkhBytes)
	buf.Write(Uint32ToBytes(txData.Nonce))
	buf.Write(Uint32ToBytes(txData.AccountId))
	buf.Write(txData.EthAuthData.getBytes())
	return buf.Bytes(), nil
}

func getTransferMessagePart(to string, amount, fee *big.Int, token *Token) (string, error) {
	var res string
	if big.NewInt(0).Cmp(amount) != 0 {
		res = fmt.Sprintf("Transfer %s %s to: %s", token.ToDecimalString(amount), token.Symbol, strings.ToLower(to))
	}
	if fee.Cmp(big.NewInt(0)) > 0 {
		if len(res) > 0 {
			res += "\n"
		}
		res += fmt.Sprintf("Fee: %s %s", token.ToDecimalString(fee), token.Symbol)
	}
	return res, nil
}

func getWithdrawMessagePart(to string, amount, fee *big.Int, token *Token) (string, error) {
	var res string
	if big.NewInt(0).Cmp(amount) != 0 {
		res = fmt.Sprintf("Withdraw %s %s to: %s", token.ToDecimalString(amount), token.Symbol, strings.ToLower(to))
	}
	if fee.Cmp(big.NewInt(0)) > 0 {
		if len(res) > 0 {
			res += "\n"
		}
		res += fmt.Sprintf("Fee: %s %s", token.ToDecimalString(fee), token.Symbol)
	}
	return res, nil
}

func getForcedExitMessagePart(to string, fee *big.Int, token *Token) (string, error) {
	var res string
	res = fmt.Sprintf("ForcedExit %s to: %s", token.Symbol, strings.ToLower(to))
	if fee.Cmp(big.NewInt(0)) > 0 {
		res += fmt.Sprintf("\nFee: %s %s", token.ToDecimalString(fee), token.Symbol)
	}
	return res, nil
}

func getMintNFTMessagePart(contentHash [32]byte, to string, fee *big.Int, token *Token) (string, error) {
	var res string
	contentHashStr := HEX_PREFIX + hex.EncodeToString(contentHash[:])
	res = fmt.Sprintf("MintNFT %s for: %s", contentHashStr, strings.ToLower(to))
	if fee.Cmp(big.NewInt(0)) > 0 {
		res += fmt.Sprintf("\nFee: %s %s", token.ToDecimalString(fee), token.Symbol)
	}
	return res, nil
}

func getWithdrawNFTMessagePart(to string, tokenId uint32, fee *big.Int, token *Token) (string, error) {
	var res string
	res = fmt.Sprintf("WithdrawNFT %d to: %s", tokenId, strings.ToLower(to))
	if fee.Cmp(big.NewInt(0)) > 0 {
		res += fmt.Sprintf("\nFee: %s %s", token.ToDecimalString(fee), token.Symbol)
	}
	return res, nil
}

func getOrderMessagePart(recipient string, amount *big.Int, sell, buy *Token, ratio []*big.Int) (string, error) {
	if len(ratio) != 2 {
		return "", errors.New("invalid ratio")
	}
	var res string
	if amount.Cmp(big.NewInt(0)) == 0 {
		res = fmt.Sprintf("Limit order for %s -> %s", sell.Symbol, buy.Symbol)
	} else {
		res = fmt.Sprintf("Order for %s %s -> %s", sell.ToDecimalString(amount), sell.Symbol, buy.Symbol)
	}
	res += fmt.Sprintf("\nRatio: %s:%s\nAddress: %s", ratio[0].String(), ratio[1].String(), strings.ToLower(recipient))
	return res, nil
}

func getSwapMessagePart(token *Token, fee *big.Int) string {
	return fmt.Sprintf("Swap fee: %s %s", token.ToDecimalString(fee), token.Symbol)
}

func getNonceMessagePart(nonce uint32) string {
	return fmt.Sprintf("Nonce: %d", nonce)
}
