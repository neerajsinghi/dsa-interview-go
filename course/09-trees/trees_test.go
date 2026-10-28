package trees

import (
	"math/rand"
	"slices"
	"testing"
)

// randomTree builds a random tree of n nodes with values drawn from vals
// (called once per node, in creation order).
func randomTree(r *rand.Rand, n int, val func() int) *TreeNode {
	if n == 0 {
		return nil
	}
	root := &TreeNode{Val: val()}
	nodes := []*TreeNode{root}
	for len(nodes) < n {
		parent := nodes[r.Intn(len(nodes))]
		child := &TreeNode{Val: val()}
		if r.Intn(2) == 0 {
			if parent.Left != nil {
				continue
			}
			parent.Left = child
		} else {
			if parent.Right != nil {
				continue
			}
			parent.Right = child
		}
		nodes = append(nodes, child)
	}
	return root
}

// randomUniqueTree builds a random tree whose n values are distinct.
func randomUniqueTree(r *rand.Rand, n int) *TreeNode {
	perm := r.Perm(n * 3)
	i := 0
	return randomTree(r, n, func() int { i++; return perm[i-1] - n })
}

// randomSmallTree builds a random tree with values in [lo, hi).
func randomSmallTree(r *rand.Rand, n, lo, hi int) *TreeNode {
	return randomTree(r, n, func() int { return lo + r.Intn(hi-lo) })
}

func TestLevelOrderRoundTrip(t *testing.T) {
	tests := [][]int{
		nil,
		{1},
		{1, 2, 3},
		{1, 2, 3, Null, 4},
		{1, Null, 2, Null, 3},
		{3, 9, 20, Null, Null, 15, 7},
	}
	for _, vals := range tests {
		if got := ToLevelOrder(FromLevelOrder(vals...)); !slices.Equal(got, vals) {
			t.Errorf("round trip %v = %v", vals, got)
		}
	}
	r := rand.New(rand.NewSource(9))
	for i := 0; i < 200; i++ {
		tree := randomSmallTree(r, r.Intn(20), -5, 5)
		if back := FromLevelOrder(ToLevelOrder(tree)...); !Equal(back, tree) {
			t.Fatalf("random round trip failed for %v", ToLevelOrder(tree))
		}
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		a, b []int
		want bool
	}{
		{nil, nil, true},
		{[]int{1}, nil, false},
		{[]int{1, 2, 3}, []int{1, 2, 3}, true},
		{[]int{1, 2}, []int{1, Null, 2}, false},
		{[]int{1, 2, 1}, []int{1, 1, 2}, false},
	}
	for _, tt := range tests {
		if got := Equal(FromLevelOrder(tt.a...), FromLevelOrder(tt.b...)); got != tt.want {
			t.Errorf("Equal(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestTraversals(t *testing.T) {
	//        1
	//       / \
	//      2   3
	//     / \   \
	//    4   5   6
	root := FromLevelOrder(1, 2, 3, 4, 5, Null, 6)
	wantPre := []int{1, 2, 4, 5, 3, 6}
	wantIn := []int{4, 2, 5, 1, 3, 6}
	wantPost := []int{4, 5, 2, 6, 3, 1}
	check := func(name string, got, want []int) {
		t.Helper()
		if !slices.Equal(got, want) {
			t.Errorf("%s = %v, want %v", name, got, want)
		}
	}
	check("PreorderRecursive", PreorderRecursive(root), wantPre)
	check("PreorderIterative", PreorderIterative(root), wantPre)
	check("InorderRecursive", InorderRecursive(root), wantIn)
	check("InorderIterative", InorderIterative(root), wantIn)
	check("PostorderRecursive", PostorderRecursive(root), wantPost)
	check("PostorderIterative", PostorderIterative(root), wantPost)

	levels := LevelOrder(root)
	wantLevels := [][]int{{1}, {2, 3}, {4, 5, 6}}
	if !slices.EqualFunc(levels, wantLevels, slices.Equal[[]int]) {
		t.Errorf("LevelOrder = %v, want %v", levels, wantLevels)
	}
	if LevelOrder(nil) != nil || PreorderIterative(nil) != nil || InorderIterative(nil) != nil {
		t.Error("traversals of nil should be empty")
	}
}

func TestTraversalsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(10))
	for i := 0; i < 300; i++ {
		root := randomSmallTree(r, r.Intn(25), 0, 10)
		if !slices.Equal(PreorderIterative(root), PreorderRecursive(root)) ||
			!slices.Equal(InorderIterative(root), InorderRecursive(root)) ||
			!slices.Equal(PostorderIterative(root), PostorderRecursive(root)) {
			t.Fatalf("iterative != recursive for %v", ToLevelOrder(root))
		}
		var flat []int
		for _, lvl := range LevelOrder(root) {
			flat = append(flat, lvl...)
		}
		var want []int
		for _, v := range ToLevelOrder(root) {
			if v != Null {
				want = append(want, v)
			}
		}
		if !slices.Equal(flat, want) {
			t.Fatalf("LevelOrder flattened = %v, want %v", flat, want)
		}
	}
}
