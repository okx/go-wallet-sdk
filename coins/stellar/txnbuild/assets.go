package txnbuild

// Assets represents a list of Stellar assets. Importantly, it is sortable.
type Assets []Asset

func (s Assets) Len() int           { return len(s) }
func (s Assets) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Assets) Less(i, j int) bool { return s[i].LessThan(s[j]) }
