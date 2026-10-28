package trees

// Problem: Validate Binary Search Tree
// Difficulty: Medium
//
// Decide whether a binary tree is a valid BST: for every node, all values in
// its left subtree are strictly smaller, and all values in its right subtree
// are strictly larger. Duplicates make it invalid.
//
// Examples:
//
//	[2,1,3]               -> true
//	[5,1,4,null,null,3,6] -> false (4 is in 5's right subtree but 4 < 5)
//	[5,4,6,null,null,3,7] -> false (3 is in 5's right subtree; checking only
//	                                 parent/child pairs misses it)
//
// Reference: https://leetcode.com/problems/validate-binary-search-tree/
//
// Pattern: pass constraints down. Each node must lie in an open interval
// (lo, hi) inherited from its ancestors: going left tightens hi, going right
// tightens lo. Equivalent view: the inorder sequence must be strictly
// increasing.
//
// Time: brute O(n^2); both optimal versions O(n). Space: O(h).

// IsValidBSTBrute checks every node against every node in its subtrees.
func IsValidBSTBrute(root *TreeNode) bool {
	if root == nil {
		return true
	}
	for _, v := range PreorderRecursive(root.Left) {
		if v >= root.Val {
			return false
		}
	}
	for _, v := range PreorderRecursive(root.Right) {
		if v <= root.Val {
			return false
		}
	}
	return IsValidBSTBrute(root.Left) && IsValidBSTBrute(root.Right)
}

// IsValidBST passes (lo, hi) bounds down. Bounds are pointers so that "no
// bound" is nil, not a sentinel like math.MinInt that a real node could hold.
func IsValidBST(root *TreeNode) bool {
	var valid func(n *TreeNode, lo, hi *int) bool
	valid = func(n *TreeNode, lo, hi *int) bool {
		if n == nil {
			return true
		}
		if (lo != nil && n.Val <= *lo) || (hi != nil && n.Val >= *hi) {
			return false
		}
		return valid(n.Left, lo, &n.Val) && valid(n.Right, &n.Val, hi)
	}
	return valid(root, nil, nil)
}

// IsValidBSTInorder walks the tree in order and checks that each value is
// larger than the previous one, stopping early on the first violation.
func IsValidBSTInorder(root *TreeNode) bool {
	var prev *TreeNode
	var stack []*TreeNode
	cur := root
	for cur != nil || len(stack) > 0 {
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.Left
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if prev != nil && cur.Val <= prev.Val {
			return false
		}
		prev = cur
		cur = cur.Right
	}
	return true
}
