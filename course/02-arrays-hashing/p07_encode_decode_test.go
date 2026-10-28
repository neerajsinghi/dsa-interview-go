package arrayshashing

import (
	"errors"
	"math/rand"
	"slices"
	"strings"
	"testing"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		strs []string
	}{
		{"empty list", []string{}},
		{"one empty string", []string{""}},
		{"several empties", []string{"", "", ""}},
		{"words", []string{"lint", "code", "love", "you"}},
		{"contains delimiter", []string{"#", "##", "a#b", "3#abc"}},
		{"digits", []string{"12", "0", "007"}},
		{"unicode", []string{"héllo", "世界", "🙂"}},
		{"long", []string{strings.Repeat("x", 1234)}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			enc := Encode(tc.strs)
			got, err := Decode(enc)
			if err != nil {
				t.Fatalf("Decode(%q) error: %v", enc, err)
			}
			if !slices.Equal(got, tc.strs) {
				t.Fatalf("round trip = %q, want %q (encoded %q)", got, tc.strs, enc)
			}
		})
	}
}

func TestEncodeFormat(t *testing.T) {
	if got := Encode([]string{"", "#"}); got != "0#1##" {
		t.Fatalf("Encode = %q", got)
	}
}

func TestDecodeMalformed(t *testing.T) {
	for _, in := range []string{"5#abc", "abc", "#abc", "x#", "3", "-1#"} {
		if _, err := Decode(in); !errors.Is(err, ErrMalformed) {
			t.Errorf("Decode(%q) err = %v, want ErrMalformed", in, err)
		}
	}
}

func TestEncodeDecodeRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	alphabet := []byte("ab#0123")
	for range 500 {
		strs := make([]string, rng.Intn(6))
		for i := range strs {
			b := make([]byte, rng.Intn(5))
			for j := range b {
				b[j] = alphabet[rng.Intn(len(alphabet))]
			}
			strs[i] = string(b)
		}
		got, err := Decode(Encode(strs))
		if err != nil || !slices.Equal(got, strs) {
			t.Fatalf("round trip of %q gave %q, %v", strs, got, err)
		}
	}
}
