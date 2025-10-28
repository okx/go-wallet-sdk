package mempool

const (
	// maxStandardP2SHSigOps is the maximum number of signature operations
	// that are considered standard in a pay-to-script-hash script.
	maxStandardP2SHSigOps = 15

	// maximumStandardSignatureScriptSize is the maximum size allowed for a
	// transaction input signature script to be considered standard. This
	// value allows for a 15-of-15 CHECKMULTISIG pay-to-script-hash with
	// compressed keys.
	//
	// The form of the overall script is: OP_0 <15 signatures> OP_PUSHDATA2
	// <2 bytes len> [OP_15 <15 pubkeys> OP_15 OP_CHECKMULTISIG]
	//
	// For the p2sh script portion, each of the 15 compressed pubkeys are
	// 33 bytes (plus one for the OP_DATA_33 opcode), and the thus it totals
	// to (15*34)+3 = 513 bytes. Next, each of the 15 signatures is a max
	// of 73 bytes (plus one for the OP_DATA_73 opcode). Also, there is one
	// extra byte for the initial extra OP_0 push and 3 bytes for the
	// OP_PUSHDATA2 needed to specify the 513 bytes for the script push.
	// That brings the total to 1+(15*74)+3+513 = 1627. This value also
	// adds a few extra bytes to provide a little buffer.
	// (1 + 15*74 + 3) + (15*34 + 3) + 23 = 1650
	maximumStandardSignatureScriptSize = 1650

	// MaximumStandardTransactionMass is the maximum mass allowed for transactions that
	// are considered standard and will therefore be relayed and considered for mining.
	MaximumStandardTransactionMass = 100_000
)
