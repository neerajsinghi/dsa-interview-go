package tries

import (
	"math/rand"
	"testing"
)

func TestFindMaximumXOR(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"classic", []int{3, 10, 5, 25, 2, 8}, 28},
		{"single", []int{0}, 0},
		{"empty", nil, 0},
		{"pair", []int{1, 2}, 3},
		{"duplicates", []int{7, 7, 7}, 0},
		{"top bit", []int{1<<31 - 1, 0}, 1<<31 - 1},
		{"larger", []int{14, 70, 53, 83, 49, 91, 36, 80, 92, 51, 66, 70}, 127},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) int{"trie": FindMaximumXOR, "brute": FindMaximumXORBrute} {
			if got := fn(tt.nums); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestFindMaximumXORRandom(t *testing.T) {
	r := rand.New(rand.NewSource(16))
	for iter := 0; iter < 300; iter++ {
		nums := make([]int, 1+r.Intn(30))
		limit := 1 << (1 + r.Intn(31)) // vary magnitude so all bit levels get exercised
		for i := range nums {
			nums[i] = r.Intn(limit)
		}
		if got, want := FindMaximumXOR(nums), FindMaximumXORBrute(nums); got != want {
			t.Fatalf("%v: trie=%d brute=%d", nums, got, want)
		}
	}
}
