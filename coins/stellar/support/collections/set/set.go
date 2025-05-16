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

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](capacity int) Set[T] {
	return make(map[T]struct{}, capacity)
}

func (set Set[T]) Add(item T) {
	set[item] = struct{}{}
}

func (set Set[T]) AddSlice(items []T) {
	for _, item := range items {
		set[item] = struct{}{}
	}
}

func (set Set[T]) Remove(item T) {
	delete(set, item)
}

func (set Set[T]) Contains(item T) bool {
	_, ok := set[item]
	return ok
}

func (set Set[T]) Slice() []T {
	slice := make([]T, 0, len(set))
	for key := range set {
		slice = append(slice, key)
	}
	return slice
}

var _ ISet[int] = (*Set[int])(nil) // ensure conformity to the interface
