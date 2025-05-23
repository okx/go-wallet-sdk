package nervos

import (
	"errors"
	"github.com/okx/go-wallet-sdk/coins/nervos/crypto"
	"github.com/okx/go-wallet-sdk/coins/nervos/types"
)

var (
	Secp256k1EmptyWitnessArgPlaceholder = make([]byte, 85)
	Secp256k1SignaturePlaceholder       = make([]byte, 65)
)

type TransactionBuilder struct {
	tx                 *types.Transaction
	chain              string
	privateKeyGroupMap map[string][]int
}

func NewTxBuild() *TransactionBuilder {
	scripts := types.NewSystemScripts("ckb")
	tx := NewSecp256k1SingleSigTx(scripts)
	return NewTxBuildWithTx(tx, "ckb")
}

func NewTestnetTxBuild() *TransactionBuilder {
	scripts := types.NewSystemScripts("ckb_testnet")
	tx := NewSecp256k1SingleSigTx(scripts)
	return NewTxBuildWithTx(tx, "ckb_testnet")
}

func NewTxBuildWithTx(tx *types.Transaction, chain string) *TransactionBuilder {
	builder := &TransactionBuilder{tx, chain, make(map[string][]int)}
	return builder
}

// AddInput adds an input to the transaction.
// The field since prevents a transaction being mined before a specific time. It already has its own RFC.
func (t *TransactionBuilder) AddInput(prevhash string, index uint, since uint64) error {
	if prevhash == "" {
		return errors.New("prevhash is empty")
	}
	input := &types.CellInput{
		Since: since,
		PreviousOutput: &types.OutPoint{
			TxHash: types.HexToHash(prevhash),
			Index:  index,
		},
	}
	inputs := []*types.CellInput{input}
	err := AddInputsForTransaction(t.tx, inputs, uint(len(Secp256k1SignaturePlaceholder)))
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionBuilder) AddInputWithPrivateKey(prevhash string, index uint, since uint64, privateKey string) error {
	inputIndex := len(t.tx.Inputs)
	t.privateKeyGroupMap[privateKey] = append(t.privateKeyGroupMap[privateKey], inputIndex)
	return t.AddInput(prevhash, index, since)
}

func (t *TransactionBuilder) AddOutput(address string, amount uint64) error {
	script, err := Parse(address)
	if err != nil {
		return err
	}
	t.AddOutputScript(amount, script.Script)
	return nil
}

func (t *TransactionBuilder) AddOutputScript(amount uint64, script *types.Script) {
	t.tx.Outputs = append(t.tx.Outputs, &types.CellOutput{
		Capacity: amount,
		Lock:     script,
	})
	t.tx.OutputsData = append(t.tx.OutputsData, []byte{})
}

func (t *TransactionBuilder) AddWitness(signature []byte) {
	t.tx.Witnesses = append(t.tx.Witnesses, signature)
}

func (t *TransactionBuilder) AddWitnessArgs(lock, inputType, outputType []byte) {
	t.tx.Witnesses = append(t.tx.Witnesses, append(lock, append(inputType, outputType...)...))
}

func (t *TransactionBuilder) Build() (*types.Transaction, error) {
	for _, output := range t.tx.Outputs {
		if output.Capacity < 60*OneCKBShannon {
			return nil, errors.New("output capacity must be greater than 60 ckb")
		}
	}
	return t.tx, nil
}

func (t *TransactionBuilder) Sign(keys ...string) error {
	if len(keys) == 1 {
		return t.SignByPrivateKey(keys[0])
	}
	emptyWitnessArgs := NewEmptyWitnessArg(uint(len(Secp256k1SignaturePlaceholder)))
	for keyText, group := range t.privateKeyGroupMap {
		key, err := crypto.HexToKey(keyText)
		if err != nil {
			return err
		}
		err = SingleSignTransaction(t.tx, group, emptyWitnessArgs, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TransactionBuilder) AddPrivateKeyGroups(keyText string, group []int) {
	t.privateKeyGroupMap[keyText] = group
}

func (t *TransactionBuilder) SignByPrivateKey(keyText string) error {
	emptyWitnessArgs := NewEmptyWitnessArg(uint(len(Secp256k1SignaturePlaceholder)))
	key, err := crypto.HexToKey(keyText)
	if err != nil {
		return err
	}
	var groups []int
	for i := 0; i < len(t.tx.Inputs); i++ {
		groups = append(groups, i)
	}
	if err := SingleSignTransaction(t.tx, groups, emptyWitnessArgs, key); err != nil {
		return err
	}
	return nil
}

func (t *TransactionBuilder) DumpTx() string {
	text, _ := types.MarshalTx(t.tx)
	return string(text)
}

// Deprecated: use AddOutput instead
func (t *TransactionBuilder) AddChangeOutput(address string) error {
	return nil
}
