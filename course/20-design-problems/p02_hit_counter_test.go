package design

import (
	"math/rand"
	"sync"
	"testing"
)

func TestHitCounter(t *testing.T) {
	var h HitCounter
	h.Hit(1)
	h.Hit(2)
	h.Hit(3)
	if got := h.GetHits(4); got != 3 {
		t.Fatalf("GetHits(4) = %d, want 3", got)
	}
	h.Hit(300)
	if got := h.GetHits(300); got != 4 {
		t.Fatalf("GetHits(300) = %d, want 4", got)
	}
	if got := h.GetHits(301); got != 3 {
		t.Fatalf("GetHits(301) = %d, want 3", got)
	}
	if got := h.GetHits(10_000); got != 0 {
		t.Fatalf("GetHits(10000) = %d, want 0", got)
	}
}

func TestHitCounterRandom(t *testing.T) {
	r := rand.New(rand.NewSource(2))
	var h HitCounter
	var log []int // every hit timestamp: the naive model
	ts := 1
	for op := 0; op < 3000; op++ {
		ts += r.Intn(40)
		if r.Intn(3) > 0 {
			h.Hit(ts)
			log = append(log, ts)
			continue
		}
		want := 0
		for _, x := range log {
			if ts-x < HitWindow {
				want++
			}
		}
		if got := h.GetHits(ts); got != want {
			t.Fatalf("GetHits(%d) = %d, want %d", ts, got, want)
		}
	}
}

func TestHitCounterConcurrent(t *testing.T) {
	var h HitCounter
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				h.Hit(42)
				_ = h.GetHits(42)
			}
		}()
	}
	wg.Wait()
	if got := h.GetHits(42); got != 8*500 {
		t.Fatalf("GetHits = %d, want %d", got, 8*500)
	}
}
