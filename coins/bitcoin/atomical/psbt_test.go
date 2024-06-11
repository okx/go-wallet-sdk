package atomical

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

// https://mempool.space/testnet/tx/766d32e560fca19f22ee3cf526edd00c5fbfd2f9833b5ff51fb64d3564321833
func TestMergePsbt(t *testing.T) {
	sellerPsbt := "cHNidP8BAN0CAAAAApV3TG/w54SD9fdtU6W5j2cBOlpdj33WuIUbrckXp3w6AQAAAAD/////3Hae3avlOA0ThY1rsXFiaFygFIX3WBHy1cNXJe7cVKwBAAAAAP////8DZQAAAAAAAAAiUSDZM1dVZgQG1pEcxQ2hKvQFtwfyJx6IV+z8hSN6zbZwRxAnAAAAAAAAIlEg2TNXVWYEBtaRHMUNoSr0BbcH8iceiFfs/IUjes22cEeWmRoAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAAAAAAABASulwRoAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAQMEAQAAAAETQA+Y1MmQQtNh33YYoUrZ5l276lwWmwGUtLijAespxz17HQpf3XfA5BE5x9ZJ8o5o+x+KAP1Zkdmwo1KznWj88lcBFyAplEsuIf4iLmoDGOkWHBC0Cn7HA2v7Rl0gCiCXK31XngAAAAAA"
	sp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPsbt)), true)
	assert.NoError(t, err)
	buyPsbt := "cHNidP8BALICAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/////3Hae3avlOA0ThY1rsXFiaFygFIX3WBHy1cNXJe7cVKwBAAAAAP////8CAAAAAAAAAAAiUSDBJU2tQKrqEh5T4M7AmYY2L2zUIieW0ZSMlHkbNSgctRAnAAAAAAAAIlEg2TNXVWYEBtaRHMUNoSr0BbcH8iceiFfs/IUjes22cEcAAAAAAAEBKwAAAAAAAAAAIlEgwSVNrUCq6hIeU+DOwJmGNi9s1CInltGUjJR5GzUoHLUAAQErZQAAAAAAAAAiUSDZM1dVZgQG1pEcxQ2hKvQFtwfyJx6IV+z8hSN6zbZwRwEDBIMAAAABE0EyqMsv2PszUAkidLXerXjgrP1m/gUoTdbfK1xcaxuONjWKReHZH9FCal3gcQg0uYik6ukg4N7gyLwT9Ptb+szagwEXICmUSy4h/iIuagMY6RYcELQKfscDa/tGXSAKIJcrfVeeAAAA"
	bp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(buyPsbt)), true)
	assert.NoError(t, err)
	sp.Inputs[1] = bp.Inputs[1]
	for k, _ := range sp.Inputs {
		err = psbt.Finalize(sp, k)
		assert.NoError(t, err)
	}
	err = psbt.MaybeFinalizeAll(sp)
	assert.NoError(t, err)

	buyerSignedTx, err := psbt.Extract(sp)
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = buyerSignedTx.Serialize(&buf)
	assert.NoError(t, err)
	require.Equal(t, int64(255), bitcoin.GetTxVirtualSize(btcutil.NewTx(buyerSignedTx)))
	require.Equal(t, "0200000000010295774c6ff0e78483f5f76d53a5b98f67013a5a5d8f7dd6b8851badc917a77c3a0100000000ffffffffdc769eddabe5380d13858d6bb17162685ca01485f75811f2d5c35725eedc54ac0100000000ffffffff036500000000000000225120d9335755660406d6911cc50da12af405b707f2271e8857ecfc85237acdb670471027000000000000225120d9335755660406d6911cc50da12af405b707f2271e8857ecfc85237acdb6704796991a0000000000225120d9335755660406d6911cc50da12af405b707f2271e8857ecfc85237acdb6704701400f98d4c99042d361df7618a14ad9e65dbbea5c169b0194b4b8a301eb29c73d7b1d0a5fdd77c0e41139c7d649f28e68fb1f8a00fd5991d9b0a352b39d68fcf257014132a8cb2fd8fb3350092274b5dead78e0acfd66fe05284dd6df2b5c5c6b1b8e36358a45e1d91fd1426a5de0710834b988a4eae920e0dee0c8bc13f4fb5bfaccda8300000000", hex.EncodeToString(buf.Bytes()))
}

