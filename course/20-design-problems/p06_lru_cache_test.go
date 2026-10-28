package design

import (
	"math/rand"
	"slices"
	"sync"
	"testing"
)

func TestLRUCache(t *testing.T) {
	c := NewLRUCache[int, int](2)
	c.Put(1, 1)
	c.Put(2, 2)
	if v, ok := c.Get(1); !ok || v != 1 {
		t.Fatal("Get(1)")
	}
	c.Put(3, 3) // evicts 2
	if _, ok := c.Get(2); ok {
		t.Fatal("2 should be evicted")
	}
	c.Put(4, 4) // evicts 1
	if _, ok := c.Get(1); ok {
		t.Fatal("1 should be evicted")
	}
	if v, _ := c.Get(3); v != 3 {
		t.Fatal("Get(3)")
	}
	if v, _ := c.Get(4); v != 4 {
		t.Fatal("Get(4)")
	}
	c.Put(3, 30) // update refreshes recency
	c.Put(5, 5)  // evicts 4
	if _, ok := c.Get(4); ok {
		t.Fatal("4 should be evicted")
	}
	if v, _ := c.Get(3); v != 30 || c.Len() != 2 {
		t.Fatal("update lost")
	}
}

func TestLRUZeroCapacity(t *testing.T) {
	c := NewLRUCache[string, int](0)
	c.Put("a", 1)
	if _, ok := c.Get("a"); ok || c.Len() != 0 {
		t.Fatal("zero-capacity cache stored a value")
	}
}

// Model: a slice ordered most → least recent, searched linearly.
func TestLRURandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	type kv struct{ k, v int }
	for _, capacity := range []int{1, 2, 3, 5} {
		c := NewLRUCache[int, int](capacity)
		var model []kv
		find := func(k int) int {
			return slices.IndexFunc(model, func(e kv) bool { return e.k == k })
		}
		for op := 0; op < 2000; op++ {
			k := r.Intn(8)
			if r.Intn(2) == 0 {
				v := r.Intn(100)
				c.Put(k, v)
				if i := find(k); i >= 0 {
					model = slices.Delete(model, i, i+1)
				} else if len(model) == capacity {
					model = model[:len(model)-1]
				}
				model = slices.Insert(model, 0, kv{k, v})
				continue
			}
			got, ok := c.Get(k)
			i := find(k)
			if ok != (i >= 0) || (ok && got != model[i].v) {
				t.Fatalf("cap %d Get(%d) = %d,%v; model %v", capacity, k, got, ok, model)
			}
			if ok {
				e := model[i]
				model = slices.Insert(slices.Delete(model, i, i+1), 0, e)
			}
		}
	}
}

func TestLRUConcurrent(t *testing.T) {
	c := NewLRUCache[int, int](50)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			for i := 0; i < 1000; i++ {
				k := r.Intn(100)
				if r.Intn(2) == 0 {
					c.Put(k, k*10)
				} else if v, ok := c.Get(k); ok && v != k*10 {
					t.Errorf("Get(%d) = %d", k, v)
				}
			}
		}(int64(g))
	}
	wg.Wait()
	if c.Len() > 50 {
		t.Fatalf("Len = %d exceeds capacity", c.Len())
	}
}
