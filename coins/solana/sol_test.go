package solana

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
	"github.com/okx/go-wallet-sdk/coins/solana/token"
	"github.com/okx/go-wallet-sdk/crypto/ed25519"
)

func Test_NewAddess(t *testing.T) {
	pk, _ := base.NewRandomPrivateKey()
	address, err := NewAddress(hex.EncodeToString(pk.Bytes()))
	require.NoError(t, err)
	require.True(t, ValidateAddress(address))
}

func Test_TransferTransaction(t *testing.T) {
	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := fromPrivate.PublicKey().String()
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendTransferInstruction(1000000000, from, to)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, _ := rawTransaction.Sign(true)
	expected := "4jijgudzgfQtujYrGnN66tv95LGUbnBv21vSWwyNQ185atbrW9b2pJdQsXXGBk3NMzrA7DcxNzkfFb3exJ11JG3JWj2WpWamCuDqza2Xg2Eh4ZhFKgYLhnXjyVdFDFtxjPa2t3xNUvLi1x1g2oE8jTcmq3ZjyQ2EFi1aNQVTwtg8eJLkFjr5kLjzn6tjnzstscj1A495KAWR3FETjHk2dTU6itaMJiSZ8sxMUZSEWKiJPDvD4MWN4vu8FwHtdWYABavzMzAxowskqevbiGKaezzAoN3zr5hJrEjQj"
	require.Equal(t, expected, tx)
}

func Test_TokenTransferTransaction(t *testing.T) {
	hash := "H6TNM3fDg5wTYT4eiv2PnGdd1555a45FEJtxVLtzv9dJ"

	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	from := fromPrivate.PublicKey().String()
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"

	mint := "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU"

	fromAssociated, _, _ := base.FindAssociatedTokenAddress(base.MustPublicKeyFromBase58(from), base.MustPublicKeyFromBase58(mint))
	toAssociated, _, _ := base.FindAssociatedTokenAddress(base.MustPublicKeyFromBase58(to), base.MustPublicKeyFromBase58(mint))

	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendAssociatedTokenAccountCreateInstruction(from, to, mint)
	rawTransaction.AppendTokenTransferInstruction(1000000, fromAssociated.String(), toAssociated.String(), from)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, _ := rawTransaction.Sign(true)
	expected := "7r9muRWaFEQC5wYCaXqtrr6BbPtZfm3pUsAFdzVrqaunHk1f6vjgi4GFa7d8ABppS9y6p6uCWLv7rraoTuA5FxEkmfBx6dNu3wgAGmxeQahgK91quoDKfQrCEGnsi3TV8pfykomPxejDdczHdq8LnCTQ5uskWyJknDuCrJDw2JH68yN5BpgwBy5k5UmAvmU7CMxaWwhNRRXv8sxVhNHkvFc4EaLuEttoaQ8CPiN85rqX4qVK3MRBMUVUtBoWDUSgEsFhBJVzXtcpEZ6htdqHqLPevJomKgfrLE2Wz7e52P4rr6dAst2nXKRHLvaTJTqhwG8d5YJ5SZpfuALXm8GN7VDogAzDzZrjXg6LXBWAiUNuBNXREWdd4mqyNZcjoUxn4Af4GkwX1fSZgMzmVym7otWeStW35me1CqwT6rqBgEtbn1UKMT4rqKMdCBTA4MZHruHfiJwwh5WpWrdUTCGa67jRvTXRdHnUAAfUkQpMNccrptZdCSqWHCiE9C2xMwqtZTTTW1avQ8t3sBsmQz775KxeKDin7aXaE5TAopjqDry6FG4FZwXFuq"
	require.Equal(t, expected, tx)
}

