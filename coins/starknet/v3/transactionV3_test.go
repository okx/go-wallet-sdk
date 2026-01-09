package v3

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/okx/go-wallet-sdk/coins/starknet"
)

const pri = "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da"

func TestDeployTxV3(t *testing.T) {
	curve := starknet.SC()

	starkPub, err := starknet.GetPubKey(curve, pri)
	assert.NoError(t, err)

	nonce := big.NewInt(0)

	tip := "0x0"
	l1DataGasMaxAmount := doubleHex("0x100")
	l1DataMaxPricePerUnit := doubleHex("0x65a")

	l1GasMaxAmount := doubleHex("0x0")
	l1GasMaxPricePerUnit := doubleHex("0xa6670a161547")

	l2GasMaxAmount := doubleHex("0xa2bc0")
	l2GasMaxPricePerUnit := doubleHex("0x10d23e189")

	tx, err := CreateSignedDeployAccountTxV3(curve, starkPub, starknet.OKXAccountClassHashCairo1, "", nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, starknet.SEPOLIA_ID, pri)
	assert.NoError(t, err)

	res := tx.GetTxRequestJson()
	assert.NotEmpty(t, res)

	hash, err := GetTxHashV3(res, starknet.SEPOLIA_ID)
	assert.NoError(t, err)
	assert.Equal(t, "{\"type\":\"DEPLOY_ACCOUNT\",\"version\":\"0x3\",\"signature\":[\"0x165484e33d6d730d426196a0b707d7f5be608a3561333fb7637d50aa45c0bfb\",\"0x86c546ffbc488d971d1f7e1408b57699422400895c83882cd12f0726b99793\"],\"nonce\":\"0x0\",\"contract_address_salt\":\"0x72ff9867ba607f204042c328cde87ddefe405b830e6515563fbe3ced9342109\",\"constructor_calldata\":[\"0x72ff9867ba607f204042c328cde87ddefe405b830e6515563fbe3ced9342109\",\"0x0\"],\"class_hash\":\"0x1c0bb51e2ce73dc007601a1e7725453627254016c28f118251a71bbb0507fcb\",\"resource_bounds\":{\"L1_GAS\":{\"max_amount\":\"0x0\",\"max_price_per_unit\":\"0x14cce142c2a8e\"},\"L1_DATA_GAS\":{\"max_amount\":\"0x200\",\"max_price_per_unit\":\"0xcb4\"},\"L2_GAS\":{\"max_amount\":\"0x145780\",\"max_price_per_unit\":\"0x21a47c312\"}},\"tip\":\"0x0\",\"paymaster_data\":[],\"nonce_data_availability_mode\":0,\"fee_data_availability_mode\":0}", res)
	assert.Equal(t, "0x302dac8a8d27dd6b6b125a1d8c01b982c08f8519ebfeb37b226522e1ae4f420", starknet.BigToHex(tx.GetTxHash()))
	assert.Equal(t, "0x302dac8a8d27dd6b6b125a1d8c01b982c08f8519ebfeb37b226522e1ae4f420", hash)
}

