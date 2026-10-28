package stacks

import "testing"

func TestEvalRPN(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []string
		want    int
		wantErr bool
	}{
		{"single number", []string{"42"}, 42, false},
		{"add then multiply", []string{"2", "1", "+", "3", "*"}, 9, false},
		{"division", []string{"4", "13", "5", "/", "+"}, 6, false},
		{"negative truncation", []string{"-7", "2", "/"}, -3, false},
		{"operand order", []string{"3", "10", "-"}, -7, false},
		{"long", []string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"}, 22, false},
		{"empty", nil, 0, true},
		{"missing operand", []string{"1", "+"}, 0, true},
		{"leftover operands", []string{"1", "2"}, 0, true},
		{"bad token", []string{"1", "x", "+"}, 0, true},
		{"divide by zero", []string{"1", "0", "/"}, 0, true},
	}
	for _, tt := range tests {
		got, err := EvalRPN(tt.tokens)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: err = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("%s: EvalRPN = %d, want %d", tt.name, got, tt.want)
		}
	}
}