func Test_TokenApproveTransaction(t *testing.T) {
	hash := "H6TNM3fDg5wTYT4eiv2PnGdd1555a45FEJtxVLtzv9dJ"

	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	from := fromPrivate.PublicKey().String()

	mint := "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU"
	delegate := base.MustPublicKeyFromBase58(from)

	fromAssociated, _, _ := base.FindAssociatedTokenAddress(base.MustPublicKeyFromBase58(from), base.MustPublicKeyFromBase58(mint))

	inst := token.NewApproveInstruction(1000000, fromAssociated, delegate, fromPrivate.PublicKey(), []base.PublicKey{}).Build()

	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendInstruction(inst)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, _ := rawTransaction.Sign(true)
	expected := "DR4zXTf95VAmywv5CjEi4tvUHm6sAWz4UByrkxKqF7z6sQ2JYZta2cMRTBhwGvZSuEaEhZSB8oWDKPEUKJKFRecQm8RZU93qee9KP9X7JXBJmZABuy5Q79Fpz3gQUJc8nqzPcyyFTAfvhXU3mcKoRCeHoaTruaLEVdXyjuarVcJR6izd89NZX728pYnKyzKTkDPCC92VnPHBsi9RnGFAyr6SfBAEQrKbuczcPdGyMzcTdFFMZatxDwuk8QRssZA9nomrRR6mPX6M6u5FuVTKjcEuWG7CDx2PTm"
	require.Equal(t, expected, tx)
}

func Test_UnMarshall(t *testing.T) {
	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := fromPrivate.PublicKey().String()

	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendTransferInstruction(1000000000, from, to)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	rawTx, _ := rawTransaction.Sign(true)
	signers := make([]string, 0)
	signers = append(signers, hex.EncodeToString(fromPrivate.Bytes()))
	tx, _ := DecodeAndSign(rawTx, signers, hash, true)
	expected := "4jijgudzgfQtujYrGnN66tv95LGUbnBv21vSWwyNQ185atbrW9b2pJdQsXXGBk3NMzrA7DcxNzkfFb3exJ11JG3JWj2WpWamCuDqza2Xg2Eh4ZhFKgYLhnXjyVdFDFtxjPa2t3xNUvLi1x1g2oE8jTcmq3ZjyQ2EFi1aNQVTwtg8eJLkFjr5kLjzn6tjnzstscj1A495KAWR3FETjHk2dTU6itaMJiSZ8sxMUZSEWKiJPDvD4MWN4vu8FwHtdWYABavzMzAxowskqevbiGKaezzAoN3zr5hJrEjQj"
	require.Equal(t, expected, tx)
}

