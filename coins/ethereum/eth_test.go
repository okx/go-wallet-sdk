package ethereum

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/okx/go-wallet-sdk/coins/ethereum/token"
	"github.com/okx/go-wallet-sdk/util"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
)

func TestEcRecover(t *testing.T) {
	_, err := EcRecoverPubKey("", "", false)
	assert.Equal(t, err, errors.New("signature too short"))

	sig := "d87758593e0b89f8a2deef5e053ce484fe971a75124bf5d89d6f4d4f586604120d0110d03c91260fec9ec917354caae50c1744d246e30ff48def277d7d9aec831b"
	msg := "0x49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"

	pubKey, err := EcRecoverPubKey(sig, msg, true)
	assert.NoError(t, err)
	assert.Equal(t, "04c847f6dd9e4fd3ce75c61614c838d9a54a5482b46e439b99aec5ebe26f9681510eab4e8116df5cb889d48194010633e83dd9ccbbffa6942a6768412293a70f41", util.EncodeHex(pubKey.SerializeUncompressed()))

	addr, err := EcRecover(sig, msg, true)
	assert.NoError(t, err)
	assert.Equal(t, "0xd74c65ad81aa8537327e9ba943011a8cec7a7b6b", addr)
}

func TestEcRecoverBytes(t *testing.T) {
	_, err := EcRecoverPubKeyBytes([]byte{0x01, 0x02}, []byte{0x01, 0x02}, false)
	assert.Equal(t, err, errors.New("signature too short"))

	sig, err := util.DecodeHexStringErr("d87758593e0b89f8a2deef5e053ce484fe971a75124bf5d89d6f4d4f586604120d0110d03c91260fec9ec917354caae50c1744d246e30ff48def277d7d9aec831b")
	assert.NoError(t, err)
	msg, err := util.DecodeHexStringErr("0x49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280")
	assert.NoError(t, err)

	pubKey, err := EcRecoverPubKeyBytes(sig, msg, true)
	assert.NoError(t, err)
	assert.Equal(t, "04c847f6dd9e4fd3ce75c61614c838d9a54a5482b46e439b99aec5ebe26f9681510eab4e8116df5cb889d48194010633e83dd9ccbbffa6942a6768412293a70f41", util.EncodeHex(pubKey))

	addr, err := EcRecoverBytes(sig, msg, true)
	assert.NoError(t, err)
	assert.Equal(t, "0xd74c65ad81aa8537327e9ba943011a8cec7a7b6b", addr)

}

func TestCalTxHash(t *testing.T) {
	t.Run("Dynamic Fee Tx", func(t *testing.T) {
		rawTx := "0x02f87083aa36a710830668a08504a817c800830668a0942de4898dd458d6dce097e29026d446300e3815fa8204d280c001a035d05e31efa695ad2a95f2e822ccffe07988aca85ebfe6076078c40069a49ca0a031ff765d66b6d8d4a2db9d70f49f3cd5dcff0e78c1f4c3feb8c89acd4d24079c"
		expectedHash := "a3ae6d08481f8f9dff5c94a19dabfff70e186867459c8e201de9e6ae5b79dfb6"
		txHash, err := CalTxHash(rawTx)
		assert.Nil(t, err)
		assert.Equal(t, expectedHash, util.EncodeHex(txHash))
	})

	t.Run("Legacy Tx", func(t *testing.T) {
		rawTx := "0xf86b0f8504a817c800830668a0942de4898dd458d6dce097e29026d446300e3815fa8204d2808401546d71a0a499527ec900a6e840e0aa82863449138e23aec3a3309fffa982a6857a83099da0587be28d4a00f344faad517207cb92315d3aa27cd4210905aa73561e65f31f69"
		expectedHash := "4427f2ecc1dfb3191e69e9405e68907af1163cdc5c46a2da84072144e55057c2"
		txHash, err := CalTxHash(rawTx)
		assert.Nil(t, err)
		assert.Equal(t, expectedHash, util.EncodeHex(txHash))
	})
}

