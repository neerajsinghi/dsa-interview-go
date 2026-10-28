package shortestpaths

import (
	"math/rand"
	"reflect"
	"testing"
)

func randomEdges(r *rand.Rand, n, m, minW, maxW int) []Edge {
	edges := make([]Edge, m)
	for i := range edges {
		edges[i] = Edge{r.Intn(n), r.Intn(n), minW + r.Intn(maxW-minW+1)}
	}
	return edges
}

func TestDijkstraSmall(t *testing.T) {
	//  0 --4--> 1
	//  |        ^
	//  1        2
	//  v        |
	//  2 --1--> 3      4 unreachable
	adj := BuildAdj(5, []Edge{{0, 1, 4}, {0, 2, 1}, {2, 3, 1}, {3, 1, 2}}, true)
	if got, want := Dijkstra(adj, 0), []int{0, 4, 1, 2, Inf}; !reflect.DeepEqual(got, want) {
		t.Errorf("Dijkstra = %v, want %v", got, want)
	}
}

func TestShortestPathAlgorithmsAgree(t *testing.T) {
	r := rand.New(rand.NewSource(41))
	for iter := 0; iter < 200; iter++ {
		n := 1 + r.Intn(7)
		edges := randomEdges(r, n, r.Intn(15), 0, 9)
		d1 := Dijkstra(BuildAdj(n, edges, true), 0)
		d2, ok := BellmanFord(n, edges, 0)
		d3 := FloydWarshall(n, edges)[0]
		if !ok || !reflect.DeepEqual(d1, d2) || !reflect.DeepEqual(d1, d3) {
			t.Fatalf("edges %v: dijkstra=%v bellman=%v(ok=%v) floyd=%v", edges, d1, d2, ok, d3)
		}
	}
}

func TestZeroOneBFS(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	for iter := 0; iter < 200; iter++ {
		n := 1 + r.Intn(8)
		adj := BuildAdj(n, randomEdges(r, n, r.Intn(16), 0, 1), r.Intn(2) == 0)
		if a, b := ZeroOneBFS(adj, 0), Dijkstra(adj, 0); !reflect.DeepEqual(a, b) {
			t.Fatalf("0-1 BFS %v vs Dijkstra %v", a, b)
		}
	}
}

func TestBellmanFordNegative(t *testing.T) {
	// Negative edge, no cycle: 0->1 (4), 0->2 (5), 2->1 (-3) gives dist[1] = 2.
	dist, ok := BellmanFord(3, []Edge{{0, 1, 4}, {0, 2, 5}, {2, 1, -3}}, 0)
	if !ok || !reflect.DeepEqual(dist, []int{0, 2, 5}) {
		t.Errorf("got %v ok=%v", dist, ok)
	}
	// Negative cycle 1->2->1 with total weight -1.
	if _, ok := BellmanFord(3, []Edge{{0, 1, 1}, {1, 2, -2}, {2, 1, 1}}, 0); ok {
		t.Error("negative cycle not detected")
	}
	// A negative cycle that src cannot reach does not matter.
	if _, ok := BellmanFord(3, []Edge{{1, 2, -2}, {2, 1, 1}}, 0); !ok {
		t.Error("unreachable negative cycle should be ignored")
	}
	if fw := FloydWarshall(3, []Edge{{0, 1, 1}, {1, 2, -2}, {2, 1, 1}}); fw[1][1] >= 0 {
		t.Error("Floyd-Warshall should expose the negative cycle on the diagonal")
	}
}

func TestMSTAlgorithmsAgree(t *testing.T) {
	// Edges 0-1 (1), 1-2 (2), 0-2 (3), 0-3 (4), 3-2 (5).
	// Kruskal keeps 1 and 2, skips 3 (0 and 2 already joined), keeps 4:
	// MST weight = 1 + 2 + 4 = 7.
	edges := []Edge{{0, 1, 1}, {1, 2, 2}, {0, 2, 3}, {0, 3, 4}, {3, 2, 5}}
	total, mst, ok := Kruskal(4, edges)
	if total != 7 || len(mst) != 3 || !ok {
		t.Errorf("Kruskal = %d %v %v, want 7", total, mst, ok)
	}
	if total, ok := Prim(BuildAdj(4, edges, false)); total != 7 || !ok {
		t.Errorf("Prim = %d %v, want 7", total, ok)
	}
	if _, _, ok := Kruskal(3, []Edge{{0, 1, 1}}); ok {
		t.Error("Kruskal should report a disconnected graph")
	}
	if _, ok := Prim(BuildAdj(3, []Edge{{0, 1, 1}}, false)); ok {
		t.Error("Prim should report a disconnected graph")
	}

	r := rand.New(rand.NewSource(43))
	for iter := 0; iter < 200; iter++ {
		n := 1 + r.Intn(8)
		edges := randomEdges(r, n, r.Intn(20), 1, 9)
		k, _, ok1 := Kruskal(n, edges)
		p, ok2 := Prim(BuildAdj(n, edges, false))
		if ok1 != ok2 || (ok1 && k != p) {
			t.Fatalf("edges %v: kruskal=%d(%v) prim=%d(%v)", edges, k, ok1, p, ok2)
		}
	}
}
