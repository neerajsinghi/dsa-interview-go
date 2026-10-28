// Package stacks covers the stack and the monotonic stack: last-in-first-out
// structures that let you match, undo, or "look back" at the most recent
// unresolved item in O(1).
package stacks

// Stack is a generic LIFO stack backed by a slice. The zero value is an empty
// stack ready to use:
//
//	var s Stack[int]
//	s.Push(1)
//	v, ok := s.Pop() // 1, true
//
// In interviews you will usually inline a plain slice instead
// (st = append(st, x); top := st[len(st)-1]; st = st[:len(st)-1]); this type
// exists so the problems read clearly and the operations are named.
type Stack[T any] struct {
	items []T
}

// Push adds v to the top of the stack. Amortised O(1).
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the top element. ok is false when the stack is empty.
func (s *Stack[T]) Pop() (v T, ok bool) {
	if len(s.items) == 0 {
		return v, false
	}
	last := len(s.items) - 1
	v = s.items[last]
	var zero T
	s.items[last] = zero // drop the reference so the GC can reclaim pointer types
	s.items = s.items[:last]
	return v, true
}

// Peek returns the top element without removing it.
func (s *Stack[T]) Peek() (v T, ok bool) {
	if len(s.items) == 0 {
		return v, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of elements.
func (s *Stack[T]) Len() int { return len(s.items) }

// IsEmpty reports whether the stack has no elements.
func (s *Stack[T]) IsEmpty() bool { return len(s.items) == 0 }
