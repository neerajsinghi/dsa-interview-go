package stacks

import "strings"

// Problem: Decode String
// Difficulty: Medium
//
// A string is encoded with the rule k[text], meaning "text repeated k times".
// Encodings can nest, k is a positive integer (possibly several digits), and
// letters outside any brackets appear as-is. Expand the string. The input is
// assumed to be well-formed.
//
// Examples:
//
//	"3[a]2[bc]"  -> "aaabcbc"
//	"3[a2[c]]"   -> "accaccacc"
//	"2[ab]3[x]y" -> "ababxxxy"
//
// Reference: https://leetcode.com/problems/decode-string/
//
// Pattern: a stack of saved contexts. On '[' push (what we built so far,
// repeat count) and start fresh; on ']' pop the context and append the inner
// result k times to it. It's the same shape as evaluating nested parentheses.
//
// Time: O(output length). Space: O(output length).

type decodeFrame struct {
	prefix string // text built before the '['
	k      int    // how many times to repeat the bracketed part
}

// DecodeString expands a k[...] encoded string.
func DecodeString(s string) string {
	var st Stack[decodeFrame]
	var cur strings.Builder
	k := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			k = k*10 + int(c-'0') // multi-digit counts like "12[a]"
		case c == '[':
			st.Push(decodeFrame{prefix: cur.String(), k: k})
			cur.Reset()
			k = 0
		case c == ']':
			f, _ := st.Pop()
			inner := cur.String()
			cur.Reset()
			cur.WriteString(f.prefix)
			cur.WriteString(strings.Repeat(inner, f.k))
		default:
			cur.WriteByte(c)
		}
	}
	return cur.String()
}
