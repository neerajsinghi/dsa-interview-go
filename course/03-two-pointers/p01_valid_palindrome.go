package twopointers

// Problem 01 — Valid Palindrome (Easy)
//
// Ignoring case and every character that is not an ASCII letter or digit,
// report whether s reads the same forwards and backwards.
//
//	"A man, a plan, a canal: Panama" → true   ("amanaplanacanalpanama")
//	"race a car"                     → false
//
// Reference: https://leetcode.com/problems/valid-palindrome/
//
// Pattern: opposite-ends two pointers. Skip characters that don't count,
// compare the two ends, and move both inwards. No cleaned copy is needed.
//
//	IsPalindromeBrute: Time O(n), Space O(n) (cleaned copy)
//	IsPalindrome:      Time O(n), Space O(1)

// IsPalindromeBrute builds the cleaned string, then compares it with its
// reverse.
func IsPalindromeBrute(s string) bool {
	clean := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if isAlphaNum(s[i]) {
			clean = append(clean, toLower(s[i]))
		}
	}
	rev := make([]byte, len(clean))
	copy(rev, clean)
	Reverse(rev)
	return string(clean) == string(rev)
}

// IsPalindrome compares in place.
func IsPalindrome(s string) bool {
	l, r := 0, len(s)-1
	for l < r {
		if !isAlphaNum(s[l]) {
			l++
			continue
		}
		if !isAlphaNum(s[r]) {
			r--
			continue
		}
		if toLower(s[l]) != toLower(s[r]) {
			return false
		}
		l++
		r--
	}
	return true
}
