package trees

// Problem: Binary Tree Right Side View
// Difficulty: Medium
//
// Imagine standing to the right of the tree. List the values you can see from
// top to bottom: the rightmost node on each level.
//
// Examples:
//
//	[1,2,3,null,5,null,4] -> [1,3,4]
//	[1,2,3,4]             -> [1,3,4]   (4 is visible: nothing is to its right)
//
// Reference: https://leetcode.com/problems/binary-tree-right-side-view/
//
// Pattern: level-order traversal, and keep the last node of each level.
// The DFS alternative visits right before left and records the first node
// seen at each new depth.
//
// Time: O(n). Space: O(width) for BFS, O(h) for DFS.

// RightSideView uses BFS level by level.
func RightSideView(root *TreeNode) []int {
	var view []int
	if root == nil {
		return view
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		size := len(queue)
		view = append(view, queue[size-1].Val) // last node on this level
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
	return view
}

// RightSideViewDFS visits right subtrees first; the first node reached at a
// given depth is the visible one.
func RightSideViewDFS(root *TreeNode) []int {
	var view []int
	var dfs func(n *TreeNode, depth int)
	dfs = func(n *TreeNode, depth int) {
		if n == nil {
			return
		}
		if depth == len(view) { // first time we reach this depth
			view = append(view, n.Val)
		}
		dfs(n.Right, depth+1)
		dfs(n.Left, depth+1)
	}
	dfs(root, 0)
	return view
}
