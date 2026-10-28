package graphs

// Problem: Number of Connected Components in an Undirected Graph
// Difficulty: Medium
//
// Given n nodes labelled 0..n-1 and a list of undirected edges, count the
// connected components (groups of nodes that can reach each other).
//
// Examples:
//
//	n=5, [[0,1],[1,2],[3,4]]       -> 2   ({0,1,2} and {3,4})
//	n=5, [[0,1],[1,2],[2,3],[3,4]] -> 1
//
// Reference: https://leetcode.com/problems/number-of-connected-components-in-an-undirected-graph/
//
// Pattern: union-find. Start with n singleton groups; every edge that joins
// two different groups reduces the count by one. DFS from each unvisited
// node is the equally valid traversal answer — union-find shines when edges
// stream in and you need the count after each one.
//
// Time: union-find O(E·α(n)) ≈ O(E); DFS O(V + E).
// Space: O(V) for union-find; O(V + E) for DFS (adjacency list).

// CountComponents uses union-find.
func CountComponents(n int, edges [][]int) int {
	uf := NewUnionFind(n)
	for _, e := range edges {
		uf.Union(e[0], e[1])
	}
	return uf.Count()
}

// CountComponentsDFS starts a DFS from every node not yet reached.
func CountComponentsDFS(n int, edges [][]int) int {
	adj := AdjList(n, edges, false)
	visited := make([]bool, n)
	var visit func(u int)
	visit = func(u int) {
		visited[u] = true
		for _, v := range adj[u] {
			if !visited[v] {
				visit(v)
			}
		}
	}
	components := 0
	for u := 0; u < n; u++ {
		if !visited[u] {
			components++
			visit(u)
		}
	}
	return components
}
