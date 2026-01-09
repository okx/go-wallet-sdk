package cardano

import (
	"fmt"

	"github.com/okx/go-wallet-sdk/coins/cardano/crypto"
	"golang.org/x/crypto/blake2b"
)

// TxBuilder is a transaction builder.
type TxBuilder struct {
	tx       *Tx
	protocol *ProtocolParams

	changeReceiver *Address
	txBuilt        bool
	change         uint64
	max            bool
}

// NewTxBuilder returns a new instance of TxBuilder.
func NewTxBuilder(protocol *ProtocolParams) *TxBuilder {
	return &TxBuilder{
		protocol: protocol,
		tx: &Tx{
			IsValid: true,
		},
	}
}

func (tb *TxBuilder) GetChange() uint64 {
	return tb.change
}

// AddInputs adds inputs to the transaction.
func (tb *TxBuilder) AddInputs(inputs ...*TxInput) {
	tb.tx.Body.Inputs = append(tb.tx.Body.Inputs, inputs...)
}

// AddOutputs adds outputs to the transaction.
func (tb *TxBuilder) AddOutputs(outputs ...*TxOutput) {
	tb.tx.Body.Outputs = append(tb.tx.Body.Outputs, outputs...)
}

// SetTtl sets the transaction's time to live.
func (tb *TxBuilder) SetTTL(ttl uint64) {
	tb.tx.Body.TTL = NewUint64(ttl)
}

// SetFee sets the transactions's fee.
func (tb *TxBuilder) SetFee(fee Coin) {
	tb.tx.Body.Fee = fee
}

// SetMax sets whether to calculate maximum ADA for the first output.
func (tb *TxBuilder) SetMax() {
	tb.max = true

	// For max mode, only keep the first output and set the ADA to 0
	if len(tb.tx.Body.Outputs) > 1 {
		tb.tx.Body.Outputs = tb.tx.Body.Outputs[:1]
	}
	if tb.tx.Body.Outputs[0].Amount.Coin != 0 {
		tb.tx.Body.Outputs[0].Amount.Coin = 0
	}
}

// AddAuxiliaryData adds auxiliary data to the transaction.
func (tb *TxBuilder) AddAuxiliaryData(data *AuxiliaryData) {
	tb.tx.AuxiliaryData = data
}

// AddCertificate adds a certificate to the transaction.
func (tb *TxBuilder) AddCertificate(cert Certificate) {
	tb.tx.Body.Certificates = append(tb.tx.Body.Certificates, cert)
}

// AddNativeScript adds a native script to the transaction.
func (tb *TxBuilder) AddNativeScript(script NativeScript) {
	tb.tx.WitnessSet.Scripts = append(tb.tx.WitnessSet.Scripts, script)
}

// Mint adds a new multiasset to mint.
func (tb *TxBuilder) Mint(asset *Mint) {
	tb.tx.Body.Mint = asset
}

// AddChangeIfNeeded instructs the builder to calculate the required fee for the
// transaction and to add an aditional output for the change if there is any.
func (tb *TxBuilder) AddChangeIfNeeded(changeAddr Address) {
	tb.changeReceiver = &changeAddr
}

func (tb *TxBuilder) AddSignature(pubKey crypto.PubKey, signature []byte) {
	witnessSet := make([]VKeyWitness, 0)
	witnessSet = append(witnessSet, VKeyWitness{
		VKey:      pubKey,
		Signature: signature,
	})
	tb.tx.WitnessSet.VKeyWitnessSet = witnessSet
}

func (tb *TxBuilder) AddEmptySignature() {
	pubKey := make([]byte, 32)
	signature := make([]byte, 64)
	tb.AddSignature(pubKey, signature)
}

func (tb *TxBuilder) calculateAmounts() (*Value, *Value) {
	input, output := NewValue(0), NewValue(tb.totalDeposits())
	for _, in := range tb.tx.Body.Inputs {
		input = input.Add(in.Amount)
	}
	for _, out := range tb.tx.Body.Outputs {
		output = output.Add(out.Amount)
	}
	if tb.tx.Body.Mint != nil {
		input = input.Add(NewValueWithAssets(0, tb.tx.Body.Mint.MultiAsset()))
	}
	return input, output
}

