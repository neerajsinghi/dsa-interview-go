package shortestpaths

// Problem: Cheapest Flights Within K Stops
// Difficulty: Medium
//
// There are n cities 0..n-1 and one-way flights {from, to, price}. Find the
// cheapest total price from src to dst using at most k intermediate stops
// (so at most k+1 flights). Return -1 if no such route exists.
//
// Example:
//
//	n=4, flights [[0,1,100],[1,2,100],[2,0,100],[1,3,600],[2,3,200]]
//	src=0, dst=3, k=1 -> 700   (0→1→3; the cheaper 0→1→2→3 needs 2 stops)
//
// Reference: https://leetcode.com/problems/cheapest-flights-within-k-stops/
//
// Pattern: Bellman-Ford with a hop limit. After round i of Bellman-Ford,
// dist[v] is the cheapest price using at most i edges — but only if each
// round relaxes from a SNAPSHOT of the previous round. Relaxing in place can
// chain several flights in one round and silently exceed the stop limit.
// Run exactly k+1 rounds.
//
// Plain Dijkstra on cost alone is wrong here: the cheapest way to reach a
// city may use too many stops, while a pricier way with fewer stops is the
// one that can still continue to dst.
//
// Time: O(k · E). Space: O(n).

// FindCheapestPrice uses hop-limited Bellman-Ford.
func FindCheapestPrice(n int, flights [][]int, src, dst, k int) int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = Inf
	}
	dist[src] = 0
	for round := 0; round <= k; round++ { // k stops = k+1 flights
		prev := append([]int(nil), dist...) // snapshot: only use paths with < round+1 edges
		for _, f := range flights {
			from, to, price := f[0], f[1], f[2]
			if prev[from] != Inf && prev[from]+price < dist[to] {
				dist[to] = prev[from] + price
			}
		}
	}
	if dist[dst] == Inf {
		return -1
	}
	return dist[dst]
}

// FindCheapestPriceBrute tries every route of at most k+1 flights with DFS.
// Exponential; small inputs only.
func FindCheapestPriceBrute(n int, flights [][]int, src, dst, k int) int {
	adj := make([][]Arc, n)
	for _, f := range flights {
		adj[f[0]] = append(adj[f[0]], Arc{f[1], f[2]})
	}
	best := Inf
	var dfs func(u, cost, flightsLeft int)
	dfs = func(u, cost, flightsLeft int) {
		if u == dst {
			best = min(best, cost)
			// Keep going is pointless: prices are non-negative.
			return
		}
		if flightsLeft == 0 {
			return
		}
		for _, a := range adj[u] {
			dfs(a.To, cost+a.Weight, flightsLeft-1)
		}
	}
	dfs(src, 0, k+1)
	if best == Inf {
		return -1
	}
	return best
}
