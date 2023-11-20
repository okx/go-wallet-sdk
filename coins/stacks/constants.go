package stacks

const (
	SerializeP2PKH  = 0
	SerializeP2SH   = 1
	SerializeP2WPKH = 2
	SerializeP2WSH  = 3
)

const (
	Int               = 0
	Uint              = 1
	Buffer            = 2
	BoolTrue          = 3
	BoolFalse         = 4
	PrincipalStandard = 5
	PrincipalContract = 6
	ResponseOk        = 7
	ResponseErr       = 8
	OptionalNone      = 9
	OptionalSome      = 10
	List              = 11
	Tuple             = 12
	IntASCII          = 13
	IntUTF8           = 14
)

const (
	ADDRESS              = 0
	PRINCIPAL            = 1
	LENGTHPREFIXEDSTRING = 2
	MEMOSTRING           = 3
	ASSETINFO            = 4
	POSTCONDITION        = 5
	PUBLICKEY            = 6
	LENGTHPREFIXEDLIST   = 7
	PAYLOAD              = 8
	MESSAGESIGNATURE     = 9
	TRANSACTIONAUTHFIELD = 10
)

const (
	privatekeybytes1 = 32
	publicKeyBytes   = 33
	signatureBytes   = 65
	c32              = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	constHex         = "0123456789abcdef"
	MaxBufferSize    = 1024 * 4
)

const (
	PostConditionModeAllow = 1
	PostConditionModeDeny  = 2
)

const (
	PostConditionPrincipalIDORIGIN   = 1
	PostConditionPrincipalIDSTANDARD = 2
	PostConditionPrincipalIDCONTRACT = 3
)

const (
	Origin   = 0x01
	Standard = 0x02
	Contract = 0x03
)

const (
	STX         = 0x00
	Fungible    = 0x01
	NonFungible = 0x02
)
