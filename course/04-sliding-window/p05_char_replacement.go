package slidingwindow

// Problem 05 — Longest Repeating Character Replacement (Medium)
//
// s contains only uppercase English letters. You may change at most k
// characters to any other uppercase letter. Return the length of the longest
// substring that can be made of one repeated letter.
//
//	s = "ABAB",    k = 2 → 4   (change both B's → "AAAA")
//	s = "AABABBA", k = 1 → 4   ("AABA" → "AAAA")
//
// Reference: https://leetcode.com/problems/longest-repeating-character-replacement/
//
// Pattern: a window is fixable when (window length − count of its most
// frequent letter) ≤ k: replace everything that isn't the majority letter.
// Grow right; if the window becomes unfixable, slide left by one.
//
// Subtle trick: when the window becomes unfixable we SLIDE it (left and
// right both advance) instead of shrinking it, so its length never decreases.
// That lets us never decrease maxFreq either: a stale maxFreq just keeps the
// window sliding at the best size found so far, and the window can only grow
// when some letter's count beats the old maxFreq — exactly when a longer
// valid window exists. The final window length, len(s) - left, is the answer.
//
//	CharacterReplacementBrute: Time O(n²), Space O(1)
//	CharacterReplacement:      Time O(n),  Space O(1)

// CharacterReplacementBrute checks every substring incrementally.
func CharacterReplacementBrute(s string, k int) int {
	best := 0
	for i := 0; i < len(s); i++ {
		var count [26]int
		maxFreq := 0
		for j := i; j < len(s); j++ {
			count[s[j]-'A']++
			maxFreq = max(maxFreq, count[s[j]-'A'])
			if (j-i+1)-maxFreq <= k {
				best = max(best, j-i+1)
			}
		}
	}
	return best
}

// CharacterReplacement slides a window that never shrinks below the best size.
func CharacterReplacement(s string, k int) int {
	var count [26]int
	maxFreq, left := 0, 0
	for right := 0; right < len(s); right++ {
		count[s[right]-'A']++
		maxFreq = max(maxFreq, count[s[right]-'A'])
		if (right-left+1)-maxFreq > k {
			// Too many replacements needed: slide (not shrink) the window.
			count[s[left]-'A']--
			left++
		}
	}
	return len(s) - left
}
