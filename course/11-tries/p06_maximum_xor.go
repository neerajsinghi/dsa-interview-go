package tries

// Problem: Maximum XOR of Two Numbers in an Array
// Difficulty: Medium
//
// Given a list of non-negative integers below 2^31, return the largest value
// of a XOR b over all pairs (a may equal b, so a one-element list gives 0).
//
// Examples:
//
//	[3, 10, 5, 25, 2, 8] -> 28   (5 XOR 25 = 00101 ^ 11001 = 11100)
//	[0]                  -> 0
//
// Reference: https://leetcode.com/problems/maximum-xor-of-two-numbers-in-an-array/
//
// Pattern: bitwise trie. Store every number as a 31-level path of bits, most
// significant bit first. To maximise x XOR y for a fixed x, walk the trie
// greedily from the top bit and prefer the child with the *opposite* bit:
// winning a higher bit is worth more than all lower bits combined
// (2^k > 2^k - 1), so greed is optimal.
//
//	x = 00101, at each level ask "is there a stored number with the other bit?"
//	bit4: x has 0, want 1 -> 25 has it, take it   (XOR gains 16)
//	bit3: x has 0, want 1 -> 25 has it            (gains 8)
//	bit2: x has 1, want 0 -> 25 has 0             (gains 4)  ... = 28
//
// Time: O(n · B) with B = 31 bits; brute force O(n²). Space: O(n · B) nodes.

// maxBit is the highest bit position of a value below 2^31.
const maxBit = 30

// bitNode is a binary trie node: children[0] and children[1].
type bitNode struct {
	children [2]*bitNode
}

// FindMaximumXOR uses the bitwise trie.
func FindMaximumXOR(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	root := &bitNode{}
	best := 0
	for _, x := range nums {
		insertBits(root, x)
		// Query after inserting, so x can pair with itself (XOR 0) and with
		// every number seen before it; each pair is still considered once.
		if v := bestXOR(root, x); v > best {
			best = v
		}
	}
	return best
}

func insertBits(root *bitNode, x int) {
	cur := root
	for b := maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		if cur.children[bit] == nil {
			cur.children[bit] = &bitNode{}
		}
		cur = cur.children[bit]
	}
}

// bestXOR returns the maximum x XOR y over the y values stored in the trie.
func bestXOR(root *bitNode, x int) int {
	cur, xor := root, 0
	for b := maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		if want := bit ^ 1; cur.children[want] != nil {
			xor |= 1 << b // this bit differs: it shows up in the XOR
			cur = cur.children[want]
		} else {
			cur = cur.children[bit]
		}
	}
	return xor
}

// FindMaximumXORBrute tries every pair.
func FindMaximumXORBrute(nums []int) int {
	best := 0
	for i := range nums {
		for j := i; j < len(nums); j++ {
			best = max(best, nums[i]^nums[j])
		}
	}
	return best
}
