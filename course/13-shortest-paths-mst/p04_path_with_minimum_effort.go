package shortestpaths

import (
	"container/heap"
	"sort"
)

// Problem: Path With Minimum Effort
// Difficulty: Medium
//
// A hiker crosses a grid of heights from the top-left to the bottom-right
// cell, moving up/down/left/right. The effort of a route is the LARGEST
// absolute height difference between two consecutive cells on it. Return the
// minimum possible effort.
//
// Examples:
//
//	1 2 2
//	3 8 2    -> 2   (1→3→5→3→5 has max step 2)
//	5 3 5
//
//	1 2 3
//	3 8 4    -> 1   (go right along the top, then down the right side)
//	5 3 5
//
// Reference: https://leetcode.com/problems/path-with-minimum-effort/
//
// Pattern: Dijkstra with a different path cost. The cost of a path is a max,
// not a sum, but it still never decreases as the path grows — that
// monotonicity is all Dijkstra needs. Relax with
// newCost = max(cost[u], |h[u]-h[v]|) instead of cost[u] + w.
//
// Alternative: binary search the answer. "Is there a route with effort ≤ k?"
// is a plain BFS that only uses steps ≤ k, and the answer is monotone in k.
//
// Time: Dijkstra O(m·n·log(m·n)); binary search O(m·n·log(maxHeight)).
// Space: O(m·n).

var dirs4 = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

// MinimumEffortPath runs "minimax" Dijkstra over grid cells.
func MinimumEffortPath(heights [][]int) int {
	m, n := len(heights), len(heights[0])
	effort := make([]int, m*n) // cell (r, c) has id r*n + c
	for i := range effort {
		effort[i] = Inf
	}
	effort[0] = 0
	h := &minHeap{{0, 0}}
	for h.Len() > 0 {
		cur := heap.Pop(h).(item)
		if cur.dist > effort[cur.node] {
			continue
		}
		r, c := cur.node/n, cur.node%n
		if r == m-1 && c == n-1 {
			return cur.dist // first time the target is popped, it is final
		}
		for _, d := range dirs4 {
			nr, nc := r+d[0], c+d[1]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue
			}
			next := max(cur.dist, abs(heights[nr][nc]-heights[r][c])) // max, not +
			if id := nr*n + nc; next < effort[id] {
				effort[id] = next
				heap.Push(h, item{id, next})
			}
		}
	}
	return effort[m*n-1]
}

// MinimumEffortPathBinarySearch finds the smallest limit for which a BFS
// using only steps ≤ limit reaches the target.
func MinimumEffortPathBinarySearch(heights [][]int) int {
	m, n := len(heights), len(heights[0])
	reachable := func(limit int) bool {
		seen := make([]bool, m*n)
		seen[0] = true
		queue := []int{0}
		for len(queue) > 0 {
			id := queue[0]
			queue = queue[1:]
			r, c := id/n, id%n
			if id == m*n-1 {
				return true
			}
			for _, d := range dirs4 {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < m && nc >= 0 && nc < n && !seen[nr*n+nc] && abs(heights[nr][nc]-heights[r][c]) <= limit {
					seen[nr*n+nc] = true
					queue = append(queue, nr*n+nc)
				}
			}
		}
		return false
	}
	maxH := 0
	for _, row := range heights {
		for _, v := range row {
			maxH = max(maxH, v)
		}
	}
	// sort.Search returns the smallest limit in [0, maxH] that is reachable.
	return sort.Search(maxH+1, reachable)
}
