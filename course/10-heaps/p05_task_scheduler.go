package heaps

// Problem: Task Scheduler
// Difficulty: Medium
//
// A CPU runs tasks labelled 'A'..'Z', one per time unit, in any order. Two runs
// of the same label must be at least n units apart; in between the CPU can run
// other tasks or sit idle. Return the minimum number of time units needed to
// finish every task.
//
// Examples:
//
//	tasks=AAABBB, n=2 -> 8   (A B _ A B _ A B)
//	tasks=AAABBB, n=0 -> 6
//	tasks=AAAAAABCDEFG, n=2 -> 16
//
// Reference: https://leetcode.com/problems/task-scheduler/
//
// Pattern: greedy with a max-heap. At each tick run the available task with
// the most remaining copies (it's the one most likely to cause idling later).
// Tasks on cooldown wait in a FIFO queue with the time they become ready.
//
// Insight for the O(n) formula: the most frequent task (count m) forces m-1
// full frames of length n+1, plus one slot per task tied for that max count.
// If there are enough other tasks to fill every gap, there is no idling and the
// answer is just len(tasks).
//
//	A _ _ | A _ _ | A        m=3, n=2  ->  (m-1)*(n+1) + 1 = 7
//
// Time: heap simulation O(T log 26) where T is the answer; formula O(n).
// Space: O(26).

// LeastIntervalHeap simulates the schedule tick by tick.
func LeastIntervalHeap(tasks []byte, n int) int {
	var counts [26]int
	for _, t := range tasks {
		counts[t-'A']++
	}
	ready := NewMax[int]() // remaining counts of tasks that can run now
	for _, c := range counts {
		if c > 0 {
			ready.Push(c)
		}
	}
	type cooling struct{ remaining, readyAt int }
	var queue []cooling // FIFO: tasks finish cooling in the order they ran

	time := 0
	for ready.Len() > 0 || len(queue) > 0 {
		time++
		if ready.Len() > 0 {
			if left := ready.Pop() - 1; left > 0 {
				queue = append(queue, cooling{left, time + n})
			}
		} // else: the CPU idles this tick
		if len(queue) > 0 && queue[0].readyAt == time {
			ready.Push(queue[0].remaining)
			queue = queue[1:]
		}
	}
	return time
}

// LeastInterval computes the answer directly from the counts.
func LeastInterval(tasks []byte, n int) int {
	var counts [26]int
	maxCount := 0
	for _, t := range tasks {
		counts[t-'A']++
		maxCount = max(maxCount, counts[t-'A'])
	}
	if maxCount == 0 {
		return 0
	}
	numMax := 0 // how many labels share the max count
	for _, c := range counts {
		if c == maxCount {
			numMax++
		}
	}
	return max(len(tasks), (maxCount-1)*(n+1)+numMax)
}
