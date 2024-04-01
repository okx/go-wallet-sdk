package doginals

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"math/big"
	"time"
)

const (
	PubKeyHashAddrID = 0x1e
	ScriptHashAddrID = 0x16
	PrivateKeyID     = 0x9e
)

var (
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// mainPowLimit is the highest proof of work value a Bitcoin block can
	// have for the main network.  It is the value 2^224 - 1.
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
)

// MainNetParams defines the network parameters for the main Bitcoin network.
var DogeMainNetParams = chaincfg.Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "8333",
	// Chain parameters
	PowLimit:                 mainPowLimit,
	PowLimitBits:             0x1d00ffff,
	BIP0034Height:            227931, // 000000000000024b89b42a942fe0d9fea3bb44ab7bd1b19115dd6a759c0808b8
	BIP0065Height:            388381, // 000000000000000004c2b624ed5d7756c508d90fd0da2c7c679febfa6c4735f0
	BIP0066Height:            363725, // 00000000000000000379eaa19dce8c9b722d46ae6a57c2f1a988119488b50931
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //
	// Mempool parameters
	RelayNonStdTxs: false,

	Bech32HRPSegwit: "doge",

	PubKeyHashAddrID:        PubKeyHashAddrID,
	ScriptHashAddrID:        ScriptHashAddrID,
	PrivateKeyID:            PrivateKeyID,
	WitnessPubKeyHashAddrID: 0x00,
	WitnessScriptHashAddrID: 0x00,

	HDPublicKeyID:  [4]byte{0x02, 0xfa, 0xca, 0xfd},
	HDPrivateKeyID: [4]byte{0x02, 0xfa, 0xc3, 0x98},

	HDCoinType: 3,
}
