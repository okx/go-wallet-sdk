
# Change Log

All notable changes to this project will be documented in this file.

# [0.0.3](https://github.com/okx/go-wallet-sdk) (2024-01-01)

### updates

- **stacks-sdk:** add support for multiple networks (mainnet/testnet) with new version constants
  - Added network version constants: `MainnetSingleSig`, `MainnetMultiSig`, `TestnetSingleSig`, `TestnetMultiSig`
  - Updated `GetAddressFromPublicKey` function to accept `version` parameter
  - Updated `NewAddress` function to accept `version` parameter
  - Updated tests and README examples to use new API
  - **BREAKING CHANGE:** Function signatures changed - version parameter now required

# [0.0.2](https://github.com/okx/go-wallet-sdk) (2023-11-20)

### updates

- **stacks-sdk:** change some files name and remove some unused libs ([21](https://github.com/okx/go-wallet-sdk/pull/21))
