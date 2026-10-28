package arrayshashing

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestGroupAnagrams(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  [][]string
	}{
		{"empty", []string{}, [][]string{}},
		{"single empty string", []string{""}, [][]string{{""}}},
		{"single", []string{"a"}, [][]string{{"a"}}},
		{"classic", []string{"eat", "tea", "tan", "ate", "nat", "bat"},
			[][]string{{"eat", "tea", "ate"}, {"tan", "nat"}, {"bat"}}},
		{"duplicates", []string{"ab", "ba", "ab"}, [][]string{{"ab", "ba", "ab"}}},
		{"same letters different counts", []string{"aab", "abb"}, [][]string{{"aab"}, {"abb"}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := GroupAnagramsSorted(tc.words); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("sorted-key = %v, want %v", got, tc.want)
			}
			if got := GroupAnagrams(tc.words); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("count-key = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGroupAnagramsRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(4))
	for range 300 {
		words := make([]string, rng.Intn(12))
		for i := range words {
			b := make([]byte, rng.Intn(4))
			for j := range b {
				b[j] = byte('a' + rng.Intn(3))
			}
			words[i] = string(b)
		}
		a, b := GroupAnagramsSorted(words), GroupAnagrams(words)
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("%v: sorted=%v count=%v", words, a, b)
		}
	}
}
