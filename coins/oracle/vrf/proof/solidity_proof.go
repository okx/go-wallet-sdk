package proof

// Logic for providing the precomputed values required by the solidity verifier,
// in binary-blob format.

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okx/go-wallet-sdk/crypto/vrf/secp256k1"
	"math/big"

	"go.dedis.ch/kyber/v3"
)

// SolidityProof contains precalculations which VRF.sol needs to verify proofs
type SolidityProof struct {
	P                           *Proof         // The core proof
	UWitness                    common.Address // Address of P.C*P.PK+P.S*G
	CGammaWitness, SHashWitness kyber.Point    // P.C*P.Gamma, P.S*HashToCurve(P.Seed)
	ZInv                        *big.Int       // Inverse of Z coord from ProjectiveECAdd(CGammaWitness, SHashWitness)
}

// String returns the values in p, in hexadecimal format
func (p *SolidityProof) String() string {
	return fmt.Sprintf(
		"SolidityProof{P: %s, UWitness: %x, CGammaWitness: %s, SHashWitness: %s, ZInv: %x}",
		p.P, p.UWitness, p.CGammaWitness, p.SHashWitness, p.ZInv)
}

func point() kyber.Point {
	return Secp256k1Curve.Point()
}

// SolidityPrecalculations returns the precomputed values needed by the solidity
// verifier, or an error on failure.
func SolidityPrecalculations(p *Proof) (*SolidityProof, error) {
	var rv SolidityProof
	rv.P = p
	c := secp256k1.IntToScalar(p.C)
	s := secp256k1.IntToScalar(p.S)
	u := point().Add(point().Mul(c, p.PublicKey), point().Mul(s, Generator))
	var err error
	rv.UWitness = secp256k1.EthereumAddress(u)
	rv.CGammaWitness = point().Mul(c, p.Gamma)
	hash, err := HashToCurve(p.PublicKey, p.Seed, func(*big.Int) {})
	if err != nil {
		return nil, err
	}
	rv.SHashWitness = point().Mul(s, hash)
	_, _, z := ProjectiveECAdd(rv.CGammaWitness, rv.SHashWitness)
	rv.ZInv = z.ModInverse(z, FieldSize)
	return &rv, nil
}
