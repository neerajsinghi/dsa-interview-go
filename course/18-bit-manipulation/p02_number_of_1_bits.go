package bits

import mbits "math/bits"

// Problem: Number of 1 Bits (Hamming weight / popcount)
// Difficulty: Easy
//
// Count how many bits are 1 in the binary form of an unsigned 32-bit integer.
//
// Examples:
//
//	11 (0b1011)  -> 3
//	128          -> 1
//	0xFFFFFFFF   -> 32
//
// Reference: https://leetcode.com/problems/number-of-1-bits/
//
// Pattern: Brian Kernighan's trick. n & (n-1) clears the lowest set bit, so
// the loop runs once per SET bit instead of once per bit position.
//
//	n     = 1011 0000
//	n-1   = 1010 1111
//	n&n-1 = 1010 0000   (lowest 1 gone)
//
// Time: O(k) for k set bits (≤ 32). Space: O(1).
// In production just call math/bits.OnesCount32.

// HammingWeightShift checks every bit position.
func HammingWeightShift(n uint32) int {
	count := 0
	for ; n != 0; n >>= 1 { // unsigned >> shifts in zeros, so this terminates
		count += int(n & 1)
	}
	return count
}

// HammingWeight uses Kernighan's trick.
func HammingWeight(n uint32) int {
	count := 0
	for n != 0 {
		n &= n - 1
		count++
	}
	return count
}

// HammingWeightStdlib is what you'd write in real code.
func HammingWeightStdlib(n uint32) int { return mbits.OnesCount32(n) }
