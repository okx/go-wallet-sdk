package solana

import (
	"testing"

	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/stretchr/testify/assert"
)

func TestNewAddressFromPubkey(t *testing.T) {
	addr := "2uWejjxZtzuqLrQeCH4gwh3C5TNn2rhHTdvC26dWzKfM"
	pubkey := base58.Decode(addr)

	addr1, err := NewAddressFromPubkey(pubkey)
	assert.NoError(t, err)
	assert.Equal(t, addr, addr1)

	invalidPubkey := pubkey[:31]
	addr2, err := NewAddressFromPubkey(invalidPubkey)
	assert.Error(t, err, assert.AnError)
	assert.Equal(t, "", addr2)
}

func TestNewTxFromRaw(t *testing.T) {
	rawTx := "3SabRqqxKTdwNqsoa2P459eZX6RpFhonzCpUbjvv6oKj9yFpQ9D8URYW6GffzxfFgYzZg4f5ujbdZH5XigVpVzhQmzshWyFuoGFfSRZszHwzUpWA2Cxb2ugHfKwhvu8Ma7C67Q7wXQta8TpA2WfkaNieWHiC3YNMgW338StjVQnUXuoKQRofP5aYikTc3V6nAR1jGpwqwtopr5HFwddfrKqLzGJivuGc4oXZyfWJM6DH8BPzTWAtvvr8e1HZh4Bf4U1FWUnPa3WPoZvGV6RGTEsTGMYgrGRQ7FHaMefp147fnLH8HgWVWYSjib31kDa9MjSqXV9vvdcGXihNNogwRGC1C46b7z6ozy6WTzsxhr3oVMJ9PCzGx26EvZ96fheZpY7QRxbFuc2gQyzV5mM8qymt4Xu9aqjip8Nsxx8kEXY2cGJP9o1NWHc2EuvkwkJPrddMH7xN3U4xvwpCWKCy4sjYyLrzp1n6iuVBjXGJm74rJCxFDitSPZVh2cnFNHPwrR8eNBJA9bZK4d9Gbuiv2yjvzzfip8qF3S7ASS25rXXXMfvYAQ9Z4KM6LBbgksxU3VE7FBTsZ58cNN2VqsQe6b1h9FQFzByydYXKBJF3h2MGxFEvi1iKsurN6nrvTisK1qK5M5TmfrMyafvowuHE4xSpQ1tL6qjBjKovfBgX1jbbNNNTwQgPoS4cgU9CUihbTyUjWNdPrpPdDW7Sc762VdbwcNCAotDC8QMkWGQpiPumK5UMNgUgn1YSA2yrC5zjLzis7dFcTJC2oMSfyzcZsYqaFkHffe266w9rrbNqdz2azBpmfbG3VcDTCwKM7JDnBEjfAtSULdSfDMRJZx1wZYmv9ujFVeagfvmsUGAxeB1F29guxbKf3NYcKf4s4vMp95GNhtLu3VGc8AtBKGkm4bEdbxUcrLkNNEeumt6RUF9srQtf8Dw6J4AFtDQSnRvjCctmTiBAez93to2z3suvx45YZZNmmUsKbceekYYQoPe5j8hi5QLmXHwoePR3pX5sdGcwEYyzkPtW2ys5ithcLzhnBt2mZNo8gz69eyagaUat8dx6p7HEKMpok9FUEtsWvmDug6giH1mPgrS"
	tx, err := NewTxFromRaw(rawTx, "base58")
	assert.NoError(t, err)
	assert.Equal(t, string(tx.Message.Version), "v0")
	assert.Equal(t, tx.Message.RecentBlockHash, "CwT5Apc8MiSE136oEwAuFUxQypPQJapTo2TTtJ3TsLS4")

	rawTxBase64 := "AbhcR6grmy8wUqM5IjMeXjUssEV+H8j8WBaFiZEsALabE5CJpg8hdYPq2w81OMyPXB7V7VCtT9xYLJNadWL5mwaAAQAEDQb9t482B6Dp4eZmxmTwVerNKk+OwaSUZehZAvDgo3fo/Qmown5STyyzJiRldhIYUPhMmlh1aRghqofs9rpIfKNT0gyZOo8cmDgVSFx2rwETLsQTYwOxYg2Hni/4ArsTQN8SXv0Pxxy0L+u/AEQUSB9h8rR9wWyaR0HLBr84TUM24U7pDQJ3qLATOdN8ojKFnYTfpL0skOzD/vGIo14EesXKX6i3/hwYl0QnQ8SOvoaW9lSvg/Uf0GsFL1K8qMUs75BBI2CMM/ENcFbVoFL6MFBC0G4Vtvld9KTbbbbwzsRkT9g6gerxK/aJl1o85kPnAi2Vy4LJdqxASFsSm7N0eNWgYAJEx6u4cFcpvamZ4gg9vTn/LfzgEYLt4fLV7NsY11ZbeLrsW9b/BmMzGOog5/Y5jS8ygOj3+MOATiiOeKeNAwZGb+UhFzL/7K26csOb57yM5bvF9xJrLEObOkAAAACMlyWPTiSJ8bs9ECkUjg2DC1oTmdr/EIQEjnvY2+n4WQxCb1EGGE1ur77aVZO2N6GaoIOvyfwlGCpaRjT9SzPssWNmQqDAIu3Cx57nYNxoPXQHqJhIOPavZDE/16qBrckFCgAFAhCGBQAKAAkDAAAAAAAAAAALBgABAAkXGAEBCwYADxYdFxgBAQwoAAIBHQkPGBoAAg4SHhUUExEdHBgYGRANDgEDEAQFCRwGEBgYGwcIF1Wtg04mlqV7D6wmAAAAAAAAggMAAAAAAAApAwAAAAAAAAEAAACsJgAAAAAAAAEAAAACAAAAAQAAAAwBAAAAZAEAAAANAQAAAGSAlpgAAQAAAAAAAAAAA2m+hY167BMiFYkA29j+OxYOg7Z/De2IodldI9hoQqA/AgECAMRI+LuCwsD1DXxJ6Vq20Ph+U32vHeJRXh6AKgQLlh0OAgNZCAFDRUdYWnt85/+Pn+x58pofBeEA06ZfJD37btWPR0TdRQ5ls3dN5ZEFYFtcXV4BXw=="
	tx, err = NewTxFromRaw(rawTxBase64, "base64")
	assert.NoError(t, err)
	assert.Equal(t, string(tx.Message.Version), "v0")
	assert.Equal(t, tx.Message.RecentBlockHash, "CwT5Apc8MiSE136oEwAuFUxQypPQJapTo2TTtJ3TsLS4")
}

