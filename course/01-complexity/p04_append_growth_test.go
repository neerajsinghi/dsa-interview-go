package complexity

import (
	"fmt"
	"slices"
	"testing"
)

func TestAppendVariants(t *testing.T) {
	for _, n := range []int{0, 1, 5, 100} {
		want := make([]int, n)
		for i := range want {
			want[i] = i
		}
		if got := AppendN(n); !slices.Equal(got, want) {
			t.Errorf("AppendN(%d) = %v", n, got)
		}
		pre := AppendNPreallocated(n)
		if !slices.Equal(pre, want) {
			t.Errorf("AppendNPreallocated(%d) = %v", n, pre)
		}
		if cap(pre) != n {
			t.Errorf("AppendNPreallocated(%d) cap = %d, want %d", n, cap(pre), n)
		}
		slices.Reverse(want)
		if got := PrependN(n); !slices.Equal(got, want) {
			t.Errorf("PrependN(%d) = %v, want %v", n, got, want)
		}
	}
}

// The number of regrowths is logarithmic in n, not linear: that is what makes
// append O(1) amortized. The exact count depends on the runtime's growth
// policy, so we assert generous bounds rather than exact values.
func TestCountGrowthsIsLogarithmic(t *testing.T) {
	if got := CountGrowths(0); got != 0 {
		t.Errorf("CountGrowths(0) = %d, want 0", got)
	}
	if got := CountGrowths(1); got != 1 {
		t.Errorf("CountGrowths(1) = %d, want 1", got)
	}
	for _, n := range []int{1_000, 50_000, 200_000} {
		got := CountGrowths(n)
		// Doubling would give ~log2(n); the ~1.25× growth for big slices
		// adds more steps but stays far below n.
		if got < 5 || got > 60 {
			t.Errorf("CountGrowths(%d) = %d, expected a small logarithmic count", n, got)
		}
		t.Logf("n=%d growths=%d", n, got)
	}
}

// Append and preallocated append scale linearly; prepend quadratically.
func BenchmarkAppend(b *testing.B) {
	for _, n := range []int{1_000, 4_000, 16_000} {
		b.Run(fmt.Sprintf("Append/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				AppendN(n)
			}
		})
		b.Run(fmt.Sprintf("Prealloc/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				AppendNPreallocated(n)
			}
		})
		b.Run(fmt.Sprintf("Prepend/n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				PrependN(n)
			}
		})
	}
}
