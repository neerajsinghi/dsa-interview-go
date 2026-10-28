package linkedlists

// Problem: Reorder List
// Difficulty: Medium
//
// Rearrange the nodes of L0 -> L1 -> ... -> Ln into
// L0 -> Ln -> L1 -> Ln-1 -> L2 -> ... in place (relink nodes; values don't move).
//
// Examples:
//
//	1->2->3->4    -> 1->4->2->3
//	1->2->3->4->5 -> 1->5->2->4->3
//
// Reference: https://leetcode.com/problems/reorder-list/
//
// Pattern: three classic moves composed. (1) find the middle with slow/fast,
// (2) reverse the second half, (3) merge the two halves alternately.
//
//	1->2->3 | 4->5   reverse 2nd half   1->2->3   5->4   weave   1->5->2->4->3
//
// Time: O(n) both. Space: brute O(n) (node slice), optimal O(1).

// ReorderListBrute copies node pointers into a slice and relinks by index.
func ReorderListBrute(head *ListNode) {
	var nodes []*ListNode
	for n := head; n != nil; n = n.Next {
		nodes = append(nodes, n)
	}
	i, j := 0, len(nodes)-1
	for i < j {
		nodes[i].Next = nodes[j]
		i++
		if i == j {
			break
		}
		nodes[j].Next = nodes[i]
		j--
	}
	if len(nodes) > 0 {
		nodes[i].Next = nil // the middle node becomes the new tail
	}
}

// ReorderList does it in O(1) extra space.
func ReorderList(head *ListNode) {
	if head == nil || head.Next == nil {
		return
	}
	// 1. Middle: slow stops at the end of the first half (first half is the
	//    longer one for odd lengths).
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
	}
	second := ReverseList(slow.Next) // 2. reverse second half (from p01)
	slow.Next = nil                  //    and cut the lists apart

	// 3. Weave: first, second, first, second, ...
	first := head
	for second != nil {
		fNext, sNext := first.Next, second.Next
		first.Next = second
		second.Next = fNext
		first, second = fNext, sNext
	}
}
