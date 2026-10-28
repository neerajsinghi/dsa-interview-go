package design

import "cmp"

// Problem: Min Stack / Max Stack
// Difficulty: Medium
//
// MinStack: a stack supporting Push, Pop, Top and GetMin, all in O(1).
// MaxStack: a stack supporting Push, Pop, Top, PeekMax and PopMax, where
// PopMax removes the largest element (the top-most one if there are ties).
//
// Examples:
//
//	MinStack: push -2, push 0, push -3, min -> -3, pop, top -> 0, min -> -2
//	MaxStack: push 5, push 1, push 5, popMax -> 5 (the top one), top -> 1
//
// Reference: https://leetcode.com/problems/min-stack/
// Reference: https://leetcode.com/problems/max-stack/
//
// Pattern: store the running aggregate WITH each element. Each frame
// remembers min(value, min-of-everything-below), so popping restores the
// previous minimum for free.
//
//	push 3   [ (3,3) ]
//	push 1   [ (3,3) (1,1) ]
//	push 2   [ (3,3) (1,1) (2,1) ]   GetMin = top.min = 1
//
// PopMax is the twist: buffer elements above the max, drop it, push them
// back — O(n). (An O(log n) version keeps a doubly linked list plus a heap
// with lazy deletion; mention it, but the buffered version is what most
// interviews accept.)
//
// Time: O(1) for every operation except PopMax, O(n). Space: O(n).
//
// Not safe for concurrent use; wrap in a mutex if shared.

type frame[T cmp.Ordered] struct {
	val, agg T // agg = min (MinStack) or max (MaxStack) of this frame and below
}

// MinStack is a stack that also reports its minimum in O(1).
type MinStack[T cmp.Ordered] struct {
	frames []frame[T]
}

// Push adds v on top.
func (s *MinStack[T]) Push(v T) {
	agg := v
	if len(s.frames) > 0 {
		agg = min(v, s.frames[len(s.frames)-1].agg)
	}
	s.frames = append(s.frames, frame[T]{v, agg})
}

// Pop removes and returns the top element; ok is false if empty.
func (s *MinStack[T]) Pop() (v T, ok bool) {
	if len(s.frames) == 0 {
		return v, false
	}
	top := s.frames[len(s.frames)-1]
	s.frames = s.frames[:len(s.frames)-1]
	return top.val, true
}

// Top returns the top element without removing it.
func (s *MinStack[T]) Top() (v T, ok bool) {
	if len(s.frames) == 0 {
		return v, false
	}
	return s.frames[len(s.frames)-1].val, true
}

// GetMin returns the smallest element currently on the stack.
func (s *MinStack[T]) GetMin() (v T, ok bool) {
	if len(s.frames) == 0 {
		return v, false
	}
	return s.frames[len(s.frames)-1].agg, true
}

// Len reports the number of elements.
func (s *MinStack[T]) Len() int { return len(s.frames) }

// MaxStack is a stack that can also remove its maximum.
type MaxStack[T cmp.Ordered] struct {
	frames []frame[T]
}

// Push adds v on top.
func (s *MaxStack[T]) Push(v T) {
	agg := v
	if len(s.frames) > 0 {
		agg = max(v, s.frames[len(s.frames)-1].agg)
	}
	s.frames = append(s.frames, frame[T]{v, agg})
}

// Pop removes and returns the top element.
func (s *MaxStack[T]) Pop() (v T, ok bool) {
	if len(s.frames) == 0 {
		return v, false
	}
	top := s.frames[len(s.frames)-1]
	s.frames = s.frames[:len(s.frames)-1]
	return top.val, true
}

// Top returns the top element without removing it.
func (s *MaxStack[T]) Top() (v T, ok bool) {
	if len(s.frames) == 0 {
		return v, false
	}
	return s.frames[len(s.frames)-1].val, true
}

// PeekMax returns the largest element without removing it.
func (s *MaxStack[T]) PeekMax() (v T, ok bool) {
	if len(s.frames) == 0 {
		return v, false
	}
	return s.frames[len(s.frames)-1].agg, true
}

// PopMax removes and returns the largest element (top-most among ties).
func (s *MaxStack[T]) PopMax() (v T, ok bool) {
	target, ok := s.PeekMax()
	if !ok {
		return v, false
	}
	var buffer []T
	for {
		top, _ := s.Pop()
		if top == target {
			break
		}
		buffer = append(buffer, top)
	}
	// Re-push in original order so every frame's aggregate is recomputed.
	for i := len(buffer) - 1; i >= 0; i-- {
		s.Push(buffer[i])
	}
	return target, true
}

// Len reports the number of elements.
func (s *MaxStack[T]) Len() int { return len(s.frames) }
