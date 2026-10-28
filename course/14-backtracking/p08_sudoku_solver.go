package backtracking

import "math/bits"

// Problem: Sudoku Solver
// Difficulty: Hard
//
// Fill a 9×9 Sudoku grid in place. Empty cells hold '.', filled cells hold
// '1'..'9'. In the finished grid every row, every column and every 3×3 box
// contains each digit exactly once. Return whether a solution was found
// (the board is left unchanged when it was not).
//
// Example (first rows of a classic puzzle → solution):
//
//	5 3 . | . 7 . | . . .        5 3 4 | 6 7 8 | 9 1 2
//	6 . . | 1 9 5 | . . .   ->   6 7 2 | 1 9 5 | 3 4 8
//	. 9 8 | . . . | . 6 .        1 9 8 | 3 4 2 | 5 6 7
//	...                          ...
//
// Reference: https://leetcode.com/problems/sudoku-solver/
//
// Pattern: constraint backtracking. Keep, for every row, column and box, a
// 9-bit mask of digits already used; the candidates for cell (r, c) are the
// bits missing from rows[r] | cols[c] | boxes[b], computed in O(1).
//
// The big win is the choice of WHICH cell to fill next: always pick the
// empty cell with the fewest candidates ("minimum remaining values"). A cell
// with one candidate is forced, and a cell with zero candidates proves the
// current branch dead immediately, so the search tree stays tiny.
//
// Time: exponential in the number of empty cells in the worst case;
// milliseconds on real puzzles with MRV. Space: O(81) for the masks and
// recursion.

// SolveSudoku solves board in place and reports success.
func SolveSudoku(board [][]byte) bool {
	var rows, cols, boxes [9]uint16 // bit d set = digit d+1 used
	var empty [][2]int
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] == '.' {
				empty = append(empty, [2]int{r, c})
				continue
			}
			bit := uint16(1) << (board[r][c] - '1')
			b := (r/3)*3 + c/3
			if rows[r]&bit != 0 || cols[c]&bit != 0 || boxes[b]&bit != 0 {
				return false // the given clues already conflict
			}
			rows[r] |= bit
			cols[c] |= bit
			boxes[b] |= bit
		}
	}

	const all = 1<<9 - 1
	var solve func(remaining int) bool
	solve = func(remaining int) bool {
		if remaining == 0 {
			return true
		}
		// MRV: move the most constrained empty cell to position remaining-1.
		best, bestCount := -1, 10
		for i := 0; i < remaining; i++ {
			r, c := empty[i][0], empty[i][1]
			free := all &^ (rows[r] | cols[c] | boxes[(r/3)*3+c/3])
			if n := bits.OnesCount16(free); n < bestCount {
				best, bestCount = i, n
				if n == 0 {
					return false // dead end: some cell has no legal digit
				}
			}
		}
		last := remaining - 1
		empty[best], empty[last] = empty[last], empty[best]
		r, c := empty[last][0], empty[last][1]
		b := (r/3)*3 + c/3
		free := uint16(all) &^ (rows[r] | cols[c] | boxes[b])
		for free != 0 {
			bit := free & -free // lowest candidate digit
			free &^= bit
			rows[r] |= bit
			cols[c] |= bit
			boxes[b] |= bit
			board[r][c] = byte('1' + bits.TrailingZeros16(bit))
			if solve(remaining - 1) {
				return true
			}
			rows[r] &^= bit // unchoose
			cols[c] &^= bit
			boxes[b] &^= bit
			board[r][c] = '.'
		}
		return false
	}
	return solve(len(empty))
}
