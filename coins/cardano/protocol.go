package cardano

// ProtocolParams is a Cardano Protocol Parameters.
type ProtocolParams struct {
	MinFeeA              Coin
	MinFeeB              Coin
	MaxBlockBodySize     uint
	MaxTxSize            uint
	MaxBlockHeaderSize   uint
	KeyDeposit           Coin
	PoolDeposit          Coin
	MaxEpoch             uint
	NOpt                 uint
	PoolPledgeInfluence  Rational
	ExpansionRate        UnitInterval
	TreasuryGrowthRate   UnitInterval
	D                    UnitInterval
	ExtraEntropy         []byte
	ProtocolVersion      ProtocolVersion
	MinPoolCost          Coin
	CoinsPerUTXOByte     Coin
	CostModels           interface{}
	ExecutionCosts       interface{}
	MaxTxExUnits         interface{}
	MaxBlockTxExUnits    interface{}
	MaxValueSize         uint
	CollateralPercentage uint
	MaxCollateralInputs  uint
}

// ProtocolVersion is the protocol version number.
type ProtocolVersion struct {
	_     struct{} `cbor:"_,toarray"`
	Major uint
	Minor uint
}

var protocolParams = &ProtocolParams{
	MinFeeA:              44,
	MinFeeB:              155381,
	KeyDeposit:           2000000,
	PoolDeposit:          500000000,
	MaxValueSize:         5000,
	MaxTxSize:            16384,
	CoinsPerUTXOByte:     4310,
	CollateralPercentage: 150,
	MaxCollateralInputs:  3,
}
