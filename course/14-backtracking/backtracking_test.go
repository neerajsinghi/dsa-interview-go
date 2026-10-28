package backtracking

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// normalize sorts each inner slice's position in the outer list so results
// whose order is unspecified can be compared. Inner slices are left as-is
// (sort them first with sortEach when their order doesn't matter either).
func normalize(sets [][]int) [][]int {
	out := make([][]int, len(sets))
	for i, s := range sets {
		out[i] = clone(s)
		if out[i] == nil {
			out[i] = []int{} // treat nil and empty alike
		}
	}
	sort.Slice(out, func(i, j int) bool { return lessInts(out[i], out[j]) })
	return out
}

// sortEach sorts every inner slice in place, then normalizes the outer list.
func sortEach(sets [][]int) [][]int {
	out := normalize(sets)
	for _, s := range out {
		sort.Ints(s)
	}
	return normalize(out)
}

func lessInts(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

// assertNoDuplicates fails if any inner slice appears twice.
func assertNoDuplicates(t *testing.T, sets [][]int) {
	t.Helper()
	seen := map[string]bool{}
	for _, s := range sets {
		key := fmt.Sprint(s)
		if seen[key] {
			t.Fatalf("duplicate %v in %v", s, sets)
		}
		seen[key] = true
	}
}

func binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func TestCombine(t *testing.T) {
	want := [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}
	if got := Combine(4, 2); !reflect.DeepEqual(got, want) {
		t.Errorf("Combine(4,2) = %v, want %v", got, want)
	}
	if got := Combine(3, 0); !reflect.DeepEqual(got, [][]int{{}}) {
		t.Errorf("Combine(3,0) = %v, want [[]]", got)
	}
	for n := 0; n <= 7; n++ {
		for k := 0; k <= n; k++ {
			got := Combine(n, k)
			if len(got) != binom(n, k) {
				t.Fatalf("Combine(%d,%d) has %d results, want %d", n, k, len(got), binom(n, k))
			}
			assertNoDuplicates(t, got)
		}
	}
}
