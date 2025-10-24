package transactionhelper

import (
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/subnetworks"
)

// CoinbaseTransactionIndex is the index of the coinbase transaction in every block
const CoinbaseTransactionIndex = 0

// IsCoinBase determines whether or not a transaction is a coinbase transaction. A coinbase
// transaction is a special transaction created by miners that distributes fees and block subsidy
// to the previous blocks' miners, and to specify the scriptPubKey that will be used to pay the current
// miner in future blocks. Each input of the coinbase transaction should set index to maximum
// value and reference the relevant block id, instead of previous transaction id.
func IsCoinBase(tx *externalapi.DomainTransaction) bool {
	// A coinbase transaction must have subnetwork id SubnetworkIDCoinbase
	return tx.SubnetworkID == subnetworks.SubnetworkIDCoinbase
}
