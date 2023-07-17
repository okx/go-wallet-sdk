package bip32

import (
	bip32 "github.com/tyler-smith/go-bip32"
	"strconv"
	"strings"
)

const (
	FirstHardenedChild = bip32.FirstHardenedChild

	PublicKeyCompressedLength = bip32.PublicKeyCompressedLength
)

type Key struct {
	*bip32.Key
}

func NewMasterKey(seed []byte) (*Key, error) {
	k, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	return &Key{k}, nil
}

func (key *Key) NewChildKeyByChainId(id uint32) (*Key, error) {
	return key.NewChildKeyByPath(bip32.FirstHardenedChild+44, id|bip32.FirstHardenedChild, bip32.FirstHardenedChild, 0, 0)
}

func (key *Key) NewChildKeyByPath(childPath ...uint32) (*Key, error) {
	currentKey := key.Key
	for _, childIdx := range childPath {
		newKey, err := currentKey.NewChildKey(childIdx)
		if err != nil {
			return nil, err
		}
		currentKey = newKey
	}
	return &Key{currentKey}, nil
}

func (key *Key) NewChildKeyByPathString(childPath string) (*Key, error) {
	arr := strings.Split(childPath, "/")
	currentKey := key.Key
	for _, part := range arr {
		if part == "m" {
			continue
		}

		var harden = false
		if strings.HasSuffix(part, "'") {
			harden = true
			part = strings.TrimSuffix(part, "'")
		}

		id, err := strconv.ParseUint(part, 10, 31)
		if err != nil {
			return nil, err
		}

		var uid = uint32(id)
		if harden {
			uid |= bip32.FirstHardenedChild
		}

		newKey, err := currentKey.NewChildKey(uid)
		if err != nil {
			return nil, err
		}
		currentKey = newKey
	}
	return &Key{currentKey}, nil
}
