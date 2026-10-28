// Package trees covers binary trees and binary search trees: the four
// traversals, recursion that returns information up the tree, and BST
// ordering.
package trees

import "math"

// TreeNode is a binary tree node. A nil *TreeNode is the empty tree.
type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

// Null marks a missing child in level-order slices passed to FromLevelOrder
// and returned by ToLevelOrder. It plays the role of LeetCode's "null".
const Null = math.MinInt

// FromLevelOrder builds a tree from LeetCode-style level-order values, where
// Null marks an absent child and children of absent nodes are not listed:
//
//	FromLevelOrder(1, 2, 3, Null, 4)
//
//	    1
//	   / \
//	  2   3
//	   \
//	    4
func FromLevelOrder(vals ...int) *TreeNode {
	if len(vals) == 0 || vals[0] == Null {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) && vals[i] != Null {
			node.Left = &TreeNode{Val: vals[i]}
			queue = append(queue, node.Left)
		}
		i++
		if i < len(vals) && vals[i] != Null {
			node.Right = &TreeNode{Val: vals[i]}
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}

// ToLevelOrder is the inverse of FromLevelOrder: level-order values with Null
// for absent children of present nodes, trailing Nulls trimmed.
func ToLevelOrder(root *TreeNode) []int {
	var out []int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			out = append(out, Null)
			continue
		}
		out = append(out, node.Val)
		queue = append(queue, node.Left, node.Right)
	}
	for len(out) > 0 && out[len(out)-1] == Null {
		out = out[:len(out)-1]
	}
	return out
}

// Equal reports whether two trees have the same shape and values
// (LeetCode's "Same Tree"). The pattern is to recurse on both trees together.
func Equal(a, b *TreeNode) bool {
	if a == nil || b == nil {
		return a == b // both nil -> true; exactly one nil -> false
	}
	return a.Val == b.Val && Equal(a.Left, b.Left) && Equal(a.Right, b.Right)
}

// ---------------------------------------------------------------------------
// Traversals. Depth-first orders are named by when the root is visited:
// pre = root, left, right; in = left, root, right; post = left, right, root.
// In a BST, inorder visits values in ascending order.
// ---------------------------------------------------------------------------

// PreorderRecursive returns values in root-left-right order.
func PreorderRecursive(root *TreeNode) []int {
	var out []int
	var walk func(*TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		out = append(out, n.Val)
		walk(n.Left)
		walk(n.Right)
	}
	walk(root)
	return out
}

// InorderRecursive returns values in left-root-right order.
func InorderRecursive(root *TreeNode) []int {
	var out []int
	var walk func(*TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		walk(n.Left)
		out = append(out, n.Val)
		walk(n.Right)
	}
	walk(root)
	return out
}

// PreorderNodes returns the nodes (not values) in preorder.
func PreorderNodes(root *TreeNode) []*TreeNode {
	var out []*TreeNode
	var walk func(*TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		out = append(out, n)
		walk(n.Left)
		walk(n.Right)
	}
	walk(root)
	return out
}

// PostorderRecursive returns values in left-right-root order.
func PostorderRecursive(root *TreeNode) []int {
	var out []int
	var walk func(*TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		walk(n.Left)
		walk(n.Right)
		out = append(out, n.Val)
	}
	walk(root)
	return out
}

// PreorderIterative uses an explicit stack. Push right before left so that
// left is popped (visited) first.
func PreorderIterative(root *TreeNode) []int {
	var out []int
	if root == nil {
		return out
	}
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		out = append(out, n.Val)
		if n.Right != nil {
			stack = append(stack, n.Right)
		}
		if n.Left != nil {
			stack = append(stack, n.Left)
		}
	}
	return out
}

// InorderIterative walks as far left as possible, pushing nodes; when it can't
// go left, it pops (visits) a node and then turns to its right subtree.
func InorderIterative(root *TreeNode) []int {
	var out []int
	var stack []*TreeNode
	cur := root
	for cur != nil || len(stack) > 0 {
		for cur != nil { // slide down the left spine
			stack = append(stack, cur)
			cur = cur.Left
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		out = append(out, cur.Val)
		cur = cur.Right
	}
	return out
}

// PostorderIterative produces root-right-left with a stack (a mirrored
// preorder) and reverses it, which gives left-right-root.
func PostorderIterative(root *TreeNode) []int {
	var out []int
	if root == nil {
		return out
	}
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		out = append(out, n.Val)
		if n.Left != nil { // left pushed first, so right is visited first
			stack = append(stack, n.Left)
		}
		if n.Right != nil {
			stack = append(stack, n.Right)
		}
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// LevelOrder returns values grouped by depth (breadth-first search). The
// trick is to snapshot len(queue) at the start of each level.
func LevelOrder(root *TreeNode) [][]int {
	var levels [][]int
	if root == nil {
		return levels
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		size := len(queue) // nodes on this level
		level := make([]int, 0, size)
		for i := 0; i < size; i++ {
			n := queue[i]
			level = append(level, n.Val)
			if n.Left != nil {
				queue = append(queue, n.Left)
			}
			if n.Right != nil {
				queue = append(queue, n.Right)
			}
		}
		queue = queue[size:]
		levels = append(levels, level)
	}
	return levels
}
