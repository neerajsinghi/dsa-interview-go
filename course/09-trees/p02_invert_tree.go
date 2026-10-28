package trees

// Problem: Invert Binary Tree
// Difficulty: Easy
//
// Mirror a binary tree: at every node, swap the left and right subtrees.
// Return the root.
//
// Example:
//
//	    4                 4
//	   / \               / \
//	  2   7     ->      7   2
//	 / \ / \           / \ / \
//	1  3 6  9         9  6 3  1
//
// Reference: https://leetcode.com/problems/invert-binary-tree/
//
// Pattern: visit every node once and do a local operation. Any traversal
// order works (pre, post, or BFS) as long as every node gets swapped.
//
// Time: O(n). Space: O(h) recursion, O(width) for BFS.

// InvertTree mirrors the tree recursively.
func InvertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	root.Left, root.Right = InvertTree(root.Right), InvertTree(root.Left)
	return root
}

// InvertTreeBFS mirrors the tree level by level with a queue.
func InvertTreeBFS(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		n.Left, n.Right = n.Right, n.Left
		if n.Left != nil {
			queue = append(queue, n.Left)
		}
		if n.Right != nil {
			queue = append(queue, n.Right)
		}
	}
	return root
}
