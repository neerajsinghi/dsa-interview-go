package design

import "sync"

// Problem: Design Hit Counter
// Difficulty: Medium
//
// Count hits received in the past 5 minutes. Hit(ts) records one hit at
// second ts; GetHits(ts) returns how many hits happened in (ts-300, ts].
// Timestamps arrive in non-decreasing order. Several hits may share a
// timestamp, and the counter must handle a huge number of hits per second.
//
// Example:
//
//	Hit(1) Hit(2) Hit(3) GetHits(4) -> 3
//	Hit(300) GetHits(300) -> 4
//	GetHits(301) -> 3   (the hit at 1 has expired)
//
// Reference: https://leetcode.com/problems/design-hit-counter/
//
// Pattern: fixed ring of time buckets. Bucket i = ts % 300 stores the last
// timestamp written there and its count. A stale bucket (different
// timestamp) is reset on write and ignored on read. Memory is O(window),
// independent of traffic — unlike a queue of every hit.
//
//	ts:      1   2   3  ...  300  301
//	bucket:  1   2   3  ...   0    1  <- 301 reuses bucket 1, evicting ts=1
//
// Time: Hit O(1), GetHits O(window). Space: O(window).
//
// Safe for concurrent use (a sync.Mutex guards the buckets).

// HitWindow is the window length in seconds.
const HitWindow = 300

// HitCounter counts hits in a sliding 300-second window.
type HitCounter struct {
	mu     sync.Mutex
	times  [HitWindow]int
	counts [HitWindow]int
}

// Hit records one hit at second ts.
func (h *HitCounter) Hit(ts int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	i := ts % HitWindow
	if h.times[i] != ts { // bucket holds an expired second: recycle it
		h.times[i] = ts
		h.counts[i] = 0
	}
	h.counts[i]++
}

// GetHits returns the number of hits in (ts-300, ts].
func (h *HitCounter) GetHits(ts int) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	total := 0
	for i := range h.times {
		if ts-h.times[i] < HitWindow {
			total += h.counts[i]
		}
	}
	return total
}
