# go-wallet-sdk

This is a Go language wallet solution that supports offline transactions. We currently support various mainstream public blockchains, and will gradually release the source codes for each blockchain.

## Supported chains

- BTC: Supports transaction creation and signing. Also Supports BRC20-related functions, including inscription creation, BRC20 buying and selling.
- Ethereum: Supports transaction creation and signing.
- Filecoin: Supports transaction creation and signing.
- Polkadot: Supports transaction creation and signing.
- Starknet: Supports transaction creation and signing.

## Main modules

- coins: Implements transaction creation and signature in each coin type.
- crypto: Handles general security and signature algorithms.
- util: Provides various utility class methods.

## Example

For specific usage examples of each coin type, please refer to the corresponding test files. Remember to replace the placeholder private key with your own private key, which is generally in hex format.

## Feedback

You can provide feedback directly in GitHub Issues, and we will respond promptly.