func TestGetSigningData(t *testing.T) {
	rawTx := "3SabRqqxKTdwNqsoa2P459eZX6RpFhonzCpUbjvv6oKj9yFpQ9D8URYW6GffzxfFgYzZg4f5ujbdZH5XigVpVzhQmzshWyFuoGFfSRZszHwzUpWA2Cxb2ugHfKwhvu8Ma7C67Q7wXQta8TpA2WfkaNieWHiC3YNMgW338StjVQnUXuoKQRofP5aYikTc3V6nAR1jGpwqwtopr5HFwddfrKqLzGJivuGc4oXZyfWJM6DH8BPzTWAtvvr8e1HZh4Bf4U1FWUnPa3WPoZvGV6RGTEsTGMYgrGRQ7FHaMefp147fnLH8HgWVWYSjib31kDa9MjSqXV9vvdcGXihNNogwRGC1C46b7z6ozy6WTzsxhr3oVMJ9PCzGx26EvZ96fheZpY7QRxbFuc2gQyzV5mM8qymt4Xu9aqjip8Nsxx8kEXY2cGJP9o1NWHc2EuvkwkJPrddMH7xN3U4xvwpCWKCy4sjYyLrzp1n6iuVBjXGJm74rJCxFDitSPZVh2cnFNHPwrR8eNBJA9bZK4d9Gbuiv2yjvzzfip8qF3S7ASS25rXXXMfvYAQ9Z4KM6LBbgksxU3VE7FBTsZ58cNN2VqsQe6b1h9FQFzByydYXKBJF3h2MGxFEvi1iKsurN6nrvTisK1qK5M5TmfrMyafvowuHE4xSpQ1tL6qjBjKovfBgX1jbbNNNTwQgPoS4cgU9CUihbTyUjWNdPrpPdDW7Sc762VdbwcNCAotDC8QMkWGQpiPumK5UMNgUgn1YSA2yrC5zjLzis7dFcTJC2oMSfyzcZsYqaFkHffe266w9rrbNqdz2azBpmfbG3VcDTCwKM7JDnBEjfAtSULdSfDMRJZx1wZYmv9ujFVeagfvmsUGAxeB1F29guxbKf3NYcKf4s4vMp95GNhtLu3VGc8AtBKGkm4bEdbxUcrLkNNEeumt6RUF9srQtf8Dw6J4AFtDQSnRvjCctmTiBAez93to2z3suvx45YZZNmmUsKbceekYYQoPe5j8hi5QLmXHwoePR3pX5sdGcwEYyzkPtW2ys5ithcLzhnBt2mZNo8gz69eyagaUat8dx6p7HEKMpok9FUEtsWvmDug6giH1mPgrS"
	tx, _ := NewTxFromRaw(rawTx, "base58")
	signingData, err := GetSigningData(tx)
	str := base58.Encode(signingData)

	assert.NoError(t, err)

	expected := "93MoPYVR3dSUVH81Dy28jm6xz9yxxvFmfpXSXD5phqFJrBiGEWYnUMJzJmB4D2gnCo6A3LYxWgDahzxBMWdGQyXCq2YAB2X4FUijoatjEF8kq7wZzbq92ukiHNX3TK6M7Jeaprcb9rN49uwNFyxdYjAtfdo9B8sg2uEUMn9PyHQN5czi9Ga7vQ2NLGocJroMPwfXKjWFhvzc1DGEJbZEtTSqjkJ6d71uNeCQYfUh3WXbTWJJC43YdFhf5CyuuFj2X3hkFzuQRJRtbwQdA7VhZbDHUGwLWjVH9S2Ccpn6AcqvUTdfeqAJ8xdmeEgMpKVTpfyyT738aMyXCSyUY9PaknJqqA7SnBVbuDyQ5S93iVHR79CLA6B6dras75hJQHosiH7R1WoYNbYg3JwHaj9L1odBsH2TBKSxRUfTT8fEeBUg7ffBYLhXHof4vHpgkLsKcLJ5tdrdkAcZ6s4LhMDqsxCdHHMxgXGsxr3QHEngjbxZRKWY9oyyBuyz9hKCqYCjQyWDBdPLNcL1rxGQfQhvrtujp21wqCcVcixN6gsLmfHd4fJAg4YPicxPyyZ1iJYh1nTQhtGajR4WTC1pYRtjiXYQcw2JWbi2p6n9R3g8wdvhfFfCHX3VsS73JVNxPA6HAoyi9aVacQUv7f9TFjsf2t5zkHVn5msvj5fCi9DpQdVwmafwv14n3BgMfCAoMxzdpX7TQYZ6pf5JWA9SvucKagNpH1QoZmHKdLBkEj96CWzJLjyVazKvLwWyPaPJaJVpGzyjpgWm84mBGJGhKZ5KBfnbYa3qwEG9qpfEXy1wpdc98KLGX5jmNz5SpbWR2jTuhmn5YQ7cQEvvQcsPE1cUtYBBiTb77UXoki1rKq23cYa8aX4Mx6u14aKJLyku9uVxvw9vbLBBpE7CJepUNeRULsqsikgiTZUujVs2TDyzijhctu8TDWfLbvFauMjXwT5mW82ReqeJ4pKB2dKJWB6HfH4nrjK4wAK2CtgZj2vT8NQ9fPDK6Y287HY"
	assert.Equal(t, expected, str)
}

