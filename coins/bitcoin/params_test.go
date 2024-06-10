package bitcoin

import (
	"encoding/hex"
	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewZECAddr(t *testing.T) {
	pub, err := hex.DecodeString("030f85be34bee3303a14f0e39c5ed25cbccb14d7761e1a3093589d750297640397")
	assert.NoError(t, err)
	addr := NewZECAddr(pub)
	assert.Equal(t, "t1ZvCvAD8jH35dZAYgCVUShT69ksyUx9xDk", addr)
	addr2, err := btcutil.DecodeAddress("t1ZvCvAD8jH35dZAYgCVUShT69ksyUx9xDk", GetZECMainNetParams())
	assert.Error(t, err)
	assert.Nil(t, addr2)
}

func TestNewBsv(t *testing.T) {
	pub, err := hex.DecodeString("0303c3fdc94ba51bfe562d82116973f1da334cc65c6ea9ba59ee5af02e978d7731")
	assert.NoError(t, err)
	pk, err := btcec.ParsePubKey(pub)
	assert.NoError(t, err)
	addr, err := btcutil.NewAddressPubKey(pk.SerializeCompressed(), &chaincfg.MainNetParams)
	assert.NoError(t, err)
	assert.Equal(t, "1KY9sjgdeGEerBcGUM1QFbsU3sh8sYKzN2", addr.EncodeAddress())
	_, err = btcutil.DecodeAddress("1KY9sjgdeGEerBcGUM1QFbsU3sh8sYKzN2", &chaincfg.MainNetParams)
	assert.NoError(t, err)
}
