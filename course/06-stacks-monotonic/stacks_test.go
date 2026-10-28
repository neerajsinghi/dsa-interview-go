package stacks

import "testing"

func TestStack(t *testing.T) {
	var s Stack[string]
	if !s.IsEmpty() || s.Len() != 0 {
		t.Fatal("zero value should be empty")
	}
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop on empty stack should report ok=false")
	}
	if _, ok := s.Peek(); ok {
		t.Fatal("Peek on empty stack should report ok=false")
	}
	for _, v := range []string{"a", "b", "c"} {
		s.Push(v)
	}
	if top, _ := s.Peek(); top != "c" || s.Len() != 3 {
		t.Fatalf("Peek = %q, Len = %d; want c, 3", top, s.Len())
	}
	for _, want := range []string{"c", "b", "a"} {
		got, ok := s.Pop()
		if !ok || got != want {
			t.Fatalf("Pop = %q, %v; want %q, true", got, ok, want)
		}
	}
	if !s.IsEmpty() {
		t.Fatal("stack should be empty after popping everything")
	}
}
