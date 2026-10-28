package design

import "sync"

// Problem: LRU Cache
// Difficulty: Medium
//
// A fixed-capacity key-value cache. Get returns a value and marks the key
// as recently used. Put inserts or updates a key (also marking it used);
// if that pushes the size over capacity, the LEAST recently used key is
// evicted. Both operations must run in O(1).
//
// Example (capacity 2):
//
//	Put(1,1) Put(2,2) Get(1)->1      order now: 1, 2  (most → least recent)
//	Put(3,3)                         evicts 2
//	Get(2) -> miss   Put(4,4) evicts 1   Get(1) -> miss  Get(3)->3  Get(4)->4
//
// Reference: https://leetcode.com/problems/lru-cache/
//
// Pattern: hash map + doubly linked list.
//
//	map[key] ──► node      list: root ⇄ [most recent] ⇄ … ⇄ [least recent] ⇄ root
//
// The map finds a node in O(1); the list reorders it in O(1) because a
// node knows both neighbours. Eviction pops the back. The node stores its
// key so eviction can delete the map entry.
//
// Time: O(1) Get/Put. Space: O(capacity).
//
// Safe for concurrent use. Note that Get MUTATES the recency list, so a
// sync.RWMutex with RLock in Get would be a data race — use a plain Mutex
// (or shard the cache to reduce contention).
//
// See also the cache-aside pattern:
// https://github.com/neerajsinghi/system-design-patterns-go (patterns/data/cache-aside)

// LRUCache is a thread-safe least-recently-used cache.
type LRUCache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	items    map[K]*node[K, V]
	order    *dlist[K, V] // front = most recently used
}

// NewLRUCache returns a cache holding at most capacity entries. A capacity
// ≤ 0 yields a cache that stores nothing.
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*node[K, V], max(capacity, 0)),
		order:    newDList[K, V](),
	}
}

// Get returns the value for key and whether it was present.
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.order.moveToFront(n)
	return n.val, true
}

// Put inserts or updates key, evicting the least recently used entry if full.
func (c *LRUCache[K, V]) Put(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity <= 0 {
		return
	}
	if n, ok := c.items[key]; ok {
		n.val = val
		c.order.moveToFront(n)
		return
	}
	if c.order.size == c.capacity {
		oldest := c.order.back()
		c.order.remove(oldest)
		delete(c.items, oldest.key) // why the node stores its key
	}
	n := &node[K, V]{key: key, val: val}
	c.order.pushFront(n)
	c.items[key] = n
}

// Len returns the number of cached entries.
func (c *LRUCache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.size
}
