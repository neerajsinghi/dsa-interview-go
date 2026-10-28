package stringalgos

// Problem: Valid Palindrome II
// Difficulty: Easy
//
// Return true if s reads the same forwards and backwards after deleting AT
// MOST one character.
//
// Examples:
//
//	"aba"  -> true
//	"abca" -> true   (delete 'b' or 'c')
//	"abc"  -> false
//
// Reference: https://leetcode.com/problems/valid-palindrome-ii/
//
// Pattern: two pointers from both ends. At the first mismatch s[l] != s[r]
// the deleted char must be one of those two, so check whether s[l+1..r] or
// s[l..r-1] is a palindrome. Only one branch point, so it stays linear.
//
// Time: O(n). Space: O(1). Brute force (try every deletion): O(n²).

// ValidPalindromeIIBrute tries deleting nothing, then every index.
func ValidPalindromeIIBrute(s string) bool {
	if isPalindromeRange(s, 0, len(s)-1) {
		return true
	}
	for i := range s {
		if isPalindromeRange(s[:i]+s[i+1:], 0, len(s)-2) {
			return true
		}
	}
	return false
}

// ValidPalindromeII uses two pointers with one allowed skip.
func ValidPalindromeII(s string) bool {
	l, r := 0, len(s)-1
	for l < r {
		if s[l] != s[r] {
			return isPalindromeRange(s, l+1, r) || isPalindromeRange(s, l, r-1)
		}
		l, r = l+1, r-1
	}
	return true
}

// isPalindromeRange reports whether s[l..r] (inclusive) is a palindrome.
func isPalindromeRange(s string, l, r int) bool {
	for l < r {
		if s[l] != s[r] {
			return false
		}
		l, r = l+1, r-1
	}
	return true
}