func TestCreateTransferTxV3(t *testing.T) {
	curve := starknet.SC()
	contractAddr := starknet.STARK
	from := "0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d"
	to := "0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d"
	amount := big.NewInt(100)
	nonce := big.NewInt(5)

	tip := doubleHex("0x0")
	l1GasMaxAmount := doubleHex("0x0")
	l1GasMaxPricePerUnit := doubleHex("0x12fb799a8d45")

	l2GasMaxAmount := doubleHex("0xded00")
	l2GasMaxPricePerUnit := doubleHex("0x1f19c38d")

	l1DataGasMaxAmount := doubleHex("0x80")
	l1DataMaxPricePerUnit := doubleHex("0x130293f")

	tx, err := CreateSignedTransferTxV3(curve, contractAddr, from, to, amount, nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, starknet.MAINNET_ID, 1, "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da")
	assert.NoError(t, err)

	res := tx.GetTxRequestJson()
	assert.NotEmpty(t, res)

	hash, err := GetTxHashV3(res, starknet.MAINNET_ID)
	assert.NoError(t, err)

	assert.Equal(t, "{\"type\":\"INVOKE_FUNCTION\",\"sender_address\":\"0x109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d\",\"calldata\":[\"0x1\",\"0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d\",\"0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e\",\"0x3\",\"0x109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d\",\"0x64\",\"0x0\"],\"version\":\"0x3\",\"signature\":[\"0x4cd8e27ee572dd5e71025034865901e2361d5e89bac6511a24d234609d22fa2\",\"0x579017c150ed4e81110ae778083dbacb707f7c4dd8f9d42beba2055d5d8bf5f\"],\"nonce\":\"0x5\",\"resource_bounds\":{\"L1_GAS\":{\"max_amount\":\"0x0\",\"max_price_per_unit\":\"0x25f6f3351a8a\"},\"L1_DATA_GAS\":{\"max_amount\":\"0x100\",\"max_price_per_unit\":\"0x260527e\"},\"L2_GAS\":{\"max_amount\":\"0x1bda00\",\"max_price_per_unit\":\"0x3e33871a\"}},\"tip\":\"0x0\",\"paymaster_data\":[],\"account_deployment_data\":[],\"nonce_data_availability_mode\":0,\"fee_data_availability_mode\":0}", res)
	assert.Equal(t, "0x5f8ef99f579f14b71e3d95d51c1053eb310623c6c0b4f3b150290a73b2590b2", starknet.BigToHex(tx.GetTxHash()))
	assert.Equal(t, "0x5f8ef99f579f14b71e3d95d51c1053eb310623c6c0b4f3b150290a73b2590b2", hash)
}

func TestCreateMultiContractV3(t *testing.T) {
	curve := starknet.SC()
	from := "0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d"
	txs := []starknet.Calls{
		{
			ContractAddress: starknet.STARK,
			Entrypoint:      "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
			Calldata:        []string{"0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d", "0x64", "0x0"},
		}, {
			ContractAddress: starknet.STARK,
			Entrypoint:      "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
			Calldata:        []string{"0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d", "0x64", "0x0"},
		},
	}
	nonce := big.NewInt(6)

	tip := doubleHex("0x0")
	l1GasMaxAmount := doubleHex("0x0")
	l1GasMaxPricePerUnit := doubleHex("0x119e72c54321")

	l2GasMaxAmount := doubleHex("0x134ac0")
	l2GasMaxPricePerUnit := doubleHex("0x1cddeb25")

	l1DataGasMaxAmount := doubleHex("0x80")
	l1DataMaxPricePerUnit := doubleHex("0xa9fb76ef14")

	tx, err := CreateSignedMultiContractTxV3(curve, from, txs, nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, starknet.MAINNET_ID, 1, "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da")
	assert.NoError(t, err)

	res := tx.GetTxRequestJson()
	assert.NotEmpty(t, res)

	hash, err := GetTxHashV3(res, starknet.MAINNET_ID)
	assert.NoError(t, err)

	assert.Equal(t, "{\"type\":\"INVOKE_FUNCTION\",\"sender_address\":\"0x109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d\",\"calldata\":[\"0x2\",\"0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d\",\"0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e\",\"0x3\",\"0x109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d\",\"0x64\",\"0x0\",\"0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d\",\"0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e\",\"0x3\",\"0x109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d\",\"0x64\",\"0x0\"],\"version\":\"0x3\",\"signature\":[\"0x7e2c479888b4ebac22f7e38f6ec83f2f30d87fbf5eb64ef6557f8d3ad9a207e\",\"0x11c8fad4a2294f4701a12f0d636655c45ba803c6d3d249d55d0cc93329baf59\"],\"nonce\":\"0x6\",\"resource_bounds\":{\"L1_GAS\":{\"max_amount\":\"0x0\",\"max_price_per_unit\":\"0x233ce58a8642\"},\"L1_DATA_GAS\":{\"max_amount\":\"0x100\",\"max_price_per_unit\":\"0x153f6edde28\"},\"L2_GAS\":{\"max_amount\":\"0x269580\",\"max_price_per_unit\":\"0x39bbd64a\"}},\"tip\":\"0x0\",\"paymaster_data\":[],\"account_deployment_data\":[],\"nonce_data_availability_mode\":0,\"fee_data_availability_mode\":0}", res)
	assert.Equal(t, "0x62b2e3ac261bad8181a95398a9032a6480cab7677a153b5a1936715fbf272e1", starknet.BigToHex(tx.GetTxHash()))
	assert.Equal(t, "0x62b2e3ac261bad8181a95398a9032a6480cab7677a153b5a1936715fbf272e1", hash)
}

