package greedy

// Problem: Gas Station
// Difficulty: Medium
//
// Stations sit on a circular road. At station i you may fill gas[i] units;
// driving from i to i+1 burns cost[i]. Starting with an empty tank, return
// the index of a station from which you can drive one full lap clockwise,
// or -1 if none exists. (When an answer exists it is unique.)
//
// Examples:
//
//	gas [1,2,3,4,5], cost [3,4,5,1,2] -> 3
//	gas [2,3,4],     cost [3,4,3]     -> -1
//
// Reference: https://leetcode.com/problems/gas-station/
//
// Pattern: greedy restart. Two facts:
//  1. If total gas < total cost, no start works.
//  2. If starting at s you first run dry before station k, then NO station
//     between s and k works either (each began with ≥ 0 extra fuel and still
//     failed), so jump the candidate straight to k+1.
//
// Time: O(n). Space: O(1). Brute force: O(n²).

// CanCompleteCircuitBrute simulates a lap from every start.
func CanCompleteCircuitBrute(gas, cost []int) int {
	n := len(gas)
	for s := 0; s < n; s++ {
		tank, ok := 0, true
		for k := 0; k < n; k++ {
			i := (s + k) % n
			tank += gas[i] - cost[i]
			if tank < 0 {
				ok = false
				break
			}
		}
		if ok {
			return s
		}
	}
	return -1
}

// CanCompleteCircuit is the one-pass greedy.
func CanCompleteCircuit(gas, cost []int) int {
	total, tank, start := 0, 0, 0
	for i := range gas {
		diff := gas[i] - cost[i]
		total += diff
		tank += diff
		if tank < 0 { // can't reach i+1 from start: none of start..i can
			start, tank = i+1, 0
		}
	}
	if total < 0 {
		return -1
	}
	return start
}
