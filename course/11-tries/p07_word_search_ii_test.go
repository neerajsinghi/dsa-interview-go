package tries

import (
	"math/rand"
	"reflect"
	"testing"
)

func gridOf(rows ...string) [][]byte {
	g := make([][]byte, len(rows))
	for i, r := range rows {
		g[i] = []byte(r)
	}
	return g
}

func TestFindWords(t *testing.T) {
	tests := []struct {
		name  string
		board [][]byte
		words []string
		want  []string
	}{
		{"classic", gridOf("oaan", "etae", "ihkr", "iflv"), []string{"oath", "pea", "eat", "rain"}, []string{"eat", "oath"}},
		{"cannot reuse a cell", gridOf("ab", "cd"), []string{"abcb"}, nil},
		{"single cell", gridOf("a"), []string{"a", "aa"}, []string{"a"}},
		{"duplicate words", gridOf("ab"), []string{"ab", "ab", "ba"}, []string{"ab", "ba"}},
		{"shared prefixes", gridOf("abc", "aed", "afg"), []string{"abcdefg", "gfedcbea", "eaabcdgfa", "befa", "dgc", "ade"}, []string{"abcdefg", "befa", "eaabcdgfa"}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]byte, []string) []string{"trie": FindWords, "brute": FindWordsBrute} {
			if got := fn(tt.board, tt.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s/%s = %q, want %q", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestFindWordsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(17))
	for iter := 0; iter < 200; iter++ {
		rows, cols := 1+r.Intn(3), 1+r.Intn(3)
		board := make([][]byte, rows)
		for i := range board {
			board[i] = make([]byte, cols)
			for j := range board[i] {
				board[i][j] = byte('a' + r.Intn(3))
			}
		}
		words := make([]string, 1+r.Intn(6))
		for i := range words {
			b := make([]byte, 1+r.Intn(4))
			for j := range b {
				b[j] = byte('a' + r.Intn(3))
			}
			words[i] = string(b)
		}
		want := FindWordsBrute(board, words)
		if got := FindWords(board, words); !reflect.DeepEqual(got, want) {
			t.Fatalf("board %q words %q: trie=%q brute=%q", board, words, got, want)
		}
	}
}
