package stacks

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"()", true},
		{"()[]{}", true},
		{"([]{})", true},
		{"{[()()]}", true},
		{"(]", false},
		{"([)]", false},
		{"((", false},
		{"))", false},
		{"(", false},
		{"]", false},
		{"(a)", false},
		{"(((())))", true},
	}
	for _, tt := range tests {
		if got := IsValid(tt.in); got != tt.want {
			t.Errorf("IsValid(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
