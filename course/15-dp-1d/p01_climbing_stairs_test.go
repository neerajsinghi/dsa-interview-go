package dp1d

import "testing"

func TestClimbStairs(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 1}, {1, 1}, {2, 2}, {3, 3}, {4, 5}, {5, 8}, {10, 89}, {20, 10946},
	}
	impls := map[string]func(int) int{
		"recursive": ClimbStairsRecursive, "memo": ClimbStairsMemo,
		"tab": ClimbStairsTab, "optimized": ClimbStairs,
	}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.n); got != tt.want {
				t.Errorf("%s(%d) = %d, want %d", name, tt.n, got, tt.want)
			}
		}
	}
	// The DP versions must agree well beyond brute-force range.
	for n := 0; n <= 80; n++ {
		a, b, c := ClimbStairsMemo(n), ClimbStairsTab(n), ClimbStairs(n)
		if a != b || b != c {
			t.Fatalf("n=%d: memo=%d tab=%d optimized=%d", n, a, b, c)
		}
	}
}
