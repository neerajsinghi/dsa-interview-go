package heaps

import (
	"math/rand"
	"testing"
)

func TestMedianFinder(t *testing.T) {
	m := NewMedianFinder()
	steps := []struct {
		add  int
		want float64
	}{
		{1, 1}, {2, 1.5}, {3, 2}, {10, 2.5}, {-5, 2}, {2, 2}, {2, 2},
	}
	for _, s := range steps {
		m.AddNum(s.add)
		if got := m.FindMedian(); got != s.want {
			t.Fatalf("after AddNum(%d): FindMedian = %v, want %v", s.add, got, s.want)
		}
	}
}

func TestMedianFinderRandom(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	for iter := 0; iter < 50; iter++ {
		m := NewMedianFinder()
		var ref SortedMedianFinder
		for i := 0; i < 200; i++ {
			v := r.Intn(100) - 50
			m.AddNum(v)
			ref.AddNum(v)
			if got, want := m.FindMedian(), ref.FindMedian(); got != want {
				t.Fatalf("iter %d step %d: FindMedian = %v, want %v", iter, i, got, want)
			}
		}
	}
}
