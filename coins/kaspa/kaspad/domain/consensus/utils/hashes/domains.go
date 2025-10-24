package hashes

import (
	"crypto/sha256"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

const (
	transcationHashDomain         = "TransactionHash"
	transcationIDDomain           = "TransactionID"
	transcationSigningDomain      = "TransactionSigningHash"
	transcationSigningECDSADomain = "TransactionSigningHashECDSA"
	blockDomain                   = "BlockHash"
	proofOfWorkDomain             = "ProofOfWorkHash"
	heavyHashDomain               = "HeavyHash"
	merkleBranchDomain            = "MerkleBranchHash"
)

// transactionSigningECDSADomainHash is a hashed version of transcationSigningECDSADomain that is used
// to make it a constant size. This is needed because this domain is used by sha256 hash writer, and
// sha256 doesn't support variable size domain separation.
var transactionSigningECDSADomainHash = sha256.Sum256([]byte(transcationSigningECDSADomain))

// NewTransactionHashWriter Returns a new HashWriter used for transaction hashes
func NewTransactionHashWriter() HashWriter {
	blake, err := blake2b.New256([]byte(transcationHashDomain))
	if err != nil {
		panic(errors.Wrapf(err, "this should never happen. %s is less than 64 bytes", transcationHashDomain))
	}
	return HashWriter{blake}
}

// NewTransactionIDWriter Returns a new HashWriter used for transaction IDs
func NewTransactionIDWriter() HashWriter {
	blake, err := blake2b.New256([]byte(transcationIDDomain))
	if err != nil {
		panic(errors.Wrapf(err, "this should never happen. %s is less than 64 bytes", transcationIDDomain))
	}
	return HashWriter{blake}
}

// NewTransactionSigningHashWriter Returns a new HashWriter used for signing on a transaction
func NewTransactionSigningHashWriter() HashWriter {
	blake, err := blake2b.New256([]byte(transcationSigningDomain))
	if err != nil {
		panic(errors.Wrapf(err, "this should never happen. %s is less than 64 bytes", transcationSigningDomain))
	}
	return HashWriter{blake}
}

// NewTransactionSigningHashECDSAWriter Returns a new HashWriter used for signing on a transaction with ECDSA
func NewTransactionSigningHashECDSAWriter() HashWriter {
	hashWriter := HashWriter{sha256.New()}
	hashWriter.InfallibleWrite(transactionSigningECDSADomainHash[:])
	return hashWriter
}

// NewBlockHashWriter Returns a new HashWriter used for hashing blocks
func NewBlockHashWriter() HashWriter {
	blake, err := blake2b.New256([]byte(blockDomain))
	if err != nil {
		panic(errors.Wrapf(err, "this should never happen. %s is less than 64 bytes", blockDomain))
	}
	return HashWriter{blake}
}

// NewPoWHashWriter Returns a new HashWriter used for the PoW function
func NewPoWHashWriter() ShakeHashWriter {
	shake256 := sha3.NewCShake256(nil, []byte(proofOfWorkDomain))
	return ShakeHashWriter{shake256}
}

// NewHeavyHashWriter Returns a new HashWriter used for the HeavyHash function
func NewHeavyHashWriter() ShakeHashWriter {
	shake256 := sha3.NewCShake256(nil, []byte(heavyHashDomain))
	return ShakeHashWriter{shake256}
}

// NewMerkleBranchHashWriter Returns a new HashWriter used for a merkle tree branch
func NewMerkleBranchHashWriter() HashWriter {
	blake, err := blake2b.New256([]byte(merkleBranchDomain))
	if err != nil {
		panic(errors.Wrapf(err, "this should never happen. %s is less than 64 bytes", merkleBranchDomain))
	}
	return HashWriter{blake}
}
