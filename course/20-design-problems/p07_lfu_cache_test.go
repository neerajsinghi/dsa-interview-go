package design

import (
	"math/rand"
	"sync"
	"testing"
)

func TestLFUCache(t *testing.T) {
	c := NewLFUCache[int, int](2)
	c.Put(1, 1)
	c.Put(2, 2)
	if v, _ := c.Get(1); v != 1 {
		t.Fatal("Get(1)")
	}
	c.Put(3, 3) // evicts 2
	if _, ok := c.Get(2); ok {
		t.Fatal("2 should be evicted")
	}
	if v, _ := c.Get(3); v != 3 {
		t.Fatal("Get(3)")
	}
	c.Put(4, 4) // 1 and 3 both freq 2; 1 used longer ago
	if _, ok := c.Get(1); ok {
		t.Fatal("1 should be evicted")
	}
	if v, _ := c.Get(3); v != 3 {
		t.Fatal("Get(3) after evict")
	}
	if v, _ := c.Get(4); v != 4 {
		t.Fatal("Get(4)")
	}
}

func TestLFUZeroCapacity(t *testing.T) {
	c := NewLFUCache[int, int](0)
	c.Put(1, 1)
	if _, ok := c.Get(1); ok {
		t.Fatal("zero-capacity cache stored a value")
	}
}

// Model: entries with explicit freq and last-used tick; evict by scan.
func TestLFURandom(t *testing.T) {
	type entry struct{ v, freq, tick int }
	r := rand.New(rand.NewSource(7))
	for _, capacity := range []int{1, 2, 3, 4} {
		c := NewLFUCache[int, int](capacity)
		model := map[int]*entry{}
		for tick := 0; tick < 3000; tick++ {
			k := r.Intn(7)
			if r.Intn(2) == 0 {
				v := r.Intn(100)
				c.Put(k, v)
				if e, ok := model[k]; ok {
					e.v, e.freq, e.tick = v, e.freq+1, tick
					continue
				}
				if len(model) == capacity {
					victim, best := -1, (*entry)(nil)
					for key, e := range model {
						if best == nil || e.freq < best.freq || (e.freq == best.freq && e.tick < best.tick) {
							victim, best = key, e
						}
					}
					delete(model, victim)
				}
				model[k] = &entry{v, 1, tick}
				continue
			}
			got, ok := c.Get(k)
			e, want := model[k]
			if ok != want || (ok && got != e.v) {
				t.Fatalf("cap %d tick %d Get(%d) = %d,%v want %v", capacity, tick, k, got, ok, want)
			}
			if ok {
				e.freq, e.tick = e.freq+1, tick
			}
		}
		if c.Len() != len(model) {
			t.Fatalf("Len = %d, model %d", c.Len(), len(model))
		}
	}
}

func TestLFUConcurrent(t *testing.T) {
	c := NewLFUCache[int, int](20)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			for i := 0; i < 1000; i++ {
				k := r.Intn(40)
				if r.Intn(2) == 0 {
					c.Put(k, -k)
				} else if v, ok := c.Get(k); ok && v != -k {
					t.Errorf("Get(%d) = %d", k, v)
				}
			}
		}(int64(g))
	}
	wg.Wait()
	if c.Len() > 20 {
		t.Fatalf("Len = %d exceeds capacity", c.Len())
	}
}
