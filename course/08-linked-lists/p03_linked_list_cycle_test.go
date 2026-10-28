package linkedlists

import "testing"

// buildWithCycle builds a list from vals and links the tail to the node at
// index pos (pos < 0 means no cycle). It returns the head and that node.
func buildWithCycle(vals []int, pos int) (head, start *ListNode) {
	head = FromSlice(vals)
	if pos < 0 {
		return head, nil
	}
	var tail *ListNode
	for i, n := 0, head; n != nil; i, n = i+1, n.Next {
		if i == pos {
			start = n
		}
		tail = n
	}
	tail.Next = start
	return head, start
}

func TestCycle(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		pos  int
	}{
		{"classic", []int{3, 2, 0, -4}, 1},
		{"two nodes", []int{1, 2}, 0},
		{"self loop", []int{1}, 0},
		{"single no cycle", []int{1}, -1},
		{"empty", nil, -1},
		{"long no cycle", []int{1, 2, 3, 4, 5, 6, 7}, -1},
		{"cycle at tail", []int{1, 2, 3, 4, 5, 6, 7}, 6},
		{"cycle at head", []int{1, 2, 3, 4, 5, 6, 7}, 0},
		{"duplicate values", []int{1, 1, 1, 1, 1}, 3},
	}
	for _, tt := range tests {
		head, start := buildWithCycle(tt.vals, tt.pos)
		want := tt.pos >= 0
		if got := HasCycle(head); got != want {
			t.Errorf("%s: HasCycle = %v, want %v", tt.name, got, want)
		}
		if got := HasCycleSet(head); got != want {
			t.Errorf("%s: HasCycleSet = %v, want %v", tt.name, got, want)
		}
		if got := DetectCycle(head); got != start {
			t.Errorf("%s: DetectCycle returned wrong node", tt.name)
		}
	}
}

func TestCycleAllPositions(t *testing.T) {
	vals := make([]int, 20)
	for pos := -1; pos < len(vals); pos++ {
		head, start := buildWithCycle(vals, pos)
		if got := DetectCycle(head); got != start {
			t.Fatalf("pos %d: DetectCycle returned wrong node", pos)
		}
	}
}
