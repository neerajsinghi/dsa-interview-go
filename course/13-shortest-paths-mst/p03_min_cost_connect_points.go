package shortestpaths

// Problem: Min Cost to Connect All Points
// Difficulty: Medium
//
// You are given points on a 2-D integer grid. Connecting two points costs
// their Manhattan distance |x1-x2| + |y1-y2|. Return the minimum total cost
// to connect all points so that every point can reach every other one.
//
// Example:
//
//	[[0,0],[2,2],[3,10],[5,2],[7,0]] -> 20
//
// Reference: https://leetcode.com/problems/min-cost-to-connect-all-points/
//
// Pattern: minimum spanning tree on a COMPLETE graph (every pair of points is
// a candidate edge, E = n(n-1)/2). That density decides the algorithm:
//   - Array-based Prim: keep minDist[v] = cheapest edge from the tree to v,
//     and each round scan for the closest outside point. O(n²) with no heap
//     and no edge list — optimal for dense graphs.
//   - Kruskal: materialise all O(n²) edges, sort, union-find. O(n² log n)
//     time and O(n²) memory; fine for small n and good to compare.
//
// Time: Prim O(n²); Kruskal O(n² log n). Space: Prim O(n); Kruskal O(n²).

// MinCostConnectPoints uses dense Prim.
func MinCostConnectPoints(points [][]int) int {
	n := len(points)
	if n == 0 {
		return 0
	}
	inTree := make([]bool, n)
	minDist := make([]int, n) // cheapest known edge from the tree to each point
	for i := range minDist {
		minDist[i] = Inf
	}
	minDist[0] = 0
	total := 0
	for round := 0; round < n; round++ {
		u := -1
		for v := 0; v < n; v++ { // closest point not yet in the tree
			if !inTree[v] && (u == -1 || minDist[v] < minDist[u]) {
				u = v
			}
		}
		inTree[u] = true
		total += minDist[u]
		for v := 0; v < n; v++ { // u may offer cheaper edges to the rest
			if !inTree[v] {
				minDist[v] = min(minDist[v], manhattan(points[u], points[v]))
			}
		}
	}
	return total
}

// MinCostConnectPointsKruskal builds every edge and runs Kruskal.
func MinCostConnectPointsKruskal(points [][]int) int {
	n := len(points)
	edges := make([]Edge, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, Edge{i, j, manhattan(points[i], points[j])})
		}
	}
	total, _, _ := Kruskal(n, edges)
	return total
}

func manhattan(a, b []int) int {
	return abs(a[0]-b[0]) + abs(a[1]-b[1])
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