func TestEth(t *testing.T) {
	p := util.DecodeHexString("5dfce364a4e9020d1bc187c9c14060e1a2f8815b3b0ceb40f45e7e39eb122103")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	address := GetNewAddress(prvKey.PubKey())
	assert.Equal(t, "0x2de4898dd458d6dce097e29026d446300e3815fa", address)

	transaction := NewEthTransaction(
		util.ToBigInt("0"),
		util.ToBigInt("420000"),
		util.ToBigInt("200000000000"),
		util.ToBigInt("100000000000000000"),
		"2de4898dd458d6dce097e29026d446300e3815fa", "0x",
	)

	hash, raw, _ := transaction.GetSigningHash(util.ToBigInt("10"))
	assert.Equal(t, "07f7adf7bf9efaf9442f792f3c7cd36b4505ec114c63493effa30b10b72d23e5", hash)
	assert.Equal(t, "ed80852e90edd000830668a0942de4898dd458d6dce097e29026d446300e3815fa88016345785d8a0000800a8080", raw)

	tx := transaction.SignTransaction(util.ToBigInt("10"), prvKey)
	assert.Equal(t, "0xf86d80852e90edd000830668a0942de4898dd458d6dce097e29026d446300e3815fa88016345785d8a00008038a0c115f04d7dd555746085b01a4c8add779353f7033560b8eb61b244fe37772138a00c333c52c61c2e1ad2af83343b5f1b396a8b0fb60252ff8529b30f2e90967580", tx)

	b, err := util.DecodeHexStringErr("2e3390fc71f35035b2ec378cced62632ef19c8d54b6b2f447e1f809c3d11ed0e")
	assert.NoError(t, err)
	d := SignAsRecoverable(b, prvKey)
	signature := d.ToHex()
	assert.Equal(t, "94f48c1dc793960d4c0e0ea0b34b95a8975d8d254edad8e25ab92a085914b37f30c943419fb1188244a8ffbdb7312003f8365ac381df5cd673a14cda28cf9b4f1b", signature)
}

func TestEthToken(t *testing.T) {
	p := util.DecodeHexString("5dfce364a4e9020d1bc187c9c14060e1a2f8815b3b0ceb40f45e7e39eb122103")
	prvKey, _ := btcec.PrivKeyFromBytes(p)

	transfer, _ := token.Transfer("2de4898dd458d6dce097e29026d446300e3815fa", util.ToBigInt("100000000000000000"))
	transaction := NewEthTransaction(
		util.ToBigInt("0"),
		util.ToBigInt("420000"),
		util.ToBigInt("200000000000"),
		util.ToBigInt("0"),
		"2ca70e7d0c396c36e8b9d206d988607a013483cf", util.EncodeHex(transfer),
	)
	tx := transaction.SignTransaction(util.ToBigInt("10"), prvKey)
	assert.Equal(t, "0xf8aa80852e90edd000830668a0942ca70e7d0c396c36e8b9d206d988607a013483cf80b844a9059cbb0000000000000000000000002de4898dd458d6dce097e29026d446300e3815fa000000000000000000000000000000000000000000000000016345785d8a000037a0afe573c296b30e9c4ef664ec64f63c48112e84167cdb8ef1ea567efc651f63c6a0283d2cc14a342464d2f62e83166bb9bf8d7ecdbe02d9b518eac02d458e158ccc", tx)
}

func TestSignMessage(t *testing.T) {
	prv := "49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
	prvB := util.DecodeHexString(prv)
	d := "0x49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
	prvKey, pub := btcec.PrivKeyFromBytes(prvB)
	addr := GetNewAddress(pub)
	sig, err := SignEthTypeMessage(d, prvKey, true)
	assert.NoError(t, err)
	assert.Equal(t, `d87758593e0b89f8a2deef5e053ce484fe971a75124bf5d89d6f4d4f586604120d0110d03c91260fec9ec917354caae50c1744d246e30ff48def277d7d9aec831b`, sig)
	addr2, err := EcRecover(sig, d, true)
	assert.NoError(t, err)
	assert.Equal(t, addr2, addr)
}

