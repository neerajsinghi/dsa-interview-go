package backtracking

import (
	"reflect"
	"testing"
)

func board(rows ...string) [][]byte {
	b := make([][]byte, len(rows))
	for i, r := range rows {
		b[i] = []byte(r)
	}
	return b
}

func TestExist(t *testing.T) {
	classic := []string{"ABCE", "SFCS", "ADEE"}
	tests := []struct {
		name  string
		board [][]byte
		word  string
		want  bool
	}{
		{"snake path", board(classic...), "ABCCED", true},
		{"short", board(classic...), "SEE", true},
		{"would reuse a cell", board(classic...), "ABCB", false},
		{"single cell hit", board("a"), "a", true},
		{"single cell miss", board("a"), "b", false},
		{"too long", board("ab"), "aba", false},
		{"empty word", board("a"), "", true},
		{"reversed search path", board("aab"), "baa", true},
		{"needs every cell", board("ab", "dc"), "abcd", true},
		{"letter count prune", board("aaaa", "aaaa", "aaaa"), "aaaaaaaaaaaab", false},
		{"no diagonal moves", board("ab", "cd"), "ad", false},
	}
	for _, tt := range tests {
		before := make([][]byte, len(tt.board))
		for i, row := range tt.board {
			before[i] = append([]byte(nil), row...)
		}
		if got := Exist(tt.board, tt.word); got != tt.want {
			t.Errorf("%s: Exist(%q) = %v, want %v", tt.name, tt.word, got, tt.want)
		}
		if !reflect.DeepEqual(before, tt.board) {
			t.Errorf("%s: board was not restored: %q", tt.name, tt.board)
		}
	}
}
