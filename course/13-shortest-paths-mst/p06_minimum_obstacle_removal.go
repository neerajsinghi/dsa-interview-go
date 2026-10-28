package shortestpaths

// Problem: Minimum Obstacle Removal to Reach Corner
// Difficulty: Hard
//
// A grid has 0 for an empty cell and 1 for an obstacle you may remove. Moving
// up/down/left/right from the top-left to the bottom-right cell, what is the
// fewest obstacles you must remove? (Start and end are always empty.)
//
// Examples:
//
//	0 1 1
//	1 1 0   -> 2
//	1 1 0
//
//	0 1 0 0 0
//	0 1 0 1 0   -> 0   (a path around exists)
//	0 0 0 1 0
//
// Reference: https://leetcode.com/problems/minimum-obstacle-removal-to-reach-corner/
//
// Pattern: 0-1 BFS. Stepping into a cell costs its value — 0 or 1. With only
// two possible weights, a deque replaces Dijkstra's heap: a 0-cost step goes
// to the FRONT (same distance as now), a 1-cost step to the BACK (distance
// + 1). The deque then always holds at most two consecutive distances in
// sorted order, which is exactly the invariant Dijkstra needs — without the
// log factor.
//
//	deque: [d d d | d+1 d+1]   pop front; 0-edge -> push front, 1-edge -> push back
//
// Time: O(m·n) for 0-1 BFS vs O(m·n·log(m·n)) for Dijkstra. Space: O(m·n).

// MinimumObstacles uses 0-1 BFS on the grid.
func MinimumObstacles(grid [][]int) int {
	adj := gridGraph(grid)
	return ZeroOneBFS(adj, 0)[len(adj)-1]
}

// MinimumObstaclesDijkstra solves the same graph with Dijkstra.
func MinimumObstaclesDijkstra(grid [][]int) int {
	adj := gridGraph(grid)
	return Dijkstra(adj, 0)[len(adj)-1]
}

// gridGraph turns the grid into a directed weighted graph: the edge into
// cell (r, c) weighs grid[r][c].
func gridGraph(grid [][]int) [][]Arc {
	m, n := len(grid), len(grid[0])
	adj := make([][]Arc, m*n)
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			for _, d := range dirs4 {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < m && nc >= 0 && nc < n {
					adj[r*n+c] = append(adj[r*n+c], Arc{nr*n + nc, grid[nr][nc]})
				}
			}
		}
	}
	return adj
}
