package stringalgos

import (
	"math/rand"
	"testing"
)

func TestReverseWords(t *testing.T) {
	tests := []struct{ in, want string }{
		{"the sky is blue", "blue is sky the"},
		{"  hello world  ", "world hello"},
		{"a good   example", "example good a"},
		{"", ""},
		{"     ", ""},
		{"single", "single"},
	}
	for _, tt := range tests {
		if got := ReverseWords(tt.in); got != tt.want {
			t.Errorf("ReverseWords(%q) = %q, want %q", tt.in, got, tt.want)
		}
		if got := string(ReverseWordsInPlace([]byte(tt.in))); got != tt.want {
			t.Errorf("ReverseWordsInPlace(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestReverseWordsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for iter := 0; iter < 500; iter++ {
		s := randText(r, r.Intn(20), "ab  ")
		if a, b := ReverseWords(s), string(ReverseWordsInPlace([]byte(s))); a != b {
			t.Fatalf("%q: fields=%q inplace=%q", s, a, b)
		}
	}
}
