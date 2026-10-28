package twopointers

import (
	"slices"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct{ in, want []int }{
		{nil, nil},
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{1, 2, 3}, []int{3, 2, 1}},
	}
	for _, tc := range tests {
		got := slices.Clone(tc.in)
		Reverse(got)
		if !slices.Equal(got, tc.want) {
			t.Errorf("Reverse(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
