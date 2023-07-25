package terra

import (
	"encoding/hex"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types"
	"github.com/okx/go-wallet-sdk/crypto/bip32"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"testing"
)

//https://bombay-lcd.terra.dev/cosmos/tx/v1beta1/txs/F39950309C6DF0A81305503B81E03BD3947EA4012957DD865988B67815C21D20
//
//https://bombay-fcd.terra.dev/v1/tx/AA2413DCBCED6C02560F7E411B980BF2A3617B7C271EE28B061EDB7F7B5918A4
//
//
//## Check account details
//https://bombay-lcd.terra.dev/cosmos/auth/v1beta1/accounts/terra1xmkczk59xgjhzgwhfg8l5tgs2uftpuj9cgazr4
//
//
//curl -X POST -d '{"tx_bytes":"CrUHCqwHCiYvdGVycmEud2FzbS52MWJldGExLk1zZ0V4ZWN1dGVDb250cmFjdBKBBwosdGVycmExeG1rY3prNTl4Z2poemd3aGZnOGw1dGdzMnVmdHB1ajljZ2F6cjQSLHRlcnJhMTR6ODByd3BkMGFsemo0eGR0Z3FkbWNxdDl3ZDl4ajVmZmQ2MHdwGpAGewogICJleGVjdXRlX3N3YXBfb3BlcmF0aW9ucyI6IHsKICAgICJtaW5pbXVtX3JlY2VpdmUiOiAiMzQxMTc0NjgiLAogICAgIm9mZmVyX2Ftb3VudCI6ICIxMDAwMDAwIiwKICAgICJvcGVyYXRpb25zIjogWwogICAgICB7CiAgICAgICAgInRlcnJhX3N3YXAiOiB7CiAgICAgICAgICAiYXNrX2Fzc2V0X2luZm8iOiB7CiAgICAgICAgICAgICJ0b2tlbiI6IHsKICAgICAgICAgICAgICAiY29udHJhY3RfYWRkciI6ICJ0ZXJyYTF1MHQzNWRyenl5MG11amo4cmtkeXpoZTI2NHVsczR1ZzN3ZHAzeCIKICAgICAgICAgICAgfQogICAgICAgICAgfSwKICAgICAgICAgICJvZmZlcl9hc3NldF9pbmZvIjogewogICAgICAgICAgICAibmF0aXZlX3Rva2VuIjogewogICAgICAgICAgICAgICJkZW5vbSI6ICJ1bHVuYSIKICAgICAgICAgICAgfQogICAgICAgICAgfQogICAgICAgIH0KICAgICAgfSwKICAgICAgewogICAgICAgICJ0ZXJyYV9zd2FwIjogewogICAgICAgICAgImFza19hc3NldF9pbmZvIjogewogICAgICAgICAgICAibmF0aXZlX3Rva2VuIjogewogICAgICAgICAgICAgICJkZW5vbSI6ICJ1dXNkIgogICAgICAgICAgICB9CiAgICAgICAgICB9LAogICAgICAgICAgIm9mZmVyX2Fzc2V0X2luZm8iOiB7CiAgICAgICAgICAgICJ0b2tlbiI6IHsKICAgICAgICAgICAgICAiY29udHJhY3RfYWRkciI6ICJ0ZXJyYTF1MHQzNWRyenl5MG11amo4cmtkeXpoZTI2NHVsczR1ZzN3ZHAzeCIKICAgICAgICAgICAgfQogICAgICAgICAgfQogICAgICAgIH0KICAgICAgfQogICAgXQogIH0KfSoQCgV1bHVuYRIHMTAwMDAwMBIEdGVzdBJpClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECc7guyn8yqBe5cTB5L5LA9V+uwv+KYQEpb8ScaxbGGmESBAoCCAEYCRIVCg8KBXVsdW5hEgYxMDAwMDAQwIQ9GkADbF62HMSF+RMX/1cAVLExsk/XDvDyjSCNAeWshy5ZsGIUXF/CfaynfZR60Kj3d9n/Ufby/4esG6I9ypXHCNsd","mode":"BROADCAST_MODE_SYNC"}' https://bombay-lcd.terra.dev/cosmos/tx/v1beta1/txs
//
//curl -X POST -d '{"tx_bytes":"CrQECrEECiYvdGVycmEud2FzbS52MWJldGExLk1zZ0V4ZWN1dGVDb250cmFjdBKGBAosdGVycmExZnJjZjM2anZxdmo0N2NyOWR5Z2Z4M3ZoNnB1cWY5ZWxtbm5zbDISLHRlcnJhMTR6ODByd3BkMGFsemo0eGR0Z3FkbWNxdDl3ZDl4ajVmZmQ2MHdwGqcDeyJleGVjdXRlX3N3YXBfb3BlcmF0aW9ucyI6eyJtaW5pbXVtX3JlY2VpdmUiOiI1OTk3MTA2NyIsIm9mZmVyX2Ftb3VudCI6IjEwMDAwMDAiLCJvcGVyYXRpb25zIjpbeyJ0ZXJyYV9zd2FwIjp7ImFza19hc3NldF9pbmZvIjp7InRva2VuIjp7ImNvbnRyYWN0X2FkZHIiOiJ0ZXJyYTF1MHQzNWRyenl5MG11amo4cmtkeXpoZTI2NHVsczR1ZzN3ZHAzeCJ9fSwib2ZmZXJfYXNzZXRfaW5mbyI6eyJuYXRpdmVfdG9rZW4iOnsiZGVub20iOiJ1bHVuYSJ9fX19LHsidGVycmFfc3dhcCI6eyJhc2tfYXNzZXRfaW5mbyI6eyJuYXRpdmVfdG9rZW4iOnsiZGVub20iOiJ1dXNkIn19LCJvZmZlcl9hc3NldF9pbmZvIjp7InRva2VuIjp7ImNvbnRyYWN0X2FkZHIiOiJ0ZXJyYTF1MHQzNWRyenl5MG11amo4cmtkeXpoZTI2NHVsczR1ZzN3ZHAzeCJ9fX19XX19ElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQKwapSLujUw96LqDzoNC8IPvknM/CMOrCfQndwqW0esaxIECgIIfxgGEgQQwIQ9GkC1B1CZIxMo9npRaMugwWQysf2HuPkZZw3jy5dhqcSiRmmWn9vqEHmOYltmOLTdkoHt+OdNQAUvM4+P/Xm5/VYL","mode":"BROADCAST_MODE_SYNC"}' https://bombay-lcd.terra.dev/cosmos/tx/v1beta1/txs

