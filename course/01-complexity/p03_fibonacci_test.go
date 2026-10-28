package complexity

import (
	"fmt"
	"testing"
)

func TestFib(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0}, {1, 1}, {2, 1}, {3, 2}, {10, 55}, {20, 6765},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprint(tc.n), func(t *testing.T) {
			if got := FibNaive(tc.n); got != tc.want {
				t.Errorf("FibNaive(%d) = %d, want %d", tc.n, got, tc.want)
			}
			if got := FibMemo(tc.n); got != tc.want {
				t.Errorf("FibMemo(%d) = %d, want %d", tc.n, got, tc.want)
			}
			if got := FibIter(tc.n); got != tc.want {
				t.Errorf("FibIter(%d) = %d, want %d", tc.n, got, tc.want)
			}
		})
	}
}

func TestFibLarge(t *testing.T) {
	// F(92) is the largest Fibonacci number that fits in int64.
	const f92 = 7540113804746346429
	if got := FibIter(92); got != f92 {
		t.Errorf("FibIter(92) = %d, want %d", got, f92)
	}
	if got := FibMemo(92); got != f92 {
		t.Errorf("FibMemo(92) = %d, want %d", got, f92)
	}
}

func TestFibAgree(t *testing.T) {
	for n := range 25 {
		if a, b, c := FibNaive(n), FibMemo(n), FibIter(n); a != b || b != c {
			t.Fatalf("n=%d: naive=%d memo=%d iter=%d", n, a, b, c)
		}
	}
}

// Naive grows ~11× per +5 (φ⁵ ≈ 11.1); memo and iter grow linearly.
func BenchmarkFib(b *testing.B) {
	for _, n := range []int{10, 15, 20, 25} {
		b.Run(fmt.Sprintf("Naive/n=%d", n), func(b *testing.B) {
			for range b.N {
				FibNaive(n)
			}
		})
		b.Run(fmt.Sprintf("Memo/n=%d", n), func(b *testing.B) {
			for range b.N {
				FibMemo(n)
			}
		})
		b.Run(fmt.Sprintf("Iter/n=%d", n), func(b *testing.B) {
			for range b.N {
				FibIter(n)
			}
		})
	}
}
