package shortestpaths

import (
	"container/heap"
	"sort"
)

// Problem: Swim in Rising Water
// Difficulty: Hard
//
// An n×n grid holds distinct elevations 0..n²-1. At time t the water level
// is t everywhere, and you can swim between adjacent cells (up/down/left/
// right) only if both are at elevation ≤ t; swimming itself is instant.
// Starting at the top-left, return the earliest time you can reach the
// bottom-right.
//
// Examples:
//
//	[[0,2],[1,3]] -> 3   (you must stand on 3 itself)
//
//	 0  1  2  3  4
//	24 23 22 21  5
//	12 13 14 15 16   -> 16
//	11 17 18 19 20
//	10  9  8  7  6
//
// Reference: https://leetcode.com/problems/swim-in-rising-water/
//
// Pattern: minimise the maximum on a path ("bottleneck path"). Two views:
//   - Dijkstra where a path's cost is the highest elevation on it:
//     relax with max(cost[u], grid[v]).
//   - Union-find / Kruskal view: raise the water one level at a time,
//     switching on the cell at that elevation and joining it with already-on
//     neighbours; the answer is the first time start and end are connected.
//
// Time: O(n² log n) for both. Space: O(n²).

// SwimInWater uses minimax Dijkstra.
func SwimInWater(grid [][]int) int {
	n := len(grid)
	best := make([]int, n*n)
	for i := range best {
		best[i] = Inf
	}
	best[0] = grid[0][0]
	h := &minHeap{{0, grid[0][0]}}
	for h.Len() > 0 {
		cur := heap.Pop(h).(item)
		if cur.dist > best[cur.node] {
			continue
		}
		r, c := cur.node/n, cur.node%n
		if r == n-1 && c == n-1 {
			return cur.dist
		}
		for _, d := range dirs4 {
			nr, nc := r+d[0], c+d[1]
			if nr < 0 || nr >= n || nc < 0 || nc >= n {
				continue
			}
			t := max(cur.dist, grid[nr][nc])
			if id := nr*n + nc; t < best[id] {
				best[id] = t
				heap.Push(h, item{id, t})
			}
		}
	}
	return best[n*n-1]
}

// SwimInWaterUnionFind floods cells in order of elevation and stops as soon
// as the two corners join.
func SwimInWaterUnionFind(grid [][]int) int {
	n := len(grid)
	cells := make([]int, n*n)
	for i := range cells {
		cells[i] = i
	}
	// Sorting works for any heights; with a permutation of 0..n²-1 you could
	// instead index cells by elevation directly.
	sort.Slice(cells, func(a, b int) bool {
		return grid[cells[a]/n][cells[a]%n] < grid[cells[b]/n][cells[b]%n]
	})
	on := make([]bool, n*n)
	d := newDSU(n * n)
	for _, id := range cells {
		r, c := id/n, id%n
		on[id] = true
		for _, dir := range dirs4 {
			nr, nc := r+dir[0], c+dir[1]
			if nr >= 0 && nr < n && nc >= 0 && nc < n && on[nr*n+nc] {
				d.union(id, nr*n+nc)
			}
		}
		if on[0] && on[n*n-1] && d.find(0) == d.find(n*n-1) {
			return grid[r][c] // the water level that just switched this cell on
		}
	}
	return -1 // unreachable for a valid grid
}
