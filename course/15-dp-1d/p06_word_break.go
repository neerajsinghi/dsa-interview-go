package dp1d

// Problem: Word Break
// Difficulty: Medium
//
// Given a string and a dictionary of words, decide whether the string can be
// split into a sequence of one or more dictionary words (words may be reused).
//
// Examples:
//
//	"leetcode",      [leet code]               -> true
//	"applepenapple", [apple pen]               -> true  (apple is reused)
//	"catsandog",     [cats dog sand and cat]   -> false
//
// Reference: https://leetcode.com/problems/word-break/
//
// Pattern: prefix DP. The last word of a valid split is some dictionary word
// w that ends the string; everything before it must itself be splittable.
//
//	State:      ok[i] = the prefix s[:i] can be segmented
//	Transition: ok[i] = OR over j < i of (ok[j] AND s[j:i] in dict)
//	Base:       ok[0] = true (empty prefix)
//	Answer:     ok[n]
//
//	s = "leetcode":  ok[0]=T  ok[4]=T ("leet")  ok[8]=T ("code" after ok[4])
//
// Only lengths that actually occur in the dictionary need to be tried, which
// bounds the inner loop by the longest word.
//
// Time: brute force exponential; DP O(n · L) substring checks with L the
// longest word length (each check O(L) to hash). Space: O(n + dict).

// WordBreakRecursive tries every dictionary word as the next piece.
func WordBreakRecursive(s string, wordDict []string) bool {
	var can func(start int) bool
	can = func(start int) bool {
		if start == len(s) {
			return true
		}
		for _, w := range wordDict {
			if len(w) > 0 && len(s)-start >= len(w) && s[start:start+len(w)] == w && can(start+len(w)) {
				return true
			}
		}
		return false
	}
	return can(0)
}

// WordBreakMemo caches can(start): 0 unknown, 1 true, 2 false.
func WordBreakMemo(s string, wordDict []string) bool {
	memo := make([]int8, len(s)+1)
	var can func(start int) bool
	can = func(start int) bool {
		if start == len(s) {
			return true
		}
		if memo[start] != 0 {
			return memo[start] == 1
		}
		result := false
		for _, w := range wordDict {
			if len(w) > 0 && len(s)-start >= len(w) && s[start:start+len(w)] == w && can(start+len(w)) {
				result = true
				break
			}
		}
		memo[start] = 2
		if result {
			memo[start] = 1
		}
		return result
	}
	return can(0)
}

// WordBreak is the bottom-up prefix DP with a hash set of words.
func WordBreak(s string, wordDict []string) bool {
	dict := make(map[string]bool, len(wordDict))
	maxLen := 0
	for _, w := range wordDict {
		dict[w] = true
		maxLen = max(maxLen, len(w))
	}
	ok := make([]bool, len(s)+1)
	ok[0] = true
	for i := 1; i <= len(s); i++ {
		// Only look back as far as the longest word.
		for j := i - 1; j >= 0 && i-j <= maxLen; j-- {
			if ok[j] && dict[s[j:i]] {
				ok[i] = true
				break
			}
		}
	}
	return ok[len(s)]
}
