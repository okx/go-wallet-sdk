/*
*
The MIT License (MIT)

Copyright (c) 2018 SmartContract ChainLink, Ltd.
*/
package proof

import (
	"fmt"
	"github.com/emresenyuva/go-wallet-sdk/crypto/vrf/secp256k1"
	"github.com/emresenyuva/go-wallet-sdk/crypto/vrf/utils"
	"math/big"

	"go.dedis.ch/kyber/v3"
)

// Proof represents a proof that Gamma was constructed from the Seed
// according to the process mandated by the PublicKey.
//
// N.B.: The kyber.Point fields must contain secp256k1.secp256k1Point values, C,
// S and Seed must be secp256k1Point, and Output must be at
// most 256 bits. See Proof.WellFormed.
type Proof struct {
	PublicKey kyber.Point // secp256k1 public key of private key used in proof
	Gamma     kyber.Point
	C         *big.Int
	S         *big.Int
	Seed      *big.Int // Seed input to verifiable random function
	Output    *big.Int // verifiable random function output;, uniform uint256 sample
}

func (p *Proof) String() string {
	return fmt.Sprintf(
		"vrf.Proof{PublicKey: %s, Gamma: %s, C: %x, S: %x, Seed: %x, Output: %x}",
		p.PublicKey, p.Gamma, p.C, p.S, p.Seed, p.Output)
}

// WellFormed is true iff p's attributes satisfy basic domain checks
func (p *Proof) WellFormed() bool {
	return (secp256k1.ValidPublicKey(p.PublicKey) &&
		secp256k1.ValidPublicKey(p.Gamma) && secp256k1.RepresentsScalar(p.C) &&
		secp256k1.RepresentsScalar(p.S) && p.Output.BitLen() <= 256)
}

// VerifyProof is true iff gamma was generated in the mandated way from the
// given publicKey and seed, and no error was encountered
func (p *Proof) VerifyVRFProof() (bool, error) {
	if !p.WellFormed() {
		return false, fmt.Errorf("badly-formatted proof")
	}
	h, err := HashToCurve(p.PublicKey, p.Seed, func(*big.Int) {})
	if err != nil {
		return false, err
	}
	err = checkCGammaNotEqualToSHash(p.C, p.Gamma, p.S, h)
	if err != nil {
		return false, fmt.Errorf("c*γ = s*hash (disallowed in solidity verifier)")
	}
	// publicKey = secretKey*Generator. See GenerateProof for u, v, m, s
	// c*secretKey*Generator + (m - c*secretKey)*Generator = m*Generator = u
	uPrime := linearCombination(p.C, p.PublicKey, p.S, Generator)
	// c*secretKey*h + (m - c*secretKey)*h = m*h = v
	vPrime := linearCombination(p.C, p.Gamma, p.S, h)
	uWitness := secp256k1.EthereumAddress(uPrime)
	cPrime, err := ScalarFromCurvePoints(h, p.PublicKey, p.Gamma, uWitness, vPrime)
	if err != nil {
		return false, fmt.Errorf("vrf.verifyProof#ScalarFromCurvePoints, %v", err)
	}
	output := utils.MustHash(string(append(
		RandomOutputHashPrefix, secp256k1.LongMarshal(p.Gamma)...)))
	return utils.Equal(p.C, cPrime) && utils.Equal(p.Output, output.Big()), nil
}
