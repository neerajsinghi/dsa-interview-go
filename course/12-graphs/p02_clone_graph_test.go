package graphs

import (
	"math/rand"
	"testing"
)

// buildGraph creates nodes 1..len(adj) where adj[i] lists the (1-based)
// neighbours of node i+1, and returns node 1.
func buildGraph(adj [][]int) *GraphNode {
	if len(adj) == 0 {
		return nil
	}
	nodes := make([]*GraphNode, len(adj))
	for i := range nodes {
		nodes[i] = &GraphNode{Val: i + 1}
	}
	for i, nbs := range adj {
		for _, v := range nbs {
			nodes[i].Neighbors = append(nodes[i].Neighbors, nodes[v-1])
		}
	}
	return nodes[0]
}

// checkClone walks original and clone in lockstep, verifying equal values,
// equal neighbour lists and that no node pointer is shared.
func checkClone(t *testing.T, orig, cl *GraphNode) {
	t.Helper()
	if orig == nil || cl == nil {
		if orig != cl {
			t.Fatalf("nil mismatch: orig=%v clone=%v", orig, cl)
		}
		return
	}
	origSet := map[*GraphNode]bool{}
	pair := map[*GraphNode]*GraphNode{}
	var walk func(a, b *GraphNode)
	walk = func(a, b *GraphNode) {
		if a == b {
			t.Fatalf("node %d is shared, not copied", a.Val)
		}
		if prev, ok := pair[a]; ok {
			if prev != b {
				t.Fatalf("node %d maps to two different copies", a.Val)
			}
			return
		}
		pair[a] = b
		origSet[a] = true
		if a.Val != b.Val || len(a.Neighbors) != len(b.Neighbors) {
			t.Fatalf("node %d: clone has val %d and %d neighbours, want %d", a.Val, b.Val, len(b.Neighbors), len(a.Neighbors))
		}
		for i := range a.Neighbors {
			walk(a.Neighbors[i], b.Neighbors[i])
		}
	}
	walk(orig, cl)
	for _, c := range pair {
		if origSet[c] {
			t.Fatalf("clone reuses original node %d", c.Val)
		}
	}
}

func TestCloneGraph(t *testing.T) {
	tests := []struct {
		name string
		adj  [][]int
	}{
		{"four-cycle", [][]int{{2, 4}, {1, 3}, {2, 4}, {1, 3}}},
		{"single node", [][]int{{}}},
		{"empty", nil},
		{"pair", [][]int{{2}, {1}}},
		{"triangle", [][]int{{2, 3}, {1, 3}, {1, 2}}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(*GraphNode) *GraphNode{"dfs": CloneGraph, "bfs": CloneGraphBFS} {
			g := buildGraph(tt.adj)
			t.Run(tt.name+"/"+name, func(t *testing.T) { checkClone(t, g, fn(g)) })
		}
	}
}

func TestCloneGraphRandom(t *testing.T) {
	r := rand.New(rand.NewSource(24))
	for iter := 0; iter < 100; iter++ {
		n := 1 + r.Intn(8)
		adj := make([][]int, n)
		link := func(a, b int) { adj[a] = append(adj[a], b+1); adj[b] = append(adj[b], a+1) }
		for v := 1; v < n; v++ {
			link(v, r.Intn(v)) // random spanning tree keeps it connected
		}
		for k := r.Intn(n); k > 0; k-- {
			if a, b := r.Intn(n), r.Intn(n); a != b {
				link(a, b)
			}
		}
		g := buildGraph(adj)
		checkClone(t, g, CloneGraph(g))
		checkClone(t, g, CloneGraphBFS(g))
	}
}
