// Package greedy contains solutions for the greedy algorithms and interval
// problems module.
//
// A greedy algorithm commits to the locally best choice and never revisits
// it. It is only correct when you can argue — usually with an exchange
// argument — that some optimal solution agrees with that choice.
package greedy

import (
	"cmp"
	"slices"
)

// Interval is a closed range [Start, End] unless a problem says otherwise.
type Interval struct {
	Start, End int
}

// sortedByStart returns a sorted copy so callers' slices are not mutated.
func sortedByStart(in []Interval) []Interval {
	out := slices.Clone(in)
	slices.SortFunc(out, func(a, b Interval) int {
		return cmp.Or(cmp.Compare(a.Start, b.Start), cmp.Compare(a.End, b.End))
	})
	return out
}

// sortedByEnd returns a copy sorted by end point — the classic order for
// "keep as many non-overlapping intervals as possible".
func sortedByEnd(in []Interval) []Interval {
	out := slices.Clone(in)
	slices.SortFunc(out, func(a, b Interval) int {
		return cmp.Or(cmp.Compare(a.End, b.End), cmp.Compare(a.Start, b.Start))
	})
	return out
}
