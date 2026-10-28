// Package bits contains solutions for the bit manipulation module.
//
// The helpers in this file are the vocabulary every problem uses. Bit i is
// counted from the least-significant end, starting at 0.
//
// Note: this package is itself named "bits", so the standard library's
// math/bits is imported under the alias mbits to keep the two apart.
package bits

import mbits "math/bits"

// Has reports whether bit i of x is 1.
func Has(x uint, i uint) bool { return x&(1<<i) != 0 }

// Set returns x with bit i forced to 1.
func Set(x uint, i uint) uint { return x | 1<<i }

// Clear returns x with bit i forced to 0. &^ is Go's AND NOT ("bit clear").
func Clear(x uint, i uint) uint { return x &^ (1 << i) }

// Toggle returns x with bit i flipped.
func Toggle(x uint, i uint) uint { return x ^ 1<<i }

// DropLowest clears the lowest set bit: 0b10110 -> 0b10100.
// x-1 flips the lowest 1 and every 0 below it; AND keeps the rest.
func DropLowest(x uint) uint { return x & (x - 1) }

// LowestBit isolates the lowest set bit: 0b10110 -> 0b00010.
// In two's complement -x == ^x + 1, which flips every bit above the lowest 1.
// Unsigned negation in Go wraps, so -x is well defined for uint.
func LowestBit(x uint) uint { return x & -x }

// PopCount counts set bits using the standard library (a single CPU
// instruction on most hardware).
func PopCount(x uint) int { return mbits.OnesCount(x) }

// LowMask returns a mask with the lowest n bits set: LowMask(3) == 0b111.
// For n == width, 1<<n is 0 in Go (no undefined behaviour), and 0-1 wraps to
// all ones, so the result is still correct.
func LowMask(n uint) uint { return 1<<n - 1 }
