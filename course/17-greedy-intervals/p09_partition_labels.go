package greedy

// Problem: Partition Labels
// Difficulty: Medium
//
// Cut a lowercase string into as many consecutive pieces as possible so that
// every letter appears in at most one piece. Return the piece lengths in
// order.
//
// Examples:
//
//	"ababcbacadefegdehijhklij" -> [9 7 8]  ("ababcbaca" "defegde" "hijhklij")
//	"eccbbbbdec"               -> [10]
//
// Reference: https://leetcode.com/problems/partition-labels/
//
// Pattern: intervals in disguise. Each letter spans [first, last]; the pieces
// are the merged intervals. Record each letter's last index, then sweep,
// extending the current piece's end to the last index of every letter seen;
// when i reaches that end, cut.
//
// Time: O(n). Space: O(1) — 26 counters.

// PartitionLabelsBrute tries every set of cut points and keeps the valid
// one with the most pieces. O(2^n · n); tests only.
func PartitionLabelsBrute(s string) []int {
	if s == "" {
		return nil
	}
	var best []int
	n := len(s)
	for mask := 0; mask < 1<<(n-1); mask++ { // bit i set = cut after s[i]
		var sizes []int
		owner := [26]int{}
		for i := range owner {
			owner[i] = -1
		}
		piece, start, ok := 0, 0, true
		for i := 0; i < n && ok; i++ {
			c := s[i] - 'a'
			if owner[c] != -1 && owner[c] != piece {
				ok = false
			}
			owner[c] = piece
			if i == n-1 || mask&(1<<i) != 0 {
				sizes = append(sizes, i-start+1)
				piece, start = piece+1, i+1
			}
		}
		if ok && len(sizes) > len(best) {
			best = sizes
		}
	}
	return best
}

// PartitionLabels is the last-occurrence sweep. Input must be 'a'..'z'.
func PartitionLabels(s string) []int {
	var last [26]int
	for i := 0; i < len(s); i++ {
		last[s[i]-'a'] = i
	}
	var sizes []int
	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		end = max(end, last[s[i]-'a']) // this piece must reach this letter's last spot
		if i == end {
			sizes = append(sizes, end-start+1)
			start = i + 1
		}
	}
	return sizes
}
