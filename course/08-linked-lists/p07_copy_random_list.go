package linkedlists

// Problem: Copy List with Random Pointer
// Difficulty: Medium
//
// Each node has a Next pointer and a Random pointer that may point to any node
// in the list or be nil. Build a deep copy: brand-new nodes whose Next and
// Random pointers mirror the original's structure but point only at new nodes.
// The original list must be left unchanged.
//
// Example (value, index of Random):
//
//	[(7,nil) (13,0) (11,4) (10,2) (1,0)]  ->  an identical, fully separate list
//
// Reference: https://leetcode.com/problems/copy-list-with-random-pointer/
//
// Pattern: map old node -> new node. The hash-map version is two passes:
// create every clone, then wire Next/Random through the map. The O(1)-space
// trick stores the map *inside the list*: interleave clones as
// A -> A' -> B -> B' so that X.Next is X's clone, set Random via
// X'.Random = X.Random.Next, then split the lists apart.
//
// Time: O(n) both. Space: O(n) map vs O(1) extra (besides the output).

// RandomNode is a list node with an extra arbitrary pointer.
type RandomNode struct {
	Val    int
	Next   *RandomNode
	Random *RandomNode
}

// CopyRandomListMap deep-copies using an old->new map.
func CopyRandomListMap(head *RandomNode) *RandomNode {
	clone := map[*RandomNode]*RandomNode{nil: nil} // nil maps to nil: no special case
	for n := head; n != nil; n = n.Next {
		clone[n] = &RandomNode{Val: n.Val}
	}
	for n := head; n != nil; n = n.Next {
		clone[n].Next = clone[n.Next]
		clone[n].Random = clone[n.Random]
	}
	return clone[head]
}

// CopyRandomListInterleave deep-copies with O(1) extra space.
func CopyRandomListInterleave(head *RandomNode) *RandomNode {
	if head == nil {
		return nil
	}
	// 1. Weave: A -> A' -> B -> B' -> ...
	for n := head; n != nil; n = n.Next.Next {
		n.Next = &RandomNode{Val: n.Val, Next: n.Next}
	}
	// 2. Random: the clone of X.Random is X.Random.Next.
	for n := head; n != nil; n = n.Next.Next {
		if n.Random != nil {
			n.Next.Random = n.Random.Next
		}
	}
	// 3. Unweave, restoring the original Next pointers.
	newHead := head.Next
	for n := head; n != nil; n = n.Next {
		c := n.Next
		n.Next = c.Next
		if c.Next != nil {
			c.Next = c.Next.Next
		}
	}
	return newHead
}
