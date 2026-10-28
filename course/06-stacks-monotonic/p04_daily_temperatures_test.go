package stacks

import (
	"math/rand"
	"slices"
	"testing"
)

func TestDailyTemperatures(t *testing.T) {
	tests := []struct {
		in, want []int
	}{
		{[]int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{[]int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
		{[]int{30, 20, 10}, []int{0, 0, 0}},
		{[]int{50, 50, 50}, []int{0, 0, 0}},
		{[]int{70}, []int{0}},
		{[]int{}, []int{}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]int) []int{
			"brute": DailyTemperaturesBrute, "stack": DailyTemperatures,
		} {
			if got := fn(tt.in); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.in, got, tt.want)
			}
		}
	}
}

func TestDailyTemperaturesRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for iter := 0; iter < 500; iter++ {
		temps := make([]int, r.Intn(30))
		for i := range temps {
			temps[i] = 30 + r.Intn(10) // narrow range forces many ties
		}
		if got, want := DailyTemperatures(temps), DailyTemperaturesBrute(temps); !slices.Equal(got, want) {
			t.Fatalf("DailyTemperatures(%v) = %v, want %v", temps, got, want)
		}
	}
}
