package bits

import mbits "math/bits"

// Problem: Reverse Bits
// Difficulty: Easy
//
// Reverse the order of all 32 bits of an unsigned 32-bit integer: bit 0
// swaps with bit 31, bit 1 with bit 30, and so on.
//
// Examples:
//
//	0b00000010100101000001111010011100 -> 0b00111001011110000010100101000000
//	1                                  -> 0x80000000
//
// Reference: https://leetcode.com/problems/reverse-bits/
//
// Pattern 1 (loop): pop the low bit of n, push it onto the low end of the
// result; after 32 steps the first bit popped is at the top.
//
// Pattern 2 (divide and conquer): swap halves, then quarters, … down to
// adjacent bits, using masks. Five steps, no loop:
//
//	swap 16-bit halves, 8-bit bytes, 4-bit nibbles, 2-bit pairs, 1-bit neighbours
//
// Time: O(32) loop / O(5) masks. Space: O(1). Stdlib: math/bits.Reverse32.

// ReverseBitsLoop shifts one bit at a time. It must run all 32 iterations —
// stopping when n becomes 0 would lose the leading zeros' positions.
func ReverseBitsLoop(n uint32) uint32 {
	var out uint32
	for i := 0; i < 32; i++ {
		out = out<<1 | n&1
		n >>= 1
	}
	return out
}

// ReverseBits swaps progressively smaller blocks with masks.
func ReverseBits(n uint32) uint32 {
	n = n>>16 | n<<16
	n = (n&0xFF00FF00)>>8 | (n&0x00FF00FF)<<8
	n = (n&0xF0F0F0F0)>>4 | (n&0x0F0F0F0F)<<4
	n = (n&0xCCCCCCCC)>>2 | (n&0x33333333)<<2 // 1100… / 0011…
	n = (n&0xAAAAAAAA)>>1 | (n&0x55555555)<<1 // 1010… / 0101…
	return n
}

// ReverseBitsStdlib is what you'd ship.
func ReverseBitsStdlib(n uint32) uint32 { return mbits.Reverse32(n) }
