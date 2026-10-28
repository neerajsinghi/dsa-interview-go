package complexity

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestContainsDuplicate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{"empty", nil, false},
		{"single", []int{1}, false},
		{"pair equal", []int{7, 7}, true},
		{"distinct", []int{1, 2, 3, 4}, false},
		{"dup far apart", []int{1, 2, 3, 1}, true},
		{"negatives", []int{-1, -2, -1}, true},
		{"extremes", []int{math.MinInt, math.MaxInt, 0}, false},
		{"all same", []int{4, 4, 4, 4}, true},
	}
	impls := map[string]func([]int) bool{
		"Brute": ContainsDuplicateBrute,
		"Sort":  ContainsDuplicateSort,
		"Hash":  ContainsDuplicate,
	}
	for _, tc := range tests {
		for name, fn := range impls {
			t.Run(tc.name+"/"+name, func(t *testing.T) {
				if got := fn(tc.nums); got != tc.want {
					t.Errorf("%s(%v) = %v, want %v", name, tc.nums, got, tc.want)
				}
			})
		}
	}
}

func TestContainsDuplicateDoesNotMutate(t *testing.T) {
	nums := []int{3, 1, 2}
	ContainsDuplicateSort(nums)
	if nums[0] != 3 || nums[1] != 1 || nums[2] != 2 {
		t.Fatalf("input mutated: %v", nums)
	}
}

func TestContainsDuplicateRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	for range 500 {
		nums := make([]int, rng.Intn(30))
		for i := range nums {
			nums[i] = rng.Intn(60) - 30
		}
		want := ContainsDuplicateBrute(nums)
		if ContainsDuplicateSort(nums) != want || ContainsDuplicate(nums) != want {
			t.Fatalf("disagreement on %v (brute=%v)", nums, want)
		}
	}
}

// Quadratic grows ~16× per 4× n; the hash set ~4×.
func BenchmarkContainsDuplicate(b *testing.B) {
	for _, n := range []int{128, 512, 2048} {
		nums := DistinctInts(n)
		b.Run(fmt.Sprintf("Brute/n=%d", n), func(b *testing.B) {
			for range b.N {
				ContainsDuplicateBrute(nums)
			}
		})
		b.Run(fmt.Sprintf("Sort/n=%d", n), func(b *testing.B) {
			for range b.N {
				ContainsDuplicateSort(nums)
			}
		})
		b.Run(fmt.Sprintf("Hash/n=%d", n), func(b *testing.B) {
			for range b.N {
				ContainsDuplicate(nums)
			}
		})
	}
}
