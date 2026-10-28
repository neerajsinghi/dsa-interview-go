// Package slidingwindow covers problems about contiguous subarrays and
// substrings, solved by maintaining a window [left, right] and updating its
// state incrementally as the edges move — instead of recomputing every window
// from scratch.
package slidingwindow

// Deque is a double-ended queue of ints backed by a slice. Sliding-window
// maximum/minimum problems use it to hold indices in monotonic order.
//
// PopFront re-slices from the front, which is O(1); the dropped prefix is
// reclaimed when append next reallocates. That is fine for interview-sized
// inputs; a ring buffer would avoid the slack in long-running code.
type Deque struct {
	items []int
}

// Len returns the number of elements.
func (d *Deque) Len() int { return len(d.items) }

// PushBack appends x at the back.
func (d *Deque) PushBack(x int) { d.items = append(d.items, x) }

// Front returns the first element. It panics if the deque is empty.
func (d *Deque) Front() int { return d.items[0] }

// Back returns the last element. It panics if the deque is empty.
func (d *Deque) Back() int { return d.items[len(d.items)-1] }

// PopFront removes and returns the first element.
func (d *Deque) PopFront() int {
	x := d.items[0]
	d.items = d.items[1:]
	return x
}

// PopBack removes and returns the last element.
func (d *Deque) PopBack() int {
	x := d.items[len(d.items)-1]
	d.items = d.items[:len(d.items)-1]
	return x
}
