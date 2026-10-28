package heaps

import (
	"math/rand"
	"testing"
)

func TestLeastInterval(t *testing.T) {
	tests := []struct {
		tasks string
		n     int
		want  int
	}{
		{"AAABBB", 2, 8},
		{"AAABBB", 0, 6},
		{"AAAAAABCDEFG", 2, 16},
		{"A", 5, 1},
		{"", 3, 0},
		{"ABCDE", 4, 5},
		{"AAAB", 3, 9},
		{"ACABDB", 1, 6},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]byte, int) int{"heap": LeastIntervalHeap, "formula": LeastInterval} {
			if got := fn([]byte(tt.tasks), tt.n); got != tt.want {
				t.Errorf("%s(%q, %d) = %d, want %d", name, tt.tasks, tt.n, got, tt.want)
			}
		}
	}
}

func TestLeastIntervalRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for iter := 0; iter < 500; iter++ {
		tasks := make([]byte, r.Intn(30))
		for i := range tasks {
			tasks[i] = byte('A' + r.Intn(5))
		}
		n := r.Intn(5)
		if got, want := LeastInterval(tasks, n), LeastIntervalHeap(tasks, n); got != want {
			t.Fatalf("LeastInterval(%q, %d) = %d, simulation says %d", tasks, n, got, want)
		}
	}
}
