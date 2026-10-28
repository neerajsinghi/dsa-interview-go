package design

import (
	"math/rand"
	"slices"
	"testing"
)

func TestSnapshotArray(t *testing.T) {
	a := NewSnapshotArray(3)
	a.Set(0, 5)
	if id := a.Snap(); id != 0 {
		t.Fatalf("first snap id = %d", id)
	}
	a.Set(0, 6)
	if got := a.Get(0, 0); got != 5 {
		t.Fatalf("Get(0,0) = %d, want 5", got)
	}
	if got := a.Get(1, 0); got != 0 {
		t.Fatalf("unset index = %d, want 0", got)
	}
	a.Set(2, 1)
	a.Set(2, 2) // overwrite within the same version
	a.Snap()
	if a.Get(2, 0) != 0 || a.Get(2, 1) != 2 || a.Get(0, 1) != 6 {
		t.Fatal("wrong values in snapshot 1")
	}
}

// Model: store a full copy of the array at every Snap.
func TestSnapshotArrayRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	const n = 6
	a := NewSnapshotArray(n)
	cur := make([]int, n)
	var snaps [][]int
	for op := 0; op < 3000; op++ {
		switch r.Intn(3) {
		case 0:
			i, v := r.Intn(n), r.Intn(100)
			a.Set(i, v)
			cur[i] = v
		case 1:
			if id := a.Snap(); id != len(snaps) {
				t.Fatalf("snap id = %d, want %d", id, len(snaps))
			}
			snaps = append(snaps, slices.Clone(cur))
		default:
			if len(snaps) == 0 {
				continue
			}
			i, s := r.Intn(n), r.Intn(len(snaps))
			if got := a.Get(i, s); got != snaps[s][i] {
				t.Fatalf("Get(%d,%d) = %d, want %d", i, s, got, snaps[s][i])
			}
		}
	}
}
