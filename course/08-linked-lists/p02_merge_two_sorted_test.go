package linkedlists

import (
	"math/rand"
	"slices"
	"testing"
)

func TestMergeTwoLists(t *testing.T) {
	tests := []struct{ a, b, want []int }{
		{[]int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{nil, nil, nil},
		{nil, []int{0}, []int{0}},
		{[]int{5}, nil, []int{5}},
		{[]int{1, 2, 3}, []int{4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{4, 5, 6}, []int{1, 2, 3}, []int{1, 2, 3, 4, 5, 6}},
	}
	for _, tt := range tests {
		if got := ToSlice(MergeTwoLists(FromSlice(tt.a), FromSlice(tt.b))); !slices.Equal(got, tt.want) {
			t.Errorf("MergeTwoLists(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMergeTwoListsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for iter := 0; iter < 300; iter++ {
		a, b := randomSlice(r, 15, 20), randomSlice(r, 15, 20)
		slices.Sort(a)
		slices.Sort(b)
		want := slices.Sorted(slices.Values(append(slices.Clone(a), b...)))
		if got := ToSlice(MergeTwoLists(FromSlice(a), FromSlice(b))); !slices.Equal(got, want) {
			t.Fatalf("MergeTwoLists(%v, %v) = %v, want %v", a, b, got, want)
		}
	}
}
