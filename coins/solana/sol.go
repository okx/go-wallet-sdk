package solana

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	associatedtokenaccount "github.com/okx/go-wallet-sdk/coins/solana/associated-token-account"
	computebudget "github.com/okx/go-wallet-sdk/coins/solana/compute-budget"
	"github.com/tyler-smith/go-bip39"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
	"github.com/okx/go-wallet-sdk/coins/solana/system"
	"github.com/okx/go-wallet-sdk/coins/solana/token"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

type RawTransaction struct {
	instructions []base.Instruction
	blockHash    string
	payer        base.PublicKey
	signers      []string
}

func NewAddress(privateKeyHex string) (string, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}
	return base.PrivateKey(privateKeyBytes).PublicKey().String(), nil
}

func NewAddressByPublic(pubKey string) (string, error) {
	pubKeyByte, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", err
	}
	return base58.Encode(pubKeyByte), nil
}

func ValidateAddress(address string) bool {
	_, err := base.PublicKeyFromBase58(address)
	return err == nil
}

func NewRawTransaction(blockHash string, payer string) *RawTransaction {
	t := RawTransaction{}
	t.instructions = make([]base.Instruction, 0)
	t.blockHash = blockHash
	t.payer = base.MustPublicKeyFromBase58(payer)
	t.signers = make([]string, 0)
	return &t
}

