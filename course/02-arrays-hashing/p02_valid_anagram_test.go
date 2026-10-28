package arrayshashing

import (
	"math/rand"
	"testing"
)

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		s, t string
		want bool
	}{
		{"", "", true},
		{"a", "a", true},
		{"a", "b", false},
		{"listen", "silent", true},
		{"rat", "car", false},
		{"aacc", "ccac", false}, // same letters, different counts
		{"ab", "abc", false},
		{"aabbcc", "cbacba", true},
	}
	for _, tc := range tests {
		for name, fn := range map[string]func(string, string) bool{
			"Sort": IsAnagramSort, "Count": IsAnagram, "Unicode": IsAnagramUnicode,
		} {
			if got := fn(tc.s, tc.t); got != tc.want {
				t.Errorf("%s(%q, %q) = %v, want %v", name, tc.s, tc.t, got, tc.want)
			}
		}
	}
}

func TestIsAnagramUnicode(t *testing.T) {
	tests := []struct {
		s, t string
		want bool
	}{
		{"héllo", "olléh", true},
		{"世界", "界世", true},
		{"é", "e", false},
		{"ab", "a", false},
	}
	for _, tc := range tests {
		if got := IsAnagramUnicode(tc.s, tc.t); got != tc.want {
			t.Errorf("IsAnagramUnicode(%q, %q) = %v, want %v", tc.s, tc.t, got, tc.want)
		}
	}
}

func TestIsAnagramRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	word := func() string {
		b := make([]byte, rng.Intn(6))
		for i := range b {
			b[i] = byte('a' + rng.Intn(3))
		}
		return string(b)
	}
	for range 1000 {
		s, u := word(), word()
		want := IsAnagramSort(s, u)
		if IsAnagram(s, u) != want || IsAnagramUnicode(s, u) != want {
			t.Fatalf("disagreement on %q, %q", s, u)
		}
	}
}
