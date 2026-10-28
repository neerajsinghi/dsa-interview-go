package tries

import "strings"

// Problem: Replace Words
// Difficulty: Medium
//
// You are given a list of word roots and a sentence of space-separated
// lowercase words. Replace every word that starts with some root by the
// shortest such root; leave other words unchanged.
//
// Examples:
//
//	roots ["cat","bat","rat"], "the cattle was rattled by the battery"
//	  -> "the cat was rat by the bat"
//	roots ["a","aa"], "aadsfasf absbs bbab" -> "a a bbab"
//
// Reference: https://leetcode.com/problems/replace-words/
//
// Pattern: shortest-prefix lookup. Put the roots in a trie and walk each word
// from the root; the first IsEnd node you hit is the shortest root. Stop as
// soon as the walk falls off the trie.
//
// Time: O(R + S) where R is the total length of the roots and S the length
// of the sentence. Brute force: O(words · roots · L).
// Space: O(R) for the trie.

// ReplaceWords replaces each word by its shortest root using a trie.
func ReplaceWords(roots []string, sentence string) string {
	t := New()
	for _, r := range roots {
		t.Insert(r)
	}
	words := strings.Fields(sentence)
	for i, w := range words {
		words[i] = shortestRoot(t, w)
	}
	return strings.Join(words, " ")
}

// shortestRoot returns the shortest stored prefix of w, or w itself.
func shortestRoot(t *Trie, w string) string {
	cur := t.root
	for i := 0; i < len(w); i++ {
		cur = cur.Children[w[i]-'a']
		if cur == nil {
			return w // no root continues this way
		}
		if cur.IsEnd {
			return w[:i+1]
		}
	}
	return w
}

// ReplaceWordsBrute tries every root on every word and keeps the shortest
// match.
func ReplaceWordsBrute(roots []string, sentence string) string {
	words := strings.Fields(sentence)
	for i, w := range words {
		best := w
		for _, r := range roots {
			if strings.HasPrefix(w, r) && len(r) < len(best) {
				best = r
			}
		}
		words[i] = best
	}
	return strings.Join(words, " ")
}
