package slidingwindow

// Problem 04 — Longest Substring Without Repeating Characters (Medium)
//
// Return the length of the longest substring of s in which no byte occurs
// twice. (Treat s as bytes; convert to []rune first if you need Unicode.)
//
//	"abcabcbb" → 3   ("abc")
//	"pwwkew"   → 3   ("wke")
//
// Reference: https://leetcode.com/problems/longest-substring-without-repeating-characters/
//
// Pattern: variable window with a "last seen at" table. When s[right] already
// occurs inside the window, jump left to just past that occurrence — no need
// to shrink one step at a time.
//
//	s = a b c a b c b b
//	    l     r          'a' last seen at 0 → left = 1
//	      l     r        'b' last seen at 1 → left = 2
//
//	LengthOfLongestSubstringBrute: Time O(n²), Space O(σ) (σ = alphabet size)
//	LengthOfLongestSubstring:      Time O(n),  Space O(σ)

// LengthOfLongestSubstringBrute extends a window from every start until a
// repeat appears.
func LengthOfLongestSubstringBrute(s string) int {
	best := 0
	for i := 0; i < len(s); i++ {
		var seen [256]bool
		for j := i; j < len(s); j++ {
			if seen[s[j]] {
				break
			}
			seen[s[j]] = true
			best = max(best, j-i+1)
		}
	}
	return best
}

// LengthOfLongestSubstring keeps, for every byte, one past the index where it
// was last seen (0 means "never"), so a fresh [256]int needs no initialisation.
func LengthOfLongestSubstring(s string) int {
	var nextAfter [256]int
	best, left := 0, 0
	for right := 0; right < len(s); right++ {
		c := s[right]
		// Only jump forward: an old occurrence left of the window is irrelevant.
		left = max(left, nextAfter[c])
		nextAfter[c] = right + 1
		best = max(best, right-left+1)
	}
	return best
}
