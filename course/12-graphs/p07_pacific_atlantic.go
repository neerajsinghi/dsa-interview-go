package graphs

// Problem: Pacific Atlantic Water Flow
// Difficulty: Medium
//
// A grid gives the height of each cell of an island. The Pacific touches the
// top and left edges, the Atlantic the bottom and right edges. Rain flows
// from a cell to a neighbour (up/down/left/right) whose height is less than
// or equal to its own, and from edge cells into the adjacent ocean. Return
// every cell from which water can reach BOTH oceans, in row-major order.
//
// Example:
//
//	heights            result (cells marked *)
//	1 2 2 3 5          . . . . *
//	3 2 3 4 4          . . . * *
//	2 4 5 3 1          . . * . .
//	6 7 1 4 5          * * . . .
//	5 1 1 2 4          * . . . .
//
// Reference: https://leetcode.com/problems/pacific-atlantic-water-flow/
//
// Pattern: reverse the flow and search from the targets. Instead of asking
// "where can this cell drain?" for every cell (a search per cell), start at
// each ocean's border and climb UPHILL (to neighbours at least as high).
// Everything the Pacific climb reaches can drain into the Pacific; same for
// the Atlantic. The answer is the intersection. Two multi-source searches
// replace m·n single-source ones.
//
// Time: O(m·n); brute force O((m·n)²). Space: O(m·n).

// PacificAtlantic climbs uphill from both coastlines with BFS.
func PacificAtlantic(heights [][]int) [][]int {
	m := len(heights)
	if m == 0 {
		return nil
	}
	n := len(heights[0])
	climb := func(starts [][2]int) [][]bool {
		seen := make([][]bool, m)
		for i := range seen {
			seen[i] = make([]bool, n)
		}
		for _, s := range starts {
			seen[s[0]][s[1]] = true
		}
		queue := starts
		for len(queue) > 0 {
			cell := queue[0]
			queue = queue[1:]
			for _, d := range dirs4 {
				r, c := cell[0]+d[0], cell[1]+d[1]
				if r < 0 || r >= m || c < 0 || c >= n || seen[r][c] {
					continue
				}
				if heights[r][c] >= heights[cell[0]][cell[1]] { // uphill or level
					seen[r][c] = true
					queue = append(queue, [2]int{r, c})
				}
			}
		}
		return seen
	}

	var pacific, atlantic [][2]int
	for r := 0; r < m; r++ {
		pacific = append(pacific, [2]int{r, 0})
		atlantic = append(atlantic, [2]int{r, n - 1})
	}
	for c := 0; c < n; c++ {
		pacific = append(pacific, [2]int{0, c})
		atlantic = append(atlantic, [2]int{m - 1, c})
	}
	// Duplicate corner starts are harmless: seen is set before the BFS.
	toP, toA := climb(pacific), climb(atlantic)

	var res [][]int
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if toP[r][c] && toA[r][c] {
				res = append(res, []int{r, c})
			}
		}
	}
	return res
}

// PacificAtlanticBrute runs a downhill DFS from every cell and records which
// oceans it touches.
func PacificAtlanticBrute(heights [][]int) [][]int {
	m := len(heights)
	if m == 0 {
		return nil
	}
	n := len(heights[0])
	var res [][]int
	for sr := 0; sr < m; sr++ {
		for sc := 0; sc < n; sc++ {
			seen := map[[2]int]bool{}
			pac, atl := false, false
			var flow func(r, c int)
			flow = func(r, c int) {
				seen[[2]int{r, c}] = true
				pac = pac || r == 0 || c == 0
				atl = atl || r == m-1 || c == n-1
				for _, d := range dirs4 {
					nr, nc := r+d[0], c+d[1]
					if nr >= 0 && nr < m && nc >= 0 && nc < n && !seen[[2]int{nr, nc}] && heights[nr][nc] <= heights[r][c] {
						flow(nr, nc)
					}
				}
			}
			flow(sr, sc)
			if pac && atl {
				res = append(res, []int{sr, sc})
			}
		}
	}
	return res
}
