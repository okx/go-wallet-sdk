package bitcoin

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

const (
	zecNet = 0x6427e924
)

// GetBTCMainNetParams BTC
func GetBTCMainNetParams() *chaincfg.Params {
	return &chaincfg.MainNetParams
}

func GetBTCTestNetParams() *chaincfg.Params {
	return &chaincfg.TestNet3Params
}

// GetDGBMainNetParams DGB
func GetDGBMainNetParams() *chaincfg.Params {
	params := chaincfg.MainNetParams
	params.Net = 0xdab6c3fa

	// Address encoding magics
	params.PubKeyHashAddrID = 30 // base58 prefix: D
	params.ScriptHashAddrID = 63 // base58 prefix: 3
	params.Bech32HRPSegwit = "dgb"
	return &params
}

// GetQTUMMainNetParams QTUM
func GetQTUMMainNetParams() *chaincfg.Params {
	params := chaincfg.MainNetParams
	params.Net = 0xf1cfa6d3

	// Address encoding magics
	params.PubKeyHashAddrID = 58 // base58 prefix: Q
	params.ScriptHashAddrID = 50 // base58 prefix: P
	params.Bech32HRPSegwit = "qc"

	return &params
}

// GetRVNMainNetParams RVN
func GetRVNMainNetParams() *chaincfg.Params {
	params := chaincfg.MainNetParams
	params.Net = 0x4e564152

	// Address encoding magics
	params.PubKeyHashAddrID = 60  // base58 prefix: R
	params.ScriptHashAddrID = 122 // base58 prefix: r
	return &params
}

// GetBTGMainNetParams BTG
func GetBTGMainNetParams() *chaincfg.Params {
	mainnetparams := chaincfg.MainNetParams
	mainnetparams.Net = 0x446d47e1

	// Address encoding magics
	mainnetparams.PubKeyHashAddrID = 38 // base58 prefix: G
	mainnetparams.ScriptHashAddrID = 23 // base58 prefix: A

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	// see https://github.com/satoshilabs/slips/blob/master/slip-0173.md
	mainnetparams.Bech32HRPSegwit = "btg"

	return &mainnetparams
}

// GetBCHmainNetParams BCH
func GetBCHmainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = 0xe8f3e1e3

	// Address encoding magics
	mainNetParams.PubKeyHashAddrID = 0
	mainNetParams.ScriptHashAddrID = 5
	return &mainNetParams
}

// GetLTCMainNetParams LTC
func GetLTCMainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = 0xdbb6c0fb

	// Address encoding magics
	mainNetParams.PubKeyHashAddrID = 48
	mainNetParams.ScriptHashAddrID = 50
	mainNetParams.Bech32HRPSegwit = "ltc"
	return &mainNetParams
}

// GetDASHMainNetParams DASH
func GetDASHMainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = 0xbd6b0cbf

	// Address encoding magics
	mainNetParams.PubKeyHashAddrID = 76
	mainNetParams.ScriptHashAddrID = 16
	return &mainNetParams
}

// GetDOGEMainNetParams DOGE
func GetDOGEMainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = 0xc0c0c0c0

	// Address encoding magics
	mainNetParams.PubKeyHashAddrID = 30
	mainNetParams.ScriptHashAddrID = 22 // base58 prefix: 9
	return &mainNetParams
}

func NewZECAddr(pubBytes []byte) string {
	version := []byte{0x1C, 0xB8}
	return NewOldAddr(version, btcutil.Hash160(pubBytes))
}

func NewOldAddr(version []byte, data []byte) string {
	var buf []byte
	buf = append(buf, version[1:]...)
	buf = append(buf, data...)
	return base58.CheckEncode(buf, version[0])
}

// GetZECMainNetParams ZEC
func GetZECMainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = zecNet

	mainNetParams.PubKeyHashAddrID = 0x1C
	mainNetParams.ScriptHashAddrID = 0xBD

	return &mainNetParams
}
