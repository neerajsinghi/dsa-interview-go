package graphs

// Problem: Rotting Oranges
// Difficulty: Medium
//
// A grid holds 0 (empty), 1 (fresh orange) or 2 (rotten orange). Every
// minute, each fresh orange next to a rotten one (up, down, left, right)
// turns rotten. Return the minutes until no fresh orange remains, or -1 if
// some fresh orange can never be reached.
//
// Examples:
//
//	2 1 1
//	1 1 0   -> 4
//	0 1 1
//
//	2 1 1
//	0 1 1   -> -1   (bottom-left orange is walled off)
//	1 0 1
//
// Reference: https://leetcode.com/problems/rotting-oranges/
//
// Pattern: multi-source BFS. Put EVERY rotten orange in the queue at time 0;
// BFS then expands all fronts simultaneously, one ring per minute, and the
// BFS depth of the last orange reached is the answer. Running a separate BFS
// from each rotten orange would be both slower and wrong (the fronts race).
//
//	minute 0      minute 1      minute 2      ...
//	2 1 1         2 2 1         2 2 2
//	1 1 0   ->    2 1 0   ->    2 2 0   ->  ...
//	0 1 1         0 1 1         0 1 1
//
// Time: O(m·n). Space: O(m·n) for the queue.
// The simulation brute force is O((m·n)²).

// OrangesRotting runs a level-by-level multi-source BFS.
func OrangesRotting(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])
	g := make([][]int, m) // work on a copy so the caller's grid is untouched
	var queue [][2]int
	fresh := 0
	for r := range grid {
		g[r] = append([]int(nil), grid[r]...)
		for c, v := range grid[r] {
			switch v {
			case 2:
				queue = append(queue, [2]int{r, c})
			case 1:
				fresh++
			}
		}
	}
	minutes := 0
	for len(queue) > 0 && fresh > 0 {
		var next [][2]int // everything that rots during this minute
		for _, cell := range queue {
			for _, d := range dirs4 {
				r, c := cell[0]+d[0], cell[1]+d[1]
				if r >= 0 && r < m && c >= 0 && c < n && g[r][c] == 1 {
					g[r][c] = 2 // rot it now so it is enqueued only once
					fresh--
					next = append(next, [2]int{r, c})
				}
			}
		}
		queue = next
		minutes++
	}
	if fresh > 0 {
		return -1
	}
	return minutes
}

// OrangesRottingSimulate replays the process literally: each round, find the
// fresh oranges adjacent to rotten ones, then rot them all at once.
func OrangesRottingSimulate(grid [][]int) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])
	g := make([][]int, m)
	for r := range grid {
		g[r] = append([]int(nil), grid[r]...)
	}
	for minutes := 0; ; minutes++ {
		var toRot [][2]int
		fresh := 0
		for r := 0; r < m; r++ {
			for c := 0; c < n; c++ {
				if g[r][c] != 1 {
					continue
				}
				fresh++
				for _, d := range dirs4 {
					nr, nc := r+d[0], c+d[1]
					if nr >= 0 && nr < m && nc >= 0 && nc < n && g[nr][nc] == 2 {
						toRot = append(toRot, [2]int{r, c})
						break
					}
				}
			}
		}
		if fresh == 0 {
			return minutes
		}
		if len(toRot) == 0 {
			return -1 // fresh oranges left but nothing spreads any more
		}
		for _, cell := range toRot {
			g[cell[0]][cell[1]] = 2
		}
	}
}
