package linkedlists

import (
	"math/rand"
	"slices"
	"testing"
)

func TestRemoveNthFromEnd(t *testing.T) {
	tests := []struct {
		vals []int
		n    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3, 5}},
		{[]int{1}, 1, nil},
		{[]int{1, 2}, 1, []int{1}},
		{[]int{1, 2}, 2, []int{2}},
		{[]int{1, 2, 3}, 3, []int{2, 3}},
	}
	impls := map[string]func(*ListNode, int) *ListNode{
		"twoPass": RemoveNthFromEndTwoPass, "onePass": RemoveNthFromEnd,
	}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := ToSlice(fn(FromSlice(tt.vals), tt.n)); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v, %d) = %v, want %v", name, tt.vals, tt.n, got, tt.want)
			}
		}
	}
}

func TestRemoveNthFromEndRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for iter := 0; iter < 300; iter++ {
		vals := randomSlice(r, 20, 100)
		if len(vals) == 0 {
			continue
		}
		n := 1 + r.Intn(len(vals))
		got := ToSlice(RemoveNthFromEnd(FromSlice(vals), n))
		want := ToSlice(RemoveNthFromEndTwoPass(FromSlice(vals), n))
		if !slices.Equal(got, want) {
			t.Fatalf("RemoveNthFromEnd(%v, %d) = %v, want %v", vals, n, got, want)
		}
	}
}
