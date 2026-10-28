// Package stringalgos contains solutions for the strings module.
//
// (It is not called "strings" so it never shadows the standard library
// package of that name.)
//
// The three linear-time building blocks below — prefix function, Z-function
// and a polynomial rolling hash — solve most "find a pattern" and
// "prefix that is also a suffix" questions. All work on bytes; convert to
// []rune first if you need to treat multi-byte UTF-8 characters as units.
package stringalgos

// PrefixFunction returns pi where pi[i] is the length of the longest PROPER
// prefix of s[:i+1] that is also a suffix of it (the KMP failure table).
//
//	s  = a b a b c a b a b a
//	pi = 0 0 1 2 0 1 2 3 4 3
//
// Runs in O(n): k grows by at most one per step, so it can only shrink O(n)
// times in total.
func PrefixFunction(s string) []int {
	pi := make([]int, len(s))
	for i := 1; i < len(s); i++ {
		k := pi[i-1] // length of the border we try to extend
		for k > 0 && s[i] != s[k] {
			k = pi[k-1] // fall back to the next shorter border
		}
		if s[i] == s[k] {
			k++
		}
		pi[i] = k
	}
	return pi
}

// ZFunction returns z where z[i] is the length of the longest common prefix
// of s and s[i:]. By convention z[0] = len(s).
//
//	s = a a b x a a b
//	z = 7 1 0 0 3 1 0
//
// [l, r) is the rightmost window known to match a prefix of s; inside it z
// values can be copied from earlier positions instead of recomputed. O(n).
func ZFunction(s string) []int {
	n := len(s)
	z := make([]int, n)
	if n == 0 {
		return z
	}
	z[0] = n
	l, r := 0, 0
	for i := 1; i < n; i++ {
		if i < r {
			z[i] = min(r-i, z[i-l]) // reuse the mirror inside the window
		}
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			z[i]++
		}
		if i+z[i] > r {
			l, r = i, i+z[i]
		}
	}
	return z
}

// Rolling-hash parameters. With mod < 2^31 every product of two residues
// fits in a uint64, so no overflow tricks are needed.
const (
	hashBase = 131
	hashMod  = 1_000_000_007
)

// RollingHash answers "hash of s[l:r]" in O(1) after O(n) preprocessing.
//
//	h(s[l:r]) = prefix[r] - prefix[l]·base^(r-l)   (mod p)
//
// Equal strings always have equal hashes; unequal strings collide with
// probability about (r-l)/p, so verify a hit when correctness matters.
type RollingHash struct {
	prefix []uint64 // prefix[i] = hash of s[:i]
	pow    []uint64 // pow[i] = base^i mod p
}

// NewRollingHash precomputes prefix hashes and powers for s.
func NewRollingHash(s string) *RollingHash {
	h := &RollingHash{
		prefix: make([]uint64, len(s)+1),
		pow:    make([]uint64, len(s)+1),
	}
	h.pow[0] = 1
	for i := 0; i < len(s); i++ {
		h.prefix[i+1] = (h.prefix[i]*hashBase + uint64(s[i])) % hashMod
		h.pow[i+1] = h.pow[i] * hashBase % hashMod
	}
	return h
}

// Hash returns the hash of s[l:r].
func (h *RollingHash) Hash(l, r int) uint64 {
	// Add hashMod before subtracting so the unsigned value never wraps.
	return (h.prefix[r] + hashMod - h.prefix[l]*h.pow[r-l]%hashMod) % hashMod
}
