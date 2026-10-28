package arrayshashing

import "slices"

// Problem 02 — Valid Anagram (Easy)
//
// Two strings are anagrams if one is a rearrangement of the other's letters.
// Given s and t made of lowercase English letters, report whether they are
// anagrams.
//
//	s = "listen", t = "silent" → true
//	s = "rat",    t = "car"    → false
//
// Reference: https://leetcode.com/problems/valid-anagram/
//
// Pattern: anagrams have identical letter frequencies. With a 26-letter
// alphabet a [26]int array beats a map, and Go arrays compare with ==.
// Follow-up: for arbitrary Unicode, count runes in a map[rune]int.
//
//	IsAnagramSort:    Time O(n log n), Space O(n)
//	IsAnagram:        Time O(n),       Space O(1) (26 counters)
//	IsAnagramUnicode: Time O(n),       Space O(k) distinct runes

// IsAnagramSort sorts both strings' bytes and compares.
func IsAnagramSort(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	a, b := []byte(s), []byte(t)
	slices.Sort(a)
	slices.Sort(b)
	return string(a) == string(b)
}

// IsAnagram compares letter counts.
func IsAnagram(s, t string) bool {
	if len(s) != len(t) {
		return false // cheap early exit
	}
	return letterCounts(s) == letterCounts(t)
}

// IsAnagramUnicode handles any UTF-8 input by counting runes, not bytes:
// "é" is two bytes but one letter.
func IsAnagramUnicode(s, t string) bool {
	counts := make(map[rune]int)
	for _, r := range s {
		counts[r]++
	}
	for _, r := range t {
		counts[r]--
		if counts[r] < 0 {
			return false // t has more of r than s
		}
	}
	for _, c := range counts {
		if c != 0 {
			return false // s has letters t lacks
		}
	}
	return true
}
