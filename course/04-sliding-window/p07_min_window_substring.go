package slidingwindow

// Problem 07 — Minimum Window Substring (Hard)
//
// Return the shortest substring of s that contains every byte of t, counting
// multiplicity (if t has two 'a's, the window needs two). If several are
// equally short, return the leftmost. Return "" if no window exists or t is
// empty.
//
//	s = "ADOBECODEBANC", t = "ABC" → "BANC"
//	s = "a",             t = "aa"  → ""
//
// Reference: https://leetcode.com/problems/minimum-window-substring/
//
// Pattern: variable window, "expand until valid, then shrink while valid".
// Keep need[c] = how many more of c the window still requires (negative means
// surplus) and `missing` = total required bytes not yet covered. The window is
// valid exactly when missing == 0.
//
//	ADOBECODEBANC
//	[ADOBEC]          first valid window (len 6)
//	 ...shrink: drop A → invalid; expand to …BA
//	     [CODEBA]     valid, len 6
//	         [BANC]   valid, len 4  ← answer
//
//	MinWindowBrute: Time O(n²·σ), Space O(σ)
//	MinWindow:      Time O(n + m), Space O(σ)  (σ = 256 byte values)

// MinWindowBrute tries every start and extends until the window covers t.
func MinWindowBrute(s, t string) string {
	if t == "" {
		return ""
	}
	var need [256]int
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	covers := func(have *[256]int) bool {
		for c := range 256 {
			if have[c] < need[c] {
				return false
			}
		}
		return true
	}
	best := ""
	for i := 0; i < len(s); i++ {
		var have [256]int
		for j := i; j < len(s); j++ {
			have[s[j]]++
			if covers(&have) {
				if best == "" || j-i+1 < len(best) {
					best = s[i : j+1]
				}
				break // longer windows from i are never better
			}
		}
	}
	return best
}

// MinWindow expands right until valid, then shrinks left while valid.
func MinWindow(s, t string) string {
	if t == "" {
		return ""
	}
	var need [256]int
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	missing := len(t)
	bestStart, bestLen := 0, -1
	left := 0
	for right := 0; right < len(s); right++ {
		c := s[right]
		if need[c] > 0 {
			missing-- // this byte was still required
		}
		need[c]--
		for missing == 0 {
			if bestLen == -1 || right-left+1 < bestLen {
				bestStart, bestLen = left, right-left+1
			}
			d := s[left]
			need[d]++
			if need[d] > 0 {
				missing++ // we just gave up a required byte
			}
			left++
		}
	}
	if bestLen == -1 {
		return ""
	}
	return s[bestStart : bestStart+bestLen]
}