func TestPsbt(t *testing.T) {
	network := &chaincfg.TestNet3Params
	// seller
	txInput := &TxInput{
		TxId:       "ac54dcee2557c3d5f21158f78514a05c686271b16b8d85130d38e5abdd9e76dc",
		VOut:       uint32(1),
		Amount:     int64(101),
		Address:    "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	}

	txOutput := &TxOutput{
		Address: "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		Amount:  int64(10000),
	}

	sellerPsbt, err := GenerateAtomicalSignedListingPSBTBase64(txInput, txOutput, network)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "cHNidP8BALICAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/////3Hae3avlOA0ThY1rsXFiaFygFIX3WBHy1cNXJe7cVKwBAAAAAP////8CAAAAAAAAAAAiUSAqfZ5VuWBDpp6/Ye7m7xjbUgvubGLNU95k/nksx3r2ihAnAAAAAAAAIlEg2TNXVWYEBtaRHMUNoSr0BbcH8iceiFfs/IUjes22cEcAAAAAAAEBKwAAAAAAAAAAIlEgKn2eVblgQ6aev2Hu5u8Y21IL7mxizVPeZP55LMd69ooAAQErZQAAAAAAAAAiUSDZM1dVZgQG1pEcxQ2hKvQFtwfyJx6IV+z8hSN6zbZwRwEDBIMAAAABE0DfnGn8QAfXzFUX9y9BcoMt1BUzA/EuOkFTHuecLJ2ibXTta1rTHNStalxaj1tz06no1htcBb0Pc76NLoS1zHsLARcgV7uy1KnLiiNXYz8gG5xRjCeV3taCt5E8a+7z/iO9bS8AAAA=", sellerPsbt)
	// buyer
	var inputs []*TxInput
	inputs = append(inputs, &TxInput{
		TxId:       "3a7ca717c9ad1b85b8d67d8f5d5a3a01678fb9a5536df7f58384e7f06f4c7795",
		VOut:       1,
		Amount:     1753509,
		Address:    "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})

	// seller input
	inputs = append(inputs, txInput)
	var outputs []*TxOutput
	outputs = append(outputs, &TxOutput{
		Address: "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		Amount:  int64(101),
	})
	// seller output
	outputs = append(outputs, txOutput)
	outputs = append(outputs, &TxOutput{
		Address: "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		Amount:  int64(0),
	})

	fee, buyerTx, err := GenerateAtomicalSignedBuyingTx(inputs, outputs, 1, 1, sellerPsbt, network)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, int64(255), fee)
	require.Equal(t, "cHNidP8BAN0CAAAAApV3TG/w54SD9fdtU6W5j2cBOlpdj33WuIUbrckXp3w6AQAAAAD/////3Hae3avlOA0ThY1rsXFiaFygFIX3WBHy1cNXJe7cVKwBAAAAAP////8DZQAAAAAAAAAiUSDZM1dVZgQG1pEcxQ2hKvQFtwfyJx6IV+z8hSN6zbZwRxAnAAAAAAAAIlEg2TNXVWYEBtaRHMUNoSr0BbcH8iceiFfs/IUjes22cEeWmRoAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAAAAAAABASulwRoAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAQMEAQAAAAETQMmEMSawuRaa8vcEgKut3e0rIxZbCs7HD8w8D8Vz+Q6WkIXpfmXbKB+Ky4bZgZT6KxAfd+yM/l9VuMoH55C0Gs4BFyBXu7LUqcuKI1djPyAbnFGMJ5Xe1oK3kTxr7vP+I71tLwABAStlAAAAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAQMEgwAAAAETQN+cafxAB9fMVRf3L0Fygy3UFTMD8S46QVMe55wsnaJtdO1rWtMc1K1qXFqPW3PTqejWG1wFvQ9zvo0uhLXMewsBFyBXu7LUqcuKI1djPyAbnFGMJ5Xe1oK3kTxr7vP+I71tLwAAAAA=", buyerTx)
}
