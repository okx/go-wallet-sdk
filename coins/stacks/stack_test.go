package stacks

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

// Deprecated
func TestStack(t *testing.T) {
	senderKey := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	stack, err := Stack("SP000000000000000000002Q6VF78", "pox", "stack-stx", "36Y1UJBWGGreKCKNYQPVPr41rgG2sQF7SC",
		big.NewInt(420041303000), big.NewInt(668000), big.NewInt(6), senderKey, big.NewInt(1), big.NewInt(5000))
	require.NoError(t, err)
	require.Equal(t, "000000000104006ecfff9cee8ac5367c83ad0819e4c500b6c475d60000000000000001000000000000138800006bb6e58aa6befcc08c74f5fa3dc5f0b78eb3b133bf7042846f1e05d45f86492a48b9a814b4ed5e82bdaddaefe7d2fe28cc6ce0b128c3e05a9f94fa91eed4cd760302000000000216000000000000000000000000000000000000000003706f7809737461636b2d7374780000000401000000000000000000000061cc69a3d80c00000002096861736862797465730200000014352481ec2fecfde0c5cdc635a383c4ac27b9f71e0776657273696f6e02000000010101000000000000000000000000000a31600100000000000000000000000000000006", stack)
}

func TestContractCall(t *testing.T) {
	key := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	from := "SP1QCZZWWXT5CADKWGEPGG6F4RM0BDH3NTTNM86ZG"
	to := "SP3HXJJMJQ06GNAZ8XWDN1QM48JEDC6PP6W3YZPZJ"
	memo := "110317"
	contract := "SP466FNC0P7JWTNM2R9T199QRZN1MYEDTAR0KP27"
	contractName := "miamicoin-token"
	functionName := "transfer"
	tokenName := "miamicoin"

	amount := big.NewInt(21)
	nonce := big.NewInt(21)

	tx, err := ContractCall(key, from, to, memo, amount, contract, contractName, tokenName, functionName, nonce, big.NewInt(200000))
	require.NoError(t, err)
	var s TransactionRes
	err = json.Unmarshal([]byte(tx), &s)
	require.NoError(t, err)
	require.Equal(t, "000000000104006ecfff9cee8ac5367c83ad0819e4c500b6c475d600000000000000150000000000030d40000010d823fd03339e8678e14e2f421735dfb4c6758452d5e96d7ec92c510f91892d3a0c55d435848fa7de6710ee15c72ea9daca460863617af389221dc504743e560302000000010102166ecfff9cee8ac5367c83ad0819e4c500b6c475d61608633eac058f2e6ab41613a0a537c7ea1a79cdd20f6d69616d69636f696e2d746f6b656e096d69616d69636f696e010000000000000015021608633eac058f2e6ab41613a0a537c7ea1a79cdd20f6d69616d69636f696e2d746f6b656e087472616e7366657200000004010000000000000000000000000000001505166ecfff9cee8ac5367c83ad0819e4c500b6c475d60516e3d94a92b80d0aabe8ef1b50de84449cd61ad6370a0200000006313130333137", s.TxSerialize)
}

