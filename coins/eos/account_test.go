package eos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/eoscanada/eos-go/ecc"
	atomic_market "github.com/okx/go-wallet-sdk/coins/eos/atomic-market"
	"github.com/okx/go-wallet-sdk/coins/eos/types"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

const (
	p1            = "5JvW9FSHci6MQcnoHjNnfv5T4Pfi5pj2weAEFQvq1TFaxs8Kbnt" // dubuqing1111 active key
	p2            = "5KD9qEUV3UT7KfipPkLW9wCfqdyAXTMmqcfKLq6egX3W6fmeV2g" // dubuqingfeng active key
	p4            = "5JFectyTfE5Uf34dU7T437Uenf8ysKvhKm7wmnc7fjiGYFc9t8J" // dubuqingfeng owner key
	n1            = "dubuqing1111"
	n2            = "dubuqingfeng"
	enabledPushTx = false
)

func TestGenerateKeyPair(t *testing.T) {
	gotPrivKey, gotPubKey := GenerateKeyPair()
	t.Log(gotPrivKey)
	privKey, err := ecc.NewPrivateKey(gotPrivKey)
	require.NoError(t, err)
	require.Equal(t, gotPubKey, privKey.PublicKey().String())
}

func TestCreateErrorAccountTransaction(t *testing.T) {
	tx := NewAccountTransaction("dubuqing111111111", "dubuqing4444", ecc.PublicKey{}, types.NewWAXAsset(100000000),
		types.NewWAXAsset(100000000), types.NewWAXAsset(100000000), false, nil)
	require.Nil(t, tx)
}

func TestCreateAccountTransaction(t *testing.T) {
	opt := &types.TxOptions{
		ChainID:     hexToChecksum256("f16b1833c747c43682f4386fca9cbb327929334a762755ebec17f6f23c9b8a12"),
		HeadBlockID: hexToChecksum256("0964056645a351623c9080af9bde15d54fcd630b84c55ac6f53f0e152f1fbf5b"),
	}
	gotPrivKey, _ := GenerateKeyPair()
	privateKey, err := ecc.NewPrivateKey(gotPrivKey)
	require.NoError(t, err)
	tx := NewAccountTransaction("test1", "test2", privateKey.PublicKey(), types.NewWAXAsset(100000000),
		types.NewWAXAsset(100000000), types.NewWAXAsset(100000000), false, opt)
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	t.Log("signedTx : ", signedTx)
	var packedTxBytes []byte
	var packedTxHex string
	packedTxBytes, err = json.Marshal(packedTx)
	require.NoError(t, err)
	packedTxHex = string(packedTxBytes)
	t.Log("packedTxHex : ", packedTxHex)
	hash, err := packedTx.ID()
	require.NoError(t, err)
	t.Log("hash : ", hash.String())
}

func TestNewAccount(t *testing.T) {
	userName := "test3"
	if len(userName) > 12 {
		return
	}
	gotPrivKey, _ := GenerateKeyPair()
	p, err := ecc.NewPrivateKey(gotPrivKey)
	require.NoError(t, err)
	creator := "eosio"
	ram := uint32(1000000)
	cpu := uint64(1000000)
	net := uint64(1000000)
	actions := []*types.Action{
		types.NewNewAccount(creator, userName, p.PublicKey()),
		types.NewBuyRAMBytes(creator, userName, ram),
		types.NewDelegateBW(
			creator,
			userName,
			types.NewEOSAsset(int64(cpu*10000)),
			types.NewEOSAsset(int64(net*10000)),
			false,
		),
	}
	chainId := []byte("e70aaab8997e1dfce58fbfac80cbbb8fecec7b99cf982a9444273cbc64c41473")
	opts := &types.TxOptions{
		ChainID: chainId,
	}
	tx := NewTransaction(actions, opts)
	if tx != nil {
		// sign the transaction
		signedTx, packedTx, err := SignTransaction(gotPrivKey, tx, chainId, types.CompressionNone)
		require.NoError(t, err)
		t.Log("signed transaction signatures : ", signedTx.Signatures)
		t.Log("broadcast transaction : ", packedTx.PackedTransaction)
	}
}
func hexToChecksum256(data string) types.Checksum256 {
	return types.Checksum256(hexToHexBytes(data))
}
func hexToHexBytes(data string) types.HexBytes {
	bytes, _ := hex.DecodeString(data)
	return types.HexBytes(bytes)
}

// note following test for other coins fork of EOS, such as WAX

var opt *types.TxOptions

