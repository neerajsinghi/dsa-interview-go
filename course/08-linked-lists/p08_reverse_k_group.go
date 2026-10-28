package linkedlists

// Problem: Reverse Nodes in k-Group
// Difficulty: Hard
//
// Reverse the list's nodes k at a time and return the new head. If the number
// of nodes left at the end is less than k, leave them in their original order.
// Relink the nodes (don't just swap values) and use O(1) extra memory. k >= 1.
//
// Examples:
//
//	1->2->3->4->5, k=2 -> 2->1->4->3->5
//	1->2->3->4->5, k=3 -> 3->2->1->4->5
//
// Reference: https://leetcode.com/problems/reverse-nodes-in-k-group/
//
// Pattern: reverse a sublist between two fixed anchors, repeated. Keep
// groupPrev (the node before the group). Check that k nodes exist, reverse
// them with the usual prev/cur loop but start prev at the node *after* the
// group so the reversed group links straight onto the rest, then reconnect
// groupPrev and advance it to the group's old first node (now its last).
//
//	dummy -> [1 -> 2 -> 3] -> 4 ...      k=3
//	dummy -> [3 -> 2 -> 1] -> 4 ...      groupPrev moves to 1
//
// Time: O(n) both. Space: brute O(n), optimal O(1).

// ReverseKGroupBrute copies values out, reverses chunks, and rebuilds.
func ReverseKGroupBrute(head *ListNode, k int) *ListNode {
	vals := ToSlice(head)
	for i := 0; i+k <= len(vals); i += k {
		for l, r := i, i+k-1; l < r; l, r = l+1, r-1 {
			vals[l], vals[r] = vals[r], vals[l]
		}
	}
	return FromSlice(vals)
}

// ReverseKGroup relinks nodes in place.
func ReverseKGroup(head *ListNode, k int) *ListNode {
	dummy := &ListNode{Next: head}
	groupPrev := dummy
	for {
		// Find the k-th node of the group; stop if the group is short.
		kth := groupPrev
		for i := 0; i < k && kth != nil; i++ {
			kth = kth.Next
		}
		if kth == nil {
			return dummy.Next
		}
		groupNext := kth.Next

		// Reverse the group. prev starts at groupNext so the reversed
		// group's tail already points at the rest of the list.
		prev, cur := groupNext, groupPrev.Next
		for cur != groupNext {
			next := cur.Next
			cur.Next = prev
			prev, cur = cur, next
		}

		first := groupPrev.Next // old first node, now the group's tail
		groupPrev.Next = kth    // kth is the group's new head
		groupPrev = first
	}
}
