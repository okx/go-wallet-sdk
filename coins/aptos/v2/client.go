package v2

// NetworkConfig a configuration for the Client and which
// network to use.  Use one of the preconfigured LocalnetConfig,
// DevnetConfig, TestnetConfig, or MainnetConfig unless you have
// your own full node
type NetworkConfig struct {
	Name       string
	ChainId    uint8
	NodeUrl    string
	IndexerUrl string
	FaucetUrl  string
}

var LocalnetConfig = NetworkConfig{
	Name:    "localnet",
	ChainId: 4,
	// We use 127.0.0.1 as it is more foolproof than localhost
	//NodeUrl:    "http://127.0.0.1:8080/v1",
	//IndexerUrl: "http://127.0.0.1:8090/v1/graphql",
	//FaucetUrl:  "http://127.0.0.1:8081/v1",
}
var DevnetConfig = NetworkConfig{
	Name:    "devnet",
	ChainId: 135,
	//NodeUrl:    "https://api.devnet.aptoslabs.com/v1",
	//IndexerUrl: "https://api.devnet.aptoslabs.com/v1/graphql",
	//FaucetUrl:  "https://faucet.devnet.aptoslabs.com/",
}
var TestnetConfig = NetworkConfig{
	Name:    "testnet",
	ChainId: 2,
	//NodeUrl:    "https://api.testnet.aptoslabs.com/v1",
	//IndexerUrl: "https://api.testnet.aptoslabs.com/v1/graphql",
	//FaucetUrl:  "https://faucet.testnet.aptoslabs.com/",
}
var MainnetConfig = NetworkConfig{
	Name:    "mainnet",
	ChainId: 1,
	//NodeUrl:    "https://api.mainnet.aptoslabs.com/v1",
	//IndexerUrl: "https://api.mainnet.aptoslabs.com/v1/graphql",
	//FaucetUrl:  "",
}

// NamedNetworks Map from network name to NetworkConfig
var NamedNetworks map[string]NetworkConfig

func init() {
	NamedNetworks = make(map[string]NetworkConfig, 4)
	setNN := func(nc NetworkConfig) {
		NamedNetworks[nc.Name] = nc
	}
	setNN(LocalnetConfig)
	setNN(DevnetConfig)
	setNN(TestnetConfig)
	setNN(MainnetConfig)
}

// Client is a facade over the multiple types of underlying clients, as the user doesn't actually care where the data
// comes from.  It will be then handled underneath
type Client struct {
	nodeClient *NodeClient
	//faucetClient  *FaucetClient
	//indexerClient *IndexerClient
}

// NewClient Creates a new client with a specific network config that can be extended in the future
func NewClient(config NetworkConfig) (client *Client, err error) {
	nodeClient, err := NewNodeClient(config.NodeUrl, config.ChainId)
	if err != nil {
		return nil, err
	}
	// Indexer may not be present
	/*var indexerClient *IndexerClient = nil
	if config.IndexerUrl != "" {
		indexerClient = NewIndexerClient(nodeClient.client, config.IndexerUrl)
	}

	// Faucet may not be present
	var faucetClient *FaucetClient = nil
	if config.FaucetUrl != "" {
		faucetClient, err = NewFaucetClient(nodeClient, config.FaucetUrl)
		if err != nil {
			return nil, err
		}
	}

	// Fetch the chain Id if it isn't in the config
	if config.ChainId == 0 {
		_, _ = nodeClient.GetChainId()
	}*/

	client = &Client{
		nodeClient,
		/*faucetClient,
		indexerClient,*/
	}
	return
}

