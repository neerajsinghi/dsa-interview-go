package design

import "iter"

// Problem: Iterator patterns — Peeking Iterator and Flatten Nested List
// Difficulty: Medium
//
// PeekingIterator wraps any iterator and adds Peek(), which returns the next
// element WITHOUT consuming it.
//
// NestedIterator walks a nested list such as [1,[4,[6]],[]] and yields its
// integers in order: 1, 4, 6. Lists may be empty or deeply nested.
//
// Examples:
//
//	peek over [1,2,3]: Next 1, Peek 2, Next 2, Next 3, HasNext false
//	flatten [[1,1],2,[1,1]] -> 1 1 2 1 1
//
// Reference: https://leetcode.com/problems/peeking-iterator/
// Reference: https://leetcode.com/problems/flatten-nested-list-iterator/
//
// Pattern:
//   - Peeking: cache one element ahead. Peek fills the cache; Next drains
//     it first.
//   - Flatten: an explicit stack of (list, position) frames replaces
//     recursion so the iterator can pause between elements. HasNext does the
//     work of skipping empty lists and descending into sub-lists, so Next is
//     trivial and HasNext is idempotent.
//   - Go 1.23 range-over-func: the same walk as an iter.Seq is just a
//     recursive function calling yield — compare how much shorter it is.
//
// Time: O(1) amortized per element. Space: O(depth) for the stack.
//
// Not safe for concurrent use — iterators are inherently single-consumer.

// Iterator is the classic pull-style interface used in these problems.
type Iterator[T any] interface {
	HasNext() bool
	Next() T
}

// SliceIterator iterates over a slice.
type SliceIterator[T any] struct {
	items []T
	pos   int
}

// NewSliceIterator returns an iterator over items.
func NewSliceIterator[T any](items []T) *SliceIterator[T] {
	return &SliceIterator[T]{items: items}
}

// HasNext reports whether Next will succeed.
func (s *SliceIterator[T]) HasNext() bool { return s.pos < len(s.items) }

// Next returns the next element. It panics when exhausted, like the
// underlying slice index would.
func (s *SliceIterator[T]) Next() T {
	v := s.items[s.pos]
	s.pos++
	return v
}

// PeekingIterator adds Peek to any Iterator.
type PeekingIterator[T any] struct {
	it     Iterator[T]
	peeked bool
	cached T
}

// NewPeekingIterator wraps it.
func NewPeekingIterator[T any](it Iterator[T]) *PeekingIterator[T] {
	return &PeekingIterator[T]{it: it}
}

// Peek returns the next element without consuming it.
func (p *PeekingIterator[T]) Peek() T {
	if !p.peeked {
		p.cached = p.it.Next()
		p.peeked = true
	}
	return p.cached
}

// Next consumes and returns the next element.
func (p *PeekingIterator[T]) Next() T {
	if p.peeked {
		p.peeked = false
		v := p.cached
		var zero T
		p.cached = zero // don't retain a reference the caller now owns
		return v
	}
	return p.it.Next()
}

// HasNext reports whether Next will succeed.
func (p *PeekingIterator[T]) HasNext() bool { return p.peeked || p.it.HasNext() }

// NestedInteger is either a single integer or a list of NestedIntegers.
type NestedInteger struct {
	isInt bool
	val   int
	list  []NestedInteger
}

// Int makes an integer element.
func Int(v int) NestedInteger { return NestedInteger{isInt: true, val: v} }

// List makes a list element.
func List(items ...NestedInteger) NestedInteger { return NestedInteger{list: items} }

// NestedIterator flattens a nested list lazily.
type NestedIterator struct {
	stack []nestedFrame
}

type nestedFrame struct {
	list []NestedInteger
	pos  int
}

// NewNestedIterator iterates over the top-level list items.
func NewNestedIterator(items []NestedInteger) *NestedIterator {
	return &NestedIterator{stack: []nestedFrame{{list: items}}}
}

// HasNext advances to the next integer (if any) without consuming it.
func (n *NestedIterator) HasNext() bool {
	for len(n.stack) > 0 {
		top := &n.stack[len(n.stack)-1]
		if top.pos == len(top.list) {
			n.stack = n.stack[:len(n.stack)-1] // finished this list
			continue
		}
		cur := top.list[top.pos]
		if cur.isInt {
			return true // positioned on an integer
		}
		top.pos++ // step past the sub-list in the parent, then descend
		// Note: append may reallocate n.stack, so `top` is not used after this.
		n.stack = append(n.stack, nestedFrame{list: cur.list})
	}
	return false
}

// Next returns the next integer. Call HasNext first.
func (n *NestedIterator) Next() int {
	if !n.HasNext() {
		panic("design: NestedIterator.Next called with no elements left")
	}
	top := &n.stack[len(n.stack)-1]
	v := top.list[top.pos].val
	top.pos++
	return v
}

// All returns a Go 1.23 push iterator over the same integers:
//
//	for v := range design.All(items) { ... }
func All(items []NestedInteger) iter.Seq[int] {
	return func(yield func(int) bool) {
		var walk func([]NestedInteger) bool
		walk = func(list []NestedInteger) bool {
			for _, x := range list {
				if x.isInt {
					if !yield(x.val) {
						return false // consumer broke out of the loop
					}
				} else if !walk(x.list) {
					return false
				}
			}
			return true
		}
		walk(items)
	}
}
