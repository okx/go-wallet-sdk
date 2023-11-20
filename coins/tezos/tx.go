package tezos

import (
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/tezos/types"
)

const GasSafetyMargin int64 = 100

var (
	// for reveal
	defaultRevealLimits = types.Limits{
		Fee:      1000,
		GasLimit: 1000,
	}
	// for transfers to tz1/2/3
	defaultTransferLimitsEOA = types.Limits{
		Fee:      1000,
		GasLimit: 1420, // 1820 when source is emptied
	}
	// for transfers to manager.tz
	defaultTransferLimitsKT1 = types.Limits{
		Fee:      1000,
		GasLimit: 2078,
	}
	// for delegation
	defaultDelegationLimitsEOA = types.Limits{
		Fee:      1000,
		GasLimit: 1000,
	}
	// for simulating contract calls and other operations
	// used when no explicit costs are set
	defaultSimulationLimits = types.Limits{
		GasLimit:     types.DefaultParams.HardGasLimitPerOperation,
		StorageLimit: types.DefaultParams.HardStorageLimitPerOperation,
	}
)

func NewTransaction(from, to string, amount int64, opts *CallOptions) (*types.Op, error) {
	return NewTransactionByOperation(from, to, amount, types.NewOp(), opts)
}

func NewJakartanetTransaction(from, to string, amount int64, opts *CallOptions) (*types.Op, error) {
	return NewTransactionByOperation(from, to, amount, types.NewJakartanetOp(), opts)
}

func InitOperationSourceBranch(from string, op *types.Op, opts *CallOptions) error {
	valid, err := ValidAddress(from)
	if err != nil || !valid {
		return fmt.Errorf("invalid address %q: %v", from, err)
	}
	sent, err := types.ParseAddress(from)
	if err != nil {
		return fmt.Errorf("Invalid sender error: %q: %v", from, err)
	}
	if !sent.IsValid() {
		return fmt.Errorf("Invalid sender %q: %v", from, err)
	}
	op.WithSource(sent)
	if opts == nil {
		opts = &DefaultOptions
	}
	if opts.BlockHash.IsValid() {
		op.WithBranch(opts.BlockHash)
	}
	return nil
}

func NewTransactionByOperation(from, to string, amount int64, op *types.Op, opts *CallOptions) (*types.Op, error) {
	if op == nil {
		op = types.NewOp()
	}
	if err := InitOperationSourceBranch(from, op, opts); err != nil {
		return nil, err
	}
	recv, err := types.ParseAddress(to)
	if err != nil || !recv.IsValid() {
		return nil, fmt.Errorf("Invalid receiver %q: %v", to, err)
	}
	op.WithTransfer(recv, amount)
	return op, nil
}

func NewDelegationTransaction(from, to string, opts *CallOptions) (*types.Op, error) {
	return NewDelegationTransactionByOperation(from, to, types.NewOp(), opts)
}

func NewJakartanetDelegationTransaction(from, to string, opts *CallOptions) (*types.Op, error) {
	return NewDelegationTransactionByOperation(from, to, types.NewJakartanetOp(), opts)
}

func NewDelegationTransactionByOperation(from, to string, op *types.Op, opts *CallOptions) (*types.Op, error) {
	if op == nil {
		op = types.NewOp()
	}
	if err := InitOperationSourceBranch(from, op, opts); err != nil {
		return nil, err
	}
	toDelegation, err := types.ParseAddress(to)
	if err != nil || !toDelegation.IsValid() {
		return nil, fmt.Errorf("Invalid receiver %q: %v", to, err)
	}
	op.WithDelegation(toDelegation)
	return op, nil
}

func NewUnDelegationTransaction(from string, opts *CallOptions) (*types.Op, error) {
	return NewUnDelegationTransactionByOperation(from, types.NewOp(), opts)
}

func NewJakartanetUnDelegationTransaction(from string, opts *CallOptions) (*types.Op, error) {
	return NewUnDelegationTransactionByOperation(from, types.NewJakartanetOp(), opts)
}

func NewUnDelegationTransactionByOperation(from string, op *types.Op, opts *CallOptions) (*types.Op, error) {
	if op == nil {
		op = types.NewOp()
	}
	if err := InitOperationSourceBranch(from, op, opts); err != nil {
		return nil, err
	}
	op.WithUndelegation()
	return op, nil
}

func AddTransferOpTransaction(to string, amount int64, op *types.Op) (*types.Op, error) {
	if op == nil {
		op = types.NewOp()
	}
	recv, err := types.ParseAddress(to)
	if err != nil || !recv.IsValid() {
		return nil, fmt.Errorf("Invalid receiver %q: %v", to, err)
	}
	op.WithTransfer(recv, amount)
	return op, nil
}

// CompleteTransaction ensures an operation is compatible with the current source account's
// on-chain state. Sets branch for TTL control, replay counters, and reveals
// the sender's pubkey if not published yet.
func CompleteTransaction(o *types.Op, key types.Key, opts *CallOptions) error {
	needBranch := !o.Branch.IsValid()
	needCounter := o.NeedCounter()
	mayNeedReveal := len(o.Contents) > 0 && o.Contents[0].Kind() != types.OpTypeReveal

	if !needBranch && !mayNeedReveal && !needCounter {
		return nil
	}

	// add branch for TTL control
	if needBranch {
		return errors.New("needBranch")
	}

	if needCounter || mayNeedReveal {
		//add reveal if necessary
		if mayNeedReveal && opts.NeedReveal {
			reveal := &types.Reveal{
				Manager: types.Manager{
					Source: key.Address(),
				},
				PublicKey: key,
			}
			reveal.WithLimits(defaultRevealLimits)
			o.WithContentsFront(reveal)
			needCounter = true
		}
		// add counters
		if needCounter {
			nextCounter := opts.Counter + 1
			for _, op := range o.Contents {
				// skip non-manager ops
				if op.GetCounter() < 0 {
					continue
				}
				op.WithCounter(nextCounter)
				nextCounter++
			}
		}
	}
	return nil
}

func BuildTransaction(op *types.Op, fee int64, key types.Key, opts *CallOptions) error {
	if err := CompleteTransaction(op, key, opts); err != nil {
		return err
	}
	if opts == nil || !opts.IgnoreLimits {
		// use default gas/storage limits, set min fee
		for _, oc := range op.Contents {
			l := oc.Limits()
			if l.GasLimit == 0 {
				l.GasLimit = defaultSimulationLimits.GasLimit / int64(len(op.Contents))
			}
			if l.StorageLimit == 0 {
				l.StorageLimit = defaultSimulationLimits.StorageLimit / int64(len(op.Contents))
			}
			l.Fee = fee
			oc.WithLimits(l)
		}
	}
	// check minFee calc against maxFee if set
	if opts.MaxFee > 0 {
		if l := op.Limits(); l.Fee > opts.MaxFee {
			return fmt.Errorf("estimated cost %d > max %d", l.Fee, opts.MaxFee)
		}
	}
	return nil
}

func SignTransaction(op *types.Op, privateKey string, opts *CallOptions) ([]byte, error) {
	key, err := types.ParsePrivateKey(privateKey)
	if err != nil {
		return []byte{}, err
	}
	err = op.Sign(key)
	if err != nil {
		return nil, err
	}
	sig := op.Signature
	op.WithSignature(sig)
	return op.Bytes(), nil
}
