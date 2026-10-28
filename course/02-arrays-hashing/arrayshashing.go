// Package arrayshashing covers the most common first step in interview
// problems: trading O(n) memory for O(1) lookups with Go maps, sets and
// fixed-size count arrays.
package arrayshashing

// Counter returns how many times each value occurs in xs.
// It is the Go equivalent of Python's collections.Counter.
func Counter[T comparable](xs []T) map[T]int {
	counts := make(map[T]int, len(xs))
	for _, x := range xs {
		counts[x]++ // the zero value of a missing key is 0, so this just works
	}
	return counts
}

// Set is a hash set built on map[T]struct{}; struct{} occupies no memory.
type Set[T comparable] map[T]struct{}

// NewSet builds a set containing every value in xs.
func NewSet[T comparable](xs []T) Set[T] {
	s := make(Set[T], len(xs))
	for _, x := range xs {
		s[x] = struct{}{}
	}
	return s
}

// Add inserts x into the set.
func (s Set[T]) Add(x T) { s[x] = struct{}{} }

// Has reports whether x is in the set.
func (s Set[T]) Has(x T) bool {
	_, ok := s[x]
	return ok
}

// letterCounts counts 'a'..'z' in s. A [26]int array is comparable, so it can
// be used directly as a map key or compared with ==. Non-lowercase bytes are
// ignored; callers document that inputs are lowercase ASCII.
func letterCounts(s string) [26]int {
	var c [26]int
	for i := 0; i < len(s); i++ {
		if ch := s[i]; ch >= 'a' && ch <= 'z' {
			c[ch-'a']++
		}
	}
	return c
}
