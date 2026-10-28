package linkedlists

// Problem: Merge Two Sorted Lists
// Difficulty: Easy
//
// Two lists are each sorted ascending. Splice their nodes together into one
// sorted list (reuse the nodes; don't allocate new ones) and return its head.
//
// Examples:
//
//	1->2->4  +  1->3->4   ->  1->1->2->3->4->4
//	nil      +  0         ->  0
//
// Reference: https://leetcode.com/problems/merge-two-sorted-lists/
//
// Pattern: dummy head + tail pointer. A dummy node removes the "which list
// provides the head?" special case; at the end, attach whatever is left.
//
// Time: O(m + n). Space: O(1).

// MergeTwoLists merges two sorted lists.
func MergeTwoLists(a, b *ListNode) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	for a != nil && b != nil {
		if a.Val <= b.Val { // <= keeps the merge stable
			tail.Next, a = a, a.Next
		} else {
			tail.Next, b = b, b.Next
		}
		tail = tail.Next
	}
	if a != nil {
		tail.Next = a // the remainder is already sorted
	} else {
		tail.Next = b
	}
	return dummy.Next
}
