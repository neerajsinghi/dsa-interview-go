package trees

// Problem: Lowest Common Ancestor (BST and general binary tree)
// Difficulty: Medium
//
// Given two nodes p and q that are both in the tree, return their lowest
// common ancestor: the deepest node that has both p and q in its subtree. A
// node counts as being in its own subtree, so if p is an ancestor of q the
// answer is p.
//
// Example (BST [6,2,8,0,4,7,9,null,null,3,5]):
//
//	        6
//	      /   \
//	     2     8
//	    / \   / \
//	   0   4 7   9
//	      / \
//	     3   5
//
//	LCA(2, 8) = 6    LCA(2, 4) = 2    LCA(3, 5) = 4
//
// References: https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-search-tree/
// and https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-tree/
//
// Pattern:
//   - BST: the LCA is the first node whose value lies between p and q. If both
//     are smaller go left, if both are larger go right, otherwise you're at the split.
//   - General tree: post-order search. Each call returns p or q if it found one
//     below, otherwise nil. The first node that gets a non-nil result from *both*
//     sides is the LCA. If only one side finds something, pass it up.
//
// Time: BST O(h); general O(n). Space: BST O(1) iterative; general O(h).

// LowestCommonAncestorBST walks down from the root using BST ordering.
func LowestCommonAncestorBST(root, p, q *TreeNode) *TreeNode {
	for n := root; n != nil; {
		switch {
		case p.Val < n.Val && q.Val < n.Val:
			n = n.Left
		case p.Val > n.Val && q.Val > n.Val:
			n = n.Right
		default:
			return n // p and q split here (or one of them is n)
		}
	}
	return nil
}

// LowestCommonAncestor works on any binary tree. It compares node pointers,
// so duplicate values are fine.
func LowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == p || root == q {
		return root
	}
	left := LowestCommonAncestor(root.Left, p, q)
	right := LowestCommonAncestor(root.Right, p, q)
	switch {
	case left != nil && right != nil:
		return root // p on one side, q on the other
	case left != nil:
		return left
	default:
		return right
	}
}
