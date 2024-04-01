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
	from := fromPrivate.PublicKey().String()
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendTransferInstruction(1000000000, from, to)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, err := rawTransaction.Sign(true)
	require.NoError(t, err)
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
	tx, err := rawTransaction.Sign(true)
	if err != nil {
		// todo
	}
	require.NoError(t, err)
	expected := "7r9muRWaFEQC5wYCaXqtrr6BbPtZfm3pUsAFdzVrqaunHk1f6vjgi4GFa7d8ABppS9y6p6uCWLv7rraoTuA5FxEkmfBx6dNu3wgAGmxeQahgK91quoDKfQrCEGnsi3TV8pfykomPxejDdczHdq8LnCTQ5uskWyJknDuCrJDw2JH68yN5BpgwBy5k5UmAvmU7CMxaWwhNRRXv8sxVhNHkvFc4EaLuEttoaQ8CPiN85rqX4qVK3MRBMUVUtBoWDUSgEsFhBJVzXtcpEZ6htdqHqLPevJomKgfrLE2Wz7e52P4rr6dAst2nXKRHLvaTJTqhwG8d5YJ5SZpfuALXm8GN7VDogAzDzZrjXg6LXBWAiUNuBNXREWdd4mqyNZcjoUxn4Af4GkwX1fSZgMzmVym7otWeStW35me1CqwT6rqBgEtbn1UKMT4rqKMdCBTA4MZHruHfiJwwh5WpWrdUTCGa67jRvTXRdHnUAAfUkQpMNccrptZdCSqWHCiE9C2xMwqtZTTTW1avQ8t3sBsmQz775KxeKDin7aXaE5TAopjqDry6FG4FZwXFuq"
	require.Equal(t, expected, tx)
}

func Test_Token2022TransferTransaction(t *testing.T) {
	hash := "HqpUiCHdybmpK91LF9pVwTtCkMfcXhwBVAbxpsNPUuFk"
	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	from := fromPrivate.PublicKey().String()
	to := "GbDq1KMiTmSys7SPwNTJVF3oSvnpirihdZyqpNTBnf3R"
	mint := "FTDMffVuqMpPPTdfaDTNgMTx7A8xe2jpPQBzMq3D85yi"
	fromAssociated, _, _ := base.FindAssociatedTokenAddress(base.MustPublicKeyFromBase58(from), base.MustPublicKeyFromBase58(mint), base.TOKEN2022)
	toAssociated, _, _ := base.FindAssociatedTokenAddress(base.MustPublicKeyFromBase58(to), base.MustPublicKeyFromBase58(mint), base.TOKEN2022)
	rawTransaction := NewRawTransaction(hash, from)
	// create token account
	rawTransaction.AppendAssociatedTokenAccountCreateInstruction(from, to, mint, base.TOKEN2022)
	rawTransaction.AppendTokenTransferInstruction(1000000, fromAssociated.String(), toAssociated.String(), from, base.TOKEN2022)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, err := rawTransaction.Sign(true)
	require.NoError(t, err)
	expected := "ACoQ83r9fiirEcQSaoP4sTi13FbeUtx7mD6CWHQCcTMj6RuCA3xoJ8eZygnrWFsSxJvBMDFbKmk2waUM35NLR6MuwifUvykvdZXBEU1Kx2ejkbjmcJwMq4wXRXXNgYV1A1W7frVmpqiPnAuhLASLSCw6LFLFcaytQJb76hee6X4cr3nzzPSrn4mapgtwyVBeTRWZiNpENUWPmKSXvcwgtfR3SKddJ4GLX9N1QaHAZKnoQe629VbWvpAJh8RmFq58wkDGPPdbmpiSDpzALDJzEQXVCMYYikeSJSiuNaXtaqVnGgvDp751CLci5NfoqXDnppTP7ENGVz7KG5vPqj4B4EZsWPbpazq4obRqPKU3dCPESB6qLY8GdxFgSnrVxfFFsttLdSyK2u8wxqLMuSxEcLEXFmHHSLkPdGo3BSmHrZwq4eLPr5P5kH95PYCHRv7L1drLHwVwruAmj2SBBXjQQ3xZPCPBqwbT7AbAPepQy8DQgpfiTpyYWLQwQfgGSxrWK8W914fxGBA2sfHkR6irLZRv33z3jD7uZMtYfoyu2TgusTVkewVVif"
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

	unsignedTx := "{\"bizType\":\"okxdex\",\"data\":\"4jijgudzgfQtujYrGnN66tv95LGUbnBv21vSWwyNQ185atbrW9b2pJdQsXXGBk3NMzrA7DcxNzkfFb3exJ11JG3JWj2WpWamCuDqza2Xg2Eh4ZhFKgYLhnXjyVdFDFtxjPa2t3xNUvLi1x1g2oE8jTcmq3ZjyQ2EFi1aNQVTwtg8eJLkFjr5kLjzn6tjnzstscj1A495KAWR3FETjHk2dTU6itaMJiSZ8sxMUZSEWKiJPDvD4MWN4vu8FwHtdWYABavzMzAxowskqevbiGKaezzAoN3zr5hJrEjQj\",\"bizId\":[\"1666182406446\"],\"from\":\"3cUbuUEJkcgtzGxvsukksNzmgqaUK9jwFS5pqRpoevtN\",\"to\":\"4itxBe4qBAwhB9zpEAw31d7w8o7gTQscYpxhRtUemjF9\",\"accountLength\":\"293\"}"
	tx, err := DecodeAndMultiSign(unsignedTx, hex.EncodeToString(fromPrivate.Bytes()), hash, true)
	require.NoError(t, err)
	expected := "4jijgudzgfQtujYrGnN66tv95LGUbnBv21vSWwyNQ185atbrW9b2pJdQsXXGBk3NMzrA7DcxNzkfFb3exJ11JG3JWj2WpWamCuDqza2Xg2Eh4ZhFKgYLhnXjyVdFDFtxjPa2t3xNUvLi1x1g2oE8jTcmq3ZjyQ2EFi1aNQVTwtg8eJLkFjr5kLjzn6tjnzstscj1A495KAWR3FETjHk2dTU6itaMJiSZ8sxMUZSEWKiJPDvD4MWN4vu8FwHtdWYABavzMzAxowskqevbiGKaezzAoN3zr5hJrEjQj"
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
