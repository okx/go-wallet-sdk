/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

import (
	"encoding/json"
	"time"
)

var (
	Mainnet    = MustParseChainIdHash("NetXdQprcVkpaWU")
	Ithacanet  = MustParseChainIdHash("NetXnHfVqm9iesp")
	Jakartanet = MustParseChainIdHash("NetXLH1uAxK7CCh")

	// DefaultParams defines the blockchain configuration for Mainnet under the latest
	// protocol.
	DefaultParams = NewParams().ForNetwork(Mainnet).ForProtocol(ProtoV012_2).
			Mixin(&Params{
			OperationTagsVersion:         2,
			MaxOperationsTTL:             120,
			HardGasLimitPerOperation:     1040000,
			HardGasLimitPerBlock:         5200000,
			OriginationSize:              257,
			CostPerByte:                  250,
			HardStorageLimitPerOperation: 60000,
			MinimalBlockDelay:            30 * time.Second,
		})

	IthacanetParams = NewParams().ForNetwork(Ithacanet).ForProtocol(ProtoV012_2).
			Mixin(&Params{
			OperationTagsVersion:         2,
			MaxOperationsTTL:             120,
			HardGasLimitPerOperation:     1040000,
			HardGasLimitPerBlock:         5200000,
			OriginationSize:              257,
			CostPerByte:                  250,
			HardStorageLimitPerOperation: 60000,
			MinimalBlockDelay:            15 * time.Second,
		})

	JakartanetParams = NewParams().ForNetwork(Jakartanet).ForProtocol(ProtoV013_2).
				Mixin(&Params{
			OperationTagsVersion:         2,
			MaxOperationsTTL:             120,
			HardGasLimitPerOperation:     1040000,
			HardGasLimitPerBlock:         5200000,
			OriginationSize:              257,
			CostPerByte:                  250,
			HardStorageLimitPerOperation: 60000,
			MinimalBlockDelay:            15 * time.Second,
		})
)

type Params struct {
	// chain identity, not part of RPC
	Name        string       `json:"name"`
	Network     string       `json:"network,omitempty"`
	Symbol      string       `json:"symbol"`
	Deployment  int          `json:"deployment"`
	Version     int          `json:"version"`
	ChainId     ChainIdHash  `json:"chain_id"`
	Protocol    ProtocolHash `json:"protocol"`
	StartHeight int64        `json:"start_height"`
	EndHeight   int64        `json:"end_height"`
	Decimals    int          `json:"decimals"`
	Token       int64        `json:"units"` // atomic units per token

	// Per-protocol configs
	SecurityDepositRampUpCycles  int64            `json:"security_deposit_ramp_up_cycles"`
	PreservedCycles              int64            `json:"preserved_cycles"`
	BlocksPerCycle               int64            `json:"blocks_per_cycle"`
	BlocksPerCommitment          int64            `json:"blocks_per_commitment"`
	BlocksPerRollSnapshot        int64            `json:"blocks_per_roll_snapshot"`
	BlocksPerVotingPeriod        int64            `json:"blocks_per_voting_period"`
	TimeBetweenBlocks            [2]time.Duration `json:"time_between_blocks"`
	EndorsersPerBlock            int              `json:"endorsers_per_block"`
	HardGasLimitPerOperation     int64            `json:"hard_gas_limit_per_operation"`
	HardGasLimitPerBlock         int64            `json:"hard_gas_limit_per_block"`
	OriginationSize              int64            `json:"origination_size"`
	CostPerByte                  int64            `json:"cost_per_byte"`
	HardStorageLimitPerOperation int64            `json:"hard_storage_limit_per_operation"`
	TestChainDuration            int64            `json:"test_chain_duration"`

	MinimalBlockDelay time.Duration `json:"minimal_block_delay"`

	// extra features to follow protocol upgrades
	MaxOperationsTTL     int64 `json:"max_operations_ttl"`               // in block meta until v011, explicit from v012+
	OperationTagsVersion int   `json:"operation_tags_version,omitempty"` // 1 after v005
	NumVotingPeriods     int   `json:"num_voting_periods,omitempty"`     // 5 after v008, 4 before
	StartBlockOffset     int64 `json:"start_block_offset,omitempty"`     // correct start/end cycle since Granada
	StartCycle           int64 `json:"start_cycle,omitempty"`            // correction since Granada v10
	VoteBlockOffset      int64 `json:"vote_block_offset,omitempty"`      // correction for Edo + Florence Mainnet-only +1 bug
}