func TestEth2(t *testing.T) {
	p := util.DecodeHexString("5dfce364a4e9020d1bc187c9c14060e1a2f8815b3b0ceb40f45e7e39eb122103")
	prvKey, _ := btcec.PrivKeyFromBytes(p)

	chainId := util.ToBigInt("10")
	transaction := NewEthTransaction(
		util.ToBigInt("0"),
		util.ToBigInt("420000"),
		util.ToBigInt("200000000000"),
		util.ToBigInt("100000000000000000"),
		"2de4898dd458d6dce097e29026d446300e3815fa", "",
	)

	unSignedHex := transaction.UnSignedTx(chainId)
	unSignedBytes, err := util.DecodeHexStringErr(unSignedHex)
	assert.NoError(t, err)
	sig := SignMessage(unSignedBytes, prvKey)
	tx := transaction.SignedTx(chainId, sig)
	assert.Equal(t, "0xf86d80852e90edd000830668a0942de4898dd458d6dce097e29026d446300e3815fa88016345785d8a00008038a0c115f04d7dd555746085b01a4c8add779353f7033560b8eb61b244fe37772138a00c333c52c61c2e1ad2af83343b5f1b396a8b0fb60252ff8529b30f2e90967580", tx)
}

func TestEIP712(t *testing.T) {
	typedData := TypedData{}
	str := `{"domain":{"name":"AuthTransfer","chainId":1,"verifyingContract":"0x1243C09717e4441341472c4b142B8ac0B71F7672"},"message":{"details":[{"token":"0x0000000000000000000000000000000000000000","expiration":1853395200}],"spenders":["0x1B256B89462710a6b459540B999AbE5771d45A6e"],"nonce":0},"primaryType":"Permits","types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"Permits":[{"name":"details","type":"PermitDetails[]"},{"name":"spenders","type":"address[]"},{"name":"nonce","type":"uint256"}],"PermitDetails":[{"name":"token","type":"address"},{"name":"expiration","type":"uint256"}]}}`
	err := json.Unmarshal([]byte(str), &typedData)
	assert.NoError(t, err)
	hash, _, err := TypedDataAndHash(typedData)
	assert.NoError(t, err)
	assert.Equal(t, "3d697a8b530f96c6d7fc222ee6a43c7976ac2ac52dede33207a4758f5d502eac", util.EncodeHex(hash))
}

