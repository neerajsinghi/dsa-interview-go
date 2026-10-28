package design

import "sort"

// Problem: Snapshot Array
// Difficulty: Medium
//
// An array of fixed length starts filled with zeros. Set(i, v) writes a
// value; Snap() takes a snapshot and returns its id (0, 1, 2, …);
// Get(i, snapID) returns what index i held when that snapshot was taken.
//
// Example:
//
//	a := NewSnapshotArray(3)
//	a.Set(0, 5); a.Snap() -> 0
//	a.Set(0, 6); a.Get(0, 0) -> 5
//
// Reference: https://leetcode.com/problems/snapshot-array/
//
// Pattern: copy-on-write per index. Copying the whole array on every Snap
// is O(length) per call. Instead each index keeps its own history of
// (snapID, value) records — only indices that change pay for storage —
// and Get binary-searches for the last record with snapID ≤ the query.
//
//	index 0: [(0,5) (1,6)]     Get(0, 0) -> 5;  Get(0, 7) -> 6
//
// Several Sets between two Snaps overwrite the same record.
//
// Time: Set O(1), Snap O(1), Get O(log s). Space: O(length + sets).
//
// Not safe for concurrent use.

type snapRecord struct {
	snap, val int
}

// SnapshotArray is an array with cheap point-in-time snapshots.
type SnapshotArray struct {
	history [][]snapRecord // history[i] sorted by snap
	snapID  int            // id the NEXT Snap will return
}

// NewSnapshotArray returns an array of length zeros.
func NewSnapshotArray(length int) *SnapshotArray {
	return &SnapshotArray{history: make([][]snapRecord, length)}
}

// Set writes val at index i in the current (unsnapped) version.
func (a *SnapshotArray) Set(i, val int) {
	h := a.history[i]
	if n := len(h); n > 0 && h[n-1].snap == a.snapID {
		h[n-1].val = val // already written since the last Snap: overwrite
		return
	}
	a.history[i] = append(h, snapRecord{a.snapID, val})
}

// Snap freezes the current version and returns its id.
func (a *SnapshotArray) Snap() int {
	a.snapID++
	return a.snapID - 1
}

// Get returns the value at index i as of snapshot snapID.
func (a *SnapshotArray) Get(i, snapID int) int {
	h := a.history[i]
	// First record written AFTER snapID; the one before it is the answer.
	j := sort.Search(len(h), func(k int) bool { return h[k].snap > snapID })
	if j == 0 {
		return 0 // never set before that snapshot
	}
	return h[j-1].val
}
