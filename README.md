# OKX Web3 Go Wallet SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/okx/go-wallet-sdk.svg)](https://pkg.go.dev/github.com/okx/go-wallet-sdk)
[![License](https://img.shields.io/github/license/okx/go-wallet-sdk)](https://github.com/okx/go-wallet-sdk/blob/main/LICENSE)

The OKX Web3 Go Wallet SDK is a comprehensive solution for building wallet applications with offline transaction capabilities across multiple blockchain networks. It provides a unified interface for account management, transaction creation, and signing across various mainstream public chains.

## 🚀 Features

-   **Multi-chain support:** Seamlessly interact with major blockchains.
-   **Offline transaction signing:** Ensure security with local signing.
-   **Account generation and management:** Derive addresses with ease.
-   **Customizable transaction creation:** Flexible parameters for all supported chains.
-   **BRC20/Atomical/Runes support:** Full Bitcoin token standard compatibility.
-   **Extensible architecture:** Modular design for future blockchain integration.

## 📚 Documentation

For detailed documentation and API references, please refer to the README files located within each blockchain directory under the `coins` folder. Each directory, such as `aptos`, `bitcoin`, and others, contains specific usage instructions and implementation details.

Example:

-   [Aptos README](https://github.com/okx/go-wallet-sdk/tree/main/coins/aptos)
-   [Bitcoin README](https://github.com/okx/go-wallet-sdk/tree/main/coins/bitcoin)

## 🌐 Supported Chains

The OKX Web3 Go Wallet SDK supports a wide range of blockchain networks. EVM-compatible chains (e.g., BSC, Polygon,
Arbitrum) and Solana-based chains can seamlessly reuse the same code structure for streamlined integration.

| Blockchain | Generate Address | Sign Transaction | Sign Message |
| ---------- | ---------------- | ---------------- | ------------ |
| Aptos      | ✅               | ✅               | ✅           |
| Bitcoin    | ✅               | ✅               | ✅           |
| Cardano    | ✅               | ✅               | ✅           |
| Cosmos     | ✅               | ✅               | ✅           |
| Ethereum   | ✅               | ✅               | ✅           |
| Kaspa      | ✅               | ✅               | ✅           |
| Near       | ✅               | ✅               | ✅           |
| Solana     | ✅               | ✅               | ✅           |
| Starknet   | ✅               | ✅               | ✅           |
| Stacks     | ✅               | ✅               | ✅           |
| Sui        | ✅               | ✅               | ✅           |
| Ton        | ✅               | ✅               | ✅           |
| Tron       | ✅               | ✅               | ✅           |

_Note: Bitcoin support includes BRC20, Atomicals, and Runes-related functions, such as deployment, minting, transfer, and trading._

## 🛠️ Architecture

The Go Wallet SDK follows a modular architecture, comprising the following core components:

1. **`coins`**: Implements transaction creation and signing for each blockchain.
2. **`crypto`**: Manages general cryptographic operations and signature algorithms.
3. **`util`**: Provides helper utilities for common operations.

This structure allows for easy integration and extension of new blockchains.

## 📦 Installation

To install the OKX Web3 Go Wallet SDK, ensure you have Go 1.22+ installed, then run:

```shell
# Install SDK
go get -u github.com/okx/go-wallet-sdk/coins/bitcoin
```

## ⚙️ Build and Test

To build and test all blockchain modules, use the `build.sh` script located in the project root. This script automatically discovers all modules under `coins/`, `crypto/`, `util/`, and `example/` directories, then runs a 4-step build process for each module.

### Build Steps

Each module goes through the following steps:

1. **go mod tidy** - Clean and verify dependencies
2. **go mod edit** - Set toolchain configuration
3. **go build** - Compile the module and all subpackages
4. **go test** - Run all tests

If any step fails, subsequent steps are skipped for that module.

### Basic Usage

```shell
# Interactive mode (prompts if previous failures exist)
sh build.sh

# Run all modules
sh build.sh all

# Run only previously failed modules
sh build.sh failed

# Run specific modules only
sh build.sh bitcoin
sh build.sh bitcoin,ethereum,ton

# Ignore specific modules
sh build.sh all -i=zksync,zkspace
sh build.sh failed -i=zksync,zkspace

```

### Command Options

```
Usage: ./build.sh [all|failed|mod1,mod2,...] [-i=module1,module2,...]
```

| Option                   | Description                                                                                 |
| ------------------------ | ------------------------------------------------------------------------------------------- |
| (no args)                | Interactive mode - prompts to run all or only failed modules if `build_failures.log` exists |
| `all`                    | Run all modules                                                                             |
| `failed`                 | Run only previously failed modules                                                          |
| `mod1,mod2,...`          | Run only specified modules (comma-separated list)                                           |
| `-i=module1,module2,...` | Ignore specific modules (comma-separated list of module names to skip)                      |

### Output

Each module displays step-by-step progress:

```
==========================================
[coins/bitcoin]
==========================================
  [1/4] go mod tidy    ... ✓ PASS
  [2/4] go mod edit    ... ✓ PASS
  [3/4] go build       ... ✓ PASS
  [4/4] go test        ... ✓ PASS

  ✓ bitcoin ALL STEPS PASSED
```

-   A summary shows total success/failure/ignored counts
-   Failed modules are logged to `build_failures.log` with detailed output
-   Ignored modules are also tracked in the log file
-   Use `sh build.sh failed` to quickly re-run only the failed modules after fixing issues

## 💬 Feedback and Support

The OKX Web3 Go Wallet SDK shares common design principles and usage patterns with the JS SDK. While each blockchain's specific usage can be found in the corresponding `coins` directory README, users can refer to the [JS SDK demo](https://okx.github.io/wallet-sdk-demo/) and [documentation](https://okx.github.io/js-wallet-sdk/) for additional guidance. If you encounter any issues or have suggestions, please submit them through [GitHub Issues](https://github.com/okx/go-wallet-sdk/issues), and we will address them promptly.

## Change Log

[detail](./CHANGELOG.md)

## 🔒 Security

If you find security risks, it is recommended to report them through the following channels and get your reward!
Submit on HackerOne platform: [https://hackerone.com/okg](https://hackerone.com/okg) or on our OKX feedback submission page: [https://www.okx.com/feedback/submit](https://www.okx.com/feedback/submit).

## 📜 License

The OKX Web3 Go Wallet SDK is open-source software licensed under the [MIT license](LICENSE).
