package heaps

// Problem: Merge k Sorted Lists
// Difficulty: Hard
//
// You are given k linked lists, each sorted ascending. Merge them into one
// sorted linked list (reusing the nodes) and return its head.
//
// Example:
//
//	[1->4->5, 1->3->4, 2->6] -> 1->1->2->3->4->4->5->6
//
// Reference: https://leetcode.com/problems/merge-k-sorted-lists/
//
// Pattern: k-way merge with a min-heap of "current heads". The smallest
// remaining value overall is always one of the k heads, so pop it, append it,
// and push its successor. The heap never holds more than k nodes. The
// divide-and-conquer alternative merges lists in pairs, log k rounds of O(N) each.
//
// Time: N = total nodes. Naive "merge one at a time" O(N·k); heap O(N log k);
// pairwise O(N log k). Space: heap O(k); pairwise O(1) iterative.

// ListNode is a singly linked list node (same shape as module 08's).
type ListNode struct {
	Val  int
	Next *ListNode
}

// MergeKListsHeap merges with a min-heap of list heads.
func MergeKListsHeap(lists []*ListNode) *ListNode {
	h := New(func(a, b *ListNode) bool { return a.Val < b.Val })
	for _, l := range lists {
		if l != nil { // never push nil heads
			h.Push(l)
		}
	}
	dummy := &ListNode{}
	tail := dummy
	for h.Len() > 0 {
		n := h.Pop()
		tail.Next = n
		tail = n
		if n.Next != nil {
			h.Push(n.Next)
		}
	}
	tail.Next = nil
	return dummy.Next
}

// MergeKListsPairwise merges neighbouring lists in rounds: k -> k/2 -> ... -> 1.
func MergeKListsPairwise(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	for step := 1; step < len(lists); step *= 2 {
		for i := 0; i+step < len(lists); i += 2 * step {
			lists[i] = mergeTwo(lists[i], lists[i+step])
		}
	}
	return lists[0]
}

func mergeTwo(a, b *ListNode) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	for a != nil && b != nil {
		if a.Val <= b.Val {
			tail.Next, a = a, a.Next
		} else {
			tail.Next, b = b, b.Next
		}
		tail = tail.Next
	}
	if a != nil {
		tail.Next = a
	} else {
		tail.Next = b
	}
	return dummy.Next
}
