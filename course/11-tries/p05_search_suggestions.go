package tries

import (
	"sort"
	"strings"
)

// Problem: Search Suggestions System
// Difficulty: Medium
//
// You have a catalogue of distinct lowercase product names and a search word
// typed one character at a time. After each character, suggest up to three
// products that start with what has been typed so far, choosing the
// lexicographically smallest ones. Return one suggestion list per character.
//
// Example:
//
//	products ["mobile","mouse","moneypot","monitor","mousepad"], search "mouse"
//	  m     -> [mobile moneypot monitor]
//	  mo    -> [mobile moneypot monitor]
//	  mou   -> [mouse mousepad]
//	  mous  -> [mouse mousepad]
//	  mouse -> [mouse mousepad]
//
// Reference: https://leetcode.com/problems/search-suggestions-system/
//
// Pattern: autocomplete. Walk one trie edge per typed character; from the
// current node a DFS that visits children in 'a'..'z' order produces words in
// sorted order, so stop after three. Once the walk falls off the trie every
// remaining answer is empty.
//
// A second classic approach sorts the products and binary-searches each
// prefix: the matches are the (up to) three entries starting at the insertion
// point. It uses no extra structure and is often the quickest to write.
//
// Time: trie O(P + L·k) where P is total product length, L the search length
// and k the work to collect three words; sort+search O(P log n + L·log n).
// Space: O(P) for the trie.

// SuggestedProducts answers with a trie.
func SuggestedProducts(products []string, searchWord string) [][]string {
	t := New()
	for _, p := range products {
		t.Insert(p)
	}
	res := make([][]string, len(searchWord))
	cur := t.root
	for i := 0; i < len(searchWord); i++ {
		if cur != nil {
			cur = cur.Children[searchWord[i]-'a']
		}
		if cur == nil {
			continue // fell off: res[i] stays nil for the rest
		}
		res[i] = t.WithPrefix(searchWord[:i+1], 3)
	}
	return res
}

// SuggestedProductsSorted sorts once and binary-searches every prefix.
func SuggestedProductsSorted(products []string, searchWord string) [][]string {
	sorted := append([]string(nil), products...) // don't reorder the caller's slice
	sort.Strings(sorted)
	res := make([][]string, len(searchWord))
	for i := range searchWord {
		prefix := searchWord[:i+1]
		start := sort.SearchStrings(sorted, prefix) // first entry >= prefix
		for j := start; j < len(sorted) && j < start+3; j++ {
			if !strings.HasPrefix(sorted[j], prefix) {
				break
			}
			res[i] = append(res[i], sorted[j])
		}
	}
	return res
}
