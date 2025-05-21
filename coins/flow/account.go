package flow

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/coins/flow/core"
	"github.com/okx/go-wallet-sdk/util"
	"golang.org/x/crypto/sha3"
)

func GenerateKeyPair() (privKey, pubKey string) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubKeyBytes := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return hex.EncodeToString(privateKey.D.Bytes()), hex.EncodeToString(pubKeyBytes)
}

func SignTx(signerAddr, privKeyHex string, tx *core.Transaction) error {
	envelopeMessage := tx.EnvelopeMessage()
	transactionDomainTag := new([32]byte)
	copy(transactionDomainTag[:], "FLOW-V0.0-transaction")
	message := append(transactionDomainTag[:], envelopeMessage...)
	hashBytes := hashSha256(message)
	sig, err := signEcdsaP256(hashBytes, privKeyHex)
	if err != nil {
		return err
	}
	signature := core.TransactionSignature{
		Address:     core.HexToAddress(signerAddr),
		SignerIndex: 0,
		KeyIndex:    0,
		Signature:   sig,
	}
	tx.EnvelopeSignatures = []core.TransactionSignature{signature}
	return nil
}

func hashSha256(message []byte) []byte {
	hasher := sha3.New256()
	hasher.Write(message)
	return hasher.Sum(nil)
}

func signEcdsaP256(hash []byte, privateKeyHex string) ([]byte, error) {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	privateKeyEcdsa := privateKey.ToECDSA()
	r, s, err := ecdsa.Sign(rand.Reader, privateKeyEcdsa, hash)
	if err != nil {
		return nil, err
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	Nlen := bitsToBytes((privateKeyEcdsa.PublicKey.Curve.Params().N).BitLen())
	signature := make([]byte, 2*Nlen)
	// pad the signature with zeroes
	copy(signature[Nlen-len(rBytes):], rBytes)
	copy(signature[2*Nlen-len(sBytes):], sBytes)
	return signature, nil
}

func bitsToBytes(bits int) int {
	return (bits + 7) >> 3
}
func ValidateAddress(address string) bool {
	bytes, err := util.DecodeHexString(address)
	return err == nil && len(bytes) == 8
}
