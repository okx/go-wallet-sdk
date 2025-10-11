package types

import (
	"encoding/json"
	"github.com/emresenyuva/go-wallet-sdk/coins/nervos/crypto"
)

const HashLength = 32

type ScriptHashType string
type DepType string

// Hash ckb hash, '0x' prefix hex string
// https://github.com/nervosnetwork/ckb-sdk-go/blob/master/types/common.go#L18
type Hash [32]byte

type OutPoint struct {
	TxHash Hash `json:"tx_hash"`
	Index  uint `json:"index"`
}

type CellInput struct {
	Since          uint64    `json:"since"`
	PreviousOutput *OutPoint `json:"previous_output"`
}

type Script struct {
	CodeHash Hash           `json:"code_hash"`
	HashType ScriptHashType `json:"hash_type"`
	Args     []byte         `json:"args"`
}

type CellOutput struct {
	Capacity uint64  `json:"capacity"`
	Lock     *Script `json:"lock"`
	Type     *Script `json:"type"`
}

type CellDep struct {
	OutPoint *OutPoint `json:"out_point"`
	DepType  DepType   `json:"dep_type"`
}

type WitnessArgs struct {
	Lock       []byte `json:"lock"`
	InputType  []byte `json:"input_type"`
	OutputType []byte `json:"output_type"`
}

type Transaction struct {
	Version     uint          `json:"version"`
	Hash        Hash          `json:"hash"`
	CellDeps    []*CellDep    `json:"cell_deps"`
	HeaderDeps  []Hash        `json:"header_deps"`
	Inputs      []*CellInput  `json:"inputs"`
	Outputs     []*CellOutput `json:"outputs"`
	OutputsData [][]byte      `json:"outputs_data"`
	Witnesses   [][]byte      `json:"witnesses"`
}

type outPoint struct {
	TxHash Hash `json:"tx_hash"`
	Index  Uint `json:"index"`
}

type cellDep struct {
	OutPoint outPoint `json:"out_point"`
	DepType  DepType  `json:"dep_type"`
}

type cellInput struct {
	Since          Uint64   `json:"since"`
	PreviousOutput outPoint `json:"previous_output"`
}

type script struct {
	CodeHash Hash           `json:"code_hash"`
	HashType ScriptHashType `json:"hash_type"`
	Args     Bytes          `json:"args"`
}

type cellOutput struct {
	Capacity Uint64  `json:"capacity"`
	Lock     *script `json:"lock"`
	Type     *script `json:"type"`
}

type inTransaction struct {
	Version     Uint         `json:"version"`
	CellDeps    []cellDep    `json:"cell_deps"`
	HeaderDeps  []Hash       `json:"header_deps"`
	Inputs      []cellInput  `json:"inputs"`
	Outputs     []cellOutput `json:"outputs"`
	OutputsData []Bytes      `json:"outputs_data"`
	Witnesses   []Bytes      `json:"witnesses"`
}

func fromTransaction(tx *Transaction) inTransaction {
	result := inTransaction{
		Version:     Uint(tx.Version),
		HeaderDeps:  tx.HeaderDeps,
		CellDeps:    fromCellDeps(tx.CellDeps),
		Inputs:      fromInputs(tx.Inputs),
		Outputs:     fromOutputs(tx.Outputs),
		OutputsData: fromBytesArray(tx.OutputsData),
		Witnesses:   fromBytesArray(tx.Witnesses),
	}
	return result
}

func MarshalTx(tx *Transaction) ([]byte, error) {
	return json.Marshal(fromTransaction(tx))
}

func fromInputs(inputs []*CellInput) []cellInput {
	result := make([]cellInput, len(inputs))
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		result[i] = cellInput{
			Since: Uint64(input.Since),
			PreviousOutput: outPoint{
				TxHash: input.PreviousOutput.TxHash,
				Index:  Uint(input.PreviousOutput.Index),
			},
		}
	}
	return result
}

func fromCellDeps(deps []*CellDep) []cellDep {
	result := make([]cellDep, len(deps))
	for i := 0; i < len(deps); i++ {
		dep := deps[i]
		result[i] = cellDep{
			OutPoint: outPoint{
				TxHash: dep.OutPoint.TxHash,
				Index:  Uint(dep.OutPoint.Index),
			},
			DepType: dep.DepType,
		}
	}
	return result
}

func fromOutputs(outputs []*CellOutput) []cellOutput {
	result := make([]cellOutput, len(outputs))
	for i := 0; i < len(outputs); i++ {
		output := outputs[i]
		result[i] = cellOutput{
			Capacity: Uint64(output.Capacity),
			Lock: &script{
				CodeHash: output.Lock.CodeHash,
				HashType: output.Lock.HashType,
				Args:     output.Lock.Args,
			},
		}
		if output.Type != nil {
			result[i].Type = &script{
				CodeHash: output.Type.CodeHash,
				HashType: output.Type.HashType,
				Args:     output.Type.Args,
			}
		}
	}
	return result
}

func fromBytesArray(bytes [][]byte) []Bytes {
	result := make([]Bytes, len(bytes))
	for i, data := range bytes {
		result[i] = data
	}
	return result
}

func (h Hash) Hex() string {
	return Encode(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func (h Hash) String() string {
	return h.Hex()
}

func (h Hash) MarshalText() ([]byte, error) {
	return Bytes(h[:]).MarshalText()
}

func (t *Transaction) ComputeHash() (Hash, error) {
	data, err := t.Serialize()
	if err != nil {
		return Hash{}, err
	}

	hash, err := crypto.Blake256(data)
	if err != nil {
		return Hash{}, err
	}

	return BytesToHash(hash), nil
}

// Serialize transaction
func (t *Transaction) Serialize() ([]byte, error) {
	v := SerializeUint(t.Version)

	// Ok, no way around this
	deps := make([]Serializer, len(t.CellDeps))
	for i, v := range t.CellDeps {
		deps[i] = v
	}
	cds, err := SerializeArray(deps)
	if err != nil {
		return nil, err
	}
	cdsBytes := SerializeFixVec(cds)

	hds := make([][]byte, len(t.HeaderDeps))
	for i := 0; i < len(t.HeaderDeps); i++ {
		hd, err := t.HeaderDeps[i].Serialize()
		if err != nil {
			return nil, err
		}

		hds[i] = hd
	}
	hdsBytes := SerializeFixVec(hds)

	ips := make([][]byte, len(t.Inputs))
	for i := 0; i < len(t.Inputs); i++ {
		ip, err := t.Inputs[i].Serialize()
		if err != nil {
			return nil, err
		}

		ips[i] = ip
	}
	ipsBytes := SerializeFixVec(ips)

	ops := make([][]byte, len(t.Outputs))
	for i := 0; i < len(t.Outputs); i++ {
		op, err := t.Outputs[i].Serialize()
		if err != nil {
			return nil, err
		}

		ops[i] = op
	}
	opsBytes := SerializeDynVec(ops)

	ods := make([][]byte, len(t.OutputsData))
	for i := 0; i < len(t.OutputsData); i++ {
		od := SerializeBytes(t.OutputsData[i])

		ods[i] = od
	}
	odsBytes := SerializeDynVec(ods)

	fields := [][]byte{v, cdsBytes, hdsBytes, ipsBytes, opsBytes, odsBytes}
	return SerializeTable(fields), nil
}
