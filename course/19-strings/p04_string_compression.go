package stringalgos

import "strconv"

// Problem: String Compression
// Difficulty: Medium
//
// Compress a byte slice in place using run-length encoding: each run of a
// repeated character becomes the character followed by the run length in
// decimal — unless the run has length 1, in which case only the character
// is written. Return the new length; the compressed form occupies
// chars[:newLen]. Use O(1) extra space.
//
// Examples:
//
//	"aabbccc"      -> 6, "a2b2c3"
//	"a"            -> 1, "a"
//	"abbbbbbbbbbbb" -> 4, "ab12"
//
// Reference: https://leetcode.com/problems/string-compression/
//
// Pattern: read/write two pointers. The write pointer never overtakes the
// read pointer, because a run of length k ≥ 2 needs 1 + digits(k) ≤ k bytes.
// Write the count's digits directly; strconv.AppendInt on a scratch array
// keeps it allocation-free.
//
// Time: O(n). Space: O(1).

// CompressNaive builds a new string — easy to get right, used as a reference.
func CompressNaive(chars []byte) string {
	var out []byte
	for i := 0; i < len(chars); {
		j := i
		for j < len(chars) && chars[j] == chars[i] {
			j++
		}
		out = append(out, chars[i])
		if j-i > 1 {
			out = strconv.AppendInt(out, int64(j-i), 10)
		}
		i = j
	}
	return string(out)
}

// Compress rewrites chars in place and returns the new length.
func Compress(chars []byte) int {
	write := 0
	for read := 0; read < len(chars); {
		ch, start := chars[read], read
		for read < len(chars) && chars[read] == ch {
			read++
		}
		chars[write] = ch
		write++
		if run := read - start; run > 1 {
			var digits [20]byte // enough for any int64
			for _, d := range strconv.AppendInt(digits[:0], int64(run), 10) {
				chars[write] = d
				write++
			}
		}
	}
	return write
}
