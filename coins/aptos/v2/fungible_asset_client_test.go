package v2

import (
	"crypto/ed25519"
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test expects network connection to mainnet")
	}

	// Create a new Aptos client
	aptosClient, err := NewClient(MainnetConfig)
	assert.NoError(t, err)

	// Owner address
	ownerAddress := AccountAddress{}
	err = ownerAddress.ParseStringRelaxed(defaultOwner)
	assert.NoError(t, err)

	// TODO: This flow seems awkward and I made mistakes by running Parse on the same address multiple times
	metadataAddress := AccountAddress{}
	err = metadataAddress.ParseStringRelaxed(defaultMetadata)
	assert.NoError(t, err)

	primaryStoreAddress := AccountAddress{}
	err = primaryStoreAddress.ParseStringRelaxed(defaultStore)
	assert.NoError(t, err)

	// Create a fungible asset client
	cli, err := NewFungibleAssetClient(aptosClient, metadataAddress)
	assert.NoError(t, err)
	assert.Equal(t, metadataAddress.String(), cli.metadataAddress.String())

	// Primary store by direct access
	/*balance, err := faClient.Balance(primaryStoreAddress)
	assert.NoError(t, err)
	println("BALANCE: ", balance)

	name, err := faClient.Name()
	assert.NoError(t, err)
	println("NAME: ", name)
	symbol, err := faClient.Symbol()
	assert.NoError(t, err)
	println("Symbol: ", symbol)

	supply, err := faClient.Supply()
	assert.NoError(t, err)
	println("Supply: ", supply.String())

	maximum, err := faClient.Maximum()
	assert.NoError(t, err)
	println("Maximum: ", maximum.String())

	storeExists, err := faClient.StoreExists(primaryStoreAddress)
	assert.NoError(t, err)
	assert.True(t, storeExists)

	// This should hold
	storeNotExist, err := faClient.StoreExists(AccountOne)
	assert.NoError(t, err)
	assert.False(t, storeNotExist)

	storeMetadataAddress, err := faClient.StoreMetadata(primaryStoreAddress)
	assert.NoError(t, err)
	assert.Equal(t, metadataAddress, storeMetadataAddress)

	decimals, err := faClient.Decimals()
	assert.NoError(t, err)
	println("DECIMALS: ", decimals)

	storePrimaryStoreAddress, err := faClient.PrimaryStoreAddress(ownerAddress)
	assert.NoError(t, err)
	assert.Equal(t, primaryStoreAddress, storePrimaryStoreAddress)

	primaryStoreExists, err := faClient.PrimaryStoreExists(ownerAddress)
	assert.NoError(t, err)
	assert.True(t, primaryStoreExists)

	// Primary store by default
	primaryBalance, err := faClient.PrimaryBalance(ownerAddress)
	assert.NoError(t, err)
	println("PRIMARY BALANCE: ", primaryBalance)

	isFrozen, err := faClient.IsFrozen(primaryStoreAddress)
	assert.NoError(t, err)
	assert.False(t, isFrozen)*/
}

func TestTransferPrimaryStore(t *testing.T) {
	// Create a new Aptos client
	aptosClient, err := NewClient(TestnetConfig)
	assert.NoError(t, err)

	// sender
	key := crypto.Ed25519PrivateKey{}
	err = key.FromHex("9a2c5e9515410f90502fedbbb2ee2deb6eb51571d295ead093dcea8588c66ddb9a2c5e9515410f90502fedbbb2ee2deb6eb51571d295ead093dcea8588c66ddb")
	assert.NoError(t, err)
	senderAccount, err := NewAccountFromSigner(&key)
	assert.NoError(t, err)

	metadataAddress := AccountAddress{}
	err = metadataAddress.ParseStringRelaxed("0x9a2c5e9515410f90502fedbbb2ee2deb6eb51571d295ead093dcea8588c66ddb")
	assert.NoError(t, err)

	// Create a fungible asset client
	fungibleAssetClient, err := NewFungibleAssetClient(aptosClient, metadataAddress)
	assert.NoError(t, err)

	receiverAddress := AccountAddress{}
	err = receiverAddress.ParseStringRelaxed("0x45ec25f81a591f546a789d5c545d223ee40881daeba82bb77e9bde68ca6f4406")
	amount := uint64(1)
	expirationTimestampSeconds := 1717123597 + 300
	signedTx, err := fungibleAssetClient.TransferPrimaryStore(senderAccount, receiverAddress, amount, ChainIdOption(135), SequenceNumber(20), MaxGasAmount(200_000), GasUnitPrice(100), ExpirationSeconds(expirationTimestampSeconds))
	assert.NoError(t, err)
	ser := &bcs.Serializer{}
	signedTx.MarshalBCS(ser)
	err = ser.Error()
	assert.NoError(t, err)
	signedTxBs := ser.ToBytes()
	expected := "c0dc16d8ebf0153340e19ce3ebe1f6461b934da5c3bf62fcc4393b775eb4feea1400000000000000020000000000000000000000000000000000000000000000000000000000000001167072696d6172795f66756e6769626c655f73746f7265087472616e73666572010700000000000000000000000000000000000000000000000000000000000000010e66756e6769626c655f6173736574084d657461646174610003209a2c5e9515410f90502fedbbb2ee2deb6eb51571d295ead093dcea8588c66ddb2045ec25f81a591f546a789d5c545d223ee40881daeba82bb77e9bde68ca6f4406080100000000000000400d0300000000006400000000000000393b5966000000008700209a2c5e9515410f90502fedbbb2ee2deb6eb51571d295ead093dcea8588c66ddb4042b1eaf6dd889ad38af4f9719ece46b7b1217db415e6aeec407edce1dd6045ab95303a58fe053095be6e87074fc91884f7bedd7e54523e72e68838916581f50f"
	assert.Equal(t, expected, hex.EncodeToString(signedTxBs))
}