func TestDeployTxV3Cairo0(t *testing.T) {
	curve := starknet.SC()
	pri := "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da"

	starkPub, err := starknet.GetPubKey(curve, pri)
	assert.NoError(t, err)

	nonce := big.NewInt(0)

	tip := "0x0"
	l1DataGasMaxAmount := doubleHex("0x160")
	l1DataMaxPricePerUnit := doubleHex("0x1e37433d")

	l1GasMaxAmount := doubleHex("0x0")
	l1GasMaxPricePerUnit := doubleHex("0x1347ba7d4908")

	l2GasMaxAmount := doubleHex("0xbec80")
	l2GasMaxPricePerUnit := doubleHex("0x1f96b292")

	tx, err := CreateSignedDeployAccountTxV3(curve, starkPub, starknet.OKXAccountClassHashCairo0, starknet.OKXProxyAccountClassHashCairo0, nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, starknet.MAINNET_ID, pri)
	assert.NoError(t, err)

	res := tx.GetTxRequestJson()
	assert.NotEmpty(t, res)

	hash, err := GetTxHashV3(res, starknet.MAINNET_ID)
	assert.NoError(t, err)

	assert.Equal(t, `{"type":"DEPLOY_ACCOUNT","version":"0x3","signature":["0x450af1a42899a77683e56cac6e1464502e4f98a8bf07effe6e08f3377ce8e05","0x35e77eebe849e5949d17879112ae06d4166098aaa2f22145176e2c3fdff92ec"],"nonce":"0x0","contract_address_salt":"0x72ff9867ba607f204042c328cde87ddefe405b830e6515563fbe3ced9342109","constructor_calldata":["0x309c042d3729173c7f2f91a34f04d8c509c1b292d334679ef1aabf8da0899cc","0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463","0x2","0x72ff9867ba607f204042c328cde87ddefe405b830e6515563fbe3ced9342109","0x0"],"class_hash":"0x3530cc4759d78042f1b543bf797f5f3d647cde0388c33734cf91b7f7b9314a9","resource_bounds":{"L1_GAS":{"max_amount":"0x0","max_price_per_unit":"0x268f74fa9210"},"L1_DATA_GAS":{"max_amount":"0x2c0","max_price_per_unit":"0x3c6e867a"},"L2_GAS":{"max_amount":"0x17d900","max_price_per_unit":"0x3f2d6524"}},"tip":"0x0","paymaster_data":[],"nonce_data_availability_mode":0,"fee_data_availability_mode":0}`, res)
	assert.Equal(t, "0xdd3b97470a24a9552a2446adcfed9ba966d06e2e8d53e5ee3e7ced5ee7938e", starknet.BigToHex(tx.GetTxHash()))
	assert.Equal(t, "0xdd3b97470a24a9552a2446adcfed9ba966d06e2e8d53e5ee3e7ced5ee7938e", hash)
}

func TestCreateTransferTxV3Cairo0(t *testing.T) {
	curve := starknet.SC()
	contractAddr := starknet.STARK
	from := "0x3a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034"
	to := "0x3a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034"
	amount := big.NewInt(100)
	nonce := big.NewInt(1)

	tip := doubleHex("0x0")
	l1GasMaxAmount := doubleHex("0x0")
	l1GasMaxPricePerUnit := doubleHex("0x1094fefd0176")

	l2GasMaxAmount := doubleHex("0xded00")
	l2GasMaxPricePerUnit := doubleHex("0x1b2b0064")

	l1DataGasMaxAmount := doubleHex("0x80")
	l1DataMaxPricePerUnit := doubleHex("0xa721590a9c")

	tx, err := CreateSignedTransferTxV3(curve, contractAddr, from, to, amount, nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, starknet.MAINNET_ID, 0, pri)
	assert.NoError(t, err)

	res := tx.GetTxRequestJson()
	assert.NotEmpty(t, res)

	hash, err := GetTxHashV3(res, starknet.MAINNET_ID)
	assert.NoError(t, err)

	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x3a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034","calldata":["0x1","0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d","0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e","0x0","0x3","0x3","0x3a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034","0x64","0x0"],"version":"0x3","signature":["0x53ef1aca1c3dda54e1a24543afab15b75086939cc111b23bf85ed64d0ae9b85","0xe7b98b07791dd98379937e200c382ef07592a2277ea2f553ac0e0d343b8a31"],"nonce":"0x1","resource_bounds":{"L1_GAS":{"max_amount":"0x0","max_price_per_unit":"0x2129fdfa02ec"},"L1_DATA_GAS":{"max_amount":"0x100","max_price_per_unit":"0x14e42b21538"},"L2_GAS":{"max_amount":"0x1bda00","max_price_per_unit":"0x365600c8"}},"tip":"0x0","paymaster_data":[],"account_deployment_data":[],"nonce_data_availability_mode":0,"fee_data_availability_mode":0}`, res)
	assert.Equal(t, "0x636c63aff2353b3029f2146b9798dc76ee70d4b984030864da19d53f36ffd17", starknet.BigToHex(tx.GetTxHash()))
	assert.Equal(t, "0x636c63aff2353b3029f2146b9798dc76ee70d4b984030864da19d53f36ffd17", hash)
}