func getPublicRPCEndpoint() string {
	return "http://api.waxtest.waxgalaxy.io/v1"
}

/*
	func getTxOptions() (*types.TxOptions, error) {
		url := getPublicRPCEndpoint() + "/chain/get_info"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		// http response
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// parse response
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}
		// check is exists response
		if _, ok := response["head_block_id"]; !ok {
			return nil, errors.New("head_block_id not found")
		}
		// check is exists response
		if _, ok := response["chain_id"]; !ok {
			return nil, errors.New("chain_id not found")
		}
		return &types.TxOptions{
			ChainID:     hexToChecksum256(response["chain_id"].(string)),
			HeadBlockID: hexToChecksum256(response["head_block_id"].(string)),
		}, nil
	}
*/
// curl https://api.waxtest.waxgalaxy.io/v1/chain/get_info
// get information on the blockchain WAX
func getTxOptions() (*types.TxOptions, error) {
	return &types.TxOptions{
		ChainID:     hexToChecksum256("f16b1833c747c43682f4386fca9cbb327929334a762755ebec17f6f23c9b8a12"),
		HeadBlockID: hexToChecksum256("0ed154c668a8a4f1916b30d9feead7e99edc9015635ba0feb3f8cd4368e8649b"),
	}, nil
}
func DumpGetRequiredKeyContent(tx *types.Transaction, publicKey string, t *testing.T) {
	m := make(map[string]interface{})
	m["transaction"] = tx
	m["available_keys"] = []string{publicKey}
	getRequiredKey, err := json.Marshal(m)
	require.NoError(t, err)
	t.Log("getRequiredKey : ", string(getRequiredKey))
}

func DumpPackedTx(t *testing.T, packedTx *types.PackedTransaction) {
	var packedTxBytes []byte
	var packedTxHex string
	var err error
	if packedTxBytes, err = json.Marshal(packedTx); err != nil {
		t.Error(err)
		return
	}
	packedTxHex = string(packedTxBytes)
	t.Log(packedTxHex)
	hash, _ := packedTx.ID()
	t.Log("Hash: ", hash)
}

func PushTransaction(packedTx *types.PackedTransaction) ([]byte, error) {
	url := getPublicRPCEndpoint() + "/chain/push_transaction"
	body, err := enc(packedTx)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// http response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// parse response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return nil, nil
}

func enc(v interface{}) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// TestWax test the wax send transaction
func TestNewWaxTransaction(t *testing.T) {
	privateKey, err := ecc.NewPrivateKey(p1)
	require.NoError(t, err)
	opt, err = getTxOptions()
	require.NoError(t, err)
	nilTx := NewTransactionWithParams(n1+"abcdefghijklmn", n2, "test", types.NewWAXAsset(500000000), opt)
	if nilTx != nil {
		t.Error("NewTransactionWithParams error")
		return
	}
	tx := NewTransactionWithParams(n1, n2, "test", types.NewWAXAsset(500000000), opt)
	require.NotNil(t, tx)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	t.Log("signedTx : ", signedTx)
	DumpPackedTx(t, packedTx)
	// send to the blockchain
	if enabledPushTx {
		PushTransaction(packedTx)
	}
}

func TestNewContractTransaction(t *testing.T) {
	privateKey, _ := ecc.NewPrivateKey(p1)
	opt, err := getTxOptions()
	require.NoError(t, err)
	contractName := "wax.token"
	tx := NewContractTransaction(contractName, n1, n2, "test", types.NewWAXAsset(500000000), opt)
	require.NotNil(t, tx)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	t.Log("signedTx : ", signedTx)
	DumpPackedTx(t, packedTx)
	// send to the blockchain
	if enabledPushTx {
		PushTransaction(packedTx)
	}
}

func TestNewBuyRAMBytesTransaction(t *testing.T) {
	privateKey, _ := ecc.NewPrivateKey(p1)
	opt, err := getTxOptions()
	require.NoError(t, err)
	nilTx := NewBuyRAMBytesTransaction(n1+"abcdefghijklmn", n2, 1000, opt)
	require.Nil(t, nilTx)
	tx := NewBuyRAMBytesTransaction(n1, n2, 1000, opt)
	if tx == nil {
		t.Error("NewBuyRAMBytesTransaction error")
		return
	}
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	t.Log("signedTx : ", signedTx)
	DumpPackedTx(t, packedTx)
	if enabledPushTx {
		PushTransaction(packedTx)
	}
}

