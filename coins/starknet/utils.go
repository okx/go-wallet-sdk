package starknet

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// given x will find corresponding public key coordinate on curve
func (sc StarkCurve) XToPubKey(x string) (*big.Int, *big.Int) {
	xin, _ := HexToBN(x)

	yout := sc.GetYCoordinate(xin)

	return xin, yout
}

// convert utf8 string to big int
func UTF8StrToBig(str string) *big.Int {
	hexStr := hex.EncodeToString([]byte(str))
	b, _ := new(big.Int).SetString(hexStr, 16)

	return b
}

// convert decimal string to big int
func StrToBig(str string) *big.Int {
	b, _ := new(big.Int).SetString(str, 10)

	return b
}

func HexToShortStr(hexStr string) string {
	numStr := strings.Replace(hexStr, "0x", "", -1)
	hb, _ := new(big.Int).SetString(numStr, 16)

	return string(hb.Bytes())
}

// trim "0x" prefix(if exists) and converts hexidecimal string to big int
func HexToBN(hexString string) (*big.Int, error) {
	numStr := strings.Replace(hexString, "0x", "", -1)

	n, ok := new(big.Int).SetString(numStr, 16)
	if !ok {
		return nil, fmt.Errorf("please input a rigth hex string")
	}
	return n, nil
}

// trim "0x" prefix(if exists) and converts hexidecimal string to byte slice
func HexToBytes(hexString string) ([]byte, error) {
	numStr := strings.Replace(hexString, "0x", "", -1)
	if (len(numStr) % 2) != 0 {
		numStr = fmt.Sprintf("%s%s", "0", numStr)
	}

	return hex.DecodeString(numStr)
}

func BytesToBig(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(bytes)
}

// convert big int to hexidecimal string
func BigToHex(in *big.Int) string {
	return fmt.Sprintf("0x%x", in)
}

func BigToHexWithPadding(in *big.Int) string {
	return fmt.Sprintf("0x%064x", in)
}

// convert hexidecimal string to big int
func HexToBig(in string) *big.Int {
	if strings.HasPrefix(in, "0x") {
		in = in[2:]
	}
	o := new(big.Int)
	o.SetString(in, 16)
	return o
}

/**
https://github.com/NethermindEth/starknet.go/blob/main/LICENSE

MIT License

Copyright (c) 2021 Don't Panic DAO

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
/*
	Although the library adheres to the 'elliptic/curve' interface.
	All testing has been done against library function explicity.
	It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).

*/

// obtain random primary key on stark curve
// NOTE: to be used for testing purposes
func (sc StarkCurve) GetRandomPrivateKey() (priv *big.Int, err error) {
	max := new(big.Int).Sub(sc.Max, big.NewInt(1))

	priv, err = rand.Int(rand.Reader, max)
	if err != nil {
		return priv, err
	}

	x, y, err := sc.PrivateToPoint(priv)
	if err != nil {
		return priv, err
	}

	if !sc.IsOnCurve(x, y) {
		return priv, fmt.Errorf("key gen is not on stark cruve")
	}

	return priv, nil
}

// obtain public key coordinates from stark curve given the private key
func (sc StarkCurve) PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	if privKey.Cmp(big.NewInt(0)) != 1 || privKey.Cmp(sc.N) != -1 {
		return x, y, fmt.Errorf("private key not in curve range")
	}
	x, y = sc.EcMult(privKey, sc.EcGenX, sc.EcGenY)
	return x, y, nil
}

func (sc StarkCurve) PrivateToPublic(privKey *big.Int) (pubKey *big.Int, err error) {
	pubKey, _, err = sc.PrivateToPoint(privKey)
	return pubKey, err
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rolen int) []byte {
	out := v.Bytes()

	// pad with zeros if it's too short
	if len(out) < rolen {
		out2 := make([]byte, rolen)
		copy(out2[rolen-len(out):], out)
		return out2
	}

	// drop most significant bytes if it's too long
	if len(out) > rolen {
		out2 := make([]byte, rolen)
		copy(out2, out[len(out)-rolen:])
		return out2
	}

	return out
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(in, q *big.Int, qlen, rolen int) []byte {
	z1 := bits2int(in, qlen)
	z2 := new(big.Int).Sub(z1, q)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.2
func bits2int(in *big.Int, qlen int) *big.Int {
	blen := len(in.Bytes()) * 8
	if blen > qlen {
		return new(big.Int).Rsh(in, uint(blen-qlen))
	}
	return in
}

// mac returns an HMAC of the given key and message.
func mac(alg func() hash.Hash, k, m, buf []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(buf[:0])
}

func GetSelectorFromName(funcName string) *big.Int {
	kec := Keccak256([]byte(funcName))

	maskedKec := MaskBits(250, 8, kec)

	return new(big.Int).SetBytes(maskedKec)
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

// NewKeccakState creates a new KeccakState
// (ref: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go)
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

func MaskBits(mask, wordSize int, slice []byte) (ret []byte) {
	excess := len(slice)*wordSize - mask
	for _, by := range slice {
		if excess > 0 {
			if excess > wordSize {
				excess = excess - wordSize
				continue
			}
			by <<= excess
			by >>= excess
			excess = 0
		}
		ret = append(ret, by)
	}
	return ret
}

func ComputeFact(programHash *big.Int, programOutputs []*big.Int) *big.Int {
	var progOutBuf []byte
	for _, programOutput := range programOutputs {
		inBuf := FmtKecBytes(programOutput, 32)
		progOutBuf = append(progOutBuf[:], inBuf...)
	}

	kecBuf := FmtKecBytes(programHash, 32)
	kecBuf = append(kecBuf[:], Keccak256(progOutBuf)...)

	return new(big.Int).SetBytes(Keccak256(kecBuf))
}

func SplitFactStr(fact string) (fact_low, fact_high string) {
	factBN, _ := HexToBN(fact)
	factBytes := factBN.Bytes()
	low := BytesToBig(factBytes[16:])
	high := BytesToBig(factBytes[:16])
	return BigToHex(low), BigToHex(high)
}

func FmtKecBytes(in *big.Int, rolen int) (buf []byte) {
	buf = append(buf, in.Bytes()...)

	// pad with zeros if too short
	if len(buf) < rolen {
		padded := make([]byte, rolen)
		copy(padded[rolen-len(buf):], buf)

		return padded
	}

	return buf
}

func jsToBN(str string) *big.Int {
	bn, _ := HexToBN(str)
	if strings.Contains(str, "0x") {
		return bn
	} else {
		return StrToBig(str)
	}
}

func FmtExecuteCalldata(txs []Transaction) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(txs)))}

	for _, tx := range txs {
		callArray = append(callArray, tx.ContractAddress, tx.EntryPointSelector)
		if len(tx.Calldata) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))
		} else {
			callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.Calldata))))
			calldataArray = append(calldataArray, tx.Calldata...)
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	return callArray
}

func OldFmtExecuteCalldata(nonce *big.Int, txs []Transaction) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(txs)))}

	for _, tx := range txs {
		callArray = append(callArray, tx.ContractAddress, tx.EntryPointSelector)
		if len(tx.Calldata) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))
		} else {
			callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.Calldata))))
			calldataArray = append(calldataArray, tx.Calldata...)
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	callArray = append(callArray, nonce)
	return callArray
}

func FmtExecuteCalldataStrings(txs []Transaction) (calldataStrings []string) {
	callArray := FmtExecuteCalldata(txs)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, data.String())
	}
	return calldataStrings
}
