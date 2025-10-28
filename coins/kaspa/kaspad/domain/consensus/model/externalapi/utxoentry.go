package externalapi

// UTXOEntry houses details about an individual transaction output in a utxo
// set such as whether or not it was contained in a coinbase tx, the daa
// score of the block that accepts the tx, its public key script, and how
// much it pays.
type UTXOEntry interface {
	Amount() uint64                    // Utxo amount in Sompis
	ScriptPublicKey() *ScriptPublicKey // The public key script for the output.
	BlockDAAScore() uint64             // Daa score of the block accepting the tx.
	IsCoinbase() bool
	Equal(other UTXOEntry) bool
}

// OutpointAndUTXOEntryPair is an outpoint along with its
// respective UTXO entry
type OutpointAndUTXOEntryPair struct {
	Outpoint  *DomainOutpoint
	UTXOEntry UTXOEntry
}
