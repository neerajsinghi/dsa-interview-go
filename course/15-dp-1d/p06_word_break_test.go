package dp1d

import (
	"math/rand"
	"strings"
	"testing"
)

var wordBreakImpls = map[string]func(string, []string) bool{
	"recursive": WordBreakRecursive, "memo": WordBreakMemo, "tab": WordBreak,
}

func TestWordBreak(t *testing.T) {
	tests := []struct {
		name string
		s    string
		dict []string
		want bool
	}{
		{"two words", "leetcode", []string{"leet", "code"}, true},
		{"reuse", "applepenapple", []string{"apple", "pen"}, true},
		{"dead end", "catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},
		{"empty string", "", []string{"a"}, true},
		{"empty dict", "a", nil, false},
		{"needs backtracking", "aaaab", []string{"a", "aa", "aaa", "aaaa"}, false},
		{"choose the shorter word", "cars", []string{"car", "ca", "rs"}, true},
	}
	for _, tt := range tests {
		for name, fn := range wordBreakImpls {
			if got := fn(tt.s, tt.dict); got != tt.want {
				t.Errorf("%s/%s = %v, want %v", tt.name, name, got, tt.want)
			}
		}
	}
	// A classic memoization stress test: exponential without the cache.
	long := strings.Repeat("a", 60) + "b"
	dict := []string{"a", "aa", "aaa", "aaaa", "aaaaa"}
	if WordBreakMemo(long, dict) || WordBreak(long, dict) {
		t.Error("long 'aaa...ab' should not be breakable")
	}
}

func TestWordBreakRandom(t *testing.T) {
	r := rand.New(rand.NewSource(75))
	word := func(maxLen int) string {
		b := make([]byte, 1+r.Intn(maxLen))
		for i := range b {
			b[i] = byte('a' + r.Intn(2))
		}
		return string(b)
	}
	for iter := 0; iter < 300; iter++ {
		dict := make([]string, r.Intn(5))
		for i := range dict {
			dict[i] = word(3)
		}
		s := ""
		for k := r.Intn(5); k > 0; k-- {
			s += word(3)
		}
		want := WordBreakRecursive(s, dict)
		for name, fn := range wordBreakImpls {
			if got := fn(s, dict); got != want {
				t.Fatalf("%q %q: %s=%v brute=%v", s, dict, name, got, want)
			}
		}
	}
}
