package greedy

import (
	"math/rand"
	"testing"
)

func TestCanCompleteCircuit(t *testing.T) {
	tests := []struct {
		gas, cost []int
		want      int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{3, 4, 5, 1, 2}, 3},
		{[]int{2, 3, 4}, []int{3, 4, 3}, -1},
		{[]int{5}, []int{4}, 0},
		{[]int{4}, []int{5}, -1},
		{[]int{3, 1, 1}, []int{1, 2, 2}, 0},
	}
	for _, tt := range tests {
		if got := CanCompleteCircuit(tt.gas, tt.cost); got != tt.want {
			t.Errorf("CanCompleteCircuit(%v,%v) = %d, want %d", tt.gas, tt.cost, got, tt.want)
		}
		if got := CanCompleteCircuitBrute(tt.gas, tt.cost); got != tt.want {
			t.Errorf("Brute(%v,%v) = %d, want %d", tt.gas, tt.cost, got, tt.want)
		}
	}
}

func TestCanCompleteCircuitRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for iter := 0; iter < 500; iter++ {
		n := 1 + r.Intn(8)
		gas, cost := make([]int, n), make([]int, n)
		for i := range gas {
			gas[i], cost[i] = r.Intn(6), r.Intn(6)
		}
		greedy := CanCompleteCircuit(gas, cost)
		brute := CanCompleteCircuitBrute(gas, cost)
		// With ties (total == 0 and several zero-diff stations) more than one
		// start can work; the greedy must return SOME valid start.
		if (greedy == -1) != (brute == -1) {
			t.Fatalf("gas=%v cost=%v: greedy=%d brute=%d", gas, cost, greedy, brute)
		}
		if greedy != -1 && CanCompleteCircuitBrute(rotate(gas, greedy), rotate(cost, greedy)) != 0 {
			t.Fatalf("gas=%v cost=%v: greedy start %d does not complete a lap", gas, cost, greedy)
		}
	}
}

func rotate(xs []int, k int) []int {
	return append(append([]int{}, xs[k:]...), xs[:k]...)
}
