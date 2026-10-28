package backtracking

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestCombinationSum(t *testing.T) {
	tests := []struct {
		name       string
		candidates []int
		target     int
		want       [][]int
	}{
		{"classic", []int{2, 3, 6, 7}, 7, [][]int{{2, 2, 3}, {7}}},
		{"three ways", []int{2, 3, 5}, 8, [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}}},
		{"impossible", []int{2}, 1, nil},
		{"unsorted input", []int{7, 3, 2, 6}, 7, [][]int{{2, 2, 3}, {7}}},
		{"single coin repeated", []int{1}, 3, [][]int{{1, 1, 1}}},
	}
	for _, tt := range tests {
		got := CombinationSum(tt.candidates, tt.target)
		if !reflect.DeepEqual(normalize(got), normalize(tt.want)) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCombinationSumRandom(t *testing.T) {
	r := rand.New(rand.NewSource(64))
	for iter := 0; iter < 200; iter++ {
		var cands []int
		for _, v := range r.Perm(8)[:1+r.Intn(4)] {
			cands = append(cands, v+1) // distinct values 1..8
		}
		target := 1 + r.Intn(15)
		got := CombinationSum(cands, target)
		if len(got) != countCombinations(cands, target) {
			t.Fatalf("%v target %d: %d answers, DP says %d", cands, target, len(got), countCombinations(cands, target))
		}
		assertNoDuplicates(t, got)
		for _, combo := range got {
			sum := 0
			for i, v := range combo {
				sum += v
				if i > 0 && combo[i-1] > v {
					t.Fatalf("%v is not non-decreasing", combo)
				}
			}
			if sum != target {
				t.Fatalf("%v sums to %d, want %d", combo, sum, target)
			}
		}
	}
}
