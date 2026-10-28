package dp2d

// Problem: Unique Paths (with obstacles)
// Difficulty: Medium
//
// A robot stands in the top-left cell of a grid and wants to reach the
// bottom-right cell. Each move goes one cell right or one cell down. Some
// cells are blocked (value 1) and cannot be entered. Count the distinct
// routes. A grid with no obstacles gives the classic "unique paths" count.
//
// Examples:
//
//	3×3 grid, no obstacles            -> 6
//	[[0,0,0],[0,1,0],[0,0,0]]         -> 2  (go around the centre either way)
//	[[0,1],[1,0]]                     -> 0
//
// Reference: https://leetcode.com/problems/unique-paths-ii/
//
// Pattern: grid DP. The only ways into (i, j) are from above and from the
// left, so ways[i][j] = ways[i-1][j] + ways[i][j-1], and 0 on an obstacle.
// Each row only needs the row above, so one rolling row suffices.
//
// Time: O(m·n) for every DP variant. Space: O(m·n) memo/table, O(n) rolling.
// Brute force: O(2^(m+n)) time.

// UniquePathsBrute walks every route recursively. Exponential; tests only.
func UniquePathsBrute(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	var walk func(i, j int) int
	walk = func(i, j int) int {
		if i >= m || j >= n || grid[i][j] == 1 {
			return 0
		}
		if i == m-1 && j == n-1 {
			return 1
		}
		return walk(i+1, j) + walk(i, j+1)
	}
	return walk(0, 0)
}

// UniquePathsMemo is the top-down version: the brute force plus a cache.
func UniquePathsMemo(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	memo := newTable(m, n, -1) // -1 = not computed yet
	var ways func(i, j int) int
	ways = func(i, j int) int {
		if i >= m || j >= n || grid[i][j] == 1 {
			return 0
		}
		if i == m-1 && j == n-1 {
			return 1
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		memo[i][j] = ways(i+1, j) + ways(i, j+1)
		return memo[i][j]
	}
	return ways(0, 0)
}

// UniquePathsTab is the bottom-up table: dp[i][j] = routes from (0,0) to (i,j).
func UniquePathsTab(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	dp := newTable(m, n, 0)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			switch {
			case grid[i][j] == 1:
				dp[i][j] = 0
			case i == 0 && j == 0:
				dp[i][j] = 1
			default:
				if i > 0 {
					dp[i][j] += dp[i-1][j]
				}
				if j > 0 {
					dp[i][j] += dp[i][j-1]
				}
			}
		}
	}
	return dp[m-1][n-1]
}

// UniquePathsRolling keeps one row: before the update row[j] still holds the
// value from the row above, and row[j-1] already holds this row's value.
func UniquePathsRolling(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	n := len(grid[0])
	row := make([]int, n)
	row[0] = 1
	for _, cells := range grid {
		for j, cell := range cells {
			if cell == 1 {
				row[j] = 0
			} else if j > 0 {
				row[j] += row[j-1] // above (old row[j]) + left (new row[j-1])
			}
		}
	}
	return row[n-1]
}

// UniquePaths counts routes on an empty m×n grid.
func UniquePaths(m, n int) int {
	if m <= 0 || n <= 0 {
		return 0
	}
	row := make([]int, n)
	for j := range row {
		row[j] = 1
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			row[j] += row[j-1]
		}
	}
	return row[n-1]
}
