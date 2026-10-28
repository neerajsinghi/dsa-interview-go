package binarysearch

import "slices"

// Problem: Koko Eating Bananas
// Difficulty: Medium
//
// There are piles of bananas, piles[i] in pile i. Each hour Koko picks one
// pile and eats up to k bananas from it (if the pile has fewer, she finishes
// it and waits out the rest of the hour). Find the smallest integer speed k
// that lets her finish every pile within h hours. Assume len(piles) <= h.
//
// Examples:
//
//	piles=[3,6,7,11], h=8       -> 4
//	piles=[30,11,23,4,20], h=5  -> 30
//	piles=[30,11,23,4,20], h=6  -> 23
//
// Reference: https://leetcode.com/problems/koko-eating-bananas/
//
// Pattern: binary search on the answer. Hours needed only goes down as k
// grows, so "can finish at speed k" is false...false true...true over
// k in [1, max(piles)]. Find the first true with FirstTrue.
//
// Time: brute O(max·n); optimal O(n log max). Space: O(1).

// hoursAt returns the hours needed to eat every pile at speed k.
func hoursAt(piles []int, k int) int {
	hours := 0
	for _, p := range piles {
		hours += (p + k - 1) / k // ceil(p / k) without floats
	}
	return hours
}

// MinEatingSpeedBrute tries every speed from 1 upward.
func MinEatingSpeedBrute(piles []int, h int) int {
	for k := 1; ; k++ {
		if hoursAt(piles, k) <= h {
			return k
		}
	}
}

// MinEatingSpeed binary-searches the smallest feasible speed.
func MinEatingSpeed(piles []int, h int) int {
	if len(piles) == 0 {
		return 1
	}
	// Speed max(piles) always works (one hour per pile), so search [1, max].
	// FirstTrue's hi is exclusive, hence +1.
	return FirstTrue(1, slices.Max(piles)+1, func(k int) bool {
		return hoursAt(piles, k) <= h
	})
}
