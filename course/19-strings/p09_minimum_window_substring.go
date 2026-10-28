package stringalgos

// Problem: Minimum Window Substring
// Difficulty: Hard
//
// Given strings s and t, return the shortest substring of s that contains
// every character of t, counting duplicates (t="aab" needs two 'a's). If no
// window works return "". If several shortest windows exist return the
// leftmost.
//
// Examples:
//
//	s="ADOBECODEBANC", t="ABC" -> "BANC"
//	s="a", t="aa"              -> ""
//
// Reference: https://leetcode.com/problems/minimum-window-substring/
//
// Pattern: variable-size sliding window with a "missing" counter.
// need[c] = how many more c's the window still needs (can go negative when
// the window has surplus). missing = total chars still needed.
//
//   - Move right: if need[c] > 0 this char was useful, missing--. need[c]--.
//   - While missing == 0 the window is valid: record it, then shrink from
//     the left; when a dropped char makes need[c] > 0, the window breaks.
//
// Each index enters and leaves the window once.
//
// Time: O(|s| + |t|). Space: O(1) — a 256-entry array for byte strings.
// Brute force: O(|s|² · alphabet).

// MinWindowBrute checks every start, extending until the window is valid.
func MinWindowBrute(s, t string) string {
	if t == "" {
		return ""
	}
	best := ""
	for i := range s {
		var need [256]int
		missing := len(t)
		for k := 0; k < len(t); k++ {
			need[t[k]]++
		}
		for j := i; j < len(s); j++ {
			if need[s[j]] > 0 {
				missing--
			}
			need[s[j]]--
			if missing == 0 {
				if best == "" || j-i+1 < len(best) {
					best = s[i : j+1]
				}
				break
			}
		}
	}
	return best
}

// MinWindow is the linear sliding window.
func MinWindow(s, t string) string {
	if t == "" || len(t) > len(s) {
		return ""
	}
	var need [256]int
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	missing := len(t)
	bestL, bestLen := 0, len(s)+1
	left := 0
	for right := 0; right < len(s); right++ {
		if need[s[right]] > 0 {
			missing--
		}
		need[s[right]]--
		for missing == 0 { // window s[left..right] is valid: try to shrink
			if right-left+1 < bestLen {
				bestL, bestLen = left, right-left+1
			}
			need[s[left]]++
			if need[s[left]] > 0 { // we just dropped a required char
				missing++
			}
			left++
		}
	}
	if bestLen > len(s) {
		return ""
	}
	return s[bestL : bestL+bestLen]
}