func TestAllowContractCaller(t *testing.T) {
	contractAddress := "SP000000000000000000002Q6VF78"
	contractName := "pox-3"
	functionName := "allow-contract-caller"
	caller := "SP21YTSM60CAY6D011EZVEVNKXVW8FVZE198XEFFP.pox-fast-pool-v2"
	untilBurnBlockHeight := big.NewInt(206600)
	privateKey := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	fee := big.NewInt(3000)
	nonce := big.NewInt(57)
	addr, err := NewContractPrincipalCV(caller)
	require.NoError(t, err)
	functionArgs := []ClarityValue{
		addr,
		// If untilBurnBlockHeight is empty, pass nil
		NewUntilBurnBlockHeight(untilBurnBlockHeight),
	}
	functionArgsJson := `
{
  "functionArgs": [
    {
      "type": 6,
      "address": {
        "type": 0,
        "version": 22,
        "hash160": "83ed66860315e334010bbfb76eb3eef887efee0a"
      },
      "contractName": {
        "content": "pox-fast-pool-v2",
        "lengthPrefixBytes": 1,
        "maxLengthBytes": 128,
        "type": 2
      }
    },
    {
      "type": 10,
      "value": {
        "type": 1,
        "value": 206600
      }
    }
  ]
}
`
	args := getFunctionArgs(functionArgsJson)
	functionArgs2 := DeserializeJson(args)
	txOption := &SignedContractCallOptions{
		ContractAddress: contractAddress,
		ContractName:    contractName,
		FunctionName:    functionName,
		FunctionArgs:    functionArgs2,
		SendKey:         privateKey,
		ValidateWithAbi: false,
		Fee:             *fee,
		Nonce:           *nonce,
		AnchorMode:      3,
	}

	tx, err := MakeContractCall(txOption)
	require.NoError(t, err)
	txSerialized := hex.EncodeToString(Serialize(*tx))
	txId := Txid(*tx)
	assert.Equal(t, true, reflect.DeepEqual(functionArgs, functionArgs2))
	assert.Equal(t, "000000000104006ecfff9cee8ac5367c83ad0819e4c500b6c475d600000000000000390000000000000bb800010994f51eef9b6b806f441fb935db71d69bb444cd7894b9811bab0f94337324c017ebf34bb0c39967b90580bd12419d87ceb22eb9fdb5710e28e0bc4eca87e1ee0302000000000216000000000000000000000000000000000000000005706f782d3315616c6c6f772d636f6e74726163742d63616c6c657200000002061683ed66860315e334010bbfb76eb3eef887efee0a10706f782d666173742d706f6f6c2d76320a0100000000000000000000000000032708", txSerialized)
	assert.Equal(t, "5ece84017f0c67dd8ee2d4136bb032a9e49bb902efe85bb47b8daea6a79e4d8c", txId)
}

func TestDelegateStx(t *testing.T) {
	contractAddress := "SP000000000000000000002Q6VF78"
	contractName := "pox-3"
	functionName := "delegate-stx"
	delegateTo := "SP3TDKYYRTYFE32N19484838WEJ25GX40Z24GECPZ"
	poxAddress := "36Y1UJBWGGreKCKNYQPVPr41rgG2sQF7SC"
	amountMicroStx := big.NewInt(100000000000)
	untilBurnBlockHeight := big.NewInt(2000)
	privateKey := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	fee := big.NewInt(3000)
	nonce := big.NewInt(58)

	poxAddressCV, err := GetPoxAddress(poxAddress)
	require.NoError(t, err)
	functionArgs := []ClarityValue{
		NewUintCV(amountMicroStx),
		*NewStandardPrincipalCV(delegateTo),
		&SomeCV{OptionalSome, NewUintCV(untilBurnBlockHeight)},
		&SomeCV{OptionalSome, poxAddressCV},
	}

	functionArgsJson := `
{
  "functionArgs": [
    {
      "type": 1,
      "value": 100000000000
    },
    {
      "type": 5,
      "address": {
        "type": 0,
        "version": 22,
        "hash160": "f4d9fbd8d79ee18aa14910440d1c7484587480f8"
      }
    },
    {
      "type": 10,
      "value": {
        "type": 1,
        "value": 2000
      }
    },
    {
      "type": 10,
      "value": {
        "type": 12,
        "data": {
          "hashbytes": {
            "buffer": "NSSB7C/s/eDFzcY1o4PErCe59x4=",
            "type": 2
          },
          "version": {
            "buffer": "AQ==",
            "type": 2
          }
        }
      }
    }
  ]
}`
	args := getFunctionArgs(functionArgsJson)
	functionArgs2 := DeserializeJson(args)
	assert.Equal(t, true, reflect.DeepEqual(functionArgs, functionArgs2))
	txOption := &SignedContractCallOptions{
		ContractAddress: contractAddress,
		ContractName:    contractName,
		FunctionName:    functionName,
		FunctionArgs:    functionArgs,
		SendKey:         privateKey,
		ValidateWithAbi: false,
		Fee:             *fee,
		Nonce:           *nonce,
		AnchorMode:      3,
	}

	tx, err := MakeContractCall(txOption)
	require.NoError(t, err)
	txSerialized := hex.EncodeToString(Serialize(*tx))
	txId := Txid(*tx)
	assert.Equal(t, "000000000104006ecfff9cee8ac5367c83ad0819e4c500b6c475d6000000000000003a0000000000000bb80001d6036a78231aecf9ab1dacf3556bda0666ece5ff858319b83dc20dddd75df64e30c26bac6a0dc6f78e96777ce64c0c5174d911d0aec2a8c82e386f0c845c74310302000000000216000000000000000000000000000000000000000005706f782d330c64656c65676174652d73747800000004010000000000000000000000174876e8000516f4d9fbd8d79ee18aa14910440d1c7484587480f80a01000000000000000000000000000007d00a0c00000002096861736862797465730200000014352481ec2fecfde0c5cdc635a383c4ac27b9f71e0776657273696f6e020000000101", txSerialized)
	assert.Equal(t, "173fae3b3effd46f803a08e27371980373fcb31fe44c1437686e2e95b61b72ad", txId)
}

