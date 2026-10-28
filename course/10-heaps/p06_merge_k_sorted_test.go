package heaps

import (
	"math/rand"
	"slices"
	"testing"
)

func listFrom(vals []int) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	for _, v := range vals {
		tail.Next = &ListNode{Val: v}
		tail = tail.Next
	}
	return dummy.Next
}

func listValues(head *ListNode) []int {
	var out []int
	for ; head != nil; head = head.Next {
		out = append(out, head.Val)
	}
	return out
}

func buildLists(spec [][]int) []*ListNode {
	lists := make([]*ListNode, len(spec))
	for i, s := range spec {
		lists[i] = listFrom(s)
	}
	return lists
}

func TestMergeKLists(t *testing.T) {
	tests := []struct {
		spec [][]int
		want []int
	}{
		{[][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}}, []int{1, 1, 2, 3, 4, 4, 5, 6}},
		{nil, nil},
		{[][]int{nil}, nil},
		{[][]int{nil, {1}, nil}, []int{1}},
		{[][]int{{5}, {1}, {3}, {2}, {4}}, []int{1, 2, 3, 4, 5}},
	}
	impls := map[string]func([]*ListNode) *ListNode{"heap": MergeKListsHeap, "pairwise": MergeKListsPairwise}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := listValues(fn(buildLists(tt.spec))); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.spec, got, tt.want)
			}
		}
	}
}

func TestMergeKListsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	for iter := 0; iter < 300; iter++ {
		spec := make([][]int, r.Intn(8))
		var all []int
		for i := range spec {
			spec[i] = randomInts(r, r.Intn(10), 50)
			slices.Sort(spec[i])
			all = append(all, spec[i]...)
		}
		slices.Sort(all) // brute force: gather everything and sort
		if got := listValues(MergeKListsHeap(buildLists(spec))); !slices.Equal(got, all) {
			t.Fatalf("heap(%v) = %v, want %v", spec, got, all)
		}
		if got := listValues(MergeKListsPairwise(buildLists(spec))); !slices.Equal(got, all) {
			t.Fatalf("pairwise(%v) = %v, want %v", spec, got, all)
		}
	}
}
