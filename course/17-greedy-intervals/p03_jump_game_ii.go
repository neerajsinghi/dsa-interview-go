package greedy

// Problem: Jump Game II
// Difficulty: Medium
//
// Same board as Jump Game, but the last index is guaranteed reachable.
// Return the minimum number of jumps needed to get there.
//
// Examples:
//
//	[2,3,1,1,4] -> 2  (0 → 1 → 4)
//	[2,3,0,1,4] -> 2
//
// Reference: https://leetcode.com/problems/jump-game-ii/
//
// Pattern: implicit BFS by levels. Indices reachable with k jumps form a
// contiguous window [start, end]. While scanning that window, track the
// farthest index reachable with one more jump; when the scan hits `end`,
// take a jump and the new window ends at that farthest point.
//
//	nums:   2  3  1  1  4
//	jumps:  0 [1  1] [2  2]
//	        ^  window 1  window 2 contains the last index
//
// Time: O(n). Space: O(1). DP alternative: O(n²).

// JumpDP computes the min jumps to every index.
func JumpDP(nums []int) int {
	const inf = 1 << 62
	dp := make([]int, len(nums))
	for i := 1; i < len(dp); i++ {
		dp[i] = inf
	}
	for i := range nums {
		if dp[i] == inf {
			continue
		}
		for j := i + 1; j <= i+nums[i] && j < len(nums); j++ {
			dp[j] = min(dp[j], dp[i]+1)
		}
	}
	return dp[len(nums)-1]
}

// Jump is the greedy level-by-level scan.
func Jump(nums []int) int {
	jumps, end, farthest := 0, 0, 0
	for i := 0; i < len(nums)-1; i++ { // never jump FROM the last index
		farthest = max(farthest, i+nums[i])
		if i == end { // finished scanning this level
			jumps++
			end = farthest
		}
	}
	return jumps
}
