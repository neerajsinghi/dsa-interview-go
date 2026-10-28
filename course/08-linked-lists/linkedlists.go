// Package linkedlists covers singly linked list manipulation: pointer
// rewiring, the dummy-head trick, and fast/slow pointers.
package linkedlists

// ListNode is a node of a singly linked list of ints. A nil *ListNode is the
// empty list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// FromSlice builds a list with the values in order and returns its head
// (nil for an empty slice).
func FromSlice(vals []int) *ListNode {
	dummy := &ListNode{} // dummy head: no special case for the first node
	tail := dummy
	for _, v := range vals {
		tail.Next = &ListNode{Val: v}
		tail = tail.Next
	}
	return dummy.Next
}

// ToSlice returns the list's values in order. It must not be called on a list
// that contains a cycle (it would never terminate).
func ToSlice(head *ListNode) []int {
	var out []int
	for n := head; n != nil; n = n.Next {
		out = append(out, n.Val)
	}
	return out
}

// Len returns the number of nodes in an acyclic list.
func Len(head *ListNode) int {
	n := 0
	for ; head != nil; head = head.Next {
		n++
	}
	return n
}
