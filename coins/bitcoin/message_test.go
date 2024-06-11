package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignBip0322(t *testing.T) {
	sig, err := SignBip0322("Hello World", "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22")
	require.NoError(t, err)
	require.Equal(t, "AkgwRQIhAM0XBsjxc8MH0n0l9NBaqpv0yVgb1zRk5zgrpAw4Y+mxAiAtA0OdpiuxuQ6U9itEeRZZO2hAxeXSo1gBp1gXEXFLRwEhA1e7stSpy4ojV2M/IBucUYwnld7WgreRPGvu8/4jvW0v", sig)
}

func TestSignBip0322TapRoot(t *testing.T) {
	sig, err := SignBip0322("Hello World", "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr", "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22")
	require.NoError(t, err)
	require.Equal(t, "AUA62WElIGyeXycCIIuyOgB9sn/Y7Jjk0yDfhu83qWCuEO6wib+ScHrpm/GilVZPWnVyI+i3r0RDZ0L3qkEmiCyy", sig)
}

func TestBuildToSpend(t *testing.T) {
	network := &chaincfg.MainNetParams
	txId, err := BuildToSpend("Hello World", "bc1q9vza2e8x573nczrlzms0wvx3gsqjx7vavgkx0l", network)
	require.NoError(t, err)
	require.Equal(t, txId, "b79d196740ad5217771c1098fc4a4b51e0535c32236c71f1ea4d61a2d603352b")
	txId, err = BuildToSpend("", "bc1q9vza2e8x573nczrlzms0wvx3gsqjx7vavgkx0l", network)
	require.NoError(t, err)
	require.Equal(t, txId, "c5680aa69bb8d860bf82d4e9cd3504b55dde018de765a91bb566283c545a99a7")
}

func TestBuildToSpendTapRoot(t *testing.T) {
	network := &chaincfg.MainNetParams

	txId, err := BuildToSpend("Hello World", "bc1ppv609nr0vr25u07u95waq5lucwfm6tde4nydujnu8npg4q75mr5sxq8lt3", network)
	require.NoError(t, err)
	require.Equal(t, txId, string(reverseBytes([]byte("0679db23166a7ca5a37998ba7836c33198bba97552657b128db108d29f6e6621"))))

	txId, err = BuildToSpend("", "bc1ppv609nr0vr25u07u95waq5lucwfm6tde4nydujnu8npg4q75mr5sxq8lt3", network)
	require.NoError(t, err)
	require.Equal(t, txId, string(reverseBytes([]byte("86273d66a2c1c682748f0cfa4d8bfa15e4ca1ef7d5add663b1aa2ae612bd9d72"))))
}

func TestBip0322Hash(t *testing.T) {
	hash1 := Bip0322Hash("hello world")
	require.Equal(t, "3467b0020ea3500b767790e76c108d412e655c905a3b9fa7dcc2026134ec2156", hash1)

	hash2 := Bip0322Hash("")
	require.Equal(t, "c90c269c4f8fcbe6880f72a721ddfbf1914268a794cbb21cfafee13770ae19f1", hash2)

	hash3 := Bip0322Hash("Hello World")
	require.Equal(t, "f0eb03b1a75ac6d9847f55c624a99169b5dccba2a31f5b23bea77ba270de0a7a", hash3)
}

func reverseBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i += 2 {
		j := len(b) - i - 2
		b[i], b[j] = b[j], b[i]
		b[i+1], b[j+1] = b[j+1], b[i+1]
	}
	return b
}

