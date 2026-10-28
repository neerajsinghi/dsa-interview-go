package stacks

import "testing"

func TestDecodeString(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", ""},
		{"abc", "abc"},
		{"3[a]2[bc]", "aaabcbc"},
		{"3[a2[c]]", "accaccacc"},
		{"2[ab]3[x]y", "ababxxxy"},
		{"10[z]", "zzzzzzzzzz"},
		{"ab2[c3[d]]e", "abcdddcddde"},
		{"1[1[1[q]]]", "q"},
	}
	for _, tt := range tests {
		if got := DecodeString(tt.in); got != tt.want {
			t.Errorf("DecodeString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
