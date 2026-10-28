package dp2d

// Problem: Interleaving String
// Difficulty: Medium
//
// Decide whether s3 can be formed by shuffling s1 and s2 together while
// keeping the internal order of each — like riffling two decks of cards.
//
// Examples:
//
//	s1="aabcc", s2="dbbca", s3="aadbbcbcac" -> true
//	s1="aabcc", s2="dbbca", s3="aadbbbaccc" -> false
//
// Reference: https://leetcode.com/problems/interleaving-string/
//
// Pattern: two-sequence DP. After using i chars of s1 and j chars of s2 we
// are at s3[i+j]. dp[i][j] = "s1[:i] and s2[:j] can interleave into
// s3[:i+j]". The last char came from s1 or from s2:
//
//	dp[i][j] = (dp[i-1][j] && s1[i-1] == s3[i+j-1]) ||
//	           (dp[i][j-1] && s2[j-1] == s3[i+j-1])
//
// Greedy fails: when both strings offer the next char you must try both.
//
// Time: O(m·n). Space: O(m·n) table, O(n) rolling row.

// IsInterleaveMemo is top-down; memo values: 0 unknown, 1 true, 2 false.
func IsInterleaveMemo(s1, s2, s3 string) bool {
	if len(s1)+len(s2) != len(s3) {
		return false
	}
	memo := newTable(len(s1)+1, len(s2)+1, 0)
	var can func(i, j int) bool
	can = func(i, j int) bool {
		if i == len(s1) && j == len(s2) {
			return true
		}
		if memo[i][j] != 0 {
			return memo[i][j] == 1
		}
		k := i + j
		ok := (i < len(s1) && s1[i] == s3[k] && can(i+1, j)) ||
			(j < len(s2) && s2[j] == s3[k] && can(i, j+1))
		memo[i][j] = 2
		if ok {
			memo[i][j] = 1
		}
		return ok
	}
	return can(0, 0)
}

// IsInterleaveTab is the bottom-up boolean table.
func IsInterleaveTab(s1, s2, s3 string) bool {
	m, n := len(s1), len(s2)
	if m+n != len(s3) {
		return false
	}
	dp := newBoolTable(m+1, n+1)
	dp[0][0] = true
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i > 0 && dp[i-1][j] && s1[i-1] == s3[i+j-1] {
				dp[i][j] = true
			}
			if j > 0 && dp[i][j-1] && s2[j-1] == s3[i+j-1] {
				dp[i][j] = true
			}
		}
	}
	return dp[m][n]
}

// IsInterleaveRow keeps a single row indexed by j.
func IsInterleaveRow(s1, s2, s3 string) bool {
	m, n := len(s1), len(s2)
	if m+n != len(s3) {
		return false
	}
	row := make([]bool, n+1)
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			switch {
			case i == 0 && j == 0:
				row[j] = true
			case i == 0:
				row[j] = row[j-1] && s2[j-1] == s3[j-1]
			case j == 0:
				row[j] = row[j] && s1[i-1] == s3[i-1]
			default:
				row[j] = (row[j] && s1[i-1] == s3[i+j-1]) || // from above
					(row[j-1] && s2[j-1] == s3[i+j-1]) // from the left
			}
		}
	}
	return row[n]
}
