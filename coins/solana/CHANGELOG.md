# Change Log

All notable changes to this project will be documented in this file.

# [0.0.5](https://github.com/okx/go-wallet-sdk) (2025-09-22)

-   **solana-sdk:** add `encoding` parameter to `NewTxFromRaw` and `AddSignature` functions to support base64 encoding([121](https://github.com/okx/go-wallet-sdk/pull/121))
    -   `NewTxFromRaw(rawTx string)` → `NewTxFromRaw(rawTx string, encoding string)`
    -   `AddSignature(tx types.Transaction, sig []byte)` → `AddSignature(tx types.Transaction, sig []byte, encoding string)`
    -   Supported encodings: "base58" (default behavior) and "base64"

# [0.0.4](https://github.com/okx/go-wallet-sdk) (2025-02-14)

### New features

-   **solana-sdk:** add message signing and verification and hash calculation ([80](https://github.com/okx/go-wallet-sdk/pull/80))

# [0.0.2](https://github.com/okx/go-wallet-sdk) (2023-11-20)

### updates

-   **solana-sdk:** change some files name and remove some unused libs ([21](https://github.com/okx/go-wallet-sdk/pull/21))
