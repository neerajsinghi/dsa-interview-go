package graphs

// Problem: Number of Islands
// Difficulty: Medium
//
// A grid holds '1' for land and '0' for water. Land cells that touch
// horizontally or vertically belong to the same island. Count the islands.
//
// Example:
//
//	1 1 0 0 0
//	1 1 0 0 0     -> 3
//	0 0 1 0 0
//	0 0 0 1 1
//
// Reference: https://leetcode.com/problems/number-of-islands/
//
// Pattern: connected components on an implicit grid graph. Scan every cell;
// each time you meet unvisited land, that is a new island — flood-fill it
// (DFS or BFS) so its other cells are not counted again. Union-find gives the
// same answer by merging adjacent land cells and counting groups, which is
// handy when land arrives incrementally ("Number of Islands II").
//
// Time: O(m·n) for both. Space: O(m·n) for the visited grid / union-find.

// NumIslands flood-fills each island with DFS. The input is not modified.
func NumIslands(grid [][]byte) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])
	seen := make([][]bool, m)
	for i := range seen {
		seen[i] = make([]bool, n)
	}
	var sink func(r, c int)
	sink = func(r, c int) {
		if r < 0 || r >= m || c < 0 || c >= n || grid[r][c] != '1' || seen[r][c] {
			return
		}
		seen[r][c] = true
		for _, d := range dirs4 {
			sink(r+d[0], c+d[1])
		}
	}
	islands := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] == '1' && !seen[r][c] {
				islands++
				sink(r, c)
			}
		}
	}
	return islands
}

// NumIslandsUnionFind merges each land cell with its right and lower land
// neighbours (enough to cover every adjacency once), then counts groups.
func NumIslandsUnionFind(grid [][]byte) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])
	uf := NewUnionFind(m * n)
	water := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] != '1' {
				water++ // each water cell is its own group: subtract it later
				continue
			}
			id := r*n + c // flatten (r, c) into one index
			if c+1 < n && grid[r][c+1] == '1' {
				uf.Union(id, id+1)
			}
			if r+1 < m && grid[r+1][c] == '1' {
				uf.Union(id, id+n)
			}
		}
	}
	return uf.Count() - water
}
