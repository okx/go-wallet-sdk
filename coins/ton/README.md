# ton-sdk

Ton SDK is used to interact with the Ton blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/ton
```

## Usage

### New Address (v4r2 wallet address)

```go
    seedHex := "45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40"
seed, _ := hex.DecodeString(seedHex)
address, err := NewAddress(seed)
fmt.Println(address)
```

### Transfer TON

```go
seed, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40") //your private key
address, err := NewAddress(seed, wallet.V4R2) //generate address
fmt.Println(address)
to := "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR" //to ,which address receives the TON.
amount := "100000000" //amount to transfer . 1 TON means 100000000
comment := ""         //your comment .optional.
seqno := uint32(0) //your seqno.
pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
expireAt := time.Now().Unix() + 600
//false means that now it is  Not a simulative transaction.
simulate := false
signedTx, err := Transfer(seed, pubKey, to, amount, comment, seqno, expireAt, 3, simulate, wallet.V4R2) //3 is recommended pattern. false means that now it is  Not a simulative transaction.
assert.Nil(t, err)
t.Log(signedTx.Tx)
fmt.Println(signedTx.Tx) // the tx is in signedTx.Tx
fmt.Println(signedTx.Hash) //this hash is used to query the transaction in the browser.
```

### Transfer Jetton token

```go
seed, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40") //your private key
fromJettonAccount := "UQD7w9qG8Cq0PgX0hnp5nVpHPeDYL0QlhcLtjFMmna43sMxz" //your jetton account
to := "UQC27fdnAFQcQDaXDrR89OKx-lW_Zyxuzcy5CjfPrS9A6vZf"                //to ,which address receives the TON.
amount := "1"      //jetton amount to transfer . Each token has a different precision.1 USDT may means 1000000
seqno := uint32(0) //your seqno.
pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
//"jetton test",, your comment .optional.
comment :="jetton test"
//recommended value
messageAttachedTons := "50000000"

expireAt := time.Now().Unix() + 600
//recommended value
invokeNotificationFee := "1"
//false means that now it is  Not a simulative transaction.
simulate := false
customPayload := ""
stateInit := ""
signedTx, err := TransferJetton(seed, pubKey, fromJettonAccount, to, amount, 9, seqno, messageAttachedTons, invokeNotificationFee, customPayload, stateInit, comment, expireAt, 0, simulate, wallet.V4R2)

assert.Nil(t, err)
fmt.Println( signedTx.Tx) // the tx is in signedTx.Tx
fmt.Println( signedTx.Hash) //this hash is used to query the transaction in the browser.
```

### Simulate Transferring TON(Used to estimate gas charges)

```go
    seed, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40") //your private key
