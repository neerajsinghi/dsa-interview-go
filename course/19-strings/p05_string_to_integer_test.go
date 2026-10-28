package stringalgos

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
)

func TestMyAtoi(t *testing.T) {
	tests := []struct {
		s    string
		want int32
	}{
		{"42", 42},
		{"   -042abc", -42},
		{"1337c0d3", 1337},
		{"0-1", 0},
		{"words 99", 0},
		{"", 0},
		{"   ", 0},
		{"+", 0},
		{"-", 0},
		{"+-12", 0},
		{"  +7", 7},
		{"2147483647", math.MaxInt32},
		{"2147483648", math.MaxInt32},
		{"-2147483648", math.MinInt32},
		{"-2147483649", math.MinInt32},
		{"-91283472332", math.MinInt32},
		{"99999999999999999999999", math.MaxInt32},
		{"00000000000000000000123", 123},
		{"\t5", 0}, // only ' ' counts as leading whitespace
	}
	for _, tt := range tests {
		if got := MyAtoi(tt.s); got != tt.want {
			t.Errorf("MyAtoi(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

// Compare with strconv on random digit strings (big.Int-free reference).
func TestMyAtoiRandom(t *testing.T) {
	r := rand.New(rand.NewSource(5))
	for iter := 0; iter < 2000; iter++ {
		digits := randText(r, 1+r.Intn(12), "0123456789")
		sign := []string{"", "+", "-"}[r.Intn(3)]
		want, err := strconv.ParseInt(sign+digits, 10, 32)
		if err != nil { // out of range: ParseInt returns the clamped bound
			if ne, ok := err.(*strconv.NumError); !ok || ne.Err != strconv.ErrRange {
				t.Fatalf("unexpected error %v", err)
			}
		}
		if got := MyAtoi("  " + sign + digits + "x9"); int64(got) != want {
			t.Fatalf("MyAtoi(%q) = %d, want %d", sign+digits, got, want)
		}
	}
}
