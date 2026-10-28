// Package dp2d contains solutions for the 2-D dynamic programming module.
//
// Every problem here has a state with two indices — dp[i][j] — usually
// "the answer for prefix a[:i] and prefix b[:j]" (two sequences), "the answer
// for cell (i, j)" (grids) or "the answer for the interval a[i..j]".
package dp2d

// newTable allocates a rows×cols table with every cell set to fill.
// One backing slice keeps the rows contiguous, which is friendlier to the
// cache than rows separate allocations.
func newTable(rows, cols, fill int) [][]int {
	backing := make([]int, rows*cols)
	if fill != 0 {
		for i := range backing {
			backing[i] = fill
		}
	}
	t := make([][]int, rows)
	for i := range t {
		t[i] = backing[i*cols : (i+1)*cols : (i+1)*cols]
	}
	return t
}

// newBoolTable allocates a rows×cols table of false.
func newBoolTable(rows, cols int) [][]bool {
	backing := make([]bool, rows*cols)
	t := make([][]bool, rows)
	for i := range t {
		t[i] = backing[i*cols : (i+1)*cols : (i+1)*cols]
	}
	return t
}
