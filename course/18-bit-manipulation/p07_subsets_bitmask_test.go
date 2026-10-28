package bits

import (
	"fmt"
	"slices"
	"testing"
)

func TestSubsets(t *testing.T) {
	got := Subsets([]int{1, 2, 3})
	want := [][]int{{}, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}}
	if !slices.EqualFunc(got, want, slices.Equal[[]int]) {
		t.Fatalf("Subsets([1 2 3]) = %v, want %v", got, want)
	}
	if got := Subsets(nil); len(got) != 1 || len(got[0]) != 0 {
		t.Fatalf("Subsets(nil) = %v, want [[]]", got)
	}
	// Every subset distinct, count 2^n.
	for n := 0; n <= 10; n++ {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i * 7
		}
		subs := Subsets(nums)
		seen := map[string]bool{}
		for _, s := range subs {
			seen[fmt.Sprint(s)] = true
		}
		if len(subs) != 1<<n || len(seen) != 1<<n {
			t.Fatalf("n=%d: %d subsets, %d distinct", n, len(subs), len(seen))
		}
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

func TestCombinations(t *testing.T) {
	got := Combinations([]int{1, 2, 3, 4}, 2)
	want := [][]int{{1, 2}, {1, 3}, {2, 3}, {1, 4}, {2, 4}, {3, 4}}
	if !slices.EqualFunc(got, want, slices.Equal[[]int]) {
		t.Fatalf("Combinations = %v, want %v", got, want)
	}
	for n := 0; n <= 10; n++ {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = i
		}
		for k := -1; k <= n+1; k++ {
			combos := Combinations(nums, k)
			if len(combos) != binom(n, k) {
				t.Fatalf("n=%d k=%d: %d combos, want %d", n, k, len(combos), binom(n, k))
			}
			for _, c := range combos {
				if len(c) != k {
					t.Fatalf("n=%d k=%d: combo %v has wrong size", n, k, c)
				}
			}
		}
	}
}