/*// SetTimeout adjusts the HTTP client timeout
func (client *Client) SetTimeout(timeout time.Duration) {
	client.nodeClient.client.Timeout = timeout
}

// Info Retrieves the node info about the network and it's current state
func (client *Client) Info() (info NodeInfo, err error) {
	return client.nodeClient.Info()
}

// Account Retrieves information about the account such as SequenceNumber and AuthKey
func (client *Client) Account(address AccountAddress, ledgerVersion ...int) (info AccountInfo, err error) {
	return client.nodeClient.Account(address, ledgerVersion...)
}

// AccountResource Retrieves a single resource given its struct name.
// Can also fetch at a specific ledger version
func (client *Client) AccountResource(address AccountAddress, resourceType string, ledgerVersion ...int) (data map[string]any, err error) {
	return client.nodeClient.AccountResource(address, resourceType, ledgerVersion...)
}

// AccountResources fetches resources for an account into a JSON-like map[string]any in AccountResourceInfo.Data
// For fetching raw Move structs as BCS, See #AccountResourcesBCS
func (client *Client) AccountResources(address AccountAddress, ledgerVersion ...int) (resources []AccountResourceInfo, err error) {
	return client.nodeClient.AccountResources(address, ledgerVersion...)
}

// AccountResourcesBCS fetches account resources as raw Move struct BCS blobs in AccountResourceRecord.Data []byte
func (client *Client) AccountResourcesBCS(address AccountAddress, ledgerVersion ...int) (resources []AccountResourceRecord, err error) {
	return client.nodeClient.AccountResourcesBCS(address, ledgerVersion...)
}

// BlockByHeight fetches a block by height
func (client *Client) BlockByHeight(blockHeight uint64, withTransactions bool) (data map[string]any, err error) {
	return client.nodeClient.BlockByHeight(blockHeight, withTransactions)
}

// BlockByVersion fetches a block by ledger version
func (client *Client) BlockByVersion(ledgerVersion uint64, withTransactions bool) (data map[string]any, err error) {
	return client.nodeClient.BlockByVersion(ledgerVersion, withTransactions)
}

// TransactionByHash gets info on a transaction
// The transaction may be pending or recently committed.
//
//	data, err := c.TransactionByHash("0xabcd")
//	if err != nil {
//		if httpErr, ok := err.(aptos.HttpError) {
//			if httpErr.StatusCode == 404 {
//				// if we're sure this has been submitted, assume it is still pending elsewhere in the mempool
//			}
//		}
//	} else {
//		if data["type"] == "pending_transaction" {
//			// known to local mempool, but not committed yet
//		}
//	}
func (client *Client) TransactionByHash(txnHash string) (data map[string]any, err error) {
	return client.nodeClient.TransactionByHash(txnHash)
}

// TransactionByVersion gets info on a transaction from its LedgerVersion.  It must have been
// committed to have a ledger version
func (client *Client) TransactionByVersion(version uint64) (data map[string]any, err error) {
	return client.nodeClient.TransactionByVersion(version)
}

// PollForTransactions Waits up to 10 seconds for transactions to be done, polling at 10Hz
// Accepts options PollPeriod and PollTimeout which should wrap time.Duration values.
func (client *Client) PollForTransactions(txnHashes []string, options ...any) error {
	return client.nodeClient.PollForTransactions(txnHashes, options...)
}

// WaitForTransaction Do a long-GET for one transaction and wait for it to complete
func (client *Client) WaitForTransaction(txnHash string) (data map[string]any, err error) {
	return client.nodeClient.WaitForTransaction(txnHash)
}

// Transactions Get recent transactions.
// Start is a version number. Nil for most recent transactions.
// Limit is a number of transactions to return. 'about a hundred' by default.
func (client *Client) Transactions(start *uint64, limit *uint64) (data []map[string]any, err error) {
	return client.nodeClient.Transactions(start, limit)
}

// SubmitTransaction Submits an already signed transaction to the blockchain
func (client *Client) SubmitTransaction(signedTransaction *SignedTransaction) (data map[string]any, err error) {
	return client.nodeClient.SubmitTransaction(signedTransaction)
}

// GetChainId Retrieves the ChainId of the network
// Note this will be cached forever, or taken directly from the config
func (client *Client) GetChainId() (chainId uint8, err error) {
	return client.nodeClient.GetChainId()
}

// Fund Uses the faucet to fund an address, only applies to non-production networks
func (client *Client) Fund(address AccountAddress, amount uint64) error {
	return client.faucetClient.Fund(address, amount)
}

// BuildTransaction Builds a raw transaction from the payload and fetches any necessary information
// from on-chain
func (client *Client) BuildTransaction(sender AccountAddress, payload TransactionPayload, options ...any) (rawTxn *RawTransaction, err error) {
	return client.nodeClient.BuildTransaction(sender, payload, options...)
}

// BuildSignAndSubmitTransaction Convenience function to do all three in one
// for more configuration, please use them separately
func (client *Client) BuildSignAndSubmitTransaction(sender *Account, payload TransactionPayload, options ...any) (hash string, err error) {
	return client.nodeClient.BuildSignAndSubmitTransaction(sender, payload, options...)
}

// View Runs a view function on chain returning a list of return values.
// TODO: support ledger version
func (client *Client) View(payload *ViewPayload) (vals []any, err error) {
	return client.nodeClient.View(payload)
}

// EstimateGasPrice Retrieves the gas estimate from the network.
func (client *Client) EstimateGasPrice() (info EstimateGasInfo, err error) {
	return client.nodeClient.EstimateGasPrice()
}

// QueryIndexer queries the indexer using GraphQL to fill the `query` struct with data.  See examples in the indexer
// client on how to make queries
func (client *Client) QueryIndexer(query any, variables map[string]any, options ...graphql.Option) error {
	return client.indexerClient.Query(query, variables, options...)
}

// GetProcessorStatus returns the ledger version up to which the processor has processed
func (client *Client) GetProcessorStatus(processorName string) (uint64, error) {
	return client.indexerClient.GetProcessorStatus(processorName)
}

// GetCoinBalances gets the balances of all coins associated with a given address
func (client *Client) GetCoinBalances(address AccountAddress) ([]CoinBalance, error) {
	return client.indexerClient.GetCoinBalances(address)
}*/
