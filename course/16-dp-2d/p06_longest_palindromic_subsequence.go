package dp2d

// Problem: Longest Palindromic Subsequence
// Difficulty: Medium
//
// Return the length of the longest subsequence of s that reads the same
// forwards and backwards. (Subsequence: delete any characters, keep order.)
//
// Examples:
//
//	"bbbab" -> 4  ("bbbb")
//	"cbbd"  -> 2  ("bb")
//
// Reference: https://leetcode.com/problems/longest-palindromic-subsequence/
//
// Pattern: interval DP. dp[i][j] = answer for s[i..j] (inclusive).
//
//	s[i] == s[j] -> dp[i+1][j-1] + 2       (wrap the inner answer)
//	otherwise    -> max(dp[i+1][j], dp[i][j-1])
//
// dp[i][j] needs the row below (i+1), so fill i from the bottom up and j
// left to right. Equivalent trick: LPS(s) == LCS(s, reverse(s)).
//
// Time: O(n²). Space: O(n²) table, O(n) with one row.

// LPSMemo is top-down over the interval [i, j].
func LPSMemo(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	memo := newTable(n, n, -1)
	var lps func(i, j int) int
	lps = func(i, j int) int {
		if i > j {
			return 0
		}
		if i == j {
			return 1
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		if s[i] == s[j] {
			memo[i][j] = 2 + lps(i+1, j-1)
		} else {
			memo[i][j] = max(lps(i+1, j), lps(i, j-1))
		}
		return memo[i][j]
	}
	return lps(0, n-1)
}

// LPSTab fills the interval table from short intervals to long ones.
func LPSTab(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	dp := newTable(n, n, 0)
	for i := n - 1; i >= 0; i-- {
		dp[i][i] = 1
		for j := i + 1; j < n; j++ {
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1] + 2 // when j == i+1 this reads dp[i+1][i] == 0
			} else {
				dp[i][j] = max(dp[i+1][j], dp[i][j-1])
			}
		}
	}
	return dp[0][n-1]
}

// LPSRow keeps one row: row[j] holds dp[i+1][j] until overwritten with dp[i][j].
func LPSRow(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	row := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		row[i] = 1
		diag := 0 // dp[i+1][i] (empty interval)
		for j := i + 1; j < n; j++ {
			below := row[j] // dp[i+1][j]
			if s[i] == s[j] {
				row[j] = diag + 2
			} else {
				row[j] = max(below, row[j-1])
			}
			diag = below
		}
	}
	return row[n-1]
}
