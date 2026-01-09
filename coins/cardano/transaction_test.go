package cardano

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {
	mnemonic := "north bulb crunch need badge orient tissue web east scan invite energy canal solar eight"
	path := "m/1852'/1815'/0'/0/0"

	prvKey, err := DerivePrvKey(mnemonic, path)
	assert.Nil(t, err)

	txData := &TxData{
		PrvKey: prvKey,
		Inputs: []*TxIn{
			{
				TxId:   "7f6a09b3eb7ea3942b788c7aa086a43124021136f9ea4afe9ac705bc28e0cf17",
				Index:  1,
				Amount: 1000000000,
			},
			{
				TxId:   "7f6a09b3eb7ea3942b788c7aa086a43124021136f9ea4afe9ac705bc28e0cf17",
				Index:  2,
				Amount: 1150770,
				MultiAsset: []*MultiAssets{
					{
						PolicyId: "29d222ce763455e3d7a09a665ce554f00ac89d2e99a1a83d267170c6",
						Assets: []*Asset{
							{
								AssetName: "4d494e",
								Amount:    11614448,
							},
						},
					},
				},
			},
			{
				TxId:   "f2b78093ca7be37d24a8f6462991745552f80f4610d1777c456a7ce24f2b3e02",
				Index:  1,
				Amount: 2000000,
				MultiAsset: []*MultiAssets{
					{
						PolicyId: "29d222ce763455e3d7a09a665ce554f00ac89d2e99a1a83d267170c6",
						Assets: []*Asset{
							{
								AssetName: "4d494e",
								Amount:    1741609,
							},
						},
					},
				},
			},
		},
		ToAddress: "addr1qxpuwk209ey9t2pqmd79tg05kyyykxsjc5j0eq6p56aru8ua7qwcmquymkmt749uuz8e60s642ynhjvvrlk0ldgjez4qxpfydr",
		Amount:    1150770,
		MultiAsset: []*MultiAssets{
			{
				PolicyId: "29d222ce763455e3d7a09a665ce554f00ac89d2e99a1a83d267170c6",
				Assets: []*Asset{
					{
						AssetName: "4d494e",
						Amount:    100000,
					},
				},
			},
		},
		ChangeAddress: "addr1qyxdgpwcqsrfsfv7gs3el47ym205hxaxnnpvs550czrjr8gr7z40zns2zm4kdd5jgxhawpstcgnyt4zdwzn4e9g6qmksvhsufu",
		TTL:           999999999,
	}

	minAda, err := MinAda(txData.ToAddress, txData.MultiAsset)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1150770), minAda)

	minFee, err := MinFee(txData)
	assert.Nil(t, err)
	assert.True(t, minFee.Valid)
	assert.Equal(t, uint64(177865), minFee.Fee)
	assert.Equal(t, uint64(1150770), minFee.Change)

	tx, err := Transfer(txData)
	assert.Nil(t, err)
	expected := "hKQAg4JYIH9qCbPrfqOUK3iMeqCGpDEkAhE2+epK/prHBbwo4M8XAYJYIH9qCbPrfqOUK3iMeqCGpDEkAhE2+epK/prHBbwo4M8XAoJYIPK3gJPKe+N9JKj2RimRdFVS+A9GENF3fEVqfOJPKz4CAQGDglg5AYPHWU8uSFWoINt8VaH0sQhLGhLFJPyDQaa6Ph+d8B2Ng4Tdtr9UvOCPnT4aqok7yYwf7P+1EsiqghoAEY8yoVgcKdIiznY0VePXoJpmXOVU8ArInS6Zoag9JnFwxqFDTUlOGgABhqCCWDkBDNQF2AQGmCWeRCOf18Tan0ubppzCyFKPwIchnQPwqvFOChbrZraSQa/XBgvCJkXUTXCnXJUaBu2CGgARjzKhWBwp0iLOdjRV49egmmZc5VTwCsidLpmhqD0mcXDGoUNNSU4aAMpFeYJYOQEM1AXYBAaYJZ5EI5/XxNqfS5umnMLIUo/AhyGdA/Cq8U4KFutmtpJBr9cGC8ImRdRNcKdclRoG7Ro7pQiFAhoAArbJAxo7msn/oQCBglggKDkK4+UZY8yTWs1bUiodu7bgYRmuOp2TFDrNMLyLKlNYQMx7FYmNxn2xsBaS8QQif7hjW47Kdi5YYofoE0TKY543CBQ138TZgExyrGNLGvCKyxHVfybkZhBjOFgLbdmNAw319g=="
	assert.Equal(t, expected, tx)

	txid, err := GetTxHash(tx)
	assert.Nil(t, err)
	assert.Equal(t, "426ff2e325d04baa9e7e65e622c13447f3852ad8c70897cf735ed20bfb5a0a6d", txid)
}

func TestMinAda(t *testing.T) {
	address := "addr1qxpuwk209ey9t2pqmd79tg05kyyykxsjc5j0eq6p56aru8ua7qwcmquymkmt749uuz8e60s642ynhjvvrlk0ldgjez4qxpfydr"
	multiAsset := []*MultiAssets{
		{
			PolicyId: "29d222ce763455e3d7a09a665ce554f00ac89d2e99a1a83d267170c6",
			Assets: []*Asset{
				{
					AssetName: "4d494e",
					Amount:    1000000,
				},
			},
		},
	}

	minAda, err := MinAda(address, multiAsset)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1150770), minAda)
}

func TestMinFee(t *testing.T) {
	txData := &TxData{
		Inputs: []*TxIn{
			{
				TxId:   "7f6a09b3eb7ea3942b788c7aa086a43124021136f9ea4afe9ac705bc28e0cf17",
				Index:  2,
				Amount: 1150770,
				MultiAsset: []*MultiAssets{
					{
						PolicyId: "29d222ce763455e3d7a09a665ce554f00ac89d2e99a1a83d267170c6",
						Assets: []*Asset{
							{
								AssetName: "4d494e",
								Amount:    11614448,
							},
						},
					},
				},
			},
			{
				TxId:   "f2b78093ca7be37d24a8f6462991745552f80f4610d1777c456a7ce24f2b3e02",
				Index:  1,
				Amount: 2000000,
				MultiAsset: []*MultiAssets{
					{
						PolicyId: "29d222ce763455e3d7a09a665ce554f00ac89d2e99a1a83d267170c6",
						Assets: []*Asset{
							{
								AssetName: "4d494e",
								Amount:    1741609,
							},
						},
					},
				},
			},
		},
		ToAddress: "addr1qyxdgpwcqsrfsfv7gs3el47ym205hxaxnnpvs550czrjr8gr7z40zns2zm4kdd5jgxhawpstcgnyt4zdwzn4e9g6qmksvhsufu",
		Amount:    1150770,
		MultiAsset: []*MultiAssets{
			{
				PolicyId: "29d222ce763455e3d7a09a665ce554f00ac89d2e99a1a83d267170c6",
				Assets: []*Asset{
					{
						AssetName: "4d494e",
						Amount:    1000000,
					},
				},
			},
		},
		ChangeAddress: "addr1qyxdgpwcqsrfsfv7gs3el47ym205hxaxnnpvs550czrjr8gr7z40zns2zm4kdd5jgxhawpstcgnyt4zdwzn4e9g6qmksvhsufu",
		TTL:           999999999,
	}

	minFee, err := MinFee(txData)
	assert.Nil(t, err)
	assert.True(t, minFee.Valid)
	assert.Equal(t, uint64(173421), minFee.Fee)
	assert.Equal(t, uint64(1150770), minFee.Change)
}
