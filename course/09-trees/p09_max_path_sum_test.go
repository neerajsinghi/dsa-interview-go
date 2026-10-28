package trees

import (
	"math/rand"
	"testing"
)

func TestMaxPathSum(t *testing.T) {
	tests := []struct {
		vals []int
		want int
	}{
		{[]int{1, 2, 3}, 6},
		{[]int{-10, 9, 20, Null, Null, 15, 7}, 42},
		{[]int{-3}, -3},
		{[]int{-2, -1}, -1},
		{[]int{2, -1}, 2},
		{[]int{5, 4, 8, 11, Null, 13, 4, 7, 2, Null, Null, Null, 1}, 48},
		{[]int{1, -2, -3, 1, 3, -2, Null, -1}, 3},
	}
	for _, tt := range tests {
		root := FromLevelOrder(tt.vals...)
		for name, fn := range map[string]func(*TreeNode) int{"brute": MaxPathSumBrute, "onePass": MaxPathSum} {
			if got := fn(root); got != tt.want {
				t.Errorf("%s(%v) = %d, want %d", name, tt.vals, got, tt.want)
			}
		}
	}
}

func TestMaxPathSumRandom(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	for i := 0; i < 300; i++ {
		root := randomSmallTree(r, 1+r.Intn(25), -10, 10)
		if got, want := MaxPathSum(root), MaxPathSumBrute(root); got != want {
			t.Fatalf("MaxPathSum(%v) = %d, want %d", ToLevelOrder(root), got, want)
		}
	}
}
