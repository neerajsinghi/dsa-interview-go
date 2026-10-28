package graphs

// Problem: Word Ladder
// Difficulty: Hard
//
// Transform beginWord into endWord by changing one letter at a time, where
// every intermediate word (and endWord itself) must appear in wordList.
// Return the number of words in the shortest such sequence, counting both
// ends, or 0 if no sequence exists. All words have the same length and are
// lowercase.
//
// Examples:
//
//	"hit" -> "cog", list [hot dot dog lot log cog] -> 5   (hit hot dot dog cog)
//	"hit" -> "cog", list [hot dot dog lot log]     -> 0   (cog is not allowed)
//
// Reference: https://leetcode.com/problems/word-ladder/
//
// Pattern: BFS on an implicit graph. Words are nodes; two words are adjacent
// if they differ in one letter. "Fewest steps" in an unweighted graph means
// BFS. Do not build the graph up front (O(N²·L)); generate neighbours on the
// fly by trying all 25 replacements at each position and keeping those in
// the word set.
//
// Bidirectional BFS grows one frontier from each end and always expands the
// smaller one. With branching factor b and distance d it explores about
// 2·b^(d/2) nodes instead of b^d.
//
//	hit ──> hot ──> dot ──> dog ──> cog
//	 front ───────────>   <─────── back   (they meet in the middle)
//
// Time: O(N · L · 26) where N = len(wordList), L = word length.
// Space: O(N · L).

// LadderLength runs a plain BFS from beginWord.
func LadderLength(beginWord, endWord string, wordList []string) int {
	words := make(map[string]bool, len(wordList))
	for _, w := range wordList {
		words[w] = true
	}
	if !words[endWord] {
		return 0
	}
	delete(words, beginWord) // never step back onto the start
	queue := []string{beginWord}
	for steps := 1; len(queue) > 0; steps++ {
		var next []string
		for _, w := range queue {
			if w == endWord {
				return steps
			}
			b := []byte(w)
			for i := range b {
				orig := b[i]
				for c := byte('a'); c <= 'z'; c++ {
					if c == orig {
						continue
					}
					b[i] = c
					if cand := string(b); words[cand] {
						delete(words, cand) // deleting = marking visited
						next = append(next, cand)
					}
				}
				b[i] = orig
			}
		}
		queue = next
	}
	return 0
}

// LadderLengthBidirectional expands the smaller of two frontiers each round.
func LadderLengthBidirectional(beginWord, endWord string, wordList []string) int {
	words := make(map[string]bool, len(wordList))
	for _, w := range wordList {
		words[w] = true
	}
	if !words[endWord] {
		return 0
	}
	if beginWord == endWord {
		return 1
	}
	front := map[string]bool{beginWord: true}
	back := map[string]bool{endWord: true}
	delete(words, beginWord)
	delete(words, endWord)
	for steps := 2; len(front) > 0 && len(back) > 0; steps++ {
		if len(front) > len(back) {
			front, back = back, front // always grow the cheaper side
		}
		next := map[string]bool{}
		for w := range front {
			b := []byte(w)
			for i := range b {
				orig := b[i]
				for c := byte('a'); c <= 'z'; c++ {
					if c == orig {
						continue
					}
					b[i] = c
					cand := string(b)
					if back[cand] {
						return steps // the frontiers touch
					}
					if words[cand] {
						delete(words, cand)
						next[cand] = true
					}
				}
				b[i] = orig
			}
		}
		front = next
	}
	return 0
}
