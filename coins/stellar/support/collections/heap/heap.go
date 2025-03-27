package heap

import "container/heap"

// A Heap is a min-heap backed by a slice.
type Heap[E any] struct {
	s sliceHeap[E]
}

// New constructs a new Heap with a comparison function.
func New[E any](less func(E, E) bool, capacity int) *Heap[E] {
	return &Heap[E]{sliceHeap[E]{
		s:    make([]E, 0, capacity),
		less: less,
	}}
}

// Push pushes an element onto the heap. The complexity is O(log n)
// where n = h.Len().
func (h *Heap[E]) Push(elem E) {
	heap.Push(&h.s, elem)
}

// Pop removes and returns the minimum element (according to the less function)
// from the heap. Pop panics if the heap is empty.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Pop() E {
	return heap.Pop(&h.s).(E)
}

// Peek returns the minimum element (according to the less function) in the heap.
// Peek panics if the heap is empty.
// The complexity is O(1).
func (h *Heap[E]) Peek() E {
	return h.s.s[0]
}

// Len returns the number of elements in the heap.
func (h *Heap[E]) Len() int {
	return len(h.s.s)
}

// sliceHeap just exists to use the existing heap.Interface as the
// implementation of Heap.
type sliceHeap[E any] struct {
	s    []E
	less func(E, E) bool
}

func (s *sliceHeap[E]) Len() int { return len(s.s) }

func (s *sliceHeap[E]) Swap(i, j int) {
	s.s[i], s.s[j] = s.s[j], s.s[i]
}

func (s *sliceHeap[E]) Less(i, j int) bool {
	return s.less(s.s[i], s.s[j])
}

func (s *sliceHeap[E]) Push(x interface{}) {
	s.s = append(s.s, x.(E))
}

func (s *sliceHeap[E]) Pop() interface{} {
	var zero E
	e := s.s[len(s.s)-1]
	// avoid memory leak by clearing out popped value in slice
	s.s[len(s.s)-1] = zero
	s.s = s.s[:len(s.s)-1]
	return e
}
