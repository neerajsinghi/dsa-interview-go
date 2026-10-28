package trees

import "math"

// Problem: Binary Tree Maximum Path Sum
// Difficulty: Hard
//
// A path is any sequence of nodes where each adjacent pair is joined by an
// edge, and no node repeats. It has at least one node and need not pass
// through the root. Values may be negative. Return the largest possible sum
// of a path's values. The tree is non-empty.
//
// Examples:
//
//	[1,2,3]                      -> 6    (2-1-3)
//	[-10,9,20,null,null,15,7]    -> 42   (15-20-7)
//	[-3]                         -> -3
//
// Reference: https://leetcode.com/problems/binary-tree-maximum-path-sum/
//
// Pattern: return one thing, record another, same as Diameter (p03). The
// helper returns the best *downward* path starting at the node (the only shape
// a parent can extend), and records the best path that *turns* at the node:
// node + left gain + right gain. A negative gain is dropped: max(gain, 0).
//
// Time: brute O(n^2); optimal O(n). Space: O(n) brute, O(h) optimal.

// MaxPathSumBrute treats the tree as an undirected graph and runs a DFS from
// every node, checking the sum of every simple path.
func MaxPathSumBrute(root *TreeNode) int {
	adj := map[*TreeNode][]*TreeNode{}
	nodes := PreorderNodes(root)
	for _, n := range nodes {
		for _, c := range []*TreeNode{n.Left, n.Right} {
			if c != nil {
				adj[n] = append(adj[n], c)
				adj[c] = append(adj[c], n)
			}
		}
	}
	best := math.MinInt
	var dfs func(n, from *TreeNode, sum int)
	dfs = func(n, from *TreeNode, sum int) {
		sum += n.Val
		best = max(best, sum)
		for _, m := range adj[n] {
			if m != from { // in a tree, not going back is enough to avoid repeats
				dfs(m, n, sum)
			}
		}
	}
	for _, n := range nodes {
		dfs(n, nil, 0)
	}
	return best
}

// MaxPathSum solves it in one post-order pass.
func MaxPathSum(root *TreeNode) int {
	best := math.MinInt
	var gain func(*TreeNode) int
	gain = func(n *TreeNode) int {
		if n == nil {
			return 0
		}
		l := max(gain(n.Left), 0) // a negative branch is better left out
		r := max(gain(n.Right), 0)
		best = max(best, n.Val+l+r) // path that turns at n
		return n.Val + max(l, r)    // parent can extend only one side
	}
	gain(root)
	return best
}