func TestCreateMultiContractV3Cairo0(t *testing.T) {
	curve := starknet.SC()
	from := "0x3a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034"
	txs := []starknet.Calls{
		{
			ContractAddress: starknet.STARK,
			Entrypoint:      "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
			Calldata:        []string{"0x3a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034", "0x64", "0x0"},
		}, {
			ContractAddress: starknet.STARK,
			Entrypoint:      "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e",
			Calldata:        []string{"0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d", "0x164", "0x0"},
		},
	}
	nonce := big.NewInt(2)

	tip := doubleHex("0x0")
	l1GasMaxAmount := doubleHex("0x0")
	l1GasMaxPricePerUnit := doubleHex("0x119e72c54321")

	l2GasMaxAmount := doubleHex("0x134ac0")
	l2GasMaxPricePerUnit := doubleHex("0x1cddeb25")

	l1DataGasMaxAmount := doubleHex("0x80")
	l1DataMaxPricePerUnit := doubleHex("0xa9fb76ef14")

	tx, err := CreateSignedMultiContractTxV3(curve, from, txs, nonce, tip, l1GasMaxAmount, l1GasMaxPricePerUnit, l2GasMaxAmount, l2GasMaxPricePerUnit, l1DataGasMaxAmount, l1DataMaxPricePerUnit, starknet.MAINNET_ID, 0, pri)
	assert.NoError(t, err)

	res := tx.GetTxRequestJson()
	assert.NotEmpty(t, res)

	hash, err := GetTxHashV3(res, starknet.MAINNET_ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	assert.Equal(t, "{\"type\":\"INVOKE_FUNCTION\",\"sender_address\":\"0x3a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034\",\"calldata\":[\"0x2\",\"0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d\",\"0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e\",\"0x0\",\"0x3\",\"0x4718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d\",\"0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e\",\"0x3\",\"0x3\",\"0x6\",\"0x3a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034\",\"0x64\",\"0x0\",\"0x109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d\",\"0x164\",\"0x0\"],\"version\":\"0x3\",\"signature\":[\"0x6cf6cb937307cb1724c6906420111153b82e0380f1a1df22a7277d9e6a6a3a0\",\"0x1acbc77701583b9090f969e33da128c465242ddb87398120a783a3843654cfc\"],\"nonce\":\"0x2\",\"resource_bounds\":{\"L1_GAS\":{\"max_amount\":\"0x0\",\"max_price_per_unit\":\"0x233ce58a8642\"},\"L1_DATA_GAS\":{\"max_amount\":\"0x100\",\"max_price_per_unit\":\"0x153f6edde28\"},\"L2_GAS\":{\"max_amount\":\"0x269580\",\"max_price_per_unit\":\"0x39bbd64a\"}},\"tip\":\"0x0\",\"paymaster_data\":[],\"account_deployment_data\":[],\"nonce_data_availability_mode\":0,\"fee_data_availability_mode\":0}", res)
	assert.Equal(t, "0x74d46e5213e790b078754e489a7d2ea6858e0668c80ce6b208986791396dbe4", starknet.BigToHex(tx.GetTxHash()))
	assert.Equal(t, "0x74d46e5213e790b078754e489a7d2ea6858e0668c80ce6b208986791396dbe4", hash)
}

func doubleHex(hexStr string) string {
	n := new(big.Int)
	n.SetString(hexStr[2:], 16)
	n.Mul(n, big.NewInt(2))
	return fmt.Sprintf("0x%x", n)
}
