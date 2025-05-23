package base

type Hash PublicKey
type Signature [64]byte
type Base58 []byte

type AccountsSettable interface {
	SetAccounts(accounts []*AccountMeta) error
}

type AccountsGettable interface {
	GetAccounts() (accounts []*AccountMeta)
}

type AccountMeta struct {
	PublicKey  PublicKey
	IsWritable bool
	IsSigner   bool
}

// Meta initializes a new AccountMeta with the provided pubKey.
func Meta(
	pubKey PublicKey,
) *AccountMeta {
	return &AccountMeta{
		PublicKey: pubKey,
	}
}

// WRITE sets IsWritable to true.
func (meta *AccountMeta) WRITE() *AccountMeta {
	meta.IsWritable = true
	return meta
}

// SIGNER sets IsSigner to true.
func (meta *AccountMeta) SIGNER() *AccountMeta {
	meta.IsSigner = true
	return meta
}

func NewAccountMeta(
	pubKey PublicKey,
	WRITE bool,
	SIGNER bool,
) *AccountMeta {
	return &AccountMeta{
		PublicKey:  pubKey,
		IsWritable: WRITE,
		IsSigner:   SIGNER,
	}
}

func (a *AccountMeta) less(act *AccountMeta) bool {
	if a.IsSigner != act.IsSigner {
		return a.IsSigner
	}
	if a.IsWritable != act.IsWritable {
		return a.IsWritable
	}
	return false
}

type AccountMetaSlice []*AccountMeta

func (slice *AccountMetaSlice) Append(account *AccountMeta) {
	*slice = append(*slice, account)
}

func (slice *AccountMetaSlice) SetAccounts(accounts []*AccountMeta) error {
	*slice = accounts
	return nil
}

func (slice AccountMetaSlice) GetAccounts() []*AccountMeta {
	out := make([]*AccountMeta, 0, len(slice))
	for i := range slice {
		if slice[i] != nil {
			out = append(out, slice[i])
		}
	}
	return out
}

// Get returns the AccountMeta at the desired index.
// If the index is not present, it returns nil.
func (slice AccountMetaSlice) Get(index int) *AccountMeta {
	if len(slice) > index {
		return slice[index]
	}
	return nil
}

// GetSigners returns the accounts that are signers.
func (slice AccountMetaSlice) GetSigners() []*AccountMeta {
	signers := make([]*AccountMeta, 0, len(slice))
	for _, ac := range slice {
		if ac.IsSigner {
			signers = append(signers, ac)
		}
	}
	return signers
}

// GetKeys returns the pubkeys of all AccountMeta.
func (slice AccountMetaSlice) GetKeys() PublicKeySlice {
	keys := make(PublicKeySlice, 0, len(slice))
	for _, ac := range slice {
		keys = append(keys, ac.PublicKey)
	}
	return keys
}

func (slice AccountMetaSlice) Len() int {
	return len(slice)
}

func (slice AccountMetaSlice) SplitFrom(index int) (AccountMetaSlice, AccountMetaSlice) {
	if index < 0 {
		panic("negative index")
	}
	if index == 0 {
		return AccountMetaSlice{}, slice
	}
	if index > len(slice)-1 {
		return slice, AccountMetaSlice{}
	}

	firstLen, secondLen := calcSplitAtLengths(len(slice), index)

	first := make(AccountMetaSlice, firstLen)
	copy(first, slice[:index])

	second := make(AccountMetaSlice, secondLen)
	copy(second, slice[index:])

	return first, second
}

func calcSplitAtLengths(total int, index int) (int, int) {
	if index == 0 {
		return 0, total
	}
	if index > total-1 {
		return total, 0
	}
	return index, total - index
}
