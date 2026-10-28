package slidingwindow

import "testing"

func TestDeque(t *testing.T) {
	var d Deque
	if d.Len() != 0 {
		t.Fatal("new deque should be empty")
	}
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)
	if d.Front() != 1 || d.Back() != 3 || d.Len() != 3 {
		t.Fatalf("front=%d back=%d len=%d", d.Front(), d.Back(), d.Len())
	}
	if d.PopFront() != 1 || d.PopBack() != 3 {
		t.Fatal("unexpected pop values")
	}
	if d.Len() != 1 || d.Front() != 2 || d.Back() != 2 {
		t.Fatal("expected single element 2")
	}
	d.PopBack()
	if d.Len() != 0 {
		t.Fatal("expected empty deque")
	}
}
