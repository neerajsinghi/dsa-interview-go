package design

import (
	"slices"
	"sync"
)

// Problem: Time-Based Key-Value Store
// Difficulty: Medium
//
// Store several values per key, each stamped with a timestamp. Get(key, ts)
// returns the value whose timestamp is the largest one ≤ ts, or "" if the
// key has no value that old. Setting the same key and timestamp again
// overwrites it.
//
// Example:
//
//	Set("foo","bar",1)
//	Get("foo",1) -> "bar"   Get("foo",3) -> "bar"
//	Set("foo","bar2",4)
//	Get("foo",4) -> "bar2"  Get("foo",5) -> "bar2"  Get("foo",0) -> ""
//
// Reference: https://leetcode.com/problems/time-based-key-value-store/
//
// Pattern: map of key → slice of (timestamp, value) kept sorted by time;
// Get is a binary search for the last timestamp ≤ ts ("floor"). The usual
// version assumes timestamps arrive increasing so Set is an append; this one
// also accepts out-of-order writes by inserting at the search position.
//
//	"foo": [(1,bar) (4,bar2)]     Get(3): first ts > 3 is index 1 → answer index 0
//
// Time: Set O(log n) search + O(1) amortized append (O(n) for an
// out-of-order insert). Get O(log n). Space: O(total entries).
//
// Safe for concurrent use. Get only reads, so a sync.RWMutex lets readers
// proceed in parallel.

type timedValue struct {
	ts    int
	value string
}

// TimeMap is a versioned key-value store.
type TimeMap struct {
	mu   sync.RWMutex
	data map[string][]timedValue
}

// NewTimeMap returns an empty store.
func NewTimeMap() *TimeMap {
	return &TimeMap{data: make(map[string][]timedValue)}
}

func cmpTS(e timedValue, ts int) int {
	switch {
	case e.ts < ts:
		return -1
	case e.ts > ts:
		return 1
	}
	return 0
}

// Set stores value for key at timestamp ts.
func (m *TimeMap) Set(key, value string, ts int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	entries := m.data[key]
	// Fast path: the common case of increasing timestamps.
	if n := len(entries); n == 0 || entries[n-1].ts < ts {
		m.data[key] = append(entries, timedValue{ts, value})
		return
	}
	i, found := slices.BinarySearchFunc(entries, ts, cmpTS)
	if found {
		entries[i].value = value
		return
	}
	m.data[key] = slices.Insert(entries, i, timedValue{ts, value})
}

// Get returns the value for key with the greatest timestamp ≤ ts.
func (m *TimeMap) Get(key string, ts int) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entries := m.data[key]
	// i = first index with entries[i].ts >= ts (or len).
	i, found := slices.BinarySearchFunc(entries, ts, cmpTS)
	if found {
		return entries[i].value
	}
	if i == 0 {
		return "" // every stored timestamp is newer than ts
	}
	return entries[i-1].value
}
