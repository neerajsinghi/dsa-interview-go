package linkedlists

import (
	"math/rand"
	"slices"
	"testing"
)

func TestFromToSlice(t *testing.T) {
	for _, vals := range [][]int{nil, {1}, {1, 2, 3}, {5, 5, 5, 5}} {
		head := FromSlice(vals)
		if got := ToSlice(head); !slices.Equal(got, vals) {
			t.Errorf("round trip %v = %v", vals, got)
		}
		if Len(head) != len(vals) {
			t.Errorf("Len(%v) = %d", vals, Len(head))
		}
	}
	if FromSlice(nil) != nil {
		t.Error("FromSlice(nil) should be the nil list")
	}
}

// randomSlice is shared by the randomized tests in this package.
func randomSlice(r *rand.Rand, maxLen, maxVal int) []int {
	xs := make([]int, r.Intn(maxLen+1))
	for i := range xs {
		xs[i] = r.Intn(maxVal)
	}
	return xs
}
