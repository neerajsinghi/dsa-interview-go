// Package design contains solutions for the object-oriented / coding-design
// round: small data structures with a public API, invariants to protect,
// and (where noted) safe concurrent use.
//
// This file holds the intrusive doubly linked list shared by the LRU and
// LFU caches. It is written by hand rather than with container/list so the
// pointer surgery — the part interviewers actually watch — is visible, and
// so nodes are typed (no interface{} assertions).
package design

// node is one cache entry. freq is only used by the LFU cache.
type node[K comparable, V any] struct {
	key        K
	val        V
	freq       int
	prev, next *node[K, V]
}

// dlist is a circular doubly linked list with a sentinel root:
// root.next is the front (most recent), root.prev is the back (oldest).
// The sentinel removes every "is this the head/tail?" special case.
type dlist[K comparable, V any] struct {
	root node[K, V]
	size int
}

func newDList[K comparable, V any]() *dlist[K, V] {
	l := &dlist[K, V]{}
	l.root.next = &l.root
	l.root.prev = &l.root
	return l
}

// pushFront links n right after the sentinel.
func (l *dlist[K, V]) pushFront(n *node[K, V]) {
	n.prev = &l.root
	n.next = l.root.next
	l.root.next.prev = n
	l.root.next = n
	l.size++
}

// remove unlinks n. n must currently be in l.
func (l *dlist[K, V]) remove(n *node[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev, n.next = nil, nil // help the GC and catch double-removes early
	l.size--
}

// back returns the oldest node, or nil if the list is empty.
func (l *dlist[K, V]) back() *node[K, V] {
	if l.size == 0 {
		return nil
	}
	return l.root.prev
}

// moveToFront marks n as most recently used.
func (l *dlist[K, V]) moveToFront(n *node[K, V]) {
	l.remove(n)
	l.pushFront(n)
}
