package trees

// Problem: Construct Binary Tree from Preorder and Inorder Traversal
// Difficulty: Medium
//
// You are given the preorder and inorder traversals of a binary tree whose
// values are all distinct. Rebuild the tree and return its root.
//
// Example:
//
//	preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]
//
//	    3
//	   / \
//	  9   20
//	     /  \
//	    15   7
//
// Reference: https://leetcode.com/problems/construct-binary-tree-from-preorder-and-inorder-traversal/
//
// Pattern: divide and conquer on traversal properties. preorder[0] is the
// root. Finding the root in inorder splits it into left subtree | root |
// right subtree, and the left part's size tells you how much of preorder
// belongs to the left subtree. A value -> inorder index map makes each lookup
// O(1). Consume preorder left-to-right with a shared cursor.
//
// Time: brute O(n^2) (linear search + slicing); optimal O(n). Space: O(n).

// BuildTreeBrute recurses on sub-slices and scans for the root each time.
func BuildTreeBrute(preorder, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}
	root := &TreeNode{Val: preorder[0]}
	mid := 0
	for inorder[mid] != root.Val {
		mid++
	}
	// mid = size of the left subtree
	root.Left = BuildTreeBrute(preorder[1:mid+1], inorder[:mid])
	root.Right = BuildTreeBrute(preorder[mid+1:], inorder[mid+1:])
	return root
}

// BuildTree uses an index map and a preorder cursor.
func BuildTree(preorder, inorder []int) *TreeNode {
	pos := make(map[int]int, len(inorder))
	for i, v := range inorder {
		pos[v] = i
	}
	next := 0 // index of the next root in preorder

	// build returns the subtree whose values are inorder[lo:hi].
	var build func(lo, hi int) *TreeNode
	build = func(lo, hi int) *TreeNode {
		if lo >= hi {
			return nil
		}
		root := &TreeNode{Val: preorder[next]}
		next++
		mid := pos[root.Val]
		root.Left = build(lo, mid) // must come first: preorder is root, LEFT, right
		root.Right = build(mid+1, hi)
		return root
	}
	return build(0, len(inorder))
}
