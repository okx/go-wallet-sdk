/*
*

	Copyright Cosmos-SDK Authors
	Copyright 2016 All in Bits, Inc

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package types

// StdTx is the legacy transaction format for wrapping a Msg with Fee and Signatures.
// It only works with Amino, please prefer the new protobuf Tx in types/tx.
// NOTE: the first signature is the fee payer (Signatures must not be nil).
type StdTx struct {
	Msgs          []StdAny       `json:"msg" yaml:"msg"`
	Fee           StdFee         `json:"fee" yaml:"fee"`
	Signatures    []StdSignature `json:"signatures" yaml:"signatures"`
	Memo          string         `json:"memo" yaml:"memo"`
	TimeoutHeight uint64         `json:"timeout_height" yaml:"timeout_height"`
}

// StdFee includes the amount of coins paid in fees and the maximum
// gas to be used by the transaction. The ratio yields an effective "gasprice",
// which must be above some miminum to be accepted into the mempool.
type StdFee struct {
	Amount  Coins  `json:"amount" yaml:"amount"`
	Gas     string `json:"gas" yaml:"gas"`
	Payer   string `json:"payer,omitempty" yaml:"payer"`
	Granter string `json:"granter,omitempty" yaml:"granter"`
}

// StdSignature represents a sig
type StdSignature struct {
	Signature     string `json:"signature" yaml:"signature"`
	AccountNumber uint64 `json:"account_number" yaml:"account_number"`
	Sequence      uint64 `json:"sequence" yaml:"sequence"`
	PubKey        StdAny `json:"pub_key" yaml:"pub_key"`
}

type StdAny struct {
	T string      `json:"type" yaml:"type"`
	V interface{} `json:"value" yaml:"value"`
}

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
// Convert StdSignDoc to bytes using Amino/JSON, sign it, fill in StdTx, serialize StdTx with Amino/JSON, and then broadcast the transaction bytes
type StdSignDoc struct {
	AccountNumber string   `json:"account_number" yaml:"account_number"`
	Sequence      string   `json:"sequence" yaml:"sequence"`
	ChainID       string   `json:"chain_id" yaml:"chain_id"`
	Fee           StdFee   `json:"fee" yaml:"fee"`
	Msgs          []StdAny `json:"msgs" yaml:"msgs"`
	Memo          string   `json:"memo" yaml:"memo"`
	TimeoutHeight string   `json:"timeout_height,omitempty" yaml:"timeout_height"`
}
