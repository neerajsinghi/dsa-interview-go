package backtracking

// Problem: Letter Combinations of a Phone Number
// Difficulty: Medium
//
// On an old phone keypad, digits 2-9 each map to a few letters (2 → abc,
// 3 → def, …, 7 → pqrs, 9 → wxyz). Given a string of such digits, return
// every string you could type by picking one letter per digit. An empty
// input yields no combinations.
//
// Examples:
//
//	"23" -> [ad ae af bd be bf cd ce cf]
//	""   -> []
//	"2"  -> [a b c]
//
// Reference: https://leetcode.com/problems/letter-combinations-of-a-phone-number/
//
// Pattern: backtracking over positions. Depth i of the recursion picks the
// letter for digit i; every leaf at depth len(digits) is one answer. Since
// every branch is valid there is no pruning — this is the purest form of
// "enumerate the Cartesian product".
//
//	           ""
//	     /     |     \
//	    a      b      c        digit '2'
//	  / | \  / | \  / | \
//	 ad ae af ...              digit '3'
//
// Time: O(4^n · n) — up to 4 letters per digit, n to build each string.
// Space: O(n) recursion besides the output.

var keypad = [10]string{"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"}

// LetterCombinations enumerates with recursive backtracking, in keypad order.
func LetterCombinations(digits string) []string {
	if digits == "" {
		return nil
	}
	var res []string
	buf := make([]byte, 0, len(digits))
	var backtrack func(i int)
	backtrack = func(i int) {
		if i == len(digits) {
			res = append(res, string(buf)) // string() copies the bytes
			return
		}
		for _, ch := range []byte(keypad[digits[i]-'0']) {
			buf = append(buf, ch)  // choose
			backtrack(i + 1)       // explore
			buf = buf[:len(buf)-1] // unchoose
		}
	}
	backtrack(0)
	return res
}

// LetterCombinationsIterative builds the answers level by level: start with
// [""] and, for each digit, extend every partial string with each letter.
func LetterCombinationsIterative(digits string) []string {
	if digits == "" {
		return nil
	}
	res := []string{""}
	for i := 0; i < len(digits); i++ {
		var next []string
		for _, prefix := range res {
			for _, ch := range keypad[digits[i]-'0'] {
				next = append(next, prefix+string(ch))
			}
		}
		res = next
	}
	return res
}
