// Package graphs holds reusable graph templates — adjacency lists, BFS, DFS,
// topological sort and union-find — and the interview problems built on them.
//
// Conventions: nodes are the integers 0..n-1 unless a problem says otherwise;
// an edge is a two-element slice {u, v}.
package graphs

// dirs4 are the four grid moves: down, up, right, left.
var dirs4 = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

// AdjList turns an edge list into an adjacency list for nodes 0..n-1.
// For an undirected graph each edge is recorded in both directions.
func AdjList(n int, edges [][]int, directed bool) [][]int {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		if !directed {
			adj[v] = append(adj[v], u)
		}
	}
	return adj
}

// BFS explores the graph level by level from start. It returns the nodes in
// visit order and dist[v] = fewest edges from start to v (-1 if unreachable).
// In an unweighted graph BFS distance IS the shortest-path distance.
//
// Time: O(V + E). Space: O(V).
func BFS(adj [][]int, start int) (order, dist []int) {
	dist = make([]int, len(adj))
	for i := range dist {
		dist[i] = -1
	}
	dist[start] = 0
	queue := []int{start}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			if dist[v] == -1 { // mark when ENQUEUED, not when dequeued
				dist[v] = dist[u] + 1
				queue = append(queue, v)
			}
		}
	}
	return order, dist
}

// DFS returns the nodes reachable from start in recursive pre-order.
//
// Time: O(V + E). Space: O(V) for visited plus recursion depth.
func DFS(adj [][]int, start int) []int {
	visited := make([]bool, len(adj))
	var order []int
	var visit func(u int)
	visit = func(u int) {
		visited[u] = true
		order = append(order, u)
		for _, v := range adj[u] {
			if !visited[v] {
				visit(v)
			}
		}
	}
	visit(start)
	return order
}

// DFSIterative produces the same order as DFS using an explicit stack, which
// avoids deep recursion on huge graphs. Neighbours are pushed in reverse so
// the first neighbour is popped first, and a node is marked visited when it
// is popped (a node can sit on the stack more than once).
func DFSIterative(adj [][]int, start int) []int {
	visited := make([]bool, len(adj))
	var order []int
	stack := []int{start}
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[u] {
			continue
		}
		visited[u] = true
		order = append(order, u)
		for i := len(adj[u]) - 1; i >= 0; i-- {
			if v := adj[u][i]; !visited[v] {
				stack = append(stack, v)
			}
		}
	}
	return order
}

// TopoSortKahn orders the nodes of a directed graph so every edge u->v has u
// before v. It repeatedly removes a node with in-degree 0. If some nodes are
// never freed, they sit on a cycle: ok is false and order is partial.
//
// Time: O(V + E). Space: O(V).
func TopoSortKahn(adj [][]int) (order []int, ok bool) {
	n := len(adj)
	indeg := make([]int, n)
	for u := range adj {
		for _, v := range adj[u] {
			indeg[v]++
		}
	}
	var queue []int
	for v := 0; v < n; v++ {
		if indeg[v] == 0 {
			queue = append(queue, v)
		}
	}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			if indeg[v]--; indeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	return order, len(order) == n
}

// Node colours for DFS cycle detection.
const (
	white = iota // not visited yet
	gray         // on the current DFS path
	black        // fully explored
)

// TopoSortDFS is the DFS flavour: a node is appended after all its
// descendants (post-order), so the reversed post-order is a topological
// order. Reaching a gray node means a back edge, i.e. a cycle.
//
// Time: O(V + E). Space: O(V).
func TopoSortDFS(adj [][]int) (order []int, ok bool) {
	color := make([]int, len(adj))
	var visit func(u int) bool
	visit = func(u int) bool {
		color[u] = gray
		for _, v := range adj[u] {
			if color[v] == gray {
				return false // back edge: cycle
			}
			if color[v] == white && !visit(v) {
				return false
			}
		}
		color[u] = black
		order = append(order, u)
		return true
	}
	for u := range adj {
		if color[u] == white && !visit(u) {
			return nil, false
		}
	}
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order, true
}

// UnionFind (disjoint-set union) tracks a partition of 0..n-1 into groups.
// With path compression and union by rank, each operation costs amortised
// O(α(n)) — the inverse Ackermann function, effectively constant.
type UnionFind struct {
	parent []int
	rank   []int // upper bound on tree height
	count  int   // number of disjoint groups
}

// NewUnionFind starts with every element in its own group.
func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{parent: make([]int, n), rank: make([]int, n), count: n}
	for i := range uf.parent {
		uf.parent[i] = i
	}
	return uf
}

// Find returns the representative of x's group, pointing every node on the
// way directly at the root (path compression).
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

// Union merges the groups of a and b. It returns false if they were already
// in the same group — exactly the signal for "this edge closes a cycle".
func (uf *UnionFind) Union(a, b int) bool {
	ra, rb := uf.Find(a), uf.Find(b)
	if ra == rb {
		return false
	}
	// Union by rank: hang the shorter tree under the taller one.
	switch {
	case uf.rank[ra] < uf.rank[rb]:
		uf.parent[ra] = rb
	case uf.rank[ra] > uf.rank[rb]:
		uf.parent[rb] = ra
	default:
		uf.parent[rb] = ra
		uf.rank[ra]++
	}
	uf.count--
	return true
}

// Connected reports whether a and b are in the same group.
func (uf *UnionFind) Connected(a, b int) bool {
	return uf.Find(a) == uf.Find(b)
}

// Count returns the number of disjoint groups.
func (uf *UnionFind) Count() int {
	return uf.count
}
