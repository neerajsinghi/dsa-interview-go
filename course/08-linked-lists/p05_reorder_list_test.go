package linkedlists

import (
	"math/rand"
	"slices"
	"testing"
)

func TestReorderList(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{1, 2, 3, 4}, []int{1, 4, 2, 3}},
		{[]int{1, 2, 3, 4, 5}, []int{1, 5, 2, 4, 3}},
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{1, 2}},
		{[]int{1, 2, 3}, []int{1, 3, 2}},
		{nil, nil},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(*ListNode){"brute": ReorderListBrute, "inPlace": ReorderList} {
			head := FromSlice(tt.in)
			fn(head)
			if got := ToSlice(head); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.in, got, tt.want)
			}
		}
	}
}

func TestReorderListRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for iter := 0; iter < 300; iter++ {
		vals := randomSlice(r, 25, 50)
		a, b := FromSlice(vals), FromSlice(vals)
		ReorderList(a)
		ReorderListBrute(b)
		if got, want := ToSlice(a), ToSlice(b); !slices.Equal(got, want) {
			t.Fatalf("ReorderList(%v) = %v, want %v", vals, got, want)
		}
	}
}
