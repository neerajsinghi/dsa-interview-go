package stacks

import (
	"errors"
	"fmt"
	"strconv"
)

// Problem: Evaluate Reverse Polish Notation
// Difficulty: Medium
//
// An arithmetic expression is given in postfix form as a list of tokens: each
// token is an integer or one of + - * /. An operator applies to the two values
// produced just before it. Division truncates toward zero. Return the value,
// or an error if the expression is malformed (bad token, missing operand,
// leftover operands, or division by zero).
//
// Examples:
//
//	["2","1","+","3","*"]        -> 9   ((2 + 1) * 3)
//	["4","13","5","/","+"]       -> 6   (4 + 13/5)
//	["-7","2","/"]               -> -3  (truncates toward zero)
//
// Reference: https://leetcode.com/problems/evaluate-reverse-polish-notation/
//
// Pattern: operand stack. Numbers are pushed; an operator pops its two
// operands (right first!), computes, and pushes the result.
//
// Time: O(n). Space: O(n).

// ErrMalformed is returned for expressions that do not reduce to one value.
var ErrMalformed = errors.New("malformed RPN expression")

// EvalRPN evaluates a postfix expression.
func EvalRPN(tokens []string) (int, error) {
	var st Stack[int]
	for _, tok := range tokens {
		switch tok {
		case "+", "-", "*", "/":
			b, okB := st.Pop() // right operand is on top
			a, okA := st.Pop()
			if !okA || !okB {
				return 0, fmt.Errorf("%w: operator %q needs two operands", ErrMalformed, tok)
			}
			switch tok {
			case "+":
				st.Push(a + b)
			case "-":
				st.Push(a - b)
			case "*":
				st.Push(a * b)
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				st.Push(a / b) // Go's integer division already truncates toward zero
			}
		default:
			n, err := strconv.Atoi(tok)
			if err != nil {
				return 0, fmt.Errorf("%w: bad token %q", ErrMalformed, tok)
			}
			st.Push(n)
		}
	}
	if st.Len() != 1 {
		return 0, fmt.Errorf("%w: %d values left on the stack", ErrMalformed, st.Len())
	}
	v, _ := st.Pop()
	return v, nil
}
