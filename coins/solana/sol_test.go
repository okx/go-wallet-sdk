package solana

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
	"github.com/okx/go-wallet-sdk/coins/solana/token"
	"github.com/okx/go-wallet-sdk/crypto/ed25519"
)

func Test_Addess(t *testing.T) {
	pk, _ := base.NewRandomPrivateKey()
	address, _ := NewAddress(hex.EncodeToString(pk.Bytes()))
	fmt.Println(address)
	fmt.Println(ValidateAddress(address))
}

func Test_TransferTransaction(t *testing.T) {
	fromPrivate, _ := base.PrivateKeyFromBase58("//todo please replace your key")
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := fromPrivate.PublicKey().String()

	// https://api.testnet.solana.com
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendTransferInstruction(1000000000, from, to)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, _ := rawTransaction.Sign(true)
	fmt.Println(tx)
}

func Test_TokenTransferTransaction(t *testing.T) {
	hash := "H6TNM3fDg5wTYT4eiv2PnGdd1555a45FEJtxVLtzv9dJ"

	fromPrivate, _ := base.PrivateKeyFromBase58("//todo please replace your key")
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
	fmt.Println(tx)
}

func Test_TokenApproveTransaction(t *testing.T) {
	hash := "H6TNM3fDg5wTYT4eiv2PnGdd1555a45FEJtxVLtzv9dJ"

	fromPrivate, _ := base.PrivateKeyFromBase58("//todo please replace your key")
	from := fromPrivate.PublicKey().String()

	mint := "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU"
	delegate := base.MustPublicKeyFromBase58(from)

	fromAssociated, _, _ := base.FindAssociatedTokenAddress(base.MustPublicKeyFromBase58(from), base.MustPublicKeyFromBase58(mint))

	inst := token.NewApproveInstruction(1000000, fromAssociated, delegate, fromPrivate.PublicKey(), []base.PublicKey{}).Build()

	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendInstruction(inst)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, _ := rawTransaction.Sign(true)
	fmt.Println(tx)
}

func Test_UnMarshall(t *testing.T) {
	fromPrivate, _ := base.PrivateKeyFromBase58("//todo please replace your key")
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := fromPrivate.PublicKey().String()

	// https://api.testnet.solana.com
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendTransferInstruction(1000000000, from, to)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	rawTx, _ := rawTransaction.Sign(true)
	fmt.Println("rawTX:=" + rawTx)

	signers := make([]string, 0)
	signers = append(signers, hex.EncodeToString(fromPrivate.Bytes()))
	tx, _ := DecodeAndSign(rawTx, signers, hash, true)
	fmt.Println("tx:=" + tx)
}

func Test_UnMarshall2(t *testing.T) {
	fromPrivate, _ := base.PrivateKeyFromBase58("//todo please replace your key")
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"

	unsignedTx := "{\"bizType\":\"okxdex\",\"data\":\"DvRxX9KypbY3wo9CW9sDNwryG4vrHhMWEHr5n2uickAN9FHzUn2rgDRvDKvRa67zfLcoU9B7VYz2Lmw3rmFziKzuUNLKQPZGGEmKAnwYfuPf4SoEisUN8msqiHRjWeaWdFWer6V9DvJDSht8byGYp6LBV7VJK2J7YWbMM592pnodkvWH343R38BdxVBFE9SBjAr1vdnvrGRbHG5f5zrZbhoeAssHXSe2p2Sf8dsb6rRdMnTWPmK12e4rPx7PGCvGJzJ5iRmu7gKENE3HrE4pBuTTQqdq2UA2SRXXH3g8vZ8RQKg6spAYC2EU3Q5A9ehfr1mKAVxb2xmMeQ1zJ4uicHC6zbwYfi9BrvaCiScqJgJZLSpZAZXvHSeZUwcnvnhh3gJCQXN8hmmBJuToi5A1eYWUAALyLc3ZrPHDUD5iYwtXiLZ8RUhfhj9siJFSXFvdBJXnoHpwvaSAEYggjWrFCHob6di7zMC5s9pAHDALApc27yVgq5fR7Bq99QFTpXHCveHtEHycVV3Lfdu4UPyooemjGFLqmqd8aNND4SFraXoAwopPnocAfG7TStiq6jrfXdkvt8MN25D9ZxSjHMy6dCvgVgAtKD4kiMonPkEix4dTx3gBs3UqAo4GriZHKqjuUoBiHxiVhK1cnQGEvV9HFbvNXtUi2KRVa1i8ZzhfeNj7h6n4mtEchzGS9HYDZdUTNz1d5HcRifi9vD5Y7zgPkQx1aFjWFHJ8zZ5RNPshsxrwFSjbfEVM6p7o3iJwD5qKjsmmnMUMV77XkBNCBVkuW8wgeDsoYRCKRejUa13SKhAaWhhPfqnNW6GjmZfNG51RrEcMmN4BVTuk4nX2ApoeiL5WQjotzdTQHHBudpFTUv1MYkJXfigA8YWmxQQMgfwS29m9SHxGGS5VcN137aHFqgYdxt6oe6wXwQi5sEB2AK8LiZ3tXbB9Tf58kuPso5zdYxJaKfKT8S17ZzcpqPpmMAbmpbYcEzhr4ccc5fQxwMJRVWQWGjSeqsWvPPJwSWLoHWRu6U4qv\",\"bizId\":[\"1666182406446\"],\"from\":\"3cUbuUEJkcgtzGxvsukksNzmgqaUK9jwFS5pqRpoevtN\",\"to\":\"4itxBe4qBAwhB9zpEAw31d7w8o7gTQscYpxhRtUemjF9\",\"accountLength\":\"293\"}"
	tx, err := DecodeAndMultiSign(unsignedTx, hex.EncodeToString(fromPrivate.Bytes()), hash, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
}
func Test_TransferTransaction1(t *testing.T) {
	fromPrivate, _ := base.PrivateKeyFromBase58("//todo please replace your key")
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := fromPrivate.PublicKey().String()
	nonceAddress := "FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1"
	sign := "c4f08e638a6735ae13f758aa2e72794ec84296c317b661c4814d83f16e2aa5dfbab7ffc34f6d0589e9d001f3ef432fa38bd0eb1c6a864f4b6348c379585e4103"
	tx, _ := SignedTx(hash, from, to, nonceAddress, 10000, sign)

	fmt.Println(tx)
}
func Test_TransferWithNoce(t *testing.T) {

	private := "//todo please replace your key"
	privateBytes, _ := hex.DecodeString(private)
	key, _ := ed25519.PrivateKeyFromSeed(hex.EncodeToString(privateBytes))
	nonceAddress := "29odEnJWGSCcWx3o7hoAPdpaDuZfyjFdDEs3q5WsfJVp"
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "8awFZzqF8KuYuXjRKWibsehoiJrt9qJXFXBNSDvkHyi8"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := "8yYSNvcxVLqVjhHJvDSXpFndkqiKVpq6w1KKvkpZGzmM"

	// https://api.testnet.solana.com
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendAdvanceNonceInstruction(from, nonceAddress)
	rawTransaction.AppendTransferInstruction(10000000, from, to)
	rawTransaction.AppendSigner(hex.EncodeToString(key[:]))
	//rawTransaction.AppendAdvanceNonceInstruction()
	tx, _ := rawTransaction.Sign(true)
	fmt.Println(tx)
}
