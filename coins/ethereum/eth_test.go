package ethereum

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"

	"github.com/okx/go-wallet-sdk/coins/ethereum/token"
)

func TestEcRecoverPubKey(t *testing.T) {
	_, err := EcRecoverPubKey("", "", false)
	assert.Equal(t, err, errors.New("signature too short"))

	_, err = EcRecoverPubKey("d87758593e0b89f8a2deef5e053ce484fe971a75124bf5d89d6f4d4f586604120d0110d03c91260fec9ec917354caae50c1744d246e30ff48def277d7d9aec831b", "0x49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280", true)
	assert.Nil(t, err)
}

func TestCalTxHash(t *testing.T) {
	tx := NewEthDynamicFeeTx(big.NewInt(int64(11155111)),
		16,
		big.NewInt(int64(420000)),
		big.NewInt(int64(20000000000)),
		420000,
		big.NewInt(int64(1234)), "2de4898dd458d6dce097e29026d446300e3815fa", "", AccessList{})
	p, _ := hex.DecodeString("5dfce364a4e9020d1bc187c9c14060e1a2f8815b3b0ceb40f45e7e39eb122103")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	txStr, err := tx.SignTransaction(prvKey)
	assert.Nil(t, err)
	txHash, err := CalTxHash(txStr)
	assert.NoError(t, err)
	assert.Equal(t, txHash, "a3ae6d08481f8f9dff5c94a19dabfff70e186867459c8e201de9e6ae5b79dfb6")
}

func TestPubKeyToAddr(t *testing.T) {
	p, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	address := GetNewAddress(prvKey.PubKey())
	require.Equal(t, "0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", address)
	addr, err := PubKeyToAddr(prvKey.PubKey().SerializeUncompressed())
	require.Nil(t, err)
	require.Equal(t, address, addr)
}

func TestTransfer(t *testing.T) {
	p, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	address := GetNewAddress(prvKey.PubKey())
	require.Equal(t, "0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", address)

	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000000000)),
		"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", "0x",
	)
	hash, raw, _ := transaction.GetSigningHash(big.NewInt(int64(10)))
	require.Equal(t, "790f2b826ad9dfa7f2a53ec68e37ea51dc58652ecfde812da37c96a1069fcdbb", hash)
	require.Equal(t, "ed80852e90edd000830668a0941ca96f8cfe7276bb053b25e57188f1b5ec6a472888016345785d8a0000800a8080", raw)

	tx, err := transaction.SignTransaction(big.NewInt(int64(10)), prvKey)
	require.Nil(t, err)
	require.Equal(t, "0xf86d80852e90edd000830668a0941ca96f8cfe7276bb053b25e57188f1b5ec6a472888016345785d8a00008037a0afd10738449dd9ab4f95b6f49244dc076ae5f1251397c7f010ba529edecf8517a03eb5492b35278b2636870843550040edb60f6b1026bff42ee5a803c6de1b0e04", tx)

	b, _ := hex.DecodeString("2e3390fc71f35035b2ec378cced62632ef19c8d54b6b2f447e1f809c3d11ed0e")
	d, err := SignAsRecoverable(b, prvKey)
	require.Nil(t, err)
	signature := d.ToHex()
	require.Equal(t, "32466d55329625198458901517ccae23f0162fc42b333f770e8e59ab62d3d40e6c2e85072ad2fd4273d8be86af5e005b1c9df39bf3f2014897347ec81ce6bc7f1b", signature)
}

func TestTokenTransfer(t *testing.T) {
	p, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	transfer, _ := token.Transfer("0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", big.NewInt(int64(100000000000000000)))
	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(0)),
		"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", hex.EncodeToString(transfer),
	)
	tx, err := transaction.SignTransaction(big.NewInt(int64(10)), prvKey)
	require.Nil(t, err)
	require.Equal(t, "0xf8aa80852e90edd000830668a0941ca96f8cfe7276bb053b25e57188f1b5ec6a472880b844a9059cbb0000000000000000000000001ca96f8cfe7276bb053b25e57188f1b5ec6a4728000000000000000000000000000000000000000000000000016345785d8a000038a0ad7d69a4eeb889a2bdd82e2c62d4063467936350f7d3cc466aa513e7abcbb077a071b5b06e8253352f3e1aed57a6db4fdf5113b00c961fedecf7fbde96c94cb66f", tx)
}

