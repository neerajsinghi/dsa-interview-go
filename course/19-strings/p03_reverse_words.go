package stringalgos

import "strings"

// Problem: Reverse Words in a String
// Difficulty: Medium
//
// A word is a maximal run of non-space characters. Return the words in
// reverse order joined by single spaces, with no leading/trailing spaces —
// however many spaces the input had between words.
//
// Examples:
//
//	"the sky is blue"    -> "blue is sky the"
//	"  hello   world  "  -> "world hello"
//
// Reference: https://leetcode.com/problems/reverse-words-in-a-string/
//
// Pattern: "reverse twice". Reverse the whole buffer, then reverse each word
// back in place; squeeze extra spaces while copying. That's the follow-up
// interviewers want for a mutable buffer with O(1) extra space.
//
//	"  the sky "  -> squeeze -> "the sky" -> reverse all -> "yks eht"
//	              -> reverse each word -> "sky the"
//
// Time: O(n). Space: O(n) for the library version; O(1) extra for the
// in-place []byte version.

// ReverseWords is the idiomatic one-liner-ish version.
func ReverseWords(s string) string {
	words := strings.Fields(s) // splits on runs of whitespace, drops empties
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	return strings.Join(words, " ")
}

// ReverseWordsInPlace rewrites buf (spaces = ' ' only) and returns the
// prefix that holds the result.
func ReverseWordsInPlace(buf []byte) []byte {
	// 1. Squeeze: copy words left, one space between them.
	w := 0
	for r := 0; r < len(buf); {
		if buf[r] == ' ' {
			r++
			continue
		}
		if w > 0 {
			buf[w] = ' '
			w++
		}
		for r < len(buf) && buf[r] != ' ' {
			buf[w] = buf[r]
			w, r = w+1, r+1
		}
	}
	buf = buf[:w]
	// 2. Reverse everything, 3. reverse each word back.
	reverseBytes(buf)
	for start := 0; start < len(buf); {
		end := start
		for end < len(buf) && buf[end] != ' ' {
			end++
		}
		reverseBytes(buf[start:end])
		start = end + 1
	}
	return buf
}

func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
