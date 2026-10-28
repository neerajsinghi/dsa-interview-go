package slidingwindow

import (
	"math/rand"
	"testing"
)

func TestMaxSumSubarray(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		k      int
		want   int
		wantOK bool
	}{
		{"empty", nil, 1, 0, false},
		{"k zero", []int{1, 2}, 0, 0, false},
		{"k too big", []int{1, 2}, 3, 0, false},
		{"single", []int{5}, 1, 5, true},
		{"classic", []int{2, 1, 5, 1, 3, 2}, 3, 9, true},
		{"all negative", []int{-1, -2, -3}, 2, -3, true},
		{"whole slice", []int{1, 2, 3}, 3, 6, true},
		{"large values", []int{1 << 40, 1 << 40, 1}, 2, 1 << 41, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for name, fn := range map[string]func([]int, int) (int, bool){
				"Brute": MaxSumSubarrayBrute, "Window": MaxSumSubarray,
			} {
				got, ok := fn(tc.nums, tc.k)
				if got != tc.want || ok != tc.wantOK {
					t.Errorf("%s = (%d, %v), want (%d, %v)", name, got, ok, tc.want, tc.wantOK)
				}
			}
		})
	}
}

func TestMaxSumSubarrayRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for range 1000 {
		nums := make([]int, rng.Intn(15))
		for i := range nums {
			nums[i] = rng.Intn(21) - 10
		}
		k := rng.Intn(len(nums)+2) - 1
		a, aok := MaxSumSubarrayBrute(nums, k)
		b, bok := MaxSumSubarray(nums, k)
		if a != b || aok != bok {
			t.Fatalf("%v k=%d: brute=(%d,%v) window=(%d,%v)", nums, k, a, aok, b, bok)
		}
	}
}
