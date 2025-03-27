package set

import (
	"sync"

	"golang.org/x/exp/constraints"
)

// safeSet is a simple, thread-safe set implementation. Note that it *must* be
// created via NewSafeSet.
type safeSet[T constraints.Ordered] struct {
	Set[T]
	lock sync.RWMutex
}

func NewSafeSet[T constraints.Ordered](capacity int) *safeSet[T] {
	return &safeSet[T]{
		Set:  NewSet[T](capacity),
		lock: sync.RWMutex{},
	}
}

func (s *safeSet[T]) Add(item T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Set.Add(item)
}

func (s *safeSet[T]) AddSlice(items []T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Set.AddSlice(items)
}

func (s *safeSet[T]) Remove(item T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Set.Remove(item)
}

func (s *safeSet[T]) Contains(item T) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Set.Contains(item)
}

func (s *safeSet[T]) Slice() []T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Set.Slice()
}
