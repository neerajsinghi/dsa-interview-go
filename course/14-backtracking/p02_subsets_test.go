package backtracking

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestSubsets(t *testing.T) {
	tests := []struct {
		nums []int
		want [][]int
	}{
		{[]int{1, 2, 3}, [][]int{{}, {1}, {2}, {3}, {1, 2}, {1, 3}, {2, 3}, {1, 2, 3}}},
		{[]int{0}, [][]int{{}, {0}}},
		{nil, [][]int{{}}},
	}
	for _, tt := range tests {
		want := sortEach(tt.want)
		for name, fn := range map[string]func([]int) [][]int{"backtrack": Subsets, "bitmask": SubsetsBitmask} {
			if got := sortEach(fn(tt.nums)); !reflect.DeepEqual(got, want) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.nums, got, want)
			}
		}
	}
}

func TestSubsetsWithDup(t *testing.T) {
	tests := []struct {
		nums []int
		want [][]int
	}{
		{[]int{1, 2, 2}, [][]int{{}, {1}, {2}, {1, 2}, {2, 2}, {1, 2, 2}}},
		{[]int{0}, [][]int{{}, {0}}},
		{[]int{4, 4, 4}, [][]int{{}, {4}, {4, 4}, {4, 4, 4}}},
		{[]int{2, 1, 2}, [][]int{{}, {1}, {2}, {1, 2}, {2, 2}, {1, 2, 2}}}, // unsorted input
	}
	for _, tt := range tests {
		orig := append([]int(nil), tt.nums...)
		if got, want := sortEach(SubsetsWithDup(tt.nums)), sortEach(tt.want); !reflect.DeepEqual(got, want) {
			t.Errorf("SubsetsWithDup(%v) = %v, want %v", tt.nums, got, want)
		}
		if !reflect.DeepEqual(orig, tt.nums) {
			t.Errorf("input was modified: %v -> %v", orig, tt.nums)
		}
	}
}

// dedupeBrute builds all 2^n subsets by bitmask, sorts each, and keeps
// distinct ones.
func dedupeBrute(nums []int) [][]int {
	seen := map[string]bool{}
	var res [][]int
	for _, s := range SubsetsBitmask(nums) {
		sort.Ints(s)
		if key := fmt.Sprint(s); !seen[key] {
			seen[key] = true
			res = append(res, s)
		}
	}
	return res
}

func TestSubsetsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(62))
	for iter := 0; iter < 200; iter++ {
		nums := make([]int, r.Intn(8))
		for i := range nums {
			nums[i] = r.Intn(4) // small range => many duplicates
		}
		got := SubsetsWithDup(nums)
		assertNoDuplicates(t, sortEach(got))
		if a, b := sortEach(got), sortEach(dedupeBrute(nums)); !reflect.DeepEqual(a, b) {
			t.Fatalf("%v: backtrack=%v brute=%v", nums, a, b)
		}
		perm := r.Perm(r.Intn(8)) // distinct values for part I
		if a, b := sortEach(Subsets(perm)), sortEach(SubsetsBitmask(perm)); !reflect.DeepEqual(a, b) {
			t.Fatalf("%v: backtrack=%v bitmask=%v", perm, a, b)
		}
	}
}
