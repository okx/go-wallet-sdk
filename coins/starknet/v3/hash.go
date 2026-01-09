package v3

import (
	"errors"
	"github.com/okx/go-wallet-sdk/coins/starknet"
	"github.com/okx/go-wallet-sdk/coins/starknet/juno_core/crypto"
	"github.com/okx/go-wallet-sdk/coins/starknet/juno_core/felt"
	"math/big"
)

var (
	ErrNotAllParametersSet = errors.New("not all neccessary parameters have been set")
	ErrFeltToBigInt        = errors.New("felt to BigInt error")
)

func DataAvailabilityModeConc(feeDAMode, nonceDAMode DataAvailabilityMode) (uint64, error) {
	const dataAvailabilityModeBits = 32
	fee64, err := feeDAMode.UInt64()
	if err != nil {
		return 0, err
	}
	nonce64, err := nonceDAMode.UInt64()
	if err != nil {
		return 0, err
	}
	return fee64 + nonce64<<dataAvailabilityModeBits, nil
}

func TipAndResourcesHash(tip uint64, resourceBounds ResourceBoundsMapping) (*felt.Felt, error) {
	l1Bytes, err := resourceBounds.L1Gas.Bytes(ResourceL1Gas)
	if err != nil {
		return nil, err
	}
	l2Bytes, err := resourceBounds.L2Gas.Bytes(ResourceL2Gas)
	if err != nil {
		return nil, err
	}
	l1DataGasBytes, err := resourceBounds.L1DataGas.Bytes(ResourceL1DataGas)
	if err != nil {
		return nil, err
	}
	l1Bounds := new(felt.Felt).SetBytes(l1Bytes)
	l2Bounds := new(felt.Felt).SetBytes(l2Bytes)
	l1DataGasBounds := new(felt.Felt).SetBytes(l1DataGasBytes)
	return crypto.PoseidonArray(new(felt.Felt).SetUint64(tip), l1Bounds, l2Bounds, l1DataGasBounds), nil
}

// TransactionHashInvokeV3 calculates the transaction hash for a invoke V3 transaction.
//
// Parameters:
//   - txn: The invoke V3 transaction to calculate the hash for
//   - chainId: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashInvokeV3(txn *InvokeTxnV3, chainId *felt.Felt) (*felt.Felt, error) {
	// https://github.com/starknet-io/SNIPs/blob/main/SNIPS/snip-8.md#protocol-changes
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation
	if txn.Version == "" || txn.ResourceBounds == (ResourceBoundsMapping{}) || len(txn.Calldata) == 0 || txn.Nonce == nil || txn.SenderAddress == nil || txn.PayMasterData == nil || txn.AccountDeploymentData == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := DataAvailabilityModeConc(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}
	tipAndResourceHash, err := TipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}
	return crypto.PoseidonArray(
		PREFIX_TRANSACTION,
		txnVersionFelt,
		txn.SenderAddress,
		tipAndResourceHash,
		crypto.PoseidonArray(txn.PayMasterData...),
		chainId,
		txn.Nonce,
		new(felt.Felt).SetUint64(DAUint64),
		crypto.PoseidonArray(txn.AccountDeploymentData...),
		crypto.PoseidonArray(txn.Calldata...),
	), nil
}

// TransactionHashDeployAccountV3 calculates the transaction hash for a deploy account V3 transaction.
//
// Parameters:
//   - txn: The deploy account V3 transaction to calculate the hash for
//   - contractAddress: The contract address as parameters as a *felt.Felt
//   - chainId: The chain ID as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func TransactionHashDeployAccountV3(txn *DeployAccountTxnV3, contractAddress, chainId *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#v3_hash_calculation_3
	if txn.Version == "" || txn.ResourceBounds == (ResourceBoundsMapping{}) || txn.Nonce == nil || txn.PayMasterData == nil {
		return nil, ErrNotAllParametersSet
	}

	txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
	if err != nil {
		return nil, err
	}
	DAUint64, err := DataAvailabilityModeConc(txn.FeeMode, txn.NonceDataMode)
	if err != nil {
		return nil, err
	}
	tipUint64, err := txn.Tip.ToUint64()
	if err != nil {
		return nil, err
	}
	tipAndResourceHash, err := TipAndResourcesHash(tipUint64, txn.ResourceBounds)
	if err != nil {
		return nil, err
	}
	return crypto.PoseidonArray(
		PREFIX_DEPLOY_ACCOUNT,
		txnVersionFelt,
		contractAddress,
		tipAndResourceHash,
		crypto.PoseidonArray(txn.PayMasterData...),
		chainId,
		txn.Nonce,
		new(felt.Felt).SetUint64(DAUint64),
		crypto.PoseidonArray(txn.ConstructorCalldata...),
		txn.ClassHash,
		txn.ContractAddressSalt,
	), nil
}

func BigIntToFelt(b *big.Int) *felt.Felt {
	return felt.HexToFelt(starknet.BigToHex(b))
}

func BigIntToFelts(b []*big.Int) []*felt.Felt {
	f := make([]*felt.Felt, len(b))
	for k, v := range b {
		f[k] = BigIntToFelt(v)
	}
	return f
}
