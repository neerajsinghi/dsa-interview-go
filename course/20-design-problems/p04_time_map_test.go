package design

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func TestTimeMap(t *testing.T) {
	m := NewTimeMap()
	m.Set("foo", "bar", 1)
	checks := []struct {
		key  string
		ts   int
		want string
	}{{"foo", 1, "bar"}, {"foo", 3, "bar"}, {"foo", 0, ""}, {"missing", 5, ""}}
	for _, c := range checks {
		if got := m.Get(c.key, c.ts); got != c.want {
			t.Errorf("Get(%q,%d) = %q, want %q", c.key, c.ts, got, c.want)
		}
	}
	m.Set("foo", "bar2", 4)
	if m.Get("foo", 4) != "bar2" || m.Get("foo", 5) != "bar2" || m.Get("foo", 3) != "bar" {
		t.Fatal("wrong value after second set")
	}
	m.Set("foo", "early", 2) // out of order
	m.Set("foo", "BAR", 1)   // overwrite
	if m.Get("foo", 2) != "early" || m.Get("foo", 1) != "BAR" || m.Get("foo", 3) != "early" {
		t.Fatal("out-of-order insert or overwrite failed")
	}
}

func TestTimeMapRandom(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	m := NewTimeMap()
	model := map[string]map[int]string{}
	for op := 0; op < 3000; op++ {
		key := fmt.Sprint("k", r.Intn(3))
		ts := r.Intn(50)
		if r.Intn(2) == 0 {
			v := fmt.Sprint(op)
			m.Set(key, v, ts)
			if model[key] == nil {
				model[key] = map[int]string{}
			}
			model[key][ts] = v
			continue
		}
		want, bestTS := "", -1
		for t2, v := range model[key] { // linear scan for the floor
			if t2 <= ts && t2 > bestTS {
				want, bestTS = v, t2
			}
		}
		if got := m.Get(key, ts); got != want {
			t.Fatalf("Get(%q,%d) = %q, want %q", key, ts, got, want)
		}
	}
}

func TestTimeMapConcurrent(t *testing.T) {
	m := NewTimeMap()
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				m.Set("shared", fmt.Sprint(g), g*1000+i)
				_ = m.Get("shared", i)
			}
		}(g)
	}
	wg.Wait()
	if got := m.Get("shared", 7*1000+199); got != "7" {
		t.Fatalf("latest = %q, want 7", got)
	}
}
