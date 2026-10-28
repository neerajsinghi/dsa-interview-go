package heaps

import (
	"cmp"
	"slices"
)

// Problem: Top K Frequent Words
// Difficulty: Medium
//
// Given a list of words and k, return the k most frequent words, ordered from
// most to least frequent. Words with equal frequency are ordered
// alphabetically. k is at most the number of distinct words.
//
// Examples:
//
//	["i","love","go","i","love","code"], k=2 -> ["i","love"]
//	["the","day","is","sunny","the","the","the","sunny","is","is"], k=4
//	    -> ["the","is","sunny","day"]
//
// Reference: https://leetcode.com/problems/top-k-frequent-words/
//
// Pattern: count with a map, then top-k with a size-k heap whose top is the
// *worst* of the current k (lowest count; for equal counts, the alphabetically
// LAST word). Pop everything and reverse to get best-first order.
//
// Time: sort O(n + u log u); heap O(n + u log k), u = distinct words.
// Space: O(u).

type wordCount struct {
	word  string
	count int
}

// better reports whether a ranks ahead of b in the final answer.
func better(a, b wordCount) bool {
	if a.count != b.count {
		return a.count > b.count
	}
	return a.word < b.word
}

func countWords(words []string) map[string]int {
	counts := make(map[string]int)
	for _, w := range words {
		counts[w]++
	}
	return counts
}

// TopKFrequentSort ranks every distinct word and takes the first k.
func TopKFrequentSort(words []string, k int) []string {
	var all []wordCount
	for w, c := range countWords(words) {
		all = append(all, wordCount{w, c})
	}
	slices.SortFunc(all, func(a, b wordCount) int {
		if c := cmp.Compare(b.count, a.count); c != 0 {
			return c
		}
		return cmp.Compare(a.word, b.word)
	})
	out := make([]string, 0, k)
	for _, wc := range all[:k] {
		out = append(out, wc.word)
	}
	return out
}

// TopKFrequent uses a size-k heap with the worst candidate on top.
func TopKFrequent(words []string, k int) []string {
	// "Less" for this heap means "worse", so the worst sits at the top.
	h := New(func(a, b wordCount) bool { return better(b, a) })
	for w, c := range countWords(words) {
		h.Push(wordCount{w, c})
		if h.Len() > k {
			h.Pop() // evict the worst
		}
	}
	out := make([]string, h.Len())
	for i := len(out) - 1; i >= 0; i-- { // pops come out worst-first
		out[i] = h.Pop().word
	}
	return out
}
