package stacks

import "slices"

// Problem: Car Fleet
// Difficulty: Medium
//
// Cars drive toward a target on a one-lane road. Car i starts at position[i]
// (all distinct, all < target) with constant speed[i]. A faster car that
// catches a slower one cannot pass; it slows down and they travel together as
// a "fleet" from then on. Catching up exactly at the target also counts as
// joining. How many fleets arrive at the target?
//
// Examples:
//
//	target=12, pos=[10,8,0,5,3], speed=[2,4,1,1,3] -> 3
//	target=10, pos=[3], speed=[3]                   -> 1
//
// Reference: https://leetcode.com/problems/car-fleet/
//
// Pattern: sort, then a monotonic stack of arrival times. Process cars from
// the one closest to the target backwards. A car whose solo arrival time is
// <= the fleet ahead of it merges into that fleet; a strictly slower one
// starts a new fleet (push). The stack ends up strictly increasing.
//
// Go gotcha: arrival time is a fraction (distance/speed). Comparing float64
// values can misjudge exact ties, so we compare fractions by cross-multiplying.
//
// Time: O(n log n) for the sort. Space: O(n).

// arrival is distance/speed kept as an exact fraction.
type arrival struct{ dist, speed int }

// slowerThan reports a > b, i.e. a.dist/a.speed > b.dist/b.speed.
func (a arrival) slowerThan(b arrival) bool {
	return a.dist*b.speed > b.dist*a.speed
}

// CarFleet returns the number of fleets that reach target.
func CarFleet(target int, position, speed []int) int {
	idx := make([]int, len(position))
	for i := range idx {
		idx[i] = i
	}
	// Closest to the target first.
	slices.SortFunc(idx, func(a, b int) int { return position[b] - position[a] })

	var fleets Stack[arrival] // arrival time of each fleet; times grow from bottom to top
	for _, i := range idx {
		t := arrival{dist: target - position[i], speed: speed[i]}
		if ahead, ok := fleets.Peek(); ok && !t.slowerThan(ahead) {
			continue // catches the fleet ahead before (or at) the target: merges
		}
		fleets.Push(t)
	}
	return fleets.Len()
}
