package shortestpaths

// Problem: Find the City With the Smallest Number of Neighbours at a Threshold Distance
// Difficulty: Medium
//
// There are n cities 0..n-1 joined by undirected weighted roads {u, v, w}.
// For each city, count the other cities reachable with total distance at
// most threshold. Return the city with the smallest count; on a tie, return
// the one with the largest label.
//
// Example:
//
//	n=4, roads [[0,1,3],[1,2,1],[1,3,4],[2,3,1]], threshold 4
//	  city 0 reaches {1,2}   city 1 reaches {0,2,3}
//	  city 2 reaches {0,1,3} city 3 reaches {1,2}       -> 3 (tie with 0, larger wins)
//
// Reference: https://leetcode.com/problems/find-the-city-with-the-smallest-number-of-neighbors-at-a-threshold-distance/
//
// Pattern: all-pairs shortest paths. With n ≤ ~100, Floyd-Warshall's O(n³)
// triple loop is the simplest correct tool: it gives every pairwise distance
// in one go. Running Dijkstra from every city is the alternative, better for
// large sparse graphs: O(n · E log n).
//
// Time: O(n³) Floyd-Warshall; O(n·E log n) repeated Dijkstra. Space: O(n²).

// FindTheCity uses Floyd-Warshall.
func FindTheCity(n int, roads [][]int, threshold int) int {
	edges := make([]Edge, 0, 2*len(roads))
	for _, r := range roads {
		// Undirected road = two directed edges.
		edges = append(edges, Edge{r[0], r[1], r[2]}, Edge{r[1], r[0], r[2]})
	}
	return pickCity(FloydWarshall(n, edges), threshold)
}

// FindTheCityDijkstra runs Dijkstra from every city.
func FindTheCityDijkstra(n int, roads [][]int, threshold int) int {
	adj := make([][]Arc, n)
	for _, r := range roads {
		adj[r[0]] = append(adj[r[0]], Arc{r[1], r[2]})
		adj[r[1]] = append(adj[r[1]], Arc{r[0], r[2]})
	}
	dist := make([][]int, n)
	for src := range dist {
		dist[src] = Dijkstra(adj, src)
	}
	return pickCity(dist, threshold)
}

// pickCity applies the counting and tie-break rule to a distance matrix.
func pickCity(dist [][]int, threshold int) int {
	best, bestCount := -1, len(dist)+1
	for i := range dist {
		count := 0
		for j, d := range dist[i] {
			if i != j && d <= threshold {
				count++
			}
		}
		if count <= bestCount { // <= so a later (larger) city wins ties
			best, bestCount = i, count
		}
	}
	return best
}
