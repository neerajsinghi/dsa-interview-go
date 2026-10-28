package dp2d

import (
	"math/rand"
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	tests := []struct {
		s, p string
		want bool
	}{
		{"aa", "a", false},
		{"aa", "a*", true},
		{"ab", ".*", true},
		{"aab", "c*a*b", true},
		{"mississippi", "mis*is*p*.", false},
		{"mississippi", "mis*is*ip*.", true},
		{"", "", true},
		{"", "a*b*", true},
		{"", "a", false},
		{"a", "", false},
		{"ab", ".*c", false},
		{"aaa", "a*a", true},
	}
	impls := map[string]func(string, string) bool{"brute": RegexBrute, "memo": RegexMemo, "tab": RegexTab}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.s, tt.p); got != tt.want {
				t.Errorf("%s(%q,%q) = %v, want %v", name, tt.s, tt.p, got, tt.want)
			}
		}
	}
}

// randPattern builds a valid pattern: each element is a char or '.', maybe followed by '*'.
func randPattern(r *rand.Rand, elems int) string {
	b := []byte{}
	for i := 0; i < elems; i++ {
		b = append(b, "ab."[r.Intn(3)])
		if r.Intn(2) == 0 {
			b = append(b, '*')
		}
	}
	return string(b)
}

// TestRegexAgainstStdlib cross-checks with Go's regexp, anchored with ^...$.
func TestRegexAgainstStdlib(t *testing.T) {
	r := rand.New(rand.NewSource(8))
	for iter := 0; iter < 500; iter++ {
		s, p := randString(r, r.Intn(7), "ab"), randPattern(r, r.Intn(5))
		want := regexp.MustCompile("^(?:" + p + ")$").MatchString(s)
		if a, b, c := RegexBrute(s, p), RegexMemo(s, p), RegexTab(s, p); a != want || b != want || c != want {
			t.Fatalf("(%q,%q): stdlib=%v brute=%v memo=%v tab=%v", s, p, want, a, b, c)
		}
	}
}
