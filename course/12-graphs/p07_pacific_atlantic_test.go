package graphs

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestPacificAtlantic(t *testing.T) {
	tests := []struct {
		name    string
		heights [][]int
		want    [][]int
	}{
		{
			"classic",
			[][]int{{1, 2, 2, 3, 5}, {3, 2, 3, 4, 4}, {2, 4, 5, 3, 1}, {6, 7, 1, 4, 5}, {5, 1, 1, 2, 4}},
			[][]int{{0, 4}, {1, 3}, {1, 4}, {2, 2}, {3, 0}, {3, 1}, {4, 0}},
		},
		{"single cell", [][]int{{7}}, [][]int{{0, 0}}},
		{"flat", [][]int{{1, 1}, {1, 1}}, [][]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}},
		{"empty", nil, nil},
		{"valley", [][]int{{3, 3, 3}, {3, 1, 3}, {3, 3, 3}}, [][]int{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int) [][]int{"reverse-bfs": PacificAtlantic, "brute": PacificAtlanticBrute} {
			if got := fn(tt.heights); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s/%s = %v, want %v", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestPacificAtlanticRandom(t *testing.T) {
	r := rand.New(rand.NewSource(29))
	for iter := 0; iter < 200; iter++ {
		m, n := 1+r.Intn(5), 1+r.Intn(5)
		h := make([][]int, m)
		for i := range h {
			h[i] = make([]int, n)
			for j := range h[i] {
				h[i][j] = r.Intn(4)
			}
		}
		if a, b := PacificAtlantic(h), PacificAtlanticBrute(h); !reflect.DeepEqual(a, b) {
			t.Fatalf("%v: fast=%v brute=%v", h, a, b)
		}
	}
}
