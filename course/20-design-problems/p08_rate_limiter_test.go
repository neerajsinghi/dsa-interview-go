package design

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock is a manually advanced, goroutine-safe clock.
type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

func TestTokenBucket(t *testing.T) {
	clk := &fakeClock{t: time.Unix(0, 0)}
	b := NewTokenBucket(3, 1, clk.Now)
	for i, want := range []bool{true, true, true, false} {
		if got := b.Allow(); got != want {
			t.Fatalf("burst call %d = %v, want %v", i, got, want)
		}
	}
	clk.Advance(500 * time.Millisecond)
	if b.Allow() {
		t.Fatal("half a token should not be enough")
	}
	clk.Advance(500 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("one token should have refilled")
	}
	clk.Advance(10 * time.Second)
	if got := b.Tokens(); got != 3 {
		t.Fatalf("tokens = %v, want capped at 3", got)
	}
	if b.AllowN(5) {
		t.Fatal("AllowN above capacity must fail")
	}
	if !b.AllowN(3) || b.Allow() {
		t.Fatal("AllowN(3) should drain the bucket exactly")
	}
	clk.Advance(-time.Hour) // clock skew backwards must not add or remove tokens
	if b.Allow() {
		t.Fatal("backwards clock should not refill")
	}
}

func TestTokenBucketRealClock(t *testing.T) {
	b := NewTokenBucket(2, 1, nil)
	if !b.Allow() || !b.Allow() {
		t.Fatal("fresh bucket should allow its burst")
	}
}

// With no refill, exactly `capacity` of many concurrent requests succeed.
// Run with -race to check the locking.
func TestTokenBucketConcurrent(t *testing.T) {
	clk := &fakeClock{t: time.Unix(0, 0)}
	b := NewTokenBucket(100, 0, clk.Now)
	var allowed atomic.Int64
	var wg sync.WaitGroup
	for g := 0; g < 50; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				if b.Allow() {
					allowed.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	if got := allowed.Load(); got != 100 {
		t.Fatalf("allowed %d requests, want exactly 100", got)
	}
}