func (t *RawTransaction) AppendAdvanceNonceInstruction(authorized string, nonceAccount string) {
	authorizedPublicKey := base.MustPublicKeyFromBase58(authorized)
	nonceAccountPublicKey := base.MustPublicKeyFromBase58(nonceAccount)
	inst := system.NewAdvanceNonceAccountInstruction(nonceAccountPublicKey, base.SysVarRecentBlockHashesPubkey, authorizedPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendAllocateInstruction(
	space uint64,
	newAccount string) {
	newAccountPublicKey := base.MustPublicKeyFromBase58(newAccount)
	inst := system.NewAllocateInstruction(space, newAccountPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendAssignInstruction(
	owner string,
	assignedAccount string) {
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	assignedAccountPublicKey := base.MustPublicKeyFromBase58(assignedAccount)
	inst := system.NewAssignInstruction(ownerPublicKey, assignedAccountPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendAuthorizeNonceAccountInstruction(
	authorized string,
	nonceAccount string,
	nonceAuthorityAccount string) {
	authorizedPublicKey := base.MustPublicKeyFromBase58(authorized)
	nonceAccountPublicKey := base.MustPublicKeyFromBase58(nonceAccount)
	nonceAuthorityAccountPublicKey := base.MustPublicKeyFromBase58(nonceAuthorityAccount)
	inst := system.NewAuthorizeNonceAccountInstruction(authorizedPublicKey, nonceAccountPublicKey, nonceAuthorityAccountPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendCreateAccountInstruction(
	lamports uint64,
	space uint64,
	owner string,
	fundingAccount string,
	newAccount string) {
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	fundingAccountPublicKey := base.MustPublicKeyFromBase58(fundingAccount)
	newAccountPublicKey := base.MustPublicKeyFromBase58(newAccount)
	inst := system.NewCreateAccountInstruction(lamports, space, ownerPublicKey, fundingAccountPublicKey, newAccountPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendInitializeNonceAccountInstruction(authorized string, nonceAccount string) {
	authorizedPublicKey := base.MustPublicKeyFromBase58(authorized)
	nonceAccountPublicKey := base.MustPublicKeyFromBase58(nonceAccount)
	inst := system.NewInitializeNonceAccountInstruction(authorizedPublicKey, nonceAccountPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendPriorityFeeInstruction(computeUnitLimit uint32, computeUnitPrice uint64) {
	if computeUnitLimit > 0 && computeUnitPrice > 0 {
		computeUnitLimitInstruction := computebudget.NewSetComputeUnitLimitInstruction(computeUnitLimit).Build()
		computeUnitPriceInstruction := computebudget.NewSetComputeUnitPriceInstruction(computeUnitPrice).Build()
		t.instructions = append(t.instructions, computeUnitLimitInstruction)
		t.instructions = append(t.instructions, computeUnitPriceInstruction)
	}
}

func (t *RawTransaction) AppendTransferInstruction(
	lamports uint64,
	fundingAccount string,
	recipientAccount string) {
	fundingAccountPublicKey := base.MustPublicKeyFromBase58(fundingAccount)
	recipientAccountPublicKey := base.MustPublicKeyFromBase58(recipientAccount)
	inst := system.NewTransferInstruction(lamports, fundingAccountPublicKey, recipientAccountPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendWithdrawNonceAccountInstruction(
	lamports uint64,
	nonceAccount string,
	recipientAccount string,
	nonceAuthorityAccount string) error {
	nonceAccountPublicKey, err := base.PublicKeyFromBase58(nonceAccount)
	if err != nil {
		return err
	}
	recipientAccountPublicKey, err := base.PublicKeyFromBase58(recipientAccount)
	if err != nil {
		return err
	}
	nonceAuthorityAccountPublicKey, err := base.PublicKeyFromBase58(nonceAuthorityAccount)
	if err != nil {
		return err
	}

	inst := system.NewWithdrawNonceAccountInstruction(lamports, nonceAccountPublicKey,
		recipientAccountPublicKey, nonceAuthorityAccountPublicKey).Build()
	t.instructions = append(t.instructions, inst)
	return nil
}

func (t *RawTransaction) AppendTokenApproveInstruction(
	amount uint64,
	// Accounts:
	source string,
	delegate string,
	owner string) {
	sourcePublicKey := base.MustPublicKeyFromBase58(source)
	delegatePublicKey := base.MustPublicKeyFromBase58(delegate)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewApproveInstruction(amount, sourcePublicKey, delegatePublicKey, ownerPublicKey, []base.PublicKey{}).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendTokenApproveCheckedInstruction(
	amount uint64,
	decimals uint8,
	source string,
	mint string,
	delegate string,
	owner string) {
	sourcePublicKey := base.MustPublicKeyFromBase58(source)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	delegatePublicKey := base.MustPublicKeyFromBase58(delegate)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewApproveCheckedInstruction(amount, decimals, sourcePublicKey, mintPublicKey, delegatePublicKey, ownerPublicKey, []base.PublicKey{}).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendTokenBurnInstruction(amount uint64, source string, mint string, owner string, options ...string) {
	sourcePublicKey := base.MustPublicKeyFromBase58(source)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewBurnInstruction(amount, sourcePublicKey, mintPublicKey, ownerPublicKey, []base.PublicKey{}).Build()
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		inst.SetProgramID(base.Token2022ProgramID)
	}
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendTokenBurnCheckedInstruction(
	// Parameters:
	amount uint64,
	decimals uint8,
	// Accounts:
	source string,
	mint string,
	owner string,
	options ...string) {
	sourcePublicKey := base.MustPublicKeyFromBase58(source)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewBurnCheckedInstruction(amount, decimals, sourcePublicKey, mintPublicKey, ownerPublicKey, []base.PublicKey{}).Build()
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		inst.SetProgramID(base.Token2022ProgramID)
	}
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendCloseAccountInstruction(
	account string,
	destination string,
	owner string) {
	accountPublicKey := base.MustPublicKeyFromBase58(account)
	destinationPublicKey := base.MustPublicKeyFromBase58(destination)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewCloseAccountInstruction(accountPublicKey, destinationPublicKey, ownerPublicKey, []base.PublicKey{}).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendFreezeAccountInstruction(
	account string,
	mint string,
	authority string) {
	accountPublicKey := base.MustPublicKeyFromBase58(account)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	authorityPublicKey := base.MustPublicKeyFromBase58(authority)
	inst := token.NewFreezeAccountInstruction(accountPublicKey, mintPublicKey, authorityPublicKey, []base.PublicKey{}).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendInitializeAccountInstruction(
	account string,
	mint string,
	owner string) {
	accountPublicKey := base.MustPublicKeyFromBase58(account)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewInitializeAccountInstruction(accountPublicKey, mintPublicKey, ownerPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendInitializeAccount2Instruction(
	account string,
	mint string,
	owner string) {
	accountPublicKey := base.MustPublicKeyFromBase58(account)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewInitializeAccount2Instruction(ownerPublicKey, accountPublicKey, mintPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendInitializeAccount3Instruction(
	account string,
	mint string,
	owner string) {
	accountPublicKey := base.MustPublicKeyFromBase58(account)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewInitializeAccount3Instruction(ownerPublicKey, accountPublicKey, mintPublicKey).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendInitializeMintInstruction(
	decimals uint8,
	mint_authority string,
	freeze_authority string,
	mint string,
	options ...string) {
	mintAuthPublicKey := base.MustPublicKeyFromBase58(mint_authority)
	freezeAuthPublicKey := base.MustPublicKeyFromBase58(freeze_authority)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	inst := token.NewInitializeMintInstruction(decimals, mintAuthPublicKey, freezeAuthPublicKey, mintPublicKey).Build()
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		inst.SetProgramID(base.Token2022ProgramID)
	}
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendTokenInitializeMint2Instruction(
	decimals uint8,
	mint_authority string,
	freeze_authority string,
	mint string,
	options ...string) {
	mintAuthPublicKey := base.MustPublicKeyFromBase58(mint_authority)
	freezeAuthPublicKey := base.MustPublicKeyFromBase58(freeze_authority)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	inst := token.NewInitializeMint2Instruction(decimals, mintAuthPublicKey, freezeAuthPublicKey, mintPublicKey).Build()
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		inst.SetProgramID(base.Token2022ProgramID)
	}
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendTokenMintToInstruction(
	amount uint64,
	mint string,
	destination string,
	authority string,
	options ...string) {
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	destinationPublicKey := base.MustPublicKeyFromBase58(destination)
	authorityPublicKey := base.MustPublicKeyFromBase58(authority)
	inst := token.NewMintToInstruction(amount, mintPublicKey, destinationPublicKey, authorityPublicKey, []base.PublicKey{}).Build()
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		inst.SetProgramID(base.Token2022ProgramID)
	}
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendTokenMintToCheckedInstruction(
	amount uint64,
	decimals uint8,
	// Accounts:
	mint string,
	destination string,
	authority string,
	options ...string) {
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	destinationPublicKey := base.MustPublicKeyFromBase58(destination)
	authorityPublicKey := base.MustPublicKeyFromBase58(authority)
	inst := token.NewMintToCheckedInstruction(amount, decimals, mintPublicKey, destinationPublicKey, authorityPublicKey, []base.PublicKey{}).Build()
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		inst.SetProgramID(base.Token2022ProgramID)
	}
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendRevokeInstruction(
	source string,
	owner string) {
	sourcePublicKey := base.MustPublicKeyFromBase58(source)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewRevokeInstruction(sourcePublicKey, ownerPublicKey, []base.PublicKey{}).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendSetAuthorityInstruction(
	authority_type token.AuthorityType,
	new_authority string,
	// Accounts:
	subject string,
	authority string) {
	newAuthPublicKey := base.MustPublicKeyFromBase58(new_authority)
	subjectPublicKey := base.MustPublicKeyFromBase58(subject)
	authorityPublicKey := base.MustPublicKeyFromBase58(authority)
	inst := token.NewSetAuthorityInstruction(authority_type, newAuthPublicKey, subjectPublicKey, authorityPublicKey, []base.PublicKey{}).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendThawAccountInstruction(
	account string,
	mint string,
	authority string) {
	accountPublicKey := base.MustPublicKeyFromBase58(account)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	authorityPublicKey := base.MustPublicKeyFromBase58(authority)
	inst := token.NewThawAccountInstruction(accountPublicKey, mintPublicKey, authorityPublicKey, []base.PublicKey{}).Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendTokenTransferInstruction(
	amount uint64,
	source string,
	destination string,
	owner string,
	options ...string) {
	sourcePublicKey := base.MustPublicKeyFromBase58(source)
	destinationPublicKey := base.MustPublicKeyFromBase58(destination)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewTransferInstruction(amount, sourcePublicKey, destinationPublicKey, ownerPublicKey, []base.PublicKey{}).Build()
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		inst.SetProgramID(base.Token2022ProgramID)
	}
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendTokenTransferCheckedInstruction(amount uint64, decimals uint8,
	source string,
	mint string,
	destination string,
	owner string,
	options ...string) {
	sourcePublicKey := base.MustPublicKeyFromBase58(source)
	mintPublicKey := base.MustPublicKeyFromBase58(mint)
	destinationPublicKey := base.MustPublicKeyFromBase58(destination)
	ownerPublicKey := base.MustPublicKeyFromBase58(owner)
	inst := token.NewTransferCheckedInstruction(amount, decimals, sourcePublicKey, mintPublicKey, destinationPublicKey, ownerPublicKey, []base.PublicKey{}).Build()
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		inst.SetProgramID(base.Token2022ProgramID)
	}
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendAssociatedTokenAccountCreateInstruction(payer string, walletAddress string, splTokenMintAddress string, options ...string) {
	payerPublicKey := base.MustPublicKeyFromBase58(payer)
	walletPublicKey := base.MustPublicKeyFromBase58(walletAddress)
	splTokenMintPublicKey := base.MustPublicKeyFromBase58(splTokenMintAddress)
	create := associatedtokenaccount.NewCreateInstruction(payerPublicKey, walletPublicKey, splTokenMintPublicKey)
	if len(options) > 0 && options[0] == base.TOKEN2022 {
		create.SetTokenProgramID(base.Token2022ProgramID)
	}
	inst := create.Build()
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendInstruction(inst base.Instruction) {
	t.instructions = append(t.instructions, inst)
}

func (t *RawTransaction) AppendSigner(privateKeyHex string) {
	for _, s := range t.signers {
		if s == privateKeyHex {
			return
		}
	}
	t.signers = append(t.signers, privateKeyHex)
}

func (t *RawTransaction) Sign(base58 bool) (string, error) {
	h := base.Hash(base.MustPublicKeyFromBase58(t.blockHash))
	tx, err := base.NewTransaction(t.instructions, h, t.payer)
	if err != nil {
		return "", err
	}

	tx.Sign(func(key base.PublicKey) *base.PrivateKey {
		for _, signer := range t.signers {
			privateBytes, _ := hex.DecodeString(signer)
			privateKey := base.PrivateKey(privateBytes)
			if key.Equals(privateKey.PublicKey()) {
				return &privateKey
			}
		}
		return nil
	})

	if base58 {
		return tx.ToBase58()
	} else {
		return tx.ToBase64()
	}
}

func DecodeAndSign(txRaw string, signers []string, recentBlockHash string, base58 bool) (string, error) {
	tx := base.Transaction{}
	if base58 {
		err := tx.UnmarshalBase58(txRaw)
		if err != nil {
			return "", err
		}
	} else {
		err := tx.UnmarshalBase64(txRaw)
		if err != nil {
			return "", err
		}
	}
	tx.Signatures = make([]base.Signature, 0)
	tx.Message.RecentBlockhash = base.Hash(base.MustPublicKeyFromBase58(recentBlockHash))
	_, err := tx.Sign(func(key base.PublicKey) *base.PrivateKey {
		for _, signer := range signers {
			privateBytes, _ := hex.DecodeString(signer)
			privateKey := base.PrivateKey(privateBytes)
			if key.Equals(privateKey.PublicKey()) {
				return &privateKey
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if base58 {
		return tx.ToBase58()
	} else {
		return tx.ToBase64()
	}
}

type UnsignedTx struct {
	BizType string
	Data    string
	BizId   []string
}

func DecodeAndMultiSign(unsignedTx string, privateKey string, recentBlockHash string, base58 bool) (string, error) {
	ut := UnsignedTx{}
	err := json.Unmarshal([]byte(unsignedTx), &ut)
	if err != nil {
		return "", fmt.Errorf("unsignedTx error")
	}

	signers := make([]string, 0)
	signers = append(signers, privateKey)

	for _, s := range ut.BizId {
		seedBytes := bip39.NewSeed(ut.BizType, s)
		extra := ed25519.NewKeyFromSeed(seedBytes[:32])
		extraHex := hex.EncodeToString(base.PrivateKey(extra).Bytes())
		signers = append(signers, extraHex)
	}
	return DecodeAndSign(ut.Data, signers, recentBlockHash, base58)
}

func (t *RawTransaction) UnsignedTx() (string, error) {
	h := base.Hash(base.MustPublicKeyFromBase58(t.blockHash))
	tx, err := base.NewTransaction(t.instructions, h, t.payer)
	if err != nil {
		return "", err
	}

	txBytes, err := tx.Message.MarshalBinary()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(txBytes), nil
}

func GetSigningHash(hash, from, to, nonceAddress string, amount uint64) (string, error) {
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendAdvanceNonceInstruction(from, nonceAddress)
	rawTransaction.AppendTransferInstruction(amount, from, to)

	h := base.Hash(base.MustPublicKeyFromBase58(rawTransaction.blockHash))
	tx, err := base.NewTransaction(rawTransaction.instructions, h, rawTransaction.payer)
	if err != nil {
		return "", err
	}

	txHashByte, err := tx.Message.MarshalBinary()
	if err != nil {
		return "", err
	}

	txHash := hex.EncodeToString(txHashByte)
	return txHash, nil
}

func SignedTx(hash, from, to, nonceAddress string, amount uint64, signData string) (string, error) {
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendAdvanceNonceInstruction(from, nonceAddress)
	rawTransaction.AppendTransferInstruction(amount, from, to)

	h := base.Hash(base.MustPublicKeyFromBase58(rawTransaction.blockHash))
	tx, err := base.NewTransaction(rawTransaction.instructions, h, rawTransaction.payer)
	if err != nil {
		return "", err
	}

	signDataByte, err := hex.DecodeString(signData)
	if err != nil {
		return "", err
	}
	var signature base.Signature
	copy(signature[:], signDataByte)
	tx.Signatures = append(tx.Signatures, signature)

	return tx.ToBase58()
}
