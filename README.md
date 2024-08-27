# go-wallet-sdk

This is a Go language wallet solution that supports offline transactions. We currently support various mainstream public
blockchains, and will gradually release the source codes for each blockchain.

## Supported chains

|          | Account Generation | Transaction Creation | Transaction Signing |
|----------|-------------------|----------------------|---------------------|
| Aptos    | ✅                 | ✅                    | ✅                   |
| Avax     | ✅                 | ✅                    | ✅                   |
| Axelar   | ✅                 | ✅                    | ✅                   |
| BTC      | ✅                 | ✅                    | ✅                   | 
| Cosmos   | ✅                 | ✅                    | ✅                   |
| Cronos   | ✅                 | ✅                    | ✅                   |
| Elrond   | ✅                 | ✅                    | ✅                   |
| EOS      | ✅                 | ✅                    | ✅                   |
| Ethereum | ✅                 | ✅                    | ✅                   |
| Evmos    | ✅                 | ✅                    | ✅                   |
| Filecoin | ✅                 | ✅                    | ✅                   |
| Flow     | ✅                 | ✅                    | ✅                   |
| Harmony  | ✅                 | ✅                    | ✅                   |
| Helium   | ✅                 | ✅                    | ✅                   |
| Iris     | ✅                 | ✅                    | ✅                   |
| Juno     | ✅                 | ✅                    | ✅                   |
| Kava     | ✅                 | ✅                    | ✅                   |
| Kaspa    | ✅                 | ✅                    | ✅                   |
| Kujira   | ✅                 | ✅                    | ✅                   |
| Near     | ✅                 | ✅                    | ✅                   |
| Nervos   | ✅                 | ✅                    | ✅                   |
| Oasis    | ✅                 | ✅                    | ✅                   |
| Okc      | ✅                 | ✅                    | ✅                   |
| Osmosis  | ✅                 | ✅                    | ✅                   |
| Polkadot | ✅                 | ✅                    | ✅                   |
| Secret   | ✅                 | ✅                    | ✅                   |
| Sei      | ✅                 | ✅                    | ✅                   |
| Solana   | ✅                 | ✅                    | ✅                   |
| Starknet | ✅                 | ✅                    | ✅                   |
| Stacks   | ✅                 | ✅                    | ✅                   |
| Stargaze | ✅                 | ✅                    | ✅                   |
| SUI      | ✅                 | ✅                    | ✅                   |
| Terra    | ✅                 | ✅                    | ✅                   |
| Tezos    | ✅                 | ✅                    | ✅                   |
| Tia      | ✅                 | ✅                    | ✅                   |
| Tron     | ✅                 | ✅                    | ✅                   |
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
