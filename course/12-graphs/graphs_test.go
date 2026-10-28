package graphs

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestBFSAndDFS(t *testing.T) {
	//   0 - 1 - 3
	//   |   |
	//   2 --+     4 (isolated)
	adj := AdjList(5, [][]int{{0, 1}, {0, 2}, {1, 2}, {1, 3}}, false)
	order, dist := BFS(adj, 0)
	if want := []int{0, 1, 2, 3}; !reflect.DeepEqual(order, want) {
		t.Errorf("BFS order = %v, want %v", order, want)
	}
	if want := []int{0, 1, 1, 2, -1}; !reflect.DeepEqual(dist, want) {
		t.Errorf("BFS dist = %v, want %v", dist, want)
	}
	if got, want := DFS(adj, 0), []int{0, 1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("DFS = %v, want %v", got, want)
	}
}

func TestDFSIterativeMatchesRecursive(t *testing.T) {
	r := rand.New(rand.NewSource(21))
	for iter := 0; iter < 200; iter++ {
		n := 1 + r.Intn(8)
		var edges [][]int
		for k := r.Intn(12); k > 0; k-- {
			edges = append(edges, []int{r.Intn(n), r.Intn(n)})
		}
		adj := AdjList(n, edges, r.Intn(2) == 0)
		if a, b := DFS(adj, 0), DFSIterative(adj, 0); !reflect.DeepEqual(a, b) {
			t.Fatalf("edges %v: recursive %v iterative %v", edges, a, b)
		}
	}
}

// isTopoOrder checks that order is a permutation of 0..n-1 respecting edges.
func isTopoOrder(adj [][]int, order []int) bool {
	if len(order) != len(adj) {
		return false
	}
	pos := make([]int, len(adj))
	for i := range pos {
		pos[i] = -1
	}
	for i, v := range order {
		if pos[v] != -1 {
			return false
		}
		pos[v] = i
	}
	for u := range adj {
		for _, v := range adj[u] {
			if pos[u] > pos[v] {
				return false
			}
		}
	}
	return true
}

func TestTopoSort(t *testing.T) {
	r := rand.New(rand.NewSource(22))
	for iter := 0; iter < 300; iter++ {
		n := 1 + r.Intn(7)
		var edges [][]int
		for k := r.Intn(10); k > 0; k-- {
			edges = append(edges, []int{r.Intn(n), r.Intn(n)})
		}
		adj := AdjList(n, edges, true)
		o1, ok1 := TopoSortKahn(adj)
		o2, ok2 := TopoSortDFS(adj)
		if ok1 != ok2 {
			t.Fatalf("edges %v: Kahn ok=%v DFS ok=%v", edges, ok1, ok2)
		}
		if ok1 && (!isTopoOrder(adj, o1) || !isTopoOrder(adj, o2)) {
			t.Fatalf("edges %v: invalid orders %v %v", edges, o1, o2)
		}
	}
}

func TestUnionFind(t *testing.T) {
	uf := NewUnionFind(5)
	if !uf.Union(0, 1) || !uf.Union(3, 4) {
		t.Fatal("first unions should succeed")
	}
	if uf.Union(1, 0) {
		t.Error("Union of an existing pair should return false")
	}
	if !uf.Connected(0, 1) || uf.Connected(1, 3) || uf.Count() != 3 {
		t.Errorf("unexpected state: count=%d", uf.Count())
	}
	uf.Union(1, 4)
	if !uf.Connected(0, 3) || uf.Count() != 2 {
		t.Errorf("after merge: count=%d", uf.Count())
	}
}
