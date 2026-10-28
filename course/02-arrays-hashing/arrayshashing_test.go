package arrayshashing

import "testing"

func TestCounter(t *testing.T) {
	c := Counter([]string{"a", "b", "a"})
	if c["a"] != 2 || c["b"] != 1 || c["z"] != 0 || len(c) != 2 {
		t.Fatalf("unexpected counts %v", c)
	}
	if len(Counter[int](nil)) != 0 {
		t.Fatal("expected empty counter")
	}
}

func TestSet(t *testing.T) {
	s := NewSet([]int{1, 2, 2})
	if !s.Has(1) || !s.Has(2) || s.Has(3) || len(s) != 2 {
		t.Fatalf("unexpected set %v", s)
	}
	s.Add(3)
	if !s.Has(3) {
		t.Fatal("Add failed")
	}
}
