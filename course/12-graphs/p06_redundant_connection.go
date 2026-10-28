package graphs

// Problem: Redundant Connection
// Difficulty: Medium
//
// A tree on nodes 1..n had exactly one extra edge added, creating a single
// cycle. The edges are given in order. Return an edge whose removal turns the
// graph back into a tree; if several qualify, return the one that appears
// last in the input.
//
// Examples:
//
//	[[1,2],[1,3],[2,3]]               -> [2 3]
//	[[1,2],[2,3],[3,4],[1,4],[1,5]]   -> [1 4]
//
// Reference: https://leetcode.com/problems/redundant-connection/
//
// Pattern: union-find cycle detection. Process edges in order; the first
// edge whose endpoints are ALREADY connected closes the cycle. Every other
// cycle edge appeared earlier, so this one is also the last cycle edge in the
// input — exactly what the problem asks for.
//
// Time: O(n·α(n)); brute force O(n²). Space: O(n).

// FindRedundantConnection returns the edge that closes the cycle.
func FindRedundantConnection(edges [][]int) []int {
	uf := NewUnionFind(len(edges) + 1) // nodes are 1-based; n == len(edges)
	for _, e := range edges {
		if !uf.Union(e[0], e[1]) {
			return e
		}
	}
	return nil
}

// FindRedundantConnectionBrute tries removing each edge from last to first
// and keeps the first removal that leaves a connected graph (n-1 edges on n
// connected nodes is a tree).
func FindRedundantConnectionBrute(edges [][]int) []int {
	n := len(edges)
	for skip := n - 1; skip >= 0; skip-- {
		adj := make([][]int, n+1)
		for i, e := range edges {
			if i != skip {
				adj[e[0]] = append(adj[e[0]], e[1])
				adj[e[1]] = append(adj[e[1]], e[0])
			}
		}
		_, dist := BFS(adj, 1)
		connected := true
		for v := 1; v <= n; v++ {
			connected = connected && dist[v] >= 0
		}
		if connected {
			return edges[skip]
		}
	}
	return nil
}
