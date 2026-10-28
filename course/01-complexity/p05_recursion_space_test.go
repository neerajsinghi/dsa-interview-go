package complexity

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"empty", nil, 0},
		{"single", []int{42}, 42},
		{"several", []int{1, 2, 3, 4}, 10},
		{"negatives", []int{-5, 2, -1}, -4},
		{"odd length", []int{1, 1, 1, 1, 1}, 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for name, fn := range map[string]func([]int) int{
				"Recursive": SumRecursive, "Iterative": SumIterative, "Halves": SumHalves,
			} {
				if got := fn(tc.nums); got != tc.want {
					t.Errorf("%s(%v) = %d, want %d", name, tc.nums, got, tc.want)
				}
			}
		})
	}
}

func TestSumRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(5))
	for range 200 {
		nums := make([]int, rng.Intn(200))
		for i := range nums {
			nums[i] = rng.Intn(2001) - 1000
		}
		want := SumIterative(nums)
		if SumRecursive(nums) != want || SumHalves(nums) != want {
			t.Fatalf("sums disagree on len %d", len(nums))
		}
	}
}

// A deep recursion works in Go (stacks grow), but costs O(n) memory.
func TestSumRecursiveDeep(t *testing.T) {
	nums := make([]int, 20_000)
	for i := range nums {
		nums[i] = 1
	}
	if got := SumRecursive(nums); got != len(nums) {
		t.Fatalf("got %d", got)
	}
}

// All three are O(n) time; the recursive ones pay for call overhead and,
// for SumRecursive, for growing the goroutine stack.
func BenchmarkSum(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		nums := make([]int, n)
		b.Run(fmt.Sprintf("Iterative/n=%d", n), func(b *testing.B) {
			for range b.N {
				SumIterative(nums)
			}
		})
		b.Run(fmt.Sprintf("Recursive/n=%d", n), func(b *testing.B) {
			for range b.N {
				SumRecursive(nums)
			}
		})
		b.Run(fmt.Sprintf("Halves/n=%d", n), func(b *testing.B) {
			for range b.N {
				SumHalves(nums)
			}
		})
	}
}
