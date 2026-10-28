package trees

// Problem: Maximum Depth of Binary Tree
// Difficulty: Easy
//
// Return the number of nodes on the longest path from the root down to a leaf.
// The empty tree has depth 0.
//
// Examples:
//
//	[3,9,20,null,null,15,7] -> 3
//	[1,null,2]              -> 2
//
// Reference: https://leetcode.com/problems/maximum-depth-of-binary-tree/
//
// Pattern: the basic tree recursion. Answer for a node = combine(answer for
// left, answer for right): here 1 + max(left, right). The BFS version counts
// levels instead.
//
// Time: O(n). Space: O(h) recursion (h = height, up to n), O(width) for BFS.

// MaxDepth computes depth recursively (DFS).
func MaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + max(MaxDepth(root.Left), MaxDepth(root.Right))
}

// MaxDepthBFS counts levels with a queue. No recursion, so it is safe for
// very deep (degenerate) trees.
func MaxDepthBFS(root *TreeNode) int {
	if root == nil {
		return 0
	}
	depth := 0
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		depth++
		size := len(queue)
		for _, n := range queue[:size] {
			if n.Left != nil {
				queue = append(queue, n.Left)
			}
			if n.Right != nil {
				queue = append(queue, n.Right)
			}
		}
		queue = queue[size:]
	}
	return depth
}
