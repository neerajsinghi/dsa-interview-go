package stringalgos

// Problem: Find the Index of the First Occurrence in a String (strStr)
// Difficulty: Easy (Medium if you must do it in linear time)
//
// Return the smallest index at which needle occurs inside haystack, or -1 if
// it never does. An empty needle matches at index 0.
//
// Examples:
//
//	haystack "sadbutsad", needle "sad" -> 0
//	haystack "leetcode",  needle "leeto" -> -1
//
// Reference: https://leetcode.com/problems/find-the-index-of-the-first-occurrence-in-a-string/
//
// Pattern: linear-time pattern matching. Three classic ways:
//
//   - KMP: build the prefix function of needle; when a char mismatches, fall
//     back to the longest border instead of restarting — the haystack
//     pointer never moves backwards.
//   - Z-function on needle + "\x00" + haystack: a match is any position whose
//     z value equals len(needle). (The separator must not occur in either.)
//   - Rabin–Karp: compare rolling hashes of each window; verify on a hit.
//
// Time: O(n + m) for KMP and Z, expected O(n + m) for Rabin–Karp.
// Space: O(m) (KMP), O(n + m) (Z, Rabin–Karp). Brute force: O(n·m).
// Production Go: strings.Index (which uses Rabin–Karp/SIMD internally).

// StrStrBrute compares the needle at every start.
func StrStrBrute(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

// StrStrKMP scans haystack once using needle's failure table.
func StrStrKMP(haystack, needle string) int {
	m := len(needle)
	if m == 0 {
		return 0
	}
	pi := PrefixFunction(needle)
	k := 0 // number of needle chars currently matched
	for i := 0; i < len(haystack); i++ {
		for k > 0 && haystack[i] != needle[k] {
			k = pi[k-1]
		}
		if haystack[i] == needle[k] {
			k++
		}
		if k == m {
			return i - m + 1
		}
	}
	return -1
}

// StrStrZ runs the Z-function on needle + sep + haystack.
func StrStrZ(haystack, needle string) int {
	m := len(needle)
	if m == 0 {
		return 0
	}
	z := ZFunction(needle + "\x00" + haystack)
	for i := m + 1; i < len(z); i++ {
		if z[i] >= m { // z can never exceed m thanks to the separator
			return i - m - 1
		}
	}
	return -1
}

// StrStrRabinKarp slides a rolling hash and verifies each hash hit.
func StrStrRabinKarp(haystack, needle string) int {
	n, m := len(haystack), len(needle)
	if m == 0 {
		return 0
	}
	if m > n {
		return -1
	}
	target := NewRollingHash(needle).Hash(0, m)
	h := NewRollingHash(haystack)
	for i := 0; i+m <= n; i++ {
		if h.Hash(i, i+m) == target && haystack[i:i+m] == needle {
			return i
		}
	}
	return -1
}
