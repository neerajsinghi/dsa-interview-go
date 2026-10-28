package prefixsums

// Problem 04 — Range Addition with a Difference Array (Medium)
//
// Start with n zeros. Each update [start, end, delta] adds delta to every
// position start..end inclusive (0-based). Return the slice after applying
// all updates.
//
//	n = 5, updates = [[1, 3, 2], [2, 4, 3], [0, 2, -2]]
//	→ [-2, 0, 3, 5, 3]
//
// Reference: https://leetcode.com/problems/range-addition/ (premium);
// same idea: https://leetcode.com/problems/corporate-flight-bookings/
//
// Pattern: the DIFFERENCE ARRAY is the inverse of prefix sums. Record each
// update with two O(1) marks — +delta at start, −delta at end+1 — and a
// single prefix-sum pass at the end spreads every delta over its range.
//
//	update [1, 3, +2]:
//	diff:   0  +2   0   0  -2   (index 4 = end+1)
//	prefix: 0   2   2   2   0   ← +2 applied to 1..3
//
//	ApplyUpdatesBrute: Time O(n·q), Space O(n)   (q = number of updates)
//	ApplyUpdates:      Time O(n + q), Space O(n)

// ApplyUpdatesBrute loops over every position of every update.
func ApplyUpdatesBrute(n int, updates [][3]int) []int {
	out := make([]int, n)
	for _, u := range updates {
		for i := u[0]; i <= u[1]; i++ {
			out[i] += u[2]
		}
	}
	return out
}

// ApplyUpdates marks range boundaries, then prefix-sums once.
func ApplyUpdates(n int, updates [][3]int) []int {
	diff := make([]int, n+1) // one extra slot so end+1 never goes out of range
	for _, u := range updates {
		start, end, delta := u[0], u[1], u[2]
		diff[start] += delta
		diff[end+1] -= delta
	}
	out := make([]int, n)
	running := 0
	for i := range out {
		running += diff[i]
		out[i] = running
	}
	return out
}