func TestSignMessage(t *testing.T) {
	prv := "49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
	prvB, err := hex.DecodeString(prv)
	assert.NoError(t, err)
	msg := "0x49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
	prvKey, pub := btcec.PrivKeyFromBytes(prvB)
	addr := GetNewAddress(pub)
	sig, err := SignEthTypeMessage(msg, prvKey, true)
	assert.Equal(t, `d87758593e0b89f8a2deef5e053ce484fe971a75124bf5d89d6f4d4f586604120d0110d03c91260fec9ec917354caae50c1744d246e30ff48def277d7d9aec831b`, sig)
	addr2, err := EcRecover(sig, msg, true)
	assert.NoError(t, err)
	assert.Equal(t, addr2, addr)
	err = VerifySignMsg(sig, msg, addr, true)
	assert.NoError(t, err)
}

func TestEth2(t *testing.T) {
	p, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	prvKey, _ := btcec.PrivKeyFromBytes(p)

	chainId := big.NewInt(int64(10))
	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000000000)),
		"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", "",
	)
	unSignedHex, _ := transaction.UnSignedTx(chainId)
	unSignedBytes, _ := hex.DecodeString(unSignedHex)
	sig, err := SignMessage(unSignedBytes, prvKey)
	require.Nil(t, err)
	signedTx, err := transaction.SignedTx(chainId, sig)
	require.Nil(t, err)
	require.Equal(t, "0xf86d80852e90edd000830668a0941ca96f8cfe7276bb053b25e57188f1b5ec6a472888016345785d8a00008037a0afd10738449dd9ab4f95b6f49244dc076ae5f1251397c7f010ba529edecf8517a03eb5492b35278b2636870843550040edb60f6b1026bff42ee5a803c6de1b0e04", signedTx)
}

func TestEth3(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		tx, err := NewTransactionFromRaw("0xf86f6a8506fc23ac0082520894e1d4fd72a48af968d80f6d9ef161d57bb9293837880de0b6b3a764000080830150f5a0b64b1eb1c2f41b95dac35c4751f4070ca8c185b9a94ea2d44454f47ca7944a23a004a98084513c3233ec962d9cfdf2a45d06473fa51554a0b3d99939e1ed387ed7")
		txJson, err := json.Marshal(tx)
		require.Nil(t, err)
		assert.Equal(t, "{\"nonce\":106,\"gasPrice\":30000000000,\"gas\":21000,\"to\":\"4dT9cqSK+WjYD22e8WHVe7kpODc=\",\"value\":1000000000000000000,\"data\":\"\",\"v\":86261,\"r\":82453663816844830000916708780693313301330244341018929324864149013981109176867,\"s\":2108735539081020627590557504528962845333568746955363520031332889159139884759}", string(txJson))
	})
	t.Run("two", func(t *testing.T) {
		tx, err := NewTransactionFromRaw("0xf869698506fc23ac0082520894af133678d4188ddbfd13655cf12e8e15f28fdecb8203e880830150f6a0b64b1eb1c2f41b95dac35c4751f4070ca8c185b9a94ea2d44454f47ca7944a23a0522afd4359e208cf86c036c4272bf65d7fdddc73f33006f845c1c94fa826befc")
		require.Nil(t, err)
		txJson, err := json.Marshal(tx)
		assert.Equal(t, "{\"nonce\":105,\"gasPrice\":30000000000,\"gas\":21000,\"to\":\"rxM2eNQYjdv9E2Vc8S6OFfKP3ss=\",\"value\":1000,\"data\":\"\",\"v\":86262,\"r\":82453663816844830000916708780693313301330244341018929324864149013981109176867,\"s\":37165609118156479824344898797594740443460068889833947985890798108880694329084}", string(txJson))
	})
}

func TestEIP712(t *testing.T) {
	var typedData TypedData
	str := `{"domain":{"name":"AuthTransfer","chainId":1,"verifyingContract":"0x1243C09717e4441341472c4b142B8ac0B71F7672"},"message":{"details":[{"token":"0x0000000000000000000000000000000000000000","expiration":1853395200}],"spenders":["0x1B256B89462710a6b459540B999AbE5771d45A6e"],"nonce":0},"primaryType":"Permits","types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"Permits":[{"name":"details","type":"PermitDetails[]"},{"name":"spenders","type":"address[]"},{"name":"nonce","type":"uint256"}],"PermitDetails":[{"name":"token","type":"address"},{"name":"expiration","type":"uint256"}]}}`
	err := json.Unmarshal([]byte(str), &typedData)
	assert.NoError(t, err)
	hash, _, err := TypedDataAndHash(typedData)
	assert.NoError(t, err)
	assert.Equal(t, "3d697a8b530f96c6d7fc222ee6a43c7976ac2ac52dede33207a4758f5d502eac", hex.EncodeToString(hash))
}
