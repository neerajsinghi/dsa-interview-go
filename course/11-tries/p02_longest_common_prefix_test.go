package tries

import (
	"math/rand"
	"strings"
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want string
	}{
		{"classic", []string{"flower", "flow", "flight"}, "fl"},
		{"nothing shared", []string{"dog", "racecar", "car"}, ""},
		{"empty list", nil, ""},
		{"single word", []string{"alone"}, "alone"},
		{"one word is the prefix", []string{"interview", "inter", "internal"}, "inter"},
		{"contains empty string", []string{"abc", ""}, ""},
		{"all identical", []string{"go", "go", "go"}, "go"},
	}
	impls := map[string]func([]string) string{
		"scan": LongestCommonPrefixScan, "trie": LongestCommonPrefixTrie,
	}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.strs); got != tt.want {
				t.Errorf("%s/%s = %q, want %q", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestLongestCommonPrefixRandom(t *testing.T) {
	r := rand.New(rand.NewSource(12))
	for iter := 0; iter < 500; iter++ {
		base := strings.Repeat("a", r.Intn(4))
		strs := make([]string, 1+r.Intn(4))
		for i := range strs {
			b := []byte(base)
			for k := r.Intn(4); k > 0; k-- {
				b = append(b, byte('a'+r.Intn(2)))
			}
			strs[i] = string(b)
		}
		want := LongestCommonPrefixScan(strs)
		if got := LongestCommonPrefixTrie(strs); got != want {
			t.Fatalf("%q: trie=%q scan=%q", strs, got, want)
		}
		for _, s := range strs {
			if !strings.HasPrefix(s, want) {
				t.Fatalf("%q: %q is not a prefix of %q", strs, want, s)
			}
		}
	}
}