func TestNewTxFromParamsForV0(t *testing.T) {
	txParams := SolanaTxParams{
		FeePayer:        "UHqRK53fc4sxHnUf8fADkNTU7pP71GZE1bFGFRnipvK",
		RecentBlockHash: "53tbv2gAF7n4E9sMtEX81E5WPCLPrb3wzQuoyKzWeKAP",
		Instructions: []Instruction{
			{
				ProgramId: "11111111111111111111111111111111",
				Keys: []AccountKey{
					{
						Pubkey:     "UHqRK53fc4sxHnUf8fADkNTU7pP71GZE1bFGFRnipvK",
						IsSigner:   true,
						IsWritable: true,
					},
					{
						Pubkey:     "CEVt1pcoSyC9b3PkARytWey6Bx8cuYxexW8k6hGCrQsr",
						IsSigner:   false,
						IsWritable: true,
					},
				},
				Data: "3Bxs43ZMjSRQLs6o",
			},
		},
		LookupTables: []LookupTable{
			{
				TableAccount: "GLhDJj8DoKzy7szzsETetjgye2tzTFPS98jKhc9NA2wn",
				AddressList: []string{
					"GPaR84jGe4ARGcfLaNS1XCH4VzWj5aLUhrGhe5sQzoYE",
					"GqcBjgCXJ1KoXjnpxDE5b7ZkqTDX3yiGvJuyuNQ4V3fL",
				},
			},
		},
	}

	tx, err := NewTxFromParams(txParams)
	assert.NoError(t, err)
	assert.Equal(t, string(tx.Message.Version), "v0")
	assert.Equal(t, tx.Message.RecentBlockHash, "53tbv2gAF7n4E9sMtEX81E5WPCLPrb3wzQuoyKzWeKAP")
	assert.Equal(t, tx.Message.Accounts[0].ToBase58(), "UHqRK53fc4sxHnUf8fADkNTU7pP71GZE1bFGFRnipvK")
}

