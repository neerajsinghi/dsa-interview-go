package linkedlists

import (
	"math/rand"
	"testing"
)

// buildRandomList builds nodes from (value, randomIndex) pairs; -1 = nil.
func buildRandomList(spec [][2]int) *RandomNode {
	nodes := make([]*RandomNode, len(spec))
	for i, s := range spec {
		nodes[i] = &RandomNode{Val: s[0]}
	}
	for i, s := range spec {
		if i+1 < len(nodes) {
			nodes[i].Next = nodes[i+1]
		}
		if s[1] >= 0 {
			nodes[i].Random = nodes[s[1]]
		}
	}
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

// describe turns a list back into (value, randomIndex) pairs.
func describe(head *RandomNode) [][2]int {
	index := map[*RandomNode]int{}
	i := 0
	for n := head; n != nil; n = n.Next {
		index[n] = i
		i++
	}
	var out [][2]int
	for n := head; n != nil; n = n.Next {
		r := -1
		if n.Random != nil {
			r = index[n.Random]
		}
		out = append(out, [2]int{n.Val, r})
	}
	return out
}

func equalSpec(a, b [][2]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func checkCopy(t *testing.T, name string, spec [][2]int, fn func(*RandomNode) *RandomNode) {
	t.Helper()
	orig := buildRandomList(spec)
	cp := fn(orig)
	if got := describe(cp); !equalSpec(got, spec) {
		t.Errorf("%s: copy = %v, want %v", name, got, spec)
	}
	if got := describe(orig); !equalSpec(got, spec) {
		t.Errorf("%s: original was modified: %v", name, got)
	}
	origNodes := map[*RandomNode]bool{}
	for n := orig; n != nil; n = n.Next {
		origNodes[n] = true
	}
	for n := cp; n != nil; n = n.Next {
		if origNodes[n] || (n.Random != nil && origNodes[n.Random]) {
			t.Errorf("%s: copy shares nodes with the original", name)
			return
		}
	}
}

func TestCopyRandomList(t *testing.T) {
	tests := [][][2]int{
		nil,
		{{1, -1}},
		{{1, 0}},
		{{7, -1}, {13, 0}, {11, 4}, {10, 2}, {1, 0}},
		{{1, 1}, {2, 1}},
		{{3, -1}, {3, 0}, {3, -1}},
	}
	r := rand.New(rand.NewSource(7))
	for i := 0; i < 100; i++ {
		n := r.Intn(15)
		spec := make([][2]int, n)
		for j := range spec {
			spec[j] = [2]int{r.Intn(10), r.Intn(n+1) - 1}
		}
		tests = append(tests, spec)
	}
	for _, spec := range tests {
		checkCopy(t, "map", spec, CopyRandomListMap)
		checkCopy(t, "interleave", spec, CopyRandomListInterleave)
	}
}
