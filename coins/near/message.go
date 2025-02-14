package near

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/near/serialize"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

func SignMessage(payload *serialize.SignMessagePayload, privateKey string) (string, error) {
	pkBytes, err := serialize.TryParse(privateKey)
	if err != nil {
		return "", err
	}
	key := ed25519.PrivateKey(pkBytes)
	data, err := payload.Serialize()
	if err != nil {
		return "", err
	}
	txHash := sha256.Sum256(data)
	sig := ed25519.Sign(key, txHash[:])
	if len(sig) != 64 {
		return "", fmt.Errorf("sign error,length is not equal 64,length=%d", len(sig))
	}
	return base58.Encode(sig), nil
}
