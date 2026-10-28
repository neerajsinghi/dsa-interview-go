package dp2d

// Problem: Regular Expression Matching
// Difficulty: Hard
//
// Implement a tiny regex engine supporting two metacharacters:
//
//	'.' matches any single character
//	'*' means "zero or more of the element right before me"
//
// The pattern must match the WHOLE string. Patterns are assumed valid
// (a '*' always follows a real character or '.').
//
// Examples:
//
//	s="aa",  p="a"     -> false
//	s="aa",  p="a*"    -> true
//	s="ab",  p=".*"    -> true
//	s="aab", p="c*a*b" -> true   (c* matches nothing)
//
// Reference: https://leetcode.com/problems/regular-expression-matching/
//
// Pattern: two-sequence DP over suffixes. match(i, j) = s[i:] matches p[j:].
// Look at whether p[j+1] is '*':
//
//	star:    match(i, j+2)                    // use x* zero times
//	      || first && match(i+1, j)           // consume one char, keep x*
//	no star: first && match(i+1, j+1)
//
// where first = i < len(s) && (p[j] == s[i] || p[j] == '.').
//
// Time: O(m·n). Space: O(m·n). Plain recursion is exponential on inputs
// like s="aaaaaaaaaaaaab", p="a*a*a*a*a*c".

// RegexBrute is the direct recursion with no caching.
func RegexBrute(s, p string) bool {
	if p == "" {
		return s == ""
	}
	first := s != "" && (p[0] == s[0] || p[0] == '.')
	if len(p) >= 2 && p[1] == '*' {
		return RegexBrute(s, p[2:]) || (first && RegexBrute(s[1:], p))
	}
	return first && RegexBrute(s[1:], p[1:])
}

// RegexMemo caches match(i, j). Memo values: 0 unknown, 1 true, 2 false.
func RegexMemo(s, p string) bool {
	memo := newTable(len(s)+1, len(p)+1, 0)
	var match func(i, j int) bool
	match = func(i, j int) bool {
		if j == len(p) {
			return i == len(s)
		}
		if memo[i][j] != 0 {
			return memo[i][j] == 1
		}
		first := i < len(s) && (p[j] == s[i] || p[j] == '.')
		var ok bool
		if j+1 < len(p) && p[j+1] == '*' {
			ok = match(i, j+2) || (first && match(i+1, j))
		} else {
			ok = first && match(i+1, j+1)
		}
		memo[i][j] = 2
		if ok {
			memo[i][j] = 1
		}
		return ok
	}
	return match(0, 0)
}

// RegexTab fills dp[i][j] (= s[i:] matches p[j:]) from the end backwards.
func RegexTab(s, p string) bool {
	m, n := len(s), len(p)
	dp := newBoolTable(m+1, n+1)
	dp[m][n] = true // empty matches empty
	for i := m; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			first := i < m && (p[j] == s[i] || p[j] == '.')
			if j+1 < n && p[j+1] == '*' {
				dp[i][j] = dp[i][j+2] || (first && dp[i+1][j])
			} else {
				dp[i][j] = first && dp[i+1][j+1]
			}
		}
	}
	return dp[0][0]
}