func TestRevokeDelegateStx(t *testing.T) {
	contractAddress := "SP000000000000000000002Q6VF78"
	contractName := "pox-3"
	functionName := "revoke-delegate-stx"
	privateKey := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	fee := big.NewInt(3000)
	nonce := big.NewInt(59)
	var functionArgs []ClarityValue

	txOption := &SignedContractCallOptions{
		ContractAddress: contractAddress,
		ContractName:    contractName,
		FunctionName:    functionName,
		FunctionArgs:    functionArgs,
		SendKey:         privateKey,
		ValidateWithAbi: false,
		Fee:             *fee,
		Nonce:           *nonce,
		AnchorMode:      3,
	}

	tx, err := MakeContractCall(txOption)
	require.NoError(t, err)
	txSerialized := hex.EncodeToString(Serialize(*tx))
	txId := Txid(*tx)
	require.Equal(t, "000000000104006ecfff9cee8ac5367c83ad0819e4c500b6c475d6000000000000003b0000000000000bb80000836a3bb9ea65e4a725e2f09b73df8a2361542452e6b4bf033639668ff395f07f1bfa0a59b5737a9b26681cb2ab53982afe1ac6061e9b4e806f293314ab934b9a0302000000000216000000000000000000000000000000000000000005706f782d33137265766f6b652d64656c65676174652d73747800000000", txSerialized)
	require.Equal(t, "ae33e6fb41975052b5ca59443cd393f667bbea1379f5c2fc7a65e7bfd6a30643", txId)
}

func TestContractCallWithPostConditions(t *testing.T) {
	contractAddress := "SP3K8BC0PPEVCV7NZ6QSRWPQ2JE9E5B6N3PA0KBR9"
	contractName := "amm-swap-pool-v1-1"
	functionName := "swap-helper"
	privateKey := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	postConditionMode := 2
	postConditions := []string{"000216c03b5520cf3a0bd270d8e41e5e19a464aef6294c010000000000002710", "010316e685b016b3b6cd9ebf35f38e5ae29392e2acd51d0f616c65782d7661756c742d76312d3116e685b016b3b6cd9ebf35f38e5ae29392e2acd51d176167653030302d676f7665726e616e63652d746f6b656e04616c657803000000000078b854"}
	fee := big.NewInt(0)
	nonce := big.NewInt(0)
	var functionArgs []ClarityValue

	txOption := &SignedContractCallOptions{
		ContractAddress:         contractAddress,
		ContractName:            contractName,
		FunctionName:            functionName,
		FunctionArgs:            functionArgs,
		SendKey:                 privateKey,
		ValidateWithAbi:         false,
		Fee:                     *fee,
		Nonce:                   *nonce,
		AnchorMode:              3,
		PostConditionMode:       postConditionMode,
		SerializePostConditions: postConditions,
	}

	tx, _ := MakeContractCall(txOption)
	txSerlize := hex.EncodeToString(Serialize(*tx))
	txId := Txid(*tx)
	fmt.Println(txSerlize)
	fmt.Println(txId)
}
