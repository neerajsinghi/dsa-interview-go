package greedy

import (
	"math/rand"
	"slices"
	"testing"
)

func TestPartitionLabels(t *testing.T) {
	tests := []struct {
		s    string
		want []int
	}{
		{"ababcbacadefegdehijhklij", []int{9, 7, 8}},
		{"eccbbbbdec", []int{10}},
		{"abc", []int{1, 1, 1}},
		{"a", []int{1}},
		{"", nil},
		{"abca", []int{4}},
	}
	for _, tt := range tests {
		if got := PartitionLabels(tt.s); !slices.Equal(got, tt.want) {
			t.Errorf("PartitionLabels(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestPartitionLabelsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	for iter := 0; iter < 200; iter++ {
		b := make([]byte, 1+r.Intn(11))
		for i := range b {
			b[i] = "abcde"[r.Intn(5)]
		}
		s := string(b)
		if a, br := PartitionLabels(s), PartitionLabelsBrute(s); !slices.Equal(a, br) {
			t.Fatalf("%q: greedy=%v brute=%v", s, a, br)
		}
	}
}
