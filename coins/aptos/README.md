# @okxweb3/coin-aptos
Aptos SDK is used to interact with the Aptos blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get github.com/okx/go-wallet-sdk/coins/aptos
```

## Usage

### Generate private key

```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
let key = await wallet.getRandomPrivateKey();
```

### Private key derivation

```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
let mnemonic = "bean mountain minute enemy state always weekend accuse flag wait island tortoise";
let param = {
    mnemonic: mnemonic,
    hdPath: "m/44'/637'/0'/0'/0'"
};
let privateKey = await wallet.getDerivedPrivateKey(param);
```

### Generate address

```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
let params: NewAddressParams = {
    privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
    addressType: "short",
};
let address = await wallet.getNewAddress(params);
```

### Verify address
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
let p: ValidAddressParams = {
    address: "0x8e6d339ff6096080a4d91c291b297d3814ff9daa34e0f5562d4e7d442cafecdc"
};
let valid = await wallet.validAddress(p);
```

### Transfer APTOS
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
const ts = Math.floor(Date.now()/1000) + 3000
const param: AptosParam =  {
    type: "transfer",
    base: {
        sequenceNumber: 1n,
        chainId: 32,
        maxGasAmount: 10000n,
        gasUnitPrice: 100n,
        expirationTimestampSecs: BigInt(ts),
    },
    data: {
        recipientAddress: "0x0163f9f9f773f3b0e788559d9efcbe547889500d0891fe024e782c7224defd01",
        amount: 1000,
    }
}
let signParams: SignTxParams = {
  privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
  data: param
};
let tx =await wallet.signTransaction(signParams);
```

### Transfer token-register
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
  const ts = Math.floor(Date.now()/1000) + 3000
  const param: AptosParam =  {
    type: "tokenRegister",
    base: {
      sequenceNumber: 4n,
      chainId: 32,
      maxGasAmount: 10000n,
      gasUnitPrice: 100n,
      expirationTimestampSecs: BigInt(ts),
    },
    data: {
      tyArg: "0x02961adfe972ee3c5ce70472cddd1a69803ad45d712d95e3c65480d44305d975::moon_coin::MoonCoin",
    }
  }
  let signParams: SignTxParams = {
    privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
    data: param
  };
  let tx = await wallet.signTransaction(signParams);
```

### transfer token-transfering
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
const ts = Math.floor(Date.now()/1000) + 3000
const param: AptosParam =  {
  type: "tokenTransfer",
  base: {
    sequenceNumber: 5n,
    chainId: 32,
    maxGasAmount: 10000n,
    gasUnitPrice: 100n,
    expirationTimestampSecs: BigInt(ts),
  },
  data: {
    tyArg: "0x02961adfe972ee3c5ce70472cddd1a69803ad45d712d95e3c65480d44305d975::moon_coin::MoonCoin",
    recipientAddress: "0x0163f9f9f773f3b0e788559d9efcbe547889500d0891fe024e782c7224defd01",
    amount: 1000,
  }
}
let signParams: SignTxParams = {
  privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
  data: param
};
let tx = await wallet.signTransaction(signParams);
```

### dex contract call
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
const callData = "{\n    \"function\":\"0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::scripts::swap\",\n    \"type_arguments\":[\n        \"0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::USDT\",\n        \"0xb4d7b2466d211c1f4629e8340bb1a9e75e7f8fb38cc145c54c5c9f9d5017a318::coins_extended::USDC\",\n        \"0xb4d7b2466d211c1f4629e8340bb1a9e75e7f8fb38cc145c54c5c9f9d5017a318::lp::LP<0xb4d7b2466d211c1f4629e8340bb1a9e75e7f8fb38cc145c54c5c9f9d5017a318::coins_extended::USDC, 0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::USDT>\"\n    ],\n    \"arguments\":[\n        \"0xb4d7b2466d211c1f4629e8340bb1a9e75e7f8fb38cc145c54c5c9f9d5017a318\",\n        \"1000000\",\n        \"465087\"\n    ],\n    \"type\":\"entry_function_payload\"\n}"
const moduleData = "[]"
const ts = Math.floor(Date.now()/1000) + 3000
const param: AptosParam =  {
  type: "dapp",
  base: {
    sequenceNumber: "6",
    chainId: 32,
    maxGasAmount: "10000",
    gasUnitPrice: "100",
    expirationTimestampSecs: ts.toString(),
  },
  data: {
    abi: moduleData,
    data: callData,
  }
}
let signParams: SignTxParams = {
  privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
  data: param
};
let tx = await wallet.signTransaction(signParams);
```

### Sign message
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
let param = {
    privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc", 
    data: "aptos message"
}
let data = await wallet.signMessage(param);
```

### offerNFT
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
const ts = Math.floor(Date.now()/1000) + 3000
const param: AptosParam =  {
  type: "offerNft",
  base: {
    sequenceNumber: "6",
    chainId: 32,
    maxGasAmount: "10000",
    gasUnitPrice: "100",
    expirationTimestampSecs: ts.toString(),
  },
  data: {
    receiver: "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28",
    creator: "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28",
    collectionName: "collect_test",
    tokenName: "nft_test",
    version: "1",
    amount: "1"
  }
}
let signParams: SignTxParams = {
  privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
  data: param
};
let tx = await wallet.signTransaction(signParams);
```

### claimNFT
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
const ts = Math.floor(Date.now()/1000) + 3000
const param: AptosParam =  {
  type: "claimNft",
  base: {
    sequenceNumber: "6",
    chainId: 32,
    maxGasAmount: "10000",
    gasUnitPrice: "100",
    expirationTimestampSecs: ts.toString(),
  },
  data: {
    sender: "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28",
    creator: "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28",
    collectionName: "collect_test",
    tokenName: "nft_test",
    version: "1"
  }
}
let signParams: SignTxParams = {
  privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
  data: param
};

let tx = await wallet.signTransaction(signParams);
```

### offerNFT_simulate
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
const ts = Math.floor(Date.now()/1000) + 3000
const param: AptosParam =  {
  type: "offerNft_simulate",
  base: {
    sequenceNumber: "6",
    chainId: 32,
    maxGasAmount: "10000",
    gasUnitPrice: "100",
    expirationTimestampSecs: ts.toString(),
  },
  data: {
    receiver: "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28",
    creator: "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28",
    collectionName: "collect_test",
    tokenName: "nft_test",
    version: "1",
    amount: "1"
  }
}
let signParams: SignTxParams = {
  privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
  data: param
};
let tx = await wallet.signTransaction(signParams);
```

### claimNFT_simulate
```typescript
import {AptosWallet} from "@okxweb3/coin-aptos";

let wallet = new AptosWallet()
const ts = Math.floor(Date.now()/1000) + 3000
const param: AptosParam =  {
  type: "claimNft_simulate",
  base: {
    sequenceNumber: "6",
    chainId: 32,
    maxGasAmount: "10000",
    gasUnitPrice: "100",
    expirationTimestampSecs: ts.toString(),
  },
  data: {
    sender: "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28",
    creator: "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28",
    collectionName: "collect_test",
    tokenName: "nft_test",
    version: "1"
  }
}
let signParams: SignTxParams = {
  privateKey: "f4118e8a1193bf164ac2223f7d0e9c625d6d5ca19d2fbfea7c55d3c0d0284cd0312a81c872aad3a910157ca7b05e70fe2e62aed55b4a14ad033db4556c1547dc",
  data: param
};
let tx = await wallet.signTransaction(signParams);
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/LICENSE>) licensed, see package or folder for the respective license.
