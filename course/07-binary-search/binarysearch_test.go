package binarysearch

import (
	"math/rand"
	"slices"
	"sort"
	"testing"
)

func TestBoundsAgainstStdlib(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for iter := 0; iter < 1000; iter++ {
		xs := make([]int, r.Intn(20))
		for i := range xs {
			xs[i] = r.Intn(10)
		}
		slices.Sort(xs)
		target := r.Intn(12) - 1

		wantLo, _ := slices.BinarySearch(xs, target)
		if got := LowerBound(xs, target); got != wantLo {
			t.Fatalf("LowerBound(%v, %d) = %d, want %d", xs, target, got, wantLo)
		}
		wantHi := sort.Search(len(xs), func(i int) bool { return xs[i] > target })
		if got := UpperBound(xs, target); got != wantHi {
			t.Fatalf("UpperBound(%v, %d) = %d, want %d", xs, target, got, wantHi)
		}
		if got := FirstTrue(0, len(xs), func(i int) bool { return xs[i] >= target }); got != wantLo {
			t.Fatalf("FirstTrue(%v >= %d) = %d, want %d", xs, target, got, wantLo)
		}
	}
}

func TestFirstTrueNoneTrue(t *testing.T) {
	if got := FirstTrue(3, 9, func(int) bool { return false }); got != 9 {
		t.Fatalf("FirstTrue with no true = %d, want hi (9)", got)
	}
	if got := FirstTrue(3, 9, func(int) bool { return true }); got != 3 {
		t.Fatalf("FirstTrue with all true = %d, want lo (3)", got)
	}
}

func TestBoundsStrings(t *testing.T) {
	xs := []string{"apple", "banana", "banana", "cherry"}
	if lo, hi := LowerBound(xs, "banana"), UpperBound(xs, "banana"); lo != 1 || hi != 3 {
		t.Fatalf("bounds(banana) = %d,%d; want 1,3", lo, hi)
	}
}
