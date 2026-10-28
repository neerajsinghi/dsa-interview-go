package greedy

import (
	"math/rand"
	"slices"
	"testing"
)

func TestInsertInterval(t *testing.T) {
	tests := []struct {
		name string
		in   []Interval
		add  Interval
		want []Interval
	}{
		{"overlap one", iv(1, 3, 6, 9), Interval{2, 5}, iv(1, 5, 6, 9)},
		{"overlap many", iv(1, 2, 3, 5, 6, 7, 8, 10, 12, 16), Interval{4, 8}, iv(1, 2, 3, 10, 12, 16)},
		{"empty", nil, Interval{5, 7}, iv(5, 7)},
		{"before all", iv(3, 4, 6, 7), Interval{0, 1}, iv(0, 1, 3, 4, 6, 7)},
		{"after all", iv(1, 2), Interval{5, 6}, iv(1, 2, 5, 6)},
		{"in a gap", iv(1, 2, 8, 9), Interval{4, 5}, iv(1, 2, 4, 5, 8, 9)},
		{"covers all", iv(2, 3, 5, 6), Interval{1, 10}, iv(1, 10)},
		{"touching", iv(1, 2, 5, 6), Interval{2, 5}, iv(1, 6)},
	}
	for _, tt := range tests {
		if got := InsertInterval(tt.in, tt.add); !slices.Equal(got, tt.want) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

// Inserting must equal merging the union from scratch.
func TestInsertIntervalRandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	for iter := 0; iter < 300; iter++ {
		base := MergeIntervals(randIntervals(r, r.Intn(8), 40))
		add := randIntervals(r, 1, 40)[0]
		want := MergeIntervals(append(slices.Clone(base), add))
		if got := InsertInterval(base, add); !slices.Equal(got, want) {
			t.Fatalf("base=%v add=%v: got %v, want %v", base, add, got, want)
		}
	}
}
