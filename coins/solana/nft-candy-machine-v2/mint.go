package nft_candy_machine_v2

import "github.com/okx/go-wallet-sdk/coins/solana/base"

func GetTokenWallet(wallet base.PublicKey, mint base.PublicKey, options ...string) (base.PublicKey, error) {
	tokenProgramID := base.TokenProgramID
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		tokenProgramID = base.Token2022ProgramID
	}
	addr, _, err := base.FindProgramAddress(
		[][]byte{
			wallet.Bytes(),
			tokenProgramID.Bytes(),
			mint.Bytes(),
		},
		base.SPLAssociatedTokenAccountProgramID,
	)
	return addr, err
}

func GetCandyMachineCreator(candyMachineAddress base.PublicKey) (base.PublicKey, uint8, error) {
	return base.FindProgramAddress(
		[][]byte{
			[]byte("candy_machine"),
			candyMachineAddress.Bytes(),
		},
		base.MetaplexCandyMachineV2ProgramID,
	)
}

func GetMetadata(mint base.PublicKey) (base.PublicKey, error) {
	addr, _, err := base.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			base.MetaplexTokenMetadataProgramID.Bytes(),
			mint.Bytes(),
		},
		base.MetaplexTokenMetadataProgramID,
	)
	return addr, err
}

func GetMasterEdition(mint base.PublicKey) (base.PublicKey, error) {
	addr, _, err := base.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			base.MetaplexTokenMetadataProgramID.Bytes(),
			mint.Bytes(),
			[]byte("edition"),
		},
		base.MetaplexTokenMetadataProgramID,
	)
	return addr, err
}
