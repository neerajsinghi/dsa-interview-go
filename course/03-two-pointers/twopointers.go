// Package twopointers covers problems solved by moving two indices through a
// slice or string — towards each other (opposite ends) or in the same
// direction (slow/fast, read/write) — so each element is visited O(1) times.
package twopointers

// Reverse reverses s in place with two pointers meeting in the middle.
// It is the simplest instance of the opposite-ends pattern.
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// isAlphaNum reports whether c is an ASCII letter or digit.
func isAlphaNum(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9'
}

// toLower lower-cases an ASCII letter and leaves other bytes unchanged.
func toLower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