address, err := NewAddress(seed) //generate address
fmt.Println(address)
to := "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR" //to ,which address receives the TON.
amount := "100000000" //amount to transfer . 1 TON means 100000000
comment := ""         //your comment .optional.
seqno := uint32(0) //your seqno.
expireAt := time.Now().Unix() + 600
//true means that now it IS  a simulative transaction.
simulate := true
pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
signedTx, err := Transfer(nil, pubKey, to, amount, comment, seqno, expireAt, 3, simulate, wallet.V4R2) //3 is recommended pattern. 
assert.Nil(t, err)
fmt.Println( signedTx.Tx) // the tx is in signedTx.Tx Used to estimate gas charges
```

### Simulate  Transferring Jetton token(Used to estimate gas charges)

```go
seed, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40") //your private key
fromJettonAccount := "UQD7w9qG8Cq0PgX0hnp5nVpHPeDYL0QlhcLtjFMmna43sMxz" //your jetton account
to := "UQC27fdnAFQcQDaXDrR89OKx-lW_Zyxuzcy5CjfPrS9A6vZf"                //to ,which address receives the TON.
amount := "1"      //jetton amount to transfer . Each token has a different precision.1 USDT may means 1000000
seqno := uint32(0) //your seqno.
pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
//"jetton test",, your comment .optional.
comment :="jetton test"
//recommended value
messageAttachedTons := "50000000"
expireAt := time.Now().Unix() + 600
//recommended value
invokeNotificationFee := "1"
//true means that now it IS  a simulative transaction.
simulate := true
signedTx, err := TransferJetton(seed, pubKey, fromJettonAccount, to, amount, 9, seqno, messageAttachedTons, invokeNotificationFee, customPayload, stateInit, comment, expireAt, 0, simulate, wallet.V4R2)
assert.Nil(t, err)
fmt.Println( signedTx.Tx) // the tx is in signedTx.Tx Used to estimate gas charges
```

### Sign proof

```go
seed, err := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
assert.NoError(t, err)
addr := "EQA3_JIJKDC0qauDUEQe2KjQj1iLwQRtrEREzmfDxbCKw9Kr"
proof := &ProofData{
Timestamp: 1719309177,    // timestamp 
Domain:    "ton.org.com", //domain
Payload:   "123", //comment or something
}
r, err := SignProof(addr, seed, proof)
assert.NoError(t, err)
pub := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
sign, err := base64.StdEncoding.DecodeString(r) //sign of proof 
assert.NoError(t, err)
expect := "V1ImmDgpt4DtZYYeGeZz38w7J+dXtYbBf/Hl7DLcWLEad23TOexKCSTO1f+N7i3UDreGVfycaVNbOspJnr9aDw=="
assert.Equal(t, r, expect)
assert.NoError(t, VerifySignProof(addr, pub, sign, proof))
```

### Verify proof

```go
seed, err := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
assert.NoError(t, err)
addr := "EQA3_JIJKDC0qauDUEQe2KjQj1iLwQRtrEREzmfDxbCKw9Kr"
proof := &ProofData{
Timestamp: 1719309177,    // timestamp 
Domain:    "ton.org.com", //domain
Payload:   "123", //comment or something
}
r, err := SignProof(addr, seed, proof)
assert.NoError(t, err)
pub := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
expect := "V1ImmDgpt4DtZYYeGeZz38w7J+dXtYbBf/Hl7DLcWLEad23TOexKCSTO1f+N7i3UDreGVfycaVNbOspJnr9aDw=="
assert.Equal(t, r, expect)
assert.NoError(t, VerifySignProofStr(addr, hex.EncodeToString(pub), r, proof))
```

### Dapp sign request

```go

var r MultiRequest
//code is the data from dapp on TON chain.
code := `{
        "messages":[{
	"address": "EQARULUYsmJq1RiZ-YiH-IJLcAZUVkVff-KBPwEmmaQGH6aC",
	"amount": "195000000",
	"payload":"te6cckEBAgEAigABaw+KfqUAAABqOXlveTmJaAgA7zuZAqJxsqAciTilI8/iTnGEeq62piAAHtRKd6wOcJwQLBuBAwEAnSWThWGAHIXiG4S2uBKfvTnDWV0CiqXAwDzv4KIacQogkCmsj0NuCxykPOZl3QAo0GtsbOJdNcL0J61peOgcvbzLsvXBsnC6HO6YLfLMvtDoIHzZ"
	}],
        "from": "0:a341adb1b38974d70bd09eb5a5e3a072f6f32ecbd706c9c2e873ba60b7cb32fb",
 "valid_until": 1730335778,
        "network": "-239"
}`
//your nonce or seqno
nonce := uint32(180)
err := json.Unmarshal([]byte(code), &r)
assert.NoError(t, err)
//your private key seed 
seed, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
address, err := NewAddress(seed, wallet.V4R2)
fmt.Println(address)
assert.NoError(t, err)
assert.NoError(t, r.Check())
//simulate means that now it IS  a simulative transaction.and false means that now it is NOT a simulative transaction.
//while simulate is true,the result is  Used to estimate gas charges
simulate := true
s, err := SignMultiTransfer(seed, nil, nonce, &r, simulate, wallet.V4R2)
assert.NoError(t, err)
fmt.Println(s.Tx)
tt := &testSignedTx{
Address:      s.Address,
Body:         s.Tx,
InitData:     s.Data,
InitCode:     s.Code,
IgnoreChksig: true,
}
fmt.Println(tt.Str())
```

## Credits  This project includes code adapted from the following sources:  
- [tonutils-go](https://github.com/xssnick/tonutils-go) - Ton Go SDK
- [tongo](https://github.com/tonkeeper/tongo) - TonKeeper Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License

Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/ton/LICENSE>) licensed, see
package or folder for the respective license.