func TestTransferPrimaryStoreFromSeed(t *testing.T) {
	// Create a new Aptos client
	aptosClient, err := NewClient(TestnetConfig)
	assert.NoError(t, err)

	// sender
	pri, _ := hex.DecodeString("fbfdf86582998a7cf3b56f4a9026d3190e539d74fa27ddec23698f13a27c9ed7")
	seed := ed25519.NewKeyFromSeed(pri)
	key := crypto.Ed25519PrivateKey{}
	err = key.FromBytes(seed)

	senderAccount, err := NewAccountFromSigner(&key)
	assert.NoError(t, err)

	metadataAddress := AccountAddress{}
	err = metadataAddress.ParseStringRelaxed("0x56bb65b28d8bd323d7fc538109bfad5a0be55568d834415022d9e3cae8428791")
	assert.NoError(t, err)

	// Create a fungible asset client
	fungibleAssetClient, err := NewFungibleAssetClient(aptosClient, metadataAddress)
	assert.NoError(t, err)

	receiverAddress := AccountAddress{}
	err = receiverAddress.ParseStringRelaxed("0xda29ba63f2a675e7e18180f175bc8c982b9c02ecb7248b8afbe496c2f4bb96a7")
	amount := uint64(1200000)
	expirationTimestampSeconds := 1720084828
	signedTx, err := fungibleAssetClient.TransferPrimaryStore(senderAccount, receiverAddress, amount, ChainIdOption(2), SequenceNumber(34), MaxGasAmount(200_000), GasUnitPrice(100), ExpirationSeconds(expirationTimestampSeconds))
	assert.NoError(t, err)
	ser := &bcs.Serializer{}
	signedTx.MarshalBCS(ser)
	err = ser.Error()
	assert.NoError(t, err)
	signedTxBs := ser.ToBytes()
	expected := "986d10e6ffde60518ab0664f70339196f914f96255f31a1c6b226670c73d3f922200000000000000020000000000000000000000000000000000000000000000000000000000000001167072696d6172795f66756e6769626c655f73746f7265087472616e73666572010700000000000000000000000000000000000000000000000000000000000000010e66756e6769626c655f6173736574084d6574616461746100032056bb65b28d8bd323d7fc538109bfad5a0be55568d834415022d9e3cae842879120da29ba63f2a675e7e18180f175bc8c982b9c02ecb7248b8afbe496c2f4bb96a708804f120000000000400d03000000000064000000000000005c69866600000000020020a6084a32f84d8da8851aa2503222b8a6f5ea8d6d907f68d9a3e146da83725b05403fb53a6ba0957cb9f744accfe0ed72a1f7d815241e1c17a1b803d14baa1c8090b2954d2de17a6922a9e623cc4f27b89c6d483cc17aeba4d5ddb800b80994b70a"
	assert.Equal(t, expected, hex.EncodeToString(signedTxBs))
	//https://explorer.aptoslabs.com/txn/0x3fedbd3e9d431087cbdd79f7fdd1b9b0f1a7c30f8373836ecfea2f1857cff509?network=testnet
}
