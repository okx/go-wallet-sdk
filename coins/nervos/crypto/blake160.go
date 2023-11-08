package crypto

import "github.com/dchest/blake2b"

var ckbHashPersonalization = []byte("ckb-default-hash")

func Blake160(data []byte) ([]byte, error) {
	blake, err := Blake256(data)
	if err != nil {
		return nil, err
	}
	return blake[:20], nil
}

func Blake256(data []byte) ([]byte, error) {
	config := &blake2b.Config{
		Size:   32,
		Person: ckbHashPersonalization,
	}
	hash, err := blake2b.New(config)
	if err != nil {
		return nil, err
	}
	hash.Write(data)
	return hash.Sum(nil), nil
}
