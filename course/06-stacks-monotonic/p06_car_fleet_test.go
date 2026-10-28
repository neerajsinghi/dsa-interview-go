package stacks

import "testing"

func TestCarFleet(t *testing.T) {
	tests := []struct {
		name          string
		target        int
		position, spd []int
		want          int
	}{
		{"classic", 12, []int{10, 8, 0, 5, 3}, []int{2, 4, 1, 1, 3}, 3},
		{"single car", 10, []int{3}, []int{3}, 1},
		{"all merge", 100, []int{0, 2, 4}, []int{4, 2, 1}, 1},
		{"none merge", 10, []int{0, 4, 8}, []int{1, 2, 4}, 3},
		{"tie exactly at target", 10, []int{0, 5}, []int{2, 1}, 1},
		{"no cars", 10, nil, nil, 0},
		{"faster car behind catches up", 10, []int{6, 8}, []int{3, 1}, 1},
		{"unsorted input", 10, []int{8, 0, 4}, []int{4, 1, 2}, 3},
	}
	for _, tt := range tests {
		if got := CarFleet(tt.target, tt.position, tt.spd); got != tt.want {
			t.Errorf("%s: CarFleet = %d, want %d", tt.name, got, tt.want)
		}
	}
}
