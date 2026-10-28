package linkedlists

// Problem: Reverse Linked List
// Difficulty: Easy
//
// Reverse a singly linked list in place and return the new head.
//
// Examples:
//
//	1 -> 2 -> 3 -> 4 -> nil   becomes   4 -> 3 -> 2 -> 1 -> nil
//	nil                       stays     nil
//
// Reference: https://leetcode.com/problems/reverse-linked-list/
//
// Pattern: three-pointer rewiring. Walk with prev and cur; before you point
// cur.Next back at prev, save cur.Next so you don't lose the rest of the list.
//
//	prev   cur    next
//	 nil <- 1      2 -> 3 -> nil
//
// Time: O(n). Space: O(1) iterative, O(n) recursion stack.

// ReverseList reverses the list iteratively.
func ReverseList(head *ListNode) *ListNode {
	var prev *ListNode
	cur := head
	for cur != nil {
		next := cur.Next // save the rest before overwriting the link
		cur.Next = prev
		prev, cur = cur, next
	}
	return prev
}

// ReverseListRecursive reverses the list recursively: reverse everything after
// head, then hang head at the end of that reversed tail.
func ReverseListRecursive(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	newHead := ReverseListRecursive(head.Next)
	head.Next.Next = head // head.Next is now the tail of the reversed part
	head.Next = nil
	return newHead
}
