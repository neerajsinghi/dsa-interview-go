package graphs

// Problem: Alien Dictionary
// Difficulty: Hard
//
// An alien language uses lowercase English letters in an unknown order. You
// receive its dictionary: a list of words already sorted by the alien
// order. Deduce an ordering of all letters that appear, consistent with the
// list. Return any valid ordering as a string, or "" if the list cannot be
// sorted under any ordering.
//
// Examples:
//
//	["wrt","wrf","er","ett","rftt"] -> "wertf"
//	["z","x","z"]                   -> ""   (z<x and x<z)
//	["abc","ab"]                    -> ""   (a longer word cannot precede its own prefix)
//
// Reference: https://leetcode.com/problems/alien-dictionary/ (premium)
//
// Pattern: build a graph, then topologically sort it. Each pair of ADJACENT
// words gives at most one fact: at the first position where they differ, the
// earlier word's letter comes first. (Letters after that position tell you
// nothing.) Every letter seen is a node; a topological order of the
// resulting graph is an alphabet, and a cycle means no alphabet exists.
//
//	wrt / wrf  -> t < f
//	wrf / er   -> w < e
//	er  / ett  -> r < t
//	ett / rftt -> e < r          graph: w -> e -> r -> t -> f
//
// Time: O(C) where C is the total number of characters (the graph has at
// most 26 nodes and 26² edges). Space: O(1) beyond the input (26-letter graph).

// AlienOrder returns one letter ordering consistent with words, or "".
func AlienOrder(words []string) string {
	var present [26]bool
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			present[w[i]-'a'] = true
		}
	}
	var edge [26][26]bool // adjacency matrix dedupes repeated facts
	var indeg [26]int
	for i := 0; i+1 < len(words); i++ {
		a, b := words[i], words[i+1]
		j := 0
		for j < len(a) && j < len(b) && a[j] == b[j] {
			j++
		}
		switch {
		case j < len(a) && j < len(b):
			u, v := a[j]-'a', b[j]-'a'
			if !edge[u][v] {
				edge[u][v] = true
				indeg[v]++
			}
		case len(a) > len(b):
			return "" // "abc" before "ab": impossible in any order
		}
	}

	// Kahn's algorithm over the present letters.
	var queue []int
	letters := 0
	for c := 0; c < 26; c++ {
		if present[c] {
			letters++
			if indeg[c] == 0 {
				queue = append(queue, c)
			}
		}
	}
	order := make([]byte, 0, letters)
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, byte('a'+u))
		for v := 0; v < 26; v++ {
			if edge[u][v] {
				if indeg[v]--; indeg[v] == 0 {
					queue = append(queue, v)
				}
			}
		}
	}
	if len(order) < letters {
		return "" // leftover letters sit on a cycle
	}
	return string(order)
}
