# eos-sdk
EOS SDK is used to interact with the EOS blockchain, it contains various functions that can be used for web3 wallet.
The SDK not only supports EOS, it also supports other blockchains forked from EOS such as WAX.

- EOS
- REX
- TNT
- WAX

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/eos
```

## Usage

### Supported Functions

```golang
// NewAccount
// NewTransaction & SignTransaction
// NewContractTransaction & SignTransaction
```

### New Account
```golang
	userName := "test3"
	if len(userName) > 12 {
		return
	}
	gotPrivKey, _ := GenerateKeyPair()
	p, err := ecc.NewPrivateKey(gotPrivKey)
	if err != nil {
        // todo
		fmt.Println(err)
	}
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
        if err != nil {
            // todo
		    fmt.Println(err)
	    }
	}
```

### New Transaction and Sign
```golang
	privateKey, err := ecc.NewPrivateKey("5JvW9FSHci6MQcnoHjNnfv5T4Pfi5pj2weAEFQvq1TFaxs8Kbnt")
	opt, err = getTxOptions()
	if err != nil {
        // todo
	}
	tx := NewTransactionWithParams("dubuqing1111", "dubuqing1234", "test", types.NewWAXAsset(500000000), opt)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	DumpPackedTx(t, packedTx)
```
### New Contract Transaction and Sign
```golang
	privateKey, _ := ecc.NewPrivateKey("5JvW9FSHci6MQcnoHjNnfv5T4Pfi5pj2weAEFQvq1TFaxs8Kbnt")
	opt, err := getTxOptions()
	if err != nil {
        // todo
	}
	contractName := "wax.token"
	tx := NewContractTransaction(contractName, "dubuqing1111", "dubuqing1234", "test", types.NewWAXAsset(500000000), opt)
	DumpGetRequiredKeyContent(tx, privateKey.PublicKey().String(), t)
	signedTx, packedTx, err := SignTransaction(p1, tx, opt.ChainID, types.CompressionNone)
	if err != nil {
        // todo
	}	
	DumpPackedTx(t, packedTx)
```

## Credits  This project includes code adapted from the following sources:  
- [eos-go](https://github.com/eoscanada/eos-go) - EOS Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/eos/LICENSE>) licensed, see package or folder for the respective license.
