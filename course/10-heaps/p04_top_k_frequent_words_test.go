package heaps

import (
	"math/rand"
	"slices"
	"testing"
)

func TestTopKFrequent(t *testing.T) {
	tests := []struct {
		words []string
		k     int
		want  []string
	}{
		{[]string{"i", "love", "go", "i", "love", "code"}, 2, []string{"i", "love"}},
		{[]string{"the", "day", "is", "sunny", "the", "the", "the", "sunny", "is", "is"}, 4,
			[]string{"the", "is", "sunny", "day"}},
		{[]string{"b", "a", "c"}, 3, []string{"a", "b", "c"}},
		{[]string{"x"}, 1, []string{"x"}},
		{[]string{"b", "b", "a", "a", "c"}, 1, []string{"a"}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]string, int) []string{"sort": TopKFrequentSort, "heap": TopKFrequent} {
			if got := fn(tt.words, tt.k); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v, %d) = %v, want %v", name, tt.words, tt.k, got, tt.want)
			}
		}
	}
}

func TestTopKFrequentRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	vocab := []string{"ant", "bee", "cat", "dog", "eel", "fox", "gnu", "hen"}
	for iter := 0; iter < 300; iter++ {
		words := make([]string, 1+r.Intn(40))
		for i := range words {
			words[i] = vocab[r.Intn(len(vocab))]
		}
		k := 1 + r.Intn(len(countWords(words)))
		if got, want := TopKFrequent(words, k), TopKFrequentSort(words, k); !slices.Equal(got, want) {
			t.Fatalf("TopKFrequent(%v, %d) = %v, want %v", words, k, got, want)
		}
	}
}