func TestSignAndGenerateRawTransaction(t *testing.T) {
	t.Run("Legacy Transaction", func(t *testing.T) {
		privHex := `f4d79cecc34de14e8b43e7779acaa350060513937f420f5b91ab7f483cac6b72`
		p := util.DecodeHexString(privHex)
		prvKey, _ := btcec.PrivKeyFromBytes(p)

		to := "0x05d132975D8EfCD67262980C54f9030319C91Af0"
		value := util.ToBigInt("404993102026570")
		nonce := uint64(0)
		gasPrice := util.ToBigInt("20000000000")
		gasLimit := uint64(80000)

		tx := NewEthTransaction(
			new(big.Int).SetUint64(nonce),
			new(big.Int).SetUint64(gasLimit),
			gasPrice,
			value,
			to,
			"",
		)
		evmTx := &EVMTx{
			TxType:  0,
			ChainId: util.ToBigInt("1"),
			Tx:      tx,
		}

		// Test SignRawTransactionWithPrivateKey
		expectedSigned := "0xf86c808504a817c800830138809405d132975d8efcd67262980c54f9030319c91af087017056cdfb974a8026a03fe8e4e8ebdd1c3bfcca09b48dda58ca8427fdda81393f1582241fac02ebc51aa01713f316b4f32700d9b7535e4dcafb90a4833a78ec32f465013af95aa5d191b3"
		signedTx, err := SignTx(evmTx, prvKey)
		assert.Nil(t, err)
		assert.Equal(t, expectedSigned, signedTx)

		// Test GenerateUnsignedRawTransaction
		unsignedRawTx, err := GenUnsignedTx(evmTx)
		assert.Nil(t, err)

		// Manually provide V, R, S from pre-calculated valid signature
		V := "26"
		R := "3fe8e4e8ebdd1c3bfcca09b48dda58ca8427fdda81393f1582241fac02ebc51a"
		S := "1713f316b4f32700d9b7535e4dcafb90a4833a78ec32f465013af95aa5d191b3"
		reconstructedSignedTx, err := GenTxWithSig(0, "1", unsignedRawTx, R, S, V)
		assert.Nil(t, err)
		assert.Equal(t, util.RemoveHexPrefix(signedTx), util.RemoveHexPrefix(reconstructedSignedTx))
	})

	t.Run("EIP1559 Transaction", func(t *testing.T) {
		privHex := "12a82ca8fc838ba03427f4285d553ba26c178832de7aba1c02686f25c1b6bffd"
		p := util.DecodeHexString(privHex)
		prvKey, _ := btcec.PrivKeyFromBytes(p)

		nonce := uint64(1)
		toHex := "0x05d132975d8efcd67262980c54f9030319c91af0"
		toAddr := common.HexToAddress(toHex)
		value := util.ToBigInt("1000000000")
		gasLimit := uint64(21000)
		maxFeePerGas := util.ToBigInt("20000000000")
		maxPriorityFeePerGas := util.ToBigInt("3000000000")

		tx1559 := types.NewTx(&types.DynamicFeeTx{
			Nonce:      nonce,
			To:         &toAddr,
			Value:      value,
			Gas:        gasLimit,
			GasFeeCap:  maxFeePerGas,
			GasTipCap:  maxPriorityFeePerGas,
			Data:       []byte{},
			AccessList: types.AccessList{},
		})

		evmTx := &EVMTx{
			TxType:  DynamicFeeTxType,
			ChainId: util.ToBigInt("1"),
			Tx1559:  tx1559,
		}

		expectedSigned := "0x02f86f010184b2d05e008504a817c8008252089405d132975d8efcd67262980c54f9030319c91af0843b9aca0080c001a0e4a152cf60089b026c24ac4068c2172529010dd0d6fb98ed2b096667acae2de0a01c88ccc3e77034c5e939541dc01bc93542590e3463fc1442baf599589e45b97c"
		signedTx, err := SignTx(evmTx, prvKey)
		assert.Nil(t, err)
		assert.Equal(t, expectedSigned, signedTx)

		unsignedRawTx, err := GenUnsignedTx(evmTx)
		assert.Nil(t, err)

		V := "01"
		R := "e4a152cf60089b026c24ac4068c2172529010dd0d6fb98ed2b096667acae2de0"
		S := "1c88ccc3e77034c5e939541dc01bc93542590e3463fc1442baf599589e45b97c"

		reconstructedSignedTx, err := GenTxWithSig(DynamicFeeTxType, "1", unsignedRawTx, R, S, V)
		assert.Nil(t, err)
		assert.Equal(t, util.RemoveHexPrefix(signedTx), util.RemoveHexPrefix(reconstructedSignedTx))
	})

	t.Run("EIP7702 Transaction", func(t *testing.T) {
		privHex := "49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
		p := util.DecodeHexString(privHex)
		prvKey, _ := btcec.PrivKeyFromBytes(p)
		nonce := uint64(0)
		toHex := "0xd74c65ad81aa8537327e9ba943011a8cec7a7b6b"
		toAddr := common.HexToAddress(toHex)
		value := util.ToBigInt("0")
		gasLimit := uint64(100000)
		maxFeePerGas := util.ToBigInt("10000")
		maxPriorityFeePerGas := util.ToBigInt("10000")

		tx1559 := types.NewTx(&types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &toAddr,
			Value:     value,
			Gas:       gasLimit,
			GasFeeCap: maxFeePerGas,
			GasTipCap: maxPriorityFeePerGas,
			Data:      []byte{},
		})

		nonceAuth := util.ToBigInt("1")
		chainIdAuth := util.ToBigInt("1")
		yParityAuth := util.ToBigInt("0")
		rAuth := util.ToBigInt("0xec2a15dad70e392a3c82ac988b3d9247d871c4c83c5c355d66d9b3df9bf58a46")
		sAuth := util.ToBigInt("0x67946771adff6ecbc2685454aac7057152c0550a27a81fb1d070236a6bd56462")
		addressBytes, err := util.DecodeHexStringErr("0x89aFB3EF13c03D0A816D6CDC20fdC21a915a4c24")
		assert.NoError(t, err)
		auth := &EthAuthorization{
			Address: addressBytes,
			Nonce:   nonceAuth,
			ChainId: chainIdAuth,
			YParity: yParityAuth,
			R:       rAuth,
			S:       sAuth,
		}
		authList := []*EthAuthorization{auth}

		evmTx := &EVMTx{
			TxType:            AuthorizationTxType,
			ChainId:           util.ToBigInt("1"),
			Tx1559:            tx1559,
			AuthorizationList: authList,
		}
		expectedSigned := "0x04f8c50180822710822710830186a094d74c65ad81aa8537327e9ba943011a8cec7a7b6b8080c0f85cf85a019489afb3ef13c03d0a816d6cdc20fdc21a915a4c240180a0ec2a15dad70e392a3c82ac988b3d9247d871c4c83c5c355d66d9b3df9bf58a46a067946771adff6ecbc2685454aac7057152c0550a27a81fb1d070236a6bd5646201a0393d384336222a4b0085cc75a66e001113e6d60603eddf4d7517c46c48d36b35a07bf33b2b988fa770b0713342116643364fc817b246bf463de81f2bb33bd7f9f9"
		signedTx, err := SignTx(evmTx, prvKey)
		assert.Nil(t, err)
		assert.Equal(t, expectedSigned, signedTx)

		unsignedRawTx, err := GenUnsignedTx(evmTx)
		assert.Nil(t, err)

		V := "01"
		R := "393d384336222a4b0085cc75a66e001113e6d60603eddf4d7517c46c48d36b35"
		S := "7bf33b2b988fa770b0713342116643364fc817b246bf463de81f2bb33bd7f9f9"

		reconstructedSignedTx, err := GenTxWithSig(AuthorizationTxType, "1", unsignedRawTx, R, S, V)
		assert.Nil(t, err)
		assert.Equal(t, util.RemoveHexPrefix(signedTx), util.RemoveHexPrefix(reconstructedSignedTx))
	})
}