func Test_UnMarshall2(t *testing.T) {
	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"

	unsignedTx := "{\"bizType\":\"okxdex\",\"data\":\"DvRxX9KypbY3wo9CW9sDNwryG4vrHhMWEHr5n2uickAN9FHzUn2rgDRvDKvRa67zfLcoU9B7VYz2Lmw3rmFziKzuUNLKQPZGGEmKAnwYfuPf4SoEisUN8msqiHRjWeaWdFWer6V9DvJDSht8byGYp6LBV7VJK2J7YWbMM592pnodkvWH343R38BdxVBFE9SBjAr1vdnvrGRbHG5f5zrZbhoeAssHXSe2p2Sf8dsb6rRdMnTWPmK12e4rPx7PGCvGJzJ5iRmu7gKENE3HrE4pBuTTQqdq2UA2SRXXH3g8vZ8RQKg6spAYC2EU3Q5A9ehfr1mKAVxb2xmMeQ1zJ4uicHC6zbwYfi9BrvaCiScqJgJZLSpZAZXvHSeZUwcnvnhh3gJCQXN8hmmBJuToi5A1eYWUAALyLc3ZrPHDUD5iYwtXiLZ8RUhfhj9siJFSXFvdBJXnoHpwvaSAEYggjWrFCHob6di7zMC5s9pAHDALApc27yVgq5fR7Bq99QFTpXHCveHtEHycVV3Lfdu4UPyooemjGFLqmqd8aNND4SFraXoAwopPnocAfG7TStiq6jrfXdkvt8MN25D9ZxSjHMy6dCvgVgAtKD4kiMonPkEix4dTx3gBs3UqAo4GriZHKqjuUoBiHxiVhK1cnQGEvV9HFbvNXtUi2KRVa1i8ZzhfeNj7h6n4mtEchzGS9HYDZdUTNz1d5HcRifi9vD5Y7zgPkQx1aFjWFHJ8zZ5RNPshsxrwFSjbfEVM6p7o3iJwD5qKjsmmnMUMV77XkBNCBVkuW8wgeDsoYRCKRejUa13SKhAaWhhPfqnNW6GjmZfNG51RrEcMmN4BVTuk4nX2ApoeiL5WQjotzdTQHHBudpFTUv1MYkJXfigA8YWmxQQMgfwS29m9SHxGGS5VcN137aHFqgYdxt6oe6wXwQi5sEB2AK8LiZ3tXbB9Tf58kuPso5zdYxJaKfKT8S17ZzcpqPpmMAbmpbYcEzhr4ccc5fQxwMJRVWQWGjSeqsWvPPJwSWLoHWRu6U4qv\",\"bizId\":[\"1666182406446\"],\"from\":\"3cUbuUEJkcgtzGxvsukksNzmgqaUK9jwFS5pqRpoevtN\",\"to\":\"4itxBe4qBAwhB9zpEAw31d7w8o7gTQscYpxhRtUemjF9\",\"accountLength\":\"293\"}"
	tx, err := DecodeAndMultiSign(unsignedTx, hex.EncodeToString(fromPrivate.Bytes()), hash, true)
	require.NoError(t, err)
	expected := "DFWhCQPXxfBK6aBN2EhPJjUL6yhfvWr9N5Jeup5E2SrEEq1pDJiyyxGfpCNj8WVwer8LdUSL9bJ6UAmJR4GCwHKC6isxmkwLPJi1Pz46Qy6QKRXj1bjRtWtq9SyXV2ZdhxXmAtDMxtdvph9UmVn3vkxiz4hqXtT71okjjGUwqZQXwqRxbRpL79SVwCbw94WfYkg9ofyR6XecYv5UUFmXGKoTsM38qV9qKHLsT4CmQHPzciQ6ZpofBFEUv6CsFVPpPz8zkHY5gb7PnPBHur8n7oQ2umTY2iLsspxCzZfGTysp4Te9vmo4GkutFhWroF3qGZPoL4d5wvAr746vJ8r8k9Mkc9DB1W24SvpJnXuGMt5SonDYLLYHLM2CL2t97ugxBmFKFdDQYj57YSL2c2opaTU19R17oz8YSVnyFwqhWLQFyjyE4jnpExUc58K4Y2dbGega4gQX948tDKNX3HfZTB1gneJYMDbyDeA2E2UYDgdFymD3Tod6tMVoTvCXt1jXunisjGqUsP34TPDChWEbZFnz2WT3PuUwQ5cjhNr8Z7p1vePJwqmJrEFk7CPRigyWkERdnp7eNEjrmE9LmDsRnqaq5NbtwemPkpsSgAN3V1wu7gMGNH5FJZi775vjrGMaeEBbgDPnFmwWmTCjW1r27rjZ6mwrc2rcJe5EPvNAYzi1TrwVxYUzJXNDPkdkjMshnb4rhFB7RMwyuXghTmwN1ZRyXcdwf9eUUHr29RSbxyYuiCx4e4iVYWdL5EcAGs7rF8U8UNSZ4RWCVVZrrPKiw2vg251ELsZJt6b6eEqfW3iU3x7BjQchR6fgkr5obL57iZzjjrnbFaeq2aA4fEQGmSerXpb8v1eWJL6pV4knZpGwEP4h2r4WndrP7Q7SXQ9ik38kbLHsX5tWavJi8UF5kzBjNKasGPHRRvbj3WFmh2qjDtkS24aS1AdGSXTcfguodhGUA8XTqsNorniV8W6HHMTZF1c2zPbmVbom3Qz6cLrGhnor4Zr6jkDAmfLqQdG7Ut2VaNz7J"
	require.Equal(t, expected, tx)
}

