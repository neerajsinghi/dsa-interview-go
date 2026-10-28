package design

import "sync"

// Problem: LFU Cache
// Difficulty: Hard
//
// A fixed-capacity cache that evicts the LEAST FREQUENTLY used key. Every
// Get hit and every Put (insert or update) counts as one use. Among keys
// tied for the lowest use count, evict the least recently used one. Get and
// Put must run in O(1).
//
// Example (capacity 2):
//
//	Put(1,1) Put(2,2) Get(1)        freq: 1→2, 2→1
//	Put(3,3)                        evicts 2 (lowest freq)
//	Get(3)                          freq: 1→2, 3→2
//	Put(4,4)                        tie at freq 2: evicts 1 (older use)
//
// Reference: https://leetcode.com/problems/lfu-cache/
//
// Pattern: two maps and a minimum.
//
//	items: key  → node (value, freq)
//	freqs: freq → LRU-ordered list of nodes with that freq
//	minFreq: smallest freq that currently has a non-empty list
//
//	freq 1: [3]            ← minFreq; evict from the back
//	freq 2: [1 ⇄ 4]
//
// Touching a node moves it from list f to the front of list f+1; if list
// f became empty and f was minFreq, minFreq becomes f+1. A new key always
// has freq 1, so inserting sets minFreq = 1. No search is ever needed.
//
// Time: O(1) Get/Put. Space: O(capacity).
//
// Safe for concurrent use (sync.Mutex; Get mutates frequencies).

// LFUCache is a thread-safe least-frequently-used cache.
type LFUCache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	items    map[K]*node[K, V]
	freqs    map[int]*dlist[K, V]
	minFreq  int
}

// NewLFUCache returns a cache holding at most capacity entries.
func NewLFUCache[K comparable, V any](capacity int) *LFUCache[K, V] {
	return &LFUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*node[K, V]),
		freqs:    make(map[int]*dlist[K, V]),
	}
}

// touch records one more use of n.
func (c *LFUCache[K, V]) touch(n *node[K, V]) {
	old := c.freqs[n.freq]
	old.remove(n)
	if old.size == 0 {
		delete(c.freqs, n.freq) // keep the map from growing with dead buckets
		if c.minFreq == n.freq {
			c.minFreq++ // n is now the only candidate at freq+1
		}
	}
	n.freq++
	c.list(n.freq).pushFront(n)
}

// list returns (creating if needed) the list for freq f.
func (c *LFUCache[K, V]) list(f int) *dlist[K, V] {
	l, ok := c.freqs[f]
	if !ok {
		l = newDList[K, V]()
		c.freqs[f] = l
	}
	return l
}

// Get returns the value for key and counts a use on a hit.
func (c *LFUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.touch(n)
	return n.val, true
}

// Put inserts or updates key, evicting the LFU (then LRU) entry if full.
func (c *LFUCache[K, V]) Put(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity <= 0 {
		return
	}
	if n, ok := c.items[key]; ok {
		n.val = val
		c.touch(n)
		return
	}
	if len(c.items) == c.capacity {
		l := c.freqs[c.minFreq]
		victim := l.back() // least recent among the least frequent
		l.remove(victim)
		if l.size == 0 {
			delete(c.freqs, c.minFreq)
		}
		delete(c.items, victim.key)
	}
	n := &node[K, V]{key: key, val: val, freq: 1}
	c.list(1).pushFront(n)
	c.items[key] = n
	c.minFreq = 1
}

// Len returns the number of cached entries.
func (c *LFUCache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
