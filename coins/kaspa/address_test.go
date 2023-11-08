package kaspa

import (
	"github.com/kaspanet/kaspad/domain/dagconfig"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAddress(t *testing.T) {
	privateKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	expectAddress := "kaspadev:qrsq5qug9fq8afkkgvm2rg2c00j2qcx2sthpsg2eap67269uv2l9s9gw0a397"
	privateKeyHex2 := "790f2b826ad9dfa7f2a53ec68e37ea51dc58652ecfde812da37c96a1069fcdbb"
	expectAddress2 := "kaspa:qqqzs9jqks8ljamhefyd4d33jxmr7p3r7qve37cryduzvqyls4ffu5mherkjv"

	address, err := NewAddressWithNetParams(privateKeyHex, dagconfig.DevnetParams)
	require.Nil(t, err)
	require.Equal(t, expectAddress, address)
	require.True(t, ValidateAddressWithNetParams(address, dagconfig.DevnetParams))
	require.False(t, ValidateAddressWithNetParams(address, dagconfig.MainnetParams))
	address2, err := NewAddressWithNetParams(privateKeyHex2, dagconfig.MainnetParams)
	require.Nil(t, err)
	require.Equal(t, expectAddress2, address2)
}
