package stacks

// Problem: Valid Parentheses
// Difficulty: Easy
//
// You get a string made only of the bracket characters ( ) [ ] { }. Decide
// whether it is balanced: every opener is closed by the same kind of bracket,
// and pairs close in the reverse order they were opened (properly nested).
// Any other character makes the string invalid.
//
// Examples:
//
//	"([]{})"  -> true
//	"([)]"    -> false  (the ] closes a ( )
//	"(("      -> false  (openers left over)
//
// Reference: https://leetcode.com/problems/valid-parentheses/
//
// Pattern: matching with a stack. The most recently opened bracket is the only
// one allowed to close next, which is exactly what the top of a stack holds.
//
// Time: O(n). Space: O(n) for the stack.

// IsValid reports whether s is a balanced bracket string.
func IsValid(s string) bool {
	if len(s)%2 == 1 {
		return false // an odd length can never pair up
	}
	opener := map[byte]byte{')': '(', ']': '[', '}': '{'}
	var st Stack[byte]
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '(', '[', '{':
			st.Push(c)
		case ')', ']', '}':
			top, ok := st.Pop()
			if !ok || top != opener[c] {
				return false
			}
		default:
			return false
		}
	}
	return st.IsEmpty() // leftover openers mean unbalanced
}
