# OKX Web3 Go Wallet SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/okx/go-wallet-sdk.svg)](https://pkg.go.dev/github.com/okx/go-wallet-sdk)
[![License](https://img.shields.io/github/license/okx/go-wallet-sdk)](https://github.com/okx/go-wallet-sdk/blob/main/LICENSE)

The OKX Web3 Go Wallet SDK is a comprehensive solution for building wallet applications with offline transaction capabilities across multiple blockchain networks. It provides a unified interface for account management, transaction creation, and signing across various mainstream public chains.

## 🚀 Features

- **Multi-chain support:** Seamlessly interact with major blockchains.
- **Offline transaction signing:** Ensure security with local signing.
- **Account generation and management:** Derive addresses with ease.
- **Customizable transaction creation:** Flexible parameters for all supported chains.
- **BRC20/Atomical/Runes support:** Full Bitcoin token standard compatibility.
- **Extensible architecture:** Modular design for future blockchain integration.

## 📚 Documentation

For detailed documentation and API references, please refer to the README files located within each blockchain directory under the `coins` folder. Each directory, such as `aptos`, `bitcoin`, and others, contains specific usage instructions and implementation details.

Example:
- [Aptos README](https://github.com/okx/go-wallet-sdk/tree/main/coins/aptos)
- [Bitcoin README](https://github.com/okx/go-wallet-sdk/tree/main/coins/bitcoin)



## 🌐 Supported Chains
The OKX Web3 Go Wallet SDK supports a wide range of blockchain networks. EVM-compatible chains (e.g., BSC, Polygon, Arbitrum) and Solana-based chains can seamlessly reuse the same code structure for streamlined integration.

| Blockchain | Account Generation | Transaction Creation | Transaction Signing |
| ---------- | ------------------ | -------------------- | ------------------- |
| Aptos      | ✅                  | ✅                    | ✅                   |
| Avax       | ✅                  | ✅                    | ✅                   |
| Bitcoin    | ✅                  | ✅                    | ✅                   |
| Cardano    | ✅                  | ✅                    | ✅                   |
| Cosmos     | ✅                  | ✅                    | ✅                   |
| Ethereum   | ✅                  | ✅                    | ✅                   |
| Filecoin   | ✅                  | ✅                    | ✅                   |
| Harmony    | ✅                  | ✅                    | ✅                   |
| Kaspa      | ✅                  | ✅                    | ✅                   |
| Near       | ✅                  | ✅                    | ✅                   |
| NostrAsset | ✅                  | ✅                    | ✅                   |
| Solana     | ✅                  | ✅                    | ✅                   |
| Starknet   | ✅                  | ✅                    | ✅                   |
| Stacks     | ✅                  | ✅                    | ✅                   |
| SUI        | ✅                  | ✅                    | ✅                   |
| Ton        | ✅                  | ✅                    | ✅                   |
| Tron       | ✅                  | ✅                    | ✅                   |


*Note: Bitcoin support includes BRC20, Atomicals, and Runes-related functions, such as deployment, minting, transfer, and trading.*

## 🛠️ Architecture

The Go Wallet SDK follows a modular architecture, comprising the following core components:

1. **`coins`**: Implements transaction creation and signing for each blockchain.
2. **`crypto`**: Manages general cryptographic operations and signature algorithms.
3. **`util`**: Provides helper utilities for common operations.

This structure allows for easy integration and extension of new blockchains.

## 📦 Installation

To install the OKX Web3 Go Wallet SDK, ensure you have Go 1.22+ installed, such as run:

```shell
# Install SDK
go get -u github.com/okx/go-wallet-sdk/coins/bitcoin
```

## ⚙️ Build and Test

To build and test all blockchain modules, use the `build.sh` script located in the project root. This script iterates through each chain module under the `coins` directory, runs `go mod tidy` to clean dependencies, executes tests, and verifies successful builds.

```shell
sh build.sh
```

The output will display the build status for each chain. If a module fails, the error message will indicate the issue for further debugging.



## 💬 Feedback and Support

The OKX Web3 Go Wallet SDK shares common design principles and usage patterns with the JS SDK. While each blockchain's specific usage can be found in the corresponding `coins` directory README, users can refer to the [JS SDK demo](https://okx.github.io/wallet-sdk-demo/) and [documentation](https://okx.github.io/js-wallet-sdk/) for additional guidance. If you encounter any issues or have suggestions, please submit them through [GitHub Issues](https://github.com/okx/go-wallet-sdk/issues), and we will address them promptly.

## 🔒 Security

If you find security risks, it is recommended to feedback through the following channels and get your reward!
Submit on HackerOne platform: [https://hackerone.com/okg](https://hackerone.com/okg) or on our OKX feedback submission page: [https://www.okx.com/feedback/submit](https://www.okx.com/feedback/submit).

## 📜 License

The OKX Web3 Go Wallet SDK is open-source software licensed under the [MIT license](LICENSE).

