package shortestpaths

// Problem: Network Delay Time
// Difficulty: Medium
//
// A network has n nodes labelled 1..n and directed links {u, v, w}: a signal
// sent from u reaches v after w time units. A signal starts at node k. Return
// how long until every node has received it, or -1 if some node never does.
//
// Examples:
//
//	links [[2,1,1],[2,3,1],[3,4,1]], n=4, k=2 -> 2
//	links [[1,2,1]], n=2, k=2                 -> -1  (node 1 is unreachable)
//
// Reference: https://leetcode.com/problems/network-delay-time/
//
// Pattern: single-source shortest paths. Each node hears the signal at its
// shortest-path distance from k; everybody has heard it at the LARGEST of
// those distances. Weights are non-negative, so Dijkstra applies. Bellman-
// Ford is the slower general fallback and is shown for comparison.
//
// Time: Dijkstra O(E log V); Bellman-Ford O(V·E). Space: O(V + E).

// NetworkDelayTime uses Dijkstra.
func NetworkDelayTime(times [][]int, n, k int) int {
	adj := make([][]Arc, n+1) // index 0 unused: labels are 1-based
	for _, t := range times {
		adj[t[0]] = append(adj[t[0]], Arc{t[1], t[2]})
	}
	return slowest(Dijkstra(adj, k)[1:])
}

// NetworkDelayTimeBellmanFord uses Bellman-Ford.
func NetworkDelayTimeBellmanFord(times [][]int, n, k int) int {
	edges := make([]Edge, len(times))
	for i, t := range times {
		edges[i] = Edge{t[0], t[1], t[2]}
	}
	dist, _ := BellmanFord(n+1, edges, k) // weights are non-negative: no cycle issue
	return slowest(dist[1:])
}

// slowest returns the largest distance, or -1 if any node is unreachable.
func slowest(dist []int) int {
	worst := 0
	for _, d := range dist {
		if d == Inf {
			return -1
		}
		worst = max(worst, d)
	}
	return worst
}
