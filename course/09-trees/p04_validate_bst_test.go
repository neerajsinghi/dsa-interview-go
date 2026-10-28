package trees

import (
	"math"
	"math/rand"
	"testing"
)

func TestIsValidBST(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want bool
	}{
		{"empty", nil, true},
		{"single", []int{1}, true},
		{"simple", []int{2, 1, 3}, true},
		{"right child too small", []int{5, 1, 4, Null, Null, 3, 6}, false},
		{"grandchild violates root", []int{5, 4, 6, Null, Null, 3, 7}, false},
		{"duplicate", []int{2, 2}, false},
		{"duplicate right", []int{2, Null, 2}, false},
		{"bigger", []int{8, 4, 12, 2, 6, 10, 14, 1, 3, 5, 7}, true},
		{"extreme values", []int{math.MinInt + 1, Null, math.MaxInt}, true},
	}
	impls := map[string]func(*TreeNode) bool{
		"brute": IsValidBSTBrute, "bounds": IsValidBST, "inorder": IsValidBSTInorder,
	}
	for _, tt := range tests {
		root := FromLevelOrder(tt.vals...)
		for name, fn := range impls {
			if got := fn(root); got != tt.want {
				t.Errorf("%s/%s = %v, want %v", tt.name, name, got, tt.want)
			}
		}
	}
}

// insertBST builds BSTs for the randomized test so that many inputs are valid.
func insertBST(root *TreeNode, v int) *TreeNode {
	if root == nil {
		return &TreeNode{Val: v}
	}
	if v < root.Val {
		root.Left = insertBST(root.Left, v)
	} else {
		root.Right = insertBST(root.Right, v)
	}
	return root
}

func TestIsValidBSTRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for i := 0; i < 500; i++ {
		var root *TreeNode
		if i%2 == 0 {
			for j := r.Intn(15); j > 0; j-- {
				root = insertBST(root, r.Intn(20)) // duplicates go right -> invalid
			}
		} else {
			root = randomSmallTree(r, r.Intn(6), 0, 10)
		}
		want := IsValidBSTBrute(root)
		if got := IsValidBST(root); got != want {
			t.Fatalf("IsValidBST(%v) = %v, want %v", ToLevelOrder(root), got, want)
		}
		if got := IsValidBSTInorder(root); got != want {
			t.Fatalf("IsValidBSTInorder(%v) = %v, want %v", ToLevelOrder(root), got, want)
		}
	}
}
