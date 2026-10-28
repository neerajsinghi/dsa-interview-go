package graphs

import (
	"math/rand"
	"testing"
)

func TestCourseSchedule(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		prereqs [][]int
		want    bool
	}{
		{"simple chain", 2, [][]int{{1, 0}}, true},
		{"two-cycle", 2, [][]int{{1, 0}, {0, 1}}, false},
		{"diamond", 4, [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}, true},
		{"no prerequisites", 3, nil, true},
		{"self loop", 1, [][]int{{0, 0}}, false},
		{"cycle off to the side", 5, [][]int{{1, 0}, {3, 2}, {4, 3}, {2, 4}}, false},
	}
	for _, tt := range tests {
		if got := CanFinish(tt.n, tt.prereqs); got != tt.want {
			t.Errorf("%s: CanFinish = %v, want %v", tt.name, got, tt.want)
		}
		if got := CanFinishDFS(tt.n, tt.prereqs); got != tt.want {
			t.Errorf("%s: CanFinishDFS = %v, want %v", tt.name, got, tt.want)
		}
		order := FindOrder(tt.n, tt.prereqs)
		if tt.want && !isTopoOrder(prereqGraph(tt.n, tt.prereqs), order) {
			t.Errorf("%s: FindOrder = %v is not a valid schedule", tt.name, order)
		}
		if !tt.want && (order == nil || len(order) != 0) {
			t.Errorf("%s: FindOrder = %#v, want empty non-nil slice", tt.name, order)
		}
	}
}

func TestCourseScheduleRandom(t *testing.T) {
	r := rand.New(rand.NewSource(26))
	for iter := 0; iter < 300; iter++ {
		n := 1 + r.Intn(7)
		var prereqs [][]int
		for k := r.Intn(9); k > 0; k-- {
			prereqs = append(prereqs, []int{r.Intn(n), r.Intn(n)})
		}
		a, b := CanFinish(n, prereqs), CanFinishDFS(n, prereqs)
		if a != b {
			t.Fatalf("n=%d %v: kahn=%v dfs=%v", n, prereqs, a, b)
		}
		if a && !isTopoOrder(prereqGraph(n, prereqs), FindOrder(n, prereqs)) {
			t.Fatalf("n=%d %v: invalid order", n, prereqs)
		}
	}
}
