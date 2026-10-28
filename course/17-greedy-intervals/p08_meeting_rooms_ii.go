package greedy

import (
	"container/heap"
	"slices"
)

// Problem: Meeting Rooms II (minimum rooms)
// Difficulty: Medium
//
// Meetings are half-open time ranges [Start, End): a meeting ending at 10
// frees its room for one starting at 10. Return the minimum number of rooms
// needed so that no two overlapping meetings share a room.
//
// Examples:
//
//	[0,30) [5,10) [15,20) -> 2
//	[7,10) [2,4)          -> 1
//
// Reference: https://leetcode.com/problems/meeting-rooms-ii/
//
// Pattern: the answer is the maximum number of meetings live at any moment.
// Two ways to measure it:
//
//   - Heap: sort by start; a min-heap holds end times of rooms in use. Before
//     placing a meeting, free the room that ends earliest if it's done.
//   - Sweep: sort starts and ends separately; walk time forward, +1 on a
//     start, -1 on an end. On ties process the END first (half-open ranges).
//
//	time:   0    5   10   15   20   30
//	live:   1    2    1    2    1    0   -> peak 2
//
// Time: O(n log n). Space: O(n).

// MinMeetingRoomsBrute counts live meetings at every start time. O(n²).
func MinMeetingRoomsBrute(meetings []Interval) int {
	best := 0
	for _, m := range meetings {
		live := 0
		for _, o := range meetings {
			if o.Start <= m.Start && m.Start < o.End {
				live++
			}
		}
		best = max(best, live)
	}
	return best
}

// endHeap is a min-heap of meeting end times.
type endHeap []int

func (h endHeap) Len() int           { return len(h) }
func (h endHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h endHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *endHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *endHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// MinMeetingRoomsHeap allocates rooms greedily by start time.
func MinMeetingRoomsHeap(meetings []Interval) int {
	h := &endHeap{}
	for _, m := range sortedByStart(meetings) {
		if h.Len() > 0 && (*h)[0] <= m.Start {
			heap.Pop(h) // reuse the room that frees up earliest
		}
		heap.Push(h, m.End)
	}
	return h.Len() // rooms ever opened = heap size (never shrinks below peak)
}

// MinMeetingRoomsSweep uses two sorted arrays and two pointers.
func MinMeetingRoomsSweep(meetings []Interval) int {
	starts := make([]int, len(meetings))
	ends := make([]int, len(meetings))
	for i, m := range meetings {
		starts[i], ends[i] = m.Start, m.End
	}
	slices.Sort(starts)
	slices.Sort(ends)
	live, best, e := 0, 0, 0
	for _, s := range starts {
		for ends[e] <= s { // free rooms whose meeting ended by time s
			live--
			e++
		}
		live++
		best = max(best, live)
	}
	return best
}
