package dp2d

// Problem: Longest Common Subsequence
// Difficulty: Medium
//
// A subsequence keeps some characters of a string in their original order
// (gaps allowed). Given two strings, return the length of the longest
// subsequence they share.
//
// Examples:
//
//	"abcde", "ace" -> 3  ("ace")
//	"abc",   "def" -> 0
//
// Reference: https://leetcode.com/problems/longest-common-subsequence/
//
// Pattern: two-sequence DP over prefixes. dp[i][j] = LCS of a[:i] and b[:j].
//
//	a[i-1] == b[j-1]  ->  dp[i][j] = dp[i-1][j-1] + 1   (use both chars)
//	otherwise         ->  dp[i][j] = max(dp[i-1][j], dp[i][j-1])  (drop one)
//
// Row 0 and column 0 are the empty-prefix base case (0).
//
// Time: O(m·n). Space: O(m·n) table, O(min(m,n)) with two rows.
// Brute force: O(2^(m+n)).

// LCSBrute tries both "drop a char" branches whenever chars differ.
func LCSBrute(a, b string) int {
	if a == "" || b == "" {
		return 0
	}
	if a[0] == b[0] {
		return 1 + LCSBrute(a[1:], b[1:])
	}
	return max(LCSBrute(a[1:], b), LCSBrute(a, b[1:]))
}

// LCSMemo is top-down over suffix start indices (i, j).
func LCSMemo(a, b string) int {
	memo := newTable(len(a), len(b), -1)
	var lcs func(i, j int) int
	lcs = func(i, j int) int {
		if i == len(a) || j == len(b) {
			return 0
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		if a[i] == b[j] {
			memo[i][j] = 1 + lcs(i+1, j+1)
		} else {
			memo[i][j] = max(lcs(i+1, j), lcs(i, j+1))
		}
		return memo[i][j]
	}
	return lcs(0, 0)
}

// LCSTab is the bottom-up prefix table.
func LCSTab(a, b string) int {
	return lcsTable(a, b)[len(a)][len(b)]
}

func lcsTable(a, b string) [][]int {
	dp := newTable(len(a)+1, len(b)+1, 0)
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp
}

// LCSTwoRows keeps only the previous and current rows.
func LCSTwoRows(a, b string) int {
	if len(b) > len(a) {
		a, b = b, a // make the row length the shorter string
	}
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				cur[j] = prev[j-1] + 1
			} else {
				cur[j] = max(prev[j], cur[j-1])
			}
		}
		prev, cur = cur, prev // reuse the old row as the next scratch row
	}
	return prev[len(b)]
}

// LCSString reconstructs one longest common subsequence by walking the
// table backwards from dp[m][n]. Interviewers often ask for this follow-up.
func LCSString(a, b string) string {
	dp := lcsTable(a, b)
	out := make([]byte, 0, dp[len(a)][len(b)])
	i, j := len(a), len(b)
	for i > 0 && j > 0 {
		switch {
		case a[i-1] == b[j-1]:
			out = append(out, a[i-1])
			i, j = i-1, j-1
		case dp[i-1][j] >= dp[i][j-1]:
			i--
		default:
			j--
		}
	}
	// Built back to front; reverse it.
	for l, r := 0, len(out)-1; l < r; l, r = l+1, r-1 {
		out[l], out[r] = out[r], out[l]
	}
	return string(out)
}
