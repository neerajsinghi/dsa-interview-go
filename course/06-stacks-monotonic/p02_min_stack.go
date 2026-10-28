package stacks

// Problem: Min Stack
// Difficulty: Medium
//
// Design a stack that supports Push, Pop, Top and GetMin, where GetMin returns
// the smallest element currently in the stack. Every operation must run in
// constant time. Pop, Top and GetMin are only called on a non-empty stack.
//
// Example:
//
//	Push(5) Push(2) Push(7)  GetMin() -> 2
//	Pop()                    GetMin() -> 2
//	Pop()                    GetMin() -> 5
//
// Reference: https://leetcode.com/problems/min-stack/
//
// Pattern: store derived state alongside each element. Each entry remembers
// the minimum of itself and everything beneath it, so popping automatically
// "restores" the previous minimum.
//
// Time: O(1) for every operation. Space: O(n).

type minEntry struct {
	val int
	min int // minimum of this entry and all entries below it
}

// MinStack is a stack of ints with O(1) GetMin.
type MinStack struct {
	items []minEntry
}

// NewMinStack returns an empty MinStack.
func NewMinStack() *MinStack { return &MinStack{} }

// Push adds v to the top.
func (s *MinStack) Push(v int) {
	m := v
	if n := len(s.items); n > 0 {
		m = min(m, s.items[n-1].min)
	}
	s.items = append(s.items, minEntry{val: v, min: m})
}

// Pop removes the top element. It panics on an empty stack.
func (s *MinStack) Pop() {
	s.items = s.items[:len(s.items)-1]
}

// Top returns the top element. It panics on an empty stack.
func (s *MinStack) Top() int { return s.items[len(s.items)-1].val }

// GetMin returns the smallest element. It panics on an empty stack.
func (s *MinStack) GetMin() int { return s.items[len(s.items)-1].min }

// Len returns the number of elements.
func (s *MinStack) Len() int { return len(s.items) }