func TestCalcSignHash(t *testing.T) {
	data := []byte("hello world")

	// With prefix
	hashWithPrefix := CalcSignHash(data, true)
	expectedMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	s := sha3.NewLegacyKeccak256()
	s.Write([]byte(expectedMsg))
	expectedHash := s.Sum(nil)
	assert.Equal(t, expectedHash, hashWithPrefix)

	// Without prefix
	hashWithoutPrefix := CalcSignHash(data, false)
	assert.Equal(t, data, hashWithoutPrefix)
}

func TestGenerateSigningHash(t *testing.T) {
	txHex := "ed80852e90edd000830668a0942de4898dd458d6dce097e29026d446300e3815fa88016345785d8a0000800a8080"
	hash, err := CalTxHash(txHex)
	assert.Nil(t, err)

	expectedHashHex := "07f7adf7bf9efaf9442f792f3c7cd36b4505ec114c63493effa30b10b72d23e5"
	assert.Equal(t, expectedHashHex, util.EncodeHex(hash))
}

func TestVerifySignMsg(t *testing.T) {
	sig := "d87758593e0b89f8a2deef5e053ce484fe971a75124bf5d89d6f4d4f586604120d0110d03c91260fec9ec917354caae50c1744d246e30ff48def277d7d9aec831b"
	msg := "0x49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"

	prvHex := "49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
	prvB := util.DecodeHexString(prvHex)
	_, pub := btcec.PrivKeyFromBytes(prvB)
	addr := GetNewAddress(pub)

	err := VerifySignMsg(sig, msg, addr, true)
	assert.Nil(t, err)

	err = VerifySignMsg(sig, msg, "0xSomeOtherAddress", true)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid sign", err.Error())
}

func TestValidateAddress(t *testing.T) {
	valid := ValidateAddress("0xe688b84b23f322a994A53dbF8E15FA82CDB71127")
	assert.True(t, valid)
}
