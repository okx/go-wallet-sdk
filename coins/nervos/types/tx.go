package types

const (
	SECP256K1_BLAKE160_SIGHASH_ALL_DATA_HASH  = "0x973bdb373cbb1d752b4ac006e2bb5bdcb63431ed2b6e394b22721c8906a2ad72"
	SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH  = "0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"
	SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH = "0x5c5069eb0857efc65e1bca0c07df34c31663b3622fd3876c876320fc9634e2a8"
)

type SystemScriptCell struct {
	CellHash Hash
	OutPoint *OutPoint
	HashType ScriptHashType
	DepType  DepType
}

type SystemScripts struct {
	SecpSingleSigCell *SystemScriptCell
	SecpMultiSigCell  *SystemScriptCell
	DaoCell           *SystemScriptCell
	ACPCell           *SystemScriptCell
	SUDTCell          *SystemScriptCell
	ChequeCell        *SystemScriptCell
}

func NewSystemScripts(chain string) *SystemScripts {
	return &SystemScripts{
		SecpSingleSigCell: secpSingleSigCell(chain),
	}
}

func secpSingleSigCell(chain string) *SystemScriptCell {
	if chain == "ckb" {
		return &SystemScriptCell{
			CellHash: HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			OutPoint: &OutPoint{
				TxHash: HexToHash("0x71a7ba8fc96349fea0ed3a5c47992e3b4084b031a42264a018e0072e8172e46c"),
				Index:  0,
			},
			HashType: HashTypeType,
			DepType:  DepTypeDepGroup,
		}
	} else {
		return &SystemScriptCell{
			CellHash: HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			OutPoint: &OutPoint{
				TxHash: HexToHash("0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37"),
				Index:  0,
			},
			HashType: HashTypeType,
			DepType:  DepTypeDepGroup,
		}
	}
}
