package dp2d

import (
	"math/rand"
	"testing"
)

func TestLCS(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"abcde", "ace", 3},
		{"abc", "abc", 3},
		{"abc", "def", 0},
		{"", "abc", 0},
		{"", "", 0},
		{"bsbininm", "jmjkbkjkv", 1},
		{"oxcpqrsvwf", "shmtulqrypy", 2},
	}
	impls := map[string]func(string, string) int{"brute": LCSBrute, "memo": LCSMemo, "tab": LCSTab, "tworows": LCSTwoRows}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.a, tt.b); got != tt.want {
				t.Errorf("%s(%q,%q) = %d, want %d", name, tt.a, tt.b, got, tt.want)
			}
		}
	}
}

// isSubseq reports whether s is a subsequence of t.
func isSubseq(s, t string) bool {
	i := 0
	for j := 0; i < len(s) && j < len(t); j++ {
		if s[i] == t[j] {
			i++
		}
	}
	return i == len(s)
}

func randString(r *rand.Rand, n int, alphabet string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[r.Intn(len(alphabet))]
	}
	return string(b)
}

func TestLCSRandom(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for iter := 0; iter < 300; iter++ {
		a, b := randString(r, r.Intn(9), "abc"), randString(r, r.Intn(9), "abc")
		want := LCSBrute(a, b)
		if m, tb, tr := LCSMemo(a, b), LCSTab(a, b), LCSTwoRows(a, b); m != want || tb != want || tr != want {
			t.Fatalf("(%q,%q): brute=%d memo=%d tab=%d tworows=%d", a, b, want, m, tb, tr)
		}
		s := LCSString(a, b)
		if len(s) != want || !isSubseq(s, a) || !isSubseq(s, b) {
			t.Fatalf("LCSString(%q,%q) = %q, want a common subsequence of length %d", a, b, s, want)
		}
	}
}
