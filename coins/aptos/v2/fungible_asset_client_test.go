package v2

import (
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
