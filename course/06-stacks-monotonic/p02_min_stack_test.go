package stacks

import (
	"math/rand"
	"slices"
	"testing"
)

func TestMinStackExample(t *testing.T) {
	s := NewMinStack()
	s.Push(5)
	s.Push(2)
	s.Push(7)
	if got := s.GetMin(); got != 2 {
		t.Fatalf("GetMin = %d, want 2", got)
	}
	s.Pop()
	if got := s.Top(); got != 2 {
		t.Fatalf("Top = %d, want 2", got)
	}
	s.Pop()
	if got := s.GetMin(); got != 5 {
		t.Fatalf("GetMin = %d, want 5", got)
	}
	s.Push(5) // duplicates of the minimum must survive a pop
	s.Pop()
	if got := s.GetMin(); got != 5 {
		t.Fatalf("GetMin after duplicate pop = %d, want 5", got)
	}
}

// TestMinStackRandom compares against a plain slice that scans for the min.
func TestMinStackRandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	s := NewMinStack()
	var ref []int
	for i := 0; i < 5000; i++ {
		if len(ref) == 0 || r.Intn(3) > 0 {
			v := r.Intn(41) - 20
			s.Push(v)
			ref = append(ref, v)
		} else {
			s.Pop()
			ref = ref[:len(ref)-1]
		}
		if s.Len() != len(ref) {
			t.Fatalf("Len = %d, want %d", s.Len(), len(ref))
		}
		if len(ref) > 0 {
			if s.Top() != ref[len(ref)-1] || s.GetMin() != slices.Min(ref) {
				t.Fatalf("step %d: Top/GetMin = %d/%d, want %d/%d",
					i, s.Top(), s.GetMin(), ref[len(ref)-1], slices.Min(ref))
			}
		}
	}
}
