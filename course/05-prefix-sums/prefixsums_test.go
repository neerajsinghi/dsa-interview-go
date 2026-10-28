package prefixsums

import (
	"slices"
	"testing"
)

func TestBuildPrefix(t *testing.T) {
	tests := []struct{ nums, want []int }{
		{nil, []int{0}},
		{[]int{5}, []int{0, 5}},
		{[]int{3, 1, 4, 1}, []int{0, 3, 4, 8, 9}},
		{[]int{-1, -2}, []int{0, -1, -3}},
	}
	for _, tc := range tests {
		if got := BuildPrefix(tc.nums); !slices.Equal(got, tc.want) {
			t.Errorf("BuildPrefix(%v) = %v, want %v", tc.nums, got, tc.want)
		}
	}
	p := BuildPrefix([]int{3, 1, 4, 1})
	if got := RangeSum(p, 1, 2); got != 5 {
		t.Errorf("RangeSum(1, 2) = %d, want 5", got)
	}
}

func TestMod(t *testing.T) {
	tests := []struct{ x, k, want int }{
		{7, 6, 1}, {-7, 6, 5}, {0, 3, 0}, {-6, 6, 0}, {5, 1, 0},
	}
	for _, tc := range tests {
		if got := mod(tc.x, tc.k); got != tc.want {
			t.Errorf("mod(%d, %d) = %d, want %d", tc.x, tc.k, got, tc.want)
		}
	}
	if -7%6 != -1 {
		t.Error("documenting Go semantics: -7 % 6 should be -1")
	}
}
