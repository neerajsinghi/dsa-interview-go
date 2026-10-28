package arrayshashing

import (
	"errors"
	"strconv"
	"strings"
)

// Problem 07 — Encode and Decode Strings (Medium)
//
// Design Encode, which turns a list of strings into a single string, and
// Decode, which recovers the exact original list. Strings may contain ANY
// characters, including your delimiter, and may be empty.
//
//	["lint","code","love","you"] → "4#lint4#code4#love3#you" → back again
//	["", "#"]                     → "0#1##"                   → back again
//
// Reference: https://leetcode.com/problems/encode-and-decode-strings/ (premium)
//
// Pattern: LENGTH-PREFIX framing, the same trick network protocols use.
// A delimiter alone is ambiguous if it can appear in the data; writing the
// byte length first means the decoder never has to look inside the payload.
//
//	Encode: Time O(total length), Space O(total length)
//	Decode: Time O(total length), Space O(total length)

// ErrMalformed is returned by Decode for input Encode could not have produced.
var ErrMalformed = errors.New("malformed encoding")

// Encode writes each string as "<byte length>#<bytes>".
func Encode(strs []string) string {
	var b strings.Builder
	for _, s := range strs {
		b.WriteString(strconv.Itoa(len(s))) // len is in BYTES, which is what we slice by
		b.WriteByte('#')
		b.WriteString(s)
	}
	return b.String()
}

// Decode reverses Encode.
func Decode(encoded string) ([]string, error) {
	out := []string{}
	i := 0
	for i < len(encoded) {
		// Read digits up to the '#'. The '#' search starts at i, and digits
		// never contain '#', so a '#' inside a payload can't confuse us.
		hash := strings.IndexByte(encoded[i:], '#')
		if hash <= 0 {
			return nil, ErrMalformed
		}
		n, err := strconv.Atoi(encoded[i : i+hash])
		if err != nil || n < 0 {
			return nil, ErrMalformed
		}
		start := i + hash + 1
		if start+n > len(encoded) {
			return nil, ErrMalformed
		}
		out = append(out, encoded[start:start+n])
		i = start + n
	}
	return out, nil
}
