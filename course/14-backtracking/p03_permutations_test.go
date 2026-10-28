package backtracking

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func TestPermute(t *testing.T) {
	tests := []struct {
		nums []int
		want [][]int
	}{
		{[]int{1, 2, 3}, [][]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}}},
		{[]int{0, 1}, [][]int{{0, 1}, {1, 0}}},
		{[]int{1}, [][]int{{1}}},
		{nil, [][]int{{}}}, // one way to arrange nothing
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) [][]int{"used": Permute, "swap": PermuteSwap, "unique": PermuteUnique} {
			if got := normalize(fn(tt.nums)); !reflect.DeepEqual(got, normalize(tt.want)) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.nums, got, tt.want)
			}
		}
	}
	if got := len(Permute([]int{1, 2, 3, 4, 5, 6})); got != 720 {
		t.Errorf("6 distinct numbers: %d permutations, want 720", got)
	}
}

func TestPermuteUnique(t *testing.T) {
	tests := []struct {
		nums []int
		want [][]int
	}{
		{[]int{1, 1, 2}, [][]int{{1, 1, 2}, {1, 2, 1}, {2, 1, 1}}},
		{[]int{2, 2, 2}, [][]int{{2, 2, 2}}},
		{[]int{2, 1, 1}, [][]int{{1, 1, 2}, {1, 2, 1}, {2, 1, 1}}},
	}
	for _, tt := range tests {
		if got := normalize(PermuteUnique(tt.nums)); !reflect.DeepEqual(got, normalize(tt.want)) {
			t.Errorf("PermuteUnique(%v) = %v, want %v", tt.nums, got, tt.want)
		}
	}
}

func TestPermuteRandom(t *testing.T) {
	r := rand.New(rand.NewSource(63))
	for iter := 0; iter < 100; iter++ {
		nums := make([]int, r.Intn(6))
		for i := range nums {
			nums[i] = r.Intn(3)
		}
		// Brute force for part II: all n! orderings, deduplicated.
		seen := map[string]bool{}
		var want [][]int
		for _, p := range PermuteSwap(nums) {
			if k := fmt.Sprint(p); !seen[k] {
				seen[k] = true
				want = append(want, p)
			}
		}
		got := PermuteUnique(nums)
		assertNoDuplicates(t, got)
		if !reflect.DeepEqual(normalize(got), normalize(want)) {
			t.Fatalf("%v: unique=%v brute=%v", nums, got, want)
		}
		if len(PermuteSwap(nums)) != factorial(len(nums)) {
			t.Fatalf("%v: swap version produced wrong count", nums)
		}
	}
}
