package dp2d

// Problem: Minimum Path Sum
// Difficulty: Medium
//
// Each cell of a grid holds a non-negative cost. Moving only right or down,
// travel from the top-left to the bottom-right cell and return the smallest
// possible total of the costs you step on (both endpoints included).
//
// Examples:
//
//	[[1,3,1],[1,5,1],[4,2,1]] -> 7   (1→3→1→1→1)
//	[[1,2,3],[4,5,6]]         -> 12  (1→2→3→6)
//
// Reference: https://leetcode.com/problems/minimum-path-sum/
//
// Pattern: grid DP with min instead of +.
// best[i][j] = grid[i][j] + min(best[i-1][j], best[i][j-1]).
//
// Time: O(m·n). Space: O(m·n) memo/table, O(n) rolling row.

// MinPathSumMemo is top-down: cost(i, j) = cheapest route from (i, j) to the end.
func MinPathSumMemo(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	memo := newTable(m, n, -1)
	var cost func(i, j int) int
	cost = func(i, j int) int {
		if i == m-1 && j == n-1 {
			return grid[i][j]
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		var best int
		switch {
		case i == m-1: // bottom row: can only go right
			best = cost(i, j+1)
		case j == n-1: // right column: can only go down
			best = cost(i+1, j)
		default:
			best = min(cost(i+1, j), cost(i, j+1))
		}
		memo[i][j] = grid[i][j] + best
		return memo[i][j]
	}
	return cost(0, 0)
}

// MinPathSumTab is bottom-up: dp[i][j] = cheapest route from (0,0) to (i,j).
func MinPathSumTab(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	dp := newTable(m, n, 0)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			switch {
			case i == 0 && j == 0:
				dp[i][j] = grid[i][j]
			case i == 0:
				dp[i][j] = dp[i][j-1] + grid[i][j]
			case j == 0:
				dp[i][j] = dp[i-1][j] + grid[i][j]
			default:
				dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j]
			}
		}
	}
	return dp[m-1][n-1]
}

// MinPathSumRolling keeps a single row of the table.
func MinPathSumRolling(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}
	row := make([]int, len(grid[0]))
	for i, cells := range grid {
		for j, c := range cells {
			switch {
			case i == 0 && j == 0:
				row[j] = c
			case i == 0:
				row[j] = row[j-1] + c
			case j == 0:
				row[j] += c // only "above" exists
			default:
				row[j] = min(row[j], row[j-1]) + c
			}
		}
	}
	return row[len(row)-1]
}
