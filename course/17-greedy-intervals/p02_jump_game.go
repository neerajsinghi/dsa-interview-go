package greedy

// Problem: Jump Game
// Difficulty: Medium
//
// You start at index 0. nums[i] is the MAXIMUM jump length from index i (you
// may jump any distance from 0 up to it). Can you reach the last index?
//
// Examples:
//
//	[2,3,1,1,4] -> true
//	[3,2,1,0,4] -> false  (every route lands on the 0 at index 3)
//
// Reference: https://leetcode.com/problems/jump-game/
//
// Pattern: greedy "farthest reach". Keep the farthest index reachable so far.
// If the scan ever steps past it, you are stuck. Every index ≤ reach is
// reachable, so one pass suffices.
//
// Time: O(n). Space: O(1). DP alternative: O(n²).

// CanJumpDP marks reachable indices one by one.
func CanJumpDP(nums []int) bool {
	ok := make([]bool, len(nums))
	ok[0] = true
	for i := range nums {
		if !ok[i] {
			continue
		}
		for j := i + 1; j <= i+nums[i] && j < len(nums); j++ {
			ok[j] = true
		}
	}
	return ok[len(nums)-1]
}

// CanJump is the greedy farthest-reach scan.
func CanJump(nums []int) bool {
	reach := 0
	for i, n := range nums {
		if i > reach {
			return false // index i is unreachable, so is everything after it
		}
		reach = max(reach, i+n)
		if reach >= len(nums)-1 {
			return true
		}
	}
	return true
}