func TestMPCUnsignBip0322(t *testing.T) {
	res, err := MPCUnsignedBip0322("Hello World", "tb1qsrx5k69d92avuf4fwke8k35flywd930ems48zf", "031cf908e7712d7a1c4cee9d18c41309fbc750ca47fde5e26a52704ea7fa196a50", &chaincfg.TestNet3Params)
	require.NoError(t, err)

	require.Equal(t, "70736274ff01003d00000000012eb1ba16cbc7659ece4a151926de839f94e3e988c73ae6290af9d668f3b61061000000000000000000010000000000000000016a000000000001011f000000000000000016001480cd4b68ad2abace26a975b27b4689f91cd2c5f9010304010000000000", res.PsbtTx)
	require.Equal(t, "885d4019cae7c869d7c85567bcc8808eb8068c49d23cf1b8332ab3ba55d6c5b8", res.SignHashList[0])
}
func TestMPCSignedBip0322(t *testing.T) {
	res, err := MPCSignedBip0322("Hello World", "tb1qsrx5k69d92avuf4fwke8k35flywd930ems48zf", "031cf908e7712d7a1c4cee9d18c41309fbc750ca47fde5e26a52704ea7fa196a50", []string{"96a8a365b9502ef1f2322b6ef38c058b3d20bff26e5dab17f1c2271d643317941f116006cedc765196c6120fc835054ed6806e71c7014613579d3e04ebb8fd47"}, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	require.Equal(t, "AkgwRQIhAJaoo2W5UC7x8jIrbvOMBYs9IL/ybl2rF/HCJx1kMxeUAiAfEWAGztx2UZbGEg/INQVO1oBucccBRhNXnT4E67j9RwEhAxz5COdxLXocTO6dGMQTCfvHUMpH/eXialJwTqf6GWpQ", res.PsbtTx)
}

func TestSignMessage(t *testing.T) {
	wif := "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22"
	message := "Hello World"
	res, err := SignMessage(wif, "", message)
	require.NoError(t, err)
	require.Equal(t, "INPqXlpr5h2xvRzKk4BLCcMNgCm5Xv062zebU8JN6EMAdJegfIVF//gpna+DBR+zQztxi/d/WNFZ6QRxQDWMEAo=", res)
}

func TestMPCUnsignedMessage(t *testing.T) {
	message := "Hello World"
	hash := MPCUnsignedMessage("", message)
	require.Equal(t, "a7af0baad5ae99b97fc69b3a0d1abcf3ef17f131cc4776e1bc11933ec8550f49", hash)
}

func TestMPCSignedMessage(t *testing.T) {
	signature := "eb63e75200f97eae94ff4698f42c07de600ea164d5a96ce0594c6cd9f4cfcefe4a949e137b44fd06a2411a0f0b0db8c457b364b03c327c69e3976b5086a2dcb920"
	publicKeyHex := "037cd7ace531e991ea00d1e63f23e2e7e1397606a5483ac11308ae37e9d4e9843f"
	res, err := MPCSignedMessage(signature, publicKeyHex, nil)
	require.NoError(t, err)
	require.Equal(t, "IOtj51IA+X6ulP9GmPQsB95gDqFk1als4FlMbNn0z87+SpSeE3tE/QaiQRoPCw24xFezZLA8Mnxp45drUIai3Lk=", res)
}

func TestVerifyMessage(t *testing.T) {
	signatureB64 := "HzS54lJsvripU/LXSGrDTEsv47zy2S5M/UD2F+jRRwv4QNuncyUHraa9FAio8SZJfHM4sqKw/khRBwUkhpGLloI="
	message := "Hello World"
	publicKeyHex := "024f1bd355a61ec33cabdb7251f050fcbd922bbd0fae48743fe925b0b324493c77"
	err := VerifyMessage(signatureB64, "", message, publicKeyHex, "1LrCJN5FVSNinDvqYtRHeEVnf6Dt5e8HUz", "", &chaincfg.MainNetParams)
	require.NoError(t, err)
}

func TestMPCSignedMessageCompat(t *testing.T) {
	message := "Hello World"
	signature := "34b9e2526cbeb8a953f2d7486ac34c4b2fe3bcf2d92e4cfd40f617e8d1470bf840dba7732507ada6bd1408a8f126497c7338b2a2b0fe485107052486918b9682"
	publicKeyHex := "024f1bd355a61ec33cabdb7251f050fcbd922bbd0fae48743fe925b0b324493c77"
	res, err := MPCSignedMessageCompat("", message, signature, publicKeyHex, nil)
	require.NoError(t, err)
	require.Equal(t, "HzS54lJsvripU/LXSGrDTEsv47zy2S5M/UD2F+jRRwv4QNuncyUHraa9FAio8SZJfHM4sqKw/khRBwUkhpGLloI=", res)
}

func TestVerifySimpleForBip0322(t *testing.T) {
	// test for taproot
	signatureB64 := "AUD5MwxtURP3tAip3fS5vVRwa4L15wEyTIG0BQ3DPktJpXvQe7Sh8kf+mVaO4ldEP+vhiVZ/sXvOHEbQQnsiYpCq"
	address := "bc1ppv609nr0vr25u07u95waq5lucwfm6tde4nydujnu8npg4q75mr5sxq8lt3"
	message := "Hello World"
	publicKeyHex := "02c7f12003196442943d8588e01aee840423cc54fc1521526a3b85c2b0cbd58872"
	err := VerifySimpleForBip0322(message, address, signatureB64, publicKeyHex, &chaincfg.MainNetParams)
	require.NoError(t, err)

	// test for segwit
	signatureB64 = "AkgwRQIhAOzyynlqt93lOKJr+wmmxIens//zPzl9tqIOua93wO6MAiBi5n5EyAcPScOjf1lAqIUIQtr3zKNeavYabHyR8eGhowEhAsfxIAMZZEKUPYWI4BruhAQjzFT8FSFSajuFwrDL1Yhy"
	address = "bc1q9vza2e8x573nczrlzms0wvx3gsqjx7vavgkx0l"
	message = "Hello World"
	publicKeyHex = "02c7f12003196442943d8588e01aee840423cc54fc1521526a3b85c2b0cbd58872"
	err = VerifySimpleForBip0322(message, address, signatureB64, publicKeyHex, &chaincfg.MainNetParams)
	require.NoError(t, err)
}
