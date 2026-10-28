package trees

import (
	"fmt"
	"strconv"
	"strings"
)

// Problem: Serialize and Deserialize Binary Tree
// Difficulty: Hard
//
// Design a pair of functions: one turns any binary tree (values may repeat
// and may be negative) into a string, the other turns that string back into
// an identical tree. The format is up to you.
//
// Example:
//
//	    1
//	   / \          Serialize -> "1,2,#,#,3,4,#,#,5,#,#"
//	  2   3         Deserialize(that) -> the same tree
//	     / \
//	    4   5
//
// Reference: https://leetcode.com/problems/serialize-and-deserialize-binary-tree/
//
// Pattern: preorder with explicit null markers. Recording every nil child
// ("#") makes preorder alone unambiguous, unlike plain traversals where you
// need two (see p07). Deserialize reads tokens in the same order with a
// shared cursor: a token is either "#" (return nil) or a value (make a node,
// then build left, then right).
//
// Time: O(n) both directions. Space: O(n).

const nullToken = "#"

// Serialize encodes the tree as comma-separated preorder tokens.
func Serialize(root *TreeNode) string {
	var sb strings.Builder
	var walk func(*TreeNode)
	walk = func(n *TreeNode) {
		if sb.Len() > 0 {
			sb.WriteByte(',')
		}
		if n == nil {
			sb.WriteString(nullToken)
			return
		}
		sb.WriteString(strconv.Itoa(n.Val))
		walk(n.Left)
		walk(n.Right)
	}
	walk(root)
	return sb.String()
}

// Deserialize decodes a string produced by Serialize. It returns an error for
// malformed input instead of panicking.
func Deserialize(data string) (*TreeNode, error) {
	tokens := strings.Split(data, ",")
	i := 0
	var build func() (*TreeNode, error)
	build = func() (*TreeNode, error) {
		if i >= len(tokens) {
			return nil, fmt.Errorf("unexpected end of input after %d tokens", i)
		}
		tok := tokens[i]
		i++
		if tok == nullToken {
			return nil, nil
		}
		v, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("token %d: %w", i-1, err)
		}
		n := &TreeNode{Val: v}
		if n.Left, err = build(); err != nil {
			return nil, err
		}
		if n.Right, err = build(); err != nil {
			return nil, err
		}
		return n, nil
	}
	root, err := build()
	if err != nil {
		return nil, err
	}
	if i != len(tokens) {
		return nil, fmt.Errorf("%d trailing tokens", len(tokens)-i)
	}
	return root, nil
}
