// Package shortestpaths holds weighted-graph templates — Dijkstra,
// Bellman-Ford, 0-1 BFS, Floyd-Warshall, Kruskal and Prim — and the
// interview problems built on them.
//
// Conventions: nodes are 0..n-1; Inf marks "unreachable".
package shortestpaths

import (
	"container/heap"
	"math"
	"sort"
)

// Inf is the distance to an unreachable node. It is half of MaxInt so that
// Inf + weight never overflows when relaxing an edge.
const Inf = math.MaxInt / 2

// Edge is a weighted edge, directed From -> To unless stated otherwise.
type Edge struct {
	From, To, Weight int
}

// Arc is one entry of a weighted adjacency list: the target and the weight.
type Arc struct {
	To, Weight int
}

// BuildAdj turns an edge list into a weighted adjacency list. For an
// undirected graph each edge is stored in both directions.
func BuildAdj(n int, edges []Edge, directed bool) [][]Arc {
	adj := make([][]Arc, n)
	for _, e := range edges {
		adj[e.From] = append(adj[e.From], Arc{e.To, e.Weight})
		if !directed {
			adj[e.To] = append(adj[e.To], Arc{e.From, e.Weight})
		}
	}
	return adj
}

// item is a heap entry: a node and the tentative distance it was pushed with.
type item struct {
	node, dist int
}

// minHeap implements heap.Interface ordered by dist.
type minHeap []item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// Dijkstra returns the shortest distance from src to every node in a graph
// with NON-NEGATIVE weights (Inf if unreachable).
//
// It uses "lazy deletion": instead of a decrease-key operation, a node may be
// pushed several times and stale entries are skipped when popped.
//
// Time: O((V + E) log V). Space: O(V + E).
func Dijkstra(adj [][]Arc, src int) []int {
	dist := make([]int, len(adj))
	for i := range dist {
		dist[i] = Inf
	}
	dist[src] = 0
	h := &minHeap{{src, 0}}
	for h.Len() > 0 {
		cur := heap.Pop(h).(item)
		if cur.dist > dist[cur.node] {
			continue // stale entry: a shorter path was already settled
		}
		for _, a := range adj[cur.node] {
			if nd := cur.dist + a.Weight; nd < dist[a.To] {
				dist[a.To] = nd
				heap.Push(h, item{a.To, nd})
			}
		}
	}
	return dist
}

// BellmanFord returns shortest distances from src in a directed graph that
// may have NEGATIVE weights. ok is false if a negative cycle is reachable
// from src (then "shortest" is undefined).
//
// Idea: a shortest path uses at most n-1 edges, so relaxing every edge n-1
// times is enough. If an n-th pass still improves something, a negative
// cycle exists.
//
// Time: O(V · E). Space: O(V).
func BellmanFord(n int, edges []Edge, src int) (dist []int, ok bool) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = Inf
	}
	dist[src] = 0
	for pass := 0; pass < n-1; pass++ {
		changed := false
		for _, e := range edges {
			if dist[e.From] != Inf && dist[e.From]+e.Weight < dist[e.To] {
				dist[e.To] = dist[e.From] + e.Weight
				changed = true
			}
		}
		if !changed {
			break // already converged: stop early
		}
	}
	for _, e := range edges {
		if dist[e.From] != Inf && dist[e.From]+e.Weight < dist[e.To] {
			return dist, false
		}
	}
	return dist, true
}

// ZeroOneBFS returns shortest distances when every weight is 0 or 1.
// A deque replaces Dijkstra's heap: 0-edges go to the front (same distance),
// 1-edges to the back (distance + 1), so the deque stays sorted.
//
// Time: O(V + E). Space: O(V).
func ZeroOneBFS(adj [][]Arc, src int) []int {
	dist := make([]int, len(adj))
	for i := range dist {
		dist[i] = Inf
	}
	dist[src] = 0
	var dq deque
	dq.pushBack(src)
	for dq.len() > 0 {
		u := dq.popFront()
		for _, a := range adj[u] {
			if nd := dist[u] + a.Weight; nd < dist[a.To] {
				dist[a.To] = nd
				if a.Weight == 0 {
					dq.pushFront(a.To)
				} else {
					dq.pushBack(a.To)
				}
			}
		}
	}
	return dist
}

