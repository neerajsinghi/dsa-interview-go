package linkedlists

import (
	"math/rand"
	"slices"
	"testing"
)

func TestReverseKGroup(t *testing.T) {
	tests := []struct {
		vals []int
		k    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{2, 1, 4, 3, 5}},
		{[]int{1, 2, 3, 4, 5}, 3, []int{3, 2, 1, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 1, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 5, []int{5, 4, 3, 2, 1}},
		{[]int{1, 2, 3, 4, 5}, 6, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5, 6}, 3, []int{3, 2, 1, 6, 5, 4}},
		{nil, 2, nil},
	}
	impls := map[string]func(*ListNode, int) *ListNode{"brute": ReverseKGroupBrute, "inPlace": ReverseKGroup}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := ToSlice(fn(FromSlice(tt.vals), tt.k)); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v, %d) = %v, want %v", name, tt.vals, tt.k, got, tt.want)
			}
		}
	}
}

func TestReverseKGroupRandom(t *testing.T) {
	r := rand.New(rand.NewSource(8))
	for iter := 0; iter < 500; iter++ {
		vals := randomSlice(r, 25, 100)
		k := 1 + r.Intn(6)
		got := ToSlice(ReverseKGroup(FromSlice(vals), k))
		want := ToSlice(ReverseKGroupBrute(FromSlice(vals), k))
		if !slices.Equal(got, want) {
			t.Fatalf("ReverseKGroup(%v, %d) = %v, want %v", vals, k, got, want)
		}
	}
}
