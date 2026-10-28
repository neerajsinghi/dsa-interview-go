package prefixsums

import (
	"math/rand"
	"testing"
)

func TestNumArray(t *testing.T) {
	nums := []int{-2, 0, 3, -5, 2, -1}
	a := NewNumArray(nums)
	tests := []struct{ l, r, want int }{
		{0, 2, 1}, {2, 5, -1}, {0, 5, -3}, {3, 3, -5}, {0, 0, -2}, {5, 5, -1},
	}
	for _, tc := range tests {
		if got := a.SumRange(tc.l, tc.r); got != tc.want {
			t.Errorf("SumRange(%d, %d) = %d, want %d", tc.l, tc.r, got, tc.want)
		}
		if got := SumRangeBrute(nums, tc.l, tc.r); got != tc.want {
			t.Errorf("SumRangeBrute(%d, %d) = %d, want %d", tc.l, tc.r, got, tc.want)
		}
	}
}

func TestNumArrayIndependentOfInput(t *testing.T) {
	nums := []int{1, 2, 3}
	a := NewNumArray(nums)
	nums[0] = 100 // mutate after construction
	if got := a.SumRange(0, 2); got != 6 {
		t.Fatalf("SumRange = %d, want 6 (prefix should be a snapshot)", got)
	}
}

func TestNumArrayLargeValues(t *testing.T) {
	a := NewNumArray([]int{1 << 60, 1 << 60, -(1 << 60)})
	if got := a.SumRange(1, 2); got != 0 {
		t.Fatalf("SumRange(1, 2) = %d, want 0", got)
	}
}

func TestNumArrayRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for range 200 {
		nums := make([]int, 1+rng.Intn(20))
		for i := range nums {
			nums[i] = rng.Intn(201) - 100
		}
		a := NewNumArray(nums)
		for range 20 {
			l := rng.Intn(len(nums))
			r := l + rng.Intn(len(nums)-l)
			if x, y := SumRangeBrute(nums, l, r), a.SumRange(l, r); x != y {
				t.Fatalf("%v [%d,%d]: brute=%d prefix=%d", nums, l, r, x, y)
			}
		}
	}
}
