package trees

import (
	"math/rand"
	"testing"
)

// findVal returns the first node (preorder) with value v.
func findVal(root *TreeNode, v int) *TreeNode {
	for _, n := range PreorderNodes(root) {
		if n.Val == v {
			return n
		}
	}
	return nil
}

func TestLowestCommonAncestor(t *testing.T) {
	bst := FromLevelOrder(6, 2, 8, 0, 4, 7, 9, Null, Null, 3, 5)
	tests := []struct{ p, q, want int }{
		{2, 8, 6}, {2, 4, 2}, {3, 5, 4}, {0, 5, 2}, {7, 9, 8}, {3, 9, 6}, {6, 6, 6}, {5, 4, 4},
	}
	for _, tt := range tests {
		p, q := findVal(bst, tt.p), findVal(bst, tt.q)
		if got := LowestCommonAncestorBST(bst, p, q); got == nil || got.Val != tt.want {
			t.Errorf("LCA_BST(%d, %d) = %v, want %d", tt.p, tt.q, got, tt.want)
		}
		if got := LowestCommonAncestor(bst, p, q); got == nil || got.Val != tt.want {
			t.Errorf("LCA(%d, %d) = %v, want %d", tt.p, tt.q, got, tt.want)
		}
	}

	general := FromLevelOrder(3, 5, 1, 6, 2, 0, 8, Null, Null, 7, 4)
	for _, tt := range []struct{ p, q, want int }{{5, 1, 3}, {5, 4, 5}, {7, 8, 3}, {6, 4, 5}} {
		got := LowestCommonAncestor(general, findVal(general, tt.p), findVal(general, tt.q))
		if got == nil || got.Val != tt.want {
			t.Errorf("LCA general(%d, %d) = %v, want %d", tt.p, tt.q, got, tt.want)
		}
	}
}

// lcaBrute collects each node's root-to-node path and takes the last shared node.
func lcaBrute(root, p, q *TreeNode) *TreeNode {
	var path func(n, target *TreeNode, acc []*TreeNode) []*TreeNode
	path = func(n, target *TreeNode, acc []*TreeNode) []*TreeNode {
		if n == nil {
			return nil
		}
		acc = append(acc, n)
		if n == target {
			return acc
		}
		if found := path(n.Left, target, acc); found != nil {
			return found
		}
		return path(n.Right, target, acc)
	}
	pp, qp := path(root, p, nil), path(root, q, nil)
	var lca *TreeNode
	for i := 0; i < len(pp) && i < len(qp) && pp[i] == qp[i]; i++ {
		lca = pp[i]
	}
	return lca
}

func TestLowestCommonAncestorRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for i := 0; i < 300; i++ {
		root := randomSmallTree(r, 1+r.Intn(25), 0, 3) // many duplicate values
		nodes := PreorderNodes(root)
		p, q := nodes[r.Intn(len(nodes))], nodes[r.Intn(len(nodes))]
		if got, want := LowestCommonAncestor(root, p, q), lcaBrute(root, p, q); got != want {
			t.Fatalf("LCA mismatch on %v", ToLevelOrder(root))
		}

		var bst *TreeNode
		for _, v := range r.Perm(20)[:1+r.Intn(19)] {
			bst = insertBST(bst, v)
		}
		bn := PreorderNodes(bst)
		p, q = bn[r.Intn(len(bn))], bn[r.Intn(len(bn))]
		if got, want := LowestCommonAncestorBST(bst, p, q), lcaBrute(bst, p, q); got != want {
			t.Fatalf("LCA_BST mismatch on %v", ToLevelOrder(bst))
		}
	}
}
