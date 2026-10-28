package trees

import (
	"math/rand"
	"slices"
	"testing"
)

func TestInvertTree(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{4, 2, 7, 1, 3, 6, 9}, []int{4, 7, 2, 9, 6, 3, 1}},
		{[]int{2, 1, 3}, []int{2, 3, 1}},
		{[]int{1, 2}, []int{1, Null, 2}},
		{nil, nil},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(*TreeNode) *TreeNode{"dfs": InvertTree, "bfs": InvertTreeBFS} {
			if got := ToLevelOrder(fn(FromLevelOrder(tt.in...))); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.in, got, tt.want)
			}
		}
	}
}

// Inverting reverses the inorder sequence, and inverting twice is a no-op.
func TestInvertTreeRandom(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for i := 0; i < 200; i++ {
		vals := ToLevelOrder(randomSmallTree(r, r.Intn(25), 0, 10))
		in := InorderRecursive(FromLevelOrder(vals...))
		slices.Reverse(in)
		if got := InorderRecursive(InvertTree(FromLevelOrder(vals...))); !slices.Equal(got, in) {
			t.Fatalf("inorder of inverted %v = %v, want %v", vals, got, in)
		}
		if got := InvertTreeBFS(InvertTree(FromLevelOrder(vals...))); !Equal(got, FromLevelOrder(vals...)) {
			t.Fatalf("double inversion changed %v", vals)
		}
	}
}
