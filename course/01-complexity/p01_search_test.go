package complexity

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"empty", nil, 1, -1},
		{"single hit", []int{5}, 5, 0},
		{"single miss", []int{5}, 4, -1},
		{"first", []int{1, 3, 5, 7, 9}, 1, 0},
		{"last", []int{1, 3, 5, 7, 9}, 9, 4},
		{"middle", []int{1, 3, 5, 7, 9}, 7, 3},
		{"miss between", []int{1, 3, 5, 7, 9}, 4, -1},
		{"miss below", []int{1, 3, 5, 7, 9}, -10, -1},
		{"miss above", []int{1, 3, 5, 7, 9}, 100, -1},
		{"negatives", []int{-9, -4, -1, 0, 8}, -4, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := LinearSearch(tc.nums, tc.target); got != tc.want {
				t.Errorf("LinearSearch = %d, want %d", got, tc.want)
			}
			if got := BinarySearch(tc.nums, tc.target); got != tc.want {
				t.Errorf("BinarySearch = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestSearchRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for range 500 {
		nums := SortedInts(rng.Intn(50))
		target := rng.Intn(110) - 5
		if l, b := LinearSearch(nums, target), BinarySearch(nums, target); l != b {
			t.Fatalf("nums=%v target=%d: linear=%d binary=%d", nums, target, l, b)
		}
	}
}

// Linear time grows 10× per row; binary barely moves.
func BenchmarkSearch(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000, 1_000_000} {
		nums := SortedInts(n)
		miss := 2*n + 1 // odd, so never present: worst case for both
		b.Run(fmt.Sprintf("Linear/n=%d", n), func(b *testing.B) {
			for range b.N {
				LinearSearch(nums, miss)
			}
		})
		b.Run(fmt.Sprintf("Binary/n=%d", n), func(b *testing.B) {
			for range b.N {
				BinarySearch(nums, miss)
			}
		})
	}
}
