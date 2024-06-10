package bitcoin

// ErrorCode identifies a kind of error.
type ErrorCode int

// These constants are used to identify a specific RuleError.
const (
	// ErrDuplicateBlock indicates a block with the same hash already
	// exists.
	ErrDuplicateBlock ErrorCode = iota

	// ErrBlockTooBig indicates the serialized block size exceeds the
	// maximum allowed size.
	ErrBlockTooBig

	// ErrBlockWeightTooHigh indicates that the block's computed weight
	// metric exceeds the maximum allowed value.
	ErrBlockWeightTooHigh

	// ErrBlockVersionTooOld indicates the block version is too old and is
	// no longer accepted since the majority of the network has upgraded
	// to a newer version.
	ErrBlockVersionTooOld

	// ErrInvalidTime indicates the time in the passed block has a precision
	// that is more than one second.  The chain consensus rules require
	// timestamps to have a maximum precision of one second.
	ErrInvalidTime

	// ErrTimeTooOld indicates the time is either before the median time of
	// the last several blocks per the chain consensus rules or prior to the
	// most recent checkpoint.
	ErrTimeTooOld

	// ErrTimeTooNew indicates the time is too far in the future as compared
	// the current time.
	ErrTimeTooNew

	// ErrDifficultyTooLow indicates the difficulty for the block is lower
	// than the difficulty required by the most recent checkpoint.
	ErrDifficultyTooLow

	// ErrUnexpectedDifficulty indicates specified bits do not align with
	// the expected value either because it doesn't match the calculated
	// valued based on difficulty regarted rules or it is out of the valid
	// range.
	ErrUnexpectedDifficulty

	// ErrHighHash indicates the block does not hash to a value which is
	// lower than the required target difficultly.
	ErrHighHash

	// ErrBadMerkleRoot indicates the calculated merkle root does not match
	// the expected value.
	ErrBadMerkleRoot

	// ErrBadCheckpoint indicates a block that is expected to be at a
	// checkpoint height does not match the expected one.
	ErrBadCheckpoint

	// ErrForkTooOld indicates a block is attempting to fork the block chain
	// before the most recent checkpoint.
	ErrForkTooOld

	// ErrCheckpointTimeTooOld indicates a block has a timestamp before the
	// most recent checkpoint.
	ErrCheckpointTimeTooOld

	// ErrNoTransactions indicates the block does not have a least one
	// transaction.  A valid block must have at least the coinbase
	// transaction.
	ErrNoTransactions

	// ErrNoTxInputs indicates a transaction does not have any inputs.  A
	// valid transaction must have at least one input.
	ErrNoTxInputs

	// ErrNoTxOutputs indicates a transaction does not have any outputs.  A
	// valid transaction must have at least one output.
	ErrNoTxOutputs

	// ErrTxTooBig indicates a transaction exceeds the maximum allowed size
	// when serialized.
	ErrTxTooBig

	// ErrBadTxOutValue indicates an output value for a transaction is
	// invalid in some way such as being out of range.
	ErrBadTxOutValue

	// ErrDuplicateTxInputs indicates a transaction references the same
	// input more than once.
	ErrDuplicateTxInputs

	// ErrBadTxInput indicates a transaction input is invalid in some way
	// such as referencing a previous transaction outpoint which is out of
	// range or not referencing one at all.
	ErrBadTxInput

	// ErrMissingTxOut indicates a transaction output referenced by an input
	// either does not exist or has already been spent.
	ErrMissingTxOut

	// ErrUnfinalizedTx indicates a transaction has not been finalized.
	// A valid block may only contain finalized transactions.
	ErrUnfinalizedTx

	// ErrDuplicateTx indicates a block contains an identical transaction
	// (or at least two transactions which hash to the same value).  A
	// valid block may only contain unique transactions.
	ErrDuplicateTx

	// ErrOverwriteTx indicates a block contains a transaction that has
	// the same hash as a previous transaction which has not been fully
	// spent.
	ErrOverwriteTx

	// ErrImmatureSpend indicates a transaction is attempting to spend a
	// coinbase that has not yet reached the required maturity.
	ErrImmatureSpend

	// ErrSpendTooHigh indicates a transaction is attempting to spend more
	// value than the sum of all of its inputs.
	ErrSpendTooHigh

	// ErrBadFees indicates the total fees for a block are invalid due to
	// exceeding the maximum possible value.
	ErrBadFees

	// ErrTooManySigOps indicates the total number of signature operations
	// for a transaction or block exceed the maximum allowed limits.
	ErrTooManySigOps

	// ErrFirstTxNotCoinbase indicates the first transaction in a block
	// is not a coinbase transaction.
	ErrFirstTxNotCoinbase

	// ErrMultipleCoinbases indicates a block contains more than one
	// coinbase transaction.
	ErrMultipleCoinbases

	// ErrBadCoinbaseScriptLen indicates the length of the signature script
	// for a coinbase transaction is not within the valid range.
	ErrBadCoinbaseScriptLen

	// ErrBadCoinbaseValue indicates the amount of a coinbase value does
	// not match the expected value of the subsidy plus the sum of all fees.
	ErrBadCoinbaseValue

	// ErrMissingCoinbaseHeight indicates the coinbase transaction for a
	// block does not start with the serialized block block height as
	// required for version 2 and higher blocks.
	ErrMissingCoinbaseHeight

	// ErrBadCoinbaseHeight indicates the serialized block height in the
	// coinbase transaction for version 2 and higher blocks does not match
	// the expected value.
	ErrBadCoinbaseHeight

	// ErrScriptMalformed indicates a transaction script is malformed in
	// some way.  For example, it might be longer than the maximum allowed
	// length or fail to parse.
	ErrScriptMalformed

	// ErrScriptValidation indicates the result of executing transaction
	// script failed.  The error covers any failure when executing scripts
	// such signature verification failures and execution past the end of
	// the stack.
	ErrScriptValidation

	// ErrUnexpectedWitness indicates that a block includes transactions
	// with witness data, but doesn't also have a witness commitment within
	// the coinbase transaction.
	ErrUnexpectedWitness

	// ErrInvalidWitnessCommitment indicates that a block's witness
	// commitment is not well formed.
	ErrInvalidWitnessCommitment

	// ErrWitnessCommitmentMismatch indicates that the witness commitment
	// included in the block's coinbase transaction doesn't match the
	// manually computed witness commitment.
	ErrWitnessCommitmentMismatch

	// ErrPreviousBlockUnknown indicates that the previous block is not known.
	ErrPreviousBlockUnknown

	// ErrInvalidAncestorBlock indicates that an ancestor of this block has
	// already failed validation.
	ErrInvalidAncestorBlock

	// ErrPrevBlockNotBest indicates that the block's previous block is not the
	// current chain tip. This is not a block validation rule, but is required
	// for block proposals submitted via getblocktemplate RPC.
	ErrPrevBlockNotBest
)
