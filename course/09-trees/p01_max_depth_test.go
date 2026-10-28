package trees

import (
	"math/rand"
	"testing"
)

func TestMaxDepth(t *testing.T) {
	tests := []struct {
		vals []int
		want int
	}{
		{nil, 0},
		{[]int{1}, 1},
		{[]int{3, 9, 20, Null, Null, 15, 7}, 3},
		{[]int{1, Null, 2}, 2},
		{[]int{1, 2, Null, 3, Null, 4}, 4},
	}
	for _, tt := range tests {
		root := FromLevelOrder(tt.vals...)
		if got := MaxDepth(root); got != tt.want {
			t.Errorf("MaxDepth(%v) = %d, want %d", tt.vals, got, tt.want)
		}
		if got := MaxDepthBFS(root); got != tt.want {
			t.Errorf("MaxDepthBFS(%v) = %d, want %d", tt.vals, got, tt.want)
		}
	}
}

func TestMaxDepthRandom(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 300; i++ {
		root := randomSmallTree(r, r.Intn(30), 0, 5)
		if a, b := MaxDepth(root), MaxDepthBFS(root); a != b {
			t.Fatalf("MaxDepth = %d, MaxDepthBFS = %d for %v", a, b, ToLevelOrder(root))
		}
	}
}
