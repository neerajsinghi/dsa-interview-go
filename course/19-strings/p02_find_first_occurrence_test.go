package stringalgos

import (
	"math/rand"
	"strings"
	"testing"
)

func TestStrStr(t *testing.T) {
	tests := []struct {
		haystack, needle string
		want             int
	}{
		{"sadbutsad", "sad", 0},
		{"leetcode", "leeto", -1},
		{"hello", "ll", 2},
		{"", "", 0},
		{"abc", "", 0},
		{"", "a", -1},
		{"a", "abc", -1},
		{"aaaaab", "aab", 3},
		{"mississippi", "issip", 4},
		{"abababca", "ababca", 2},
	}
	impls := map[string]func(string, string) int{"brute": StrStrBrute, "kmp": StrStrKMP, "z": StrStrZ, "rabinkarp": StrStrRabinKarp}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.haystack, tt.needle); got != tt.want {
				t.Errorf("%s(%q,%q) = %d, want %d", name, tt.haystack, tt.needle, got, tt.want)
			}
		}
	}
}

func TestStrStrRandom(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	for iter := 0; iter < 1000; iter++ {
		h, n := randText(r, r.Intn(30), "ab"), randText(r, r.Intn(5), "ab")
		want := strings.Index(h, n)
		for name, fn := range map[string]func(string, string) int{"brute": StrStrBrute, "kmp": StrStrKMP, "z": StrStrZ, "rabinkarp": StrStrRabinKarp} {
			if got := fn(h, n); got != want {
				t.Fatalf("%s(%q,%q) = %d, want %d", name, h, n, got, want)
			}
		}
	}
}
