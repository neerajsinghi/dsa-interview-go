package backtracking

import "strings"

// Problem: N-Queens (I & II)
// Difficulty: Hard
//
// Place n chess queens on an n×n board so that no two attack each other (no
// shared row, column or diagonal). Part I: return every such board, drawn
// with 'Q' for a queen and '.' for an empty square. Part II: just count them.
//
// Example (n = 4, two solutions):
//
//	. Q . .     . . Q .
//	. . . Q     Q . . .
//	Q . . .     . . . Q
//	. . Q .     . Q . .
//
// References: https://leetcode.com/problems/n-queens/
//
//	https://leetcode.com/problems/n-queens-ii/
//
// Pattern: place one queen per row, so rows never clash by construction.
// Track the attacked columns and diagonals in sets: every square on a
// "\" diagonal shares r-c, every square on a "/" diagonal shares r+c. A
// placement is legal iff its column, r-c and r+c are all free — an O(1)
// check instead of scanning the board.
//
// The counting version packs those three sets into bitmasks: the free
// squares of the next row are ~(cols | diag1 | diag2), and the diagonal
// masks simply shift by one as you move down a row.
//
// Time: O(n!) — row r has at most n-r safe columns. Space: O(n).

// SolveNQueens returns every solution as a slice of n row strings.
func SolveNQueens(n int) [][]string {
	var res [][]string
	queenCol := make([]int, n) // queenCol[r] = column of the queen in row r
	cols := make([]bool, n)
	diag := make([]bool, 2*n) // r-c+n   (shifted to be non-negative)
	anti := make([]bool, 2*n) // r+c
	var place func(r int)
	place = func(r int) {
		if r == n {
			res = append(res, render(queenCol))
			return
		}
		for c := 0; c < n; c++ {
			if cols[c] || diag[r-c+n] || anti[r+c] {
				continue // attacked
			}
			queenCol[r] = c
			cols[c], diag[r-c+n], anti[r+c] = true, true, true
			place(r + 1)
			cols[c], diag[r-c+n], anti[r+c] = false, false, false
		}
	}
	place(0)
	return res
}

func render(queenCol []int) []string {
	n := len(queenCol)
	rows := make([]string, n)
	for r, c := range queenCol {
		rows[r] = strings.Repeat(".", c) + "Q" + strings.Repeat(".", n-c-1)
	}
	return rows
}

// TotalNQueens counts solutions with bitmasks.
func TotalNQueens(n int) int {
	full := 1<<n - 1 // n ones: every column
	var count func(cols, d1, d2 int) int
	count = func(cols, d1, d2 int) int {
		if cols == full {
			return 1 // a queen in every column => n queens placed
		}
		total := 0
		free := full &^ (cols | d1 | d2) // &^ is Go's AND NOT
		for free != 0 {
			bit := free & -free // lowest free square
			free &^= bit
			// Moving down one row, "\" attacks shift left and "/" attacks right.
			total += count(cols|bit, (d1|bit)<<1&full, (d2|bit)>>1)
		}
		return total
	}
	return count(0, 0, 0)
}
