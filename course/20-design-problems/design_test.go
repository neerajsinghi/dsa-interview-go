package design

import "testing"

func TestDList(t *testing.T) {
	l := newDList[int, string]()
	if l.back() != nil {
		t.Fatal("empty list should have no back")
	}
	a, b, c := &node[int, string]{key: 1}, &node[int, string]{key: 2}, &node[int, string]{key: 3}
	l.pushFront(a)
	l.pushFront(b)
	l.pushFront(c) // c b a
	if l.back() != a || l.size != 3 {
		t.Fatal("back should be the first pushed")
	}
	l.moveToFront(a) // a c b
	if l.back() != b {
		t.Fatal("moveToFront broke order")
	}
	l.remove(c) // a b
	if l.size != 2 || l.root.next != a || a.next != b || b.next != &l.root {
		t.Fatal("remove broke links")
	}
}
