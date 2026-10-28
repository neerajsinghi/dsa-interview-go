package tries

// Problem: Longest Common Prefix
// Difficulty: Easy
//
// Given a list of strings, return the longest string that every one of them
// starts with. If they share nothing (or the list is empty) return "".
//
// Examples:
//
//	["flower", "flow", "flight"] -> "fl"
//	["dog", "racecar", "car"]    -> ""
//
// Reference: https://leetcode.com/problems/longest-common-prefix/
//
// Pattern: in a trie of all the words, the common prefix is the single
// "trunk" from the root: keep walking while the current node has exactly one
// child and no word ends there. The vertical scan needs no extra memory and
// is what you would write for a one-off query; the trie pays off when the
// same word set is queried repeatedly (e.g. "LCP of the words with prefix p").
//
// Time: O(S) for both, where S is the total number of characters.
// Space: O(1) for the scan, O(S) for the trie.

// LongestCommonPrefixScan compares column by column: character i of the first
// word against character i of every other word, stopping at the first
// mismatch or the end of the shortest word.
func LongestCommonPrefixScan(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	first := strs[0]
	for i := 0; i < len(first); i++ {
		for _, s := range strs[1:] {
			if i == len(s) || s[i] != first[i] {
				return first[:i]
			}
		}
	}
	return first
}

// LongestCommonPrefixTrie inserts every word, then follows the trunk.
// Words must be lowercase a-z.
func LongestCommonPrefixTrie(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	t := New()
	for _, s := range strs {
		t.Insert(s)
	}
	var prefix []byte
	cur := t.root
	for !cur.IsEnd { // a word ending here caps the prefix
		only, next := -1, (*Node)(nil)
		for c, child := range cur.Children {
			if child == nil {
				continue
			}
			if only != -1 { // a second child: the words diverge here
				return string(prefix)
			}
			only, next = c, child
		}
		if next == nil {
			break
		}
		prefix = append(prefix, byte('a'+only))
		cur = next
	}
	return string(prefix)
}
