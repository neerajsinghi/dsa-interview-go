package bits

// Problem: Single Number II
// Difficulty: Medium
//
// Every value appears exactly three times except one, which appears once.
// Find it in linear time and constant extra space.
//
// Examples:
//
//	[2,2,3,2]          -> 3
//	[0,1,0,1,0,1,99]   -> 99
//
// Reference: https://leetcode.com/problems/single-number-ii/
//
// Pattern A (count per bit): for each bit position, count how many numbers
// have it set. Triples contribute multiples of 3, so count % 3 is the lone
// number's bit.
//
// Pattern B (two-bit counter): treat each bit position as a counter mod 3
// stored across two words, `ones` and `twos`:
//
//	seen 0× -> ones=0 twos=0;  1× -> ones=1 twos=0;  2× -> ones=0 twos=1;  3× -> back to 0,0
//	ones = (ones ^ x) &^ twos
//	twos = (twos ^ x) &^ ones
//
// After the scan, `ones` holds the bits seen once (mod 3) — the answer.
//
// Time: O(n) (A is O(64·n)). Space: O(1).

// SingleNumberIICount sums each of the 64 bit positions mod 3.
func SingleNumberIICount(nums []int) int {
	var result uint64
	for bit := 0; bit < 64; bit++ {
		count := 0
		for _, x := range nums {
			count += int(uint64(x) >> bit & 1) // unsigned view: no sign extension
		}
		if count%3 != 0 {
			result |= 1 << bit
		}
	}
	return int(result) // reinterpreting the bits restores negative answers
}

// SingleNumberII runs the ones/twos state machine.
func SingleNumberII(nums []int) int {
	ones, twos := 0, 0
	for _, x := range nums {
		ones = (ones ^ x) &^ twos
		twos = (twos ^ x) &^ ones
	}
	return ones
}
