package dp1d

import "testing"

func TestFibonacci(t *testing.T) {
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765}
	impls := map[string]func(int) int{
		"recursive": FibRecursive, "memo": FibMemo, "tab": FibTab, "optimized": FibOptimized,
	}
	for n, w := range want {
		for name, fn := range impls {
			if got := fn(n); got != w {
				t.Errorf("%s(%d) = %d, want %d", name, n, got, w)
			}
		}
	}
	// Only the polynomial versions can handle large n.
	const f90 = 2880067194370816120
	for name, fn := range map[string]func(int) int{"memo": FibMemo, "tab": FibTab, "optimized": FibOptimized} {
		if got := fn(90); got != f90 {
			t.Errorf("%s(90) = %d, want %d", name, got, f90)
		}
	}
}