func NewParams() *Params {
	return &Params{
		Name:             "Tezos",
		Network:          "",
		Symbol:           "XTZ",
		StartHeight:      -1,
		EndHeight:        -1,
		Decimals:         6,
		Token:            1000000, // initial, changed several times later
		NumVotingPeriods: 4,       // initial, changed once in v008
		MaxOperationsTTL: 60,      // initial, changed once in v011
	}
}

func (p *Params) ForNetwork(net ChainIdHash) *Params {
	pp := &Params{}
	*pp = *p
	pp.ChainId = net
	switch true {
	case Mainnet.Equal(net):
		pp.Network = "Mainnet"
		pp.SecurityDepositRampUpCycles = 64
	case Ithacanet.Equal(net):
		pp.Network = "Ithacanet"
		pp.Version = 11 // starts at Hangzhou
	case Jakartanet.Equal(net):
		pp.Network = "Jakartanet"
		pp.Version = 12 // starts at Ithaca
	default:
		pp.Network = "Sandbox"
	}
	return pp
}

func (p *Params) Mixin(src *Params) *Params {
	buf, _ := json.Marshal(src)
	_ = json.Unmarshal(buf, p)
	return p
}

func (p *Params) ForProtocol(proto ProtocolHash) *Params {
	pp := &Params{}
	*pp = *p
	pp.Protocol = proto
	pp.NumVotingPeriods = 4
	pp.MaxOperationsTTL = 60
	switch true {
	case ProtoBootstrap.Equal(proto):
		// retain version set in ForNetwork()
		pp.StartHeight = 1
		pp.EndHeight = 1

	case ProtoV001.Equal(proto):
		pp.Version = 1
		pp.StartHeight = 2
		pp.EndHeight = 28082

	case ProtoV002.Equal(proto):
		pp.Version = 2
		pp.StartHeight = 28083
		pp.EndHeight = 204761

	case ProtoV003.Equal(proto):
		pp.Version = 3
		pp.StartHeight = 204762
		pp.EndHeight = 458752

	case ProtoV004.Equal(proto): // Athens
		pp.Version = 4
		pp.StartHeight = 458753
		pp.EndHeight = 655360

	case PsBabyM1.Equal(proto): // Babylon
		pp.Version = 5
		pp.OperationTagsVersion = 1
		pp.StartHeight = 655361
		pp.EndHeight = 851968

	case PsCARTHA.Equal(proto): // Carthage
		pp.Version = 6
		pp.OperationTagsVersion = 1
		pp.StartHeight = 851969
		pp.EndHeight = 1212416

	case PsDELPH1.Equal(proto): // Delphi
		pp.Version = 7
		pp.OperationTagsVersion = 1
		// this is extremely hacky!
		pp.StartBlockOffset = 0
		pp.StartCycle = 0
		pp.BlocksPerCycle = 4096
		pp.BlocksPerCommitment = 32
		pp.BlocksPerRollSnapshot = 256
		pp.BlocksPerVotingPeriod = 32768
		pp.EndorsersPerBlock = 32
		pp.StartHeight = 1212417
		pp.EndHeight = 1343488

	case PtEdo2Zk.Equal(proto): // Edo
		pp.Version = 8
		pp.OperationTagsVersion = 1
		pp.NumVotingPeriods = 5
		pp.StartBlockOffset = 1343488
		pp.StartCycle = 328
		pp.VoteBlockOffset = 1
		// this is extremely hacky!
		pp.BlocksPerCycle = 4096
		pp.BlocksPerCommitment = 32
		pp.BlocksPerRollSnapshot = 256
		pp.BlocksPerVotingPeriod = 20480
		pp.EndorsersPerBlock = 32
		pp.StartHeight = 1343489
		pp.EndHeight = 1466367

	case PsFLoren.Equal(proto): // Florence
		pp.Version = 9
		pp.OperationTagsVersion = 1
		pp.NumVotingPeriods = 5
		pp.StartBlockOffset = 1466368
		pp.StartCycle = 358
		pp.VoteBlockOffset = 1 // same as Edo (!!)
		// FIXME: this is extremely hacky!
		pp.BlocksPerCycle = 4096
		pp.BlocksPerCommitment = 32
		pp.BlocksPerRollSnapshot = 256
		pp.BlocksPerVotingPeriod = 20480
		pp.StartHeight = 1466368
		pp.EndHeight = 1589247

	case PtGRANAD.Equal(proto): // Granada
		pp.Version = 10
		pp.OperationTagsVersion = 1
		pp.NumVotingPeriods = 5
		pp.MaxOperationsTTL = 120
		pp.StartBlockOffset = 1589248
		pp.StartCycle = 388
		pp.VoteBlockOffset = 0
		// FIXME: this is extremely hacky!
		pp.BlocksPerCycle = 8192
		pp.BlocksPerCommitment = 64
		pp.BlocksPerRollSnapshot = 512
		pp.BlocksPerVotingPeriod = 40960
		pp.EndorsersPerBlock = 256
		pp.StartHeight = 1589248
		pp.EndHeight = 1916928

	case PtHangz2.Equal(proto): // Hangzhou
		pp.Version = 11
		pp.OperationTagsVersion = 1
		pp.NumVotingPeriods = 5
		pp.MaxOperationsTTL = 120
		if Mainnet.Equal(p.ChainId) {
			pp.StartBlockOffset = 1916928
			pp.StartCycle = 428
			pp.VoteBlockOffset = 0
			// FIXME: this is extremely hacky!
			pp.BlocksPerCycle = 8192
			pp.BlocksPerCommitment = 64
			pp.BlocksPerRollSnapshot = 512
			pp.BlocksPerVotingPeriod = 40960
			pp.EndorsersPerBlock = 256
			pp.StartHeight = 1916929
			pp.EndHeight = 2244608
		}
	case Psithaca.Equal(proto): // Ithaca
		pp.Version = 12
		pp.OperationTagsVersion = 2
		pp.NumVotingPeriods = 5
		pp.MaxOperationsTTL = 120
		if Mainnet.Equal(p.ChainId) {
			pp.StartBlockOffset = 2244608
			pp.StartCycle = 468
			pp.VoteBlockOffset = 0
			// FIXME: this is extremely hacky!
			pp.BlocksPerCycle = 8192
			pp.BlocksPerCommitment = 64
			pp.BlocksPerRollSnapshot = 512
			pp.BlocksPerVotingPeriod = 40960
			pp.EndorsersPerBlock = 0
			pp.StartHeight = 2244609
			pp.EndHeight = -1
		} else if Ithacanet.Equal(p.ChainId) {
			pp.StartBlockOffset = 8192
			pp.StartCycle = 2
			pp.StartHeight = 8192
			pp.EndHeight = -1
		}
	case PtJakart.Equal(proto): // Jakarta
		pp.Version = 13
		pp.OperationTagsVersion = 2
		pp.NumVotingPeriods = 5
		pp.MaxOperationsTTL = 120
		if Mainnet.Equal(p.ChainId) {
			pp.StartBlockOffset = 2490368
			pp.StartCycle = 498
			pp.VoteBlockOffset = 0
			// FIXME: this is extremely hacky!
			pp.BlocksPerCycle = 8192
			pp.BlocksPerCommitment = 64
			pp.BlocksPerRollSnapshot = 512
			pp.BlocksPerVotingPeriod = 40960
			pp.EndorsersPerBlock = 0
			pp.StartHeight = 2490369
			pp.EndHeight = -1
		} else if Jakartanet.Equal(p.ChainId) {
			pp.StartBlockOffset = 8192
			pp.StartCycle = 2
			pp.StartHeight = 8193
			pp.EndHeight = -1
		}
	}
	return pp
}
