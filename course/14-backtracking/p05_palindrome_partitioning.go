package backtracking

// Problem: Palindrome Partitioning
// Difficulty: Medium
//
// Split a string into pieces so that every piece reads the same forwards and
// backwards. Return every such split.
//
// Examples:
//
//	"aab" -> [a a b] [aa b]
//	"a"   -> [a]
//
// Reference: https://leetcode.com/problems/palindrome-partitioning/
//
// Pattern: backtracking over cut positions. From index start, try every end
// such that s[start..end] is a palindrome, take that piece, and recurse from
// end+1. Reaching the end of the string means the path is a full partition.
// Precomputing isPal[i][j] with a small DP makes each check O(1):
// s[i..j] is a palindrome iff s[i] == s[j] and s[i+1..j-1] is one.
//
//	"aab":  start 0 -> "a"  -> start 1 -> "a" -> start 2 -> "b" ✓
//	                                   -> "ab" ✗
//	                -> "aa" -> start 2 -> "b" ✓
//	                -> "aab" ✗
//
// Time: O(n · 2^n) — up to 2^(n-1) partitions, O(n) to copy each.
// Space: O(n²) for isPal.

// Partition returns every palindromic partition of s.
func Partition(s string) [][]string {
	n := len(s)
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
	}
	for i := n - 1; i >= 0; i-- { // i descends so isPal[i+1] is ready
		for j := i; j < n; j++ {
			isPal[i][j] = s[i] == s[j] && (j-i < 2 || isPal[i+1][j-1])
		}
	}

	var res [][]string
	var path []string
	var backtrack func(start int)
	backtrack = func(start int) {
		if start == n {
			res = append(res, append([]string(nil), path...))
			return
		}
		for end := start; end < n; end++ {
			if !isPal[start][end] {
				continue
			}
			path = append(path, s[start:end+1])
			backtrack(end + 1)
			path = path[:len(path)-1]
		}
	}
	backtrack(0)
	return res
}

// PartitionBrute enumerates all 2^(n-1) ways to cut the string (bit i set =
// cut after character i) and keeps those whose pieces are all palindromes.
func PartitionBrute(s string) [][]string {
	n := len(s)
	if n == 0 {
		return [][]string{{}}
	}
	var res [][]string
	for mask := 0; mask < 1<<(n-1); mask++ {
		var parts []string
		start, ok := 0, true
		for i := 0; i < n && ok; i++ {
			if i == n-1 || mask&(1<<i) != 0 {
				piece := s[start : i+1]
				ok = isPalindrome(piece)
				parts = append(parts, piece)
				start = i + 1
			}
		}
		if ok {
			res = append(res, parts)
		}
	}
	return res
}

func isPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}
