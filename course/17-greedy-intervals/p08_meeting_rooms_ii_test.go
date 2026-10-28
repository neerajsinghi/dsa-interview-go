package greedy

import (
	"math/rand"
	"testing"
)

func TestMinMeetingRooms(t *testing.T) {
	tests := []struct {
		name string
		in   []Interval
		want int
	}{
		{"classic", iv(0, 30, 5, 10, 15, 20), 2},
		{"disjoint", iv(7, 10, 2, 4), 1},
		{"back to back", iv(1, 5, 5, 10, 10, 15), 1},
		{"all overlap", iv(1, 10, 2, 9, 3, 8), 3},
		{"empty", nil, 0},
		{"single", iv(3, 4), 1},
		{"same start", iv(1, 3, 1, 3, 3, 5), 2},
	}
	impls := map[string]func([]Interval) int{"brute": MinMeetingRoomsBrute, "heap": MinMeetingRoomsHeap, "sweep": MinMeetingRoomsSweep}
	for _, tt := range tests {
		for name, fn := range impls {
			if got := fn(tt.in); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestMinMeetingRoomsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(8))
	for iter := 0; iter < 300; iter++ {
		in := randIntervals(r, r.Intn(15), 30)
		for i := range in {
			in[i].End++ // meetings have positive length
		}
		want := MinMeetingRoomsBrute(in)
		if a, b := MinMeetingRoomsHeap(in), MinMeetingRoomsSweep(in); a != want || b != want {
			t.Fatalf("%v: brute=%d heap=%d sweep=%d", in, want, a, b)
		}
	}
}
