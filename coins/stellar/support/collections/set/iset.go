package set

type ISet[T comparable] interface {
	Add(item T)
	AddSlice(items []T)
	Remove(item T)
	Contains(item T) bool
	Slice() []T
}

var _ ISet[int] = (*Set[int])(nil) // ensure conformity to the interface
var _ ISet[int] = (*safeSet[int])(nil)
