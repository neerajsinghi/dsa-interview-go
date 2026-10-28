package backtracking

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestLetterCombinations(t *testing.T) {
	tests := []struct {
		digits string
		want   []string
	}{
		{"23", []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"}},
		{"", nil},
		{"2", []string{"a", "b", "c"}},
		{"79", nil}, // checked by count below
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(string) []string{"recursive": LetterCombinations, "iterative": LetterCombinationsIterative} {
			got := fn(tt.digits)
			if tt.digits == "79" {
				if len(got) != 16 {
					t.Errorf("%s(%q) has %d results, want 16", name, tt.digits, len(got))
				}
				continue
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s(%q) = %v, want %v", name, tt.digits, got, tt.want)
			}
		}
	}
}

func TestLetterCombinationsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(61))
	for iter := 0; iter < 100; iter++ {
		b := make([]byte, r.Intn(5))
		for i := range b {
			b[i] = byte('2' + r.Intn(8))
		}
		d := string(b)
		// Both produce keypad order, so the outputs must match exactly.
		if a, c := LetterCombinations(d), LetterCombinationsIterative(d); !reflect.DeepEqual(a, c) {
			t.Fatalf("%q: recursive=%v iterative=%v", d, a, c)
		}
	}
}
