package tries

import "sort"

// Problem: Word Search II
// Difficulty: Hard
//
// Given a grid of lowercase letters and a list of words, return every word
// that can be traced in the grid by moving between horizontally or
// vertically adjacent cells, using each cell at most once per word.
//
// Example:
//
//	o a a n
//	e t a e      words ["oath","pea","eat","rain"] -> ["eat","oath"]
//	i h k r
//	i f l v
//
// Reference: https://leetcode.com/problems/word-search-ii/
//
// Pattern: trie + grid DFS. Running a separate word search per word repeats
// the same prefix exploration again and again. Instead, put all words in a
// trie and run ONE DFS from every cell, stepping down the trie as you step
// across the grid. A path is abandoned the moment its letters stop being a
// prefix of any word. Three refinements make it fast:
//   - store the full word at its end node, so no string building is needed;
//   - clear that word once found, so it is reported only once;
//   - prune trie leaves that have been fully used, so later DFS runs skip them.
//
// Time: O(m·n·4·3^(L-1)) worst case for max word length L (the first step has
// 4 directions, later ones 3 since you cannot go back), usually far less
// thanks to pruning. Brute force multiplies that by the number of words.
// Space: O(total characters of words) for the trie, O(L) recursion.

// wordNode is a trie node that remembers the word ending at it.
type wordNode struct {
	children [alphabet]*wordNode
	word     string // non-empty only at the end of a not-yet-found word
}

// FindWords returns the words present in board, sorted.
func FindWords(board [][]byte, words []string) []string {
	root := &wordNode{}
	for _, w := range words {
		cur := root
		for i := 0; i < len(w); i++ {
			c := w[i] - 'a'
			if cur.children[c] == nil {
				cur.children[c] = &wordNode{}
			}
			cur = cur.children[c]
		}
		cur.word = w
	}

	var found []string
	var dfs func(r, c int, parent *wordNode)
	dfs = func(r, c int, parent *wordNode) {
		ch := board[r][c]
		node := parent.children[ch-'a']
		if node == nil {
			return
		}
		if node.word != "" {
			found = append(found, node.word)
			node.word = "" // report each word once
		}
		board[r][c] = '#' // mark visited for this path
		for _, d := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nr, nc := r+d[0], c+d[1]
			if nr >= 0 && nr < len(board) && nc >= 0 && nc < len(board[nr]) && board[nr][nc] != '#' {
				dfs(nr, nc, node)
			}
		}
		board[r][c] = ch // restore on the way back (backtracking)
		if node.word == "" && isLeaf(node) {
			parent.children[ch-'a'] = nil // nothing left to find below: prune
		}
	}

	for r := range board {
		for c := range board[r] {
			dfs(r, c, root)
		}
	}
	sort.Strings(found)
	return found
}

func isLeaf(n *wordNode) bool {
	for _, child := range n.children {
		if child != nil {
			return false
		}
	}
	return true
}

// FindWordsBrute runs a classic single-word search for every distinct word.
func FindWordsBrute(board [][]byte, words []string) []string {
	seen := map[string]bool{}
	var found []string
	for _, w := range words {
		if !seen[w] && existsInGrid(board, w) {
			found = append(found, w)
		}
		seen[w] = true
	}
	sort.Strings(found)
	return found
}

func existsInGrid(board [][]byte, w string) bool {
	var dfs func(r, c, i int) bool
	dfs = func(r, c, i int) bool {
		if r < 0 || r >= len(board) || c < 0 || c >= len(board[r]) || board[r][c] != w[i] {
			return false
		}
		if i == len(w)-1 {
			return true
		}
		ch := board[r][c]
		board[r][c] = '#'
		ok := dfs(r+1, c, i+1) || dfs(r-1, c, i+1) || dfs(r, c+1, i+1) || dfs(r, c-1, i+1)
		board[r][c] = ch
		return ok
	}
	for r := range board {
		for c := range board[r] {
			if dfs(r, c, 0) {
				return true
			}
		}
	}
	return false
}
