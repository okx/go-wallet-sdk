package proof

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okx/go-wallet-sdk/crypto/vrf/utils"
	"math/big"
)

// Seed represents a VRF seed as a serialized uint256
type Seed [32]byte

// BigToSeed returns seed x represented as a Seed, or an error if x is too big
func BigToSeed(x *big.Int) (Seed, error) {
	seed, err := utils.Uint256ToBytes(x)
	if err != nil {
		return Seed{}, err
	}
	return Seed(common.BytesToHash(seed)), nil
}

// Big returns the uint256 seed represented by s
func (s *Seed) Big() *big.Int {
	return common.Hash(*s).Big()
}

// PreSeedData contains the data the VRF provider needs to compute the final VRF
// output and marshal the proof for transmission to the VRFCoordinator contract.
type PreSeedData struct {
	PreSeed          Seed        // Seed to be mixed with hash of containing block
	BlockHash        common.Hash // Hash of block containing VRF request
	BlockNum         uint64      // Cardinal number of block containing VRF request
	SubId            uint64
	CallbackGasLimit uint32
	NumWords         uint32
	Sender           common.Address
}

// FinalSeed is the seed which is actually passed to the VRF proof generator,
// given the pre-seed and the hash of the block in which the VRFCoordinator
// emitted the log for the request this is responding to.
func FinalSeed(s PreSeedData) (finalSeed *big.Int) {
	seedHashMsg := append(s.PreSeed[:], s.BlockHash.Bytes()...)
	return utils.MustHash(string(seedHashMsg)).Big()
}
