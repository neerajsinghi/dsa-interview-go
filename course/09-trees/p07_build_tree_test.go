package trees

import (
	"math/rand"
	"slices"
	"testing"
)

func TestBuildTree(t *testing.T) {
	tests := []struct{ pre, in, want []int }{
		{[]int{3, 9, 20, 15, 7}, []int{9, 3, 15, 20, 7}, []int{3, 9, 20, Null, Null, 15, 7}},
		{[]int{-1}, []int{-1}, []int{-1}},
		{nil, nil, nil},
		{[]int{1, 2, 3}, []int{3, 2, 1}, []int{1, 2, Null, 3}},
		{[]int{1, 2, 3}, []int{1, 2, 3}, []int{1, Null, 2, Null, 3}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int, []int) *TreeNode{"brute": BuildTreeBrute, "map": BuildTree} {
			if got := ToLevelOrder(fn(tt.pre, tt.in)); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v, %v) = %v, want %v", name, tt.pre, tt.in, got, tt.want)
			}
		}
	}
}

func TestBuildTreeRandom(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for i := 0; i < 300; i++ {
		orig := randomUniqueTree(r, r.Intn(30))
		pre, in := PreorderRecursive(orig), InorderRecursive(orig)
		if got := BuildTree(pre, in); !Equal(got, orig) {
			t.Fatalf("BuildTree(%v, %v) = %v, want %v", pre, in, ToLevelOrder(got), ToLevelOrder(orig))
		}
		if got := BuildTreeBrute(pre, in); !Equal(got, orig) {
			t.Fatalf("BuildTreeBrute(%v, %v) = %v, want %v", pre, in, ToLevelOrder(got), ToLevelOrder(orig))
		}
	}
}
