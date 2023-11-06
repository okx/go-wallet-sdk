package core

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type HttpTransaction struct {
	// Base64 encoded content of the Cadence script.
	Script string `json:"script"`
	// An array containing arguments each encoded as Base64 passed in the [JSON-Cadence interchange format](https://docs.onflow.org/cadence/json-cadence-spec/).
	Arguments        []string `json:"arguments"`
	ReferenceBlockId string   `json:"reference_block_id"`
	// The limit on the amount of computation a transaction is allowed to preform.
	GasLimit           string                     `json:"gas_limit"`
	Payer              string                     `json:"payer"`
	ProposalKey        *HttpProposalKey           `json:"proposal_key"`
	Authorizers        []string                   `json:"authorizers"`
	PayloadSignatures  []HttpTransactionSignature `json:"payload_signatures"`
	EnvelopeSignatures []HttpTransactionSignature `json:"envelope_signatures"`
}

type HttpProposalKey struct {
	Address        string `json:"address"`
	KeyIndex       string `json:"key_index"`
	SequenceNumber string `json:"sequence_number"`
}

type HttpTransactionSignature struct {
	Address   string `json:"address"`
	KeyIndex  string `json:"key_index"`
	Signature string `json:"signature"`
}

type HttpTransactionSignatures []HttpTransactionSignature

func TransactionToHTTP(tx Transaction) ([]byte, error) {
	auths := make([]string, len(tx.Authorizers))
	for i, address := range tx.Authorizers {
		auths[i] = address.String()
	}

	args := make([]string, len(tx.Arguments))
	for i, argument := range tx.Arguments {
		args[i] = base64.StdEncoding.EncodeToString(argument)
	}

	payloadSignatures := make(HttpTransactionSignatures, len(tx.PayloadSignatures))
	for i, sig := range tx.PayloadSignatures {
		payloadSignatures[i] = HttpTransactionSignature{
			Address:   sig.Address.String(),
			KeyIndex:  fmt.Sprintf("%d", sig.KeyIndex),
			Signature: base64.StdEncoding.EncodeToString(sig.Signature),
		}
	}

	envelopeSignatures := make(HttpTransactionSignatures, len(tx.EnvelopeSignatures))
	for i, sig := range tx.EnvelopeSignatures {
		envelopeSignatures[i] = HttpTransactionSignature{
			Address:   sig.Address.String(),
			KeyIndex:  fmt.Sprintf("%d", sig.KeyIndex),
			Signature: base64.StdEncoding.EncodeToString(sig.Signature),
		}
	}

	return json.Marshal(HttpTransaction{
		Script:           base64.StdEncoding.EncodeToString(tx.Script),
		Arguments:        args,
		ReferenceBlockId: hex.EncodeToString(tx.ReferenceBlockID[:]),
		GasLimit:         fmt.Sprintf("%d", tx.GasLimit),
		Payer:            tx.Payer.String(),
		ProposalKey: &HttpProposalKey{
			Address:        tx.ProposalKey.Address.String(),
			KeyIndex:       fmt.Sprintf("%d", tx.ProposalKey.KeyIndex),
			SequenceNumber: fmt.Sprintf("%d", tx.ProposalKey.SequenceNumber),
		},
		Authorizers:        auths,
		PayloadSignatures:  payloadSignatures,
		EnvelopeSignatures: envelopeSignatures,
	})
}
