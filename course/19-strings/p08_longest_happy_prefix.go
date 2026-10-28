package stringalgos

// Problem: Longest Happy Prefix
// Difficulty: Hard
//
// A "happy prefix" is a non-empty prefix of s that is also a suffix of s,
// excluding s itself. Return the longest one, or "" if none exists.
//
// Examples:
//
//	"level"  -> "l"
//	"ababab" -> "abab"   (prefix and suffix may overlap)
//	"abc"    -> ""
//
// Reference: https://leetcode.com/problems/longest-happy-prefix/
//
// Pattern: this is literally the last entry of the KMP prefix function —
// pi[n-1] is the longest proper border of the whole string. Alternatively,
// grow the prefix hash and suffix hash together in one pass:
//
//	prefixHash = prefixHash·base + s[i]
//	suffixHash = s[n-1-i]·base^i + suffixHash
//
// and remember the longest length where they agree (then verify).
//
// Time: O(n) for both. Space: O(n) for pi, O(1) for the hash pass.
// Brute force: O(n²).

// LongestPrefixBrute compares prefix and suffix for every length.
func LongestPrefixBrute(s string) string {
	for k := len(s) - 1; k > 0; k-- {
		if s[:k] == s[len(s)-k:] {
			return s[:k]
		}
	}
	return ""
}

// LongestPrefixKMP reads the answer off the prefix function.
func LongestPrefixKMP(s string) string {
	if s == "" {
		return ""
	}
	pi := PrefixFunction(s)
	return s[:pi[len(s)-1]]
}

// LongestPrefixHash compares rolling hashes of s[:k] and s[n-k:] for every k
// in one pass, verifying candidates so collisions can't give wrong answers.
func LongestPrefixHash(s string) string {
	n := len(s)
	var pre, suf uint64
	pow := uint64(1) // base^i
	best := 0
	for k := 1; k < n; k++ {
		pre = (pre*hashBase + uint64(s[k-1])) % hashMod
		suf = (uint64(s[n-k])*pow + suf) % hashMod // prepend s[n-k]
		pow = pow * hashBase % hashMod
		if pre == suf && s[:k] == s[n-k:] {
			best = k
		}
	}
	return s[:best]
}
