package linkedlists

// Problem: Linked List Cycle (detect + find the start)
// Difficulty: Medium
//
// A list may loop back on itself: some node's Next points to an earlier node.
// Part 1: report whether a cycle exists. Part 2: return the first node of the
// cycle (where the tail reconnects), or nil if there is none. Use O(1) memory.
//
// Example:
//
//	3 -> 2 -> 0 -> -4
//	     ^          |
//	     +----------+      has cycle; it starts at the node with value 2
//
// References: https://leetcode.com/problems/linked-list-cycle/ and
// https://leetcode.com/problems/linked-list-cycle-ii/
//
// Pattern: Floyd's tortoise and hare. slow moves 1, fast moves 2. In a cycle,
// fast gains one step per move and must land on slow. To find the start:
// if the tail before the cycle has length a and they meet b steps into a cycle
// of length c, then 2(a+b) = a+b+kc, so a = kc - b. A pointer from head and a
// pointer from the meeting point, both moving 1 step, meet at the cycle start.
//
// Time: O(n). Space: O(1) (the hash-set version is O(n) space).

// HasCycleSet remembers every visited node. Simple, O(n) extra space.
func HasCycleSet(head *ListNode) bool {
	seen := map[*ListNode]bool{} // pointer identity, not value equality
	for n := head; n != nil; n = n.Next {
		if seen[n] {
			return true
		}
		seen[n] = true
	}
	return false
}

// HasCycle uses Floyd's fast/slow pointers.
func HasCycle(head *ListNode) bool {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}

// DetectCycle returns the first node of the cycle, or nil.
func DetectCycle(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow, fast = slow.Next, fast.Next.Next
		if slow == fast {
			// Phase 2: distance head->start equals distance meet->start (mod c).
			p := head
			for p != slow {
				p, slow = p.Next, slow.Next
			}
			return p
		}
	}
	return nil
}
