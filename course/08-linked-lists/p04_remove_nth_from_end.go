package linkedlists

// Problem: Remove Nth Node From End of List
// Difficulty: Medium
//
// Delete the n-th node counting from the end (n = 1 is the last node) and
// return the head. Assume 1 <= n <= length. Try to do it in a single pass.
//
// Examples:
//
//	1->2->3->4->5, n=2 -> 1->2->3->5
//	1, n=1             -> nil
//
// Reference: https://leetcode.com/problems/remove-nth-node-from-end-of-list/
//
// Pattern: two pointers with a fixed gap, plus a dummy head. Move fast n+1
// steps ahead of slow (both start at dummy). When fast falls off the end, slow
// sits just before the node to delete. The dummy handles deleting the head.
//
// Time: O(L) both versions (two-pass vs one-pass). Space: O(1).

// RemoveNthFromEndTwoPass counts the length first, then walks to L-n.
func RemoveNthFromEndTwoPass(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	prev := dummy
	for i := 0; i < Len(head)-n; i++ {
		prev = prev.Next
	}
	prev.Next = prev.Next.Next
	return dummy.Next
}

// RemoveNthFromEnd does it in one pass with a gap of n+1.
func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	slow, fast := dummy, dummy
	for i := 0; i <= n; i++ { // n+1 steps so slow stops *before* the target
		fast = fast.Next
	}
	for fast != nil {
		slow, fast = slow.Next, fast.Next
	}
	slow.Next = slow.Next.Next
	return dummy.Next
}
