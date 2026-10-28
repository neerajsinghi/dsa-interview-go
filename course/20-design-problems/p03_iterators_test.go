package design

import (
	"math/rand"
	"slices"
	"testing"
)

func TestPeekingIterator(t *testing.T) {
	p := NewPeekingIterator[int](NewSliceIterator([]int{1, 2, 3}))
	if p.Next() != 1 || p.Peek() != 2 || p.Peek() != 2 || p.Next() != 2 || p.Next() != 3 {
		t.Fatal("wrong sequence")
	}
	if p.HasNext() {
		t.Fatal("should be exhausted")
	}
	empty := NewPeekingIterator[string](NewSliceIterator[string](nil))
	if empty.HasNext() {
		t.Fatal("empty iterator has next")
	}
}

func TestNestedIterator(t *testing.T) {
	tests := []struct {
		name  string
		items []NestedInteger
		want  []int
	}{
		{"pairs", []NestedInteger{List(Int(1), Int(1)), Int(2), List(Int(1), Int(1))}, []int{1, 1, 2, 1, 1}},
		{"deep", []NestedInteger{Int(1), List(Int(4), List(Int(6)))}, []int{1, 4, 6}},
		{"empties", []NestedInteger{List(), List(List()), Int(3), List(List(), Int(4))}, []int{3, 4}},
		{"all empty", []NestedInteger{List(List(List()))}, nil},
		{"nil", nil, nil},
	}
	for _, tt := range tests {
		var got []int
		it := NewNestedIterator(tt.items)
		for it.HasNext() {
			if !it.HasNext() { // HasNext must be idempotent
				t.Fatal("HasNext changed its answer")
			}
			got = append(got, it.Next())
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("%s: iterator = %v, want %v", tt.name, got, tt.want)
		}
		if seq := slices.Collect(All(tt.items)); !slices.Equal(seq, tt.want) {
			t.Errorf("%s: All = %v, want %v", tt.name, seq, tt.want)
		}
	}
}

func TestAllEarlyBreak(t *testing.T) {
	items := []NestedInteger{Int(1), List(Int(2), Int(3)), Int(4)}
	var got []int
	for v := range All(items) {
		if v == 3 {
			break
		}
		got = append(got, v)
	}
	if !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func randNested(r *rand.Rand, depth int, next *int) []NestedInteger {
	var out []NestedInteger
	for k := r.Intn(4); k > 0; k-- {
		if depth > 0 && r.Intn(2) == 0 {
			out = append(out, List(randNested(r, depth-1, next)...))
		} else {
			out = append(out, Int(*next))
			*next++
		}
	}
	return out
}

func TestNestedIteratorRandom(t *testing.T) {
	r := rand.New(rand.NewSource(3))
	for iter := 0; iter < 200; iter++ {
		n := 0
		items := randNested(r, 4, &n)
		var got []int
		for it := NewNestedIterator(items); it.HasNext(); {
			got = append(got, it.Next())
		}
		// Integers were numbered in DFS order, so the flattening is 0..n-1.
		for i, v := range got {
			if v != i {
				t.Fatalf("got %v, want 0..%d", got, n-1)
			}
		}
		if len(got) != n || !slices.Equal(got, slices.Collect(All(items))) {
			t.Fatalf("got %v (len %d), want %d items", got, len(got), n)
		}
	}
}