func TestLuna1(t *testing.T) {
	mnemonic := "arena special tunnel keen skate chapter media scare injury indoor topic aware autumn lecture depth lava legal raccoon clog pulp renew diagram upper blade"
	if !bip39.IsMnemonicValid(mnemonic) {
		t.Fatal("mnemonic error")
		return
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		t.Fatal("seed error")
		return
	}

	masterKey, _ := bip32.NewMasterKey(seed)
	key, _ := masterKey.NewChildKeyByChainId(330)

	privateKeyHex := hex.EncodeToString(key.Key.Key)
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	if address != "terra1zu3nyxs7setu9l69ehdf0e9g5yc8uhw4xxjksx" {
		t.Fatal("NewAddress failed")
	}

	ret := ValidateAddress(address)
	if !ret {
		t.Fatal("ValidateAddress failed")
	}
}

func TestLuna2(t *testing.T) {
	mnemonic := "arena special tunnel keen skate chapter media scare injury indoor topic aware autumn lecture depth lava legal raccoon clog pulp renew diagram upper blade"
	if !bip39.IsMnemonicValid(mnemonic) {
		t.Fatal("mnemonic error")
		return
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		t.Fatal("seed error")
		return
	}

	masterKey, _ := bip32.NewMasterKey(seed)
	key, _ := masterKey.NewChildKeyByChainId(330)

	privateKeyHex := hex.EncodeToString(key.Key.Key)

	input := TransactionInput{}
	input.ChainId = "bombay-12"
	input.Memo = "test"
	input.Sequence = 4
	input.AccountNumber = 588053
	input.GasLimit = 100000
	input.AppendFeeCoin("uluna", big.NewInt(2000))

	sendCoins := types.NewCoins(types.NewCoin("uluna", types.NewIntFromUint64(10000)))
	input.AppendSendMsg("terra1xmkczk59xgjhzgwhfg8l5tgs2uftpuj9cgazr4", "terra1vm9pfph4syf9g3hfz29636cfw5wp9n6xwut8xu", &sendCoins)

	rawTx := NewTransaction(input, privateKeyHex)
	signHex := Sign(rawTx, privateKeyHex)
	fmt.Println(signHex)
	assert.Equal(t, signHex, "f3bf2b401b6d6048c9ffd17e6cfe08b5115c15ca5b8e938210c343161e9adb303f0b7709f0222f91965fb6f55b9f8c065f4062277e2b18464fc3fa87b4a67b0f")
	result := SignEnd(rawTx, signHex)
	assert.Equal(t, result, "CpUBCowBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmwKLHRlcnJhMXhta2N6azU5eGdqaHpnd2hmZzhsNXRnczJ1ZnRwdWo5Y2dhenI0Eix0ZXJyYTF2bTlwZnBoNHN5ZjlnM2hmejI5NjM2Y2Z3NXdwOW42eHd1dDh4dRoOCgV1bHVuYRIFMTAwMDASBHRlc3QSZwpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAufRw/M+ifgAtdVlKRnENXGP04S7avt8IIkbqRory22CEgQKAggBGAQSEwoNCgV1bHVuYRIEMjAwMBCgjQYaQPO/K0AbbWBIyf/Rfmz+CLURXBXKW46TghDDQxYemtswPwt3CfAiL5GWX7b1W5+MBl9AYid+KxhGT8P6h7Smew8=")
}

