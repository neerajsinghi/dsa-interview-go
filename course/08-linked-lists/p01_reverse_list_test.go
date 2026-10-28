package linkedlists

import (
	"math/rand"
	"slices"
	"testing"
)

func TestReverseList(t *testing.T) {
	impls := map[string]func(*ListNode) *ListNode{
		"iterative": ReverseList, "recursive": ReverseListRecursive,
	}
	tests := [][]int{nil, {1}, {1, 2}, {1, 2, 3, 4, 5}, {3, 3, 1}}
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		tests = append(tests, randomSlice(r, 30, 10))
	}
	for _, vals := range tests {
		want := slices.Clone(vals)
		slices.Reverse(want)
		for name, fn := range impls {
			if got := ToSlice(fn(FromSlice(vals))); !slices.Equal(got, want) {
				t.Errorf("%s(%v) = %v, want %v", name, vals, got, want)
			}
		}
	}
}
