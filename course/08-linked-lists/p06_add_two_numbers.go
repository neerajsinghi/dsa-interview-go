package linkedlists

// Problem: Add Two Numbers
// Difficulty: Medium
//
// Two non-negative integers are stored as linked lists of decimal digits,
// least-significant digit first (so 342 is 2 -> 4 -> 3). Return their sum in
// the same form. Inputs have no leading zeros except the number 0 itself.
//
// Examples:
//
//	2->4->3 + 5->6->4 -> 7->0->8   (342 + 465 = 807)
//	9->9    + 1       -> 0->0->1   (99 + 1 = 100)
//
// Reference: https://leetcode.com/problems/add-two-numbers/
//
// Pattern: grade-school addition with a carry, walking both lists together.
// Keep looping while either list has digits OR a carry remains. Never
// convert to an int: the numbers can be far longer than 64 bits.
//
// Time: O(max(m, n)). Space: O(max(m, n)) for the output.

// AddTwoNumbers returns the digit list of a + b.
func AddTwoNumbers(a, b *ListNode) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	carry := 0
	for a != nil || b != nil || carry > 0 { // the final carry can add a digit
		sum := carry
		if a != nil {
			sum += a.Val
			a = a.Next
		}
		if b != nil {
			sum += b.Val
			b = b.Next
		}
		carry = sum / 10
		tail.Next = &ListNode{Val: sum % 10}
		tail = tail.Next
	}
	return dummy.Next
}
