package stringalgos

import "slices"

// Problem: Repeated DNA Sequences
// Difficulty: Medium
//
// A DNA string uses only 'A', 'C', 'G', 'T'. Return every 10-letter
// substring that occurs more than once (overlaps allowed), each listed once.
// This version returns them sorted so the output is deterministic.
//
// Examples:
//
//	"AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT" -> ["AAAAACCCCC" "CCCCCAAAAA"]
//	"AAAAAAAAAAAAA"                    -> ["AAAAAAAAAA"]
//
// Reference: https://leetcode.com/problems/repeated-dna-sequences/
//
// Pattern: rolling hash with an exact encoding. Four letters fit in 2 bits,
// so a 10-letter window is a 20-bit integer — a PERFECT hash, no collisions.
// Slide the window by shifting in 2 new bits and masking off the oldest:
//
//	window = (window<<2 | code(c)) & (1<<20 - 1)
//
// Counting uint32 keys instead of 10-byte strings avoids hashing and
// allocating a substring per position.
//
// Time: O(n). Space: O(n) (at most 4^10 distinct windows).

const dnaLen = 10

// FindRepeatedDNASubstrings uses a map of substrings — simple and correct.
func FindRepeatedDNASubstrings(s string) []string {
	count := map[string]int{}
	var out []string
	for i := 0; i+dnaLen <= len(s); i++ {
		sub := s[i : i+dnaLen] // substring shares memory with s; no copy
		count[sub]++
		if count[sub] == 2 {
			out = append(out, sub)
		}
	}
	slices.Sort(out)
	return out
}

// FindRepeatedDNA uses the 2-bit rolling encoding.
func FindRepeatedDNA(s string) []string {
	code := [256]uint32{'A': 0, 'C': 1, 'G': 2, 'T': 3}
	const mask = 1<<(2*dnaLen) - 1
	count := map[uint32]int{}
	var out []string
	var window uint32
	for i := 0; i < len(s); i++ {
		window = (window<<2 | code[s[i]]) & mask
		if i < dnaLen-1 {
			continue // window not full yet
		}
		count[window]++
		if count[window] == 2 {
			out = append(out, s[i-dnaLen+1:i+1])
		}
	}
	slices.Sort(out)
	return out
}
