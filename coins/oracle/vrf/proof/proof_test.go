/*
*
The MIT License (MIT)

Copyright (c) 2018 SmartContract ChainLink, Ltd.
*/
package proof

import (
	"math/big"
	"testing"
)

func TestVRF_VerifyProof(t *testing.T) {
	sk, err := NewV2()
	if err != nil {
		t.Error(err)
	}
	seed, nonce := big.NewInt(2), big.NewInt(3)
	p, err := sk.GenerateProofWithNonce(seed, nonce)
	if err != nil {
		t.Errorf("could not generate proof, %v", err)
	}
	valid, err := p.VerifyVRFProof()
	if err != nil {
		t.Errorf("could not validate proof, %v", err)
	}
	if !valid {
		t.Error("invalid proof was found valid")
	} else {
		t.Logf(
			"vrf.Proof{PublicKey: %s, Gamma: %s, C: %x, S: %x, Seed: %x, Output: %x}",
			p.PublicKey, p.Gamma, p.C, p.S, p.Seed, p.Output)
	}
}
