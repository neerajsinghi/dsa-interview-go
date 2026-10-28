package stringalgos

import (
	"math/rand"
	"strings"
	"testing"
)

func TestCompress(t *testing.T) {
	tests := []struct{ in, want string }{
		{"aabbccc", "a2b2c3"},
		{"a", "a"},
		{"abbbbbbbbbbbb", "ab12"},
		{"", ""},
		{"abc", "abc"},
		{"aaabbaa", "a3b2a2"},
		{strings.Repeat("z", 100), "z100"},
	}
	for _, tt := range tests {
		buf := []byte(tt.in)
		n := Compress(buf)
		if got := string(buf[:n]); got != tt.want {
			t.Errorf("Compress(%q) = %q, want %q", tt.in, got, tt.want)
		}
		if got := CompressNaive([]byte(tt.in)); got != tt.want {
			t.Errorf("CompressNaive(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCompressRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for iter := 0; iter < 500; iter++ {
		var sb strings.Builder
		for k := r.Intn(6); k > 0; k-- {
			sb.WriteString(strings.Repeat(string("abc"[r.Intn(3)]), 1+r.Intn(15)))
		}
		in := sb.String()
		buf := []byte(in)
		n := Compress(buf)
		if want := CompressNaive([]byte(in)); string(buf[:n]) != want {
			t.Fatalf("%q: inplace=%q naive=%q", in, buf[:n], want)
		}
	}
}
