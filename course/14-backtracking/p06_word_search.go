package backtracking

// Problem: Word Search
// Difficulty: Medium
//
// Given a grid of letters and a word, decide whether the word can be spelled
// by a path of horizontally or vertically adjacent cells, using each cell at
// most once.
//
// Example:
//
//	A B C E
//	S F C S     "ABCCED" -> true    "SEE" -> true    "ABCB" -> false
//	A D E E
//
// Reference: https://leetcode.com/problems/word-search/
//
// Pattern: grid backtracking. From each cell matching word[0], DFS to
// neighbours matching the next letter. Mark the current cell as used before
// recursing (overwrite it with a sentinel) and restore it afterwards — that
// restore is the "unchoose" step, and forgetting it corrupts later searches.
//
// Pruning that matters on adversarial inputs:
//   - if the word is longer than the grid, or needs more of some letter
//     than the grid contains, answer false immediately;
//   - if the last letter is rarer in the grid than the first, search for
//     the reversed word instead (fewer starting points, earlier failures).
//
// Time: O(m·n·3^L) — 4 directions at the first step, then 3 (no going
// back). Space: O(L) recursion.

// Exist reports whether word can be traced in board. board is restored to
// its original contents before returning.
func Exist(board [][]byte, word string) bool {
	if len(word) == 0 {
		return true
	}
	m := len(board)
	if m == 0 {
		return false
	}
	n := len(board[0])
	if len(word) > m*n {
		return false
	}
	var have [256]int
	for _, row := range board {
		for _, ch := range row {
			have[ch]++
		}
	}
	var need [256]int
	for i := 0; i < len(word); i++ {
		if need[word[i]]++; need[word[i]] > have[word[i]] {
			return false // not enough of this letter anywhere
		}
	}
	if have[word[len(word)-1]] < have[word[0]] {
		word = reverse(word) // start from the rarer end
	}

	var dfs func(r, c, i int) bool
	dfs = func(r, c, i int) bool {
		if r < 0 || r >= m || c < 0 || c >= n || board[r][c] != word[i] {
			return false
		}
		if i == len(word)-1 {
			return true
		}
		ch := board[r][c]
		board[r][c] = 0 // choose: mark used (0 never appears in a word)
		found := dfs(r+1, c, i+1) || dfs(r-1, c, i+1) || dfs(r, c+1, i+1) || dfs(r, c-1, i+1)
		board[r][c] = ch // unchoose: restore
		return found
	}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if dfs(r, c, 0) {
				return true
			}
		}
	}
	return false
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
