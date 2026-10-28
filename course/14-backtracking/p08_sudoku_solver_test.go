package backtracking

import (
	"reflect"
	"testing"
)

// isSolved checks every row, column and box holds 1..9 once and that the
// original clues were kept.
func isSolved(board, clues [][]byte) bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if clues[r][c] != '.' && clues[r][c] != board[r][c] {
				return false
			}
		}
	}
	for i := 0; i < 9; i++ {
		var row, col, box [10]bool
		for j := 0; j < 9; j++ {
			cells := [3]byte{board[i][j], board[j][i], board[(i/3)*3+j/3][(i%3)*3+j%3]}
			for k, seen := range []*[10]bool{&row, &col, &box} {
				ch := cells[k]
				if ch < '1' || ch > '9' || seen[ch-'0'] {
					return false
				}
				seen[ch-'0'] = true
			}
		}
	}
	return true
}

func copyBoard(b [][]byte) [][]byte {
	out := make([][]byte, len(b))
	for i := range b {
		out[i] = append([]byte(nil), b[i]...)
	}
	return out
}

func TestSolveSudoku(t *testing.T) {
	tests := []struct {
		name   string
		rows   []string
		solved bool
	}{
		{"classic", []string{
			"53..7....", "6..195...", ".98....6.",
			"8...6...3", "4..8.3..1", "7...2...6",
			".6....28.", "...419..5", "....8..79",
		}, true},
		{"hard (few clues)", []string{
			"8........", "..36.....", ".7..9.2..",
			".5...7...", "....457..", "...1...3.",
			"..1....68", "..85...1.", ".9....4..",
		}, true},
		{"empty grid", []string{
			".........", ".........", ".........",
			".........", ".........", ".........",
			".........", ".........", ".........",
		}, true},
		{"conflicting clues", []string{
			"55.......", ".........", ".........",
			".........", ".........", ".........",
			".........", ".........", ".........",
		}, false},
		{"consistent clues, no solution", []string{
			"12345678.", "........9", ".........",
			".........", ".........", ".........",
			".........", ".........", ".........",
		}, false},
	}
	for _, tt := range tests {
		b := board(tt.rows...)
		clues := copyBoard(b)
		got := SolveSudoku(b)
		if got != tt.solved {
			t.Errorf("%s: SolveSudoku = %v, want %v", tt.name, got, tt.solved)
			continue
		}
		if got && !isSolved(b, clues) {
			t.Errorf("%s: board is not a valid solution: %q", tt.name, b)
		}
		if !got && !reflect.DeepEqual(b, clues) {
			t.Errorf("%s: unsolvable board was modified", tt.name)
		}
	}
}
