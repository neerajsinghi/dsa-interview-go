package dp2d

// Problem: Edit Distance (Levenshtein distance)
// Difficulty: Medium
//
// You may insert a character, delete a character, or replace one character
// with another; each costs 1. Return the fewest operations that turn word a
// into word b.
//
// Examples:
//
//	"horse", "ros"           -> 3  (horse→rorse→rose→ros)
//	"intention", "execution" -> 5
//
// Reference: https://leetcode.com/problems/edit-distance/
//
// Pattern: two-sequence DP. dp[i][j] = cost to turn a[:i] into b[:j].
//
//	dp[i][0] = i (delete everything), dp[0][j] = j (insert everything)
//	a[i-1] == b[j-1] -> dp[i-1][j-1]                    (free)
//	otherwise        -> 1 + min(dp[i-1][j],   // delete a[i-1]
//	                            dp[i][j-1],   // insert b[j-1]
//	                            dp[i-1][j-1]) // replace
//
// Time: O(m·n). Space: O(m·n) table, O(n) with one row + a saved diagonal.

// EditDistanceMemo is top-down over suffixes a[i:], b[j:].
func EditDistanceMemo(a, b string) int {
	memo := newTable(len(a)+1, len(b)+1, -1)
	var dist func(i, j int) int
	dist = func(i, j int) int {
		if i == len(a) {
			return len(b) - j // insert the rest of b
		}
		if j == len(b) {
			return len(a) - i // delete the rest of a
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		if a[i] == b[j] {
			memo[i][j] = dist(i+1, j+1)
		} else {
			memo[i][j] = 1 + min(dist(i+1, j), dist(i, j+1), dist(i+1, j+1))
		}
		return memo[i][j]
	}
	return dist(0, 0)
}

// EditDistanceTab is the bottom-up prefix table.
func EditDistanceTab(a, b string) int {
	m, n := len(a), len(b)
	dp := newTable(m+1, n+1, 0)
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return dp[m][n]
}

// EditDistanceRow uses one row. The cell we overwrite (row[j], "above") is
// needed as the diagonal for column j+1, so it is saved in diag first.
func EditDistanceRow(a, b string) int {
	n := len(b)
	row := make([]int, n+1)
	for j := range row {
		row[j] = j
	}
	for i := 1; i <= len(a); i++ {
		diag := row[0] // dp[i-1][0]
		row[0] = i
		for j := 1; j <= n; j++ {
			above := row[j]
			if a[i-1] == b[j-1] {
				row[j] = diag
			} else {
				row[j] = 1 + min(above, row[j-1], diag)
			}
			diag = above
		}
	}
	return row[n]
}
