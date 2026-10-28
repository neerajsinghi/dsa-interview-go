package trees

import (
	"math/rand"
	"testing"
)

func TestDiameter(t *testing.T) {
	tests := []struct {
		vals []int
		want int
	}{
		{nil, 0},
		{[]int{1}, 0},
		{[]int{1, 2}, 1},
		{[]int{1, 2, 3, 4, 5}, 3},
		// Longest path avoids the root: it runs through node 2's subtrees.
		{[]int{1, 2, Null, 3, 4, 5, Null, Null, 6, 7, Null, Null, 8}, 6},
	}
	for _, tt := range tests {
		root := FromLevelOrder(tt.vals...)
		for name, fn := range map[string]func(*TreeNode) int{"brute": DiameterBrute, "onePass": Diameter} {
			if got := fn(root); got != tt.want {
				t.Errorf("%s(%v) = %d, want %d", name, tt.vals, got, tt.want)
			}
		}
	}
}

func TestDiameterRandom(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for i := 0; i < 300; i++ {
		root := randomSmallTree(r, r.Intn(30), 0, 5)
		if got, want := Diameter(root), DiameterBrute(root); got != want {
			t.Fatalf("Diameter(%v) = %d, want %d", ToLevelOrder(root), got, want)
		}
	}
}
