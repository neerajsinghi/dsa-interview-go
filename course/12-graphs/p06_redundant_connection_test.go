package graphs

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestFindRedundantConnection(t *testing.T) {
	tests := []struct {
		name  string
		edges [][]int
		want  []int
	}{
		{"triangle", [][]int{{1, 2}, {1, 3}, {2, 3}}, []int{2, 3}},
		{"square with tail", [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 4}, {1, 5}}, []int{1, 4}},
		{"double edge", [][]int{{1, 2}, {2, 1}}, []int{2, 1}},
		{"cycle closes late", [][]int{{3, 4}, {1, 2}, {2, 4}, {3, 5}, {2, 5}}, []int{2, 5}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]int) []int{"unionfind": FindRedundantConnection, "brute": FindRedundantConnectionBrute} {
			if got := fn(tt.edges); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s/%s = %v, want %v", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestFindRedundantConnectionRandom(t *testing.T) {
	r := rand.New(rand.NewSource(28))
	for iter := 0; iter < 300; iter++ {
		n := 2 + r.Intn(8)
		has := map[[2]int]bool{}
		var edges [][]int
		add := func(a, b int) bool {
			if a > b {
				a, b = b, a
			}
			if a == b || has[[2]int{a, b}] {
				return false
			}
			has[[2]int{a, b}] = true
			edges = append(edges, []int{a, b})
			return true
		}
		for v := 2; v <= n; v++ {
			add(v, 1+r.Intn(v-1)) // random tree on 1..n
		}
		if n == 2 {
			edges = append(edges, []int{1, 2}) // the only possible extra edge
		} else {
			for !add(1+r.Intn(n), 1+r.Intn(n)) {
			}
		}
		r.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
		a, b := FindRedundantConnection(edges), FindRedundantConnectionBrute(edges)
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("%v: unionfind=%v brute=%v", edges, a, b)
		}
	}
}