func TestLuna3(t *testing.T) {
	mnemonic := "arena special tunnel keen skate chapter media scare injury indoor topic aware autumn lecture depth lava legal raccoon clog pulp renew diagram upper blade"
	if !bip39.IsMnemonicValid(mnemonic) {
		t.Fatal("mnemonic error")
		return
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		t.Fatal("seed error")
		return
	}

	masterKey, _ := bip32.NewMasterKey(seed)
	key, _ := masterKey.NewChildKeyByChainId(330)

	privateKeyHex := hex.EncodeToString(key.Key.Key)

	input := TransactionInput{}
	input.ChainId = "bombay-12"
	input.Memo = "test"
	input.Sequence = 9
	input.AccountNumber = 588053
	input.GasLimit = 1000000
	input.AppendFeeCoin("uluna", big.NewInt(100000))

	//input.AppendSwapMsg("terra1xmkczk59xgjhzgwhfg8l5tgs2uftpuj9cgazr4", "uusd", "uluna", big.NewInt(1000000))
	contractCoins := types.NewCoins(types.NewCoin("uluna", types.NewIntFromUint64(1000000)))
	input.AppendContractMsg("terra1xmkczk59xgjhzgwhfg8l5tgs2uftpuj9cgazr4", "terra14z80rwpd0alzj4xdtgqdmcqt9wd9xj5ffd60wp", "{\n  \"execute_swap_operations\": {\n    \"minimum_receive\": \"34117468\",\n    \"offer_amount\": \"1000000\",\n    \"operations\": [\n      {\n        \"terra_swap\": {\n          \"ask_asset_info\": {\n            \"token\": {\n              \"contract_addr\": \"terra1u0t35drzyy0mujj8rkdyzhe264uls4ug3wdp3x\"\n            }\n          },\n          \"offer_asset_info\": {\n            \"native_token\": {\n              \"denom\": \"uluna\"\n            }\n          }\n        }\n      },\n      {\n        \"terra_swap\": {\n          \"ask_asset_info\": {\n            \"native_token\": {\n              \"denom\": \"uusd\"\n            }\n          },\n          \"offer_asset_info\": {\n            \"token\": {\n              \"contract_addr\": \"terra1u0t35drzyy0mujj8rkdyzhe264uls4ug3wdp3x\"\n            }\n          }\n        }\n      }\n    ]\n  }\n}", &contractCoins)

	rawTx := NewTransaction(input, privateKeyHex)
	signHex := Sign(rawTx, privateKeyHex)
	result := SignEnd(rawTx, signHex)
	assert.Equal(t, result, "CrUHCqwHCiYvdGVycmEud2FzbS52MWJldGExLk1zZ0V4ZWN1dGVDb250cmFjdBKBBwosdGVycmExeG1rY3prNTl4Z2poemd3aGZnOGw1dGdzMnVmdHB1ajljZ2F6cjQSLHRlcnJhMTR6ODByd3BkMGFsemo0eGR0Z3FkbWNxdDl3ZDl4ajVmZmQ2MHdwGpAGewogICJleGVjdXRlX3N3YXBfb3BlcmF0aW9ucyI6IHsKICAgICJtaW5pbXVtX3JlY2VpdmUiOiAiMzQxMTc0NjgiLAogICAgIm9mZmVyX2Ftb3VudCI6ICIxMDAwMDAwIiwKICAgICJvcGVyYXRpb25zIjogWwogICAgICB7CiAgICAgICAgInRlcnJhX3N3YXAiOiB7CiAgICAgICAgICAiYXNrX2Fzc2V0X2luZm8iOiB7CiAgICAgICAgICAgICJ0b2tlbiI6IHsKICAgICAgICAgICAgICAiY29udHJhY3RfYWRkciI6ICJ0ZXJyYTF1MHQzNWRyenl5MG11amo4cmtkeXpoZTI2NHVsczR1ZzN3ZHAzeCIKICAgICAgICAgICAgfQogICAgICAgICAgfSwKICAgICAgICAgICJvZmZlcl9hc3NldF9pbmZvIjogewogICAgICAgICAgICAibmF0aXZlX3Rva2VuIjogewogICAgICAgICAgICAgICJkZW5vbSI6ICJ1bHVuYSIKICAgICAgICAgICAgfQogICAgICAgICAgfQogICAgICAgIH0KICAgICAgfSwKICAgICAgewogICAgICAgICJ0ZXJyYV9zd2FwIjogewogICAgICAgICAgImFza19hc3NldF9pbmZvIjogewogICAgICAgICAgICAibmF0aXZlX3Rva2VuIjogewogICAgICAgICAgICAgICJkZW5vbSI6ICJ1dXNkIgogICAgICAgICAgICB9CiAgICAgICAgICB9LAogICAgICAgICAgIm9mZmVyX2Fzc2V0X2luZm8iOiB7CiAgICAgICAgICAgICJ0b2tlbiI6IHsKICAgICAgICAgICAgICAiY29udHJhY3RfYWRkciI6ICJ0ZXJyYTF1MHQzNWRyenl5MG11amo4cmtkeXpoZTI2NHVsczR1ZzN3ZHAzeCIKICAgICAgICAgICAgfQogICAgICAgICAgfQogICAgICAgIH0KICAgICAgfQogICAgXQogIH0KfSoQCgV1bHVuYRIHMTAwMDAwMBIEdGVzdBJpClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEC59HD8z6J+AC11WUpGcQ1cY/ThLtq+3wgiRupGivLbYISBAoCCAEYCRIVCg8KBXVsdW5hEgYxMDAwMDAQwIQ9GkBfsMTIjEEKt2AMH+AP3FeYmHXJwoG3CujwQB061NU/wW8phJe/msODAvtk+R49K75oVUS6E5lZr+LXejbxo9Z2")
}
