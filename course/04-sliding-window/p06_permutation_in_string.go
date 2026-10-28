package slidingwindow

// Problem 06 — Permutation in String (Medium)
//
// s1 and s2 are lowercase English strings. Report whether some rearrangement
// of s1 appears as a contiguous substring of s2.
//
//	s1 = "ab", s2 = "eidbaooo" → true   ("ba")
//	s1 = "ab", s2 = "eidboaoo" → false
//
// Reference: https://leetcode.com/problems/permutation-in-string/
//
// Pattern: fixed-size window of len(s1) over s2, comparing letter counts.
// Instead of comparing all 26 counters per step, track how many of the 26
// letters currently have matching counts; each slide touches only 2 letters.
//
//	CheckInclusionBrute: Time O(n·m),  Space O(1)  (n = len(s2), m = len(s1))
//	CheckInclusion:      Time O(n + m), Space O(1)

// CheckInclusionBrute recounts every window.
func CheckInclusionBrute(s1, s2 string) bool {
	m := len(s1)
	if m > len(s2) {
		return false
	}
	want := counts(s1)
	for start := 0; start+m <= len(s2); start++ {
		if counts(s2[start:start+m]) == want {
			return true
		}
	}
	return false
}

func counts(s string) [26]int {
	var c [26]int
	for i := 0; i < len(s); i++ {
		c[s[i]-'a']++
	}
	return c
}

// CheckInclusion maintains window counts and a "matching letters" tally.
func CheckInclusion(s1, s2 string) bool {
	m := len(s1)
	if m > len(s2) {
		return false
	}
	need, have := counts(s1), counts(s2[:m])
	matches := 0
	for i := range 26 {
		if need[i] == have[i] {
			matches++
		}
	}
	// update adjusts letter c by delta and keeps matches in sync.
	update := func(c byte, delta int) {
		i := c - 'a'
		if have[i] == need[i] {
			matches-- // was matching, is about to change
		}
		have[i] += delta
		if have[i] == need[i] {
			matches++
		}
	}
	for right := m; right < len(s2); right++ {
		if matches == 26 {
			return true
		}
		update(s2[right], +1)
		update(s2[right-m], -1)
	}
	return matches == 26
}