func Test_TransferTransaction1(t *testing.T) {
	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := fromPrivate.PublicKey().String()
	nonceAddress := "FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1"
	sign := "c4f08e638a6735ae13f758aa2e72794ec84296c317b661c4814d83f16e2aa5dfbab7ffc34f6d0589e9d001f3ef432fa38bd0eb1c6a864f4b6348c379585e4103"
	tx, _ := SignedTx(hash, from, to, nonceAddress, 10000, sign)

	expected := "7DhA4Xf5cvq8B7CawxbdaJCeNVDZ8MxzY4dEykNTnFG1Yycu45nZsjqzSngGuF7WSMbwGRz4eTpsayBs86CqU7PtDufJ7nzZP3s9gRM2qjB2P5Lyq2uxFG4RvTcHEbB2m45JhiELsB1759br4zZdXNHEbJPVPGhPitgNfLG7Hyxoqcmze2uuk9Vdg1Lviiw2SbGDycnY9KFqySXBeyFUQR3WrYM1XaFgJ9c8RPfz9WHyqKnot3nqaP2kNjv1Ps5s8r49hd96JE7ArEZCS5WyoNUS9dVjUmhryER1e1TcrZ87ceTCmUVFoZNSauTZXnoYfq8WfZ1mekvQhbGJsriKTxfvE2gVWnT8caWPDwVX1fcNnY5XbR778F3NuBfsFb8CQLsNrLJUa3"
	require.Equal(t, expected, tx)
}
func Test_TransferWithNoce(t *testing.T) {

	private := "b90ae8f3c465425f561ebad958dd2e385ce9aeb95259f07af1550cfb6c7c90ec"
	privateBytes, err := hex.DecodeString(private)
	require.NoError(t, err)
	key, err := ed25519.PrivateKeyFromSeed(hex.EncodeToString(privateBytes))
	require.NoError(t, err)
	nonceAddress := "29odEnJWGSCcWx3o7hoAPdpaDuZfyjFdDEs3q5WsfJVp"
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "8awFZzqF8KuYuXjRKWibsehoiJrt9qJXFXBNSDvkHyi8"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := "5vWSQFWuHuwz3cCHY3MYXB3twp6w4UtXAFG2VeqALGUq"

	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendAdvanceNonceInstruction(from, nonceAddress)
	rawTransaction.AppendTransferInstruction(10000000, from, to)
	rawTransaction.AppendSigner(hex.EncodeToString(key[:]))
	tx, err := rawTransaction.Sign(true)
	require.NoError(t, err)
	expected := "6B3zHpjnSZU1jDVQv5pjJdkoKToBQYUAAGq2VNRHTHfiPL3JYY35o9x4U18gVMykHajHvNK9BNeLrzPBiZnWRRf6MTP5G5wTgseWCkCzcDh4WkMfZtDmDcYyhEPbcAWL1rtu474d9NVC3ypbsPwNyHo2oZ7a7hCwXa4p2idtvffWhwZ36df9xzVJzEPWBWFJasqkcRXR3SsG4DvEdv7S9BTySZTvpKUz5rX3FRhP6PdRtiPXrpHKjfP9AuqVvpgsbdCkz8wE1HyJa6ihgWap1zmqFbT5uny9mmgLas657jgKSedKoxLSeeiVRcEkBJZDMYsy1JH7soPYzY7PgHwWMjos7BNdaKt2fqksqV8yW7RQGj7FpzaBxNwJX2GcsYZumsGwgjhNyV"
	require.Equal(t, expected, tx)
}
