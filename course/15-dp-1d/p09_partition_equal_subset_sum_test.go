package dp1d

import (
	"math/rand"
	"testing"
)

var partitionImpls = map[string]func([]int) bool{
	"recursive": CanPartitionRecursive, "memo": CanPartitionMemo,
	"tab": CanPartitionTab, "optimized": CanPartition,
}

func TestCanPartition(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{"classic", []int{1, 5, 11, 5}, true},
		{"odd total", []int{1, 2, 3, 5}, false},
		{"even total, no split", []int{1, 2, 5}, false},
		{"pair", []int{3, 3}, true},
		{"single", []int{2}, false},
		{"empty: two empty halves", nil, true},
		{"reuse trap", []int{1, 3}, false}, // iterating s upward would use 1 twice to reach 2
		{"many small", []int{1, 1, 1, 1, 1, 1, 1, 1}, true},
	}
	for _, tt := range tests {
		for name, fn := range partitionImpls {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s/%s = %v, want %v", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestCanPartitionRandom(t *testing.T) {
	r := rand.New(rand.NewSource(78))
	for iter := 0; iter < 300; iter++ {
		nums := make([]int, r.Intn(14))
		for i := range nums {
			nums[i] = 1 + r.Intn(12)
		}
		want := CanPartitionRecursive(nums)
		for name, fn := range partitionImpls {
			if got := fn(nums); got != want {
				t.Fatalf("%v: %s=%v brute=%v", nums, name, got, want)
			}
		}
	}
}