func (tb *TxBuilder) totalDeposits() Coin {
	certs := tb.tx.Body.Certificates
	var deposit Coin
	if len(certs) != 0 {
		for _, cert := range certs {
			if cert.Type == StakeRegistration {
				deposit += tb.protocol.KeyDeposit
			}
		}
	}
	return deposit
}

// MinFee computes the minimal fee required for the transaction.
// This assumes that the inputs-outputs are defined and signing keys are present.
func (tb *TxBuilder) MinFee() (Coin, error) {
	// Set a temporary realistic fee in order to serialize a valid transaction
	currentFee := tb.tx.Body.Fee
	tb.tx.Body.Fee = 200000
	tb.AddEmptySignature()
	minFee := tb.calculateMinFee()
	tb.tx.Body.Fee = currentFee
	return minFee, nil
}

func (tb *TxBuilder) Fee() (Coin, error) {
	return tb.tx.Body.Fee, nil
}

// MinCoinsForTxOut computes the minimal amount of coins required for a given transaction output.
func (tb *TxBuilder) MinCoinsForTxOut(txOut *TxOutput) Coin {
	bytes, err := cborEnc.Marshal(txOut)
	if err != nil {
		panic(err)
	}
	return Coin(uint64(len(bytes))+160) * tb.protocol.CoinsPerUTXOByte
}

// calculateMinFee computes the minimal fee required for the transaction.
func (tb *TxBuilder) calculateMinFee() Coin {
	txBytes := tb.tx.Bytes()
	txLength := uint64(len(txBytes))
	return tb.protocol.MinFeeA*Coin(txLength) + tb.protocol.MinFeeB
}

// Reset resets the builder to its initial state.
func (tb *TxBuilder) Reset() {
	tb.tx = &Tx{IsValid: true}
	tb.changeReceiver = nil
}

func (tb *TxBuilder) hasSufficientAda(output *TxOutput) bool {
	minAda := tb.MinCoinsForTxOut(output)
	return output.Amount.Coin >= minAda
}

// createChangeOutput creates a TxOutput for change with the given value
func (tb *TxBuilder) createChangeOutput(value *Value) *TxOutput {
	return &TxOutput{
		Address: *tb.changeReceiver,
		Amount:  value,
	}
}

// lastOutputIndex returns the index of the last output
func (tb *TxBuilder) lastOutputIndex() int {
	return len(tb.tx.Body.Outputs) - 1
}

// finalizeWithoutAdachg completes the transaction when no separate ADA change output exists.
// If machg exists, remaining ADA is added to it. Otherwise, remaining ADA is burned as fee.
func (tb *TxBuilder) finalizeWithoutAdachg(adaChange Coin, machgIndex int) (bool, error) {
	tb.AddEmptySignature()
	minFee := tb.calculateMinFee()

	if adaChange < minFee {
		tb.tx.Body.Fee = minFee
		return false, nil
	}

	if machgIndex != -1 {
		tb.tx.Body.Outputs[machgIndex].Amount.Coin += adaChange - minFee
		tb.tx.Body.Fee = minFee
	} else {
		tb.tx.Body.Fee = adaChange
	}
	return true, nil
}

