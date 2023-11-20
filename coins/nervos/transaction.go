package nervos

import (
	"encoding/binary"
	"errors"
	"github.com/okx/go-wallet-sdk/coins/nervos/crypto"
	"github.com/okx/go-wallet-sdk/coins/nervos/types"
)

func NewSecp256k1SingleSigTx(scripts *types.SystemScripts) *types.Transaction {
	return &types.Transaction{
		Version:    0,
		HeaderDeps: []types.Hash{},
		CellDeps: []*types.CellDep{
			{
				OutPoint: scripts.SecpSingleSigCell.OutPoint,
				DepType:  types.DepTypeDepGroup,
			},
		},
	}
}

func NewEmptyWitnessArg(LockScriptLength uint) *types.WitnessArgs {
	return &types.WitnessArgs{
		Lock:       make([]byte, LockScriptLength),
		InputType:  nil,
		OutputType: nil,
	}
}

func AddInputsForTransaction(transaction *types.Transaction, inputs []*types.CellInput,
	signatureLengthInBytes uint) error {
	if len(inputs) == 0 {
		return errors.New("input cells empty")
	}
	start := len(transaction.Inputs)
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		transaction.Inputs = append(transaction.Inputs, input)
		transaction.Witnesses = append(transaction.Witnesses, []byte{})
	}
	emptyWitnessArgs := NewEmptyWitnessArg(signatureLengthInBytes)
	emptyWitnessArgsBytes, err := emptyWitnessArgs.Serialize()
	if err != nil {
		return err
	}
	transaction.Witnesses[start] = emptyWitnessArgsBytes
	return nil
}

// SingleSignTransaction group is an array, which content is the index of input after grouping
func SingleSignTransaction(transaction *types.Transaction, group []int, witnessArgs *types.WitnessArgs, key types.Key) error {
	data, err := witnessArgs.Serialize()
	if err != nil {
		return err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, err := transaction.ComputeHash()
	if err != nil {
		return err
	}

	message := append(hash.Bytes(), length...)
	message = append(message, data...)
	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = transaction.Witnesses[group[i]]
			length := make([]byte, 8)
			binary.LittleEndian.PutUint64(length, uint64(len(data)))
			message = append(message, length...)
			message = append(message, data...)
		}
	}
	// hash witnesses which do not in any input group
	for _, witness := range transaction.Witnesses[len(transaction.Inputs):] {
		length := make([]byte, 8)
		binary.LittleEndian.PutUint64(length, uint64(len(witness)))
		message = append(message, length...)
		message = append(message, witness...)
	}

	message, err = crypto.Blake256(message)
	if err != nil {
		return err
	}

	signed, err := key.Sign(message)
	if err != nil {
		return err
	}

	wa := &types.WitnessArgs{
		Lock:       signed,
		InputType:  witnessArgs.InputType,
		OutputType: witnessArgs.OutputType,
	}

	wab, err := wa.Serialize()
	if err != nil {
		return err
	}
	transaction.Witnesses[group[0]] = wab

	return nil
}