func TestNewTxFromParamsForLegacy(t *testing.T) {
	txParams := SolanaTxParams{
		FeePayer:        "UHqRK53fc4sxHnUf8fADkNTU7pP71GZE1bFGFRnipvK",
		RecentBlockHash: "53tbv2gAF7n4E9sMtEX81E5WPCLPrb3wzQuoyKzWeKAP",
		Instructions: []Instruction{
			{
				ProgramId: "11111111111111111111111111111111",
				Keys: []AccountKey{
					{
						Pubkey:     "UHqRK53fc4sxHnUf8fADkNTU7pP71GZE1bFGFRnipvK",
						IsSigner:   true,
						IsWritable: true,
					},
					{
						Pubkey:     "CEVt1pcoSyC9b3PkARytWey6Bx8cuYxexW8k6hGCrQsr",
						IsSigner:   false,
						IsWritable: true,
					},
				},
				Data: "3Bxs43ZMjSRQLs6o",
			},
		},
	}

	tx, err := NewTxFromParams(txParams)
	assert.NoError(t, err)
	assert.Equal(t, string(tx.Message.Version), "legacy")
	assert.Equal(t, tx.Message.RecentBlockHash, "53tbv2gAF7n4E9sMtEX81E5WPCLPrb3wzQuoyKzWeKAP")
	assert.Equal(t, tx.Message.Accounts[0].ToBase58(), "UHqRK53fc4sxHnUf8fADkNTU7pP71GZE1bFGFRnipvK")
}
func TestAddSignatureWithParams(t *testing.T) {
	// [2, 160, 134, 1, 0] [3, 64, 13, 3, 0, 0, 0, 0, 0] [2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0]
	txParams := SolanaTxParams{
		FeePayer:        "7sTwHDmF2GTqPgFS37vXLeu6A2LcqASKM2y1pzz3xxPA",
		RecentBlockHash: "4fgbdQ4fgDrzyaU2vfZKv6Hvg3ZPRDrrrDnKWHrDfUV5",
		Instructions: []Instruction{
			{
				ProgramId: "11111111111111111111111111111111",
				Keys: []AccountKey{
					{
						Pubkey:     "7sTwHDmF2GTqPgFS37vXLeu6A2LcqASKM2y1pzz3xxPA",
						IsSigner:   true,
						IsWritable: true,
					},
					{
						Pubkey:     "12hXjGwDRu4xWGeWyPwBVkzEcjKCh6KgoC6YnDQGT7U9",
						IsSigner:   false,
						IsWritable: true,
					},
				},
				Data: "3Bxs412MvVNQj175",
			},
		},
	}

	tx, err := NewTxFromParams(txParams)
	assert.Nil(t, err)

	sig := "4dHGp38njfNnLDSQaFoJJjP47ZWZgzKYsupaCMzbsRDHUrC4TqhCm63wptPncNzLgqjbPEoLHob2dBbtQ1U7Btk"
	txData, err := AddSignature(tx, base58.Decode(sig), "base58")
	assert.Nil(t, err)

	expectedRawTx := "3oatFgjNLDpPqJqm1NwACf1UqMtYBZ6RfR3kQJk9fy2FP9moQCo4kiSgz6dAyQ7bG1jooK2LpAkt6etMG7PrHZgWhDp2Qd6dv9XdRc8dDgeCvq3wv8N3wjeqW6Cm46hQ3KL6Pupmzwh5cpuVxRCqZPYfLwTa9BAgKETNWNbzLmKi8UmhfKJpP8G91t2tujXn6WSbK7fBJj8izJk4vAM4ez8ttWJhWQWTDaYEPMjTN8dSoquN58CK44Sk4ThtP8mD7gs64nvQEhQSYFwxgXuWo1XwcdgZD15L1Jzo9"
	assert.Equal(t, expectedRawTx, txData.RawTx)
	assert.Equal(t, sig, txData.TxId)

	txData, err = AddSignature(tx, base58.Decode(sig), "base64")
	assert.Nil(t, err)

	expectedRawTx = "AQMgY2fIrveUUmmaWI/b5mxVpM59gCXg2vMBS0mfAGYPzaSU8VBWQ7jT/w0fBaJNb4DAAZbgpS7WeZpYdFCK0gEBAAEDZhOeVM2wZrmmeTmSXNeMhhL74KXFZ53kgb1KnaubJtkAb2hzS/2Kvk9Qk3ELJphLONpCG6NZvkf9+oqBBBGC3gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANnvnbilB51bCmPbBucKYrCp9obNxRXXupUvmUlH3a/gBAgIAAQwCAAAAAQAAAAAAAAA="
	assert.Equal(t, expectedRawTx, txData.RawTx)
	assert.Equal(t, sig, txData.TxId)
}

