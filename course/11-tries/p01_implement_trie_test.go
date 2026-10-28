package tries

import (
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestTrieBasics(t *testing.T) {
	tr := New()
	tr.Insert("apple")
	checks := []struct {
		name string
		got  bool
		want bool
	}{
		{"search apple", tr.Search("apple"), true},
		{"search app before insert", tr.Search("app"), false},
		{"startsWith app", tr.StartsWith("app"), true},
		{"startsWith empty", tr.StartsWith(""), true},
		{"startsWith b", tr.StartsWith("b"), false},
		{"search longer", tr.Search("apples"), false},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
	tr.Insert("app")
	if !tr.Search("app") {
		t.Errorf("Search(app) after insert = false")
	}
	if got := tr.CountPrefix("ap"); got != 2 {
		t.Errorf("CountPrefix(ap) = %d, want 2", got)
	}
	tr.Insert("app") // duplicate insert must not change counts
	if got := tr.Len(); got != 2 {
		t.Errorf("Len after duplicate insert = %d, want 2", got)
	}
}

func TestTrieDelete(t *testing.T) {
	tr := New()
	for _, w := range []string{"car", "card", "care", "cat"} {
		tr.Insert(w)
	}
	if !tr.Delete("car") {
		t.Fatal("Delete(car) = false, want true")
	}
	if tr.Search("car") || !tr.Search("card") || !tr.StartsWith("car") {
		t.Error("deleting a word must keep longer words that share its prefix")
	}
	if tr.Delete("car") {
		t.Error("second Delete(car) = true, want false")
	}
	if tr.Delete("ca") {
		t.Error("Delete of a mere prefix must fail")
	}
	tr.Delete("card")
	tr.Delete("care")
	if tr.StartsWith("car") {
		t.Error("branch should be pruned once no word uses it")
	}
	if got := tr.CountPrefix("c"); got != 1 {
		t.Errorf("CountPrefix(c) = %d, want 1", got)
	}
	// Empty word is a legal key.
	tr.Insert("")
	if !tr.Search("") || tr.Len() != 2 {
		t.Error("empty word should be storable")
	}
}

func TestTrieWithPrefix(t *testing.T) {
	tr := New()
	for _, w := range []string{"mouse", "mobile", "moneypot", "monitor", "mousepad"} {
		tr.Insert(w)
	}
	tests := []struct {
		prefix string
		limit  int
		want   []string
	}{
		{"mo", 3, []string{"mobile", "moneypot", "monitor"}},
		{"mou", 0, []string{"mouse", "mousepad"}},
		{"x", 3, nil},
		{"mousepad", 3, []string{"mousepad"}},
	}
	for _, tt := range tests {
		if got := tr.WithPrefix(tt.prefix, tt.limit); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("WithPrefix(%q,%d) = %v, want %v", tt.prefix, tt.limit, got, tt.want)
		}
	}
}

func TestMapTrieUnicode(t *testing.T) {
	tr := NewMapTrie()
	tr.Insert("héllo")
	tr.Insert("日本語")
	if !tr.Search("héllo") || !tr.StartsWith("日本") || tr.Search("日本") {
		t.Error("map trie should handle multi-byte runes")
	}
}

// TestTrieRandom drives Trie and MapTrie with random operations and checks
// every answer against a plain map used as the reference implementation.
func TestTrieRandom(t *testing.T) {
	r := rand.New(rand.NewSource(11))
	randWord := func() string {
		b := make([]byte, r.Intn(5))
		for i := range b {
			b[i] = byte('a' + r.Intn(3)) // tiny alphabet = lots of shared prefixes
		}
		return string(b)
	}
	tr, mt := New(), NewMapTrie()
	ref := map[string]bool{}
	for op := 0; op < 3000; op++ {
		w := randWord()
		switch r.Intn(4) {
		case 0:
			tr.Insert(w)
			mt.Insert(w)
			ref[w] = true
		case 1:
			if got, want := tr.Delete(w), ref[w]; got != want {
				t.Fatalf("Delete(%q) = %v, want %v", w, got, want)
			}
			delete(ref, w)
		case 2:
			if got := tr.Search(w); got != ref[w] {
				t.Fatalf("Search(%q) = %v, want %v", w, got, ref[w])
			}
		case 3:
			count := 0
			for k := range ref {
				if strings.HasPrefix(k, w) {
					count++
				}
			}
			if got := tr.CountPrefix(w); got != count {
				t.Fatalf("CountPrefix(%q) = %d, want %d", w, got, count)
			}
			if got := tr.StartsWith(w); got != (count > 0 || w == "") {
				t.Fatalf("StartsWith(%q) = %v, want %v", w, got, count > 0)
			}
		}
		if tr.Len() != len(ref) {
			t.Fatalf("Len = %d, want %d", tr.Len(), len(ref))
		}
	}
	// MapTrie never deletes, so compare it against everything ever inserted.
	all := tr.WithPrefix("", 0)
	if !sort.StringsAreSorted(all) || len(all) != len(ref) {
		t.Fatalf("WithPrefix(\"\") = %v, not the sorted key set", all)
	}
	for _, w := range all {
		if !mt.Search(w) {
			t.Fatalf("MapTrie lost %q", w)
		}
	}
}
