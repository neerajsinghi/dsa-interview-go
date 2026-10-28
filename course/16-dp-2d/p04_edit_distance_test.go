package dp2d

import (
	"math/rand"
	"testing"
)

func TestEditDistance(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"horse", "ros", 3},
		{"intention", "execution", 5},
		{"", "", 0},
		{"", "abc", 3},
		{"abc", "", 3},
		{"same", "same", 0},
		{"kitten", "sitting", 3},
		{"a", "b", 1},
	}
	impls := map[string]func(string, string) int{"memo": EditDistanceMemo, "tab": EditDistanceTab, "row": EditDistanceRow}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.a, tt.b); got != tt.want {
				t.Errorf("%s(%q,%q) = %d, want %d", name, tt.a, tt.b, got, tt.want)
			}
		}
	}
}

func TestEditDistanceAgree(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for iter := 0; iter < 300; iter++ {
		a, b := randString(r, r.Intn(10), "abcd"), randString(r, r.Intn(10), "abcd")
		x, y, z := EditDistanceMemo(a, b), EditDistanceTab(a, b), EditDistanceRow(a, b)
		if x != y || y != z {
			t.Fatalf("(%q,%q): memo=%d tab=%d row=%d", a, b, x, y, z)
		}
		// Sanity bounds: |m-n| <= dist <= max(m,n).
		lo, hi := len(a)-len(b), max(len(a), len(b))
		if lo < 0 {
			lo = -lo
		}
		if x < lo || x > hi {
			t.Fatalf("(%q,%q): dist %d outside [%d,%d]", a, b, x, lo, hi)
		}
	}
}
