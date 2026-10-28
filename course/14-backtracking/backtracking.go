// Package backtracking contains the choose / explore / unchoose template and
// the classic exhaustive-search interview problems: subsets, permutations,
// combinations, partitions, grid search, N-Queens and Sudoku.
//
// Every backtracking solution has the same skeleton:
//
//	func backtrack(state) {
//		if state is a complete answer {
//			record a COPY of state
//			return
//		}
//		for each choice available from state {
//			if choice is invalid { continue }   // prune
//			apply choice                        // choose
//			backtrack(state)                    // explore
//			undo choice                         // unchoose
//		}
//	}
//
// The recursion tree enumerates candidates; pruning cuts whole subtrees the
// moment a partial candidate cannot lead to a valid answer.
package backtracking

// Combine returns all k-element combinations of 1..n, each in increasing
// order. It is the template in its purest form and the base of many
// problems ("choose k of n").
//
//	n=4, k=2:              []
//	                 /    |    \
//	              [1]    [2]    [3]      (4 cannot start: no room for a 2nd)
//	             / | \    | \     \
//	        [1 2][1 3][1 4][2 3][2 4][3 4]
//
// Time: O(k · C(n, k)). Space: O(k) recursion, plus the output.
func Combine(n, k int) [][]int {
	var res [][]int
	path := make([]int, 0, k)
	var backtrack func(start int)
	backtrack = func(start int) {
		if len(path) == k {
			res = append(res, clone(path)) // path keeps changing: store a copy
			return
		}
		// Prune: we still need k-len(path) numbers, so start can go no higher
		// than n-(k-len(path))+1.
		for x := start; x <= n-(k-len(path))+1; x++ {
			path = append(path, x)    // choose
			backtrack(x + 1)          // explore
			path = path[:len(path)-1] // unchoose
		}
	}
	backtrack(1)
	return res
}

// clone returns an independent copy of s. Appending a slice that the
// recursion keeps mutating is the #1 Go backtracking bug: every stored
// answer would alias the same backing array.
// It always returns a non-nil slice, so an empty answer prints as [] and
// compares equal to []int{}.
func clone(s []int) []int {
	out := make([]int, len(s))
	copy(out, s)
	return out
}
