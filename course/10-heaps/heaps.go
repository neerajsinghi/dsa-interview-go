// Package heaps covers binary heaps and priority queues: top-k selection,
// k-way merging, scheduling by priority, and running medians.
//
// Two implementations are provided:
//
//   - Heap[T]: a generic priority queue on top of container/heap. This is
//     what you'd use in real Go code and in most interviews.
//   - ArrayHeap[T]: a hand-rolled min-heap with explicit siftUp/siftDown, to
//     show what container/heap does for you.
package heaps

import (
	"cmp"
	"container/heap"
)

// ---------------------------------------------------------------------------
// Generic priority queue on container/heap
// ---------------------------------------------------------------------------

// heapSlice adapts a slice + less function to heap.Interface. container/heap
// calls Push/Pop (the `any` ones) only to grow/shrink the slice; it does the
// sifting itself using Less and Swap.
type heapSlice[T any] struct {
	items []T
	less  func(a, b T) bool
}

func (h *heapSlice[T]) Len() int           { return len(h.items) }
func (h *heapSlice[T]) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *heapSlice[T]) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *heapSlice[T]) Push(x any)         { h.items = append(h.items, x.(T)) }
func (h *heapSlice[T]) Pop() any {
	n := len(h.items)
	x := h.items[n-1]
	var zero T
	h.items[n-1] = zero // let the GC reclaim pointer values
	h.items = h.items[:n-1]
	return x
}

// Heap is a generic binary heap. The element for which less is "smallest"
// is at the top: pass a < b for a min-heap, a > b for a max-heap.
type Heap[T any] struct {
	h *heapSlice[T]
}

// New returns an empty heap ordered by less.
func New[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{h: &heapSlice[T]{less: less}}
}

// NewMin returns an empty min-heap of ordered values.
func NewMin[T cmp.Ordered]() *Heap[T] { return New(func(a, b T) bool { return a < b }) }

// NewMax returns an empty max-heap of ordered values.
func NewMax[T cmp.Ordered]() *Heap[T] { return New(func(a, b T) bool { return a > b }) }

// FromSlice builds a heap from items in O(n) (heap.Init). It takes ownership
// of the slice.
func FromSlice[T any](items []T, less func(a, b T) bool) *Heap[T] {
	h := &heapSlice[T]{items: items, less: less}
	heap.Init(h)
	return &Heap[T]{h: h}
}

// Push adds v. O(log n).
func (h *Heap[T]) Push(v T) { heap.Push(h.h, v) }

// Pop removes and returns the top element. O(log n). Panics if empty.
func (h *Heap[T]) Pop() T { return heap.Pop(h.h).(T) }

// Peek returns the top element without removing it. Panics if empty.
func (h *Heap[T]) Peek() T { return h.h.items[0] }

// Len returns the number of elements.
func (h *Heap[T]) Len() int { return h.h.Len() }

// ---------------------------------------------------------------------------
// Hand-rolled array heap (for teaching)
// ---------------------------------------------------------------------------

// ArrayHeap is a min-heap stored in a slice. For node i:
//
//	parent = (i-1)/2    left = 2i+1    right = 2i+2
//
// The invariant: every parent is <= its children, so the minimum is a[0].
type ArrayHeap[T cmp.Ordered] struct {
	a []T
}

// Push appends v at the bottom and sifts it up. O(log n).
func (h *ArrayHeap[T]) Push(v T) {
	h.a = append(h.a, v)
	h.siftUp(len(h.a) - 1)
}

// Pop removes the minimum: move the last element to the root and sift it
// down. O(log n). ok is false when the heap is empty.
func (h *ArrayHeap[T]) Pop() (v T, ok bool) {
	if len(h.a) == 0 {
		return v, false
	}
	v = h.a[0]
	last := len(h.a) - 1
	h.a[0] = h.a[last]
	h.a = h.a[:last]
	h.siftDown(0)
	return v, true
}

// Peek returns the minimum without removing it.
func (h *ArrayHeap[T]) Peek() (v T, ok bool) {
	if len(h.a) == 0 {
		return v, false
	}
	return h.a[0], true
}

// Len returns the number of elements.
func (h *ArrayHeap[T]) Len() int { return len(h.a) }

// Heapify builds a heap from xs in O(n) by sifting down every internal node,
// starting from the last one. (Pushing one at a time would be O(n log n).)
func Heapify[T cmp.Ordered](xs []T) *ArrayHeap[T] {
	h := &ArrayHeap[T]{a: xs}
	for i := len(xs)/2 - 1; i >= 0; i-- {
		h.siftDown(i)
	}
	return h
}

func (h *ArrayHeap[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.a[parent] <= h.a[i] {
			return // heap property holds
		}
		h.a[parent], h.a[i] = h.a[i], h.a[parent]
		i = parent
	}
}

func (h *ArrayHeap[T]) siftDown(i int) {
	n := len(h.a)
	for {
		smallest := i
		if l := 2*i + 1; l < n && h.a[l] < h.a[smallest] {
			smallest = l
		}
		if r := 2*i + 2; r < n && h.a[r] < h.a[smallest] {
			smallest = r
		}
		if smallest == i {
			return
		}
		h.a[i], h.a[smallest] = h.a[smallest], h.a[i]
		i = smallest
	}
}
