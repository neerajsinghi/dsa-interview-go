// Package complexity contains small, deliberately contrasting algorithms whose
// tests prove they agree and whose benchmarks make their growth rates visible.
//
// Run the benchmarks with:
//
//	go test -run=^$ -bench=. -benchmem ./course/01-complexity/
package complexity

// SortedInts returns [0, 2, 4, ..., 2(n-1)]: n distinct even numbers in
// ascending order. Searching for an odd number is a guaranteed miss, which
// forces the worst case for search benchmarks.
func SortedInts(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = 2 * i
	}
	return out
}

// DistinctInts returns n distinct values in a scrambled (but deterministic)
// order. With no duplicates, duplicate checks must inspect everything — the
// worst case.
func DistinctInts(n int) []int {
	out := make([]int, n)
	for i := range out {
		// Multiplying by a large odd constant is a bijection on uint64,
		// so the values stay distinct while looking random.
		out[i] = int(uint64(i) * 0x9E3779B97F4A7C15 >> 1)
	}
	return out
}