// deque is a minimal double-ended queue of ints: front holds pushed-front
// items in reverse (top of stack = first), back is a queue read from head.
type deque struct {
	front []int
	back  []int
	head  int
}

func (d *deque) pushFront(x int) { d.front = append(d.front, x) }
func (d *deque) pushBack(x int)  { d.back = append(d.back, x) }
func (d *deque) len() int        { return len(d.front) + len(d.back) - d.head }
func (d *deque) popFront() int {
	if n := len(d.front); n > 0 {
		x := d.front[n-1]
		d.front = d.front[:n-1]
		return x
	}
	x := d.back[d.head]
	d.head++
	return x
}

// FloydWarshall returns the all-pairs shortest distance matrix of a directed
// graph (Inf if unreachable). Negative edges are fine; negative cycles show
// up as dist[i][i] < 0.
//
// The DP: after round k, dist[i][j] is the shortest path whose intermediate
// nodes are all in {0..k}. Node k is either unused or used once in the
// middle: dist[i][j] = min(dist[i][j], dist[i][k] + dist[k][j]).
// k MUST be the outer loop.
//
// Time: O(V³). Space: O(V²).
func FloydWarshall(n int, edges []Edge) [][]int {
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			if i != j {
				dist[i][j] = Inf
			}
		}
	}
	for _, e := range edges {
		dist[e.From][e.To] = min(dist[e.From][e.To], e.Weight) // keep the cheapest parallel edge
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if dist[i][k] == Inf {
				continue
			}
			for j := 0; j < n; j++ {
				if dist[k][j] != Inf && dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}
	return dist
}

// dsu is a compact union-find used by Kruskal (see module 12 for the full
// explanation).
type dsu []int

func newDSU(n int) dsu {
	d := make(dsu, n)
	for i := range d {
		d[i] = i
	}
	return d
}

func (d dsu) find(x int) int {
	for d[x] != x {
		d[x] = d[d[x]] // path halving
		x = d[x]
	}
	return x
}

func (d dsu) union(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	d[ra] = rb
	return true
}

// Kruskal builds a minimum spanning tree of an undirected graph: sort the
// edges by weight and keep each edge that joins two different components.
// ok is false if the graph is disconnected (then mst is a spanning forest).
//
// Time: O(E log E). Space: O(V + E).
func Kruskal(n int, edges []Edge) (total int, mst []Edge, ok bool) {
	sorted := append([]Edge(nil), edges...) // don't reorder the caller's slice
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Weight < sorted[j].Weight })
	d := newDSU(n)
	for _, e := range sorted {
		if d.union(e.From, e.To) {
			total += e.Weight
			mst = append(mst, e)
			if len(mst) == n-1 {
				break
			}
		}
	}
	return total, mst, len(mst) == n-1 || n == 0
}

// Prim grows a minimum spanning tree from node 0, always adding the cheapest
// edge that leaves the tree. ok is false if the graph is disconnected.
// adj must be undirected (both directions present).
//
// Time: O(E log V) with a heap. Space: O(V + E).
func Prim(adj [][]Arc) (total int, ok bool) {
	n := len(adj)
	if n == 0 {
		return 0, true
	}
	inTree := make([]bool, n)
	h := &minHeap{{0, 0}} // item.dist = weight of the edge that reaches node
	added := 0
	for h.Len() > 0 && added < n {
		cur := heap.Pop(h).(item)
		if inTree[cur.node] {
			continue // already connected more cheaply
		}
		inTree[cur.node] = true
		total += cur.dist
		added++
		for _, a := range adj[cur.node] {
			if !inTree[a.To] {
				heap.Push(h, item{a.To, a.Weight})
			}
		}
	}
	return total, added == n
}
