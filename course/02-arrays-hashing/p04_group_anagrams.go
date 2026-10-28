package arrayshashing

import "slices"

// Problem 04 — Group Anagrams (Medium)
//
// Given a list of lowercase words, partition them into groups where every
// word in a group is an anagram of the others. Groups appear in the order of
// their first word in the input; words keep their input order inside a group.
//
//	["eat","tea","tan","ate","nat","bat"] → [["eat","tea","ate"], ["tan","nat"], ["bat"]]
//	[""]                                  → [[""]]
//
// Reference: https://leetcode.com/problems/group-anagrams/
//
// Pattern: pick a CANONICAL KEY that is equal for all anagrams, and bucket by
// it in a map. Two keys work:
//   - the sorted word ("aet") — simple, O(k log k) per word;
//   - the [26]int letter counts — O(k) per word. Go arrays are comparable,
//     so the array itself is a valid map key (slices are not!).
//
// Map iteration order is random, so we remember each key's group index to
// produce deterministic output.
//
//	GroupAnagramsSorted: Time O(n·k log k), Space O(n·k)
//	GroupAnagrams:       Time O(n·k),       Space O(n·k)   (k = max word length)

// GroupAnagramsSorted keys each word by its sorted letters.
func GroupAnagramsSorted(words []string) [][]string {
	groupIdx := map[string]int{}
	groups := [][]string{}
	for _, w := range words {
		b := []byte(w)
		slices.Sort(b)
		key := string(b)
		i, ok := groupIdx[key]
		if !ok {
			i = len(groups)
			groupIdx[key] = i
			groups = append(groups, nil)
		}
		groups[i] = append(groups[i], w)
	}
	return groups
}

// GroupAnagrams keys each word by its letter-count array.
func GroupAnagrams(words []string) [][]string {
	groupIdx := map[[26]int]int{}
	groups := [][]string{}
	for _, w := range words {
		key := letterCounts(w)
		i, ok := groupIdx[key]
		if !ok {
			i = len(groups)
			groupIdx[key] = i
			groups = append(groups, nil)
		}
		groups[i] = append(groups[i], w)
	}
	return groups
}
