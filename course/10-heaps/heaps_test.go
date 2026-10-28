package heaps

import (
	"math/rand"
	"slices"
	"testing"
)

func randomInts(r *rand.Rand, n, maxVal int) []int {
	xs := make([]int, n)
	for i := range xs {
		xs[i] = r.Intn(maxVal)
	}
	return xs
}

func TestHeapSortsLikeSlicesSort(t *testing.T) {
	r := rand.New(rand.NewSource(10))
	for iter := 0; iter < 200; iter++ {
		xs := randomInts(r, r.Intn(50), 30)
		want := slices.Sorted(slices.Values(xs))

		minH := NewMin[int]()
		arr := &ArrayHeap[int]{}
		for _, v := range xs {
			minH.Push(v)
			arr.Push(v)
		}
		built := FromSlice(slices.Clone(xs), func(a, b int) bool { return a < b })
		heapified := Heapify(slices.Clone(xs))

		var got1, got2, got3, got4 []int
		for minH.Len() > 0 {
			got1 = append(got1, minH.Pop())
		}
		for built.Len() > 0 {
			got2 = append(got2, built.Pop())
		}
		for arr.Len() > 0 {
			v, _ := arr.Pop()
			got3 = append(got3, v)
		}
		for heapified.Len() > 0 {
			v, _ := heapified.Pop()
			got4 = append(got4, v)
		}
		for name, got := range map[string][]int{"Heap": got1, "FromSlice": got2, "ArrayHeap": got3, "Heapify": got4} {
			if !slices.Equal(got, want) {
				t.Fatalf("%s order = %v, want %v", name, got, want)
			}
		}

		maxH := NewMax[int]()
		for _, v := range xs {
			maxH.Push(v)
		}
		slices.Reverse(want)
		for i := 0; maxH.Len() > 0; i++ {
			if v := maxH.Pop(); v != want[i] {
				t.Fatalf("max-heap pop %d = %d, want %d", i, v, want[i])
			}
		}
	}
}

func TestHeapPeekAndEmpty(t *testing.T) {
	h := NewMin[string]()
	h.Push("pear")
	h.Push("apple")
	h.Push("fig")
	if h.Peek() != "apple" || h.Len() != 3 {
		t.Fatalf("Peek = %q, Len = %d", h.Peek(), h.Len())
	}
	var a ArrayHeap[int]
	if _, ok := a.Pop(); ok {
		t.Fatal("Pop on empty ArrayHeap should report ok=false")
	}
	if _, ok := a.Peek(); ok {
		t.Fatal("Peek on empty ArrayHeap should report ok=false")
	}
	a.Push(3)
	a.Push(1)
	if v, _ := a.Peek(); v != 1 {
		t.Fatalf("ArrayHeap Peek = %d, want 1", v)
	}
}
