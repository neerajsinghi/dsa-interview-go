package greedy

import (
	"math/rand"
	"slices"
	"testing"
)

func iv(pairs ...int) []Interval {
	var out []Interval
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Interval{pairs[i], pairs[i+1]})
	}
	return out
}

func randIntervals(r *rand.Rand, n, span int) []Interval {
	out := make([]Interval, n)
	for i := range out {
		s := r.Intn(span)
		out[i] = Interval{s, s + r.Intn(span/3+1)}
	}
	return out
}

func TestMergeIntervals(t *testing.T) {
	tests := []struct {
		name    string
		in, out []Interval
	}{
		{"classic", iv(1, 3, 2, 6, 8, 10, 15, 18), iv(1, 6, 8, 10, 15, 18)},
		{"touching", iv(1, 4, 4, 5), iv(1, 5)},
		{"unsorted", iv(8, 10, 1, 3, 2, 6), iv(1, 6, 8, 10)},
		{"contained", iv(1, 10, 2, 3, 4, 5), iv(1, 10)},
		{"single", iv(5, 7), iv(5, 7)},
		{"empty", nil, nil},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		if got := MergeIntervals(tt.in); !slices.Equal(got, tt.out) {
			t.Errorf("%s: MergeIntervals = %v, want %v", tt.name, got, tt.out)
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("%s: input was mutated", tt.name)
		}
		if got := MergeIntervalsBrute(tt.in); !slices.Equal(got, tt.out) {
			t.Errorf("%s: brute = %v, want %v", tt.name, got, tt.out)
		}
	}
}

func TestMergeIntervalsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for iter := 0; iter < 300; iter++ {
		in := randIntervals(r, 1+r.Intn(10), 30)
		if a, b := MergeIntervals(in), MergeIntervalsBrute(in); !slices.Equal(a, b) {
			t.Fatalf("%v: sweep=%v brute=%v", in, a, b)
		}
	}
}
