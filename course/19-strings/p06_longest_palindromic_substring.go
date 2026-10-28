package stringalgos

// Problem: Longest Palindromic Substring
// Difficulty: Medium
//
// Return the longest contiguous substring of s that is a palindrome. If
// several share the maximum length, return the leftmost.
//
// Examples:
//
//	"babad" -> "bab"  ("aba" is equally long but starts later)
//	"cbbd"  -> "bb"
//
// Reference: https://leetcode.com/problems/longest-palindromic-substring/
//
// Pattern 1 — expand around center: every palindrome mirrors around a
// center, and there are 2n-1 centers (each char, and each gap between
// chars for even lengths). Grow outwards while the ends match. O(n²), O(1)
// space, and the answer interviewers expect first.
//
// Pattern 2 — Manacher: insert a separator between chars so every palindrome
// has odd length ("aba" -> "#a#b#a#"), then reuse radii: inside the
// rightmost palindrome found so far, position i mirrors position
// 2·center−i, so its radius starts at min(rad[mirror], right−i) instead of
// 0. The right edge only moves forward, giving O(n).
//
//	t   = #  a  #  b  #  a  #
//	rad = 0  1  0  3  0  1  0      rad in t == palindrome length in s
//
// Time: O(n³) brute, O(n²) expand, O(n) Manacher. Space: O(1) / O(1) / O(n).

// LongestPalindromeBrute checks every substring, longest first.
func LongestPalindromeBrute(s string) string {
	for length := len(s); length > 0; length-- {
		for start := 0; start+length <= len(s); start++ {
			if isPalindromeRange(s, start, start+length-1) {
				return s[start : start+length]
			}
		}
	}
	return ""
}

// LongestPalindromeExpand tries all 2n-1 centers.
func LongestPalindromeExpand(s string) string {
	bestL, bestLen := 0, 0
	expand := func(l, r int) {
		for l >= 0 && r < len(s) && s[l] == s[r] {
			l, r = l-1, r+1
		}
		// The loop overshot by one on each side: palindrome is s[l+1:r].
		if r-l-1 > bestLen {
			bestL, bestLen = l+1, r-l-1
		}
	}
	for i := range s {
		expand(i, i)   // odd length, centered on s[i]
		expand(i, i+1) // even length, centered between s[i] and s[i+1]
	}
	return s[bestL : bestL+bestLen]
}

// LongestPalindromeManacher runs in linear time. The transformed string t
// is never built: t[i] is a separator for even i and s[i/2] for odd i.
func LongestPalindromeManacher(s string) string {
	n := len(s)
	if n == 0 {
		return ""
	}
	m := 2*n + 1
	// Positions symmetric about a center share parity, so two separators
	// always match and a separator never meets a real char.
	match := func(i, j int) bool { return i%2 == 0 || s[i/2] == s[j/2] }

	rad := make([]int, m) // t[i-rad[i] .. i+rad[i]] is a palindrome
	center, right := 0, 0 // right = center + rad[center], the rightmost edge seen
	best := 0
	for i := 0; i < m; i++ {
		r := 0
		if i < right {
			r = min(rad[2*center-i], right-i) // free radius from the mirror
		}
		for i-r-1 >= 0 && i+r+1 < m && match(i-r-1, i+r+1) {
			r++
		}
		rad[i] = r
		if i+r > right {
			center, right = i, i+r
		}
		if r > rad[best] {
			best = i
		}
	}
	start := (best - rad[best]) / 2
	return s[start : start+rad[best]]
}
