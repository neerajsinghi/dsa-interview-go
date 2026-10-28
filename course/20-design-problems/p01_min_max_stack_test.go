package design

import (
	"math/rand"
	"slices"
	"testing"
)

func TestMinStack(t *testing.T) {
	var s MinStack[int]
	if _, ok := s.GetMin(); ok {
		t.Fatal("empty stack should have no min")
	}
	if _, ok := s.Pop(); ok {
		t.Fatal("empty pop should fail")
	}
	s.Push(-2)
	s.Push(0)
	s.Push(-3)
	if m, _ := s.GetMin(); m != -3 {
		t.Fatalf("min = %d, want -3", m)
	}
	s.Pop()
	if top, _ := s.Top(); top != 0 {
		t.Fatalf("top = %d, want 0", top)
	}
	if m, _ := s.GetMin(); m != -2 {
		t.Fatalf("min = %d, want -2", m)
	}
}

func TestMaxStack(t *testing.T) {
	var s MaxStack[int]
	s.Push(5)
	s.Push(1)
	s.Push(5)
	if v, _ := s.Top(); v != 5 {
		t.Fatalf("top = %d", v)
	}
	if v, _ := s.PopMax(); v != 5 {
		t.Fatalf("popMax = %d", v)
	}
	if v, _ := s.Top(); v != 1 {
		t.Fatalf("top after popMax = %d, want 1", v)
	}
	if v, _ := s.PeekMax(); v != 5 {
		t.Fatalf("peekMax = %d, want 5", v)
	}
	if v, _ := s.Pop(); v != 1 {
		t.Fatalf("pop = %d", v)
	}
	if v, _ := s.Top(); v != 5 {
		t.Fatalf("top = %d", v)
	}
	s.PopMax()
	if _, ok := s.PopMax(); ok || s.Len() != 0 {
		t.Fatal("stack should be empty")
	}
}

// Model: a plain slice where min/max are found by scanning.
func TestMinMaxStackRandom(t *testing.T) {
	r := rand.New(rand.NewSource(20))
	var mn MinStack[int]
	var mx MaxStack[int]
	var model []int
	for op := 0; op < 3000; op++ {
		switch k := r.Intn(4); {
		case k <= 1 || len(model) == 0:
			v := r.Intn(20)
			mn.Push(v)
			mx.Push(v)
			model = append(model, v)
		case k == 2:
			want := model[len(model)-1]
			model = model[:len(model)-1]
			a, _ := mn.Pop()
			b, _ := mx.Pop()
			if a != want || b != want {
				t.Fatalf("pop: min=%d max=%d want %d", a, b, want)
			}
		default:
			// PopMax on the max stack only; mirror it on the model and rebuild mn.
			want := slices.Max(model)
			i := len(model) - 1
			for model[i] != want {
				i--
			}
			model = slices.Delete(model, i, i+1)
			if got, _ := mx.PopMax(); got != want {
				t.Fatalf("popMax = %d, want %d", got, want)
			}
			mn = MinStack[int]{}
			for _, v := range model {
				mn.Push(v)
			}
		}
		if len(model) > 0 {
			gMin, _ := mn.GetMin()
			gMax, _ := mx.PeekMax()
			gTop, _ := mx.Top()
			if gMin != slices.Min(model) || gMax != slices.Max(model) || gTop != model[len(model)-1] {
				t.Fatalf("model %v: min=%d max=%d top=%d", model, gMin, gMax, gTop)
			}
		}
	}
}
