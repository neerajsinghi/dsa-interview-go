package trees

// Problem: Diameter of Binary Tree
// Difficulty: Easy
//
// The diameter is the number of edges on the longest path between any two
// nodes. The path need not pass through the root.
//
// Examples:
//
//	[1,2,3,4,5] -> 3     (4-2-1-3 or 5-2-1-3)
//	[1,2]       -> 1
//
// Reference: https://leetcode.com/problems/diameter-of-binary-tree/
//
// Pattern: return one thing, record another. The recursive helper returns
// the node's height (what the parent needs), and on the way it updates a
// global best with leftHeight + rightHeight (the longest path that turns at this
// node). The same shape solves Max Path Sum (p09) and Balanced Tree.
//
// Time: brute O(n^2) (recomputes heights); optimal O(n). Space: O(h).

// height returns the number of nodes on the longest downward path.
func height(n *TreeNode) int {
	if n == nil {
		return 0
	}
	return 1 + max(height(n.Left), height(n.Right))
}

// DiameterBrute tries every node as the turning point, recomputing heights.
func DiameterBrute(root *TreeNode) int {
	if root == nil {
		return 0
	}
	through := height(root.Left) + height(root.Right)
	return max(through, DiameterBrute(root.Left), DiameterBrute(root.Right))
}

// Diameter computes heights once, bottom-up, tracking the best path.
func Diameter(root *TreeNode) int {
	best := 0
	var h func(*TreeNode) int
	h = func(n *TreeNode) int {
		if n == nil {
			return 0
		}
		l, r := h(n.Left), h(n.Right)
		best = max(best, l+r) // longest path that turns at n (in edges)
		return 1 + max(l, r)  // what the parent needs: height through n
	}
	h(root)
	return best
}
