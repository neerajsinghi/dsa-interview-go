package shortestpaths

import (
	"math/rand"
	"testing"
)

func TestMinCostConnectPoints(t *testing.T) {
	tests := []struct {
		name   string
		points [][]int
		want   int
	}{
		{"classic", [][]int{{0, 0}, {2, 2}, {3, 10}, {5, 2}, {7, 0}}, 20},
		{"three points", [][]int{{3, 12}, {-2, 5}, {-4, 1}}, 18},
		{"single point", [][]int{{0, 0}}, 0},
		{"empty", nil, 0},
		{"duplicates cost nothing", [][]int{{1, 1}, {1, 1}, {2, 1}}, 1},
		{"far apart", [][]int{{-1000000, -1000000}, {1000000, 1000000}}, 4000000},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int) int{"prim": MinCostConnectPoints, "kruskal": MinCostConnectPointsKruskal} {
			if got := fn(tt.points); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestMinCostConnectPointsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(46))
	for iter := 0; iter < 200; iter++ {
		points := make([][]int, 1+r.Intn(12))
		for i := range points {
			points[i] = []int{r.Intn(21) - 10, r.Intn(21) - 10}
		}
		if a, b := MinCostConnectPoints(points), MinCostConnectPointsKruskal(points); a != b {
			t.Fatalf("%v: prim=%d kruskal=%d", points, a, b)
		}
	}
}
