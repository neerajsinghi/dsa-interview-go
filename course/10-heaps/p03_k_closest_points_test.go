package heaps

import (
	"math/rand"
	"slices"
	"testing"
)

// sortedDists returns the multiset of squared distances, sorted, so answers
// that differ only in tie-breaking or order compare equal.
func sortedDists(points [][]int) []int {
	d := make([]int, len(points))
	for i, p := range points {
		d[i] = dist2(p)
	}
	slices.Sort(d)
	return d
}

func TestKClosest(t *testing.T) {
	tests := []struct {
		points [][]int
		k      int
		want   [][]int
	}{
		{[][]int{{1, 3}, {-2, 2}}, 1, [][]int{{-2, 2}}},
		{[][]int{{3, 3}, {5, -1}, {-2, 4}}, 2, [][]int{{3, 3}, {-2, 4}}},
		{[][]int{{0, 0}}, 1, [][]int{{0, 0}}},
		{[][]int{{1, 1}, {2, 2}, {3, 3}}, 3, [][]int{{1, 1}, {2, 2}, {3, 3}}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int, int) [][]int{"sort": KClosestSort, "heap": KClosest} {
			if got := fn(tt.points, tt.k); !slices.Equal(sortedDists(got), sortedDists(tt.want)) {
				t.Errorf("%s(%v, %d) = %v, want %v", name, tt.points, tt.k, got, tt.want)
			}
		}
	}
}

func TestKClosestRandom(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for iter := 0; iter < 300; iter++ {
		points := make([][]int, 1+r.Intn(30))
		for i := range points {
			points[i] = []int{r.Intn(11) - 5, r.Intn(11) - 5}
		}
		k := 1 + r.Intn(len(points))
		got := KClosest(points, k)
		if len(got) != k {
			t.Fatalf("KClosest returned %d points, want %d", len(got), k)
		}
		if want := KClosestSort(points, k); !slices.Equal(sortedDists(got), sortedDists(want)) {
			t.Fatalf("KClosest(%v, %d) = %v, want distances of %v", points, k, got, want)
		}
	}
}
