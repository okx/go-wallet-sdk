/*
 * Copyright 2016 Stellar Development Foundation and contributors.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file includes portions of third-party code from [https://github.com/stellar/go].
 * The original code is licensed under the Apache License 2.0.
 */

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
