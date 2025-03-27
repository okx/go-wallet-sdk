package strkey

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaxEncodedSize(t *testing.T) {
	assert.Equal(t, encoding.EncodedLen(maxRawSize), maxEncodedSize)
}

func TestVersion(t *testing.T) {
	cases := []struct {
		Name                string
		Address             string
		ExpectedVersionByte VersionByte
	}{
		{
			Name:                "AccountID",
			Address:             "GA3D5KRYM6CB7OWQ6TWYRR3Z4T7GNZLKERYNZGGA5SOAOPIFY6YQHES5",
			ExpectedVersionByte: VersionByteAccountID,
		},
		{
			Name:                "Seed",
			Address:             "SBU2RRGLXH3E5CQHTD3ODLDF2BWDCYUSSBLLZ5GNW7JXHDIYKXZWHOKR",
			ExpectedVersionByte: VersionByteSeed,
		},
		{
			Name:                "HashTx",
			Address:             "TBU2RRGLXH3E5CQHTD3ODLDF2BWDCYUSSBLLZ5GNW7JXHDIYKXZWHXL7",
			ExpectedVersionByte: VersionByteHashTx,
		},
		{
			Name:                "HashX",
			Address:             "XBU2RRGLXH3E5CQHTD3ODLDF2BWDCYUSSBLLZ5GNW7JXHDIYKXZWGTOG",
			ExpectedVersionByte: VersionByteHashX,
		},
		{
			Name:                "Other (0x60)",
			Address:             "MBU2RRGLXH3E5CQHTD3ODLDF2BWDCYUSSBLLZ5GNW7JXHDIYKXZWGTOG",
			ExpectedVersionByte: VersionByte(0x60),
		},
		{
			Name:                "Signed Payload",
			Address:             "PDPYP7E6NEYZSVOTV6M23OFM2XRIMPDUJABHGHHH2Y67X7JL25GW6AAAAAAAAAAAAAAJEVA",
			ExpectedVersionByte: VersionByteSignedPayload,
		},
		{
			Name:                "Contract",
			Address:             "CA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUWDA",
			ExpectedVersionByte: VersionByteContract,
		},
	}

	for _, kase := range cases {
		actual, err := Version(kase.Address)
		if assert.NoError(t, err, "An error occurred decoding case %s", kase.Name) {
			assert.Equal(t, kase.ExpectedVersionByte, actual, "Output mismatch in case %s", kase.Name)
		}
	}
}

func TestIsValidEd25519PublicKey(t *testing.T) {
	validKey := "GDWZCOEQRODFCH6ISYQPWY67L3ULLWS5ISXYYL5GH43W7YFMTLB65PYM"
	isValid := IsValidEd25519PublicKey(validKey)
	assert.Equal(t, true, isValid)

	invalidKey := "GDWZCOEQRODFCH6ISYQPWY67L3ULLWS5ISXYYL5GH43W7Y"
	isValid = IsValidEd25519PublicKey(invalidKey)
	assert.Equal(t, false, isValid)

	invalidKey = ""
	isValid = IsValidEd25519PublicKey(invalidKey)
	assert.Equal(t, false, isValid)

	invalidKey = "SBCVMMCBEDB64TVJZFYJOJAERZC4YVVUOE6SYR2Y76CBTENGUSGWRRVO"
	isValid = IsValidEd25519PublicKey(invalidKey)
	assert.Equal(t, false, isValid)

	invalidKey = "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26"
	isValid = IsValidEd25519PublicKey(invalidKey)
	assert.False(t, isValid)
}

func TestIsValidMuxedAccountEd25519PublicKey(t *testing.T) {
	validKey := "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26"
	isValid := IsValidMuxedAccountEd25519PublicKey(validKey)
	assert.True(t, isValid)

	invalidKeys := []struct {
		key    string
		reason string
	}{
		{
			key:    "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUAAAAAAAAAAAACJUR",
			reason: "The unused trailing bit must be zero in the encoding of the last three bytes (24 bits) as five base-32 symbols (25 bits)",
		},
		{
			key:    "GA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVSGZA",
			reason: "Invalid length (congruent to 1 mod 8)",
		},
		{
			key:    "G47QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVP2I",
			reason: "Invalid algorithm (low 3 bits of version byte are 7)",
		},
		{
			key:    "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLKA",
			reason: "Invalid length (congruent to 6 mod 8)",
		},
		{
			key:    "M47QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUAAAAAAAAAAAACJUQ",
			reason: "Invalid algorithm (low 3 bits of version byte are 7)",
		},
		{
			key:    "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUAAAAAAAAAAAACJUK",
			reason: "Padding bytes are not allowed",
		},
		{
			key:    "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUAAAAAAAAAAAACJUO",
			reason: "Invalid checksum",
		},
		{
			key:    "",
			reason: "Invalid length (string is empty)",
		},
		{
			key:    "SBCVMMCBEDB64TVJZFYJOJAERZC4YVVUOE6SYR2Y76CBTENGUSGWRRVO",
			reason: "Invalid key (this is a secret key)",
		},
		{
			key:    "GDWZCOEQRODFCH6ISYQPWY67L3ULLWS5ISXYYL5GH43W7YFMTLB65PYM",
			reason: "Invalid key (this is an Ed25519 G-address)",
		},
	}
	for _, invalidKey := range invalidKeys {
		isValid = IsValidMuxedAccountEd25519PublicKey(invalidKey.key)
		assert.False(t, isValid, invalidKey.reason)
	}
}

func TestIsValidEd25519SecretSeed(t *testing.T) {
	validKey := "SBCVMMCBEDB64TVJZFYJOJAERZC4YVVUOE6SYR2Y76CBTENGUSGWRRVO"
	isValid := IsValidEd25519SecretSeed(validKey)
	assert.Equal(t, true, isValid)

	invalidKey := "SBCVMMCBEDB64TVJZFYJOJAERZC4YVVUOE6SYR2Y76CBTENGUSG"
	isValid = IsValidEd25519SecretSeed(invalidKey)
	assert.Equal(t, false, isValid)

	invalidKey = ""
	isValid = IsValidEd25519SecretSeed(invalidKey)
	assert.Equal(t, false, isValid)

	invalidKey = "GDWZCOEQRODFCH6ISYQPWY67L3ULLWS5ISXYYL5GH43W7YFMTLB65PYM"
	isValid = IsValidEd25519SecretSeed(invalidKey)
	assert.Equal(t, false, isValid)

	invalidKey = "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26"
	isValid = IsValidEd25519SecretSeed(invalidKey)
	assert.False(t, isValid)
}