// BuildUnsigned builds the transaction and calculates change outputs.
// Returns (valid, error) where valid indicates if the transaction has sufficient funds.
func (tb *TxBuilder) BuildUnsigned() (bool, error) {
	if err := tb.buildAuxiliaryData(); err != nil {
		return false, err
	}
	tb.txBuilt = true
	tb.tx.Body.Fee = 3e5

	inputAmount, outputAmount := tb.calculateAmounts()
	if res := inputAmount.Cmp(outputAmount); res < 0 || res == 2 {
		return false, fmt.Errorf("not enough ada")
	}

	changeAmount := inputAmount.Sub(outputAmount)
	adaChange := changeAmount.Coin
	machgIndex := -1

	// Multi-asset change: create output with tokens and minimum required ADA
	if !changeAmount.OnlyCoin() {
		minAda := calcMinAda(*tb.changeReceiver, changeAmount.MultiAsset)
		tb.AddOutputs(tb.createChangeOutput(&Value{
			Coin:       minAda,
			MultiAsset: changeAmount.MultiAsset,
		}))
		tb.change = uint64(minAda)
		machgIndex = tb.lastOutputIndex()

		if adaChange < minAda {
			tb.AddEmptySignature()
			tb.tx.Body.Fee = tb.calculateMinFee()
			return false, nil
		}
		adaChange -= minAda
	}

	// Max mode: set first output ADA to max available minus fee
	if tb.max {
		tb.tx.Body.Outputs[0].Amount.Coin = adaChange
		tb.AddEmptySignature()
		minFee := tb.calculateMinFee()

		if adaChange < minFee {
			tb.tx.Body.Fee = minFee
			return false, nil
		}

		tb.tx.Body.Outputs[0].Amount.Coin = adaChange - minFee
		if !tb.hasSufficientAda(tb.tx.Body.Outputs[0]) {
			tb.tx.Body.Fee = minFee
			return false, nil
		}

		tb.tx.Body.Fee = minFee
		return true, nil
	}

	// ADA change: try to create separate output if sufficient
	adachgOutput := tb.createChangeOutput(NewValue(adaChange))
	if !tb.hasSufficientAda(adachgOutput) {
		return tb.finalizeWithoutAdachg(adaChange, machgIndex)
	}

	// Add ADA change output and calculate fee
	tb.AddOutputs(adachgOutput)
	tb.AddEmptySignature()
	minFee := tb.calculateMinFee()

	// Verify ADA change is still sufficient after fee deduction
	if adaChange > minFee {
		adachgIndex := tb.lastOutputIndex()
		tb.tx.Body.Outputs[adachgIndex].Amount.Coin = adaChange - minFee

		if tb.hasSufficientAda(tb.tx.Body.Outputs[adachgIndex]) {
			tb.tx.Body.Fee = minFee
			return true, nil
		}
	}

	// ADA change insufficient after fee, remove it
	tb.tx.Body.Outputs = tb.tx.Body.Outputs[:tb.lastOutputIndex()]
	return tb.finalizeWithoutAdachg(adaChange, machgIndex)
}

func (tb *TxBuilder) Build(privateKeys ...crypto.PrvKey) (*Tx, error) {
	valid, err := tb.BuildUnsigned()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("not enough ada")
	}

	if err := tb.Sign(privateKeys...); err != nil {
		return nil, err
	}

	return tb.tx, nil
}

func (tb *TxBuilder) GetTx() *Tx {
	return tb.tx
}

func (tb *TxBuilder) GetTxHash() ([]byte, error) {
	return tb.tx.Hash()
}

func (tb *TxBuilder) GetChangeAmount() *Value {
	inputAmount, outputAmount := tb.calculateAmounts()
	return inputAmount.Sub(outputAmount)
}

func (tb *TxBuilder) Sign(privateKeys ...crypto.PrvKey) error {
	witnessSet := make([]VKeyWitness, 0)
	if !tb.txBuilt {
		return fmt.Errorf("transaction is not built")
	}
	txHash, err := tb.tx.Hash()
	if err != nil {
		return err
	}

	// Create witness set
	for _, pkey := range privateKeys {
		witnessSet = append(witnessSet, VKeyWitness{
			VKey:      pkey.PubKey(),
			Signature: pkey.Sign(txHash),
		})
	}
	tb.tx.WitnessSet.VKeyWitnessSet = witnessSet
	return nil
}

func (tb *TxBuilder) buildAuxiliaryData() error {
	if tb.tx.AuxiliaryData != nil {
		auxBytes, err := cborEnc.Marshal(tb.tx.AuxiliaryData)
		if err != nil {
			return err
		}
		auxHash := blake2b.Sum256(auxBytes)
		auxHash32 := Hash32(auxHash[:])
		tb.tx.Body.AuxiliaryDataHash = &auxHash32
	}
	return nil
}

func calcMinAda(address Address, assets *MultiAsset) Coin {
	value := NewValueWithAssets(Coin(3000000), assets)
	tb := NewTxBuilder(protocolParams)
	return tb.MinCoinsForTxOut(&TxOutput{
		Address: address,
		Amount:  value,
	})
}
