# go-wallet-sdk

This is a Go language wallet solution that supports offline transactions. We currently support various mainstream public
blockchains, and will gradually release the source codes for each blockchain.

## Supported chains

|          | Account Generation | Transaction Creation | Transaction Signing |
|----------|-------------------|----------------------|---------------------|
| BTC      | ✅                 | ✅                    | ✅                   | 
| Ethereum | ✅                 | ✅                    | ✅                   |
| EOS      | ✅                 | ✅                    | ✅                   |
| Filecoin | ✅                 | ✅                    | ✅                   |
| Polkadot | ✅                 | ✅                    | ✅                   |
| Starknet | ✅                 | ✅                    | ✅                   |
| Aptos    | ✅                 | ✅                    | ✅                   |
| Near     | ✅                 | ✅                    | ✅                   |
| Solana   | ✅                 | ✅                    | ✅                   |
| Stacks   | ✅                 | ✅                    | ✅                   |
| SUI      | ✅                 | ✅                    | ✅                   |
| Tron     | ✅                 | ✅                    | ✅                   |
| Cosmos   | ✅                 | ✅                    | ✅                   |
| Axelar   | ✅                 | ✅                    | ✅                   |
| Cronos   | ✅                 | ✅                    | ✅                   |
| Evmos    | ✅                 | ✅                    | ✅                   |
| Iris     | ✅                 | ✅                    | ✅                   |
| Juno     | ✅                 | ✅                    | ✅                   |
| Kava     | ✅                 | ✅                    | ✅                   |
| Kujira   | ✅                 | ✅                    | ✅                   |
| Okc      | ✅                 | ✅                    | ✅                   |
| Osmosis  | ✅                 | ✅                    | ✅                   |
| Secret   | ✅                 | ✅                    | ✅                   |
| Sei      | ✅                 | ✅                    | ✅                   |
| Stargaze | ✅                 | ✅                    | ✅                   |
| Terra    | ✅                 | ✅                    | ✅                   |
| Tia      | ✅                 | ✅                    | ✅                   |
| Avax     | ✅                 | ✅                    | ✅                   |
| Elrond   | ✅                 | ✅                    | ✅                   |
| Flow     | ✅                 | ✅                    | ✅                   |
| Harmony  | ✅                 | ✅                    | ✅                   |
| Helium   | ✅                 | ✅                    | ✅                   |
| Kaspa    | ✅                 | ✅                    | ✅                   |
| Nervos   | ✅                 | ✅                    | ✅                   |
| Oasis    | ✅                 | ✅                    | ✅                   |
| Tezos    | ✅                 | ✅                    | ✅                   |
| Waves    | ✅                 | ✅                    | ✅                   |
| Zil      | ✅                 | ✅                    | ✅                   |
| Zkspace  | ✅                 | ✅                    | ✅                   |
| Zksync   | ✅                 | ✅                    | ✅                   |

*BTC: Supports Supports BRC20-related functions, including inscription creation, BRC20 buying and selling.

## Main modules

- coins: Implements transaction creation and signature in each coin type.
- crypto: Handles general security and signature algorithms.
- util: Provides various utility class methods.

## Example

For specific usage examples of each coin type, please refer to the corresponding test files. Remember to replace the
placeholder private key with your own private key, which is generally in hex format.

## Feedback

You can provide feedback directly in GitHub Issues, and we will respond promptly.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/LICENSE>) licensed, see package or folder for the respective license.