func TestNewBuyRamTransaction(t *testing.T) {
	privateKey, _ := ecc.NewPrivateKey(p1)
	opt, err := getTxOptions()
	require.NoError(t, err)
	nilTx := NewBuyRamTransaction(n1+"abcdefghijklmn", n1, types.NewWAXAsset(500000000), opt)
	require.Nil(t, nilTx)
	tx := NewBuyRamTransaction(n1, n1, types.NewWAXAsset(500000000), opt)
	require.NotNil(t, tx)
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	t.Log("signedTx : ", signedTx)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	t.Log("signedTx : ", signedTx)
	DumpPackedTx(t, packedTx)
	if enabledPushTx {
		PushTransaction(packedTx)
	}
}

func TestNewSellRAMTransaction(t *testing.T) {
	privateKey, _ := ecc.NewPrivateKey(p2)
	opt, err := getTxOptions()
	require.NoError(t, err)
	nilTx := NewSellRAMTransaction(n2+"abcdefghijklmn", 1000, opt)
	require.Nil(t, nilTx)
	tx := NewSellRAMTransaction(n2, 100000, opt)
	require.NotNil(t, tx)
	signedTx, packedTx, err := SignTransaction(p2, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	t.Log("signedTx : ", signedTx)
	DumpPackedTx(t, packedTx)
	if enabledPushTx {
		PushTransaction(packedTx)
	}
}

func TestNewDelegateBWTransaction(t *testing.T) {
	privateKey, _ := ecc.NewPrivateKey(p1)
	opt, err := getTxOptions()
	require.NoError(t, err)
	nilTx := NewDelegateBWTransaction(n1+"abcdefghijklmn", n2, types.NewWAXAsset(500000000),
		types.NewWAXAsset(500000000), true, opt)
	require.Nil(t, nilTx)
	tx := NewDelegateBWTransaction(n1, n2, types.NewWAXAsset(500000000),
		types.NewWAXAsset(500000000), true, opt)
	require.NotNil(t, tx)
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	t.Log("signedTx : ", signedTx)
	DumpPackedTx(t, packedTx)
	if enabledPushTx {
		PushTransaction(packedTx)
	}
}

func TestNewUndelegateBWTransaction(t *testing.T) {
	privateKey, _ := ecc.NewPrivateKey(p2)
	opt, err := getTxOptions()
	require.NoError(t, err)
	nilTx := NewUndelegateBWTransaction(n2+"abcdefghijklmn", n2, types.NewWAXAsset(500000000),
		types.NewWAXAsset(500000000), opt)
	require.Nil(t, nilTx)
	tx := NewUndelegateBWTransaction(n2, n2, types.NewWAXAsset(400000000),
		types.NewWAXAsset(400000000), opt)
	require.NotNil(t, tx)
	signedTx, packedTx, err := SignTransaction(p2, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	t.Log("signedTx : ", signedTx)
	DumpPackedTx(t, packedTx)
	if enabledPushTx {
		PushTransaction(packedTx)
	}
}

func TestNewBuyNFTTransaction(t *testing.T) {
	opt, err := getTxOptions()
	require.NoError(t, err)
	privateKey, _ := ecc.NewPrivateKey(p1)
	saleId := uint64(32276)
	listingPriceToAssert := "2.00000000 WAX"
	nftId := uint64(1099532320836)
	actions := []*types.Action{
		atomic_market.NewPurchaseSale(n1, saleId, 0, ""),
		types.NewTransfer(n1, "atomicmarket", types.NewWAXAsset(100000000), "deposit"),
		atomic_market.NewAssertSale(n1, listingPriceToAssert, "", saleId, []uint64{nftId}),
	}
	tx := NewTransaction(actions, opt)
	require.NotNil(t, tx)
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	require.NoError(t, err)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	t.Log("signedTx : ", signedTx)
	DumpPackedTx(t, packedTx)
	if enabledPushTx {
		PushTransaction(packedTx)
	}
}

// TestWax test the wax send transaction, if opt is nil, will get the default options
func TestNewWaxNilOptTransaction(t *testing.T) {
	tx := NewTransactionWithParams(n1, n2, "test", types.NewWAXAsset(500000000), nil)
	require.NotNil(t, tx)
}
func setup() {
	opt, _ = getTxOptions()
}

func teardown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
func TestSignErrorTransaction(t *testing.T) {
	tx := NewTransactionWithParams(n1, n2, "test", types.NewWAXAsset(500000000), nil)
	require.NotNil(t, tx)
	_, _, err := SignTransaction(p1+"11dsf", tx, opt.ChainID, types.CompressionNone)
	require.NotNil(t, err)
}
