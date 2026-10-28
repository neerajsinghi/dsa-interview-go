package stringalgos

import "math"

// Problem: String to Integer (atoi)
// Difficulty: Medium
//
// Parse a signed 32-bit integer from the start of s with these rules:
//  1. skip leading ' ' characters;
//  2. read an optional '+' or '-';
//  3. read as many decimal digits as follow (leading zeros fine) and stop at
//     the first non-digit;
//  4. if no digits were read the result is 0;
//  5. clamp anything outside [-2^31, 2^31-1] to the nearest bound.
//
// Examples:
//
//	"42"              -> 42
//	"   -042abc"      -> -42
//	"words 99"        -> 0
//	"-91283472332"    -> -2147483648 (clamped)
//
// Reference: https://leetcode.com/problems/string-to-integer-atoi/
//
// Pattern: a small state machine plus overflow checking BEFORE the multiply.
// Accumulate the magnitude as a positive number and test
//
//	result > (limit - digit) / 10
//
// which is the rearranged form of result*10 + digit > limit, so the check
// itself can't overflow. Use limit = MaxInt32 for '+' and MaxInt32+1 for '-'
// (|MinInt32| is one bigger than MaxInt32).
//
// Time: O(n). Space: O(1).

// MyAtoi parses s according to the rules above.
func MyAtoi(s string) int32 {
	i, n := 0, len(s)
	for i < n && s[i] == ' ' {
		i++
	}
	negative := false
	if i < n && (s[i] == '+' || s[i] == '-') {
		negative = s[i] == '-'
		i++
	}
	limit := int64(math.MaxInt32) // largest magnitude we may return
	if negative {
		limit++ // 2147483648
	}
	var result int64 // int64 only so limit+1 fits; the check never needs > 32 bits
	for ; i < n && s[i] >= '0' && s[i] <= '9'; i++ {
		d := int64(s[i] - '0')
		if result > (limit-d)/10 { // result*10 + d would exceed limit
			result = limit
			break
		}
		result = result*10 + d
	}
	if negative {
		return int32(-result)
	}
	return int32(result)
}
