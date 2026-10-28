package dp2d

import (
	"math/rand"
	"testing"
)

func TestIsInterleave(t *testing.T) {
	tests := []struct {
		s1, s2, s3 string
		want       bool
	}{
		{"aabcc", "dbbca", "aadbbcbcac", true},
		{"aabcc", "dbbca", "aadbbbaccc", false},
		{"", "", "", true},
		{"", "", "a", false},
		{"a", "", "a", true},
		{"", "b", "b", true},
		{"ab", "cd", "acbd", true},
		{"ab", "cd", "abdc", false},
		{"aa", "ab", "aaba", true}, // greedy-from-s1 would fail here
	}
	impls := map[string]func(string, string, string) bool{"memo": IsInterleaveMemo, "tab": IsInterleaveTab, "row": IsInterleaveRow}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.s1, tt.s2, tt.s3); got != tt.want {
				t.Errorf("%s(%q,%q,%q) = %v, want %v", name, tt.s1, tt.s2, tt.s3, got, tt.want)
			}
		}
	}
}

func TestIsInterleaveRandom(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for iter := 0; iter < 300; iter++ {
		s1, s2 := randString(r, r.Intn(6), "ab"), randString(r, r.Intn(6), "ab")
		// Half the time build a genuine interleaving, otherwise random noise.
		var s3 string
		if r.Intn(2) == 0 {
			i, j, b := 0, 0, []byte{}
			for i < len(s1) || j < len(s2) {
				if j == len(s2) || (i < len(s1) && r.Intn(2) == 0) {
					b, i = append(b, s1[i]), i+1
				} else {
					b, j = append(b, s2[j]), j+1
				}
			}
			s3 = string(b)
		} else {
			s3 = randString(r, len(s1)+len(s2), "ab")
		}
		a, b, c := IsInterleaveMemo(s1, s2, s3), IsInterleaveTab(s1, s2, s3), IsInterleaveRow(s1, s2, s3)
		if a != b || b != c {
			t.Fatalf("(%q,%q,%q): memo=%v tab=%v row=%v", s1, s2, s3, a, b, c)
		}
	}
}
