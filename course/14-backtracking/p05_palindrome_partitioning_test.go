package backtracking

import (
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// normalizeStrings sorts partitions by their joined form for comparison.
func normalizeStrings(parts [][]string) []string {
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = strings.Join(p, "|")
	}
	sort.Strings(out)
	return out
}

func TestPartition(t *testing.T) {
	tests := []struct {
		s    string
		want [][]string
	}{
		{"aab", [][]string{{"a", "a", "b"}, {"aa", "b"}}},
		{"a", [][]string{{"a"}}},
		{"aba", [][]string{{"a", "b", "a"}, {"aba"}}},
		{"ab", [][]string{{"a", "b"}}},
		{"aaa", [][]string{{"a", "a", "a"}, {"a", "aa"}, {"aa", "a"}, {"aaa"}}},
		{"", [][]string{{}}}, // one way to split nothing
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(string) [][]string{"backtrack": Partition, "brute": PartitionBrute} {
			if got := normalizeStrings(fn(tt.s)); !reflect.DeepEqual(got, normalizeStrings(tt.want)) {
				t.Errorf("%s(%q) = %q, want %q", name, tt.s, got, normalizeStrings(tt.want))
			}
		}
	}
}

func TestPartitionRandom(t *testing.T) {
	r := rand.New(rand.NewSource(65))
	for iter := 0; iter < 200; iter++ {
		b := make([]byte, 1+r.Intn(9))
		for i := range b {
			b[i] = byte('a' + r.Intn(2))
		}
		s := string(b)
		if a, c := normalizeStrings(Partition(s)), normalizeStrings(PartitionBrute(s)); !reflect.DeepEqual(a, c) {
			t.Fatalf("%q: backtrack=%q brute=%q", s, a, c)
		}
	}
}
