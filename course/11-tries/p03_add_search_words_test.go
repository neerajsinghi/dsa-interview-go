package tries

import (
	"math/rand"
	"testing"
)

func TestWordDictionary(t *testing.T) {
	d := NewWordDictionary()
	for _, w := range []string{"bad", "dad", "mad"} {
		d.AddWord(w)
	}
	tests := []struct {
		pattern string
		want    bool
	}{
		{"pad", false},
		{"bad", true},
		{".ad", true},
		{"b..", true},
		{"b.", false},
		{"...", true},
		{"....", false},
		{"", false},
		{"ba", false}, // prefix of a word is not a match
	}
	for _, tt := range tests {
		if got := d.Search(tt.pattern); got != tt.want {
			t.Errorf("Search(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
	d.AddWord("")
	if !d.Search("") {
		t.Error("Search(\"\") after adding the empty word = false")
	}
}

func TestWordDictionaryRandom(t *testing.T) {
	r := rand.New(rand.NewSource(13))
	randStr := func(alpha string) string {
		b := make([]byte, r.Intn(4))
		for i := range b {
			b[i] = alpha[r.Intn(len(alpha))]
		}
		return string(b)
	}
	for round := 0; round < 50; round++ {
		d, ref := NewWordDictionary(), &wordListDictionary{}
		for op := 0; op < 60; op++ {
			if r.Intn(3) == 0 {
				w := randStr("abc")
				d.AddWord(w)
				ref.AddWord(w)
				continue
			}
			p := randStr("abc..")
			if got, want := d.Search(p), ref.Search(p); got != want {
				t.Fatalf("Search(%q) = %v, want %v (words %q)", p, got, want, ref.words)
			}
		}
	}
}