func TestAddSignatureWithPriorityFeeAndParams(t *testing.T) {
	// [2, 160, 134, 1, 0] [3, 64, 13, 3, 0, 0, 0, 0, 0] [2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0]
	txParams := SolanaTxParams{
		FeePayer:        "7sTwHDmF2GTqPgFS37vXLeu6A2LcqASKM2y1pzz3xxPA",
		RecentBlockHash: "Bb5oyrmU56gp3Mcs2muqPRFK7cRRhiGvVbkrNJZyeFnJ",
		Instructions: []Instruction{
			{
				ProgramId: "ComputeBudget111111111111111111111111111111",
				Keys:      []AccountKey{},
				Data:      "JC3gyu",
			},
			{
				ProgramId: "ComputeBudget111111111111111111111111111111",
				Keys:      []AccountKey{},
				Data:      "3QAwFKa3MJAs",
			},
			{
				ProgramId: "11111111111111111111111111111111",
				Keys: []AccountKey{
					{
						Pubkey:     "7sTwHDmF2GTqPgFS37vXLeu6A2LcqASKM2y1pzz3xxPA",
						IsSigner:   true,
						IsWritable: true,
					},
					{
						Pubkey:     "12hXjGwDRu4xWGeWyPwBVkzEcjKCh6KgoC6YnDQGT7U9",
						IsSigner:   false,
						IsWritable: true,
					},
				},
				Data: "3Bxs412MvVNQj175",
			},
		},
	}

	tx, err := NewTxFromParams(txParams)
	assert.Nil(t, err)

	sig := []byte{142, 137, 214, 216, 153, 139, 117, 234, 225, 189, 222, 233, 77, 152, 241, 213, 204, 244, 17, 135, 93, 83, 244, 206, 125, 62, 163, 120, 218, 176, 79, 225, 127, 125, 214, 43, 3, 102, 118, 179, 233, 166, 9, 125, 33, 30, 60, 183, 206, 248, 75, 113, 224, 176, 70, 41, 176, 187, 248, 218, 171, 112, 248, 1}
	txData, err := AddSignature(tx, sig, "base58")
	assert.Nil(t, err)

	expectedRawTx := "5ZvojwPnKZqYpRxuVkLpKpYvV9my9BK5X8a1PNEbXAwgbrZs5oRusyKsmrzsKSmWyFEc98W66NaFEjenRKJg8z2ZWm2aDELg9aaMrdMaRCjaLcUEPbNtjCCGgvbEjQKtGjDfKc2wJheix3Z1pB8JgdL6xFUbRMbB9h81u6yxuZqWCRjUFn2rrhXeDpULWB9rfSEFxQzKccrzg1p3NraDdfwritnjyTetjYkpEgMdtxx9t6yLpkNeohWSRjLotQUNYEj8RQVtDtMn3c5XNna7VFNpK9bvoFfcRh4u86iev5bZ2oXT3B4bWD9K7b9TH5QcrKNpjW8Ga8D89QGJfuR1eQgnvSAjNNg7QuXybbfaaHKu"
	assert.Equal(t, expectedRawTx, txData.RawTx)
	assert.Equal(t, base58.Encode(sig), txData.TxId)

	txData, err = AddSignature(tx, sig, "base64")
	assert.Nil(t, err)

	expectedRawTx = "AY6J1tiZi3Xq4b3e6U2Y8dXM9BGHXVP0zn0+o3jasE/hf33WKwNmdrPppgl9IR48t874S3HgsEYpsLv42qtw+AEBAAIEZhOeVM2wZrmmeTmSXNeMhhL74KXFZ53kgb1KnaubJtkAb2hzS/2Kvk9Qk3ELJphLONpCG6NZvkf9+oqBBBGC3gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwZGb+UhFzL/7K26csOb57yM5bvF9xJrLEObOkAAAACdUA4yyf/d77UjMWxVQEprbL+MajtD2a/K35wa+YZnQwMDAAUCoIYBAAMACQNADQMAAAAAAAICAAEMAgAAAAEAAAAAAAAA"
	assert.Equal(t, expectedRawTx, txData.RawTx)
	assert.Equal(t, base58.Encode(sig), txData.TxId)
}

