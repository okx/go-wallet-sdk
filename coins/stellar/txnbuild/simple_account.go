package txnbuild

// SimpleAccount is a minimal implementation of an Account.
type SimpleAccount struct {
	AccountID string
	Sequence  int64
}

// GetAccountID returns the Account ID.
func (sa *SimpleAccount) GetAccountID() string {
	return sa.AccountID
}

// IncrementSequenceNumber increments the internal record of the
// account's sequence number by 1.
func (sa *SimpleAccount) IncrementSequenceNumber() (int64, error) {
	sa.Sequence++
	return sa.Sequence, nil
}

// GetSequenceNumber returns the sequence number of the account.
func (sa *SimpleAccount) GetSequenceNumber() (int64, error) {
	return sa.Sequence, nil
}

// NewSimpleAccount is a factory method that creates a SimpleAccount from "accountID" and "sequence".
func NewSimpleAccount(accountID string, sequence int64) SimpleAccount {
	return SimpleAccount{accountID, sequence}
}

// ensure that SimpleAccount implements Account interface.
var _ Account = &SimpleAccount{}
