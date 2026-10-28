package dp2d

import (
	"math/rand"
	"testing"
)

func TestLPS(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"bbbab", 4}, {"cbbd", 2}, {"", 0}, {"a", 1}, {"abc", 1},
		{"racecar", 7}, {"character", 5}, {"aaaa", 4},
	}
	impls := map[string]func(string) int{"memo": LPSMemo, "tab": LPSTab, "row": LPSRow}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.s); got != tt.want {
				t.Errorf("%s(%q) = %d, want %d", name, tt.s, got, tt.want)
			}
		}
	}
}

func TestLPSMatchesLCSOfReverse(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	for iter := 0; iter < 300; iter++ {
		s := randString(r, r.Intn(12), "abc")
		rev := []byte(s)
		for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
			rev[i], rev[j] = rev[j], rev[i]
		}
		want := LCSTab(s, string(rev))
		if a, b, c := LPSMemo(s), LPSTab(s), LPSRow(s); a != want || b != want || c != want {
			t.Fatalf("%q: lcs(s,rev)=%d memo=%d tab=%d row=%d", s, want, a, b, c)
		}
	}
}
