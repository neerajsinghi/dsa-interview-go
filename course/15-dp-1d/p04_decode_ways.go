package dp1d

// Problem: Decode Ways
// Difficulty: Medium
//
// Letters are encoded as numbers: A=1, B=2, …, Z=26. Given a string of
// digits, count how many letter strings could have produced it. A piece
// with a leading zero ("06") is not a valid code, and neither is "0".
//
// Examples:
//
//	"12"  -> 2   ("AB" or "L")
//	"226" -> 3   ("BZ", "VF", "BBF")
//	"06"  -> 0
//	"10"  -> 1   (only "J")
//
// Reference: https://leetcode.com/problems/decode-ways/
//
// Pattern: Climbing Stairs with conditions. The last letter of a decoding
// uses either the last digit alone (valid if it is 1-9) or the last two
// digits together (valid if they form 10-26).
//
//	State:      ways(i) = decodings of the prefix s[:i]
//	Transition: ways(i) = [s[i-1] != '0'] · ways(i-1)
//	                    + [10 <= s[i-2:i] <= 26] · ways(i-2)
//	Base:       ways(0) = 1 (the empty prefix has one decoding: nothing)
//
//	s = "226":  ways = [1, 1, 2, 3]
//	                    ""  "2" "22" "226"
//
// Time: O(2^n) brute force; O(n) DP. Space: O(1) optimized.

// oneDigit reports whether s[i-1] alone is a letter.
func oneDigit(s string, i int) bool { return s[i-1] != '0' }

// twoDigits reports whether s[i-2:i] is a letter (10..26).
func twoDigits(s string, i int) bool {
	if i < 2 || s[i-2] == '0' {
		return false
	}
	v := int(s[i-2]-'0')*10 + int(s[i-1]-'0')
	return v <= 26
}

// NumDecodingsRecursive is the brute force.
func NumDecodingsRecursive(s string) int {
	var ways func(i int) int
	ways = func(i int) int {
		if i == 0 {
			return 1
		}
		total := 0
		if oneDigit(s, i) {
			total += ways(i - 1)
		}
		if twoDigits(s, i) {
			total += ways(i - 2)
		}
		return total
	}
	return ways(len(s))
}

// NumDecodingsMemo caches ways(i).
func NumDecodingsMemo(s string) int {
	memo := make([]int, len(s)+1)
	for i := range memo {
		memo[i] = -1
	}
	var ways func(i int) int
	ways = func(i int) int {
		if i == 0 {
			return 1
		}
		if memo[i] != -1 {
			return memo[i]
		}
		total := 0
		if oneDigit(s, i) {
			total += ways(i - 1)
		}
		if twoDigits(s, i) {
			total += ways(i - 2)
		}
		memo[i] = total
		return total
	}
	return ways(len(s))
}

// NumDecodingsTab fills ways[0..n] bottom-up.
func NumDecodingsTab(s string) int {
	dp := make([]int, len(s)+1)
	dp[0] = 1
	for i := 1; i <= len(s); i++ {
		if oneDigit(s, i) {
			dp[i] += dp[i-1]
		}
		if twoDigits(s, i) {
			dp[i] += dp[i-2]
		}
	}
	return dp[len(s)]
}

// NumDecodings keeps only ways(i-1) and ways(i-2).
func NumDecodings(s string) int {
	twoBack, oneBack := 0, 1 // ways(-1) is never used; ways(0) = 1
	for i := 1; i <= len(s); i++ {
		cur := 0
		if oneDigit(s, i) {
			cur += oneBack
		}
		if twoDigits(s, i) {
			cur += twoBack
		}
		twoBack, oneBack = oneBack, cur
	}
	return oneBack
}
